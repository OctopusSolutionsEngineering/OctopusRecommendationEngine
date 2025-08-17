package security

import (
	"errors"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
)

// CustomProject is the simplest representation of a project and its version controlled settings
type CustomProject struct {
	PersistenceSettings CustomPersistenceSettings `json:"PersistenceSettings"`
	Name                string                    `json:"Name"`
}

type CustomPersistenceSettings struct {
	Type        string             `json:"Type"`
	Credentials *CustomCredentials `json:"Credentials"`
}

type CustomCredentials struct {
	Type     string  `json:"Type"`
	Username *string `json:"Username"`
}

// OctopusDuplicatedGitCredentialsCheck reports on any perpetual api keys
type OctopusDuplicatedGitCredentialsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusDuplicatedGitCredentialsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusDuplicatedGitCredentialsCheck {
	return OctopusDuplicatedGitCredentialsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusDuplicatedGitCredentialsCheck) Id() string {
	return "OctoLintSharedGitUsername"
}

func (o OctopusDuplicatedGitCredentialsCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	url := o.client.HttpSession().BaseURL.String() + "/api/" + o.client.GetSpaceID() + "/Projects?take=1000"
	allProjects, err := newclient.Get[resources.Resources[CustomProject]](o.client.HttpSession(), url)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	gitUsernameCounts := map[string]int{}
	gitUsernameProjects := map[string][]string{}
	for i, p := range allProjects.Items {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allProjects.Items))*100) + "% complete")

		if p.PersistenceSettings.Type == "VersionControlled" &&
			p.PersistenceSettings.Credentials.Type == "UsernamePassword" &&
			p.PersistenceSettings.Credentials.Username != nil {
			if _, ok := gitUsernameCounts[*p.PersistenceSettings.Credentials.Username]; !ok {
				gitUsernameCounts[*p.PersistenceSettings.Credentials.Username] = 0
			}
			gitUsernameCounts[*p.PersistenceSettings.Credentials.Username]++

			if _, ok := gitUsernameProjects[*p.PersistenceSettings.Credentials.Username]; !ok {
				gitUsernameProjects[*p.PersistenceSettings.Credentials.Username] = []string{}
			}
			gitUsernameProjects[*p.PersistenceSettings.Credentials.Username] = append(gitUsernameProjects[*p.PersistenceSettings.Credentials.Username], p.Name)
		}
	}

	duplicatedGitCredentials := map[string][]string{}
	for u, c := range gitUsernameCounts {
		if c > 1 {
			duplicatedGitCredentials[u] = gitUsernameProjects[u]
		}
	}

	if len(duplicatedGitCredentials) != 0 {
		message := []string{}
		for u, p := range duplicatedGitCredentials {
			message = append(message, u+" ("+strings.Join(p, ", ")+")")
		}

		return checks.NewOctopusCheckResultImpl(
			"The following Git usernames have been reused across the following projects:\n"+strings.Join(message, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"No Git usernames have been reused",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}
