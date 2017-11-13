#!/usr/bin/env bash
set -xe

echo 'Running tests...'
overalls -project=github.com/goeuro/ingress-generator-kit -covermode=atomic

echo 'Checking code format...'
(! gofmt -l -s -e . 2>&1 | grep -v 'vendor/' | grep .go) || exit 1

echo 'Linting...'
golint -set_exit_status .
