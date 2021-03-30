/*
Copyright © 2021 Florian Hopfensperger <f.hopfensperger@gmail.com>

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

	"github.com/rs/zerolog/log"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/spf13/cobra"
)

// cancelCmd represents the cancelScheduled command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel all scheduled AMQP messages",
	Long:  `Cancel all scheduled AMQP messages`,
	Run: func(cmd *cobra.Command, args []string) {
		cancelScheduled()
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

}

func cancelScheduled() {
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

	sequenceNumbers := make([]int64, 0)

	messageIterator, err := client.Peek(ctx)
	if err != nil {
		log.Err(err).Msgf("Could not peek message for queue %s", queueName)
	}
	for !messageIterator.Done() {
		msg, err := messageIterator.Next(ctx)
		if err != nil {
			break
		}
		log.Info().Msgf("Cancel scheduled message: \n%s\non queue %s (time: %v) with sequence number: %v and id: %s", msg.Data, queueName, msg.SystemProperties.ScheduledEnqueueTime, *msg.SystemProperties.SequenceNumber, msg.ID)
		if msg.SystemProperties.ScheduledEnqueueTime != nil {
			sequenceNumbers = append(sequenceNumbers, *msg.SystemProperties.SequenceNumber)
		}
	}
	if len(sequenceNumbers) < 1 {
		log.Info().Msgf("No scheduled message for queue %s", queueName)
		return
	}
	err = client.CancelScheduled(ctx, sequenceNumbers...)
	if err != nil {
		log.Err(err).Msgf("Could not cancel message for queue %s", queueName)
		return
	}
	log.Info().Msgf("Successfully canceled message(s) with sequence number(s) %v", sequenceNumbers)

}
