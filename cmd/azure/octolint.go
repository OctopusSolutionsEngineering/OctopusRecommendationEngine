package main

import (
	"encoding/json"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/args"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/entry"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/environment"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/reporters"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type AzureFunctionRequestDataReq struct {
	Body string `json:"Body"`
}

type AzureFunctionRequestData struct {
	Req AzureFunctionRequestDataReq `json:"req"`
}

type AzureFunctionRequest struct {
	Data AzureFunctionRequestData `json:"Data"`
}

func octoterraHandler(w http.ResponseWriter, r *http.Request) {
	// Allow the more sensitive values to be passed as headers
	apiKey := r.Header.Get("X-Octopus-ApiKey")
	accessToken := r.Header.Get("X-Octopus-AccessToken")
	url := r.Header.Get("X-Octopus-Url")
	redirectorRedirections := r.Header.Get("X_REDIRECTION_REDIRECTIONS")
	redirectorApiKey := r.Header.Get("X_REDIRECTION_API_KEY")
	redirectorServiceApiKey, _ := os.LookupEnv("REDIRECTION_SERVICE_API_KEY")
	redirectorHost, _ := os.LookupEnv("REDIRECTION_HOST")
	disableRedirector, _ := os.LookupEnv("DISABLE_REDIRECTION")

	enableRedirector, err := useRedirector(url, disableRedirector, redirectorServiceApiKey, redirectorHost, redirectorRedirections, redirectorApiKey)

	if err != nil {
		handleError(err, w)
		return
	}

	respBytes, err := io.ReadAll(r.Body)

	if err != nil {
		handleError(err, w)
		return
	}

	if len(respBytes) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("Request body is empty"))
		return
	}

	file, err := os.CreateTemp("", "*.json")

	if err != nil {
		handleError(err, w)
		return
	}

	configJson, err := sanitizeConfig(respBytes)

	if err != nil {
		handleError(err, w)
		return
	}

	err = os.WriteFile(file.Name(), configJson, 0644)

	if err != nil {
		handleError(err, w)
		return
	}

	// Clean up the file when we are done
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}(file.Name())

	filename := filepath.Base(file.Name())
	extension := filepath.Ext(filename)
	filenameWithoutExtension := filename[0 : len(filename)-len(extension)]

	commandLineArgs := []string{"-spinner=False", "-configFile", filenameWithoutExtension, "-configPath", filepath.Dir(file.Name())}

	if apiKey != "" {
		commandLineArgs = append(commandLineArgs, "-apiKey", apiKey)
	} else if accessToken != "" {
		commandLineArgs = append(commandLineArgs, "-accessToken", accessToken)
	}

	if url != "" {
		commandLineArgs = append(commandLineArgs, "-url", url)
	}

	if enableRedirector {
		commandLineArgs = append(commandLineArgs, "-useRedirector")
		commandLineArgs = append(commandLineArgs, "-redirectorHost", redirectorHost)
		commandLineArgs = append(commandLineArgs, "-redirectorServiceApiKey", redirectorServiceApiKey)
		commandLineArgs = append(commandLineArgs, "-redirecrtorApiKey", redirectorApiKey)
		commandLineArgs = append(commandLineArgs, "-redirectorRedirections", redirectorRedirections)
	}

	webArgs, err := args.ParseArgs(commandLineArgs)

	if err != nil {
		handleError(err, w)
		return
	}

	results, err := entry.Entry(webArgs)

	if err != nil {
		handleError(err, w)
		return
	}

	reporter := reporters.NewOctopusWebCheckReporter(checks.Warning)
	report, err := reporter.Generate(results)

	if err != nil {
		handleError(err, w)
		return
	}

	w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(200)
	if _, err := w.Write([]byte(report)); err != nil {
		zap.L().Error(err.Error())
	}
}

func useRedirector(octopusUrl string, disableRedirector string, redirectorServiceApiKey string, redirectorHost string, redirections string, redirectorApiKey string) (bool, error) {
	parsedUrl, err := url.Parse(octopusUrl)

	if err != nil {
		return false, err
	}

	disableRedirectorParsed, err := strconv.ParseBool(disableRedirector)

	if err != nil {
		disableRedirectorParsed = false
	}

	return !disableRedirectorParsed && redirectorServiceApiKey != "" && redirectorHost != "" &&
		(!hostIsCloudOrLocal(parsedUrl.Hostname()) ||
			(redirections != "" && redirectorApiKey != "")), nil
}

func hostIsCloudOrLocal(host string) bool {
	return strings.HasSuffix(host, ".octopus.app") ||
		strings.HasSuffix(host, ".testoctopus.com") ||
		host == "localhost" ||
		host == "127.0.0.1"
}

// sanitizeConfig removes sensitive information from the config so it is not
// persisted to the disk.
func sanitizeConfig(rawConfig []byte) ([]byte, error) {
	if len(rawConfig) == 0 {
		return rawConfig, nil
	}

	config := map[string]any{}
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		return nil, err
	}
	delete(config, "apiKey")
	delete(config, "url")
	delete(config, "redirectorServiceApiKey")
	delete(config, "redirecrtorApiKey")
	delete(config, "redirectorHost")
	delete(config, "useRedirector")
	delete(config, "redirectorRedirections")
	return json.Marshal(config)
}

func handleError(err error, w http.ResponseWriter) {
	zap.L().Error(err.Error())
	w.WriteHeader(500)
	if _, err := w.Write([]byte(err.Error())); err != nil {
		zap.L().Error(err.Error())
	}
}

func main() {
	listenAddr := ":" + environment.GetPort()
	http.HandleFunc("/api/octolint", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			octoterraHandler(writer, request)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
			w.WriteHeader(200)
			w.Write([]byte("Healthy"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
