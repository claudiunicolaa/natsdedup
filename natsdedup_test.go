package natsdedup_test

import (
	"testing"
	"time"

	"github.com/claudiunicolaa/natsdedup"
	"github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeduplicator(t *testing.T) {
	// Start a NATS server for testing
	natsServer := test.RunDefaultServer()
	defer natsServer.Shutdown()

	// Connect to the test NATS server
	nc, err := nats.Connect(natsServer.ClientURL())
	require.NoError(t, err)
	defer nc.Close()

	// Set up test subjects and deduplicator
	inputSubject := "test.source"
	outputSubject := "test.destination"
	deduplicationTTL := 100 * time.Millisecond

	deduplicator := natsdedup.NewDeduplicator(inputSubject, outputSubject, deduplicationTTL)
	require.NoError(t, deduplicator.Run(nc))

	// Test deduplication
	duplicateMessage := []byte("duplicate message")
	numMessages := 5
	numDeduplicated := 0

	outputSub, err := nc.SubscribeSync(outputSubject)
	require.NoError(t, err)

	// Send duplicate messages
	for i := 0; i < numMessages; i++ {
		require.NoError(t, nc.Publish(inputSubject, duplicateMessage))
	}

	// Try to receive messages on the output subject
	for {
		msg, err := outputSub.NextMsg(200 * time.Millisecond)
		if err != nil {
			if err == nats.ErrTimeout {
				break
			}
			require.NoError(t, err)
		}
		assert.Equal(t, duplicateMessage, msg.Data)
		numDeduplicated++
	}

	// Expect only one message to be forwarded due to deduplication
	assert.Equal(t, 1, numDeduplicated)
}
