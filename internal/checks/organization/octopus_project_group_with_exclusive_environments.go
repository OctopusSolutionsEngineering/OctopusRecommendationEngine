package organization

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/hayageek/threadsafe"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
	"strings"
)

const OctoLintProjectGroupsWithExclusiveEnvironments = "OctoLintProjectGroupsWithExclusiveEnvironments"

// OctopusProjectGroupsWithExclusiveEnvironmentsCheck checks to see if the project groups contain projects that have mutually exclusive environments.
type OctopusProjectGroupsWithExclusiveEnvironmentsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusProjectGroupsWithExclusiveEnvironmentsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusProjectGroupsWithExclusiveEnvironmentsCheck {
	return OctopusProjectGroupsWithExclusiveEnvironmentsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusProjectGroupsWithExclusiveEnvironmentsCheck) Id() string {
	return OctoLintProjectGroupsWithExclusiveEnvironments
}

func (o OctopusProjectGroupsWithExclusiveEnvironmentsCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	allProjectGroups, err := o.client.ProjectGroups.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allProjects, err := client_wrapper.GetProjectsWithFilter(
		o.client,
		o.client.GetSpaceID(),
		o.config.ExcludeProjectsExcept,
		o.config.ExcludeProjects,
		o.config.MaxExclusiveEnvironmentsProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allLifecycles, err := o.client.Lifecycles.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(concurrency)

	projectGroupsWithExclusiveEnvs := threadsafe.NewSlice[string]()

	for i, pg := range allProjectGroups {

		i := i
		pg := pg

		g.Go(func() error {

			zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allProjectGroups))*100) + "% complete")

			// Find the groups of environments captured in the default lifecyles of the projects in the project group
			envGroups := [][]string{}
			for _, p := range allProjects {
				if p.ProjectGroupID == pg.ID {
					lifecycle := o.getLifecycleById(allLifecycles, p.LifecycleID)
					projectEnvironments := o.getLifecycleEnvironments(lifecycle)
					envGroups = append(envGroups, projectEnvironments)
				}
			}

			// don't do any further processing if there was just one project
			if len(envGroups) <= 1 {
				return nil
			}

			// Attempt to find at least every environment in a lifecycle with one in another environment
			for i, eg1 := range envGroups[0 : len(envGroups)-1] {
				allExclusive := true

				for _, eg2 := range envGroups[i+1:] {
					for _, e1 := range eg1 {
						if slices.Index(eg2, e1) != -1 {
							allExclusive = false
							break
						}
					}
				}

				// if none of the environments from this lifecycle are found in any other lifecycles, we have an project with exclusive environments
				if allExclusive && !projectGroupsWithExclusiveEnvs.Contains(pg.Name) {
					projectGroupsWithExclusiveEnvs.Append(pg.Name)
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if projectGroupsWithExclusiveEnvs.Length() > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following project groups contain projects with mutually exclusive environments in their default lifecycle:\n"+strings.Join(projectGroupsWithExclusiveEnvs.Values(), "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no project groups with mutually exclusive lifecycles",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusProjectGroupsWithExclusiveEnvironmentsCheck) getLifecycleEnvironments(lifecycle *lifecycles.Lifecycle) []string {
	projectEnvironments := []string{}
	for _, phase := range lifecycle.Phases {
		projectEnvironments = append(projectEnvironments, phase.AutomaticDeploymentTargets...)
		projectEnvironments = append(projectEnvironments, phase.OptionalDeploymentTargets...)
	}
	slices.Sort(projectEnvironments)
	return projectEnvironments
}

func (o OctopusProjectGroupsWithExclusiveEnvironmentsCheck) getLifecycleById(lifecycles []*lifecycles.Lifecycle, id string) *lifecycles.Lifecycle {
	for _, l := range lifecycles {
		if l.ID == id {
			return l
		}
	}

	return nil
}
