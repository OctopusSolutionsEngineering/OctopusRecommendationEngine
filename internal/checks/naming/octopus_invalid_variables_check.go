package naming

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
)

var linkOptions = regexp.MustCompile(`\{.*?}`)

const OctoLintInvalidVariableNames = "OctoLintInvalidVariableNames"

// OctopusInvalidVariableNameCheck checks to see if any project variables are unused.
type OctopusInvalidVariableNameCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusInvalidVariableNameCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusInvalidVariableNameCheck {
	return OctopusInvalidVariableNameCheck{
		client:       client,
		errorHandler: errorHandler,
		config:       config,
	}
}

func (o OctopusInvalidVariableNameCheck) Id() string {
	return OctoLintInvalidVariableNames
}

func (o OctopusInvalidVariableNameCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
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
		o.config.MaxInvalidVariableProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Naming, err)
	}

	regex, err := regexp.Compile(o.config.VariableNameRegex)

	if err != nil {
		return checks.NewOctopusCheckResultImpl(
			"The supplied regex "+o.config.VariableNameRegex+" does not compile",
			o.Id(),
			"",
			checks.Error,
			checks.Naming), nil
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	messages := threadsafe.NewSlice[string]()
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

			for _, v := range variableSet.Variables {
				if checks.IgnoreVariable(v.Name) {
					continue
				}

				if !regex.Match([]byte(v.Name)) {
					messages.Append(p.Name + ": " + v.Name)
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
		return o.errorHandler.HandleError(o.Id(), checks.Naming, goroutineErrors.Values()[0])
	}

	if messages.Length() > 0 {

		return checks.NewOctopusCheckResultImpl(
			"The following variables do not match the regex "+o.config.VariableNameRegex+":\n"+strings.Join(messages.Values(), "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Naming), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no unused variables",
		o.Id(),
		"",
		checks.Ok,
		checks.Naming), nil
}

func (o OctopusInvalidVariableNameCheck) getDeploymentSteps(p *projects2.Project) ([]*deployments.DeploymentStep, error) {
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
func (o OctopusInvalidVariableNameCheck) naiveStepVariableScan(deploymentSteps []*deployments.DeploymentStep, variable *variables.Variable) bool {
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
func (o OctopusInvalidVariableNameCheck) naiveVariableSetVariableScan(variables variables.VariableSet, variable *variables.Variable) bool {
	for _, v := range variables.Variables {
		if strings.Index(v.Value, variable.Name) != -1 {
			return true
		}
	}

	return false
}
