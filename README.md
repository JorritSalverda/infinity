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
