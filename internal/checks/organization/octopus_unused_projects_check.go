package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tasks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"strings"
	"time"
)

const OctopusUnusedProjectsCheckName = "OctoLintUnusedProjects"

// OctopusUnusedProjectsCheck find projects that have not had a deployment in the last 30 days
type OctopusUnusedProjectsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusUnusedProjectsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusUnusedProjectsCheck {
	return OctopusUnusedProjectsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusUnusedProjectsCheck) Id() string {
	return OctopusUnusedProjectsCheckName
}

func (o OctopusUnusedProjectsCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	projects, err := client_wrapper.GetProjectsWithFilter(
		o.client,
		o.client.GetSpaceID(),
		o.config.ExcludeProjectsExcept,
		o.config.ExcludeProjects,
		o.config.MaxUnusedProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	unusedProjects := []string{}
	for i, project := range projects {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(projects))*100) + "% complete")

		// Ignore disabled projects
		if project.IsDisabled {
			continue
		}

		projectHasTask := false

		tasks, err := o.client.Tasks.Get(tasks.TasksQuery{
			Project: project.ID,
			Skip:    0,
			Take:    100,
		})

		if err != nil {
			return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
		}

		for _, task := range tasks.Items {
			if task.StartTime != nil && task.StartTime.After(time.Now().Add(-time.Hour*24*time.Duration(o.config.MaxDaysSinceLastTask))) {
				projectHasTask = true
				break
			}
		}

		if !projectHasTask {
			unusedProjects = append(unusedProjects, project.Name)
		}
	}

	daysString := fmt.Sprintf("%d", o.config.MaxDaysSinceLastTask)

	if len(unusedProjects) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following projects have not had any tasks in "+daysString+" days:\n"+strings.Join(unusedProjects, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no projects that have not had any tasks in the last "+daysString+" days",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}
