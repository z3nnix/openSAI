package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var engine = "NgramV1"

func processMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, names map[string]bool, responses []string, vocabulary []string, lastMessages *[]string, messageCount *int) {
	messageText := strings.TrimSpace(update.Message.Text)

	var msg tgbotapi.MessageConfig
	for name := range names {
		if strings.Contains(messageText, name) {
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
			return
		}
	}
	
	if names[messageText] {
		// Делаем запрос к серверу вместо использования локальных ответов
		responseText, err := getResponseFromServer(messageText)
		if err != nil {
			// Если ошибка, используем fallback из локальных ответов
			responseText = getRandomResponse(responses, lastMessages)
		}

		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typing)
		time.Sleep(2 * time.Second)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		*lastMessages = append(*lastMessages, responseText)
		if len(*lastMessages) > 15 {
			*lastMessages = (*lastMessages)[1:]
		}
	}

	*messageCount++

	if *messageCount%200 == 0 {
		// Для периодических сообщений тоже используем сервер
		responseText, err := getResponseFromServer("random")
		if err != nil {
			responseText = getRandomVocabulary(vocabulary, lastMessages)
		}

		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typing)
		time.Sleep(2 * time.Second)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		*lastMessages = append(*lastMessages, responseText)
		if len(*lastMessages) > 15 {
			*lastMessages = (*lastMessages)[1:]
		}
	}

	if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.UserName == bot.Self.UserName {
		// Для ответов на сообщения бота тоже используем сервер
		responseText, err := getResponseFromServer(messageText)
		if err != nil {
			responseText = getRandomVocabulary(vocabulary, lastMessages)
		}

		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typing)
		time.Sleep(2 * time.Second)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		*lastMessages = append(*lastMessages, responseText)
		if len(*lastMessages) > 15 {
			*lastMessages = (*lastMessages)[1:]
		}
	}
}

// Функция для получения ответа от сервера
func getResponseFromServer(input string) (string, error) {
	// Кодируем параметр для URL (пробелы -> %20)
	encodedParam := url.QueryEscape(input)
	
	// Формируем URL запроса
	apiUrl := "http://localhost:64000/api/v1/ask?param=" + encodedParam
	
	// Выполняем HTTP GET запрос
	resp, err := http.Get(apiUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return "", http.ErrNoLocation
	}
	
	return string(body), nil
}

// Вспомогательные функции для fallback
func getRandomResponse(responses []string, lastMessages *[]string) string {
	var randomResponse string
	for {
		randomResponse = responses[rand.Intn(len(responses))]
		if !contains(*lastMessages, randomResponse) {
			break
		}
	}
	return randomResponse
}

func getRandomVocabulary(vocabulary []string, lastMessages *[]string) string {
	var randomVocabulary string
	for {
		randomVocabulary = vocabulary[rand.Intn(len(vocabulary))]
		if !contains(*lastMessages, randomVocabulary) {
			break
		}
	}
	return randomVocabulary
}