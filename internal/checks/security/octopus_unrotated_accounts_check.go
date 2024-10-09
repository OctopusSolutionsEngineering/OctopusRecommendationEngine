package security

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"net/url"
	"strings"
	"time"
)

const maxTimeSinceAccountEdit = time.Hour * 24 * 90

// OctopusUnrotatedAccountsCheck checks to see if any targets have not been used in a month
type OctopusUnrotatedAccountsCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

func NewOctopusUnrotatedAccountsCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctopusUnrotatedAccountsCheck {
	return OctopusUnrotatedAccountsCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctopusUnrotatedAccountsCheck) Id() string {
	return "OctoLintUnrotatedAccounts"
}

func (o OctopusUnrotatedAccountsCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	now := time.Now()
	start := now.Add(maxTimeSinceAccountEdit * -1)
	end := now

	allAccounts, err := o.GetAll(o.client, o.client.GetSpaceID())

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}

	uneditedAccounts := []string{}
	for i, m := range allAccounts {

		zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(allAccounts))*100) + "% complete")

		// Skip OIDC accounts
		if m.GetAccountType() == "AmazonWebServicesOidcAccount" {
			continue
		}

		if m.GetAccountType() == "AzureOidc" {
			continue
		}

		audits, err := newclient.Get[resources.Resources[OctopusAudit]](o.client.HttpSession(), "/api/events?regardingAny="+m.GetID()+"&from="+url.QueryEscape(start.Format("2006-01-02T15:04:05-0700"))+"&to="+url.QueryEscape(end.Format("2006-01-02T15:04:05-0700")))

		if err != nil {
			if !o.errorHandler.ShouldContinue(err) {
				return nil, err
			}
			continue
		}

		recentEdit := false
		for _, t := range audits.Items {
			if t.Category == "Modified" && slices.Index(t.RelatedDocumentIds, m.GetID()) != -1 {
				recentEdit = true
			}
		}

		if !recentEdit {
			uneditedAccounts = append(uneditedAccounts, m.GetName())
		}

	}

	if len(uneditedAccounts) > 0 {
		return checks.NewOctopusCheckResultImpl(
			"The following accounts have not been updated in 90 days:\n"+strings.Join(uneditedAccounts, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no unedited accounts",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}

type OctopusAudit struct {
	Category           string
	RelatedDocumentIds []string
}

func (o OctopusUnrotatedAccountsCheck) GetAll(client newclient.Client, spaceID string) ([]*accounts.AccountResource, error) {
	items, err := newclient.GetAll[accounts.AccountResource](client, "/api/{spaceId}/accounts", spaceID)
	return items, err
}
