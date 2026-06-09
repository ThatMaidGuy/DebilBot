package commands

import (
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/telebot.v4"
)

func Rate(c telebot.Context, args []string) {
	rating := rand.Intn(11)

	c.Reply("Оцениваю на " + strconv.Itoa(rating) + " из 10")
}

func Time(c telebot.Context, args []string) {
	dt := time.Now()
	c.Reply("Ну, у меня " + dt.Format("15:04") + ". А у тебя?")
}
