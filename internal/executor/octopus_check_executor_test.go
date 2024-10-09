package executor

import (
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"testing"
)

type alwaysFailCheck struct {
}

func (o alwaysFailCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	return checks.NewOctopusCheckResultImpl("This check always fails", o.Id(), "", checks.Error, ""), nil
}

func (o alwaysFailCheck) Id() string {
	return "OctoRecAlwaysFail"
}

type alwaysPassCheck struct {
}

func (o alwaysPassCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	return checks.NewOctopusCheckResultImpl("This check passed ok", o.Id(), "", checks.Ok, ""), nil
}

func (o alwaysPassCheck) Id() string {
	return "OctoRecAlwaysPass"
}

func TestNoChecks(t *testing.T) {
	results, err := OctopusCheckExecutor{}.ExecuteChecks(nil, func(check checks.OctopusCheck, err error) error {
		return nil
	})

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if len(results) != 0 {
		t.Fatal("Should not have returned any results")
	}
}

func TestFailChecks(t *testing.T) {
	results, err := OctopusCheckExecutor{}.ExecuteChecks([]checks.OctopusCheck{alwaysFailCheck{}}, func(check checks.OctopusCheck, err error) error {
		return nil
	})

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if len(results) != 1 {
		t.Fatal("Should have returned 1 results")
	}
}

func TestFailAndPassChecks(t *testing.T) {
	results, err := OctopusCheckExecutor{}.ExecuteChecks([]checks.OctopusCheck{alwaysFailCheck{}, alwaysPassCheck{}}, func(check checks.OctopusCheck, err error) error {
		return nil
	})

	if err != nil {
		t.Fatal("Should not have returned an error")
	}

	if len(results) != 2 {
		t.Fatal("Should have returned 2 results")
	}
}
