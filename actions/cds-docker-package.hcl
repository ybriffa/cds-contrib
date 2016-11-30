
name = "CDS_DockerPackage"
description = "Build image and push it to docker repository"

// Requirements
requirements = {
	"docker" = {
		type = "binary"
		value = "docker"
	}
	"bash" = {
		type = "binary"
		value = "bash"
	}
}

// Parameters
parameters = {
	 "dockerfileDirectory" = {
		type = "string"
		description = "Directory which contains your Dockerfile."
	}
	"dockerOpts" = {
		type = "string"
		value = ""
		description = "Docker options, Enter --no-cache --pull if you want for example"
	}
	"dockerRegistry" = {
		type = "string"
		description = "Docker Registry. Enter myregistry for build image myregistry/myimage:mytag"
	}
	"imageName" = {
		type = "string"
		description = "Name of your docker image, without tag. Enter myimage for build image myregistry/myimage:mytag"
	}
	"imageTag" = {
		type = "string"
		description = "Tag og your docker image.
Enter mytag for build image myregistry/myimage:mytag. {{.cds.version}} is a good tag from CDS.
You can use many tags: firstTag,SecondTag
Example : {{.cds.version}},latest"
		value = "{{.cds.version}}"
	}
	"dockerPush" = {
		type = "boolean"
		description = "Docker push built image?"
		value = "true"
	}
	"dockerRMi" = {
		type = "boolean"
		description = "docker rmi built image?"
		value = "false"
	}
}

// Steps
steps = [{
	script = <<EOF
#!/bin/bash
set -e

IMG=`echo {{.imageName}}| tr '[:upper:]' '[:lower:]'`
TAG=`echo {{.imageTag}} | sed 's/\///g'`
echo "Building ${IMG}:${TAG}"

cd {{.dockerfileDirectory}}
docker build {{.dockerOpts}} -t {{.dockerRegistry}}/$IMG:$TAG .

EOF
	}, {
		script = <<EOF
#!/bin/bash

if [ "xtrue" != "x{{.dockerPush}}" ]; then
	echo "docker push: {{.dockerPush}}. So, no image pushed... "
	exit 0
fi;

IFS=', ' read -r -a tags <<< "$string"

for t in "${tags[@]}"; do

	IMG=`echo {{.imageName}}| tr '[:upper:]' '[:lower:]'`
	TAG=`echo ${t} | sed 's/\///g'`

	echo "Pushing {{.dockerRegistry}}/$IMG:$TAG"
	docker push {{.dockerRegistry}}/$IMG:$TAG

	if [ $? -ne 0 ]; then
		set -e
		echo "/!\ Error while pushing to repository. Automatic retry in 60s..."
	    sleep 60
	    docker push {{.dockerRegistry}}/$IMG:$TAG
	fi

	set -e
	echo " {{.dockerRegistry}}/$IMG:$TAG is pushed"

	if [ "xtrue" == "x{{.dockerRMi}}" ]; then
		docker rmi -f {{.dockerRegistry}}/$IMG:$TAG || true;
	else
		echo "docker rmi: {{.dockerRMi}}. So, built image is not deleted"
	fi;

done

EOF
	}]
