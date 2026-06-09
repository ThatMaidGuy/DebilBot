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

	// Хэндлер для обычного текста
	bot.Handle(telebot.OnText, OnMessage)

	// --- НОВЫЕ ХЭНДЛЕРЫ ДЛЯ МЕДИА ---
	// Вешаем одну и ту же функцию OnMessage на все типы сообщений с медиа
	bot.Handle(telebot.OnPhoto, OnMessage)
	bot.Handle(telebot.OnVideo, OnMessage)
	bot.Handle(telebot.OnAudio, OnMessage)
	bot.Handle(telebot.OnDocument, OnMessage)
	bot.Handle(telebot.OnVoice, OnMessage)
	bot.Handle(telebot.OnAnimation, OnMessage) // Гифки
	bot.Handle(telebot.OnSticker, OnMessage)   // На случай, если у стикера есть эмодзи, которые вы хотите парсить

	shutdown.Add(func() {
		log.Println("Остановка бота...")
		bot.Stop()
		log.Println("Пока :(")
		os.Exit(1)
	})

	log.Println("Лонгпул запущен")
	bot.Start()
}

// Вспомогательная функция для получения текста ИЛИ подписи к медиа
func getMessageText(msg *telebot.Message) string {
	if msg.Text != "" {
		return msg.Text
	}
	// Если это медиафайл, текст будет лежать в Caption
	return msg.Caption
}

func OnMessage(c telebot.Context) error {
	msg := c.Message()

	// Используем вспомогательную функцию вместо msg.Text
	textToProcess := getMessageText(msg)
	mText := strings.ToLower(textToProcess)
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
					for _, chat_id := range StupidMessagesChatIDs {
						if c.Chat().ID == chat_id.(int64) {
							go OnStupidMessage(c)
							return nil
						}
					}
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

	// Тоже используем извлечение текста/подписи
	rawText := strings.ToLower(getMessageText(msg))

	if appeal != "" {
		rawText = strings.Replace(rawText, appeal, "", 1)
	}

	rawText = strings.TrimSpace(rawText)
	args := strings.Split(rawText, " ")

	// Если текста и подписи нет вообще (например, отправили просто фото без текста)
	if len(args) == 0 || args[0] == "" {
		switch {
		case msg.Photo != nil:
			c.Reply("Красивое фото!")
		case msg.Video != nil:
			c.Reply("Длинное видео, не буду смотреть.")
		case msg.Voice != nil || msg.VideoNote != nil: // Заодно поймаем и круглые видео-сообщения (кружочки)
			c.Reply("Не, мне лень слушать.")
		case msg.Document != nil:
			c.Reply("Да мне как то все равно, что ты там отправляешь...")
		case msg.Animation != nil:
			c.Reply("Гифка топ!")
		case msg.Sticker != nil:
			c.Reply(&telebot.Sticker{File: telebot.File{FileID: "CAACAgIAAxkBAAFL8n5qKAc8xSr8mKA7PlkkIwtFd4p-AgAC-Z4AAgfZ8EgRAeZCFgnt6jsE"}})
		default:
			return
		}
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
