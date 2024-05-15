package telegram

import (
	"log"
	"strings"
)

const (
	StartCmd = "/start"
)

func (p *Processor) cmd(text, firstName, userName string, chatID int) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, userName)

	switch text {
	case StartCmd:
		buttons := [][]string{
			{},
			{},
			{},
		}
		return p.sendHello(chatID, firstName, buttons)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) sendHello(chatID int, name string, buttons [][]string) error {
	msgHello := "Привет, " + name + "!\n"
	path := "C:\\Goland\\ovpn\\internal\\telegram\\pics\\vpn_start.png"
	return p.tg.SendPhoto(chatID, msgHello, path, buttons)
}
