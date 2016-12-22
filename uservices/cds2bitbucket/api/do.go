package main

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/bsm/sarama-cluster.v2"

	"github.com/ovh/cds/sdk"
)

func do() {

	var config = sarama.NewConfig()
	config.Net.TLS.Enable = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = viper.GetString("event_kafka_user")
	config.Net.SASL.Password = viper.GetString("event_kafka_password")
	config.Version = sarama.V0_10_0_1

	config.ClientID = viper.GetString("event_kafka_user")

	clusterConfig := cluster.NewConfig()
	clusterConfig.Config = *config
	clusterConfig.Consumer.Return.Errors = true

	var errConsumer error
	consumer, errConsumer := cluster.NewConsumer(
		[]string{viper.GetString("event_kafka_broker_addresses")},
		viper.GetString("event_kafka_group"),
		[]string{viper.GetString("event_kafka_topic")},
		clusterConfig)

	if errConsumer != nil {
		log.Fatalf("Error creating consumer: %s", errConsumer)
	}

	// Consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Errorf("Error during consumption: %s", err)
		}
	}()

	log.Info("Ready to consume messages...")
	for msg := range consumer.Messages() {
		var event sdk.Event
		json.Unmarshal(msg.Value, &event)
		log.Debugf("Receive: type:%s all: %+v", event.EventType, event)
		if event.EventType == fmt.Sprintf("%T", sdk.EventPipelineBuild{}) {
			var eventpb sdk.EventPipelineBuild
			if err := mapstructure.Decode(event.Payload, &eventpb); err != nil {
				log.Errorf("Error during consumption: %s", err)
			} else {
				process(&eventpb)
			}
		}
	}
	return
}
