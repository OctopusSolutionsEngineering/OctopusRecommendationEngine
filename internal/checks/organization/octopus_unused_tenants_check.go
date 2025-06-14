package organization

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tasks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/hayageek/threadsafe"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"strings"
	"time"
)

const OctopusUnusedTenantsCheckName = "OctoLintUnusedTenants"

// OctopusUnusedTenantsCheck find projects that have not had a deployment in the last 30 days
type OctopusUnusedTenantsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusUnusedTenantsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusUnusedTenantsCheck {
	return OctopusUnusedTenantsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusUnusedTenantsCheck) Id() string {
	return OctopusUnusedTenantsCheckName
}

func (o OctopusUnusedTenantsCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	tenants, err := client_wrapper.GetTenants(
		o.config.MaxUnusedTenants,
		o.client,
		o.client.GetSpaceID())

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	unusedTenants := threadsafe.NewSlice[string]()
	goroutineErrors := threadsafe.NewSlice[error]()

	for i, tenant := range tenants {
		i := i
		tenant := tenant

		g.Go(func() error {
			zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(tenants))*100) + "% complete")

			// Ignore disabled projects
			if tenant.IsDisabled {
				return nil
			}

			tenantHasTask := false

			tasks, err := o.client.Tasks.Get(tasks.TasksQuery{
				Tenant: tenant.ID,
				Skip:   0,
				Take:   100,
			})

			if err != nil {
				goroutineErrors.Append(err)
				return nil
			}

			for _, task := range tasks.Items {
				if task.StartTime != nil && task.StartTime.After(time.Now().Add(-time.Hour*24*time.Duration(o.config.MaxDaysSinceLastTask))) {
					tenantHasTask = true
					break
				}
			}

			if !tenantHasTask {
				unusedTenants.Append(tenant.Name)
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

	daysString := fmt.Sprintf("%d", o.config.MaxDaysSinceLastTask)

	if unusedTenants.Length() > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following tenants have not had any tasks in "+daysString+" days:\n"+strings.Join(unusedTenants.Values(), "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no tenants that have not had any tasks in the last "+daysString+" days",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}
