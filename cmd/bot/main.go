package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startTime time.Time
var botUsername string

func main() {
	startTime = time.Now()

	token, err := readToken("config/token.bot")
	if err != nil {
		log.Panic("Error reading token file:", err)
	}

	names, err := readNames("config/names.bot")
	if err != nil {
		log.Panic("Error reading names file:", err)
	}

	responses, err := readLines("config/response.bot")
	if err != nil {
		log.Panic("Error reading response file:", err)
	}

	vocabulary, err := readLines("vocabulary.bot")
	if err != nil {
		log.Panic("Error reading vocabulary file:", err)
	}

	infoLines, err := readLines("config/info.bot")
	if err != nil {
		log.Panic("Error reading info file:", err)
	}

	if len(infoLines) < 2 {
		log.Panic("info.bot file must contain at least two lines")
	}

	infoText := infoLines[0]
	botUsername = infoLines[1]

	rand.Seed(time.Now().UnixNano())

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	messageCount := 0
	var lastMessages []string

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		messageText := strings.TrimSpace(update.Message.Text)

		if strings.HasPrefix(messageText, "/") {
			cmd := strings.Split(messageText, "@")[0]
			cmdUser := ""
			if strings.Contains(messageText, "@") {
				cmdUser = strings.Split(messageText, "@")[1]
			}

			if cmdUser == "" || cmdUser == botUsername {
				switch cmd {
				case "/fetch":
					uptime := time.Since(startTime)
					uptimeFormatted := formatDuration(uptime)
					infoTextEscaped := escapeMarkdownV2(infoText)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n————————————————————\n*Время жизни:* %s\n\n_Powered by [OpenSAI](https://github.com/z3nnix/openSAI)_", infoTextEscaped, uptimeFormatted))
					msg.ParseMode = "MarkdownV2"
					bot.Send(msg)
				}
			}
			continue
		}

		processMessage(bot, update, names, responses, vocabulary, &lastMessages, &messageCount)
	}
}