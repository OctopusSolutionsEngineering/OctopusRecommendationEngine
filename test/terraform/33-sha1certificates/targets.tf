data "octopusdeploy_machine_policies" "default_machine_policy" {
  ids          = null
  partial_name = "Default Machine Policy"
  skip         = 0
  take         = 1
}

data "octopusdeploy_worker_pools" "workerpool_default" {
  name = "Default Worker Pool"
  ids  = null
  skip = 0
  take = 1
}

resource "octopusdeploy_listening_tentacle_deployment_target" "target_example" {
  environments                      = ["${octopusdeploy_environment.development_environment.id}"]
  is_disabled                       = true
  machine_policy_id                 = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
  name                              = "sha1 certificate test target"
  roles                             = ["sha1cert-app"]
  tenanted_deployment_participation = "Untenanted"
  tentacle_url                      = "https://target-example.com:1234/"
  thumbprint                        = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
}

resource "octopusdeploy_listening_tentacle_worker" "worker_example" {
  name         = "sha1 listening_worker"
  machine_policy_id = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
  worker_pool_ids = [ "${data.octopusdeploy_worker_pools.workerpool_default.worker_pools[0].id}" ]
  thumbprint   = "96203ED84246201C26A2F4360D7CBC36AC1D232C"
  uri          = "https://worker-example.com:1234/"
}