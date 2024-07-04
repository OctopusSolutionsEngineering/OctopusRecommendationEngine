package naming

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"regexp"
	"strings"
)

const OctoLintInvalidTargetRoles = "OctoLintInvalidTargetRoles"

// OctopusInvalidTargetRole find targets that have not been healthy in the last 30 days.
type OctopusInvalidTargetRole struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusInvalidTargetRole(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusInvalidTargetRole {
	return OctopusInvalidTargetRole{
		client:       client,
		errorHandler: errorHandler,
		config:       config,
	}
}

func (o OctopusInvalidTargetRole) Id() string {
	return OctoLintInvalidTargetRoles
}

func (o OctopusInvalidTargetRole) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	if strings.TrimSpace(o.config.TargetRoleRegex) == "" {
		return nil, nil
	}

	regex, err := regexp.Compile(o.config.TargetRoleRegex)

	if err != nil {
		return checks.NewOctopusCheckResultImpl(
			"The supplied regex "+o.config.TargetNameRegex+" does not compile",
			o.Id(),
			"",
			checks.Error,
			checks.Naming), nil
	}

	allMachines, err := client_wrapper.GetMachines(o.config.MaxInvalidRoleTargets, o.client, o.client.GetSpaceID())

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Naming, err)
	}

	responses := []string{}
	for i, m := range allMachines {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allMachines))*100) + "% complete")

		invalidRoles := []string{}
		for _, r := range m.Roles {
			if !regex.Match([]byte(r)) {
				invalidRoles = append(invalidRoles, r)
			}
		}

		if len(invalidRoles) != 0 {
			responses = append(responses, m.Name+" - "+strings.Join(invalidRoles, ","))
		}
	}

	if len(responses) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following target roles do not match the regex "+o.config.TargetRoleRegex+":\n"+strings.Join(responses, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Naming), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"All targets match the regex "+o.config.TargetNameRegex,
		o.Id(),
		"",
		checks.Ok,
		checks.Naming), nil
}
