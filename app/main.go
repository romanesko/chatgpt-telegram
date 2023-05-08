package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
	"time"
)

var TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
var OpenAiToken = os.Getenv("OPENAI_TOKEN")
var AuthPassword = os.Getenv("AUTH_PASSWORD")

var requests = NewRequests(OpenAiToken)

func requestChatGpt(text string) string {

	params := NewGtpRequest(text)
	var response GptResponse
	err := requests.post("https://api.openai.com/v1/chat/completions", &params, &response)
	if err != nil {
		return err.Error()
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("Something wrong with response, please check the logs")
		}
	}()

	if response.Error.Message != "" {
		return response.Error.Message
	}

	return response.Choices[0].Message.Content
}

type void struct{}

func main() {

	if TelegramBotToken == "" || OpenAiToken == "" {
		log.Fatal("Please set TELEGRAM_BOT_TOKEN and OPENAI_TOKEN environment variables")
	}

	bot, err := tgbotapi.NewBotAPI(TelegramBotToken)

	if err != nil {
		log.Panic(err)
	}

	allowedChats := getChatsListStore()

	log.Println("Allowed chats:", allowedChats.listChats())

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message

			if strings.HasPrefix(update.Message.Text, "/auth") {
				password := strings.SplitN(update.Message.Text, "/auth", 2)
				if strings.TrimSpace(password[1]) == AuthPassword {
					allowedChats.addChat(update.Message.Chat.ID)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Ask me a question")
					bot.Send(msg)
				} else {
					time.Sleep(2 * time.Second)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong password!")
					bot.Send(msg)
				}
				continue
			}

			if !allowedChats.isChatAllowed(update.Message.Chat.ID) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know you! Please send me a message:\n/auth <password>")
				bot.Send(msg)
				continue
			}

			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			processingMessage, _ := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Processing..."))
			delMsgConfig := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, processingMessage.MessageID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, requestChatGpt(update.Message.Text))
			msg.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(delMsgConfig)
			bot.Send(msg)
		}
	}

}
