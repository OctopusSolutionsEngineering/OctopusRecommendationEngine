package organization

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	projects2 "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

const OctoLintDuplicatedVariables = "OctoLintDuplicatedVariables"

type projectVar struct {
	project1  *projects2.Project
	variable1 *variables.Variable
	project2  *projects2.Project
	variable2 *variables.Variable
}

// OctopusDuplicatedVariablesCheck checks for variables with the same value across projects. This may be an indication
// that library variable sets should be used to capture shared values.
type OctopusDuplicatedVariablesCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusDuplicatedVariablesCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusDuplicatedVariablesCheck {
	return OctopusDuplicatedVariablesCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusDuplicatedVariablesCheck) Id() string {
	return OctoLintDuplicatedVariables
}

func (o OctopusDuplicatedVariablesCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	projects, err := client_wrapper.GetProjectsWithFilter(
		o.client,
		o.client.GetSpaceID(),
		o.config.ExcludeProjectsExcept,
		o.config.ExcludeProjects,
		o.config.MaxDuplicateVariableProjects)

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Organization, err)
	}

	projectVars := map[*projects2.Project]variables.VariableSet{}
	for i, p := range projects {
		if o.config.MaxDuplicateVariableProjects != 0 && i >= o.config.MaxDuplicateVariableProjects {
			break
		}

		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(projects))*100) + "% complete")

		variableSet, err := o.client.Variables.GetAll(p.ID)

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		projectVars[p] = variableSet
	}

	duplicateVars := []projectVar{}

OuterLoop:
	for index1 := 0; index1 < len(projects); index1++ {
		project1 := projects[index1]
		for _, variable1 := range projectVars[project1].Variables {
			if o.config.MaxDuplicateVariables != 0 && len(duplicateVars) >= o.config.MaxDuplicateVariables {
				break OuterLoop
			}

			if o.shouldIgnoreVariable(variable1) {
				continue
			}

			for index2 := index1 + 1; index2 < len(projects); index2++ {
				project2 := projects[index2]
				for _, variable2 := range projectVars[project2].Variables {
					if variable1.Value == variable2.Value {
						duplicateVars = append(duplicateVars, projectVar{
							project1:  project1,
							variable1: variable1,
							project2:  project2,
							variable2: variable2,
						})
					}
				}
			}
		}
	}

	if len(duplicateVars) > 0 {
		messages := []string{}
		for _, variable := range duplicateVars {
			messages = append(messages, variable.project1.Name+"/"+variable.variable1.Name+" == "+variable.project2.Name+"/"+variable.variable2.Name)
		}

		return checks.NewOctopusCheckResultImpl(
			"The following variables are duplicated between projects. Consider moving these into library variable sets:\n"+strings.Join(messages, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Organization), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no duplicated variables",
		o.Id(),
		"",
		checks.Ok,
		checks.Organization), nil
}

func (o OctopusDuplicatedVariablesCheck) shouldIgnoreVariable(variable *variables.Variable) bool {
	_, err := strconv.Atoi(variable.Value)
	return variable.Value == "" ||
		variable.Type != "String" ||
		slices.Index(checks.SpecialVars, variable.Name) != -1 ||
		strings.ToLower(variable.Value) == "true" ||
		strings.ToLower(variable.Value) == "false" ||
		strings.ToLower(variable.Value) == "yes" ||
		strings.ToLower(variable.Value) == "no" ||
		err == nil
}
