package chatbot

import (
	"DebilBot/globals"

	"gopkg.in/telebot.v4"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

func FindAndSendAnswer(c telebot.Context, full_text string) {
	if !globals.HasAnswers {
		return
	}

	// log.Println("Ищем ответ для \"" + full_text + "\"")

	var last_similarity float64
	var last_sim_item []string

	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false

	for _, v := range globals.FullBase {
		sim := strutil.Similarity(full_text, v[0], swg)
		if sim > last_similarity {
			last_similarity = sim
			last_sim_item = v
		}
	}

	// log.Println("Отправлен ответ с \"" + last_sim_item[1] + "\" с коэфф " + fmt.Sprintf("%f", last_similarity))

	c.Reply(last_sim_item[1])
}
