package organization

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/hayageek/threadsafe"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"strings"
)

const maxStepCount = 20
const OctoLintTooManySteps = "OctoLintTooManySteps"

// OctopusProjectTooManyStepsCheck checks to see if any project has too many steps.
type OctopusProjectTooManyStepsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusProjectTooManyStepsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusProjectTooManyStepsCheck {
	return OctopusProjectTooManyStepsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusProjectTooManyStepsCheck) Id() string {
	return OctoLintTooManySteps
}

func (o OctopusProjectTooManyStepsCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
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
		o.config.MaxProjectStepsProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	complexProjects := threadsafe.NewSlice[string]()
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

			if stepCount >= maxStepCount {
				complexProjects.Append(p.Name)
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

	if complexProjects.Length() > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following projects have 20 or more steps:\n"+strings.Join(complexProjects.Values(), "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no projects with too many steps",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusProjectTooManyStepsCheck) stepsInDeploymentProcess(deploymentProcessID string) (int, error) {
	if deploymentProcessID == "" {
		return 0, nil
	}

	resource, err := o.client.DeploymentProcesses.GetByID(deploymentProcessID)

	if err != nil {
		// If we can't find the deployment process, assume zero steps
		if err.(*core.APIError).StatusCode == 404 {
			return 0, nil
		}
		return 0, err
	}

	return len(resource.Steps), nil
}
