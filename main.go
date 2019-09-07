package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/gomail"
)

func main() {
	log := logrus.New()
	cfg, err := initConfig()
	if err != nil {
		panic(fmt.Errorf("Opening config: %s\n", err.Error()))
	}
	initMail(cfg, log)
	spew.Dump(cfg)
}
