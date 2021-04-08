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
//
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var connectionString string
var queueName string

var azUsage = `A simple command line utility to send and receive AMQP message to / from Azure Service Bus:

Sending strings:
async-messaging-client sb send "my message" -q queueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
async-messaging-client sb receive -q queueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."

Using json files:
content test.json: { "key": "value" }
async-messaging-client sb send -f test.json -q queueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
async-messaging-client sb receive -q queueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."

You could also use environment variables for defining the queue and connection string:
export QUEUE=myQueueName
export CONNECTION_STRING='Endpoint=sb://host.servicebus.windows.net/'
async-messaging-client sb send -f test.json
async-messaging-client sb receive
`

var sbCmd = &cobra.Command{
	Use:   "sb",
	Short: "Interact with Azure Service Bus",
	Long:  azUsage,
}

func init() {
	rootCmd.AddCommand(sbCmd)

	flags := sbCmd.PersistentFlags()
	flags.StringP("connection-string", "c", "", "Connection String to connecto to Azure Service Bus")
	viper.BindPFlag("connection-string", flags.Lookup("connection-string"))

	flags.StringP("queue", "q", "", "Azure Service Bus Queue name")
	viper.BindPFlag("queue", flags.Lookup("queue"))

	cobra.OnInitialize(azConfig)
}

// azConfig reads in ENV variables if set.
func azConfig() {
	queueName = viper.GetString("queue")
	connectionString = viper.GetString("connection-string")
}
