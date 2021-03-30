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
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

// peekCmd represents the peek command
var peekCmd = &cobra.Command{
	Use:   "peek",
	Short: "Get message preview, but dont take it from queue",
	Long:  `Get message preview, but dont take it from queue`,
	Run: func(cmd *cobra.Command, args []string) {
		peek(queueName)
	},
}

func init() {
	rootCmd.AddCommand(peekCmd)

}

func peek(queueName string) {
	log.Info().Msgf("Peeking messages for %s", queueName)

	// Create a context to limit how long we will try to send, then push the message over the wire.
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()

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

	messageIterator, err := client.Peek(ctx)
	if err != nil {
		log.Err(err).Msgf("Could not peek message for queue %s", queueName)
	}
	for !messageIterator.Done() {
		msg, err := messageIterator.Next(ctx)
		if err != nil {
			break
		}
		if msg.SystemProperties.ScheduledEnqueueTime != nil {
			log.Info().Msgf("Scheduled message (No %v):\n%v\non %s delivery time: %s", *msg.SystemProperties.SequenceNumber, string(msg.Data), queueName, msg.SystemProperties.ScheduledEnqueueTime)
		} else {
			log.Info().Msgf("Message (No %v):\n%v\non %s", *msg.SystemProperties.SequenceNumber, string(msg.Data), queueName)
		}
	}
}
