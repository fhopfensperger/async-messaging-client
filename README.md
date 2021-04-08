# async-messaging-client
![Go](https://github.com/fhopfensperger/async-messaging-client/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fhopfensperger/async-messaging-client)](https://goreportcard.com/report/github.com/fhopfensperger/async-messaging-client)
[![Release](https://img.shields.io/github/release/fhopfensperger/async-messaging-client.svg?style=flat-square)](https://github.com/fhopfensperger/async-messaging-client/releases/latest)


Sends and receives messages in an asynchronous way to / from different Cloud messaging services.

Currently, the following Cloud services are supported:
- [Google Cloud Pub/Sub](https://cloud.google.com/pubsub)
- [Azure Service Bus](https://azure.microsoft.com/en-us/services/service-bus/)

## Installation
### Linux / Mac OS
#### Option 1 (script)

```bash
curl https://raw.githubusercontent.com/fhopfensperger/async-messaging-client/main/get.sh | bash
```

#### Option 2 (manually)

Go to [Releases](https://github.com/fhopfensperger/async-messaging-client/releases) download the latest release according to your processor architecture and operating system, then unzip and copy it to the right location

```bash
tar xvfz async-messaging-client_x.x.x_darwin_amd64.tar.gz
cd async-messaging-client_x.x.x_darwin_amd64
chmod +x async-messaging-client
sudo mv async-messaging-client /usr/local/bin/
```
### Windows
1. Go to [Releases](https://github.com/fhopfensperger/async-messaging-client/releases)
2. Download the latest release async-messaging-client_x.x.x_windows_amd64.tar.gz
3. Use [7-Zip](https://www.7-zip.org/) or your favourite file archiver to unpack the archive
4. *Optional* Add the `async-messaging-client.exe` in to your `PATH`


## Usage Examples:

### Google Pub/Sub

---
**NOTE**
Before publishing message to a Pub/Sub topic and/ or subscribing to a subscription it is important to set proper authentication.
___

#### Print usage:
````bash
$ async-messaging-client pubsub --help
Interact with Google Cloud Pub/Sub

Usage:
  async-messaging-client pubsub [command]

Available Commands:
  publish     Publish a message to a Google Pub/Sub Topic
  subscribe   Subscribe to a Google Pub/Sub Subscription

Flags:
  -h, --help             help for pubsub
  -p, --project string   Google Cloud Project ID
  -t, --topic string     Google Cloud Pub/Sub Topic

Use "async-messaging-client pubsub [command] --help" for more information about a command.
````

#### Authentication
1. Create a service account and assign proper permissions in the Google Cloud Platform console
2. Download the service account key
3. Set environment variable GOOGLE_APPLICATION_CREDENTIALS
For Linux / Mac OS
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
```
For Windows
```cmd
set GOOGLE_APPLICATION_CREDENTIALS=C:\Path\to\service-account-key.json
```

#### Publishing (Sending)
**Local files**
```bash
# Content of test.json
{
  "prop1": 123,
  "prop2": "456",
  "prop3": "789",
  "prop4": {
    "prop1": 123,
    "prop2": "456",
    "propArray": [
      "456", "789", "0"
    ]
  }
}

# Content of attributes.json
{
  "eventType": "this.is.a.sample.event.type",
  "dispatcherId": "this.is.a.sample.id",
  "schemaVersion": "v1"
}
```
Send message (test.json) with attributes (attributes.json) to a topic `topic-test`
```bash
$ async-messaging-client pubsub publish -p google-project-id -t topic-test -f test.json -a attributes.json
2021-04-08T14:33:17+02:00 INF Trying to send message: 
{
  "prop1": 123,
  "prop2": "456",
  "prop3": "789",
  "prop4": {
    "prop1": 123,
    "prop2": "456",
    "propArray": [
      "456", "789", "0"
    ]
  }
}
with attributes: map[string]string{"dispatcherId":"this.is.a.sample.id", "eventType":"this.is.a.sample.event.type", "schemaVersion":"v1"}
to topic: topic-test for project: google-project-id
2021-04-08T14:33:18+02:00 INF Successfully published message with msgId: 2258869234167928 to projects/google-project-id/topics/topic-test 
```

**Simple Strings**
```bash
$ async-messaging-client pubsub publish -p google-project-id message-string
```
#### Subscribing (Receiving)

To receive message which are sent to a topic using a subscription and acknowledge them right away
```bash
$ async-messaging-client pubsub subscribe -p google-project-id -s test-sub
2021-04-08T14:36:08+02:00 INF Subscribing to test-sub on project google-project-id

2021-04-08T14:36:29+02:00 INF Got message: {
  "prop1": 123,
  "prop2": "456",
  "prop3": "789",
  "prop4": {
    "prop1": 123,
    "prop2": "456",
    "propArray": [
      "456", "789", "0"
    ]
  }
}, from projects/google-project-id/subscriptions/test-sub with attributes map[string]string{"dispatcherId":"this.is.a.sample.id", "eventType":"this.is.a.sample.event.type", "schemaVersion":"v1"} with message id 2258869234167928

 
```
