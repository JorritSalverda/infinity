# Infinity

Infinity is a CLI to easily build your applications using a _pipeline as code_ approach. It uses an `.infinity.yaml` manifest inside a code repository that specifies the build time dependencies and commands to execute. The _infinity_ tool can execute this manifest locally, so you can build an application without needing all build time dependencies on your machine, only `docker` and `infinity`.

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

## Scaffolding a new application build manifest

In order to create an `.infinity.yaml` build template run the following:

```
infinity scaffold [template name] [application name]
```

This could be used like:

```
infinity scaffold golang myapp
```

After running this the manifest will be generated in the current working directory.

You can find a number of templates at https://github.com/JorritSalverda/infinity/tree/main/templates.

## Validate a application build manifest

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

![Build output](https://github.com/JorritSalverda/infinity/blob/main/screenshot.png?raw=true)

### Mounts and privileged mode

To run some more advanced use cases you can set `privileged: true` on a stage, and multiple mounts. This allows you for example to let _infinity_ build a dockerfile in the following manner:

```yaml
  - name: bake
    image: docker:20.10.7
    privileged: true
    mounts:
    - /var/run/docker.sock:/var/run/docker.sock
    commands:
    - docker build -t web:local .
```

With this example the `docker build` command actually uses your hosts Docker daemon. In order to build the container in isolation you can use

```yaml
  - name: bake
    image: docker:20.10.7-dind
    privileged: true
    commands:
    - ( dockerd-entrypoint.sh & )
    - ( while true ; do if [ -S /var/run/docker.sock ] ; then break ; fi ; sleep 3 ; done )
    - docker build -t web:local .
```

You can also use it to mount devices and in that way allow stages to control some connected hardware:

```yaml
  - name: test
    image: alpine:3.13
    privileged: true
    mounts:
    - /dev/ttyUSB0:/dev/ttyUSB0
    commands:
    # this runs forever, but shows serial usb port output
    - cat /dev/ttyUSB0
```

### Bare metal

In the exceptional case that a command can't run inside a Docker container a stage can be run with `bareMetal: true`; this runs the specified commands directly on the host operating system. The drawback of using this mode is that the build time dependencies either need to be preinstalled or get installed using the commands, leaving them behind on the host.

```yaml
  - name: upload
    bareMetal: true
    commands:
    - apt-get update && apt-get install -y curl
    - curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | sh -s 0.18.3
    - arduino-cli core install arduino:avr 
    - arduino-cli board list
    - arduino-cli core list
    - arduino-cli compile -b arduino:avr:uno sketches/blink
    - arduino-cli upload -b arduino:avr:uno -p /dev/cu.usbserial-1460 sketches/blink
```

### Parallel stages

Regular stages run sequentially, but in order to speed up things you can run stages in parallel by nesting them inside a named containing stage:

```yaml
  - name: parallel
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