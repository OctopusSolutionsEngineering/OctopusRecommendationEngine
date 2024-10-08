package config

import (
	"strings"
)

type OctolintConfig struct {
	Url           string
	Space         string
	ApiKey        string
	SkipTests     string
	OnlyTests     string
	VerboseErrors bool
	Version       bool
	Spinner       bool
	ConfigFile    string
	ConfigPath    string
	Verbose       bool

	// Global filters for resources
	ExcludeProjects       StringSliceArgs
	ExcludeProjectsExcept StringSliceArgs
	ExcludeProjectsRegex  StringSliceArgs

	// These values are used to configure individual checks
	MaxEnvironments                           int
	ContainerImageRegex                       string
	VariableNameRegex                         string
	TargetNameRegex                           string
	WorkerNameRegex                           string
	WorkerPoolNameRegex                       string
	TargetRoleRegex                           string
	ProjectReleaseTemplateRegex               string
	ProjectStepWorkerPoolRegex                string
	SpaceNameRegex                            string
	LibraryVariableSetNameRegex               string
	TenantNameRegex                           string
	TagSetNameRegex                           string
	TagNameRegex                              string
	FeedNameRegex                             string
	AccountNameRegex                          string
	MachinePolicyNameRegex                    string
	CertificateNameRegex                      string
	GitCredentialNameRegex                    string
	ScriptModuleNameRegex                     string
	ProjectGroupNameRegex                     string
	ProjectNameRegex                          string
	LifecycleNameRegex                        string
	MaxDaysSinceLastTask                      int
	MaxDuplicateVariables                     int
	MaxDuplicateVariableProjects              int
	MaxInvalidVariableProjects                int
	MaxInvalidReleaseTemplateProjects         int
	MaxInvalidContainerImageProjects          int
	MaxInvalidWorkerPoolProjects              int
	MaxEmptyProjectCheckProjects              int
	MaxExclusiveEnvironmentsProjects          int
	MaxProjectSpecificEnvironmentProjects     int
	MaxProjectSpecificEnvironmentEnvironments int
	MaxProjectStepsProjects                   int
	MaxUnusedVariablesProjects                int
	MaxUnusedProjects                         int
	MaxDefaultStepNameProjects                int
	MaxDeploymentsByAdminProjects             int
	MaxUnusedTargets                          int
	MaxUnhealthyTargets                       int
	MaxTenantTagsTargets                      int
	MaxTenantTagsTenants                      int
	MaxInvalidRoleTargets                     int
	MaxInvalidNameTargets                     int
	MaxInsecureK8sTargets                     int
	MaxDeploymentTasks                        int
}

type StringSliceArgs []string

func (i *StringSliceArgs) String() string {
	return "A collection of strings passed as arguments"
}

func (i *StringSliceArgs) Set(value string) error {
	trimmed := strings.TrimSpace(value)

	if len(trimmed) == 0 {
		return nil
	}

	*i = append(*i, trimmed)
	return nil
}
