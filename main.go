package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/ztrue/shutdown"
	"gopkg.in/telebot.v4" // Перешли на v4

	"DebilBot/chatbot"
	"DebilBot/globals"
)

var (
	bot *telebot.Bot
)

func main() {
	log.Println("Загрузка конфигов...")
	LoadConfig()

	log.Println("Загрузка команд...")
	LoadCommands()

	log.Println("Загрузка базы ответов...")
	LoadAnswers()

	var err error
	pref := telebot.Settings{
		Token:  globals.AccessToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err = telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	bot.Handle(telebot.OnText, OnMessage)

	shutdown.Add(func() {
		log.Println("Остановка бота...")
		bot.Stop()
		log.Println("Пока :(")
		os.Exit(1)
	})

	log.Println("Лонгпул запущен")
	bot.Start()
}

func OnMessage(c telebot.Context) error {
	msg := c.Message()
	mText := strings.ToLower(msg.Text)
	chatType := c.Chat().Type

	// 1. Если это ЛС — обрабатываем сразу без обращений
	if chatType == telebot.ChatPrivate {
		go OnMessageToBot(c, "")
		return nil
	}

	// 2. Если это группа/супергруппа
	if chatType == telebot.ChatGroup || chatType == telebot.ChatSuperGroup {

		// Проверяем: если это ответ (reply) на сообщение нашего бота
		if msg.ReplyTo != nil && msg.ReplyTo.Sender.ID == bot.Me.ID {
			// Пользователь общается напрямую с ботом внутри цепочки ответов, appeal не нужен
			go OnMessageToBot(c, "")
			return nil
		}

		if StupidMessagesEnable {
			// Сначала проверяем, нет ли в тексте триггеров для нубиков
			for _, trigger := range StupidMessagesTriggers {
				if strings.Contains(mText, trigger.(string)) {
					// Найдено ключевое слово! Обрабатываем без appeal
					go OnStupidMessage(c) // Оставил вашу функцию для триггеров, как в вашем примере
					return nil
				}
			}
		}

		// Если триггеров не нашлось, проверяем стандартное обращение (Appeals)
		for _, a := range Appeals {
			appealStr := strings.ToLower(a.(string))
			if strings.HasPrefix(mText, appealStr) {
				go OnMessageToBot(c, appealStr)
				break
			}
		}
	}

	return nil
}

func OnMessageToBot(c telebot.Context, appeal string) {
	msg := c.Message()
	rawText := strings.ToLower(msg.Text)

	if appeal != "" {
		rawText = strings.Replace(rawText, appeal, "", 1)
	}

	rawText = strings.TrimSpace(rawText)
	args := strings.Split(rawText, " ")

	if len(args) == 0 || args[0] == "" {
		return
	}

	if val, ok := commandList[args[0]]; ok {
		val.Function(c, args)
	} else {
		if globals.HasAnswers {
			chatbot.FindAndSendAnswer(c, rawText)
		}
	}
}

func OnStupidMessage(c telebot.Context) {
	c.Reply(StupidMessagesText)
}
