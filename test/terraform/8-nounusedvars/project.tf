data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "Default Lifecycle"
  skip         = 0
  take         = 1
}

data "octopusdeploy_project_groups" "default_project_group" {
  ids          = null
  partial_name = "Default Project Group"
  skip         = 0
  take         = 1
}

data "octopusdeploy_worker_pools" "workerpool_default" {
  name = "Default Worker Pool"
  ids  = null
  skip = 0
  take = 1
}

data "octopusdeploy_feeds" "built_in_feed" {
  feed_type    = "BuiltIn"
  ids          = null
  partial_name = ""
  skip         = 0
  take         = 1
}

resource "octopusdeploy_project" "project_with_no_process" {
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Test project with no process"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Test no process"
  project_group_id                     = data.octopusdeploy_project_groups.default_project_group.project_groups[0].id
  tenanted_deployment_participation    = "Untenanted"
  space_id                             = var.octopus_space_id
  included_library_variable_sets       = []
  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
  }

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
}

resource "octopusdeploy_runbook" "runbook_no_process" {
  project_id         = octopusdeploy_project.project_with_no_process.id
  name               = "Runbook"
  description        = "Test Runbook"
  multi_tenancy_mode = "Untenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  retention_policy {
    quantity_to_keep = 10
  }
  environment_scope           = "Specified"
  environments                = []
  default_guided_failure_mode = "EnvironmentDefault"
  force_package_download      = true
}

resource "octopusdeploy_project" "deploy_frontend_project" {
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Test project"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Test"
  project_group_id                     = data.octopusdeploy_project_groups.default_project_group.project_groups[0].id
  tenanted_deployment_participation    = "Untenanted"
  space_id                             = var.octopus_space_id
  included_library_variable_sets       = []
  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
  }

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
}

resource "octopusdeploy_variable" "jsonvariable" {
  owner_id     = "${octopusdeploy_project.deploy_frontend_project.id}"
  value        = "Whatever"
  name         = "propert1:property2:0"
  type         = "String"
  description  = ""
  is_sensitive = false
  depends_on = []
}

resource "octopusdeploy_variable" "runbookvar" {
  owner_id     = "${octopusdeploy_project.deploy_frontend_project.id}"
  value        = "Whatever"
  name         = "Runbook.Variable"
  type         = "String"
  description  = ""
  is_sensitive = false
  depends_on = []
}

resource "octopusdeploy_variable" "variablea" {
  owner_id     = "${octopusdeploy_project.deploy_frontend_project.id}"
  value        = "Whatever"
  name         = "VariableA"
  type         = "String"
  description  = ""
  is_sensitive = false
  depends_on = []
}

resource "octopusdeploy_variable" "variableb" {
  owner_id     = "${octopusdeploy_project.deploy_frontend_project.id}"
  value        = "Whatever"
  name         = "VariableB"
  type         = "String"
  description  = ""
  is_sensitive = false
  depends_on = []
}

resource "octopusdeploy_deployment_process" "deployment_process_project_api_gateway" {
  project_id = "${octopusdeploy_project.deploy_frontend_project.id}"

  step {
    condition           = "Success"
    name                = "Deploy Node.js app"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Deploy Node.js app"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = true
      is_required                        = false
      worker_pool_id                     = "${data.octopusdeploy_worker_pools.workerpool_default.worker_pools[0].id}"
      properties                         = {
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.Syntax" = "Bash"
        "Octopus.Action.Script.ScriptBody" = "echo \"#{VariableA} #{VariableB}\""
      }

      container {
        feed_id = ""
        image   = ""
      }

      environments          = []
      excluded_environments = []
      channels              = []
      tenant_tags           = []

      package {
        name                      = "RandomQuotes-JS"
        package_id                = "RandomQuotes-JS"
        acquisition_location      = "Server"
        extract_during_deployment = false
        feed_id                   = "${data.octopusdeploy_feeds.built_in_feed.feeds[0].id}"
        id                        = "ae4c9205-abb4-4a48-9252-bad99ec692d6"
        properties                = { Extract = "True", SelectionMode = "immediate" }
      }
      features = []
    }

    properties   = {}
    target_roles = []
  }
}

resource "octopusdeploy_runbook" "runbook" {
  project_id         = octopusdeploy_project.deploy_frontend_project.id
  name               = "Runbook"
  description        = "Test Runbook"
  multi_tenancy_mode = "Untenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  retention_policy {
    quantity_to_keep = 10
  }
  environment_scope           = "Specified"
  environments                = []
  default_guided_failure_mode = "EnvironmentDefault"
  force_package_download      = true
}

resource "octopusdeploy_runbook_process" "runbook" {
  runbook_id = octopusdeploy_runbook.runbook.id



  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Hello world (using PowerShell)"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = true
      worker_pool_id                     = ""
      properties                         = {
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.ScriptBody"   = "Write-Host 'Hello world, using PowerShell #{Runbook.Variable}'\n\n#TODO: Experiment with steps of your own :)\n\nWrite-Host '[Learn more about the types of steps available in Octopus](https://oc.to/OnboardingAddStepsLearnMore)'"
        "Octopus.Action.Script.Syntax"       = "PowerShell"
      }
      environments          = []
      excluded_environments = []
      channels              = []
      tenant_tags           = []
      features              = []

      package {
        name                      = "package1"
        package_id                = "package1"
        acquisition_location      = "Server"
        extract_during_deployment = false
        feed_id                   = data.octopusdeploy_feeds.built_in_feed.feeds[0].id
        properties                = { Extract = "True", Purpose = "", SelectionMode = "immediate" }
      }
    }

    properties   = {}
    target_roles = []
  }

  step {
    condition           = "Success"
    name                = "Test"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.TerraformApply"
      name                               = "Test"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = false
      worker_pool_id                     = ""
      worker_pool_variable               = ""
      properties                         = {
        "Octopus.Action.Terraform.PlanJsonOutput"               = "False"
        "Octopus.Action.Terraform.TemplateDirectory"            = "blah"
        "Octopus.Action.Terraform.ManagedAccount"               = "None"
        "Octopus.Action.Terraform.GoogleCloudAccount"           = "False"
        "Octopus.Action.Script.ScriptSource"                    = "Package"
        "Octopus.Action.GoogleCloud.UseVMServiceAccount"        = "True"
        "Octopus.Action.Terraform.AzureAccount"                 = "False"
        "Octopus.Action.Package.DownloadOnTentacle"             = "False"
        "Octopus.Action.Terraform.RunAutomaticFileSubstitution" = "True"
        "Octopus.Action.GoogleCloud.ImpersonateServiceAccount"  = "False"
        "Octopus.Action.Terraform.AllowPluginDownloads"         = "True"
      }
      environments          = []
      excluded_environments = []
      channels              = []
      tenant_tags           = []
      features              = []

      primary_package {
        package_id           = "terraform"
        acquisition_location = "Server"
        feed_id              = data.octopusdeploy_feeds.built_in_feed.feeds[0].id
        properties           = { SelectionMode = "immediate" }
      }
    }

    properties   = {}
    target_roles = []
  }
}