/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"net/http"

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
	Port              int
	EnableHealthCheck bool
	HealthCheckPath   string
	HealthCheckPort   int
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().Int("port", 8080, "port for example http server")
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

	config.Port = viper.GetInt("port")
	config.EnableHealthCheck = viper.GetBool("enable_health_check")
	config.HealthCheckPath = viper.GetString("health_check_path")
	config.HealthCheckPort = viper.GetInt("health_check_port")

	logging.Log.WithField("settings", fmt.Sprintf("%+v", *config)).Debug("viper settings")

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

	address := fmt.Sprintf(":%d", config.Port)
	logging.Log.WithFields(logrus.Fields{"address": address, "path": "/"}).Info("starting example server")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		logging.Log.WithError(e.Wrap(err)).Error("error running example server")
	}
}
