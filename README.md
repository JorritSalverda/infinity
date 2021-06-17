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

In order to create a `.infinity.yaml` build template run the following:

```
infinity scaffold [template name] [application name]
```

This could be use like:

```
infinity scaffold golang myapp
```

After running this a `.infinity.yaml` manifest will be generated in the current working directory.

You can find a number of templates at https://github.com/JorritSalverda/infinity/tree/main/templates.

## Validate a application build manifest

Once a `.infinity.yaml` manifest exist in the current directory it can be validated with:

```
infinity validate
```

## Build an application locally

The build stages in the `.infinity.yaml` manifest can be executed with:

```
infinity build
```

This will run each stage in a docker container into which the current directory gets mounted, so you can build, test and release your applications in a repeatable fashion.

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