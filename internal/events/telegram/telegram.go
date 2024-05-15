package telegram

import (
	"log"

	"alertbot/internal/events"
	"alertbot/internal/models"
	"alertbot/internal/usecase"
)

type Processor struct {
	tg     *usecase.Client
	offset int
}

type MetaNew struct {
	ChatID    int
	FirstName string
	UserName  string
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, err
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.ProcessMessage(event)
	default:
		log.Fatal("can't process message")
	}

	return nil
}

func (p *Processor) ProcessMessage(event events.Event) error {
	meta, err := Meta(event)
	if err != nil {
		return err
	}

	if err := p.Cmd(
		event.Text,
		meta.FirstName,
		meta.UserName,
		meta.ChatID,
	); err != nil {
		return err
	}

	return nil
}

func Meta(event events.Event) (MetaNew, error) {
	res, ok := event.Meta.(MetaNew)
	if !ok {
		return MetaNew{}, nil
	}

	return res, nil
}

func event(upd models.Update) events.Event {
	res := events.Event{
		Type: FetchType(upd),
		Text: FetchText(upd),
	}

	if FetchType(upd) == events.Message {
		res.Meta = MetaNew{
			ChatID:    upd.Message.Chat.ID,
			FirstName: upd.Message.From.FirstName,
			UserName:  upd.Message.From.UserName,
		}
	}

	return res
}

func FetchText(upd models.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func FetchType(upd models.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func NewTg(client *usecase.Client) *Processor {
	return &Processor{
		tg: client,
	}
}
