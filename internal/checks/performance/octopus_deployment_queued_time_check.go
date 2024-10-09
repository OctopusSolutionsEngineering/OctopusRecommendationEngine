package performance

import (
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/events"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/checks"
	"github.com/OctopusSolutionsEngineering/OctopusRecommendationEngine/internal/config"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"math"
	"strings"
	"time"
)

const maxQueueTimeMinutes = 1
const maxQueuedTasks = 10
const OctoLintDeploymentQueuedTime = "OctoLintDeploymentQueuedTime"

type deploymentInfo struct {
	deploymentId string
	duration     float64
	queuedAt     time.Time
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
	config       *config.OctolintConfig
}

func NewOctopusDeploymentQueuedTimeCheck(client *client.Client, config *config.OctolintConfig, url string, space string, errorHandler checks.OctopusClientErrorHandler) OctopusDeploymentQueuedTimeCheck {
	return OctopusDeploymentQueuedTimeCheck{config: config, client: client, url: url, space: space, errorHandler: errorHandler}
}

func (o OctopusDeploymentQueuedTimeCheck) Id() string {
	return OctoLintDeploymentQueuedTime
}

func (o OctopusDeploymentQueuedTimeCheck) Execute(concurrency int) (checks.OctopusCheckResult, error) {
	if o.client == nil {
		return nil, errors.New("octoclient is nil")
	}

	zap.L().Debug("Starting check " + o.Id())

	defer func() {
		zap.L().Debug("Ended check " + o.Id())
	}()

	resource, err := o.client.Events.Get(events.EventsQuery{
		EventCategories: []string{"DeploymentQueued", "DeploymentStarted"},
		Skip:            0,
		Take:            o.config.MaxDeploymentTasks,
	})

	if err != nil {
		return o.errorHandler.HandleError(o.Id(), checks.Performance, err)
	}

	deployments := []deploymentInfo{}
	if resource != nil {
		for i, r := range resource.Items {
			zap.L().Debug(o.Id() + " " + fmt.Sprintf("%.2f", float32(i+1)/float32(len(resource.Items))*100) + "% complete")

			if r.Category == "DeploymentQueued" {
				queuedDeploymentId := o.getDeploymentFromRelatedDocs(r)
				for _, r2 := range resource.Items {
					if r2.Category == "DeploymentStarted" && queuedDeploymentId == o.getDeploymentFromRelatedDocs(r2) {
						queueTime := r2.Occurred.Sub(r.Occurred)
						if queueTime.Minutes() > maxQueueTimeMinutes {
							deployments = append(deployments, deploymentInfo{
								deploymentId: queuedDeploymentId,
								duration:     queueTime.Minutes(),
								queuedAt:     r.Occurred,
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
			return item.deploymentId + " (" + item.queuedAt.Format(time.RFC822) + " " + fmt.Sprint(item.toFixed(1)) + "m)"
		}

		return o.url + "/app#/" + o.space + "/projects/" + deployment.ProjectID + "/deployments/releases/" + deployment.ReleaseID +
			"/deployments/" + item.deploymentId + " (" + item.queuedAt.Format(time.RFC822) + " " + fmt.Sprint(item.toFixed(1)) + "m)"
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
