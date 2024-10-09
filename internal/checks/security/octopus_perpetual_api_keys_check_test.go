package security

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

func TestPerpetualApiKeys(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		dir := filepath.Join("..", "..", "..", "test", "terraform")
		newSpaceId, err := testFramework.Act(t, container, dir, "1-singlespace", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		allUsers, err := newSpaceClient.Users.GetAll()

		if err != nil {
			return err
		}

		_, err = newSpaceClient.APIKeys.Create(&users.CreateAPIKey{
			APIKey:  "",
			Created: nil,
			Purpose: "",
			UserID:  allUsers[0].ID,
			Expires: nil,
		})

		if err != nil {
			return err
		}

		check := NewOctopusPerpetualApiKeysCheck(newSpaceClient, &config.OctolintConfig{}, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute(2)

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have returned a warning")
		}

		return nil
	})
}
