package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ovh/tat"
	"github.com/spf13/viper"
)

var instance *tat.Client

// getClient initializes client on tat engine
func getClient() *tat.Client {
	if instance != nil {
		return instance
	}

	tc, err := tat.NewClient(tat.Options{
		URL:      viper.GetString("url_tat_engine"),
		Username: viper.GetString("username_tat_engine"),
		Password: viper.GetString("password_tat_engine"),
		Referer:  "tatexamplecron.v." + VERSION,
	})

	if err != nil {
		log.Errorf("Error while create new Tat Client:%s", err)
	}

	switch viper.GetString("log_level") {
	case "debug":
		tat.DebugLogFunc = log.Debugf
	case "info":
		tat.DebugLogFunc = log.Infof
	case "error":
		tat.DebugLogFunc = log.Warnf
	default:
		tat.DebugLogFunc = log.Debugf
	}

	instance = tc
	return instance
}
