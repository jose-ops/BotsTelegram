package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gde "github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	Timeout              int    `json:"timeout"`
	UsdMinoristaEndpoint string `json:"usd_minorista"`
	Debug                bool   `json:"debug"`
}

type Env struct {
	BotToken string
	ApiToken string
}

//
type Coti struct {
	Value float32 `json:"v"`
	Fecha string  `json:"d"`
}

var env = &Env{}
var configuration = &Configuration{}

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := gde.Load(); err != nil {
		log.Print("No .env file found")
	}
	apiToken, exists := os.LookupEnv("API_TOKEN")
	if !exists {
		panic("No API_TOKEN")
	}
	botToken, exists := os.LookupEnv("BOT_TOKEN")
	if !exists {
		panic("No API_TOKEN")
	}
	env.BotToken = botToken
	env.ApiToken = apiToken
	file, err := os.Open("config.json")
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return
	}
}

func main() {

	bot, err := tgbotapi.NewBotAPI(env.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = configuration.Debug
	u := tgbotapi.NewUpdate(0)
	u.Timeout = configuration.Timeout
	updates, err := bot.GetUpdatesChan(u)
	client := &http.Client{}
	req, err := http.NewRequest("GET", configuration.UsdMinoristaEndpoint, nil)

	authValue := fmt.Sprintf("BEARER %s", env.ApiToken)
	req.Header.Add("Authorization", authValue)

	for update := range updates {
		if update.Message.Text == "/coti" {
			var coti Coti
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
