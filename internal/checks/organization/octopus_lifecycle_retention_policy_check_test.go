package organization

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

func deleteDefaultLifecycle(newSpaceClient *client.Client) error {
	// Delete the default lifecycle, which keeps things forever, and messes with these tests
	defaultLifecycle, err := newSpaceClient.Lifecycles.GetByName("Default Lifecycle")
	if err != nil {
		return err
	}

	err = newSpaceClient.Lifecycles.DeleteByID(defaultLifecycle.ID)
	if err != nil {
		return err
	}

	return nil
}

func TestLifecyclesMeetRecommendations(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "16-lifecyclesmeetrecommendations", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		err = deleteDefaultLifecycle(newSpaceClient)
		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, &config.OctolintConfig{}, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute(2)

		// Assert
		if result.Severity() != checks.Ok {
			return errors.New("Check should have passed")
		}

		return nil
	})
}

func TestLifecycleKeepsReleasesForever(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "17-lifecyclekeepsreleasesforever", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		err = deleteDefaultLifecycle(newSpaceClient)
		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, &config.OctolintConfig{}, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute(2)

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have produced a warning")
		}

		return nil
	})
}

func TestLifecycleKeepsFilesForever(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "18-lifecyclekeepsfilesforever", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		err = deleteDefaultLifecycle(newSpaceClient)
		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, &config.OctolintConfig{}, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute(2)

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have produced a warning")
		}

		return nil
	})
}

func TestLifecyclePhaseKeepsReleasesForever(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "19-lifecyclephasekeepsreleasesforever", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		err = deleteDefaultLifecycle(newSpaceClient)
		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, &config.OctolintConfig{}, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute(2)

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have produced a warning")
		}

		return nil
	})
}

func TestLifecyclePhaseKeepsFilesForever(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, filepath.Join("..", "..", "..", "test", "terraform"), "20-lifecyclephasekeepsfilesforever", []string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		err = deleteDefaultLifecycle(newSpaceClient)
		if err != nil {
			return err
		}

		check := NewOctopusLifecycleRetentionPolicyCheck(newSpaceClient, &config.OctolintConfig{}, checks.OctopusClientPermissiveErrorHandler{})

		result, err := check.Execute(2)

		if err != nil {
			return err
		}

		// Assert
		if result.Severity() != checks.Warning {
			return errors.New("Check should have produced a warning")
		}

		return nil
	})
}
