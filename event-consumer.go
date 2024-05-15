package alertbot

import (
	"context"
	"log"
	"time"
)

const (
	Unknown Type = iota
	Message
)

type Type int

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type processor interface {
	Process(e Event) error
}

type consumer struct {
	fetcher   fetcher
	processor processor
	batchSize int
}

func (c *consumer) start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			gotEvents, err := c.fetcher.Fetch(c.batchSize)
			if err != nil {
				log.Printf("[ERR] consumer: %s", err.Error())
				continue
			}

			if len(gotEvents) > 0 {
				if err = c.handleEvents(ctx, gotEvents); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func (c *consumer) handleEvents(ctx context.Context, events []Event) error {
	for _, event := range events {
		select {
		case <-ctx.Done():
			return nil
		default:

		}

		if event.Type == 0 {
			continue
		}

		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			return err
		}
	}

	return nil
}

func newConsumer(fetcher fetcher, processor processor, batchSize int) *consumer {
	return &consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}
