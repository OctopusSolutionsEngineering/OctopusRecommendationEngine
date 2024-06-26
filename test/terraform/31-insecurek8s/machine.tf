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

resource octopusdeploy_kubernetes_cluster_deployment_target test_eks {
  cluster_url                       = "https://cluster"
  environments                      = ["${octopusdeploy_environment.test_environment.id}"]
  name                              = "EKS Cluster"
  roles                             = ["eks"]
  cluster_certificate               = ""
  machine_policy_id                 = "${data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id}"
  namespace                         = ""
  skip_tls_verification             = true
  tenant_tags                       = []
  tenanted_deployment_participation = "Untenanted"
  tenants                           = []
  thumbprint                        = ""
  uri                               = ""
  default_worker_pool_id            = data.octopusdeploy_worker_pools.workerpool_default.worker_pools[0].id

  container {
    feed_id = ""
    image   = ""
  }

  aws_account_authentication {
    account_id        = "${octopusdeploy_aws_account.account_aws_account.id}"
    cluster_name      = "clustername"
    assume_role       = false
    use_instance_role = false
  }
}