package gcppubsub

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
)

type PubSubClient struct {
	client       *pubsub.Client
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}

func New(gcpProject string) *PubSubClient {

	c, err := pubsub.NewClient(context.Background(), gcpProject)
	if err != nil {
		log.Err(err).Msg("")
	}

	return &PubSubClient{client: c}
}

func (p *PubSubClient) SetTopic(topicName string) error {
	p.Topic = p.client.Topic(topicName)
	return nil
}

func (p *PubSubClient) SetSubscription(subscription string) error {
	p.Subscription = p.client.Subscription(subscription)
	return nil
}

func (p *PubSubClient) Publish(message []byte, attributes map[string]string) error {
	ctx := context.Background()
	msg := &pubsub.Message{
		Data:       message,
		Attributes: attributes,
	}
	if p.Topic == nil {
		return errors.New("")
	}

	msgId, err := p.Topic.Publish(ctx, msg).Get(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not publish message: %v", err))
	}

	log.Info().Msgf("Successfully published message with msgId: %s to %v", msgId, p.Topic)
	return nil
}

func (p *PubSubClient) Subscribe() error {
	sub := p.client.Subscription(p.Subscription.ID())
	err := sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		log.Info().Msgf("Got message: %s, from %s with attributes %#v with message id %s", m.Data, p.Subscription, m.Attributes, m.ID)
		m.Ack()
	})
	if err != nil {
		log.Info().Err(err).Msg("")
		return err
	}
	return nil
}
