package main

import (
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/template"
)

type TemplatePlain struct {
	template.Common
}

func (t *TemplatePlain) Name() string {
	return "template-plain"
}

func (t *TemplatePlain) Description() string {
	return `
This sample template create
- a build pipeline with	two stages: Compile Stage and Packaging Stage
- a deploy pipeline with one stage: Deploy Stage

Compile Stage :
- run git clone
- run make build

Packaging Stage:
- run docker build and docker push

Deploy Stage:
- it's en empty script

Packaging and Deploy are optional.
`
}

func (t *TemplatePlain) Identifier() string {
	return "github.com/ovh/cds-contrib/templates/template-plain/TemplatePlain"
}

func (t *TemplatePlain) Author() string {
	return "Yvonnick Esnault <yvonnick.esnault@corp.ovh.com>"
}

func (t *TemplatePlain) Type() string {
	return "BUILD"
}

func (t *TemplatePlain) Parameters() []sdk.TemplateParam {
	return []sdk.TemplateParam{
		{
			Name:        "withPackage",
			Type:        sdk.BooleanVariable,
			Value:       "withPackage",
			Description: "Do you want a Docker Package?",
		},
		{
			Name:        "withDeploy",
			Type:        sdk.BooleanVariable,
			Value:       "withDeploy",
			Description: "Do you want an deploy Pipeline?",
		},
	}
}

func (t *TemplatePlain) ActionsNeeded() []string {
	return []string{
		"CDS_GitClone",
	}
}

func (t *TemplatePlain) Apply(opts template.IApplyOptions) (sdk.Application, error) {
	//Return full application

	a := sdk.Application{
		Name:       opts.ApplicationName(),
		ProjectKey: opts.ProjetKey(),
	}

	/* Build Pipeline */
	/* Build Pipeline - Compile Stage */

	jobCompile := sdk.Action{
		Name: "Compile",
		Actions: []sdk.Action{
			sdk.Action{
				Name: "CDS_GitClone",
			},
			sdk.NewActionScript(`#!/bin/bash

set -xe

cd $(ls -1) && make`,

				[]sdk.Requirement{
					{
						Name:  "make",
						Type:  sdk.BinaryRequirement,
						Value: "make",
					},
				},
			),
		},
	}

	compileStage := sdk.Stage{
		Name:       "Compile Stage",
		BuildOrder: 0,
		Enabled:    true,
		Actions:    []sdk.Action{jobCompile},
	}

	/* Build Pipeline - Packaging Stage */

	jobDockerPackage := sdk.Action{
		Name: "Docker package",
		Actions: []sdk.Action{
			sdk.Action{
				Name: "CDS_GitClone",
			},
			sdk.NewActionScript(`#!/bin/bash
set -ex

cd $(ls -1)

docker build -t cds/{{.cds.application}}-{{.cds.version}} .
docker push cds/{{.cds.application}}-{{.cds.version}}`, []sdk.Requirement{
				{
					Name:  "bash",
					Type:  sdk.BinaryRequirement,
					Value: "bash",
				},
			},
			),
		},
	}

	packagingStage := sdk.Stage{
		Name:       "Packaging Stage",
		BuildOrder: 0,
		Enabled:    true,
		Actions:    []sdk.Action{jobDockerPackage},
	}

	/* Deploy Pipeline */
	/* Deploy Pipeline - Deploy Stage */

	jobDeploy := sdk.Action{
		Name: "Deploy",
		Actions: []sdk.Action{
			sdk.NewActionScript(`#!/bin/bash
set -ex

echo "CALL YOUR DEPLOY SCRIPT HERE"`, []sdk.Requirement{
				{
					Name:  "docker",
					Type:  sdk.BinaryRequirement,
					Value: "docker",
				},
			},
			),
		},
	}

	deployStage := sdk.Stage{
		Name:       "Deploy Stage",
		BuildOrder: 0,
		Enabled:    true,
		Actions:    []sdk.Action{jobDeploy},
	}

	/* Assemble Pipeline */

	a.Pipelines = []sdk.ApplicationPipeline{
		{
			Pipeline: sdk.Pipeline{
				Name:   "build",
				Type:   sdk.BuildPipeline,
				Stages: []sdk.Stage{compileStage},
			},
			Triggers: []sdk.PipelineTrigger{
				{
					DestPipeline: sdk.Pipeline{
						Name: "deploy",
					},
					DestEnvironment: sdk.Environment{
						Name: "Production",
					},
				},
			},
		},
	}

	if opts.Parameters().Get("withPackage") == "true" {
		a.Pipelines[0].Pipeline.Stages = append(a.Pipelines[0].Pipeline.Stages, packagingStage)
	}

	if opts.Parameters().Get("withDeploy") == "true" {
		a.Pipelines = append(a.Pipelines,
			sdk.ApplicationPipeline{
				Pipeline: sdk.Pipeline{
					Name:   "deploy",
					Type:   sdk.DeploymentPipeline,
					Stages: []sdk.Stage{deployStage},
				},
			},
		)
	}

	return a, nil
}

func main() {
	p := TemplatePlain{}
	template.Serve(&p)
}
