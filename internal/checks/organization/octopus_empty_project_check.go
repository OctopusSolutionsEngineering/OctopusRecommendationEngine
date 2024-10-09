package organization

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/hayageek/threadsafe"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"strings"
)

const OctoLintEmptyProject = "OctoLintEmptyProject"

// OctopusEmptyProjectCheck checks for projects with no steps and no runbooks.
type OctopusEmptyProjectCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusEmptyProjectCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusEmptyProjectCheck {
	return OctopusEmptyProjectCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusEmptyProjectCheck) Id() string {
	return OctoLintEmptyProject
}

func (o OctopusEmptyProjectCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
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
		o.config.MaxEmptyProjectCheckProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	runbooks, err := o.client.Runbooks.GetAll()

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	emptyProjects := threadsafe.NewSlice[string]()
	goroutineErrors := threadsafe.NewSlice[error]()

	for i, p := range projects {
		i := i
		p := p

		g.Go(func() error {
			zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(projects))*100) + "% complete")

			stepCount, err := o.stepsInDeploymentProcess(p.DeploymentProcessID)

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					goroutineErrors.Append(err)
				}
				return nil
			}

			if runbooksInProject(p.ID, runbooks) == 0 && stepCount == 0 {
				emptyProjects.Append(p.Name)
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

	if emptyProjects.Length() > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following projects have no runbooks and no deployment process:\n"+strings.Join(emptyProjects.Values(), "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no empty projects",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func runbooksInProject(projectID string, runbooks []*runbooks.Runbook) int {
	count := 0
	for _, r := range runbooks {
		if r.ProjectID == projectID {
			count++
		}
	}
	return count
}

func (o OctopusEmptyProjectCheck) stepsInDeploymentProcess(deploymentProcessID string) (int, error) {
	if deploymentProcessID == "" {
		return 0, nil
	}

	resource, err := o.client.DeploymentProcesses.GetByID(deploymentProcessID)

	if err != nil {
		return 0, err
	}

	return len(resource.Steps), nil
}
