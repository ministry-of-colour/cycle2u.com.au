package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type configData struct {
	Port    int
	Address string
	Name    string
	DBPath  string
	Mail    mailConfig
	Monitor monitorConfig
	SMS     smsConfig
}

type mailConfig struct {
	Server   string
	Username string
	Password string
	From     string
	Email    string
	CC       []string
	BCC      []string
}

type monitorConfig struct {
	Email string
}

type smsConfig struct {
	API         string
	Username    string
	Password    string
	Destination string
	Source      string
}

var configDataData *configData

func initConfig(log *logrus.Logger) (*configData, error) {
	cfg := &configData{}
	configDataData = cfg
	viper.AddConfigPath("$HOME/.cycle2u.com.au")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		println("config changed", e.Name)
		readConfig("Config Changed", cfg, log)
	})
	return cfg, readConfig("Startup", cfg, log)
}

func readConfig(reason string, cfg *configData, log *logrus.Logger) error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}
	return monitorEmail("Cycle2U Startup: "+cfg.Name, fmt.Sprintf(`
<h1>%s: %s</h1>
<ul>
	<li>Port: %d
	<li>Monitor: %s
	<li>From: %s
	<li>Server: %s
	<li>User: %s
	<li>Pass: ********
</ul>
`,
		cfg.Name,
		reason,
		cfg.Port,
		cfg.Monitor.Email,
		cfg.Mail.From,
		cfg.Mail.Server,
		cfg.Mail.Username,
	),
		cfg, log)
}
