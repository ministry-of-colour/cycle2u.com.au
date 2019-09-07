package main

import (
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/gomail"
)

func initMail(cfg configMap, log *logrus.Logger) {
	m, ok := cfg["mail"]
	if !ok {
		return
	}
	mon, ok := cfg["monitor"]
	if !ok {
		return
	}

	switch m.(type) {
	case configMap:
		mm := m.(configMap)
		monm := mon.(configMap)
		log.WithField("monitor", monm["email"]).Info("Sending email on boot")
		mailer := gomail.New(mm["server"].(string), mm["username"].(string), mm["password"].(string))
		mailer.Send(mm["from"].(string), monm["email"].(string), "Cycle2U App Startup", "booting ...")
	}
}
