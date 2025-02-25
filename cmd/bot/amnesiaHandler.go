package main

import (
	"fmt"
	"log"
	"bufio"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
) 

func amnesiaHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update,) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅")

		chatConfig := tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: update.Message.Chat.ID,
			},
		}

		admins, err := bot.GetChatAdministrators(chatConfig)
		if err != nil {
			log.Printf("%s", err)
			return
		}
		userID := update.Message.From.ID
		isAdmin := false
		for _, admin := range admins {
			if admin.User.ID == userID {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌")
			bot.Send(msg)
			return
		}

		file, err := os.Open("neglected.bot")
		if err != nil {
			log.Printf("%s", err)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		chatIDStr := fmt.Sprintf("%d", update.Message.Chat.ID)
		var lines []string
		found := false

		for scanner.Scan() {
			line := scanner.Text()
			if line == chatIDStr {
				found = true
			} else {
				lines = append(lines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("%s", err)
			return
		}

		if !found {
			lines = append(lines, chatIDStr)
			log.Println("ID added")
		}

		file, err = os.OpenFile("neglected.bot", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("%s", err)
			return
		}
		defer file.Close()

		for _, line := range lines {
			if _, err := file.WriteString(line + "\n"); err != nil {
				log.Printf("%s", err)
				return
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("%s", err)
		} else {
			log.Println("message sent")
		}
}