package getbuilds

import (
	"time"

	"github.com/concourse/atc/db"
	"github.com/concourse/atc/web/routes"
)

type PresentedBuild struct {
	ID           int
	JobName      string
	PipelineName string
	Status       string

	StartTime string
	EndTime   string

	CSSClass string
	Path     string
}

func formatDate(date time.Time) string {
	const layout = "2006-01-02 15:04:05 (MST)"
	return date.Format(layout)
}

func PresentBuilds(builds []db.Build) []PresentedBuild {
	presentedBuilds := []PresentedBuild{}

	for _, build := range builds {
		var cssClass string
		var jobName string
		var pipelineName string
		var startTime string

		if build.OneOff() {
			jobName = "[one off]"
			pipelineName = "[one off]"
			cssClass = "build-one-off"
		} else {
			jobName = build.JobName
			pipelineName = build.PipelineName
		}

		nilTime := time.Time{}
		if build.StartTime == nilTime {
			startTime = "failed to start"
		} else {
			startTime = formatDate(build.StartTime)
		}

		presentedBuilds = append(presentedBuilds, PresentedBuild{
			ID:           build.ID,
			JobName:      jobName,
			PipelineName: pipelineName,
			StartTime:    startTime,
			EndTime:      formatDate(build.EndTime),
			CSSClass:     cssClass,
			Status:       string(build.Status),
			Path:         routes.PathForBuild(build),
		})
	}

	return presentedBuilds
}
