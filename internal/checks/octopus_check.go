package checks

// OctopusCheck defines the contract for each lint check
type OctopusCheck interface {
	// Execute runs the check
	Execute(concurrency int) (OctopusCheckResult, error)
	// Id returns the unique ID of the check, used to cross-reference with documentation
	Id() string
}
