package executor

import (
	"context"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/avast/retry-go/v4"
	"golang.org/x/sync/errgroup"
)

const ParallelTasks = 15

// OctopusCheckExecutor is responsible for running each lint check and returning the results. It deals with things
// like retries and error handling.
type OctopusCheckExecutor struct {
}

func NewOctopusCheckExecutor() OctopusCheckExecutor {
	return OctopusCheckExecutor{}
}

// ExecuteChecks executes each check and collects the results.
func (o OctopusCheckExecutor) ExecuteChecks(checkCollection []checks.OctopusCheck, handleError func(checks.OctopusCheck, error) error) ([]checks.OctopusCheckResult, error) {
	if checkCollection == nil || len(checkCollection) == 0 {
		return []checks.OctopusCheckResult{}, nil
	}

	checkResults := []checks.OctopusCheckResult{}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(ParallelTasks)

	for _, c := range checkCollection {
		c := c
		g.Go(func() error {
			err := retry.Do(
				func() error {
					result, err := c.Execute()

					if err != nil {
						checkResults = append(
							checkResults,
							checks.NewOctopusCheckResultImpl(
								"The check failed to run: "+err.Error(),
								c.Id(),
								"",
								checks.Error,
								checks.GeneralError))
					}

					if result != nil {
						checkResults = append(checkResults, result)
					}

					return nil
				}, retry.Attempts(3))

			if err != nil {
				err := handleError(c, err)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return checkResults, nil
}
