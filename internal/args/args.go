package args

import (
	"bytes"
	"errors"
	"flag"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/naming"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/organization"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/performance"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/security"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/defaults"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func ParseArgs(args []string) (*config.OctolintConfig, error) {
	flags := flag.NewFlagSet("octolint", flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	config := config.OctolintConfig{}

	flags.StringVar(&config.Url, "url", "", "The Octopus URL e.g. https://myinstance.octopus.app")
	flags.StringVar(&config.Space, "space", "", "The Octopus space name or ID")
	flags.StringVar(&config.ApiKey, "apiKey", "", "The Octopus api key")
	flags.StringVar(&config.SkipTests, "skipTests", "", "A comma separated list of tests to skip")
	flags.StringVar(&config.OnlyTests, "onlyTests", "", "A comma separated list of tests to include")
	flags.StringVar(&config.ConfigFile, "configFile", "octolint", "The name of the configuration file to use. Do not include the extension. Defaults to octolint")
	flags.StringVar(&config.ConfigPath, "configPath", ".", "The path of the configuration file to use. Defaults to the current directory")
	flags.BoolVar(&config.Verbose, "verbose", false, "Print verbose logs")
	flags.BoolVar(&config.VerboseErrors, "verboseErrors", false, "Print error details as verbose logs in Octopus")
	flags.BoolVar(&config.Version, "version", false, "Print the version")
	flags.BoolVar(&config.Spinner, "spinner", true, "Display the spinner")
	flags.IntVar(&config.MaxEnvironments, "maxEnvironments", defaults.MaxEnvironments, "Maximum number of environments for the "+organization.OctopusEnvironmentCountCheckName+" check")
	flags.IntVar(&config.MaxDaysSinceLastTask, "maxDaysSinceLastTask", defaults.MaxTimeSinceLastTask, "Maximum number of days since the last project task for the "+organization.OctopusUnusedProjectsCheckName+" check")
	flags.IntVar(&config.MaxDuplicateVariables, "maxDuplicateVariables", defaults.MaxDuplicateVariables, "Maximum number of duplicate variables to report on for the "+organization.OctoLintDuplicatedVariables+" check. Set to 0 to report all duplicate variables.")
	flags.IntVar(&config.MaxDuplicateVariableProjects, "maxDuplicateVariableProjects", defaults.MaxDuplicateVariableProjects, "Maximum number of projects to check for duplicate variables for the "+organization.OctoLintDuplicatedVariables+" check. Set to 0 to check all projects.")
	flags.IntVar(&config.MaxDeploymentsByAdminProjects, "maxDeploymentsByAdminProjects", defaults.MaxDeploymentsByAdminProjects, "Maximum number of projects to check for admin deployments for the "+security.OctoLintDeploymentQueuedByAdmin+" check. Set to 0 to check all projects.")
	flags.IntVar(&config.MaxInvalidVariableProjects, "maxInvalidVariableProjects", defaults.MaxInvalidVariableProjects, "Maximum number of projects to check for invalid variables for the "+naming.OctoLintInvalidVariableNames+" check. Set to 0 to check all projects.")
	flags.IntVar(&config.MaxInvalidWorkerPoolProjects, "maxInvalidWorkerPoolProjects", defaults.MaxInvalidWorkerPoolProjects, "Maximum number of projects to check for invalid worker pools for the  "+naming.OctoLintProjectWorkerPool+" check. Set to 0 to check all projects.")
	flags.IntVar(&config.MaxInvalidContainerImageProjects, "maxInvalidContainerImageProjects", defaults.MaxInvalidContainerImageProjects, "Maximum number of projects to check for invalid container images for the "+naming.OctoLintContainerImageName+" check. Set to 0 to check all projects.")
	flags.IntVar(&config.MaxDefaultStepNameProjects, "maxDefaultStepNameProjects", defaults.MaxDefaultStepNameProjects, "Maximum number of projects to check for default step names for the "+naming.OctoLintProjectDefaultStepNames+" check. Set to 0 to report all projects")
	flags.IntVar(&config.MaxInvalidReleaseTemplateProjects, "maxInvalidReleaseTemplateProjects", defaults.MaxInvalidReleaseTemplateProjects, "Maximum number of projects to check for invalid release templates for the "+naming.OctoLintProjectReleaseTemplate+" check. Set to 0 to report all projects.")
	flags.IntVar(&config.MaxProjectSpecificEnvironmentProjects, "maxProjectSpecificEnvironmentProjects", defaults.MaxProjectSpecificEnvironmentProjects, "Maximum number of projects to check for project specific environments for the "+organization.OctoLintProjectSpecificEnvs+" check. Set to 0 to check all projects.")
	flags.IntVar(&config.MaxProjectSpecificEnvironmentEnvironments, "maxProjectSpecificEnvironmentEnvironments", defaults.MaxProjectSpecificEnvironmentEnvironments, "Maximum number of environments to check for project specific environments for the "+organization.OctoLintProjectSpecificEnvs+" check. Set to 0 to check all projects.")
	flags.IntVar(&config.MaxUnusedVariablesProjects, "maxUnusedVariablesProjects", defaults.MaxUnusedVariablesProjects, "Maximum number of projects to check for project specific environments for the "+organization.OctoLintUnusedVariables+" check. Set to 0 to report all projects for specific environments.")
	flags.IntVar(&config.MaxProjectStepsProjects, "maxProjectStepsProjects", defaults.MaxProjectStepsProjects, "Maximum number of projects to check for project step counts for the "+organization.OctoLintTooManySteps+" check. Set to 0 to report all projects for their step counts.")
	flags.IntVar(&config.MaxExclusiveEnvironmentsProjects, "maxExclusiveEnvironmentsProjects", defaults.MaxExclusiveEnvironmentsProjects, "Maximum number of projects to check for exclusive environments for the "+organization.OctoLintProjectGroupsWithExclusiveEnvironments+" check. Set to 0 to report all projects with exclusive environments.")
	flags.IntVar(&config.MaxEmptyProjectCheckProjects, "maxEmptyProjectCheckProjects", defaults.MaxEmptyProjectCheckProjects, "Maximum number of projects to check for no steps for the "+organization.OctoLintEmptyProject+" check. Set to 0 to report all empty projects.")
	flags.IntVar(&config.MaxUnusedProjects, "maxUnusedProjects", defaults.MaxUnusedProjects, "Maximum number of unused projects to check for the "+organization.OctopusUnusedProjectsCheckName+" check. Set to 0 to report all unused projects.")
	flags.IntVar(&config.MaxUnusedTargets, "maxUnusedTargets", defaults.MaxUnusedTargets, "Maximum number of unused targets to check for the "+organization.OctoLintUnusedTargets+" check. Set to 0 to report all unused targets.")
	flags.IntVar(&config.MaxUnhealthyTargets, "maxUnhealthyTargets", defaults.MaxUnhealthyTargets, "Maximum number of unhealthy targets to check for the "+organization.OctoLintUnhealthyTargets+" check. Set to 0 to report all unhealthy targets.")
	flags.IntVar(&config.MaxInvalidRoleTargets, "maxInvalidRoleTargets", defaults.MaxInvalidRoleTargets, "Maximum number of targets to check for invalid roles for the "+naming.OctoLintInvalidTargetRoles+" check. Set to 0 to report all targets.")
	flags.IntVar(&config.MaxTenantTagsTargets, "maxTenantTagsTargets", defaults.MaxTenantTagsTargets, "Maximum number of targets to check for potential tenant tags for the "+organization.OctoLintDirectTenantReferences+" check. Set to 0 to check all targets.")
	flags.IntVar(&config.MaxTenantTagsTenants, "maxTenantTagsTenants", defaults.MaxTenantTagsTenants, "Maximum number of tenants to check for potential tenant tags for the "+organization.OctoLintDirectTenantReferences+" check. Set to 0 to check all targets.")
	flags.IntVar(&config.MaxInvalidNameTargets, "maxInvalidNameTargets", defaults.MaxInvalidNameTargets, "Maximum number of targets to check for invalid names for the "+naming.OctoLintInvalidTargetNames+" check. Set to 0 to check all targets.")
	flags.IntVar(&config.MaxInsecureK8sTargets, "maxInsecureK8sTargets", defaults.MaxInsecureK8sTargets, "Maximum number of targets to check for insecure k8s configuration for the "+security.OctoLintInsecureK8sTargets+" check. Set to 0 to check all targets.")
	flags.IntVar(&config.MaxDeploymentTasks, "maxDeploymentTasks", defaults.MaxDeploymentTasks, "Maximum number of deployment tasks to scan for the "+performance.OctoLintDeploymentQueuedTime+" check. Set to 0 to check all targets.")
	flags.StringVar(&config.ContainerImageRegex, "containerImageRegex", "", "The regular expression used to validate container images for the "+naming.OctoLintContainerImageName+" check")
	flags.StringVar(&config.VariableNameRegex, "variableNameRegex", "", "The regular expression used to validate variable names for the "+naming.OctoLintInvalidVariableNames+" check")
	flags.StringVar(&config.TargetNameRegex, "targetNameRegex", "", "The regular expression used to validate target names for the "+naming.OctoLintInvalidTargetNames+" check")
	flags.StringVar(&config.TargetRoleRegex, "targetRoleRegex", "", "The regular expression used to validate target roles for the "+naming.OctoLintInvalidTargetRoles+" check")
	flags.StringVar(&config.ProjectReleaseTemplateRegex, "projectReleaseTemplateRegex", "", "The regular expression used to validate project release templates for the "+naming.OctoLintProjectReleaseTemplate+" check")
	flags.StringVar(&config.ProjectStepWorkerPoolRegex, "projectStepWorkerPoolRegex", "", "The regular expression used to validate step worker pools for the  "+naming.OctoLintProjectReleaseTemplate+" check")
	flags.StringVar(&config.LifecycleNameRegex, "lifecycleNameRegex", "", "The regular expression used to validate lifecycle names for the  "+naming.OctoLintInvalidLifecycleNames+" check")

	err := flags.Parse(args)

	if err != nil {
		return nil, err
	}

	err = overrideArgs(config.ConfigPath, config.ConfigFile)

	if err != nil {
		return nil, err
	}

	if config.Url == "" {
		config.Url = os.Getenv("OCTOPUS_CLI_SERVER")
	}

	if config.ApiKey == "" {
		config.ApiKey = os.Getenv("OCTOPUS_CLI_API_KEY")
	}

	return &config, nil
}

// Inspired by https://github.com/carolynvs/stingoftheviper
// Viper needs manual handling to implement reading settings from env vars, config files, and from the command line
func overrideArgs(configPath string, configFile string) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(configFile)

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(configPath)

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix("octolint")

	// Environment variables can't have dashes in them, so bind them to their equivalent
	// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	return bindFlags(v)
}

// Bind each flag to its associated viper configuration (config file and environment variable)
func bindFlags(v *viper.Viper) (funErr error) {
	var funcError error = nil

	flag.VisitAll(func(allFlags *flag.Flag) {
		defined := false
		flag.Visit(func(definedFlag *flag.Flag) {
			if definedFlag.Name == allFlags.Name && definedFlag.Name != "configFile" && definedFlag.Name != "configPath" {
				defined = true
			}
		})

		if !defined && v.IsSet(allFlags.Name) {
			configName := strings.ReplaceAll(allFlags.Name, "-", "")

			for _, value := range v.GetStringSlice(configName) {
				err := flag.Set(allFlags.Name, value)
				funcError = errors.Join(funcError, err)
			}
		}
	})

	return funcError
}
