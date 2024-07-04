package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/factory"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/naming"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/organization"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/performance"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/security"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/defaults"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/executor"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/reporters"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/briandowns/spinner"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var Version = "development"

func main() {
	octolintConfig, err := parseArgs()

	if err != nil {
		errorExit(err.Error())
		return
	}

	zap.ReplaceGlobals(createLogger(octolintConfig.Verbose))

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	if octolintConfig.Spinner && !octolintConfig.Verbose {
		s.Start()
	}

	defer func() {
		if octolintConfig.Spinner && !octolintConfig.Verbose {
			s.Stop()
		}
	}()

	if octolintConfig.Version {
		fmt.Println("Version: " + Version)
		os.Exit(0)
	}

	if octolintConfig.Url == "" {
		errorExit("You must specify the URL with the -url argument")
	}

	if octolintConfig.ApiKey == "" {
		errorExit("You must specify the API key with the -apiKey argument")
	}

	if octolintConfig.Space == "" {
		errorExit("You must specify the space key with the -space argument")
	}

	if !strings.HasPrefix(octolintConfig.Space, "Spaces-") {
		spaceId, err := lookupSpaceAsName(octolintConfig.Url, octolintConfig.Space, octolintConfig.ApiKey)

		if err != nil {
			errorExit("Failed to create the Octopus client_wrapper. Check that the url, api key, and space are correct.\nThe error was: " + err.Error())
		}

		octolintConfig.Space = spaceId
	}

	client, err := octoclient.CreateClient(octolintConfig.Url, octolintConfig.Space, octolintConfig.ApiKey)

	if err != nil {
		errorExit("Failed to create the Octopus client_wrapper. Check that the url, api key, and space are correct.\nThe error was: " + err.Error())
	}

	factory := factory.NewOctopusCheckFactory(client, octolintConfig.Url, octolintConfig.Space)
	checkCollection, err := factory.BuildAllChecks(octolintConfig)

	if err != nil {
		errorExit("Failed to create the checks")
	}

	// Time the execution
	startTime := time.Now().UnixMilli()
	defer func() {
		endTime := time.Now().UnixMilli()
		fmt.Println("Report took " + fmt.Sprint((endTime-startTime)/1000) + " seconds")
	}()

	executor := executor.NewOctopusCheckExecutor()
	results, err := executor.ExecuteChecks(checkCollection, func(check checks.OctopusCheck, err error) error {
		fmt.Fprintf(os.Stderr, "Failed to execute check "+check.Id())
		if octolintConfig.VerboseErrors {
			fmt.Println("##octopus[stdout-verbose]")
			fmt.Println(err.Error())
			fmt.Println("##octopus[stdout-default]")
		} else {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
		return nil
	})

	if err != nil {
		errorExit("Failed to run the checks")
	}

	reporter := reporters.NewOctopusPlainCheckReporter(checks.Warning)
	report, err := reporter.Generate(results)

	if err != nil {
		errorExit("Failed to generate the report")
	}

	fmt.Println(report)
}

func createLogger(verbose bool) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	level := zap.InfoLevel

	if verbose {
		level = zap.DebugLevel
	}

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(zapConfig.Build())
}

func errorExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func parseArgs() (*config.OctolintConfig, error) {
	config := config.OctolintConfig{}

	flag.StringVar(&config.Url, "url", "", "The Octopus URL e.g. https://myinstance.octopus.app")
	flag.StringVar(&config.Space, "space", "", "The Octopus space name or ID")
	flag.StringVar(&config.ApiKey, "apiKey", "", "The Octopus api key")
	flag.StringVar(&config.SkipTests, "skipTests", "", "A comma separated list of tests to skip")
	flag.StringVar(&config.OnlyTests, "onlyTests", "", "A comma separated list of tests to include")
	flag.StringVar(&config.ConfigFile, "configFile", "octolint", "The name of the configuration file to use. Do not include the extension. Defaults to octolint")
	flag.StringVar(&config.ConfigPath, "configPath", ".", "The path of the configuration file to use. Defaults to the current directory")
	flag.BoolVar(&config.Verbose, "verbose", false, "Print verbose logs")
	flag.BoolVar(&config.VerboseErrors, "verboseErrors", false, "Print error details as verbose logs in Octopus")
	flag.BoolVar(&config.Version, "version", false, "Print the version")
	flag.BoolVar(&config.Spinner, "spinner", true, "Display the spinner")
	flag.IntVar(&config.MaxEnvironments, "maxEnvironments", defaults.MaxEnvironments, "Maximum number of environments for the "+organization.OctopusEnvironmentCountCheckName+" check")
	flag.IntVar(&config.MaxDaysSinceLastTask, "maxDaysSinceLastTask", defaults.MaxTimeSinceLastTask, "Maximum number of days since the last project task for the "+organization.OctopusUnusedProjectsCheckName+" check")
	flag.IntVar(&config.MaxDuplicateVariables, "maxDuplicateVariables", defaults.MaxDuplicateVariables, "Maximum number of duplicate variables to report on for the "+organization.OctoLintDuplicatedVariables+" check. Set to 0 to report all duplicate variables.")
	flag.IntVar(&config.MaxDuplicateVariableProjects, "maxDuplicateVariableProjects", defaults.MaxDuplicateVariableProjects, "Maximum number of projects to check for duplicate variables for the "+organization.OctoLintDuplicatedVariables+" check. Set to 0 to check all projects.")
	flag.IntVar(&config.MaxDeploymentsByAdminProjects, "maxDeploymentsByAdminProjects", defaults.MaxDeploymentsByAdminProjects, "Maximum number of projects to check for admin deployments for the "+security.OctoLintDeploymentQueuedByAdmin+" check. Set to 0 to check all projects.")
	flag.IntVar(&config.MaxInvalidVariableProjects, "maxInvalidVariableProjects", defaults.MaxInvalidVariableProjects, "Maximum number of projects to check for invalid variables for the "+naming.OctoLintInvalidVariableNames+" check. Set to 0 to check all projects.")
	flag.IntVar(&config.MaxInvalidWorkerPoolProjects, "maxInvalidWorkerPoolProjects", defaults.MaxInvalidWorkerPoolProjects, "Maximum number of projects to check for invalid worker pools for the  "+naming.OctoLintProjectWorkerPool+" check. Set to 0 to check all projects.")
	flag.IntVar(&config.MaxInvalidContainerImageProjects, "maxInvalidContainerImageProjects", defaults.MaxInvalidContainerImageProjects, "Maximum number of projects to check for invalid container images for the "+naming.OctoLintContainerImageName+" check. Set to 0 to check all projects.")
	flag.IntVar(&config.MaxDefaultStepNameProjects, "maxDefaultStepNameProjects", defaults.MaxDefaultStepNameProjects, "Maximum number of projects to check for default step names for the "+naming.OctoLintProjectDefaultStepNames+" check. Set to 0 to report all projects")
	flag.IntVar(&config.MaxInvalidReleaseTemplateProjects, "maxInvalidReleaseTemplateProjects", defaults.MaxInvalidReleaseTemplateProjects, "Maximum number of projects to check for invalid release templates for the "+naming.OctoLintProjectReleaseTemplate+" check. Set to 0 to report all projects.")
	flag.IntVar(&config.MaxProjectSpecificEnvironmentProjects, "maxProjectSpecificEnvironmentProjects", defaults.MaxProjectSpecificEnvironmentProjects, "Maximum number of projects to check for project specific environments for the "+organization.OctoLintProjectSpecificEnvs+" check. Set to 0 to check all projects.")
	flag.IntVar(&config.MaxProjectSpecificEnvironmentEnvironments, "maxProjectSpecificEnvironmentEnvironments", defaults.MaxProjectSpecificEnvironmentEnvironments, "Maximum number of environments to check for project specific environments for the "+organization.OctoLintProjectSpecificEnvs+" check. Set to 0 to check all projects.")
	flag.IntVar(&config.MaxUnusedVariablesProjects, "maxUnusedVariablesProjects", defaults.MaxUnusedVariablesProjects, "Maximum number of projects to check for project specific environments for the "+organization.OctoLintUnusedVariables+" check. Set to 0 to report all projects for specific environments.")
	flag.IntVar(&config.MaxProjectStepsProjects, "maxProjectStepsProjects", defaults.MaxProjectStepsProjects, "Maximum number of projects to check for project step counts for the "+organization.OctoLintTooManySteps+" check. Set to 0 to report all projects for their step counts.")
	flag.IntVar(&config.MaxExclusiveEnvironmentsProjects, "maxExclusiveEnvironmentsProjects", defaults.MaxExclusiveEnvironmentsProjects, "Maximum number of projects to check for exclusive environments for the "+organization.OctoLintProjectGroupsWithExclusiveEnvironments+" check. Set to 0 to report all projects with exclusive environments.")
	flag.IntVar(&config.MaxEmptyProjectCheckProjects, "maxEmptyProjectCheckProjects", defaults.MaxEmptyProjectCheckProjects, "Maximum number of projects to check for no steps for the "+organization.OctoLintEmptyProject+" check. Set to 0 to report all empty projects.")
	flag.IntVar(&config.MaxUnusedProjects, "maxUnusedProjects", defaults.MaxUnusedProjects, "Maximum number of unused projects to check for the "+organization.OctopusUnusedProjectsCheckName+" check. Set to 0 to report all unused projects.")
	flag.IntVar(&config.MaxUnusedTargets, "maxUnusedTargets", defaults.MaxUnusedTargets, "Maximum number of unused targets to check for the "+organization.OctoLintUnusedTargets+" check. Set to 0 to report all unused targets.")
	flag.IntVar(&config.MaxUnhealthyTargets, "maxUnhealthyTargets", defaults.MaxUnhealthyTargets, "Maximum number of unhealthy targets to check for the "+organization.OctoLintUnhealthyTargets+" check. Set to 0 to report all unhealthy targets.")
	flag.IntVar(&config.MaxInvalidRoleTargets, "maxInvalidRoleTargets", defaults.MaxInvalidRoleTargets, "Maximum number of targets to check for invalid roles for the "+naming.OctoLintInvalidTargetRoles+" check. Set to 0 to report all targets.")
	flag.IntVar(&config.MaxTenantTagsTargets, "maxTenantTagsTargets", defaults.MaxTenantTagsTargets, "Maximum number of targets to check for potential tenant tags for the "+organization.OctoLintDirectTenantReferences+" check. Set to 0 to check all targets.")
	flag.IntVar(&config.MaxTenantTagsTenants, "maxTenantTagsTenants", defaults.MaxTenantTagsTenants, "Maximum number of tenants to check for potential tenant tags for the "+organization.OctoLintDirectTenantReferences+" check. Set to 0 to check all targets.")
	flag.IntVar(&config.MaxInvalidNameTargets, "maxInvalidNameTargets", defaults.MaxInvalidNameTargets, "Maximum number of targets to check for invalid names for the "+naming.OctoLintInvalidTargetNames+" check. Set to 0 to check all targets.")
	flag.IntVar(&config.MaxInsecureK8sTargets, "maxInsecureK8sTargets", defaults.MaxInsecureK8sTargets, "Maximum number of targets to check for insecure k8s configuration for the "+security.OctoLintInsecureK8sTargets+" check. Set to 0 to check all targets.")
	flag.IntVar(&config.MaxDeploymentTasks, "maxDeploymentTasks", defaults.MaxDeploymentTasks, "Maximum number of deployment tasks to scan for the "+performance.OctoLintDeploymentQueuedTime+" check. Set to 0 to check all targets.")
	flag.StringVar(&config.ContainerImageRegex, "containerImageRegex", "", "The regular expression used to validate container images for the "+naming.OctoLintContainerImageName+" check")
	flag.StringVar(&config.VariableNameRegex, "variableNameRegex", "", "The regular expression used to validate variable names for the "+naming.OctoLintInvalidVariableNames+" check")
	flag.StringVar(&config.TargetNameRegex, "targetNameRegex", "", "The regular expression used to validate target names for the "+naming.OctoLintInvalidTargetNames+" check")
	flag.StringVar(&config.TargetRoleRegex, "targetRoleRegex", "", "The regular expression used to validate target roles for the "+naming.OctoLintInvalidTargetRoles+" check")
	flag.StringVar(&config.ProjectReleaseTemplateRegex, "projectReleaseTemplateRegex", "", "The regular expression used to validate project release templates for the "+naming.OctoLintProjectReleaseTemplate+" check")
	flag.StringVar(&config.ProjectStepWorkerPoolRegex, "projectStepWorkerPoolRegex", "", "The regular expression used to validate step worker pools for the  "+naming.OctoLintProjectReleaseTemplate+" check")
	flag.StringVar(&config.LifecycleNameRegex, "lifecycleNameRegex", "", "The regular expression used to validate lifecycle names for the  "+naming.OctoLintInvalidLifecycleNames+" check")

	flag.Parse()

	err := overrideArgs(config.ConfigPath, config.ConfigFile)

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

func lookupSpaceAsName(octopusUrl string, spaceName string, apiKey string) (string, error) {
	if len(strings.TrimSpace(spaceName)) == 0 {
		return "", errors.New("space can not be empty")
	}

	requestURL := fmt.Sprintf("%s/api/Spaces?take=1000&partialName=%s", octopusUrl, url.QueryEscape(spaceName))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		return "", err
	}

	if apiKey != "" {
		req.Header.Set("X-Octopus-ApiKey", apiKey)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", nil
	}
	defer res.Body.Close()

	collection := resources.Resources[spaces.Space]{}
	err = json.NewDecoder(res.Body).Decode(&collection)

	if err != nil {
		return "", err
	}

	for _, space := range collection.Items {
		if space.Name == spaceName {
			return space.ID, nil
		}
	}

	return "", errors.New("did not find space with name " + spaceName)
}
