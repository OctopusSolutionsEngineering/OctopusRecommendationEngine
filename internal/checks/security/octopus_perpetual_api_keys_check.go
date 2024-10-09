package security

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"time"
)

type APIKeyKey struct {
	Hint *string
}

// APIKey is used because the go client_wrapper has an invalid APIKey value that prevents the usual functions for querying users keys
type APIKey struct {
	APIKey  APIKeyKey  `json:"ApiKey,omitempty"`
	Expires *time.Time `json:"Expires,omitempty"`
}

// OctopusPerpetualApiKeysCheck reports on any perpetual api keys
type OctopusPerpetualApiKeysCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusPerpetualApiKeysCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusPerpetualApiKeysCheck {
	return OctopusPerpetualApiKeysCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusPerpetualApiKeysCheck) Id() string {
	return "OctoLintPerpetualApiKeys"
}

func (o OctopusPerpetualApiKeysCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	users, err := o.client.Users.GetAll()

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	linksTemplate := regexp.MustCompile(`\{.+\}`)
	perpetualApiKeys := []string{}
	for i, u := range users {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(users))*100) + "% complete")

		apiKeysLink := linksTemplate.ReplaceAllString(u.Links["ApiKeys"], "")
		keys, err := newclient.Get[resources.Resources[APIKey]](o.client.HttpSession(), apiKeysLink)

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		for _, k := range keys.Items {
			if k.Expires == nil && k.APIKey.Hint != nil && u.Username != "guest" {
				perpetualApiKeys = append(perpetualApiKeys, *k.APIKey.Hint+"... ("+u.Username+")")
			}
		}
	}

	if len(perpetualApiKeys) != 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following API keys do not expire:\n"+strings.Join(perpetualApiKeys, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"No perpetual API keys found",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}
