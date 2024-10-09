package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
)

const maxProjectsInDefaultGroup = 10

// OctopusDefaultProjectGroupCountCheck checks to see if the default project group contains too many projects. This is
// usually an indication that additional projects groups should be created to organize the dashboard.
type OctopusDefaultProjectGroupCountCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusDefaultProjectGroupCountCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusDefaultProjectGroupCountCheck {
	return OctopusDefaultProjectGroupCountCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusDefaultProjectGroupCountCheck) Id() string {
	return "OctoLintDefaultProjectGroupChildCount"
}

func (o OctopusDefaultProjectGroupCountCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	resource, err := o.client.ProjectGroups.GetByName("Default Project Group")

	if err != nil {
		apiError, ok := err.(*core.APIError)
		if (ok && apiError.StatusCode == 404) || err.Error() == "cannot find the item" {
			return checks.NewOctopusCheckResultImpl(
				"The default project group was not found",
				o.Id(),
				"",
				checks.Ok,
				checks.Organization), nil
		}
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	if resource != nil {

		projects, err := o.client.ProjectGroups.GetProjects(resource)

		if err != nil {
			return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
		}

		if len(projects) > maxProjectsInDefaultGroup {
			return checks.NewOctopusCheckResultImpl(
				"The default project group contains "+fmt.Sprint(len(projects))+" projects. You may want to organize these projects into additional project groups.",
				o.Id(),
				"",
				checks.Warning,
				checks.Organization), nil
		}
	}

	return checks.NewOctopusCheckResultImpl(
		"The number of projects in the default project group is OK",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}
