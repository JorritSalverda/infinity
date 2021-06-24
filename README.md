# Infinity

Infinity is a CLI to easily build your applications using a _pipeline as code_ approach. It uses an `.infinity.yaml` manifest inside a code repository that specifies the build time dependencies and commands to execute. The _infinity_ tool can execute this manifest locally, so you can build an application without needing all build time dependencies on your machine, only `docker` and `infinity`. The same manifest in combination with the _infinity_ cli can be used inside CI pipelines to build the application, so `works on my machine` also means it `works in the CI pipeline`.

# Install

## With Homebrew

First install Homebrew:

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

And docker:

```
brew install --cask docker
```

Then install the `infinity` cli with

```
brew install jorritsalverda/core/infinity
```

## From source

```
go install github.com/JorritSalverda/infinity
```

# Usage

## Basics

To check what commands are available for the _infinity_ cli run

```
infinity help
```

And to check the installed version run

```
infinity version
```

## Scaffolding a new application build manifest

In order to create an `.infinity.yaml` build template run the following:

```
infinity scaffold [application type] [language] [application name]
```

This could be used like:

```
infinity scaffold library go mylib
```

After running this the manifest will be generated in the current working directory.

You can find a number of templates at https://github.com/JorritSalverda/infinity/tree/main/templates.

## Validate an application build manifest

Once an `.infinity.yaml` manifest exist in the current directory it can be validated with:

```
infinity validate
```

## Build an application locally

The build stages in the `.infinity.yaml` manifest can be executed with:

```
infinity build
```

This will run each stage's commands inside a docker container into which the current directory gets mounted, so you can build, test and release your applications in a repeatable fashion.

Having a pipeline as code gives control over build time dependency to the authors of the application.

An `.infinity.yaml` manifest looks as follows:

```yaml
build:
  stages:
  - name: audit
    image: node:16-alpine
    env:
      npm_config_update-notifier: false
    commands:
    - npm audit
  - name: restore
    image: node:16-alpine
    env:
      npm_config_update-notifier: false
    commands:
    - npm ci
```

When executed with the `infinity build` command it executes the `npm audit` and `npm ci` commands inside a `node:16-alpine` container where the current directory gets mounted to the `/work` directory. The output looks as follows:

![Build output](https://github.com/JorritSalverda/infinity/blob/main/screenshot.jpg?raw=true)

### Volumes, devices and privileged mode

To run some more advanced use cases you can set `privileged: true` on a stage and mount one or more volumes with the `volumes` array. This allows you for example to let _infinity_ build a Dockerfile in the following manner:

```yaml
  - name: bake
    image: docker:20.10.7
    privileged: true
    volumes:
    - /var/run/docker.sock:/var/run/docker.sock
    commands:
    - docker build -t web:local .
```

With this example the `docker build` command actually uses your hosts Docker daemon. 

In order to build the container in isolation you can use the _docker inside docker_ image instead which runs its own Docker daemon.

```yaml
  - name: bake
    image: docker:20.10.7-dind
    privileged: true
    commands:
    - ( dockerd-entrypoint.sh & )
    - ( while true ; do if [ -S /var/run/docker.sock ] ; then break ; fi ; sleep 3 ; done )
    - docker build -t web:local .
```

You can mount devices so commands inside the stage can connect to hardware on the host:

```yaml
  - name: test
    image: alpine:3.13
    devices:
    - /dev/ttyUSB0:/dev/ttyUSB0
    commands:
    # this runs forever, but shows serial usb port output
    - cat /dev/ttyUSB0
```

You can do the same by mounting the devices as volumes, but that needs to be combined with _privileged_ mode:

```yaml
  - name: test
    image: alpine:3.13
    privileged: true
    volumes:
    - /dev/ttyUSB0:/dev/ttyUSB0
    commands:
    # this runs forever, but shows serial usb port output
    - cat /dev/ttyUSB0
```

### Parallel stages

Regular stages run sequentially, but in order to speed up things you can run stages in parallel by nesting them inside a named containing stage:

```yaml
  - name: build-and-lint
    stages:
    - name: lint
      image: golangci/golangci-lint:latest-alpine
      commands:
      - golangci-lint run
    - name: release
      image: golang:1.16-alpine
      env:
        CGO_ENABLED: 0
      commands:
      - go test -short ./...
      - go build -a -installsuffix cgo .
```

Since these stages share the same mounted working directory do ensure they are safe to run concurrently!

### Detached stages

In order to run containers in the background, for example to be used as a service for in-pipeline integration tests you can add `detach: true` to the stage. It will make the container start and then continue to run until all stages are done. Once they're done the _detached_ stage containers will be terminated and their logs shown.

```yaml
  - name: cockroachdb-as-service
    image: cockroachdb/cockroach:v21.1.2
    detach: true
    mount: false
    env:
      COCKROACH_SKIP_ENABLING_DIAGNOSTIC_REPORTING: "true"
    commands:
    - exec /cockroach/cockroach start-single-node --insecure --advertise-addr cockroachdb-as-service
  - name: wait
    image: alpine:3.13
    commands:
    # run schema updates and then integration tests against the database
    - sleep 20s
```

Do note `mount: false` in order to prevent the working directory from getting mounted; in this particular instance the _cockroachdb_ container doesn't need access to any of the files in the working directory.

### Bare metal runner

In the exceptional case that a command can't run inside a Docker container a stage can be run with `runner: metal`; this runs the specified commands directly on the host operating system. The drawback of using this mode is that the build time dependencies either need to be preinstalled or get installed using the commands, leaving them behind on the host.

```yaml
  - name: upload
    runner: metal
    commands:
    - apt-get update && apt-get install -y curl
    - curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | sh -s 0.18.3
    - arduino-cli core install arduino:avr
    - arduino-cli board list
    - arduino-cli core list
    - arduino-cli compile -b arduino:avr:uno sketches/blink
    - arduino-cli upload -b arduino:avr:uno -p /dev/cu.usbserial-1460 sketches/blink
```

### Stage parameters

To make intermediate containers that are more friendly to be used than by passing commands you can set any property - outside of the reserved ones - and they will be passed on as environment variables in the form of `INFINITY_PARAMETER_<UPPER_SNAKE_CASE_VERSION_OF_PARAMETER_NAME>`.

As an example a container could be created that supports the `action`, `container` and `tag` parameters by using environment variable `INFINITY_PARAMETER_ACTION`, `INFINITY_PARAMETER_CONTAINER` and `INFINITY_PARAMETER_TAG` inside the docker container.

```yaml
  - name: build-docker-container
    image: jsalverda/docker:stable
    action: build
    container: web
    tag: 1.0.0
    privileged: true
```

The `jsalverda/docker:stable` container can be created with `Dockerfile`:

```dockerfile
FROM docker:20.10.7-dind

COPY ./docker-entrypoint.sh /
RUN chmod 500 /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
```

And `docker-entrypoint.sh` that looks something like:

```bash
#!/bin/sh
set -e

# start docker daemon
dockerd-entrypoint.sh &

# wait for docker daemon to be ready
while true ; do if [ -S /var/run/docker.sock ] ; then break ; fi ; sleep 3 ; done

case "${INFINITY_PARAMETER_ACTION}" in
  build)
    docker build -t ${INFINITY_PARAMETER_CONTAINER}:${INFINITY_PARAMETER_TAG} .
  ;;

  push)
    docker push ${INFINITY_PARAMETER_CONTAINER}:${INFINITY_PARAMETER_TAG}
  ;;

  scan)
    docker scan ${INFINITY_PARAMETER_CONTAINER}:${INFINITY_PARAMETER_TAG}
  ;;

  *)
    echo "action ${INFINITY_PARAMETER_ACTION} is not supported; use action: build|push|scan"
    exit 1
esac
```

# Examples

In the `examples` directory you can find the following examples highlighting specific features:

| example     | shows...                                                                                                                                             |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| arduino-uno | ...how to mount devices with the default container runner and how to use the (bare) metal runner                                                     |
| cmake       | ...the use of an intermediate docker builder image with prepared build time dependencies for improved performance                                    |
| db-test     | ...how to use detached stages to provide a service to other stages                                                                                   |
| web         | ...splitting commands with same image into multiples stages for better visibility of time spent in each command; shows how a dockerfile can be built |

# Local development

You can either install `infinity` and run `infinity build`; this will compile and install the cli in the go binary folder. Or just run `go install` instead.

# Release

To create a new release in Github first create a zipped version of the CLI with

```
INFINITY_VERSION=v0.1.8

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-X 'github.com/JorritSalverda/infinity/cmd.version=${INFINITY_VERSION}'" -o infinity-${INFINITY_VERSION}-darwin-amd64
zip infinity-${INFINITY_VERSION}-darwin-amd64.zip infinity-${INFINITY_VERSION}-darwin-amd64
rm -rf infinity-${INFINITY_VERSION}-darwin-amd64
shasum -a 256 infinity-${INFINITY_VERSION}-darwin-amd64.zip
```

Then create a release with the version as tag and release title and add the zip file as attached file. Then update `url`, `sha256`, `version` and `install` in the `Formula/infinity.rb` file in repository [github.com/JorritSalverda/homebrew-core](https://github.com/JorritSalverda/homebrew-core). Once it's pushed to Github you can run `brew upgrade` to get the latest version on your machine.

# Manifest reference

| property                    | description                                                                                                                                                                                                                  | allowed values                           | default     |
| --------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------- | ----------- |
| `application`               | application type metadata for use in a future centralized CI/CD system                                                                                                                                                       | `library\|cli\|firmware\|api\|web`       |             |
| `language`                  | language metadata                                                                                                                                                                                                            | `go\|c\|c++\|java\|csharp\|python\|node` |             |
| `name`                      | unique name for the application                                                                                                                                                                                              | `string`                                 |             |
| `build.stages[].name`       | name for the stage                                                                                                                                                                                                           | `string`                                 |             |
| `build.stages[].runner`     | runner type for the stage                                                                                                                                                                                                    | `container\|metal`                       | `container` |
| `build.stages[].image`      | docker container image path for the image to run the stage commands in                                                                                                                                                       | `string`                                 |             |
| `build.stages[].detach`     | run stage in detached mode, to provide a service in the background                                                                                                                                                           | `true\|false`                            | `false`     |
| `build.stages[].privileged` | run stage in privileged mode, to allow more privileges to the host operating system                                                                                                                                          | `true\|false`                            | `false`     |
| `build.stages[].mount`      | mount the working directory into the stage container                                                                                                                                                                         | `true\|false`                            | `true`      |
| `build.stages[].work`       | directory to which the working copy gets mounted                                                                                                                                                                             | `string`                                 | `/work`     |
| `build.stages[].volumes`    | array of volumes to mount, with source and target folder separated by `:`                                                                                                                                                    | `[]string`                               |             |
| `build.stages[].devices`    | array of devices to mount, with source and target device path separated by `:`                                                                                                                                               | `[]string`                               |             |
| `build.stages[].env`        | map of environment value keys and values to allow setting envvars in a stage                                                                                                                                                 | `map[string]string`                      |             |
| `build.stages[].devices`    | array of commands to execute inside the stage container or on bare metal                                                                                                                                                     | `[]string`                               |             |
| `build.stages[].stages`     | array of nested stages that are executed in parallel to speed up total build time                                                                                                                                            | `[]stage`                                |             |
| `build.stages[].*`          | any other property set on the stage is passed as an environment variable in the form of `INFINITY_PARAMETER_<UPPER_SNAKE_CASE_VERSION_OF_PARAMETER_NAME>` to allow for more friendly configuration of a prepared stage image |                                          |             |