package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/events"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"strings"
	"time"
)

const maxHealthCheckTime = time.Hour * 24 * 30
const OctoLintUnhealthyTargets = "OctoLintUnhealthyTargets"

// OctopusUnhealthyTargetCheck find targets that have not been healthy in the last 30 days.
type OctopusUnhealthyTargetCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusUnhealthyTargetCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusUnhealthyTargetCheck {
	return OctopusUnhealthyTargetCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusUnhealthyTargetCheck) Id() string {
	return OctoLintUnhealthyTargets
}

func (o OctopusUnhealthyTargetCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	allMachines, err := client_wrapper.GetMachines(o.config.MaxUnhealthyTargets, o.client, o.client.GetSpaceID())

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	unhealthyMachines := []string{}
	for i, m := range allMachines {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allMachines))*100) + "% complete")

		wasEverHealthy := true
		if m.HealthStatus == "Unhealthy" {
			wasEverHealthy = false

			targetEvents, err := o.client.Events.Get(events.EventsQuery{
				Regarding: m.ID,
			})

			if err != nil {
				if !o.errorHandler.ShouldContinue(err) {
					return nil, err
				}
				continue
			}

			for _, e := range targetEvents.Items {
				if e.Category == "MachineHealthy" && time.Now().Sub(e.Occurred) < maxHealthCheckTime {
					wasEverHealthy = true
					break
				}
			}
		}

		if !wasEverHealthy {
			unhealthyMachines = append(unhealthyMachines, m.Name)
		}
	}

	if len(unhealthyMachines) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following targets have not been healthy in the last 30 days:\n"+strings.Join(unhealthyMachines, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no targets that were unhealthy for all of the last 30 days",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}
