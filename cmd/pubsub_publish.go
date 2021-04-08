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
	"github.com/pkg/errors"

	"github.com/fhopfensperger/async-messaging-client/pkg/gcppubsub"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// publishCmd represents the send command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a message to a Google Pub/Sub Topic",
	Long:  `Publish a message to a Google Pub/Sub Topic`,
	Run: func(cmd *cobra.Command, args []string) {
		var message []byte
		var err error
		file := viper.GetString("pubsubfile")
		if file != "" {
			message, err = transformFileToJson(file)
			if err != nil {
				log.Err(err).Msg("")
				return
			}
		} else if len(args) == 1 {
			message = []byte(args[0])
		} else {
			log.Warn().Msgf("Nothing send because the message is empty")
			return
		}
		attributesFile := viper.GetString("attributes")
		attributes, err := transformAttributes(attributesFile)
		if err != nil {
			log.Err(err).Msg("")
			return
		}

		log.Info().Msgf("Trying to send message: \n%s\nwith attributes: %#v\nto topic: %v for project: %v", message, attributes, pubSubTopic, gcpProject)
		err = publish(message, attributes, pubSubTopic, gcpProject)
	},
}

func init() {
	pubSubCommand.AddCommand(publishCmd)
	flags := publishCmd.Flags()
	flags.StringP("file", "f", "", "Sends .json file (must be .json) as message")
	viper.BindPFlag("pubsubfile", flags.Lookup("file"))

	flags.StringP("attributes-file", "a", "", "Publish optional Pub/Sub message attributes from file")
	viper.BindPFlag("attributes", flags.Lookup("attributes-file"))
}

func publish(message []byte, attributes map[string]string, topic, gcpProject string) error {
	if topic == "" || gcpProject == "" {
		err := errors.New("topic or project could not empty")
		log.Err(err).Msg("")
		return err
	}
	client := gcppubsub.New(gcpProject)
	err := client.SetTopic(topic)
	if err != nil {
		log.Err(err).Msgf("Could not set pubSubTopic %v", topic)
		return err
	}

	err = client.Publish(message, attributes)
	if err != nil {
		log.Err(err).Msgf("Could not publish message %v to pubSubTopic %v", string(message), topic)
		return err
	}
	return nil
}
