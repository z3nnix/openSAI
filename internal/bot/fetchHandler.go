package bot

import (
	"fmt"
	"time"

	"OpenSAI/pkg/formatting"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func fetchHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, startTime time.Time, infoText string, engine string) {
	uptime := time.Since(startTime)
	uptimeFormatted := formatting.Duration(uptime)
	infoTextEscaped := formatting.EscapeMarkdownV2(infoText)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n————————————————————\n*Время жизни:* %s\n*Движок ответов:* %s\n\n_Powered by [OpenSAI](https://github.com/z3nnix/openSAI)_", infoTextEscaped, uptimeFormatted, engine))
	msg.ParseMode = "MarkdownV2"

	if _, err := bot.Send(msg); err != nil {
		fmt.Printf("Ошибка при отправке сообщения: %s\n", err)
	}
}
