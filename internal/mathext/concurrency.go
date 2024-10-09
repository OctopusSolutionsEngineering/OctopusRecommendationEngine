package mathext

func TopLevelConcurrency(maxConcurrency, numberOfChecks int) int {
	return MinInt(maxConcurrency, numberOfChecks)
}

func InternalLevelConcurrency(maxExternalConcurrency, minInternalConcurrency, numberOfChecks int) int {
	topLevelConcurrency := TopLevelConcurrency(maxExternalConcurrency, numberOfChecks)

	if topLevelConcurrency == 0 {
		return maxExternalConcurrency * minInternalConcurrency
	}

	return MaxInt(minInternalConcurrency, maxExternalConcurrency*minInternalConcurrency/topLevelConcurrency)
}
