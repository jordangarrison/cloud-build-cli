package cloudbuild

import (
	"context"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

type CloudBuildClient struct {
	Service   *cloudbuild.Service
	ProjectID string
}

func NewCloudBuildClient(ctx context.Context, projectID string) (*CloudBuildClient, error) {
	service, err := cloudbuild.NewService(ctx)
	if err != nil {
		return nil, err
	}
	client := &CloudBuildClient{
		Service:   service,
		ProjectID: projectID,
	}
	return client, nil
}

type cloudbuildResult struct {
	Build   *cloudbuild.Build
	Trigger *cloudbuild.BuildTrigger
}

func (client *CloudBuildClient) GetCurrentBuilds() ([]cloudbuildResult, error) {
	builds := cloudbuild.NewProjectsBuildsService(client.Service)
	triggers := cloudbuild.NewProjectsTriggersService(client.Service)
	call, err := builds.List(client.ProjectID).Do()
	var cloudBuildResults []cloudbuildResult
	for _, build := range call.Builds {
		trigger, err := triggers.Get(client.ProjectID, build.BuildTriggerId).Do()
		if err != nil {
			return cloudBuildResults, err
		}
		result := *&cloudbuildResult{
			Build:   build,
			Trigger: trigger,
		}
		cloudBuildResults = append(cloudBuildResults, result)
	}
	if err != nil {
		return cloudBuildResults, err
	}
	return cloudBuildResults, err
}
