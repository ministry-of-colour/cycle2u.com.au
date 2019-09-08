package main

//go:generate rice embed-go

import (
	"fmt"
	"github.com/steveoc64/memdebug"

	"github.com/sirupsen/logrus"
)

func main() {
	memdebug.GCMode(false)
	log := logrus.New()
	cfg, err := initConfig(log)
	if err != nil {
		panic(fmt.Errorf("Opening config: %s\n", err.Error()))
	}
	web := NewWebHandler(cfg, log)
	web.Run()
}
