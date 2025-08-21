package main

import (
    "fmt"
    "log"
    "math/rand"
    "strings"
    "time"
    "os"
    "bufio"
    "strconv"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var startTime time.Time
var botUsername string
var lastResponseTime time.Time
var responseCount int
const maxResponsesPerMinute = 5

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

        if int64(update.Message.Date) < startTime.Unix() {
            log.Println("The message was sent before the bot was launched, ignore it.")
            continue
        }

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
                    uptimeFormatted := formatDuration(uptime)
                    infoTextEscaped := escapeMarkdownV2(infoText)

                    phrases, err := countLines("vocabulary.bot")
                    if err != nil {
                        fmt.Printf("Ошибка: %v\n", err)
                        phrases = "<err>"
                    }

                    msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n————————————————————\n*Время жизни:* %s\n*Движок ответов:* %s\n*Словарный запас\\:* %s фраз\n\n_Powered by [OpenSAI](https://github.com/z3nnix/openSAI)_", infoTextEscaped, uptimeFormatted, engine, phrases))
                    msg.ParseMode = "MarkdownV2"
                    bot.Send(msg)

                case "/amnesia":
                    amnesiaHandler(bot, update)
                }
                continue
            }
        }

        // Проверка частоты ответов
        if time.Since(lastResponseTime) < time.Minute {
            if responseCount >= maxResponsesPerMinute {
                log.Println("Rate limit exceeded, ignoring message.")
                continue
            }
        } else {
            responseCount = 0
        }

        processMessage(bot, update, names, responses, vocabulary, &lastMessages, &messageCount)

        // Обновление времени последнего ответа и счетчика
        lastResponseTime = time.Now()
        responseCount++
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