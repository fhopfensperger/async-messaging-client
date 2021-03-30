/*
Copyright © 2020 Florian Hopfensperger <f.hopfensperger@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var connectionString string
var queueName string

var globalUsage = `A simple command line utility to send and receive AMQP message to / from Azure Service Bus:

Sending strings:
async-messaging-client send "my message" -q queueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
async-messaging-client receive -q queueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."

Using json files:
content test.json: { "key": "value" }
async-messaging-client send -f test.json -q queueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
async-messaging-client receive -q queueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."

You could also use environment variables for defining the queue and connection string:
export QUEUE=myQueueName
export CONNECTION_STRING='Endpoint=sb://host.servicebus.windows.net/'
async-messaging-client send -f test.json
async-messaging-client receive

 `

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "async-messaging-client",
	Short: "Send and receive AMQP message to / from Azure Service Bus",
	Long:  globalUsage,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) { fmt.Println("hallo from cli") },
	//Args: NoArgs,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(globalUsage)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.async-messaging-client.yaml)")
	pf := rootCmd.PersistentFlags()
	pf.StringP("connection-string", "c", "", "Connection String to connecto to Azure Service Bus")
	viper.BindPFlag("connection-string", pf.Lookup("connection-string"))

	pf.StringP("queue", "q", "", "Azure Service Bus Queue name")
	viper.BindPFlag("queue", pf.Lookup("queue"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.SetVersionTemplate(`{{printf "v%s\n" .Version}}`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Err(err).Msg("Could not find homedir")
			os.Exit(1)
		}

		// Search config in home directory with name ".async-messaging-client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".async-messaging-client")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msgf("Using config file:", viper.ConfigFileUsed())
	}

	queueName = viper.GetString("queue")
	connectionString = viper.GetString("connection-string")

}

// NoArgs returns an error if any args are included.
func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.Errorf(
			"%q accepts no arguments\n\nUsage:  %s",
			cmd.CommandPath(),
			cmd.UseLine(),
		)
	}
	return nil
}
