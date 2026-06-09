package commands

import "gopkg.in/telebot.v4"

func TestCommand(c telebot.Context, args []string) {
	c.Reply("Понг!")
}
