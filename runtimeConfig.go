package main

import "flag"

type RuntimeConfig struct {
	tickerDuration int
	configFileLocation string
}

var runtimeConfig *RuntimeConfig

func init() {
	runtimeConfig = &RuntimeConfig{}
	flag.IntVar(&runtimeConfig.tickerDuration, "t", 10, "ticker time in seconds")
	flag.StringVar(&runtimeConfig.configFileLocation, "c", "config.json", "Configuration file with definitions set up")
	flag.Parse()
}