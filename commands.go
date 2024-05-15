package alertbot

import (
	"log"
	"strings"
)

type Level int

const (
	LevelInfo Level = iota
	LevelWarning
	LevelErrorL
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
		return p.tg.sendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) send(chatID int, level Level, text string) error {
	return p.tg.sendMessage(chatID, text)
}

func (p *Processor) SendAlert(chatID int, text string) error {
	return p.send(chatID, LevelWarning, text)
}

func (p *Processor) SendHello(chatID int, name string) error {
	msgHello := "Привет, " + name + "!\n"
	return p.tg.sendMessage(chatID, msgHello)
}
