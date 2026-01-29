package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  sarama.ConsumerGroupHandler
}

// Event structure for user events
type UserEvent struct {
	Type      string                 `json:"type"` // created, updated, deleted
	UserID    int64                  `json:"user_id"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaProducer{producer: producer}, nil
}

func (p *KafkaProducer) SendMessage(topic string, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}

func (p *KafkaProducer) SendUserEvent(event *UserEvent) error {
	return p.SendMessage("user-events", fmt.Sprintf("user-%d", event.UserID), event)
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}

// Consumer Handler
type ConsumerHandler struct {
	ready   chan bool
	handler func(message *sarama.ConsumerMessage) error
}

func NewConsumerHandler(handler func(*sarama.ConsumerMessage) error) *ConsumerHandler {
	return &ConsumerHandler{
		ready:   make(chan bool),
		handler: handler,
	}
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.handler(message); err != nil {
			log.Printf("Error handling message: %v", err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string, handler func(*sarama.ConsumerMessage) error) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.V2_6_0_0

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		topics:   topics,
		handler:  NewConsumerHandler(handler),
	}, nil
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	go func() {
		for {
			if err := c.consumer.Consume(ctx, c.topics, c.handler); err != nil {
				log.Printf("Error consuming messages: %v", err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-c.handler.(*ConsumerHandler).ready
	log.Printf("Kafka consumer started for topics: %s", strings.Join(c.topics, ", "))
	return nil
}

func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}

// Example handler function
func DefaultUserEventHandler(message *sarama.ConsumerMessage) error {
	var event UserEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user event: %w", err)
	}

	log.Printf("Received user event: Type=%s, UserID=%d, Timestamp=%d",
		event.Type, event.UserID, event.Timestamp)

	// Process the event based on type
	switch event.Type {
	case "created":
		log.Printf("User created: %v", event.Data)
	case "updated":
		log.Printf("User updated: %v", event.Data)
	case "deleted":
		log.Printf("User deleted: UserID=%d", event.UserID)
	}

	return nil
}
