package main

import (
	"github.com/apex/log"
	"github.com/goeuro/kubernetes-ingressify/core"
	"strings"
)

func main() {
	data, err := core.Asset("gen/version")
	if err != nil {
		log.WithError(err).Error("error")
	}

	version := strings.TrimSpace(string(data))
	log.Infof("kubernetes-ingressify version %s", version)
}
