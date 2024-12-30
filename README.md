# Ilo

[![Go](https://github.com/fourls/ilo/actions/workflows/build.yml/badge.svg)](https://github.com/fourls/ilo/actions/workflows/build.yml)
[![Ilo Flows](https://github.com/fourls/ilo/actions/workflows/ilo-build.yml/badge.svg)](https://github.com/fourls/ilo/actions/workflows/ilo-build.yml)

Ilo is a lightweight task runner and automation server that can be run either
locally or as part of a CI process. It is designed to be a thin wrapper over
whatever commands need to be executed, without the layers of cruft and configuration
most other automation servers demand. 

## Development status

Ilo is unstable, in heavy development, and is not yet suitable for general use.

### Features

- [x] Command line task runner
  - [x] Read flows from `ilo.yml` files
    - [x] `run` flow step runs an arbitrary command
    - [x] `echo` flow step prints to the console
      - [ ] Support environment variable substitutions
    - [ ] Specify environment variables for projects, flows, and steps
  - [x] Run flows on-demand with `ilocli run`
  - [x] Register programs by name for use within flows with `ilocli tool add`
    - [ ] Register programs by name and version for use within flows
- [ ] Local automation server to schedule and run flows intermittently
- [ ] Local web interface to view projects, flows, and recent execution information

## Basic Usage

Firstly, create an `ilo.yml` and fill out your project definition. An example `ilo.yml`
for a Go project can be seen below:

```yaml
name: My Go project
flows:
  test:
    - echo: Starting tests
    - run: go test ./...
    - echo: Finished running tests
  build:
    - run: go build ./...
    - run: bash -c 'echo "Finished build at $(date)"'
```

These flows can then be executed by running `ilocli run <flow>` in the same directory.

## Examples

This repository uses Ilo for its continuous integration:

- [Project ilo.yml](ilo.yml)
- [Ilo Flows workflow on GitHub Actions](.github/workflows/ilo-build.yml)
- [Workflow runs using Ilo](https://github.com/fourls/ilo/actions/workflows/ilo-build.yml)