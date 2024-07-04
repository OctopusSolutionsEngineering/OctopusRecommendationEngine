package client_wrapper

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
)

func GetEnvironments(limit int, client newclient.Client, spaceID string) ([]*environments.Environment, error) {
	if limit == 0 {
		return environments.GetAll(client, spaceID)
	}

	result, err := environments.Get(client, spaceID, environments.EnvironmentsQuery{
		Take: limit,
	})

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}
