package main

import (
	"github.com/dgraph-io/badger"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"os"
	"time"
)

var DB *badger.DB = nil
var BOT *tgbotapi.BotAPI = nil
var ChatMotions *ChatMotionsClient = nil

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	ChatMotions = ChatMotionsOpen()

	// init db
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/xyebot_db"
	opts.ValueDir = "/tmp/xyebot_db"
	var err error
	DB, err = badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	// init bot
	BOT, err = tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	BOT.Debug = true

	log.Printf("Authorized on account %s", BOT.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := BOT.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		request := NewRequest(&update)
		request.Handle()
	}
}
