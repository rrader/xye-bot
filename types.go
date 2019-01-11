package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

type Request struct {
	update *tgbotapi.Update
}

type Command struct {
	request *Request
}

type CommandIF interface {
	Handle() error
}
