package main

import (
	"testing"
)

func TestExecHook(t *testing.T) {
	msg := "hello world"
	cmd := []string{"/bin/echo", "-n", msg}
	str, err := ExecHook(cmd)
	if err != nil {
		t.Errorf("Failed to execute echo command")
	}
	if str != "hello world" {
		t.Errorf("ExecHook did not return output of command, got: %s, expected: %s", str, msg)
	}
}
