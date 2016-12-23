package main

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/ovh/cds/sdk"
	"github.com/spf13/viper"
)

const (
	inProgress = "INPROGRESS"
	successful = "SUCCESSFUL"
	failed     = "FAILED"
)

// BitbucketRequest ...
type BitbucketRequest struct {
	Description string `json:"description"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	State       string `json:"state"`
	URL         string `json:"url"`
}

// Process send message to all notifications backend
func process(event *sdk.EventPipelineBuild) {

	log.Debugf("Process event:%+v", event)

	cdsProject := event.ProjectKey
	cdsApplication := event.ApplicationName
	cdsPipelineName := event.PipelineName
	cdsBuildNumber := event.BuildNumber
	cdsEnvironmentName := event.EnvironmentName

	key := fmt.Sprintf("%s-%s-%s",
		cdsProject,
		cdsApplication,
		cdsPipelineName,
	)

	// project/CDS/application/cds2tat/pipeline/monPipeline/build/855?env=monEnvi
	url := fmt.Sprintf("%s/project/%s/application/%s/pipeline/%s/build/%d?env=%s",
		viper.GetString("url_cds_ui"),
		cdsProject,
		cdsApplication,
		cdsPipelineName,
		cdsBuildNumber,
		cdsEnvironmentName,
	)

	r := &BitbucketRequest{
		Key:   key,
		Name:  fmt.Sprintf("%s%d", key, cdsBuildNumber),
		State: getBitbucketStateFromStatus(event.Status),
		URL:   url,
	}

	jsonStr, err := json.Marshal(r)
	if err != nil {
		log.Errorf("Error while marshalling bitbucketRequest: %s", err.Error())
		return
	}

	// http://localhost:7990/bitbucket
	// /rest/build-status/1.0/commits/9e72f04322c4a1f240e0b3158c67c3c19cdd16e7
	pathBitbucket := fmt.Sprintf("%s/rest/build-status/1.0/commits/%s", viper.GetString("url_bitbucket"), event.Hash)

	log.Debugf("bitbucket url %+v with json:%s", pathBitbucket, jsonStr)

	if _, err := reqWant(pathBitbucket, "POST", jsonStr); err != nil {
		log.Warnf("Error on bitbucket: %ss", err)
	}

	return
}

func getBitbucketStateFromStatus(status sdk.Status) string {
	switch status {
	case sdk.StatusSuccess:
		return successful
	case sdk.StatusWaiting:
		return inProgress
	case sdk.StatusBuilding:
		return inProgress
	case sdk.StatusFail:
		return failed
	default:
		return failed
	}
}
