package main

import (
	"log"

	"DebilBot/globals"

	"github.com/pelletier/go-toml"
)

var (
	BotApiUrl string

	Appeals []interface{}

	StupidMessagesEnable   bool
	StupidMessagesTriggers []interface{}
	StupidMessagesText     string
	StupidMessagesChatIDs  []interface{}
)

func LoadConfig() {

	config, err := toml.LoadFile("conf.toml")
	if err != nil {
		log.Panic(err)
	}

	globals.AccessToken = config.Get("account.access_token").(string)
	BotApiUrl = config.Get("account.url").(string)

	globals.BotSettings = config.Get("bot_settings").(*toml.Tree)
	Appeals = config.Get("bot_settings.appeal").([]interface{})

	StupidMessagesEnable = config.Get("stupid_messages.enable").(bool)
	if StupidMessagesEnable {
		StupidMessagesTriggers = config.Get("stupid_messages.triggers").([]interface{})
		StupidMessagesText = config.Get("stupid_messages.text").(string)
		StupidMessagesChatIDs = config.Get("stupid_messages.chat_ids").([]interface{})
	}
}
