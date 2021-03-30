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
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"

	servicebus "github.com/Azure/azure-service-bus-go"

	"github.com/spf13/cobra"
)

var delay time.Duration

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send AMQP message to Azure Service Bus",
	Long:  `Send AMQP message to Azure Service Bus either from a string or from a JSON file`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		delay = viper.GetDuration("schedule")
		file := viper.GetString("file")
		if file != "" {
			sendJSONFile(file)
			return
		}
		send([]byte(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	flags := sendCmd.Flags()
	flags.StringP("file", "f", "", "Sends .json file to Queue (must be .json)")
	viper.BindPFlag("file", flags.Lookup("file"))

	flags.DurationP("schedule", "s", delay, "Sends scheduled message; Delay specified as duration: 10m, 1h, 1h10m, 1h10m10s")
	viper.BindPFlag("schedule", flags.Lookup("schedule"))

}

func send(messageContent []byte) {
	log.Info().Msgf("Sending message: \n%s\nto %s", messageContent, queueName)

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		return
	}

	client, err := ns.NewQueue(queueName)
	if err != nil {
		log.Err(err).Msgf("Could not use queue %s", queueName)
		return
	}

	// Create a context to limit how long we will try to send, then push the message over the wire.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := servicebus.NewMessage(messageContent)
	message.ContentType = "application/json"

	if delay.Milliseconds() > 0 {
		expectedTime := time.Now().Add(delay)
		sequenceNo, err := client.ScheduleAt(ctx, expectedTime, message)
		if err != nil {
			log.Err(err).Msg("Could not send scheduled msg")
			return
		}
		log.Info().Msgf("Scheduled message with id: %s and sequence number %v will be delivered at %s to %s", message.ID, sequenceNo, expectedTime.UTC(), queueName)
	} else {
		if err := client.Send(ctx, message); err != nil {
			log.Err(err).Msg("Could not send msg")
			return
		}
		log.Info().Msgf("Sent message with id %s to %s", message.ID, queueName)
	}

}

func sendJSONFile(fileName string) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Err(err).Msg("Could not open file")
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Err(err).Msg("Could open file as []byte")
		return
	}

	var jsonFileContent json.RawMessage

	err = json.Unmarshal(byteValue, &jsonFileContent)
	if err != nil {
		log.Err(err).Msg("Could unmarshal file to json")
		return
	}

	jsonFileContentMsg, err := jsonFileContent.MarshalJSON()
	if err != nil {
		log.Err(err).Msg("Could marshal file to json")
		return
	}

	send(jsonFileContentMsg)
}
