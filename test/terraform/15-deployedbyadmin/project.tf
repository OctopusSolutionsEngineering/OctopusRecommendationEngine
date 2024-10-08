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


resource "octopusdeploy_project" "echo_project" {
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

resource "octopusdeploy_variable" "variablea" {
  owner_id     = "${octopusdeploy_project.echo_project.id}"
  value        = "Whatever"
  name         = "VariableA"
  type         = "String"
  description  = ""
  is_sensitive = false
  depends_on = []
}

resource "octopusdeploy_variable" "variableb" {
  owner_id     = "${octopusdeploy_project.echo_project.id}"
  value        = "Whatever"
  name         = "VariableB"
  type         = "String"
  description  = ""
  is_sensitive = false
  depends_on = []
}

resource "octopusdeploy_deployment_process" "echo_project_process" {
  project_id = "${octopusdeploy_project.echo_project.id}"

  step {
    condition           = "Success"
    name                = "Echo"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Echo"
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
      features = []
    }

    properties   = {}
    target_roles = []
  }
}

output "project_id" {
  value = octopusdeploy_project.echo_project.id
}