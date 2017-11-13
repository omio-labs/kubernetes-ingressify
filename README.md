# ingress-generator-kit [![Latest Release](https://img.shields.io/github/release/goeuro/ingress-generator-kit.svg)](https://github.com/goeuro/ingress-generator-kit/releases/latest) [![Build Status](https://travis-ci.org/goeuro/ingress-generator-kit.svg?branch=travis-ci)](https://travis-ci.org/goeuro/ingress-generator-kit) [![Go Report Card](https://goreportcard.com/badge/github.com/goeuro/ingress-generator-kit)](https://goreportcard.com/report/github.com/goeuro/ingress-generator-kit) [![codecov](https://codecov.io/gh/goeuro/ingress-generator-kit/branch/master/graph/badge.svg)](https://codecov.io/gh/goeuro/ingress-generator-kit)

> Under active development

## Development

Use docker/docker-compose to develop. You don't need to have golang installed.

* `docker-compose build` Builds image for development
* `docker-compose run --rm default /bin/bash` Gives you a terminal inside the container, from where you can run go commands like:
  * `bin/test.sh` Runs all tests
  * `gofmt -s -w .` Fix code formatting
  * `go run main.go` Compiles and runs main
