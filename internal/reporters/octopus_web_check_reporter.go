package reporters

import (
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"strings"
)

// OctopusWebCheckReporter prints the lint reports in plain text to std out.
type OctopusWebCheckReporter struct {
	minSeverity int
}

func NewOctopusWebCheckReporter(minSeverity int) OctopusWebCheckReporter {
	return OctopusWebCheckReporter{minSeverity: minSeverity}
}

func (o OctopusWebCheckReporter) Generate(results []checks.OctopusCheckResult) (string, error) {
	if results == nil || len(results) == 0 {
		return "", nil
	}

	report := []string{}

	for _, r := range results {
		if r.Severity() >= o.minSeverity {
			report = append(report, r.Description())
		}
	}

	if len(report) == 0 {
		return "No issues detected", nil
	} else {
		report = append(report, "The checks are documented at https://github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/wiki")
	}

	return strings.Join(report[:], "\n\n"), nil
}
