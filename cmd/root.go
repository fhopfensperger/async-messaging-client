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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var globalUsage = `Command line tool to send and receive message to Google Pub/Sub or a Azure Service instance`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "async-messaging-client",
	Short: "Command line tool to send and receive message to Google Pub/Sub or a Azure Service instance",
	Long:  globalUsage,
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
	rootCmd.SetVersionTemplate(`{{printf "v%s\n" .Version}}`)
}

// initConfig reads in ENV variables if set.
func initConfig() {

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

func transformFileToJson(fileName string) ([]byte, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Err(err).Msg("Could not open file")
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Err(err).Msg("Could open file as []byte")
		return nil, err
	}

	var jsonFileContent json.RawMessage

	err = json.Unmarshal(byteValue, &jsonFileContent)
	if err != nil {
		log.Err(err).Msg("Could unmarshal file to json")
		return nil, err
	}

	jsonFileContentMsg, err := jsonFileContent.MarshalJSON()
	if err != nil {
		log.Err(err).Msg("Could marshal file to json")
		return nil, err
	}
	return jsonFileContentMsg, nil
}
