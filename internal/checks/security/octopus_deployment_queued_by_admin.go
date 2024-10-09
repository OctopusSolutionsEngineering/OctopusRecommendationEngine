package security

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/events"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"strings"
	"time"
)

const OctoLintDeploymentQueuedByAdmin = "OctoLintDeploymentQueuedByAdmin"

// OctopusDeploymentQueuedByAdminCheck checks to see if any deployments were initiated by someone from the admin teams.
// This usually means that a more specific and limited user should be created to perform deployments.
type OctopusDeploymentQueuedByAdminCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusDeploymentQueuedByAdminCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusDeploymentQueuedByAdminCheck {
	return OctopusDeploymentQueuedByAdminCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusDeploymentQueuedByAdminCheck) Id() string {
	return OctoLintDeploymentQueuedByAdmin
}

func (o OctopusDeploymentQueuedByAdminCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
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
		o.config.MaxDeploymentsByAdminProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	teams, err := o.getAdminTeams()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	projectsDeployedByAdmins := []string{}

	now := time.Now()
	fromDate := now.AddDate(0, -3, 0)
	from := fromDate.Format("2006-01-02")

	for i, p := range projects {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(projects))*100) + "% complete")

		projectId := p.ID
		usersWhoDeployedProject := []string{}

		resource, err := o.client.Events.Get(events.EventsQuery{
			EventCategories: []string{"DeploymentQueued"},
			Projects:        []string{projectId},
			Skip:            0,
			Take:            100,
			From:            from,
		})

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		if resource != nil {
			for _, r := range resource.Items {
				if r.Username == "system" {
					continue
				}

				user, err := o.client.Users.Get(users.UsersQuery{
					Filter: r.Username,
					Skip:   0,
					Take:   1,
				})

				if err != nil {
					if !o.errorHandler.ShouldContinue(err) {
						return nil, err
					}
					continue
				}

				for _, u := range user.Items {
					for _, t := range teams {
						if slices.Index(t.MemberUserIDs, u.ID) != -1 && slices.Index(usersWhoDeployedProject, u.Username) == -1 {
							usersWhoDeployedProject = append(usersWhoDeployedProject, u.Username)
						}
					}
				}
			}
		}

		if len(usersWhoDeployedProject) != 0 {
			projectsDeployedByAdmins = append(projectsDeployedByAdmins, p.Name+" ("+strings.Join(usersWhoDeployedProject, ",")+")")
		}
	}

	if len(projectsDeployedByAdmins) != 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following projects were deployed by admins. Consider creating a limited user account to perform deployments:\n"+strings.Join(projectsDeployedByAdmins, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"No deployments were found",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}

func (o OctopusDeploymentQueuedByAdminCheck) getAdminTeams() ([]*teams.Team, error) {
	adminTeams := []string{"Octopus Administrators", "Space Managers", "Octopus Managers"}

	teamResources := []*teams.Team{}
	for _, adminTeam := range adminTeams {
		team, err := o.client.Teams.Get(teams.TeamsQuery{
			IDs:           nil,
			IncludeSystem: true,
			PartialName:   adminTeam,
			Skip:          0,
			Take:          1,
		})

		if err != nil {
			return nil, err
		}

		teamResources = append(teamResources, team.Items...)
	}

	return teamResources, nil
}
