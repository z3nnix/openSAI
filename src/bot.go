package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startTime time.Time

func main() {
	startTime = time.Now()

	botToken, err := ioutil.ReadFile("token.bot")
	if err != nil {
		log.Panic("Error reading token file:", err)
	}

	if len(botToken) == 0 {
		log.Panic("Token file is empty")
	}

	namesFile, err := os.Open("names.bot")
	if err != nil {
		log.Panic("Error reading names file:", err)
	}
	defer namesFile.Close()

	names := make(map[string]bool)
	scanner := bufio.NewScanner(namesFile)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			names[name] = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	responseFile, err := os.Open("response.bot")
	if err != nil {
		log.Panic("Error reading response file:", err)
	}
	defer responseFile.Close()

	var responses []string
	scanner = bufio.NewScanner(responseFile)
	for scanner.Scan() {
		response := strings.TrimSpace(scanner.Text())
		if response != "" {
			responses = append(responses, response)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	vocabularyFile, err := os.Open("vocabulary.bot")
	if err != nil {
		log.Panic("Error reading vocabulary file:", err)
	}
	defer vocabularyFile.Close()

	var vocabulary []string
	scanner = bufio.NewScanner(vocabularyFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			vocabulary = append(vocabulary, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	infoFile, err := os.Open("info.bot")
	if err != nil {
		log.Panic("Error reading info file:", err)
	}
	defer infoFile.Close()

	var infoLines []string
	scanner = bufio.NewScanner(infoFile)
	for scanner.Scan() {
		infoLines = append(infoLines, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	if len(infoLines) < 2 {
		log.Panic("info.bot file must contain at least two lines")
	}

	infoText := infoLines[0]
	botUsername := infoLines[1]

	rand.Seed(time.Now().UnixNano())

	bot, err := tgbotapi.NewBotAPI(string(botToken))
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

		var msg tgbotapi.MessageConfig

		if names[messageText] {
			var randomResponse string
			for {
				randomResponse = responses[rand.Intn(len(responses))]
				if !contains(lastMessages, randomResponse) {
					break
				}
			}

			typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
			bot.Send(typing)
			time.Sleep(2 * time.Second)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, randomResponse)
			bot.Send(msg)

			lastMessages = append(lastMessages, randomResponse)
			if len(lastMessages) > 15 {
				lastMessages = lastMessages[1:]
			}
		} else if !strings.HasPrefix(messageText, "/") {
			appendMessageToFile("vocabulary.bot", messageText)
		}

		messageCount++

		if messageCount%10 == 0 {
			var randomVocabulary string
			for {
				randomVocabulary = vocabulary[rand.Intn(len(vocabulary))]
				if !contains(lastMessages, randomVocabulary) {
					break
				}
			}

			typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
			bot.Send(typing)
			time.Sleep(2 * time.Second)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, randomVocabulary)
			bot.Send(msg)

			lastMessages = append(lastMessages, randomVocabulary)
			if len(lastMessages) > 15 {
				lastMessages = lastMessages[1:]
			}
		}

		if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.UserName == bot.Self.UserName {
			var randomVocabulary string
			for {
				randomVocabulary = vocabulary[rand.Intn(len(vocabulary))]
				if !contains(lastMessages, randomVocabulary) {
					break
				}
			}

			typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
			bot.Send(typing)
			time.Sleep(2 * time.Second)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, randomVocabulary)
			bot.Send(msg)

			lastMessages = append(lastMessages, randomVocabulary)
			if len(lastMessages) > 15 {
				lastMessages = lastMessages[1:]
			}
		}
	}
}

func appendMessageToFile(filename, message string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(message + "\n"); err != nil {
		log.Println("Error writing to file:", err)
	}
}

func formatDuration(d time.Duration) string {
	d = time.Duration(math.Ceil(d.Seconds())) * time.Second
	seconds := int(d.Seconds())

	if seconds < 60 {
		return fmt.Sprintf("%d секунд", seconds)
	}

	minutes := seconds / 60
	seconds = seconds % 60

	if minutes < 60 {
		return fmt.Sprintf("%d минут %d секунд", minutes, seconds)
	}

	hours := minutes / 60
	minutes = minutes % 60

	if hours < 24 {
		return fmt.Sprintf("%d часа %d минут %d секунд", hours, minutes, seconds)
	}

	days := hours / 24
	hours = hours % 24

	return fmt.Sprintf("%d дней %d часа %d минут %d секунд", days, hours, minutes, seconds)
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		`_`, `\_`,
		`*`, `\*`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`~`, `\~`,
		`>`, `\>`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`,
		`=`, `\=`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`.`, `\.`,
		`!`, `\!`,
	)
	return replacer.Replace(text)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}