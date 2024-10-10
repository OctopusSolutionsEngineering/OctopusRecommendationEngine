package organization

import (
	"errors"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tasks"
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

		err = startHealthCheck(newSpaceClient)

		if err != nil {
			return err
		}

		// loop for a bit until the target is unhealthy
		for i := 0; i < 24; i++ {
			if unhealthy, err := checkMachinesUnhealthy(newSpaceClient); err != nil {
				return err
			} else {
				if unhealthy {
					break
				}
			}

			time.Sleep(time.Second * 10)
		}

		// A final sanity check to make sure the machine is actually unhealthy
		if unhealthy, err := checkMachinesUnhealthy(newSpaceClient); err != nil {
			return err
		} else {
			if !unhealthy {
				return errors.New("machine was never unhealthy")
			}
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

func startHealthCheck(newSpaceClient *client.Client) error {
	machines, err := client_wrapper.GetMachines(0, newSpaceClient, newSpaceClient.GetSpaceID())

	if err != nil {
		return err
	}

	for _, machine := range machines {

		task := tasks.NewTask()
		task.Name = "Health"
		task.SpaceID = newSpaceClient.GetSpaceID()
		task.Description = machine.Name
		task.Arguments = map[string]any{
			"Timeout":        "00:05:00",
			"MachineTimeout": "00:05:00",
			"EnvironmentId":  machine.EnvironmentIDs[0],
			"MachineIds":     []string{machine.ID},
		}

		_, err := newSpaceClient.Tasks.Add(task)

		if err != nil {
			return err
		}
	}

	return nil
}

func checkMachinesUnhealthy(newSpaceClient *client.Client) (bool, error) {
	machines, err := client_wrapper.GetMachines(0, newSpaceClient, newSpaceClient.GetSpaceID())

	if err != nil {
		return false, err
	}

	if len(machines) > 0 && machines[0].HealthStatus != "Healthy" && machines[0].HealthStatus != "Unknown" {
		return true, nil
	}

	return false, nil
}
