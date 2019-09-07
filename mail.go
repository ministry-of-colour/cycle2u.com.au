package main

import (
	"github.com/sirupsen/logrus"
	"github.com/steveoc64/gomail"
)

func sendMail(subject, msg string, cfg *configData, log *logrus.Logger) error {
	if cfg.Mail.Server == "" {
		log.Info("No mail server configured")
		return nil
	}
	if cfg.Monitor.Email == "" {
		log.Info("No monitor address configured")
		return nil
	}
	mailer := gomail.New(cfg.Mail.Server, cfg.Mail.Username, cfg.Mail.Password)
	return mailer.Send(cfg.Mail.From, cfg.Monitor.Email, subject, msg)
}
