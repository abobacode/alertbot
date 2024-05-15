package alertbot

import (
	"log"
)

type Processor struct {
	tg     *client
	offset int
}

type MetaNew struct {
	ChatID    int
	FirstName string
	UserName  string
}

func (p *Processor) Fetch(limit int) ([]Event, error) {
	updates, err := p.tg.updates(p.offset, limit)
	if err != nil {
		return nil, err
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event Event) error {
	switch event.Type {
	case Message:
		return p.ProcessMessage(event)
	default:
		log.Fatal("can't process message")
	}

	return nil
}

func (p *Processor) ProcessMessage(event Event) error {
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

func Meta(event Event) (MetaNew, error) {
	res, ok := event.Meta.(MetaNew)
	if !ok {
		return MetaNew{}, nil
	}

	return res, nil
}

func event(upd Update) Event {
	res := Event{
		Type: FetchType(upd),
		Text: FetchText(upd),
	}

	if FetchType(upd) == Message {
		res.Meta = MetaNew{
			ChatID:    upd.Message.Chat.ID,
			FirstName: upd.Message.From.FirstName,
			UserName:  upd.Message.From.UserName,
		}
	}

	return res
}

func FetchText(upd Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func FetchType(upd Update) Type {
	if upd.Message == nil {
		return Unknown
	}

	return Message
}

func NewTg(client *client) *Processor {
	return &Processor{
		tg: client,
	}
}
