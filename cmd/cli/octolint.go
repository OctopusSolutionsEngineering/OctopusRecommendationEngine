package main

import (
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/args"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/entry"
)

func main() {
	octolintConfig, err := args.ParseArgs()

	if err != nil {
		entry.ErrorExit(err.Error())
		return
	}

	entry.Entry(octolintConfig)
}
