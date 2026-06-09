package chatbot

import (
	"math/rand/v2"
	"strings" // Добавили для strings.ReplaceAll и strings.ToUpper

	"DebilBot/globals"

	"github.com/adrg/strutil/metrics"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/telebot.v4"
)

func FindAndSendAnswer(c telebot.Context, full_text string) {
	if !globals.HasAnswers {
		return
	}

	// Используем наш новый метод для поиска
	match, found := findBestMatch(full_text, globals.FullBase)
	if !found {
		nf := globals.BotSettings.Get("answer_notfound").([]interface{})

		c.Reply(nf[rand.IntN(len(nf))])
		return
	}

	// Получаем сырой текст ответа
	answer := match[1]

	// 1. Получаем имя пользователя из Telegram
	username := "Пользователь" // Дефолтное имя на случай, если у юзера нет имени
	if c.Sender() != nil {
		if c.Sender().FirstName != "" {
			username = c.Sender().FirstName
		} else if c.Sender().Username != "" {
			username = c.Sender().Username
		}
	}

	// 2. Готовим варианты имени (Капсом и с Большой буквы)
	usernameUpper := strings.ToUpper(username)

	// Для безопасного приведения первой буквы к верхнему регистру (с учетом кириллицы)
	// используем пакет golang.org/x/text/cases.
	// Если пакет не установлен, выполните в терминале: go get golang.org/x/text
	caser := cases.Title(language.Russian)
	usernameTitle := caser.String(username)

	// 3. Делаем замены в тексте ответа
	answer = strings.ReplaceAll(answer, "%USERNAME%", usernameUpper)
	answer = strings.ReplaceAll(answer, "%username%", usernameTitle)

	// Отправляем измененный текст
	c.Reply(answer)
}

// findBestMatch ищет наиболее похожую фразу в базе данных.
// Возвращает массив строк (строку базы) и true, если совпадение найдено.
// Если совпадений выше порога нет, возвращает nil и false.
func findBestMatch(userText string, base [][]string) ([]string, bool) {
	// Минимальный порог схожести (0.4 — оптимально, чтобы отсеять бред.
	// Чем ближе к 1.0, тем точнее должно быть совпадение).
	const minSimilarity = 0.5

	var maxSim = minSimilarity
	var bestMatch []string

	// Настраиваем метрику один раз перед циклом
	swg := &metrics.Levenshtein{
		CaseSensitive: false,
		InsertCost:    1,
		DeleteCost:    1,
		ReplaceCost:   1,
	}

	for _, row := range base {
		// Защита: проверяем, что в строке базы есть как минимум триггер и ответ
		if len(row) < 2 {
			continue
		}

		if sim := swg.Compare(userText, row[0]); sim > maxSim {
			maxSim = sim
			bestMatch = row
		}
	}

	// Если bestMatch остался пустым, значит никто не преодолел порог minSimilarity
	if bestMatch == nil {
		return nil, false
	}

	return bestMatch, true
}
