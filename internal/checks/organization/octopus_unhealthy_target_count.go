package organization

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/events"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/hayageek/threadsafe"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

func (o OctopusUnhealthyTargetCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
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

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	unhealthyMachines := threadsafe.NewSlice[string]()
	goroutineErrors := threadsafe.NewSlice[error]()

	for i, m := range allMachines {
		i := i
		m := m

		g.Go(func() error {

			zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allMachines))*100) + "% complete")

			wasEverHealthy := true
			if m.HealthStatus == "Unhealthy" {
				wasEverHealthy = false

				targetEvents, err := o.client.Events.Get(events.EventsQuery{
					Regarding: m.ID,
				})

				if err != nil {
					if !o.errorHandler.ShouldContinue(err) {
						goroutineErrors.Append(err)
					}
					return nil
				}

				for _, e := range targetEvents.Items {
					if e.Category == "MachineHealthy" && time.Now().Sub(e.Occurred) < maxHealthCheckTime {
						wasEverHealthy = true
						break
					}
				}
			}

			if !wasEverHealthy {
				unhealthyMachines.Append(m.Name)
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

	if unhealthyMachines.Length() > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following targets have not been healthy in the last 30 days:\n"+strings.Join(unhealthyMachines.Values(), "\n"),
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
