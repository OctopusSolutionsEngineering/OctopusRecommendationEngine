package entry

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks/factory"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/client_wrapper"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/executor"
	"github.com/briandowns/spinner"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var Version = "development"

func Entry(octolintConfig *config.OctolintConfig) ([]checks.OctopusCheckResult, error) {
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
		return nil, errors.New("You must specify the URL with the -url argument")
	}

	if _, err := url.ParseRequestURI(octolintConfig.Url); err != nil {
		return nil, errors.New("The URL \"" + octolintConfig.Url + "\" is not valid")
	}

	if octolintConfig.ApiKey == "" && octolintConfig.AccessToken == "" {
		return nil, errors.New("You must specify the API key with the -apiKey argument")
	}

	if octolintConfig.Space == "" {
		return nil, errors.New("You must specify the space key with the -space argument")
	}

	if !strings.HasPrefix(octolintConfig.Space, "Spaces-") {
		spaceId, err := lookupSpaceAsName(octolintConfig.Url, octolintConfig.Space, octolintConfig.ApiKey, octolintConfig.AccessToken)

		if err != nil {
			return nil, errors.New("Failed to create the Octopus client_wrapper. Check that the url, api key, and space are correct.\nThe error was: " + err.Error())
		}

		octolintConfig.Space = spaceId
	}

	httpClient, err := createHttpClient(octolintConfig)

	if err != nil {
		return nil, errors.New("Failed to create the HTTP client. Check that the url is correct.\nThe error was: " + err.Error())
	}

	client, err := createClient(httpClient, octolintConfig)

	if err != nil {
		return nil, errors.New("Failed to create the Octopus client_wrapper. Check that the url, api key, and space are correct.\nThe error was: " + err.Error())
	}

	checkFactory := factory.NewOctopusCheckFactory(client, octolintConfig.Url, octolintConfig.Space)
	checkCollection, err := checkFactory.BuildAllChecks(octolintConfig)

	if err != nil {
		ErrorExit("Failed to create the checks")
	}

	// Time the execution
	startTime := time.Now().UnixMilli()
	defer func() {
		endTime := time.Now().UnixMilli()
		fmt.Println("Report took " + fmt.Sprint((endTime-startTime)/1000) + " seconds")
	}()

	checkExecutor := executor.NewOctopusCheckExecutor()
	results, err := checkExecutor.ExecuteChecks(checkCollection, func(check checks.OctopusCheck, err error) error {
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
		return nil, errors.New("Failed to run the checks")
	}

	return results, nil
}

func createClient(httpClient *http.Client, octolintConfig *config.OctolintConfig) (*client.Client, error) {

	parsedUrl, err := getHost(octolintConfig)

	if err != nil {
		return nil, err
	}

	if octolintConfig.ApiKey != "" {
		return createClientApiKey(httpClient, parsedUrl, octolintConfig.Space, octolintConfig.ApiKey)
	}

	return createClientAccessToken(httpClient, parsedUrl, octolintConfig.Space, octolintConfig.AccessToken)
}

func getHost(octolintConfig *config.OctolintConfig) (*url.URL, error) {
	if octolintConfig.UseRedirector {
		return url.Parse("https://" + octolintConfig.RedirectorHost)
	}

	return url.Parse(octolintConfig.Url)
}

func createHttpClient(octolintConfig *config.OctolintConfig) (*http.Client, error) {
	if octolintConfig.UseRedirector {
		parsedUrl, err := url.Parse(octolintConfig.Url)

		if err != nil {
			return nil, err
		}

		headers := map[string]string{
			"X_REDIRECTION_UPSTREAM_HOST":   parsedUrl.Hostname(),
			"X_REDIRECTION_REDIRECTIONS":    octolintConfig.RedirectorRedirections,
			"X_REDIRECTION_API_KEY":         octolintConfig.RedirecrtorApiKey,
			"X_REDIRECTION_SERVICE_API_KEY": octolintConfig.RedirectorServiceApiKey,
		}

		return &http.Client{
			Transport: &client_wrapper.HeaderRoundTripper{
				Transport: http.DefaultTransport,
				Headers:   headers,
			},
		}, nil
	}

	return &http.Client{}, nil
}

func createClientApiKey(httpClient *http.Client, apiURL *url.URL, spaceId string, apiKey string) (*client.Client, error) {
	apiKeyCredential, err := client.NewApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return client.NewClientWithCredentials(httpClient, apiURL, apiKeyCredential, spaceId, "")
}

func createClientAccessToken(httpClient *http.Client, apiURL *url.URL, spaceId string, accessToken string) (*client.Client, error) {
	accessTokenCredential, err := client.NewAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	return client.NewClientWithCredentials(httpClient, apiURL, accessTokenCredential, spaceId, "")
}

func ErrorExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func lookupSpaceAsName(octopusUrl string, spaceName string, apiKey string, accessToken string) (string, error) {
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
	} else if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
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
