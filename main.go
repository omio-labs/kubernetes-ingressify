package main

import (
	"github.com/apex/log"
  "github.com/goeuro/ingress-generator-kit/core"
  "strings"
)

func main() {
  data, err := core.Asset("gen/version")
  if err != nil {
    log.WithError(err).Error("error")
  }

  version := strings.TrimSpace(string(data))
	log.Infof("ingress-generator-kit version %s", version)
}
