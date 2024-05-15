package telegram

import (
	"log"
	"strings"
)

const (
	StartCmd          = "/start"
	msgUnknownCommand = "What Is This ? :)"
)

func (p *Processor) Cmd(text, firstName, userName string, chatID int) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, userName)

	switch text {
	case StartCmd:
		return p.SendHello(chatID, firstName)
	case Error:
		return p.SendAlert(chatID, Error)
	case Warning:
		return p.SendAlert(chatID, Warning)
	case Okay:
		return p.SendAlert(chatID, Okay)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) SendHello(chatID int, name string) error {
	msgHello := "Привет, " + name + "!\n"
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) SendAlert(chatID int, text string) error {
	return p.tg.SendMessage(chatID, text)
}
