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
	"golang.org/x/sync/errgroup"
	"regexp"
	"strings"
)

const OctoLintContainerImageName = "OctoLintProjectContainerImageName"

// OctopusProjectContainerImageRegex checks to see if any project has too many steps.
type OctopusProjectContainerImageRegex struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusProjectContainerImageRegex(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusProjectContainerImageRegex {
	return OctopusProjectContainerImageRegex{
		client:       client,
		errorHandler: errorHandler,
		config:       config,
	}
}

func (o OctopusProjectContainerImageRegex) Id() string {
	return OctoLintContainerImageName
}

func (o OctopusProjectContainerImageRegex) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	if strings.TrimSpace(o.config.ContainerImageRegex) == "" {
		return nil, nil
	}

	regex, err := regexp.Compile(o.config.ContainerImageRegex)

	if err != nil {

		return checks.NewOctopusCheckResultImpl(
			"The supplied regex "+o.config.ContainerImageRegex+" does not compile",
			o.Id(),
			"",
			checks.Error,
			checks.Naming), nil
	}

	projects, err := client_wrapper.GetProjectsWithFilter(
		o.client,
		o.client.GetSpaceID(),
		o.config.ExcludeProjectsExcept,
		o.config.ExcludeProjects,
		o.config.MaxInvalidContainerImageProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Naming, err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	actionsWithInvalidImages := threadsafe.NewSlice[string]()
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
					if a.Container == nil || strings.TrimSpace(a.Container.Image) == "" {
						continue
					}

					if !regex.Match([]byte(a.Container.Image)) {
						actionsWithInvalidImages.Append(p.Name + "/" + a.Name + ": " + a.Container.Image)
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

	if actionsWithInvalidImages.Length() > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following project actions do not match the regex "+o.config.ContainerImageRegex+":\n"+strings.Join(actionsWithInvalidImages.Values(), "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Naming), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no project actions with invalid container images",
		o.Id(),
		"",
		checks.Ok,
		checks.Naming), nil
}

func (o OctopusProjectContainerImageRegex) stepsInDeploymentProcess(deploymentProcessID string) (*deployments.DeploymentProcess, error) {
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
