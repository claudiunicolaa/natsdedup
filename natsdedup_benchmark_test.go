package natsdedup_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/claudiunicolaa/natsdedup"
	"github.com/nats-io/nats.go"
)

func BenchmarkNatsDedup(b *testing.B) {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS server: %v", err)
	}
	defer nc.Close()

	// Create a deduplicator
	inputSubject := "test.input"
	outputSubject := "test.output"
	deduplicationTTL := 1 * time.Second
	deduplicator := natsdedup.NewDeduplicator(inputSubject, outputSubject, deduplicationTTL)

	// Run the deduplicator
	go func() {
		if err := deduplicator.Run(nc); err != nil {
			log.Fatalf("Error running deduplicator: %v", err)
		}
	}()

	// Benchmark loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := fmt.Sprintf("Message-%d", i)
		if err := nc.Publish(inputSubject, []byte(msg)); err != nil {
			log.Fatalf("Error publishing message: %v", err)
		}
	}

	// Allow time for the deduplicator to process the remaining messages
	time.Sleep(2 * time.Second)
}
