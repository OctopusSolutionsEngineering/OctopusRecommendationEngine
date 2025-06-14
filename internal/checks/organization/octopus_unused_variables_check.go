package organization

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	projects2 "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/hayageek/threadsafe"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strings"
	"sync"
)

const OctoLintUnusedVariables = "OctoLintUnusedVariables"

var linkOptions = regexp.MustCompile(`\{.*?}`)

// OctopusUnusedVariablesCheck checks to see if any project variables are unused.
type OctopusUnusedVariablesCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
	mu           sync.Mutex
}

func NewOctopusUnusedVariablesCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) *OctopusUnusedVariablesCheck {
	return &OctopusUnusedVariablesCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o *OctopusUnusedVariablesCheck) Id() string {
	return OctoLintUnusedVariables
}

func (o *OctopusUnusedVariablesCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	projects, err := client_wrapper.GetProjectsWithFilter(
		o.client,
		o.client.GetSpaceID(),
		o.config.ExcludeProjectsExcept,
		o.config.ExcludeProjects,
		o.config.MaxUnusedVariablesProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	unusedVars := map[*projects2.Project][]*variables.Variable{}
	goroutineErrors := threadsafe.NewSlice[error]()

	for i, p := range projects {
		i := i
		p := p

		g.Go(func() error {
			zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(projects))*100) + "% complete")

			variableSet, err := o.client.Variables.GetAll(p.ID)

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					goroutineErrors.Append(err)
				}
				return nil
			}

			deploymentSteps, err := o.getDeploymentSteps(p)

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					goroutineErrors.Append(err)
				}
				return nil
			}

			// Lock the map so we are not writing to it concurrently
			o.mu.Lock()
			defer o.mu.Unlock()

			for _, v := range variableSet.Variables {
				if checks.IgnoreVariable(v.Name) {
					continue
				}

				used := o.naiveStepVariableScan(deploymentSteps, v) || o.naiveVariableSetVariableScan(variableSet, v)

				if !used {
					if _, ok := unusedVars[p]; !ok {
						unusedVars[p] = []*variables.Variable{}
					}
					unusedVars[p] = append(unusedVars[p], v)
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Treat the first error as the root cause
	if goroutineErrors.Length() > 0 {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, goroutineErrors.Values()[0])
	}

	if len(unusedVars) > 0 {
		messages := []string{}
		for p, variables := range unusedVars {
			if len(variables) != 0 {
				for _, variable := range variables {
					messages = append(messages, p.Name+": "+variable.Name)
				}
			}
		}

		return checks.NewOctopusCheckResultImpl(
			"The following variables may be unused (note there are edge cases that may use these variables that can't be detected, so double check these before deleting them): \n"+strings.Join(messages, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no unused variables",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o *OctopusUnusedVariablesCheck) getDeploymentSteps(p *projects2.Project) ([]*deployments.DeploymentStep, error) {
	deploymentProcesses := []*deployments.DeploymentStep{}
	deploymentProcess, err := o.client.DeploymentProcesses.GetByID(p.DeploymentProcessID)

	if err != nil {
		if !o.errorHandler.ShouldContinue(err) {
			return nil, err
		}
	} else {
		if deploymentProcess != nil && deploymentProcess.Steps != nil {
			deploymentProcesses = append(deploymentProcesses, deploymentProcess.Steps...)
		}
	}

	if link, ok := p.Links["Runbooks"]; ok {
		runbooks, err := newclient.Get[resources.Resources[runbooks.Runbook]](o.client.HttpSession(), linkOptions.ReplaceAllString(link, ""))

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
		}

		for _, runbook := range runbooks.Items {
			runbookProcess, err := o.client.RunbookProcesses.GetByID(runbook.RunbookProcessID)

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					return nil, err
				}
				continue
			} else {
				if runbookProcess != nil && runbookProcess.Steps != nil {
					deploymentProcesses = append(deploymentProcesses, runbookProcess.Steps...)
				}
			}
		}
	}

	return deploymentProcesses, nil
}

// naiveStepVariableScan does a simple text search for the variable in a steps properties. This does lead to false positives as simple variables names, like "a",
// will almost certainly appear in a step property text without necessarily being referenced as a variable.
func (o *OctopusUnusedVariablesCheck) naiveStepVariableScan(deploymentSteps []*deployments.DeploymentStep, variable *variables.Variable) bool {
	if deploymentSteps != nil {
		for _, s := range deploymentSteps {
			for _, a := range s.Actions {
				for _, p := range a.Properties {
					if strings.Index(p.Value, variable.Name) != -1 {
						return true
					}
				}

				// Packages and feeds can use variables
				for _, p := range a.Packages {
					if strings.Index(p.FeedID, variable.Name) != -1 || strings.Index(p.PackageID, variable.Name) != -1 {
						return true
					}
				}
			}
		}
	}

	return false
}

// naiveVariableSetVariableScan does a simple text search for the variable in the value of other variables
func (o *OctopusUnusedVariablesCheck) naiveVariableSetVariableScan(variables variables.VariableSet, variable *variables.Variable) bool {
	for _, v := range variables.Variables {
		if strings.Index(v.Value, variable.Name) != -1 {
			return true
		}
	}

	return false
}
