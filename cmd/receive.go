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
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"

	servicebus "github.com/Azure/azure-service-bus-go"

	"github.com/spf13/cobra"
)

var duration time.Duration

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive AMQP messages from Azure Service Bus",
	Long: `Receive AMQP messages from a Azure Service Bus queue. 
Finishes either after receiving one message or after a specific duration.
It can also listen on several queues simultaneously.`,
	Args: NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		duration = viper.GetDuration("duration")
		multipleQueues := viper.GetStringSlice("multiple-queues")
		parallel := viper.GetInt("parallel")
		var wg sync.WaitGroup
		if len(multipleQueues) > 0 {

			for _, q := range multipleQueues {
				for numberOfWorker := 0; numberOfWorker < parallel; numberOfWorker++ {
					wg.Add(1)
					if duration.Milliseconds() > 0 {
						go receiveWitDuration(q, duration, &wg)
					} else {
						go receiveOne(q, &wg)
					}
				}
			}
			wg.Wait()

		} else if duration.Milliseconds() > 0 {
			for numberOfWorker := 0; numberOfWorker < parallel; numberOfWorker++ {
				wg.Add(1)
				go receiveWitDuration(queueName, duration, &wg)
			}
			wg.Wait()
			return
		} else {
			receiveOne(queueName, nil)
		}
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	flags := receiveCmd.Flags()
	flags.DurationP("duration", "d", duration, "Listen on queue for duration, example: 10m, 1h, 1h10m, 1h10m10s")
	viper.BindPFlag("duration", flags.Lookup("duration"))

	flags.StringSliceP("multiple-queues", "m", []string{}, "Listen on multiple queues, example: queue1,queue2")
	viper.BindPFlag("multiple-queues", flags.Lookup("multiple-queues"))

	flags.IntP("parallel", "p", 1, `Run x multiple listener in parallel. Must be combined with "duration"`)
	viper.BindPFlag("parallel", flags.Lookup("parallel"))
}

func receiveOne(queueName string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	log.Info().Msgf("Receiving one message from: %s", queueName)

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		log.Err(err).Msg("Receiving failed! Check connection string.")
		return
	}

	client, err := ns.NewQueue(queueName)
	if err != nil {
		log.Err(err).Msgf("Could not open queue %s", queueName)
		return
	}

	// Define a context to limit how long we will block to receiveOne messages, then start serving our function.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	if err := client.ReceiveOne(ctx, printMessage); err != nil {
		log.Err(err).Msgf("Could not ReceiveOne from queue %s", queueName)
	}
}

func receiveWitDuration(queueName string, duration time.Duration, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	log.Info().Msgf("Receiving messages from: %s for: %s", queueName, duration)

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		log.Err(err).Msg("Receiving failed! Check connection string.")
		return
	}

	client, err := ns.NewQueue(queueName)
	if err != nil {
		log.Err(err).Msgf("Could not open queue %s", queueName)
		return
	}

	// Define a context to limit how long we will block to receive messages, then start serving our function.
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	if err := client.Receive(ctx, printMessage); err != nil {
		log.Err(err).Msgf("Could not Receive from queue %s", queueName)
	}
}

// Define a function that should be executed when a message is received.
var printMessage servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) error {

	log.Info().Msgf("Message:\n%v\nreceived from %s", string(msg.Data), queueName)

	return msg.Complete(ctx)
}
