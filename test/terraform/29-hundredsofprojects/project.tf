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


resource "octopusdeploy_project" "deploy_frontend_project" {
  count = 100
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Test project ${count.index}"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Test ${count.index}"
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
  count = 100
  owner_id     = "${octopusdeploy_project.deploy_frontend_project[count.index].id}"
  value        = "Whatever"
  name         = "VariableA"
  type         = "String"
  description  = ""
  is_sensitive = false
  depends_on = []
}

resource "octopusdeploy_variable" "variableb" {
  count = 100
  owner_id     = "${octopusdeploy_project.deploy_frontend_project[count.index].id}"
  value        = "Whatever"
  name         = "VariableB"
  type         = "String"
  description  = ""
  is_sensitive = false
  depends_on = []
}

resource "octopusdeploy_deployment_process" "deployment_process_project_api_gateway" {
  count = 100
  project_id = "${octopusdeploy_project.deploy_frontend_project[count.index].id}"

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