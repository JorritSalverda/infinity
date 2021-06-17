# Infinity

Infinity is a CLI to easily build your applications using a pipeline as code

# Install

## From source

```
go install github.com/JorritSalverda/infinity
```

## With Homebrew

First install Homebrew:

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Then install the `infinity` cli with

```
brew install jorritsalverda/core/infinity
```

# Usage


## Scaffolding a new application build manifest

In order to create a `.infinity.yaml` build template run the following:

```
infinity scaffold <template> <application name>
```

This could be use like:

```
infinity scaffold web myapp
```

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
