package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
)

const OctopusEnvironmentCountCheckName = "OctoLintEnvironmentCount"

// OctopusEnvironmentCountCheck checks to see if too many environments have been created in a space.
type OctopusEnvironmentCountCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusEnvironmentCountCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusEnvironmentCountCheck {
	return OctopusEnvironmentCountCheck{client: client, errorHandler: errorHandler, config: config}
}

func (o OctopusEnvironmentCountCheck) Id() string {
	return OctopusEnvironmentCountCheckName
}

func (o OctopusEnvironmentCountCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	query := environments.EnvironmentsQuery{
		PartialName: "",
		Skip:        0,
		Take:        1000,
	}
	resources, err := o.client.Environments.Get(query)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	if len(resources.Items) > o.config.MaxEnvironments {
		return checks.NewOctopusCheckResultImpl(
			"The recommended maximum number of environments is "+fmt.Sprint(o.config.MaxEnvironments)+". You have at least "+fmt.Sprint(len(resources.Items)),
			o.Id(),
			"https://octopus.com/docs/getting-started/best-practices/environments-and-deployment-targets-and-roles#environments",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"The number of environments in the space is OK",
		o.Id(),
		"https://octopus.com/docs/getting-started/best-practices/environments-and-deployment-targets-and-roles#environments",
		checks.Ok,
		checks.Organization), nil
}
