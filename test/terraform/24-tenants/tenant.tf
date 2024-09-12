resource "octopusdeploy_tenant" "tenant_team_a" {
  name        = "Team A"
  description = "Test tenant"
  tenant_tags = []
  depends_on = []
}

resource "octopusdeploy_tenant_project" "tenant_project_link" {
  environment_ids = [octopusdeploy_environment.development_environment.id]
  project_id      = octopusdeploy_project.deploy_frontend_project.id
  tenant_id       = octopusdeploy_tenant.tenant_team_a.id
}

resource "octopusdeploy_tenant" "tenant_team_b" {
  name        = "Team B"
  description = "Test tenant"
  tenant_tags = []
  depends_on = []
}

resource "octopusdeploy_tenant_project" "tenant_project_link" {
  environment_ids = [octopusdeploy_environment.development_environment.id]
  project_id      = octopusdeploy_project.deploy_frontend_project.id
  tenant_id       = octopusdeploy_tenant.tenant_team_b.id
}
