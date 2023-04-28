// Package natsdedup provides a deduplication solution for NATS messages.
package natsdedup

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Deduplicator is a simple message deduplication mechanism for NATS messages.
// It listens to messages on a specified input subject, deduplicates them
// using a configurable time-to-live (TTL) cache, and forwards unique messages
// to a specified output subject.
type Deduplicator struct {
	InputSubject     string
	OutputSubject    string
	DeduplicationTTL time.Duration
	recentMessages   map[string]time.Time
	mux              sync.Mutex
}

// NewDeduplicator creates a new deduplicator with the specified input subject,
// output subject, and deduplication TTL.
func NewDeduplicator(inputSubject, outputSubject string, deduplicationTTL time.Duration) *Deduplicator {
	return &Deduplicator{
		InputSubject:     inputSubject,
		OutputSubject:    outputSubject,
		DeduplicationTTL: deduplicationTTL,
		recentMessages:   make(map[string]time.Time),
	}
}

// handleMessage processes a received NATS message, deduplicates it, and
// forwards unique messages to the output subject.
func (d *Deduplicator) handleMessage(msg *nats.Msg, nc *nats.Conn) {
	d.mux.Lock()
	defer d.mux.Unlock()

	msgStr := string(msg.Data)
	if _, seen := d.recentMessages[msgStr]; seen {
		// Duplicate message, ignore it
		return
	}

	// Remember this message and remove it after the deduplication TTL
	d.recentMessages[msgStr] = time.Now()
	time.AfterFunc(d.DeduplicationTTL, func() {
		d.mux.Lock()
		defer d.mux.Unlock()
		delete(d.recentMessages, msgStr)
	})

	// Forward the message to the output subject
	// TODO: Handle errors here (e.g. if the output subject is invalid)
	_ = nc.Publish(d.OutputSubject, msg.Data)
}

// Run starts the deduplicator, subscribing to the input subject and forwarding
// unique messages to the output subject. Pass a connected nats.Conn object to
// interact with the NATS server.
func (d *Deduplicator) Run(nc *nats.Conn) error {
	_, err := nc.Subscribe(d.InputSubject, func(msg *nats.Msg) {
		d.handleMessage(msg, nc)
	})
	return err
}
