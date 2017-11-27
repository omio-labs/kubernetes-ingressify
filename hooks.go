package main

import (
	"github.com/apex/log"
	"os/exec"
)

// Hook is a struct that contains the pre-render and post-render scripts to be executed
type Hook struct {
	PreRender  []string `json:"pre_render"`
	PostRender []string `json:"post_render"`
}

// ExecHook executes an array of commands
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
