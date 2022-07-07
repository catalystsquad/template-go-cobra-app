/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"net/http"

	configutils "github.com/catalystsquad/app-utils-go/config"
	"github.com/catalystsquad/app-utils-go/logging"
	"github.com/nozzle/e"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs an example service",
	Long:  `Runs an example service with an example health check.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := initRunCmdConfig()
		maybeStartHealthServer(config)

		// TODO impelement a real server
		runExampleServer(config)
	},
}

type runCmdConfig struct {
	ExampleServerPort int    `json:"example_server_port"`
	EnableHealthCheck bool   `json:"enable_health_check"`
	HealthCheckPath   string `json:"health_check_path"`
	HealthCheckPort   int    `json:"health_check_port"`
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().Int("example_server_port", 8080, "port for example http server")
	runCmd.PersistentFlags().Bool("enable_health_check", true, "when true, runs an http server on port 6000 that can be used for a health check for things like kubernetes with GET /health")
	runCmd.PersistentFlags().String("health_check_path", "/health", "path to serve health check on when health check is enabled")
	runCmd.PersistentFlags().Int("health_check_port", 6000, "port to serve health check on when health check is enabled")

	// bind flags
	err := viper.BindPFlags(runCmd.PersistentFlags())
	// die on error
	if err != nil {
		panic(e.Wrap(err, e.Msg("error initializing configuration")))
	}
}

func initRunCmdConfig() *runCmdConfig {
	// instantiate config struct
	config := &runCmdConfig{}
	// get config from viper, including struct validation. This lets us lean on
	// viper for env var, or config file configuration with ease
	err := configutils.GetConfigFromViper(config)
	if err != nil {
		// validation error, log an error and exit
		logging.Log.WithError(err).Error("configuration is invalid")
		os.Exit(1)
	}
	return config
}

func maybeStartHealthServer(config *runCmdConfig) {
	if config.EnableHealthCheck {
		// start health server in the background
		go func() {
			http.HandleFunc(config.HealthCheckPath, func(writer http.ResponseWriter, request *http.Request) {})
			address := fmt.Sprintf(":%d", config.HealthCheckPort)
			logging.Log.WithFields(logrus.Fields{"address": address, "path": config.HealthCheckPath}).Info("starting health server")
			err := http.ListenAndServe(address, nil)
			if err != nil {
				logging.Log.WithError(e.Wrap(err)).Error("error running health server")
			}
		}()
	}
}

func runExampleServer(config *runCmdConfig) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	address := fmt.Sprintf(":%d", config.ExampleServerPort)
	logging.Log.WithFields(logrus.Fields{"address": address, "path": "/"}).Info("starting example server")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		logging.Log.WithError(e.Wrap(err)).Error("error running example server")
	}
}
