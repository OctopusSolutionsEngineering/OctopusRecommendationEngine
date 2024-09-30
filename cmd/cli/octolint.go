package main

import (
	"fmt"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/args"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/entry"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/reporters"
	"os"
)

func main() {
	octolintConfig, err := args.ParseArgs(os.Args[1:])

	if err != nil {
		entry.ErrorExit(err.Error())
		return
	}

	results := entry.Entry(octolintConfig)

	reporter := reporters.NewOctopusPlainCheckReporter(checks.Warning)
	report, err := reporter.Generate(results)

	if err != nil {
		entry.ErrorExit("Failed to generate the report")
	}

	fmt.Println(report)
}
