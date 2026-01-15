package stupidV1

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"OpenSAI/pkg/formatting"
	"OpenSAI/pkg/vocman"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var engine = "StupidV1"

func ProcessMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, names map[string]bool, responses []string, vocabulary []string, lastMessages *[]string, messageCount *int) {
	messageText := strings.TrimSpace(update.Message.Text)
	chatID := strconv.FormatInt(update.Message.Chat.ID, 10)

	var msg tgbotapi.MessageConfig

	if names[messageText] {
		var randomResponse string
		for {
			randomResponse = responses[rand.Intn(len(responses))]
			if !formatting.Contains(*lastMessages, randomResponse) {
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
		vocman.AppendMessageToFile("vocabulary.bot", messageText, chatID)
	}

	*messageCount++

	// Получаем словарь для текущего чата (основной или локальный)
	chatVocabulary := vocman.GetChatVocabulary(chatID, vocabulary)

	if *messageCount%200 == 0 {
		if len(chatVocabulary) > 0 {
			var randomVocabulary string
			for {
				randomVocabulary = chatVocabulary[rand.Intn(len(chatVocabulary))]
				if !formatting.Contains(*lastMessages, randomVocabulary) {
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

	if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.UserName == bot.Self.UserName {
		if len(chatVocabulary) > 0 {
			var randomVocabulary string
			for {
				randomVocabulary = chatVocabulary[rand.Intn(len(chatVocabulary))]
				if !formatting.Contains(*lastMessages, randomVocabulary) {
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
}
