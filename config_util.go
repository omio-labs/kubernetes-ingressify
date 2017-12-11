package main

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"time"
)

// Config represents the structure of the config file
type Config struct {
	Kubeconfig      string `json:"kubeconfig"`
	Interval        string `json:"interval"`
	InTemplate      string `json:"in_template"`
	OutTemplate     string `json:"out_file"`
	HealthCheckPort uint32 `json:"health_check_port"`
	Hooks           Hook   `json:"hooks"`
}

func (c Config) getInterval() (time.Duration, error) {
	return time.ParseDuration(c.Interval)
}

// ReadConfig is a helper function to read the config
func ReadConfig(path string) Config {
	var config Config
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(dat, &config)
	if err != nil {
		panic(err.Error())
	}
	return config
}
