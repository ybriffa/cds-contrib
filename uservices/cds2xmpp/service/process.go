package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-xmpp"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/event"
)

var cdsbot *botClient

const resource = "cds"

type botClient struct {
	creation   time.Time
	XMPPClient *xmpp.Client
}

func born() error {
	xClient, err := getNewXMPPClient()
	if err != nil {
		return fmt.Errorf("getClient >> error with getNewXMPPClient err:%s", err)
	}

	cdsbot = &botClient{
		creation:   time.Now(),
		XMPPClient: xClient,
	}

	go cdsbot.receive()

	return nil
}

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

func process(event sdk.Event) error {
	var eventNotif sdk.EventNotif
	log.Debugf("process> receive: type:%s", event.EventType)

	// skip all event != eventNotif
	if event.EventType != fmt.Sprintf("%T", sdk.EventNotif{}) {
		log.Debugf("process> receive: type:%s - skipped", event.EventType)
		return nil
	}

	if err := mapstructure.Decode(event.Payload, &eventNotif); err != nil {
		log.Warnf("process> Error during consumption. type:%s err:%s", event.EventType, err)
		return nil
	}

	log.Debugf("process> event:%+v", event)

	for _, r := range eventNotif.Recipients {
		if !strings.Contains(r, "@") {
			r += "@" + viper.GetString("xmpp_default_hostname")
		}
		log.Debugf("process> event send to :%s", r)
		cdsbot.XMPPClient.Send(xmpp.Chat{
			Remote: r,
			Type:   "chat",
			Text:   eventNotif.Subject + " " + eventNotif.Body,
		})
	}

	return nil
}

func (bot *botClient) receive() {
	for {
		chat, err := bot.XMPPClient.Recv()
		if err != nil {
			log.Errorf("receive >> err: %s", err)
		}
		switch v := chat.(type) {
		case xmpp.Chat:
			log.Debugf("receive> msg from xmpp :%+v\n", v)
		}
	}
}
