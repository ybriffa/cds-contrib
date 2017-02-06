# group-tmpl
plugin-group-tmpl is an helper to create a deployment file for an application group on mesos/marathon using golang text/template
Once the file is generated, you can simply push the file to mesos/marathon as a group update.

## How to build

Make sure go is installed and properly configured ($GOPATH must be set)

```shell
$ go test ./...
$ go build
```

## How to install

### Install as CDS plugin

As CDS admin:

- login to Web Interface,
- go to Action admin page,
- choose the plugin binary file freshly built
- click on "Add plugin"

## How to use

### Parameters

- **config**, a template marathon config to apply to every application
- **applications**, a variables file to overwrite the variables in the template file, containing both default variables to have an overall default configuration, and apps variables to be able to configure apps by apps
- **output**, generated file, takes by the default <config>.out or just trimming the .tpl extension



## Examples

Template file
```json
{
        "id": "{{.id}}",
        "image": "{{.image}}",
        "instances": {{.instances}},
        "cpus": {{.cpus}},
        "mem": {{.mem}}
}

```

Applications file :
```json
{
    "default": {
        "mem":"512",
        "cpus":"0.5",
        "image": "docker.registry/my-awesome-image",
        "instances": 1
    },
    "apps": {
        "first": {
        },
        "second": {
            "cpus": "3",
            "mem": 2048
        }
    }
}
```

Then the output file will be (also see §Alterations)
```json
{
  "apps": [
    {
       "id": "first",
       "image": "docker.registry/my-awesome-image",
       "cpus": 1,
       "instances": 1,
       "mem": 512
     }, {
       "id": "second",
       "image": "docker.registry/my-awesome-image",
       "cpus": 3,
       "instances": 1,
       "mem": 2048
     }
 ]
}
```

## Alterations

Somes keys in the applications variables can be altered by functions. This is the list of the current alterations

```json
{
        "id": the value is the map apps' keys
}

```