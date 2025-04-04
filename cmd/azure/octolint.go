package main

import (
	"encoding/json"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/args"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/entry"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/reporters"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
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
