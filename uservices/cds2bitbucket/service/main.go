package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VERSION ...
const VERSION = "0.1.0"

var mainCmd = &cobra.Command{
	Use: "cds2bitbucket",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("cds")
		viper.AutomaticEnv()

		switch viper.GetString("log_level") {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "error":
			log.SetLevel(log.WarnLevel)
			gin.SetMode(gin.ReleaseMode)
		default:
			log.SetLevel(log.DebugLevel)
		}

		router := gin.New()
		router.Use(gin.Recovery())

		router.Use(cors.Middleware(cors.Config{
			Origins:         "*",
			Methods:         "GET, PUT, POST, DELETE",
			RequestHeaders:  "Origin, Authorization, Content-Type, Accept",
			MaxAge:          50 * time.Second,
			Credentials:     true,
			ValidateHeaders: false,
		}))

		router.GET("/mon/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"version": VERSION})
		})

		s := &http.Server{
			Addr:           ":" + viper.GetString("listen_port"),
			Handler:        router,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Infof("Running cds2bitbucket on %s", viper.GetString("listen_port"))

		go do()

		if err := s.ListenAndServe(); err != nil {
			log.Errorf("Error while running ListenAndServe: %s", err.Error())
		}
	},
}

func init() {
	flags := mainCmd.Flags()

	flags.String("log-level", "", "Log Level : debug, info or warn")
	viper.BindPFlag("log_level", flags.Lookup("log-level"))

	flags.String("listen-port", "8085", "Listen Port")
	viper.BindPFlag("listen_port", flags.Lookup("listen-port"))

	flags.String("url-bitbucket", "", "URL Of bitbucket")
	viper.BindPFlag("url_bitbucket", flags.Lookup("url-bitbucket"))

	flags.String("user-bitbucket", "", "User Of bitbucket")
	viper.BindPFlag("user_bitbucket", flags.Lookup("user-bitbucket"))

	flags.String("password-bitbucket", "", "Password Of bitbucket")
	viper.BindPFlag("password_bitbucket", flags.Lookup("password-bitbucket"))

	flags.String("event-kafka-broker-addresses", "", "Ex: --event-kafka-broker-addresses=host:port,host2:port2")
	viper.BindPFlag("event_kafka_broker_addresses", flags.Lookup("event-kafka-broker-addresses"))

	flags.String("event-kafka-topic", "", "Ex: --kafka-topic=your-kafka-topic")
	viper.BindPFlag("event_kafka_topic", flags.Lookup("event-kafka-topic"))

	flags.String("event-kafka-user", "", "Ex: --kafka-user=your-kafka-user")
	viper.BindPFlag("event_kafka_user", flags.Lookup("event-kafka-user"))

	flags.String("event-kafka-password", "", "Ex: --kafka-password=your-kafka-password")
	viper.BindPFlag("event_kafka_password", flags.Lookup("event-kafka-password"))

	flags.String("event-kafka-group", "", "Ex: --kafka-group=your-kafka-group")
	viper.BindPFlag("event_kafka_group", flags.Lookup("event-kafka-group"))
}

func main() {
	mainCmd.Execute()
}
