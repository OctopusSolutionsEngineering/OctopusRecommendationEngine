package checks

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"net/http"
	"strings"
)

// OctopusClientPermissiveErrorHandler treats almost every 400 HTTP error as a permissions error and returns
// a result at the Permission level.
type OctopusClientPermissiveErrorHandler struct {
}

func (o OctopusClientPermissiveErrorHandler) HandleError(id string, group string, err error) (OctopusCheckResult, error) {
	if o.ShouldContinue(err) {
		return NewOctopusCheckResultImpl(
			"You do not have permission to run the check: "+err.Error(),
			id,
			"",
			Permission,
			group), nil
	}
	return nil, err
}

// ShouldContinue is used to determine if an error was a permissions error. Things like 404s are also treated
// as permission errors (we saw this a lot trying to get deployment processes). Interestingly we also saw a lot of
// StatusCode's set to 0, so this function also reads the error to work out what is going on.
func (o OctopusClientPermissiveErrorHandler) ShouldContinue(err error) bool {
	apiError, ok := err.(*core.APIError)
	if ok {
		return apiError.StatusCode == http.StatusUnauthorized ||
			apiError.StatusCode == http.StatusForbidden ||
			apiError.StatusCode == http.StatusNotFound ||
			strings.Index(strings.ToLower(apiError.ErrorMessage), "you do not have permission") != -1 ||
			strings.Index(strings.ToLower(apiError.ErrorMessage), "invalid username or password") != -1 ||
			strings.Index(strings.ToLower(apiError.ErrorMessage), "support for password authentication was removed") != -1
	}
	return true
}
