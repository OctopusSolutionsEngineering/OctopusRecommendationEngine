package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"go.uber.org/zap"
)

const (
	OctoLintSha1Certificates = "OctoLintSha1Certificates"
	sha1Alg                  = "sha1RSA"
)

// OctoLintSha1CertificatesCheck checks to see if any targets, workers or the server itself is using a sha1 certificate
type OctoLintSha1CertificatesCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	config       *config.OctolintConfig
}

type Sha1CertificateResult struct {
	Name string
	Type string // "Target", "Worker", or "Global"
}

type ServerCertificate struct {
	ID                 string            `json:"Id"`
	Name               string            `json:"Name"`
	Thumbprint         string            `json:"Thumbprint"`
	SignatureAlgorithm string            `json:"SignatureAlgorithm"`
	Links              map[string]string `json:"Links"`
}

func NewOctoLintSha1CertificatesCheck(client *client.Client, config *config.OctolintConfig, errorHandler checks.OctopusClientErrorHandler) OctoLintSha1CertificatesCheck {
	return OctoLintSha1CertificatesCheck{config: config, client: client, errorHandler: errorHandler}
}

func (o OctoLintSha1CertificatesCheck) Id() string {
	return OctoLintSha1Certificates
}

// fetchServerCertificate gets the server certificate object and returns it.
func fetchServerCertificate(url, apiKey, accessToken string) (*ServerCertificate, error) {
	requestURL := fmt.Sprintf("%s/api/configuration/certificates/certificate-global", url)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	if apiKey != "" {
		req.Header.Set("X-Octopus-ApiKey", apiKey)
	} else if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var cert ServerCertificate
	if err := json.NewDecoder(res.Body).Decode(&cert); err != nil {
		return nil, err
	}

	return &cert, nil
}

// hasSha1Certificate checks if an endpoint has CertificateSignatureAlgorithm == sha1Alg
func hasSha1Certificate(ep machines.IEndpoint) bool {
	if ep == nil {
		return false
	}

	switch e := ep.(type) {
	case *machines.ListeningTentacleEndpoint:
		return e.CertificateSignatureAlgorithm == sha1Alg
	case *machines.PollingTentacleEndpoint:
		return e.CertificateSignatureAlgorithm == sha1Alg
	default:
		return false
	}
}

// addSha1FromMachines is a top-level generic helper (function literals can't have type params).
func addSha1FromMachines[T any](
	results *[]Sha1CertificateResult,
	items []T,
	getName func(T) string,
	getEndpoint func(T) machines.IEndpoint,
	typ string,
) {
	for _, item := range items {
		if hasSha1Certificate(getEndpoint(item)) {
			*results = append(*results, Sha1CertificateResult{Name: getName(item), Type: typ})
		}
	}
}

func (o OctoLintSha1CertificatesCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())
	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	var results []Sha1CertificateResult

	// Check server certificate
	cert, err := fetchServerCertificate(o.config.Url, o.config.ApiKey, o.config.AccessToken)
	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}
	if cert != nil && cert.SignatureAlgorithm == sha1Alg {
		results = append(results, Sha1CertificateResult{Name: cert.Name, Type: "Global"})
	}

	// Check deployment targets
	targets, err := client_wrapper.GetMachines(o.config.MaxSha1CertificatesMachines, o.client, o.client.GetSpaceID())
	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}
	addSha1FromMachines(&results, targets,
		func(m *machines.DeploymentTarget) string { return m.Name },
		func(m *machines.DeploymentTarget) machines.IEndpoint { return m.Endpoint },
		"Target",
	)

	// Check workers
	workers, err := client_wrapper.GetWorkers(o.config.MaxSha1CertificatesMachines, o.client, o.client.GetSpaceID())
	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Security, err)
	}
	addSha1FromMachines(&results, workers,
		func(w *machines.Worker) string { return w.Name },
		func(w *machines.Worker) machines.IEndpoint { return w.Endpoint },
		"Worker",
	)

	// Provide results
	if len(results) > 0 {
		// Sort by Type then Name for stable output
		sort.Slice(results, func(i, j int) bool {
			if results[i].Type == results[j].Type {
				return results[i].Name < results[j].Name
			}
			return results[i].Type < results[j].Type
		})

		lines := make([]string, len(results))
		for i, m := range results {
			lines[i] = fmt.Sprintf("%s: %s", m.Type, m.Name)
		}

		return checks.NewOctopusCheckResultImpl(
			"The following resources use a SHA1 certificate:\n"+strings.Join(lines, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Security), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"There are no uses of SHA1 certificates in targets, workers or the main Server Certificate",
		o.Id(),
		"",
		checks.Ok,
		checks.Security), nil
}
