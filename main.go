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
	"github.com/rs/cors"
)

type User struct {
	UserName    string  `bson:"userName"`
	TgId        string  `bson:"tgId"`
	Email       string  `bson:"email"`
	TotalPoints float32 `bson:"totalPoints"`
}

func main() {
	//================================================================================== load configuration data from config package
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
		log.Fatal("🔴" + err.Error())
	}

	fmt.Println("🔵 Successfully set the bot", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("🔴" + err.Error())
	}

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Println("🟡 Request from ", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Text {
			case "/start":
				desText := "This is SmartFox Bot written by Golang  🍖"
				photo := tgbotapi.NewPhotoShare(update.Message.Chat.ID, "https://i.ibb.co/cvzgGD6/pic.png")
				photo.Caption = desText
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("Let's Go", "https://t.me/gosmartfoxbot"),
						tgbotapi.NewInlineKeyboardButtonURL("Join Smart Community", "https://t.me/gosmartfoxbot"),
					),
				)
				photo.ReplyMarkup = keyboard
				if _, err := bot.Send(photo); err != nil {
					log.Println("Error", err)
				}
			case "/help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here are the available commands:\n/start - Start the bot\n/help - Show help")
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I can't understand that command!")
				bot.Send(msg)
			}
		}
	}()

	//===================================================================================== CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://foo.com:8080"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	fmt.Println("🔵 Server is running on port: ", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, c.Handler(router))
	if err != nil {
		fmt.Println("🔴 " + err.Error())
	}
}
