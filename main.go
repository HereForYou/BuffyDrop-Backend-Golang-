package main

import (
	"fmt"
	"log"

	"go-test/config"
	"go-test/db"
	"go-test/routes/setting"
	"go-test/routes/user"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/mux"
)

type User struct {
	UserName    string  `bson:"userName"`
	TgId        string  `bson:"tgId"`
	Email       string  `bson:"email"`
	TotalPoints float32 `bson:"totalPoints"`
}

func main() {
	// load configuration data from config package
	cfg := config.LoadConfig()

	//================================================================================== setting router
	router := mux.NewRouter()
	user_router.RegisterUserRoute(router)
	setting_router.RegisterUserRoute(router)

	//================================================================================== Connect to DB
	db.Connect(cfg.DbUrl)

	//================================================================================== Setting Telegram Bot
	botToken := cfg.BotToken
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully set the bot", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Println(update.Message.From.UserName, update.Message.Text)

		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the Bot!")
			bot.Send(msg)
		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here are the available commands:\n/start - Start the bot\n/help - Show help")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I cann't understand that command!")
			bot.Send(msg)
		}
	}

	fmt.Println("Server is running on port: ", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, router)
}
