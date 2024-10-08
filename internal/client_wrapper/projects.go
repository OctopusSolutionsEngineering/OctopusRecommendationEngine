package client_wrapper

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/excluder"
	"github.com/samber/lo"
)

func GetProjects(limit int, client newclient.Client, spaceID string) ([]*projects.Project, error) {
	if limit == 0 {
		return projects.GetAll(client, spaceID)
	}

	result, err := projects.Get(client, spaceID, projects.ProjectsQuery{
		Take: limit,
	})

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func GetProjectByName(name string, client newclient.Client, spaceID string) ([]*projects.Project, error) {
	if name == "" {
		return []*projects.Project{}, nil
	}

	result, err := projects.Get(client, spaceID, projects.ProjectsQuery{
		PartialName: name,
	})

	if err != nil {
		return nil, err
	}

	for _, project := range result.Items {
		if project.Name == name {
			return []*projects.Project{project}, nil
		}
	}

	return []*projects.Project{}, nil
}

func GetProjectsWithFilter(client newclient.Client, spaceID string, excludeProjectsExcept config.StringSliceArgs, excludeProjects config.StringSliceArgs, maxItems int) ([]*projects.Project, error) {
	if len(excludeProjectsExcept) != 0 {
		return GetNamedProjects()
	}

	if allProjects, err := GetProjects(maxItems, client, spaceID); err != nil {
		return nil, err
	} else {
		defaultExcluder := excluder.DefaultExcluder{}
		return lo.Filter(allProjects, func(item *projects.Project, index int) bool {
			return !defaultExcluder.IsResourceExcluded(item.Name, false, excludeProjects, excludeProjectsExcept)
		}), nil
	}
}

func GetNamedProjects(client newclient.Client, spaceID string, excludeProjectsExcept config.StringSliceArgs) ([]*projects.Project, error) {
	projects := []*projects.Project{}

	for _, projectName := range excludeProjectsExcept {
		if project, err := GetProjectByName(projectName, client, spaceID); err != nil {
			return nil, err
		} else {
			projects = append(projects, project...)
		}
	}

	return projects, nil
}
