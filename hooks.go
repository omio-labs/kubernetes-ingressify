package main

import (
	"github.com/apex/log"
	"os/exec"
)

type Hook struct {
	PreRender  []string `json:"pre_render"`
	PostRender []string `json:"post_render"`
}

func ExecHook(hook []string) (string, error) {
	run := exec.Command(hook[0], hook[1:]...)
	log.Info("Executing hook")
	out, err := run.Output()
	if err != nil {
		log.WithError(err).Error("Failed to run hook")
		return "", err
	}
	log.Info("Hook execution successful")
	return string(out), nil
}
