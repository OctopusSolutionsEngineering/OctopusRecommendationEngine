package client_wrapper

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workers"
)

func GetWorkers(limit int, client newclient.Client, spaceID string) ([]*machines.Worker, error) {
	if limit == 0 {
		return workers.GetAll(client, spaceID)
	}

	result, err := workers.Get(client, spaceID, machines.WorkersQuery{
		Take: limit,
	})

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}
