#!/usr/bin/env bash
set -xe

if [ "$TRAVIS_BRANCH" != "master" ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
  echo "Not on master branch, skipping latest release"
  exit 0
fi

VERSION="v0.0.1-snapshot"
git tag -d $VERSION || true
git tag $VERSION
git push https://${GH_TOKEN}:x-oauth-basic@github.com/goeuro/kubernetes-ingressify.git $VERSION -f
