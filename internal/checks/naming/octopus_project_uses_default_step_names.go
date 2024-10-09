package naming

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/hayageek/threadsafe"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
	"strings"
)

const OctoLintProjectDefaultStepNames = "OctoLintProjectDefaultStepNames"

// OctopusProjectDefaultStepNames checks to see if any project has too many steps.
type OctopusProjectDefaultStepNames struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusProjectDefaultStepNames(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusProjectDefaultStepNames {
	return OctopusProjectDefaultStepNames{
		client:       client,
		errorHandler: errorHandler,
		config:       config,
	}
}

func (o OctopusProjectDefaultStepNames) Id() string {
	return OctoLintProjectDefaultStepNames
}

func (o OctopusProjectDefaultStepNames) Execute(concurrency int) (checks.OctopusCheckResult, error) {
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
		o.config.MaxDefaultStepNameProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Naming, err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	actionsWithDefaultNames := threadsafe.NewSlice[string]()
	goroutineErrors := threadsafe.NewSlice[error]()

	for i, p := range projects {
		i := i
		p := p

		g.Go(func() error {
			zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(projects))*100) + "% complete")

			deploymentProcess, err := o.stepsInDeploymentProcess(p.DeploymentProcessID)

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					goroutineErrors.Append(err)
				}
				return nil
			}

			if deploymentProcess == nil {
				return nil
			}

			for _, s := range deploymentProcess.Steps {
				for _, a := range s.Actions {
					if slices.Index(checks.DefaultStepNames, a.Name) != -1 {
						actionsWithDefaultNames.Append(p.Name + "/" + a.Name)
					}
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

	if actionsWithDefaultNames.Length() > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following project actions use the default step names:\n"+strings.Join(actionsWithDefaultNames.Values(), "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Naming), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no project actions default step names",
		o.Id(),
		"",
		checks.Ok,
		checks.Naming), nil
}

func (o OctopusProjectDefaultStepNames) stepsInDeploymentProcess(deploymentProcessID string) (*deployments.DeploymentProcess, error) {
	if deploymentProcessID == "" {
		return nil, nil
	}

	resource, err := o.client.DeploymentProcesses.GetByID(deploymentProcessID)

	if err != nil {
		// If we can't find the deployment process, assume zero steps
		if err.(*core.APIError).StatusCode == 404 {
			return nil, nil
		}
		return nil, err
	}

	return resource, nil
}
