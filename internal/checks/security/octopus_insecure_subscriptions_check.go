package security

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/services/api"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"strings"
)

// OctopusInsecureSubscriptionsCheck checks to see if any targets have not been used in a month
type OctopusInsecureSubscriptionsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusInsecureSubscriptionsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusInsecureSubscriptionsCheck {
	return OctopusInsecureSubscriptionsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusInsecureSubscriptionsCheck) Id() string {
	return "OctoLintInsecureWebhookUrls"
}

func (o OctopusInsecureSubscriptionsCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	collection := resources.Resources[*OctopusSubscription]{}
	_, err := api.ApiGet(o.client.Subscriptions.GetClient(), &collection, o.client.Subscriptions.BasePath+"?skip=0&take=2147483647")

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	insecureItems := []string{}
	for i, m := range collection.Items {
		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(collection.Items))*100) + "% complete")

		if m.EventNotificationSubscription != nil && strings.HasPrefix(m.EventNotificationSubscription.WebhookURI, "http://") {
			insecureItems = append(insecureItems, m.Name)
		}

	}

	if len(insecureItems) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following subscriptions use an insecure HTTP webhook URL:\n"+strings.Join(insecureItems, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no insecure subscriptions",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}

type OctopusSubscription struct {
	Name                          string
	EventNotificationSubscription *OctopusEventNotificationSubscription
}

type OctopusEventNotificationSubscription struct {
	WebhookURI string
}
