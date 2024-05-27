package alertbot

import (
	"strings"
)

const (
	msgUnknownCommand = "What Is This ? :)"
)

func (p *Processor) Cmd(text string) error {
	text = strings.TrimSpace(text)

	switch text {
	case Error:
		return p.SendAlert(Error)
	case Warning:
		return p.SendAlert(Warning)
	case Okay:
		return p.SendAlert(Okay)
	default:
		return p.tg.sendMessage(msgUnknownCommand)
	}
}

func (p *Processor) SendAlert(text string) error {
	return p.tg.sendMessage(text)
}
