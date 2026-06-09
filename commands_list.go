package main

import (
	"DebilBot/commands"
	"DebilBot/globals"

	"gopkg.in/telebot.v4"
)

// Теперь функция команды принимает telebot.Context вместо VK Message
type commandFunction = func(c telebot.Context, args []string)

type Command struct {
	Name        string
	Description string
	Icon        string
	Function    commandFunction
	isHidden    bool
}

var (
	commandList map[string]Command
)

func AllCommands(c telebot.Context, args []string) {
	resultText := "Все команды:\n\n"

	for _, v := range commandList {
		if !v.isHidden {
			resultText = resultText + v.Icon + " " + v.Name + " - " + v.Description + "\n"
		}
	}

	// c.Reply автоматически отвечает реплаем на исходное сообщение пользователя
	_ = c.Reply(resultText)
}

func HelpCommand(c telebot.Context, args []string) {
	infoText, ok := globals.BotSettings.Get("info_text").(string)
	if !ok {
		infoText = "Информация о боте временно недоступна."
	}

	_ = c.Reply(infoText)
}

func LoadCommands() {
	commandList = make(map[string]Command)

	commandList["помощь"] = Command{
		Name:        "помощь",
		Description: "Информация о боте",
		Icon:        "🚑",
		Function:    HelpCommand,
		isHidden:    false,
	}
	commandList["команды"] = Command{
		Name:        "команды",
		Description: "Все команды",
		Icon:        "📱",
		Function:    AllCommands,
		isHidden:    false,
	}
	commandList["тест"] = Command{
		Name:        "пинг",
		Description: "Проверяет бота",
		Icon:        "💡",
		Function:    commands.TestCommand,
		isHidden:    false,
	}
	commandList["время"] = Command{
		Name:        "время",
		Description: "Показать время у бота",
		Icon:        "⏰",
		Function:    commands.Time,
		isHidden:    false,
	}
	commandList["оцени"] = Command{
		Name:        "оцени",
		Description: "Объективная оценка от бота",
		Icon:        "💯",
		Function:    commands.Rate,
		isHidden:    false,
	}
}
