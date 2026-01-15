package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"OpenSAI/internal/bot"
	"OpenSAI/internal/engines/smlV1"
	"OpenSAI/pkg/formatting"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startTime time.Time
var botUsername string
var lastResponseTime time.Time
var responseCount int
var engine = "SMLv1"

const (
	maxResponsesPerMinute = 5
	softLimit             = 50
	hardLimit             = 80
)

func main() {
	startTime = time.Now()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in main: %v", r)
		}
	}()

	token, err := bot.ReadToken("config/token.bot")
	if err != nil {
		log.Panic("Error reading token file:", err)
	}

	names, err := bot.ReadNames("config/names.bot")
	if err != nil {
		log.Panic("Error reading names file:", err)
	}

	responses, err := bot.ReadLines("config/response.bot")
	if err != nil {
		log.Panic("Error reading response file:", err)
	}

	vocabulary, err := bot.ReadLines("vocabulary.bot")
	if err != nil {
		log.Panic("Error reading vocabulary file:", err)
	}

	infoLines, err := bot.ReadLines("config/info.bot")
	if err != nil {
		log.Panic("Error reading info file:", err)
	}

	if len(infoLines) < 2 {
		log.Panic("info.bot file must contain at least two lines")
	}

	infoText := infoLines[0]
	botUsername = infoLines[1]

	rand.Seed(time.Now().UnixNano())

	telegramBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	telegramBot.Debug = true
	log.Printf("Authorized on account %s", telegramBot.Self.UserName)

	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in update loop: %v", r)
					time.Sleep(5 * time.Second)
				}
			}()

			u := tgbotapi.NewUpdate(0)
			u.Timeout = 60

			updates := telegramBot.GetUpdatesChan(u)
			log.Println("Bot started listening for updates...")

			messageCount := 0
			var lastMessages []string

			for update := range updates {
				go func(update tgbotapi.Update) {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Recovered from panic in message processing: %v", r)
						}
					}()

					if update.Message == nil {
						return
					}

					log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

					messageText := strings.TrimSpace(update.Message.Text)

					if int64(update.Message.Date) < startTime.Unix() {
						log.Println("The message was sent before the bot was launched, ignore it.")
						return
					}

					now := time.Now()
					if now.Sub(lastResponseTime) < time.Minute {
						if responseCount >= hardLimit {
							log.Println("Hard rate limit exceeded - skipping message")
							return
						} else if responseCount >= softLimit {
							delay := time.Duration(responseCount-softLimit+1) * time.Second
							log.Printf("Soft limit exceeded, adding %v delay", delay)
							time.Sleep(delay)
						}
					} else {
						responseCount = 0
					}

					responseCount++
					lastResponseTime = now

					if strings.HasPrefix(messageText, "/") {
						cmdParts := strings.Split(messageText, "@")
						cmd := cmdParts[0]
						cmdUser := ""

						if len(cmdParts) > 1 {
							cmdUser = cmdParts[1]
						}

						if cmdUser == "" || cmdUser == botUsername {
							switch cmd {
							case "/fetch":
								uptime := time.Since(startTime)
								uptimeFormatted := formatting.Duration(uptime)
								infoTextEscaped := formatting.EscapeMarkdownV2(infoText)

								phrases, err := countLines("vocabulary.bot")
								if err != nil {
									log.Printf("Error counting lines: %v", err)
									phrases = "<err>"
								}

								msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n————————————————————\n*Время жизни:* %s\n*Движок ответов:* %s\n*Словарный запас\\:* %s фраз\n\n_Powered by [OpenSAI](https://github.com/z3nnix/openSAI)_", infoTextEscaped, uptimeFormatted, engine, phrases))
								msg.ParseMode = "MarkdownV2"
								if _, err := telegramBot.Send(msg); err != nil {
									log.Printf("Error sending message: %v", err)
									time.Sleep(2 * time.Second)
								}
								return

							case "/amnesia":
								bot.AmnesiaHandler(telegramBot, update)
								return
							}
						}
					}

					smlV1.ProcessMessage(telegramBot, update, names, responses, vocabulary, &lastMessages, &messageCount)

				}(update)
			}
		}()

		log.Println("Update channel closed, reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func countLines(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strconv.Itoa(count), nil
}

func sendMessageWithRetry(telegramBot *tgbotapi.BotAPI, msg tgbotapi.Chattable) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		_, err = telegramBot.Send(msg)
		if err == nil {
			return nil
		}
		log.Printf("Send attempt %d failed: %v", attempt+1, err)
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}
	return err
}
