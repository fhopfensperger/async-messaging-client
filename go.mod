module github.com/fhopfensperger/async-messaging-client

go 1.16

require (
	cloud.google.com/go/pubsub v1.29.0
	github.com/Azure/azure-service-bus-go v0.11.5
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.29.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.15.0
)

replace github.com/gin-gonic/gin v1.7.3 => github.com/gin-gonic/gin v1.7.4
