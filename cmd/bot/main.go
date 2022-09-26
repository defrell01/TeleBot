package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/defrell01/telebot/pkg/repository"
	"github.com/defrell01/telebot/pkg/repository/boltdb"
	"github.com/defrell01/telebot/pkg/server"
	"github.com/defrell01/telebot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5508345363:AAGNx-uv_q5r_w1EY93oXHEuEcD3GP-6pjs")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("103119-27bacb6de2d28bb7378660c")

	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()

	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "https://localhost/") // мб ошибка с url

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/pocket_defrell_bot")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
