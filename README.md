# Infinity

Infinity is a CLI to easily build your applications using a pipeline as code

# Install

## With Homebrew

First install Homebrew:

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
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