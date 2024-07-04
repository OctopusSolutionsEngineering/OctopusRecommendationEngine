package security

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strings"
)

const OctoLintInsecureK8sTargets = "OctoLintInsecureK8sTargets"

// OctopusInsecureK8sCheck checks to see if any targets have not been used in a month
type OctopusInsecureK8sCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusInsecureK8sCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusInsecureK8sCheck {
	return OctopusInsecureK8sCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusInsecureK8sCheck) Id() string {
	return OctoLintInsecureK8sTargets
}

func (o OctopusInsecureK8sCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	targets, err := client_wrapper.GetMachines(o.config.MaxInsecureK8sTargets, o.client, o.client.GetSpaceID())

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	k8sTargets := lo.Filter(targets, func(item *machines.DeploymentTarget, index int) bool {
		return item.Endpoint != nil && item.Endpoint.GetCommunicationStyle() == "Kubernetes"
	})

	insecureMachines := []string{}
	for i, m := range k8sTargets {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(k8sTargets))*100) + "% complete")

		k8sEndpoint := m.Endpoint.(*machines.KubernetesEndpoint)
		if k8sEndpoint.SkipTLSVerification || (k8sEndpoint.ClusterURL != nil && strings.HasPrefix(k8sEndpoint.ClusterURL.String(), "http://")) {
			insecureMachines = append(insecureMachines, m.Name)
		}

	}

	if len(insecureMachines) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following Kubernetes targets skip TLS validation or use an insecure HTTP endpoint:\n"+strings.Join(insecureMachines, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no insecure Kubernetes targets",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}
