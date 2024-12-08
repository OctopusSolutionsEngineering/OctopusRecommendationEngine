package factory

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/naming"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/organization"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/performance"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/security"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"strings"
)

// OctopusCheckFactory builds all the lint checks. This is where you can customize things like error handlers.
type OctopusCheckFactory struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	url          string
	space        string
}

func NewOctopusCheckFactory(client *client.Client, url string, space string) OctopusCheckFactory {
	return OctopusCheckFactory{client: client, url: url, space: space, errorHandler: checks.OctopusClientPermissiveErrorHandler{}}
}

// BuildAllChecks creates new instances of all the checks and returns them as an array.
func (o OctopusCheckFactory) BuildAllChecks(config *config.OctolintConfig) ([]checks.OctopusCheck, error) {
	skipChecksSlice := lo.FilterMap(strings.Split(config.SkipTests, ","), func(item string, index int) (string, bool) {
		itemTrimmed := strings.TrimSpace(item)
		return itemTrimmed, len(itemTrimmed) != 0
	})

	onlyChecksSlice := lo.FilterMap(strings.Split(config.OnlyTests, ","), func(item string, index int) (string, bool) {
		itemTrimmed := strings.TrimSpace(item)
		return itemTrimmed, len(itemTrimmed) != 0
	})

	allChecks := []checks.OctopusCheck{
		security.NewOctopusUnrotatedAccountsCheck(o.client, config, o.errorHandler),
		security.NewOctopusDeploymentQueuedByAdminCheck(o.client, config, o.errorHandler),
		security.NewOctopusPerpetualApiKeysCheck(o.client, config, o.errorHandler),
		security.NewOctopusDuplicatedGitCredentialsCheck(o.client, config, o.errorHandler),
		security.NewOctopusInsecureK8sCheck(o.client, config, o.errorHandler),
		security.NewOctopusInsecureFeedsCheck(o.client, config, o.errorHandler),
		security.NewOctopusInsecureSubscriptionsCheck(o.client, config, o.errorHandler),
		organization.NewOctopusEnvironmentCountCheck(o.client, config, o.errorHandler),
		organization.NewOctopusDefaultProjectGroupCountCheck(o.client, config, o.errorHandler),
		organization.NewOctopusEmptyProjectCheck(o.client, config, o.errorHandler),
		organization.NewOctopusUnusedVariablesCheck(o.client, config, o.errorHandler),
		organization.NewOctopusDuplicatedVariablesCheck(o.client, config, o.errorHandler),
		organization.NewOctopusProjectTooManyStepsCheck(o.client, config, o.errorHandler),
		organization.NewOctopusLifecycleRetentionPolicyCheck(o.client, config, o.errorHandler),
		organization.NewOctopusUnusedTargetsCheck(o.client, config, o.errorHandler),
		organization.NewOctopusProjectSpecificEnvironmentCheck(o.client, config, o.errorHandler),
		organization.NewOctopusTenantsInsteadOfTagsCheck(o.client, config, o.errorHandler),
		organization.NewOctopusProjectGroupsWithExclusiveEnvironmentsCheck(o.client, config, o.errorHandler),
		organization.NewOctopusUnhealthyTargetCheck(o.client, config, o.errorHandler),
		organization.NewOctopusUnusedProjectsCheck(o.client, config, o.errorHandler),
		organization.NewOctopusUnusedTenantsCheck(o.client, config, o.errorHandler),
		performance.NewOctopusDeploymentQueuedTimeCheck(o.client, config, o.url, o.space, o.errorHandler),
		naming.NewOctopusProjectContainerImageRegex(o.client, config, o.errorHandler),
		naming.NewOctopusInvalidVariableNameCheck(o.client, config, o.errorHandler),
		naming.NewOctopusInvalidTargetName(o.client, config, o.errorHandler),
		naming.NewOctopusInvalidTargetRole(o.client, config, o.errorHandler),
		naming.NewOctopusProjectReleaseTemplateRegex(o.client, config, o.errorHandler),
		naming.NewOctopusProjectWorkerPoolRegex(o.client, config, o.errorHandler),
		naming.NewOctopusInvalidLifecycleName(o.client, config, o.errorHandler),
		naming.NewOctopusProjectDefaultStepNames(o.client, config, o.errorHandler),
	}

	return lo.Filter(allChecks, func(item checks.OctopusCheck, index int) bool {
		return slices.Index(skipChecksSlice, item.Id()) == -1 &&
			(len(onlyChecksSlice) == 0 || slices.Index(onlyChecksSlice, item.Id()) != -1)
	}), nil
}
