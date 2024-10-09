package security

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"strings"
)

// OctopusInsecureFeedsCheck checks to see if any targets have not been used in a month
type OctopusInsecureFeedsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusInsecureFeedsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusInsecureFeedsCheck {
	return OctopusInsecureFeedsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusInsecureFeedsCheck) Id() string {
	return "OctoLintInsecureFeedsTargets"
}

func (o OctopusInsecureFeedsCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	targets, err := o.client.Feeds.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	insecureFeeds := []string{}
	for i, m := range targets {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(targets))*100) + "% complete")

		if m.GetFeedType() == "ArtifactoryGeneric" {
			typedFeed := m.(*feeds.ArtifactoryGenericFeed)
			if strings.HasPrefix(typedFeed.FeedURI, "http://") {
				insecureFeeds = append(insecureFeeds, m.GetName())
			}
		}

		if m.GetFeedType() == "NuGet" {
			typedFeed := m.(*feeds.NuGetFeed)
			if strings.HasPrefix(typedFeed.FeedURI, "http://") {
				insecureFeeds = append(insecureFeeds, m.GetName())
			}
		}

		if m.GetFeedType() == "Maven" {
			typedFeed := m.(*feeds.MavenFeed)
			if strings.HasPrefix(typedFeed.FeedURI, "http://") {
				insecureFeeds = append(insecureFeeds, m.GetName())
			}
		}

		if m.GetFeedType() == "Helm" {
			typedFeed := m.(*feeds.HelmFeed)
			if strings.HasPrefix(typedFeed.FeedURI, "http://") {
				insecureFeeds = append(insecureFeeds, m.GetName())
			}
		}

		if m.GetFeedType() == "GitHub" {
			typedFeed := m.(*feeds.GitHubRepositoryFeed)
			if strings.HasPrefix(typedFeed.FeedURI, "http://") {
				insecureFeeds = append(insecureFeeds, m.GetName())
			}
		}

		if m.GetFeedType() == "Docker" {
			typedFeed := m.(*feeds.DockerContainerRegistry)
			if strings.HasPrefix(typedFeed.FeedURI, "http://") {
				insecureFeeds = append(insecureFeeds, m.GetName())
			}
		}

	}

	if len(insecureFeeds) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following feeds use an insecure HTTP endpoint:\n"+strings.Join(insecureFeeds, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no insecure feeds",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}
