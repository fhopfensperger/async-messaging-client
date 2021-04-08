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
	"github.com/spf13/viper"

	"github.com/fhopfensperger/async-messaging-client/pkg/gcppubsub"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// subscribeCmd represents the send command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to a Google Pub/Sub Subscription",
	Long:  `Subscribe to a Google Pub/Sub Subscription`,
	Run: func(cmd *cobra.Command, args []string) {
		subscription := viper.GetString("subscription")
		log.Info().Msgf("Subscribing to %s on project %s", subscription, gcpProject)
		err := subscribe(subscription)
		if err != nil {
			log.Err(err).Msg("")
			return
		}
	},
}

func init() {
	pubSubCommand.AddCommand(subscribeCmd)
	flags := subscribeCmd.Flags()
	flags.StringP("subscription", "s", "", "Pub/Sub subscription")
	viper.BindPFlag("subscription", flags.Lookup("subscription"))
}

func subscribe(subscription string) error {
	if subscription == "" {
		err := errors.New("subscription could not empty")
		log.Err(err).Msg("")
		return err
	}
	client := gcppubsub.New(gcpProject)

	err := client.SetSubscription(subscription)
	if err != nil {
		log.Err(err).Msgf("Could not set subscription %s", subscription)
		return err
	}

	err = client.Subscribe()
	if err != nil {
		log.Err(err).Msgf("Could not subscribe to %s", subscription)
		return err
	}
	return nil

}
