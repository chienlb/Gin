package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"gin-demo/internal/config"
	"gin-demo/pkg/cache"
	"gin-demo/pkg/messaging"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Redis
	redisClient, err := cache.NewRedisClient(
		cfg.GetRedisAddr(),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize Kafka Producer
	producer, err := messaging.NewKafkaProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize Kafka Consumer
	consumer, err := messaging.NewKafkaConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.ConsumerGroup,
		[]string{cfg.Kafka.Topics.UserEvents},
		messaging.DefaultUserEventHandler,
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	ctx := context.Background()

	// Start Kafka consumer
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	// Example 1: Redis Cache Operations
	fmt.Println("\n=== Redis Cache Examples ===")

	// Set a value
	err = redisClient.Set(ctx, "user:1000", map[string]interface{}{
		"id":    1000,
		"name":  "John Doe",
		"email": "john@example.com",
	}, 1*time.Hour)
	if err != nil {
		log.Printf("Redis SET error: %v", err)
	} else {
		fmt.Println("✅ User cached successfully")
	}

	// Get a value
	var cachedUser map[string]interface{}
	err = redisClient.Get(ctx, "user:1000", &cachedUser)
	if err != nil {
		log.Printf("Redis GET error: %v", err)
	} else {
		fmt.Printf("✅ Retrieved user from cache: %v\n", cachedUser)
	}

	// Check if key exists
	exists, _ := redisClient.Exists(ctx, "user:1000")
	fmt.Printf("✅ Key exists: %v\n", exists)

	// Set with NX (only if not exists)
	set, _ := redisClient.SetNX(ctx, "user:1000", "some value", 1*time.Hour)
	fmt.Printf("✅ SetNX result (should be false): %v\n", set)

	// Example 2: Kafka Producer
	fmt.Println("\n=== Kafka Producer Examples ===")

	// Send user created event
	event := &messaging.UserEvent{
		Type:      "created",
		UserID:    1001,
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"name":  "Jane Smith",
			"email": "jane@example.com",
		},
	}
	err = producer.SendUserEvent(event)
	if err != nil {
		log.Printf("Kafka send error: %v", err)
	} else {
		fmt.Println("✅ User created event sent to Kafka")
	}

	// Send user updated event
	updateEvent := &messaging.UserEvent{
		Type:      "updated",
		UserID:    1001,
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"name":  "Jane Smith Updated",
			"email": "jane.updated@example.com",
		},
	}
	err = producer.SendUserEvent(updateEvent)
	if err != nil {
		log.Printf("Kafka send error: %v", err)
	} else {
		fmt.Println("✅ User updated event sent to Kafka")
	}

	// Example 3: Rate Limiting with Redis
	fmt.Println("\n=== Rate Limiting Example ===")

	rateLimitKey := "rate_limit:user:1001"
	for i := 1; i <= 5; i++ {
		count, _ := redisClient.Increment(ctx, rateLimitKey)
		if i == 1 {
			// Set expiration on first increment
			redisClient.Expire(ctx, rateLimitKey, 60*time.Second)
		}
		fmt.Printf("Request %d: Count = %d\n", i, count)

		if count > 3 {
			fmt.Println("❌ Rate limit exceeded!")
			break
		}
	}

	// Example 4: Cache Invalidation
	fmt.Println("\n=== Cache Invalidation Example ===")

	err = redisClient.Delete(ctx, "user:1000")
	if err != nil {
		log.Printf("Delete error: %v", err)
	} else {
		fmt.Println("✅ Cache invalidated")
	}

	// Wait a bit for consumer to process messages
	fmt.Println("\n⏳ Waiting for Kafka consumer to process messages...")
	time.Sleep(3 * time.Second)

	fmt.Println("\n✅ All examples completed!")
}
