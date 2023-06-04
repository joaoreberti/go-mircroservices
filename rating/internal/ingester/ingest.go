package ingester

import (
	"context"
	"encoding/json"
	"fmt"

	banana "github.com/confluentinc/confluent-kafka-go/kafka"
	"movieexample.com/rating/pkg/model"
)

// Ingester defines a Kafka ingester.
type Ingester struct {
	consumer *banana.Consumer
	topic    string
}

// NewIngester creates a new Kafka ingester.
func NewIngester(addr string, groupID string, topic string) (*Ingester, error) {
	consumer, err := banana.NewConsumer(&banana.ConfigMap{
		"bootstrap.servers":     addr,
		"group.id":              groupID,
		"auto.offset.reset":     "earliest",
		"broker.address.family": "v4", // Set this to v4 to disable IPv6
		"group.protocol.type":   "consumer",
	})
	if err != nil {
		return nil, err
	}
	return &Ingester{consumer, topic}, nil
}

// Ingest starts ingestion from Kafka and returns a channel // containing rating events
// representing the data consumed from the topic.
func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}

	ch := make(chan model.RatingEvent, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
			default:
			}
			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				fmt.Println("Consumer error: " + err.Error())
				continue
			}
			var event model.RatingEvent
			fmt.Printf("Received message: %s\n", msg.Value)
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Println("Unmarshal error: " + err.Error())
				continue
			}
			ch <- event
		}
	}()
	return ch, nil
}
