package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
	"time"
)

func TestUnhealthyTargets(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "27-unhealthytargets", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		// loop for a bit until the target is unhealthy
		for i := 0; i < 12; i++ {
			machines, err := client_wrapper.GetMachines(0, newSpaceClient, newSpaceClient.GetSpaceID())

			if err != nil {
				return err
			}

			if len(machines) > 0 && machines[0].HealthStatus != "Healthy" && machines[0].HealthStatus != "Unknown" {
				break
			}

			time.Sleep(time.Second * 10)
		}

		check := NewOctopusUnhealthyTargetCheck(newSpaceClient, &config.OctolintConfig{}, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute(2)

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have failed")
		}

		return nil
	})
}
