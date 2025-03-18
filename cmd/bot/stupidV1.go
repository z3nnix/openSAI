package main

import (
	"math/rand"
	"strings"
	"time"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var engine = "StupidV1"

func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, names map[string]bool, responses []string, vocabulary []string, lastMessages *[]string, messageCount *int) {
	messageText := strings.TrimSpace(update.Message.Text)

	var msg tgbotapi.MessageConfig

	if names[messageText] {
		var randomResponse string
		for {
			randomResponse = responses[rand.Intn(len(responses))]
			if !contains(*lastMessages, randomResponse) {
				break
			}
		}

		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typing)
		time.Sleep(2 * time.Second)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, randomResponse)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		*lastMessages = append(*lastMessages, randomResponse)
		if len(*lastMessages) > 15 {
			*lastMessages = (*lastMessages)[1:]
		}
	} else if !strings.HasPrefix(messageText, "/") {
		chatid := strconv.FormatInt(update.Message.Chat.ID, 10)
		appendMessageToFile("vocabulary.bot", messageText, chatid)
	}

	*messageCount++

	if *messageCount%200 == 0 {
		var randomVocabulary string
		for {
			randomVocabulary = vocabulary[rand.Intn(len(vocabulary))]
			if !contains(*lastMessages, randomVocabulary) {
				break
			}
		}

		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typing)
		time.Sleep(2 * time.Second)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, randomVocabulary)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		*lastMessages = append(*lastMessages, randomVocabulary)
		if len(*lastMessages) > 15 {
			*lastMessages = (*lastMessages)[1:]
		}
	}

	if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.UserName == bot.Self.UserName {
		var randomVocabulary string
		for {
			randomVocabulary = vocabulary[rand.Intn(len(vocabulary))]
			if !contains(*lastMessages, randomVocabulary) {
				break
			}
		}

		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typing)
		time.Sleep(2 * time.Second)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, randomVocabulary)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		*lastMessages = append(*lastMessages, randomVocabulary)
		if len(*lastMessages) > 15 {
			*lastMessages = (*lastMessages)[1:]
		}
	}
}