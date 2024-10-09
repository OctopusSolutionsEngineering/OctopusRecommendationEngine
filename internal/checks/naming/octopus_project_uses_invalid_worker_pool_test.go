package naming

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

func TestInvalidWorkerPool(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(
			t,
			container,
			filepath.Join("..", "..", "..", "test", "terraform"), "30-defaultstepnames", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		check := NewOctopusProjectWorkerPoolRegex(
			newSpaceClient,
			&config.OctolintConfig{
				ProjectStepWorkerPoolRegex: "kubernetes",
			},
			checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute(2)

		if err != nil {
			return err
		}

		// Assert
		if result == nil || result.Severity() != checks.Warning {
			return errors.New("check should have produced a warning")
		}

		return nil
	})
}
