package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/manudelca/tp1-distribuidos1/test-client/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix("cli")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Add env variables supported
	v.BindEnv("id")
	v.BindEnv("server", "address")
	v.BindEnv("loop", "period")
	v.BindEnv("loop", "lapse")
	v.BindEnv("log", "level")
	v.BindEnv("type")

	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Configuration could not be read from config file. Using env variables instead")
	}

	if _, err := time.ParseDuration(v.GetString("loop.lapse")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_LAPSE env var as time.Duration.")
	}

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_PERIOD env var as time.Duration.")
	}

	return v, nil
}

func InitLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)
	return nil
}

func PrintConfig(v *viper.Viper) {
	logrus.Infof("Client configuration")
	logrus.Infof("Client ID: %s", v.GetString("id"))
	logrus.Infof("Server Address: %s", v.GetString("server.address"))
	logrus.Infof("Loop Lapse: %v", v.GetDuration("loop.lapse"))
	logrus.Infof("Loop Period: %v", v.GetDuration("loop.period"))
	logrus.Infof("Log Level: %s", v.GetString("log.level"))
	logrus.Info("Client type: %s", v.GetString("type"))
}

func main() {
	v, err := InitConfig()
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := InitLogger(v.GetString("log.level")); err != nil {
		log.Fatalf("%s", err)
	}

	// Print program config with debugging purposes
	PrintConfig(v)

	clientConfig := common.ClientConfig{
		ServerAddress: v.GetString("server.address"),
		ID:            v.GetString("id"),
		LoopLapse:     v.GetDuration("loop.lapse"),
		LoopPeriod:    v.GetDuration("loop.period"),
		ClientType:    v.GetString("type"),
	}

	client := common.NewClient(clientConfig)
	client.StartClientLoop()
}
