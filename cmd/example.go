package cmd

import (
	configutils "github.com/catalystsquad/app-utils-go/config"
	"github.com/catalystsquad/app-utils-go/logging"
	"github.com/nozzle/e"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// ExampleCommandConfig is a struct representation of the config we expect for this command
type ExampleCommandConfig struct {
	Name    string `json:"name" valid:"required"`
	TheName string `json:"the_name" valid:"required"`
}

// exampleCmd is the cobra command configuration for the `example` command
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "An example command",
	Long: `An example command that demonstrates usage of cobra and viper.
Config can be provided by any viper compatible means. Most commonly this is through environment variables.
For example, the two options available for this command are '--name' and '--the_name'. If you set these by exporting
'NAME' and 'THE_NAME' environment variables, those environment variables will be used. Flags given on the command line
take precedence over environment variables.
`,
	// Run is the actual function that executes when the command is called. This is where your business logic should be. Or at least the entry point to it.
	Run: func(cmd *cobra.Command, args []string) {
		// instantiate config struct
		config := &ExampleCommandConfig{}
		// get config from viper, including struct validation. This lets us lean on viper for env var, or config file configuration with ease
		err := configutils.GetConfigFromViper(config)
		if err != nil {
			// validation error, log an error and exit
			logging.Log.WithError(err).Error("configuration is invalid")
			os.Exit(1)
		}
		// got a valid config, do the things, for the example command that's just printing out the config we got
		logging.Log.WithFields(logrus.Fields{"name": config.Name, "the_name": config.TheName}).Info("example command called")
		// ... your implementation below, probably just a function call into more business logic so this function doesn't get super long
	},
}

// init initializes the command and binds the flags to viper
func init() {
	rootCmd.AddCommand(exampleCmd)
	exampleCmd.PersistentFlags().String("name", "", "simple flag name")
	// we recommend using _ in flags instead of - because it makes env var usage and config file usage much easier
	// the_name flag can be assigned via the THE_NAME environment variable
	exampleCmd.PersistentFlags().String("the_name", "", "flag name with underscore")
	flags := exampleCmd.PersistentFlags()
	// bind the flags to viper
	err := viper.BindPFlags(flags)
	// die on error
	if err != nil {
		panic(e.Wrap(err, e.Msg("error initializing configuration")))
	}
}
