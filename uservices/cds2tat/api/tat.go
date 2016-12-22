package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/ovh/cds/sdk"
	"github.com/ovh/tat"
)

const (
	failed   = "#a94442"
	success  = "#3c763d"
	waiting  = "#8a6d3b"
	building = "#31708f"
	unknown  = "#333"
)

func processEventPipelineBuild(e *sdk.EventPipelineBuild) {
	eventType := "pipelineBuild"
	cdsProject := e.ProjectKey
	cdsApp := e.ApplicationName
	cdsPipeline := e.PipelineName
	cdsEnvironment := e.EnvironmentName
	version := e.Version
	branch := e.BranchName

	processMsg(eventType, cdsProject, cdsApp, cdsPipeline, cdsEnvironment, version, branch, e.Status)
}

func processEventJob(e *sdk.EventJob) {
	eventType := "job"
	cdsProject := e.ProjectKey
	cdsApp := e.ApplicationName
	cdsPipeline := e.PipelineName
	cdsEnvironment := e.EnvironmentName
	version := e.Version
	branch := e.BranchName

	processMsg(eventType, cdsProject, cdsApp, cdsPipeline, cdsEnvironment, version, branch, e.Status)
}

func processMsg(eventType, cdsProject, cdsApp, cdsPipeline, cdsEnvironment string, version int64, branch string, cdsStatus sdk.Status) {

	text := fmt.Sprintf("#cds #type:%s #project:%s #app:%s #pipeline:%s #environment:%s #version:%d #branch:%s",
		eventType, cdsProject, cdsApp, cdsPipeline, cdsEnvironment, version, branch)

	tagsReference := fmt.Sprintf("#cds,#type:%s,#project:%s,#app:%s,#pipeline:%s,#environment:%s,#version:%d,#branch:%s",
		eventType, cdsProject, cdsApp, cdsPipeline, cdsEnvironment, version, branch)

	msg := tat.MessageJSON{
		Text:         text,
		TagReference: tagsReference,
		Labels:       getLabelsFromStatus(cdsStatus),
	}

	if _, err := getClient().MessageAdd(msg); err != nil {
		log.Errorf("Error while MessageAdd:%s", err)
	}

}

func getLabelsFromStatus(status sdk.Status) []tat.Label {
	switch status {
	case sdk.StatusSuccess:
		return []tat.Label{tat.Label{Text: status.String(), Color: success}}
	case sdk.StatusWaiting:
		return []tat.Label{tat.Label{Text: status.String(), Color: waiting}}
	case sdk.StatusBuilding:
		return []tat.Label{tat.Label{Text: status.String(), Color: building}}
	case sdk.StatusFail:
		return []tat.Label{tat.Label{Text: status.String(), Color: failed}}
	default:
		return []tat.Label{tat.Label{Text: status.String(), Color: unknown}}
	}
}
