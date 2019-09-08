package main

//go:generate rice embed-go

import (
	"fmt"
	"os"

	"github.com/steveoc64/memdebug"

	"github.com/sirupsen/logrus"
)

func main() {
	memdebug.GCMode(false)
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	cfg, err := initConfig(log)
	if err != nil {
		panic(fmt.Errorf("Opening config: %s\n", err.Error()))
	}
	web := NewWebHandler(cfg, log)
	web.Run()
}
