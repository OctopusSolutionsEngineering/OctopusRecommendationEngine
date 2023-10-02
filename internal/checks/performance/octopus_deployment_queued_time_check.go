package performance

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/events"
	"github.com/mcasperson/OctopusRecommendationEngine/internal/checks"
	"github.com/samber/lo"
	"math"
	"strings"
)

const maxQueueTimeMinutes = 1
const maxQueuedTasks = 10

type deploymentInfo struct {
	deploymentId string
	duration     float64
}

func (d deploymentInfo) round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func (d deploymentInfo) toFixed(precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(d.round(d.duration*output)) / output
}

// OctopusDeploymentQueuedTimeCheck checks to see if any deployments were queued for a long period of time
type OctopusDeploymentQueuedTimeCheck struct {
	client       *client.Client
	errorHandler checks.OctopusClientErrorHandler
	url          string
	space        string
}

func NewOctopusDeploymentQueuedTimeCheck(client *client.Client, url string, space string, errorHandler checks.OctopusClientErrorHandler) OctopusDeploymentQueuedTimeCheck {
	return OctopusDeploymentQueuedTimeCheck{client: client, url: url, space: space, errorHandler: errorHandler}
}

func (o OctopusDeploymentQueuedTimeCheck) Id() string {
	return "OctoLintDeploymentQueuedTime"
}

func (o OctopusDeploymentQueuedTimeCheck) Execute() (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	resource, err := o.client.Events.Get(events.EventsQuery{
		EventCategories: []string{"DeploymentQueued", "DeploymentStarted"},
		Skip:            0,
		Take:            1000,
	})

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Performance, err)
	}

	deployments := []deploymentInfo{}
	if resource != nil {
		for _, r := range resource.Items {
			if r.Category == "DeploymentQueued" {
				queuedDeploymentId := o.getDeploymentFromRelatedDocs(r)
				for _, r2 := range resource.Items {
					if r2.Category == "DeploymentStarted" && queuedDeploymentId == o.getDeploymentFromRelatedDocs(r2) {
						queueTime := r2.Occurred.Sub(r.Occurred)
						if queueTime.Minutes() > maxQueueTimeMinutes {
							deployments = append(deployments, deploymentInfo{
								deploymentId: queuedDeploymentId,
								duration:     queueTime.Minutes(),
							})
						}
					}
				}
			}
		}
	}

	deploymentLinks := lo.Map(deployments, func(item deploymentInfo, index int) string {
		deployment, err := o.client.Deployments.GetByID(item.deploymentId)

		if err != nil {
			return item.deploymentId + " (" + fmt.Sprint(item.toFixed(1)) + "m)"
		}

		return o.url + "/app#/" + o.space + "/projects/" + deployment.ProjectID + "/deployments/releases/" + deployment.ReleaseID +
			"/deployments/" + item.deploymentId + " (" + fmt.Sprint(item.toFixed(1)) + "m)"
	})

	if len(deployments) >= maxQueuedTasks {
		return checks.NewOctopusCheckResultImpl(
			fmt.Sprint("Found "+fmt.Sprint(len(deployments)))+" deployments that were queued for longer than "+fmt.Sprint(maxQueueTimeMinutes)+" minutes. Consider increasing the task cap or adding a HA node to reduce task queue times:\n"+
				strings.Join(deploymentLinks, "\n"),
			o.Id(),
			"",
			checks.Warning,
			checks.Performance), nil
	}

	return checks.NewOctopusCheckResultImpl(
		"Found "+fmt.Sprint(len(deployments))+" deployment tasks that were queued for longer than "+fmt.Sprint(maxQueueTimeMinutes)+" minutes:\n"+
			strings.Join(deploymentLinks, ", "),
		o.Id(),
		"",
		checks.Ok,
		checks.Performance), nil
}

func (o OctopusDeploymentQueuedTimeCheck) getDeploymentFromRelatedDocs(event *events.Event) string {
	for _, d := range event.RelatedDocumentIds {
		if strings.HasPrefix(d, "Deployments-") {
			return d
		}
	}
	return ""
}
