#!/usr/bin/env bash
set -xe

# tests with coverage
overalls -project=github.com/goeuro/ingress-generator-kit -covermode=atomic

# gofmt
! (gofmt -l -s -e . 2>&1 | grep -v 'vendor/' | grep .go)

# golint
golint -set_exit_status .
