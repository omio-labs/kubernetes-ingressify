# ingress-generator-kit

> Under active development

## Development

Use docker/docker-compose to develop. You don't need to have golang installed.

* `docker-compose build` Builds image for development
* `docker-compose run --rm default /bin/bash` Gives you a terminal inside the container, from where you can run go commands like:
  * `go test ./...` Runs all tests
  * `go run main.go` Compiles and runs main
