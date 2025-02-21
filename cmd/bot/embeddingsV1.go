package main

import (
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var engine = "EmbeddingsV1"

func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, names map[string]bool, responses []string, vocabulary []string, lastMessages *[]string, messageCount *int) {
	messageText := strings.TrimSpace(update.Message.Text)

	var msg tgbotapi.MessageConfig

	if names[messageText] {
		similarWord := findMostSimilarWord(vocabulary, messageText)

		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typing)
		time.Sleep(2 * time.Second)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, similarWord)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		*lastMessages = append(*lastMessages, similarWord)
		if len(*lastMessages) > 15 {
			*lastMessages = (*lastMessages)[1:]
		}
	} else if !strings.HasPrefix(messageText, "/") {
		appendMessageToFile("vocabulary.bot", messageText)
	}

	*messageCount++

	if *messageCount%10 == 0 {
		randomVocabulary := getRandomUniqueWord(vocabulary, *lastMessages)

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
		randomVocabulary := getRandomUniqueWord(vocabulary, *lastMessages)

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

func findMostSimilarWord(vocabulary []string, input string) string {
	minDistance := -1
	var mostSimilarWord string

	for _, word := range vocabulary {
		distance := levenshteinDistance(input, word)
        if minDistance == -1 || distance < minDistance {
            minDistance = distance
            mostSimilarWord = word
        }
    }

    return mostSimilarWord
}

func levenshteinDistance(a, b string) int {
	al := len(a)
    bl := len(b)

    dp := make([][]int, al+1)
    for i := range dp {
        dp[i] = make([]int, bl+1)
    }

    for i := 0; i <= al; i++ {
        dp[i][0] = i
    }
    for j := 0; j <= bl; j++ {
        dp[0][j] = j
    }

    for i := 1; i <= al; i++ {
        for j := 1; j <= bl; j++ {
            cost := 0
            if a[i-1] != b[j-1] {
                cost = 1
            }
            dp[i][j] = min(dp[i-1][j]+1, min(dp[i][j-1]+1, dp[i-1][j-1]+cost))
        }
    }

    return dp[al][bl]
}

func getRandomUniqueWord(vocabulary []string, lastMessages []string) string {
	var randomVocabulary string

	for {
        randomVocabulary = vocabulary[rand.Intn(len(vocabulary))]
        if !contains(lastMessages, randomVocabulary) {
            break
        }
    }

    return randomVocabulary
}

func min(a, b int) int {
	if a < b {
        return a 
    }
	return b 
}
