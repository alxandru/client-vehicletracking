package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Callback func(value []byte, err error)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

// IConsumer .
type IConsumer interface {
	StartConsumer(cb Callback)
	StopConsumer()
}

// Consumer .
type Consumer struct {
	kreader *kafkago.Reader
	wg      sync.WaitGroup
}

func NewConsumer(endpoint string, topic string, groupid string) IConsumer {
	return &Consumer{
		kreader: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:        []string{endpoint},
			Topic:          topic,
			GroupID:        groupid,
			StartOffset:    kafkago.LastOffset,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second,
		}),
	}
}

func (c *Consumer) StartConsumer(cb Callback) {
	c.wg.Add(1)
	ctx, cancel = context.WithCancel(context.Background())
	go c.run(cb)
}

func (c *Consumer) StopConsumer() {
	cancel()
	c.wg.Wait()
	c.kreader.Close()
}

func (c *Consumer) run(cb Callback) {
	defer c.wg.Done()
	for {
		fmt.Println("Reading Message")
		m, err := c.kreader.ReadMessage(ctx)
		if err != nil {
			cb(nil, err)
			return
		}
		cb(m.Value, err)
	}
}
