package natsdedup

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Deduplicator struct {
	InputSubject     string
	OutputSubject    string
	DeduplicationTTL time.Duration
	recentMessages   map[string]time.Time
	mux              sync.Mutex
}

func NewDeduplicator(inputSubject, outputSubject string, deduplicationTTL time.Duration) *Deduplicator {
	return &Deduplicator{
		InputSubject:     inputSubject,
		OutputSubject:    outputSubject,
		DeduplicationTTL: deduplicationTTL,
		recentMessages:   make(map[string]time.Time),
	}
}

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
	nc.Publish(d.OutputSubject, msg.Data)
}

func (d *Deduplicator) Run(nc *nats.Conn) error {
	_, err := nc.Subscribe(d.InputSubject, func(msg *nats.Msg) {
		d.handleMessage(msg, nc)
	})
	return err
}
