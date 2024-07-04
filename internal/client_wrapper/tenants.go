package client_wrapper

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
)

func GetTenants(limit int, client newclient.Client, spaceID string) ([]*tenants.Tenant, error) {
	if limit == 0 {
		return tenants.GetAll(client, spaceID)
	}

	result, err := tenants.Get(client, spaceID, tenants.TenantsQuery{
		Take: limit,
	})

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}
