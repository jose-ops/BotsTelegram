package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//
type Coti struct {
	Value float32 `json:"v"`
	Fecha string  `json:"d"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI("830198805:AAHdRi2Xy0-rtfczMW3aC5-8MF5UKcUykpI")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	//log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Text == "/coti" {
			var coti Coti
			client := &http.Client{}
			req, err := http.NewRequest("GET", "https://api.estadisticasbcra.com/usd_of_minorista", nil)

			req.Header.Add("Authorization", "BEARER eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDA1NjMwNzQsInR5cGUiOiJleHRlcm5hbCIsInVzZXIiOiJuYWNoby5uaWV2YUBvdXRsb29rLmNvbSJ9.nsGACzjVB9CyzHS8uCk_g6FJZq9lCq9G9Sx-H1zFLWRqcCmCP1fnMynLF6h-yEj7JMDhSTBE4puCesPFjtqF4Q")
			resp, err := client.Do(req)

			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "oops this didn't work"+err.Error())
				bot.Send(msg)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			contentSlice := strings.Split(string(body), ",")

			final := strings.Replace(contentSlice[len(contentSlice)-2]+","+contentSlice[len(contentSlice)-1], "]", "", -1)

			json.Unmarshal([]byte(final), &coti)
			fmt.Println(final, coti)
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("La cotizacion es de %.2f de la fecha "+coti.Fecha, coti.Value))

			bot.Send(msg)
		} else {
			fmt.Println("ole")
		}
	}
}

//eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDA1NjMwNzQsInR5cGUiOiJleHRlcm5hbCIsInVzZXIiOiJuYWNoby5uaWV2YUBvdXRsb29rLmNvbSJ9.nsGACzjVB9CyzHS8uCk_g6FJZq9lCq9G9Sx-H1zFLWRqcCmCP1fnMynLF6h-yEj7JMDhSTBE4puCesPFjtqF4Q
//Authorization: BEARER
