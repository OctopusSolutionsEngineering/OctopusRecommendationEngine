resource "octopusdeploy_lifecycle" "example" {
  description = "This is the default lifecycle."
  name        = "Test Lifecycle (OK to Delete)"

  release_retention_policy {
    quantity_to_keep    = 0
    should_keep_forever = true
    unit                = "Items"
  }

  tentacle_retention_policy {
    quantity_to_keep    = 30
    should_keep_forever = false
    unit                = "Items"
  }
}