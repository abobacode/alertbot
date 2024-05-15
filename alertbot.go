package alertbot

import (
	"context"
)

func New(ctx context.Context, host, token string) error {
	const (
		batchSize = 100
	)

	eventProcessor := NewTg(newClient(host, token))
	cons := newConsumer(eventProcessor, eventProcessor, batchSize)

	if err := cons.start(ctx); err != nil {
		return err
	}
	return nil
}
