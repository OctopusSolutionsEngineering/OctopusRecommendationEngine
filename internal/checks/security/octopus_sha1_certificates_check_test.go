package security

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
)

func TestSha1Certificates(t *testing.T) {
	testFramework := test.OctopusContainerTest{}

	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act: Deploy Terraform scenario that sets up SHA1 certificates
		newSpaceId, err := testFramework.Act(
			t,
			container,
			filepath.Join("..", "..", "..", "test", "terraform"),
			"33-sha1certificates", // folder containing your Terraform scenario
			[]string{},
		)
		if err != nil {
			return err
		}

		// Create a client for the new space
		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		if err != nil {
			return err
		}

		// Create the check
		check := NewOctopusSha1CertificatesCheck(
			newSpaceClient,
			&config.OctolintConfig{
				Url:                         container.URI,
				ApiKey:                      test.ApiKey,
				MaxSha1CertificatesMachines: 100,
			},
			checks.OctopusClientPermissiveErrorHandler{},
		)

		// Execute the check
		result, err := check.Execute(2)
		if err != nil {
			return err
		}

		// Assert
		if result == nil || result.Severity() != checks.Ok {
			return errors.New("check should have passed")
		}

		return nil
	})
}
