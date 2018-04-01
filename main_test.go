package main

import "testing"

func TestGetEnvVars(t *testing.T) {
	c := getConfig()
	if c.LogLevel != "info" {
		t.Error("Got incorrect env var", c)
	}

	if c.Debug != false {
		t.Error("Got incorrect env var", c)
	}
}
