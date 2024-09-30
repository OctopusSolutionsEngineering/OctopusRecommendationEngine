# Octolint

<img src="https://user-images.githubusercontent.com/160104/222631936-e1ec480e-abd5-4622-978d-08259844aa14.png" width="100" height="100">

[![Github All Releases](https://img.shields.io/github/downloads/OctopusSolutionsEngineering/OctopusRecommendationEngine/total.svg)]()

This CLI tool scans an Octopus instance to find potential issues in the configuration and suggests solutions.

## Support

Feel free to report
an [issue](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/issues).

This tool is not officially supported by Octopus. Please do not contact the Octopus support channels regarding octolint.

## Usage

Download the latest binary from
the [releases](https://github.com/OctopusSalesEngineering/OctopusRecommendationEngine/releases).

```
./octolint \
    -apiKey API-YOURAPIKEY \
    -url https://yourinstance.octopus.app \
    -space Spaces-1234
```

Octolint is also distributed as a Docker image:

```
docker run -t --rm ghcr.io/octopussolutionsengineering/octolint \
    -url https://yourinstance.octopus.app \
    -apiKey API-YOURAPIKEY \
    -space Spaces-1
```

You can run octolint as part of an Octopus deployment process or runbook:

```bash
echo "##octopus[stdout-verbose]"
docker pull ghcr.io/octopussolutionsengineering/octolint 2>&1
echo "##octopus[stdout-default]"

docker run -t --rm \
    ghcr.io/octopussolutionsengineering/octolint \
    -spinner=false \
    -url #{Octopus.Web.ServerUri} \
    -apiKey #{ApiKey} \
    -space #{Octopus.Space.Id}
```

## Configuration files and environment variables

All program arguments can be defined as environment variables with the prefix `OCTOLINT_` or in a YAML file called
`octolint.yaml` saved alongside the executable.

For example, the maximum number of environments can be passed as a command line argument:

```bash
./octolint -maxEnvironments 5
```

Or defined in an environment variable (environment variables are case insensitive):

```bash
OCTOLINT_MAXENVIRONMENTS=5 ./octolint
```

Or defined in a file called `octolint.yaml`:

```yaml
maxEnvironments: 5
```

The order of precedence from lowest to highest is:
* Default values
* Config file
* Environment variable
* Command line arguments

## Default resource limits

Octolint will scan 100 projects and targets by default. This prevents the scans from taking too long in large Octopus spaces.
However, it also means that some issues may not be detected if they are in projects or targets that are not scanned.

The arguments starting with `max...`, like `maxDuplicateVariableProjects` or `maxUnhealthyTargets`, can be set to 0 to scan all projects
or targets, or set to a number larger than 0 to scan a custom number of projects or targets.

Run `octolint -h` to see all the available arguments.

## Capturing output in Octopus

The easiest way to capture the output of Octolint in Octopus is to capture the standard output in a variable and use the variable
to create an [output variable](https://octopus.com/docs/projects/variables/output-variables).

The example below shows how to achieve this in Bash:

```bash
echo "##octopus[stdout-verbose]"
docker pull ghcr.io/octopussolutionsengineering/octolint 2>&1
echo "##octopus[stdout-default]"

RESULTS=$(docker run -t --rm \
    ghcr.io/octopussolutionsengineering/octolint \
    -spinner=false \
    -url "#{Octopus.Web.ServerUri}" \
    -apiKey "#{ApiKey}" \
    -space "#{Octopus.Space.Id}")
    
set_octopusvariable "OctolintResults" "$RESULTS"
echo $RESULTS
```

## Permissions

`octolint` only requires read access - it does not modify anything on the server.

To create a read only account, deploy the Terraform module under the [serviceaccount](serviceaccount) directory:

```bash
export TF_VAR_octopus_server=https://yourinstance.octopus.app
export TF_VAR_octopus_apikey=API-apikeygoeshere
export TF_VAR_octopus_space_id=Spaces-#
cd serviceaccount
terraform init
terraform apply
```

This creates a user role, team, and service account all called `Octolint`. You can then create an API key for the service account, and use that API key with `octolint`. 

## Example output

This is an example of the tool output:

```
[OctoLintDefaultProjectGroupChildCount] The default project group contains 79 projects. You may want to organize these projects into additional project groups.
[OctoLintEmptyProject] The following projects have no runbooks and no deployment process: Azure Octopus test, CloudFormation, K8s Yaml Import, AA Training Demo, Test2, Cac Vars, Helm Demo, Helm, Package, K8s
[OctoLintUnusedVariables] The following variables are unused: App Runner/x, App Runner/thisisnotused, Vars Demo/workerpool, GAE Node.js/Email Address, CloudFormation - Lambda/APIKey, ReleaseDiffTest/Variable2, ReleaseDiffTest/Variable2, ReleaseDiffTest/Variable1, Release Diff/packages[package].files[file].diff, Release Diff/scoped variable, Release Diff/scoped variable, Release Diff/unscoped variable, Var test/Config[Databases:Name].Value, Var test/Config[General:DEBUG].Value, K8s Command Example/MyTags[Three].Name, Rolling Deployment/b, Rolling Deployment/aws, Rolling Deployment/azure, Rolling Deployment/workerpool, Rolling Deployment/gcp, Rolling Deployment/cert, CloudFormation - RDS/APIKey, Devops Tasks/APIKey, Devops Tasks/AWS, Terraform Test/DockerHub.Password, Terraform Test/New Value, Terraform Test/Scoped value, CloudFormation - Lambda Simple/APIKey, AWS Account/AWS
[OctoLintDuplicatedVariables] The following variables are duplicated between projects. Consider moving these into library variable sets: Cart/MONGO_DB_HOSTNAME == Orders/MONGO_DB_HOSTNAME, Cart/MONGO_DB_HOSTNAME == User/MONGO_DB_HOSTNAME, Orders/MONGO_DB_HOSTNAME == User/MONGO_DB_HOSTNAME
[OctoLintTooManySteps] The following projects have 20 or more steps: K8s Yaml Import 2
```

## Checks

Refer to the [wiki](https://github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/wiki) for a list of checks. 

## Debugging network issues in docker

If you get an error saying the client could not be created, and you are running octolint from a Docker container, check
that the Octopus server can be resolved from the container with the following command. This command overrides the
container entry point to run `nslookup` and passes the name of the server as the final argument. The example below
attempts to resolve `yourinstance.octopus.app`, and you must change this to reflect the hostname of the `-url` argument
you would normally pass to octolint:

```shell
docker run --rm --entrypoint "/usr/bin/nslookup" ghcr.io/octopussolutionsengineering/octolint yourinstance.octopus.app
```

This is an example of the output when the Docker container can not resolve the network address. This either indicates that
the Docker networking is not allowing the hostname to be resolved, or that the hostname is invalid:

```shell
$ docker run --rm --entrypoint "/usr/bin/nslookup" ghcr.io/octopussolutionsengineering/octolint this.address.does.not.exist
Server:		1.1.1.1
Address:	1.1.1.1:53

** server can't find this.address.does.not.exist: NXDOMAIN

** server can't find this.address.does.not.exist: NXDOMAIN
```

This is an example were the hostname can be successfully resolved:

```shell
$ docker run --rm --entrypoint "/usr/bin/nslookup" ghcr.io/octopussolutionsengineering/octolint mattc.octopus.app
Server:		1.1.1.1
Address:	1.1.1.1:53

Non-authoritative answer:

Non-authoritative answer:
Name:	mattc.octopus.app
Address: 20.53.101.130
```