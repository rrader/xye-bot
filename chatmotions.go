package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{}

type ChatMotionsClient struct {
	url string
}


func ChatMotionsOpen() *ChatMotionsClient {
	self := &ChatMotionsClient{
		url: "http://127.0.0.1:8080",
	}
	return self
}

func (self *ChatMotionsClient) HookNewMessage(message *tgbotapi.Message) {
	if message.Text == "" {
		return
	}

	url := self.url + "/message/" + strconv.FormatInt(message.Chat.ID, 10)
	values := map[string]string{
		"author": message.From.UserName,
		"message": message.Text,
		"timestamp": message.Time().Format(time.RFC3339),
	}
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
}

type ChatStatus struct {
	ShouldRespond int
}

func (self *ChatMotionsClient) ChatStatus(chatId int64) *ChatStatus {
	url := self.url + "/status/" + strconv.FormatInt(chatId, 10)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	var objmap ChatStatus
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	return &objmap
}


func (self *ChatMotionsClient) ResetShouldRespond(chatId int64) {
	url := self.url + "/reset_should_respond/" + strconv.FormatInt(chatId, 10)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
