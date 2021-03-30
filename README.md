# async-messaging-client
![Go](https://github.com/fhopfensperger/async-messaging-client/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fhopfensperger/async-messaging-client)](https://goreportcard.com/report/github.com/fhopfensperger/async-messaging-client)
[![Release](https://img.shields.io/github/release/fhopfensperger/async-messaging-client.svg?style=flat-square)](https://github.com//fhopfensperger/async-messaging-client/releases/latest)


Sends and receives AMQP messages to / from Azure Service Bus

## Installation

### Option 1 (script)

```bash
curl https://raw.githubusercontent.com/fhopfensperger/async-messaging-client/main/get.sh | bash
```

### Option 2 (manually)

Go to [Releases](https://github.com/fhopfensperger/async-messaging-client/releases) download the latest release according to your processor architecture and operating system, then unzip and copy it to the right location

```bash
tar xvfz async-messaging-client_x.x.x_darwin_amd64.tar.gz
cd async-messaging-client_x.x.x_darwin_amd64
chmod +x async-messaging-client
sudo mv async-messaging-client /usr/local/bin/
```

## Usage Examples:
### Option 1
##### **`test.json`**
```json 
{ "key1": "value1", "key2": "value2", "message" }
```
##### **Sending**
```bash
async-messaging-client send -f test.json -q myQueueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
```
##### **Receiving one message**
```bash
async-messaging-client receive -q myQueueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
```

##### **Receiving for a specific duration**
```bash
async-messaging-client receive -d 10m -q myQueueName -c "Endpoint=sb://host.servicebus.windows.net/;SharedAccessKeyName=..."
```

### Option 2 (using environment variables)
##### **Setting environment variables**
```bash
export CONNECTION_STRING='Endpoint=sb:...'
export QUEUE="myQueueName"
```
##### **Sending**
```bash
async-messaging-client send -f test.json 
```
##### **Receiving**
```bash
async-messaging-client receive -d 1h
```