package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/event"
)

func do() {
	event.ConsumeKafka(viper.GetString("event_kafka_broker_addresses"),
		viper.GetString("event_kafka_topic"),
		viper.GetString("event_kafka_group"),
		viper.GetString("event_kafka_user"),
		viper.GetString("event_kafka_password"),
		func(e sdk.Event) error {
			return process(e)
		},
		log.Errorf,
	)
}
