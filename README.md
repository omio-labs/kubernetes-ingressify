# ingress-generator-kit

[![Latest Release](https://img.shields.io/github/release/goeuro/ingress-generator-kit.svg)](https://github.com/goeuro/ingress-generator-kit/releases/latest)
[![Build Status](https://img.shields.io/travis/goeuro/ingress-generator-kit.svg?label=master)](https://travis-ci.org/goeuro/ingress-generator-kit)
[![Release Status](https://img.shields.io/travis/goeuro/ingress-generator-kit/v0.0.1-snapshot.svg?label=release)](https://travis-ci.org/goeuro/ingress-generator-kit/branches)
[![Go Report Card](https://goreportcard.com/badge/github.com/goeuro/ingress-generator-kit)](https://goreportcard.com/report/github.com/goeuro/ingress-generator-kit)
[![codecov](https://codecov.io/gh/goeuro/ingress-generator-kit/branch/master/graph/badge.svg)](https://codecov.io/gh/goeuro/ingress-generator-kit)

This is a simple binary that watches kubernetes ingress rules, renders your template and calls your script.
You can use this to generate ingress-based configuration for any backend router.

## Why?

There are multiple kubernetes ingress controller implementations.
Unfortunately they control both router implementation (how router is compiled, built, etc) and behavior.
Its not easy to extend them and add custom logic, e.g. adding custom nginx modules/plugins, custom annotations/features, overriding templates, etc.
This is an attempt to decouple router implementation from configuration and allow anyone to create ingress-aware routers.

## Current status

- [x] Bootstrap
- [ ] >> Design
- [ ] Implementation+Dogfooding
- [ ] Documentation and examples
- [ ] Release v0.1

## Usage

We have a special release called `v0.0.1-snapshot` which always reflects master build.
Please go ahead and grab it: https://github.com/goeuro/ingress-generator-kit/releases/tag/v0.0.1-snapshot

Create a configuration file:

```
# ingress.cfg
kubeconfig: <path to kubeconfig, leave it empty for in-cluster authentication>
in-template: <path to template, context provided to template will be documented, defaults to ingress.cfg.tpl>
out-file: <path to output file, defaults to ingress.cfg>
hooks:
  post-render: <path to script that will be called after template is rendered, e.g. reload nginx/haproxy/etc>
```

Run it:

```
ingress-generator-kit -f ingress.cfg
```

## Development

Use docker/docker-compose to develop. You don't need to have golang installed.

* `docker-compose build` Builds image for development
* `docker-compose run --rm default /bin/bash` Gives you a terminal inside the container, from where you can run go commands like:
  * `bin/test.sh` Runs all tests
  * `gofmt -s -w .` Fix code formatting
  * `go run main.go` Compiles and runs main
* Adding dependencies:
  * catch-22: godep has issues, dep is alpha (not all libs support it) and glide is deprecated in favor of dep
  * We use Godep as the least working option, but it means slightly additional effort when adding dependencies
  * SSH into your container (as above) and follow https://github.com/tools/godep#add-a-dependency
