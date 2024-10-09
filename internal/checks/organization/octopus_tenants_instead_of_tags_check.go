package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"strings"
)

const OctoLintDirectTenantReferences = "OctoLintDirectTenantReferences"

// OctopusTenantsInsteadOfTagsCheck checks to see if any common groups of tenants are found against common resources like accounts, targets etc
type OctopusTenantsInsteadOfTagsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusTenantsInsteadOfTagsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusTenantsInsteadOfTagsCheck {
	return OctopusTenantsInsteadOfTagsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusTenantsInsteadOfTagsCheck) Id() string {
	return OctoLintDirectTenantReferences
}

func (o OctopusTenantsInsteadOfTagsCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	allTenants, err := client_wrapper.GetTenants(o.config.MaxTenantTagsTenants, o.client, o.client.GetSpaceID())

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allAccounts, err := o.client.Accounts.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allCertificates, err := o.client.Certificates.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	allMachines, err := client_wrapper.GetMachines(o.config.MaxTenantTagsTargets, o.client, o.client.GetSpaceID())

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	tenantReferenceCounts := map[string]int{}
	tenantReferenceSources := map[string][]string{}
	for i, a := range allAccounts {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allAccounts))*33) + "% complete")

		if a.GetTenantedDeploymentMode() == core.TenantedDeploymentModeTenantedOrUntenanted {
			o.addTenants(a.GetTenantIDs(), "Account - "+a.GetName(), tenantReferenceCounts, tenantReferenceSources)
		}
	}

	for i, c := range allCertificates {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allCertificates))*33) + "% complete")

		if c.TenantedDeploymentMode == core.TenantedDeploymentModeTenantedOrUntenanted {
			o.addTenants(c.TenantIDs, "Certificate - "+c.Name, tenantReferenceCounts, tenantReferenceSources)
		}
	}

	for i, m := range allMachines {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allMachines))*33) + "% complete")

		if m.TenantedDeploymentMode == core.TenantedDeploymentModeTenantedOrUntenanted {
			o.addTenants(m.TenantIDs, "Target - "+m.Name, tenantReferenceCounts, tenantReferenceSources)
		}
	}

	// get any commonly grouped tenants
	multipleTenantReferences := []string{}
	for tenantGroups, groupCount := range tenantReferenceCounts {
		if groupCount > 1 {
			multipleTenantReferences = append(multipleTenantReferences, tenantGroups)
		}
	}

	if len(multipleTenantReferences) > 0 {

		// We have to convert the comma separated list of tenant IDs into a comma separated list of tenant names
		groupedTenants := []string{}
		for _, groupedTenant := range multipleTenantReferences {
			splitTenants := strings.Split(groupedTenant, ",")
			splitTenantNames := []string{}
			for _, splitTenant := range splitTenants {
				splitTenantNames = append(splitTenantNames, o.getTenantNameById(allTenants, splitTenant))
			}
			groupedTenants = append(groupedTenants, strings.Join(splitTenantNames, ", ")+" ("+strings.Join(tenantReferenceSources[groupedTenant], ", ")+")")
		}

		return checks.NewOctopusCheckResultImpl(
			"The following groups of tenants have been directly referenced more than once, and may be better grouped as tenant tags:\n"+strings.Join(groupedTenants, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"No duplicate groups of tenants were found",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusTenantsInsteadOfTagsCheck) getTenantNameById(tenants []*tenants.Tenant, id string) string {
	for _, l := range tenants {
		if l.ID == id {
			return l.Name
		}
	}

	return ""
}

func (o OctopusTenantsInsteadOfTagsCheck) addTenants(tenantIds []string, source string, tenantReferences map[string]int, tenantReferenceSources map[string][]string) {
	if len(tenantIds) <= 1 {
		return
	}

	slices.Sort(tenantIds)
	tenants := strings.Join(tenantIds, ",")

	if _, ok := tenantReferences[tenants]; !ok {
		tenantReferences[tenants] = 0
	}
	tenantReferences[tenants]++

	if _, ok := tenantReferenceSources[tenants]; !ok {
		tenantReferenceSources[tenants] = []string{}
	}
	tenantReferenceSources[tenants] = append(tenantReferenceSources[tenants], source)
}
