The application provides a set of checks to validate the configuration of an Octopus instance.

The checks are in directories under the `checks` directory.

The checks in `checks/naming` validate the naming conventions of various resources in Octopus.

The checks in `checks/organization` validate the relationships between various resources in Octopus.

The checks in `checks/performance` validate the performance of the Octopus instance.

The checks in `checks/security` validate the security of the Octopus instance.

The application exposes a CLI interface via the `main` package in the `cmd/cli` directory.

The application also exposes an REST API interface via the `main` package in the `cmd/azure` directory.

The REST API uses the JSON API standard. The web server is configured by the functions in the `environment` package.