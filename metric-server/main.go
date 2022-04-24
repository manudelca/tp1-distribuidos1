package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/manudelca/tp1-distribuidos1/metric-server/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() (*viper.Viper, error) {
	v := viper.New()

	// Configure viper to read env variables
	v.AutomaticEnv()
	// Use a replacer to replace env variables underscores with points. This let us
	// use nested configurations in the config file and at the same time define
	// env variables for the nested configurations
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Add env variables supported
	v.BindEnv("couriers")
	v.BindEnv("log", "level")

	// Try to read configuration from config file. If config file
	// does not exists then ReadInConfig will fail but configuration
	// can be loaded from the environment variables so we shouldn't
	// return an error in that case
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from config file. Using env variables instead")
	}

	// Parse int variables and return an error if they cannot be parsed
	if _, err := strconv.Atoi(v.GetString("couriers")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse couriers env var as int.")
	}
	if _, err := strconv.Atoi(v.GetString("metricEventsBacklog")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse metricEventsBacklog env var as int.")
	}
	if _, err := strconv.Atoi(v.GetString("queryEventsBacklog")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse queryEventsBacklog env var as int.")
	}
	if _, err := strconv.Atoi(v.GetString("metricEventsWorkers")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse metricEventsWorkers env var as int.")
	}
	if _, err := strconv.Atoi(v.GetString("queryEventsWorkers")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse queryEventsWorkers env var as int.")
	}

	return v, nil
}

// InitLogger Receives the log level to be set in logrus as a string. This method
// parses the string and set the level to the logger. If the level string is not
// valid an error is returned
func InitLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)
	return nil
}

// PrintConfig Print all the configuration parameters of the program.
// For debugging purposes only
func PrintConfig(v *viper.Viper) {
	logrus.Infof("[MAIN] Metric-server configuration")
	logrus.Infof("[MAIN] Couriers: %s", v.GetString("couriers"))
	logrus.Infof("[MAIN] Port: %s", v.GetString("port"))
	logrus.Infof("[MAIN] Log Level: %s", v.GetString("log.level"))
	logrus.Infof("[MAIN] Metric events backlog: %s", v.GetString("metricEventsBacklog"))
	logrus.Infof("[MAIN] Query events backlog: %s", v.GetString("queryEventsBacklog"))
	logrus.Infof("[MAIN] Metric events workers: %s", v.GetString("metricEventsWorkers"))
	logrus.Infof("[MAIN] Query events workers: %s", v.GetString("queryEventsWorkers"))
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

	serverConfig := common.ServerConfig{
		Port:                v.GetString("port"),
		Couriers:            v.GetInt("couriers"),
		MetricEventsBacklog: v.GetInt("metricEventsBacklog"),
		QueryEventsBacklog:  v.GetInt("queryEventsBacklog"),
		MetricEventsWorkers: v.GetInt("metricEventsWorkers"),
		QueryEventsWorkers:  v.GetInt("queryEventsWorkers"),
	}

	server, err := common.NewServer(serverConfig)
	if err != nil {
		logrus.Fatalf("[MAIN] Could not create Server. Error %s", err)
		return
	}
	server.Run()
}
