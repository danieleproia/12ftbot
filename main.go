package main

import (
        "log"
        "os"
        "path/filepath"
        "github.com/go-telegram-bot-api/telegram-bot-api"
        "github.com/joho/godotenv"
)

func main() {
        exePath, err := os.Executable()
        if err != nil {
                log.Fatal("Failed to get executable path:", err)
        }

        scriptDir := filepath.Dir(exePath)
        envFilePath := filepath.Join(scriptDir, ".env")

        err = godotenv.Load(envFilePath)
        if err != nil {
                log.Fatal("Error loading .env file:", err)
        }

        botToken := os.Getenv("BOT_TOKEN")
        if botToken == "" {
                log.Fatal("BOT_TOKEN environment variable is not set.")
        }

        bot, err := tgbotapi.NewBotAPI(botToken)
        if err != nil {
                log.Fatal("Failed to create Telegram bot:", err)
        }

        bot.Debug = true

        log.Printf("Authorized on account %s", bot.Self.UserName)

        u := tgbotapi.NewUpdate(0)
        u.Timeout = 60

        updates, err := bot.GetUpdatesChan(u)
        if err != nil {
                log.Fatal("Failed to start update channel:", err)
        }

        for update := range updates {
                if update.InlineQuery != nil {
                        go handleInlineQuery(bot, update.InlineQuery)
                } else if update.Message != nil && update.Message.Text != "" {
                        go handleMessage(bot, update.Message)
                }
        }
}

func handleInlineQuery(bot *tgbotapi.BotAPI, inlineQuery *tgbotapi.InlineQuery) {
        query := inlineQuery.Query
        if query == "" {
                return
        }

        messageURL := "https://12ft.io/" + query
        inlineConfig := tgbotapi.InlineConfig{
                InlineQueryID: inlineQuery.ID,
                Results: []interface{}{
                        tgbotapi.NewInlineQueryResultArticleHTML(
                                inlineQuery.ID,
                                query,
                                messageURL,
                        ),
                },
        }

        _, err := bot.AnswerInlineQuery(inlineConfig)
        if err != nil {
                log.Println("Failed to answer inline query:", err)
        }
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
        if message.IsCommand() {
                return
        }

        messageURL := "https://12ft.io/" + message.Text

        msg := tgbotapi.NewMessage(message.Chat.ID, messageURL)
        msg.ReplyToMessageID = message.MessageID

        _, err := bot.Send(msg)
        if err != nil {
                log.Println("Failed to send message:", err)
        }
}