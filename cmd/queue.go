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

// queueCmd represents the queue command
var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Get count details for queue",
	Long:  `Get count details for queue`,
	Run: func(cmd *cobra.Command, args []string) {
		checkQueueSize(queueName)
	},
}

func init() {
	rootCmd.AddCommand(queueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queueCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkQueueSize(queueName string) {
	log.Info().Msgf("Querying count details for %s", queueName)

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		log.Err(err).Msg("Receiving failed! Check connection string.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	qm := ns.NewQueueManager()
	q, err := qm.Get(ctx, queueName)
	if err != nil {
		log.Err(err).Msgf("Queue %s not found", queueName)
		return
	}

	countDetails := *q.CountDetails

	log.Info().Msgf("Count details for queue %s\nActiveMessageCount: %v\nDeadLetterMessageCount: %v\nScheduledMessageCount: %v", queueName,
		*countDetails.ActiveMessageCount, *countDetails.DeadLetterMessageCount,
		*countDetails.ScheduledMessageCount)
}
