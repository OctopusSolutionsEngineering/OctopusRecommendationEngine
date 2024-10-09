package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"strings"
)

type OctopusLifecycleRetentionPolicyCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusLifecycleRetentionPolicyCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusLifecycleRetentionPolicyCheck {
	return OctopusLifecycleRetentionPolicyCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusLifecycleRetentionPolicyCheck) Id() string {
	return "OctoRecLifecycleRetention"
}

func (o OctopusLifecycleRetentionPolicyCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	lifecycles, err := o.client.Lifecycles.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	keepsForever := []string{}
	for i, l := range lifecycles {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(lifecycles))*100) + "% complete")

		phaseKeepsForever, err := o.anyPhasesKeepForever(l.Phases)

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		lifecycleKeepsForever := l.ReleaseRetentionPolicy.ShouldKeepForever || l.TentacleRetentionPolicy.ShouldKeepForever

		if lifecycleKeepsForever || phaseKeepsForever {
			keepsForever = append(keepsForever, l.Name)
		}
	}

	if len(keepsForever) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following lifecycles have retention policies that keep releases or files forever:\n"+strings.Join(keepsForever, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no lifecycles with retention policies that keep releases or files forever",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusLifecycleRetentionPolicyCheck) anyPhasesKeepForever(phases []*lifecycles.Phase) (bool, error) {
	if len(phases) == 0 {
		return false, nil
	}

	for _, p := range phases {
		keepReleasesForver := p.ReleaseRetentionPolicy != nil && p.ReleaseRetentionPolicy.ShouldKeepForever
		keepFilesForever := p.TentacleRetentionPolicy != nil && p.TentacleRetentionPolicy.ShouldKeepForever

		if keepReleasesForver || keepFilesForever {
			return true, nil
		}
	}

	return false, nil
}
