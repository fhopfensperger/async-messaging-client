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
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pubSubTopic string
var gcpProject string
var gcpApiKey string

var pubSubUsage = `Interact with Google Cloud Pub/Sub`

var pubSubCommand = &cobra.Command{
	Use:   "pubsub",
	Short: "Interact with Google Cloud Pub/Sub",
	Long:  pubSubUsage,
}

func init() {
	rootCmd.AddCommand(pubSubCommand)

	flags := pubSubCommand.PersistentFlags()
	flags.StringP("project", "p", "", "Google Cloud Project ID")
	viper.BindPFlag("project", flags.Lookup("project"))

	flags.StringP("topic", "t", "", "Google Cloud Pub/Sub Topic")
	viper.BindPFlag("topic", flags.Lookup("topic"))

	cobra.OnInitialize(gcpConfig)
}

// gcpConfig reads in ENV variables if set.
func gcpConfig() {
	pubSubTopic = viper.GetString("topic")
	gcpProject = viper.GetString("project")
}

func transformAttributes(fileName string) (map[string]string, error) {
	attributes := make(map[string]string)
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Err(err).Msg("Could not open file")
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Err(err).Msg("Could open file as []byte")
		return nil, err
	}
	err = json.Unmarshal(byteValue, &attributes)
	if err != nil {
		log.Err(err).Msg("Could open file as []byte")
		return nil, err
	}
	return attributes, nil
}
