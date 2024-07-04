package client_wrapper

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
)

func GetMachines(limit int, client newclient.Client, spaceID string) ([]*machines.DeploymentTarget, error) {
	if limit == 0 {
		return machines.GetAll(client, spaceID)
	}

	result, err := machines.Get(client, spaceID, machines.MachinesQuery{
		Take: limit,
	})

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}
