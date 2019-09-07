package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type configMap map[string]interface{}

func initConfig() (configMap, error) {
	viper.AddConfigPath("$HOME/.cycle2u.com.au")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		println("config changed", e.Name)
	})
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return viper.AllSettings(), nil
}
