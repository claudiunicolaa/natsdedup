# natsdedup

[![Go Report Card](https://goreportcard.com/badge/github.com/claudiunicolaa/natsdedup)](https://goreportcard.com/report/github.com/claudiunicolaa/natsdedup)
[![Run natsdedup tests](https://github.com/claudiunicolaa/natsdedup/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/claudiunicolaa/natsdedup/actions/workflows/tests.yml)

[![Open in DevPod!](https://devpod.sh/assets/open-in-devpod.svg)](https://devpod.sh/open#https://github.com/claudiunicolaa/natsdedup)

`natsdedup` is a lightweight package for deduplicating messages on NATS subjects. 
It listens to messages on a specified input subject, deduplicates them using a configurable time-to-live (TTL) cache, and forwards unique messages to a specified output subject.

Although the `natsdedup` package provides a deduplication solution for NATS messages, it's important to note that the NATS ecosystem offers built-in deduplication capabilities through its [JetStream technology](https://docs.nats.io/using-nats/developer/develop_jetstream/model_deep_dive#message-deduplication). 
However, this feature is not available in the standalone NATS server, which is where `natsdedup` comes in handy. 
By using the `natsdedup` package, users who are not utilizing JetStream can still benefit from message deduplication, thereby enhancing the efficiency and reliability of their NATS-based applications.

It is important to mention the `natsdedup` provides a different approach to deduplication compared to the built-in deduplication capabilities in NATS JetStream. 
While JetStream's deduplication is based on the message ID header, natsdedup focuses on the message content.

## Features

- Listens to messages on a NATS input subject
- Deduplicates messages using a TTL cache
- Forwards unique messages to a NATS output subject

### The cache 

Currently employs an in-memory caching mechanism to temporarily store messages and efficiently deduplicate them. 
This approach ensures low-latency processing and minimal overhead. 
However, as the project evolves, the plan is to explore and implement more advanced caching strategies to cater to various use cases and requirements. 

Potential additions could include disk-based caching, distributed caching, and support for popular caching systems such as Redis or Memcached. 
Contributions are welcomed - not only for expanding the caching options but also for enhancing the overall functionality and performance of the package.

## Installation

```bash
go get github.com/claudiunicolaa/natsdedup
```

## Usage
To use `natsdedup` as a library in your Go project, import the package and create a new instance of the `Deduplicator`:

```go
import (
	"github.com/nats-io/nats.go"
	"github.com/claudiunicolaa/natsdedup"
)

// Connect to NATS
nc, err := nats.Connect("nats://localhost:4222")
if err != nil {
	log.Fatalf("Failed to connect to NATS: %v", err)
}
defer nc.Close()

// Create and run the deduplicator
inputSubject := "source.subject"
outputSubject := "destination.subject"
deduplicationTTL := 1 * time.Minute
deduplicator := natsdedup.NewDeduplicator(inputSubject, outputSubject, deduplicationTTL)
if err := deduplicator.Run(nc); err != nil {
	log.Fatalf("Failed to run deduplicator: %v", err)
}

// Wait for messages indefinitely
select {}

```

## Examples
Example project [natsdedup-example](https://github.com/claudiunicolaa/natsdedup-example/tree/main).

## Contributions
If you are interested in contributing, please feel free to submit pull requests, report issues, or propose new ideas for the project.

## Publications

[Effortless NATS Message Deduplication in Go](https://dev.to/claudiunicolaa/effortless-nats-message-deduplication-in-go-2ohl)
