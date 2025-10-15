resource "octopusdeploy_lifecycle" "example" {
  description = "This is the default lifecycle."
  name        = "Test Lifecycle (OK to Delete)"

  release_retention_policy {
    quantity_to_keep    = 30
    should_keep_forever = false
    unit                = "Items"
  }

  tentacle_retention_policy {
    quantity_to_keep    = 30
    should_keep_forever = false
    unit                = "Days"
  }

    phase {
    optional_deployment_targets = [octopusdeploy_environment.development_environment.id]
    name                        = "Development"

    release_retention_policy {
      quantity_to_keep    = 30
      should_keep_forever = false
      unit                = "Items"
    }

    tentacle_retention_policy {
      quantity_to_keep    = 30
      should_keep_forever = false
      unit                = "Days"
    }
  }

  phase {
    optional_deployment_targets = [octopusdeploy_environment.test_environment.id]
    name                        = "Test"

    release_retention_policy {
      quantity_to_keep    = 30
      should_keep_forever = false
      unit                = "Items"
    }

    tentacle_retention_policy {
      quantity_to_keep    = 30
      should_keep_forever = false
      unit                = "Days"
    }
  }

  phase {
    optional_deployment_targets = [octopusdeploy_environment.production_environment.id]
    name                        = "Production"

    release_retention_policy {
      quantity_to_keep    = 30
      should_keep_forever = false
      unit                = "Items"
    }

    tentacle_retention_policy {
      quantity_to_keep    = 0
      should_keep_forever = true
      unit                = "Days"
    }
  }

  depends_on  = [octopusdeploy_environment.development_environment, octopusdeploy_environment.test_environment, octopusdeploy_environment.production_environment]
}