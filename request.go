package main

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func NewRequest(update *tgbotapi.Update) (*Request) {
	self := &Request{
		update: update,
	}
	return self
}

func (self *Request) Answer(message string) {
	msg := tgbotapi.NewMessage(self.update.Message.Chat.ID, message)
	BOT.Send(msg)
}

func (self *Request) AnswerErrorWithLog(message string, err error) {
	log.Print(err)
	self.Answer(message)
}
//
func (self *Request) GetReplyIDIfNeeded() *int {
	if self.update.Message.ReplyToMessage != nil {
		if self.update.Message.ReplyToMessage.From != nil {
			if strings.Compare(self.update.Message.ReplyToMessage.From.UserName, BOT.Self.UserName) == 0 {
				return &self.update.Message.MessageID
			}
		}
	}
	return nil
}

func (self *Request) GetCommand() CommandIF {
	var command CommandIF
	switch self.update.Message.Command() {
	case "start":
		command = &CommandStart{request: self}
	case "stop":
		command = &CommandStop{request: self}
	case "help":
		command = &CommandHelp{request: self}
	case "delay":
		command = &CommandDelay{request: self}
	case "hardcore":
		command = &CommandHardcore{request: self}
	case "gentle":
		command = &CommandGentle{request: self}
	case "amount":
		command = &CommandAmount{request: self}
	default:
		command = &CommandNotFound{request: self}
	}
	return command
}

func (self *Request) Handle() {
	if self.update.Message == nil { // ignore any non-Message updates
		return
	}

	if self.update.Message.IsCommand() {
		command := self.GetCommand()
		_ = command.Handle()
	}
	if self.IsStopped() {
		return
	}
	replyID := self.GetReplyIDIfNeeded()
	if self.IsAnswerNeeded(replyID) {
		if replyID == nil {
			self.CleanDelay()
		}
		output := self.Huify()
		if output != "" {
			msg := tgbotapi.NewMessage(self.update.Message.Chat.ID, output)
			if replyID != nil {
				msg.ReplyToMessageID = *replyID
			}
			BOT.Send(msg)
		}
	}
	self.HandleDelay()
}

func (self *Request) HandleDelay() {
	delay := self.GetDelay()
	log.Printf("Delay is %d", delay)

	if delay > 0 {
		self.SetDelay(delay - 1)
	} else {
		self.SetDelay(rand.Intn(self.GetMaxDelay() + 1))
	}
}

func (self *Request) IsAnswerNeeded(replyID *int) bool {
	if replyID != nil {
		return true
	}
	delay := self.GetDelay()
	return delay == 0
}

// DELAY STORAGE

func (self *Request) GetDelay() int {
	var delay int
	err := DB.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("delay:%d", self.update.Message.Chat.ID)
		item, err := txn.Get([]byte(key))
		if err != nil {
			delay = 0
			return nil
		}
		err = item.Value(func(val []byte) error {
			delay, err = strconv.Atoi(string(val))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return delay
}

func (self *Request) SetDelay(value int) {
	err := DB.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("delay:%d", self.update.Message.Chat.ID)
		err := txn.Set([]byte(key),[]byte(strconv.Itoa(value)))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (self *Request) GetMaxDelay() int {
	var delay int
	err := DB.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("maxdelay:%d", self.update.Message.Chat.ID)
		item, err := txn.Get([]byte(key))
		if err != nil {
			delay = 4
			return nil
		}
		err = item.Value(func(val []byte) error {
			delay, err = strconv.Atoi(string(val))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return delay
}

func (self *Request) SetMaxDelay(value int) error {
	err := DB.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("maxdelay:%d", self.update.Message.Chat.ID)
		err := txn.Set([]byte(key),[]byte(strconv.Itoa(value)))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	return err
}

func (self *Request) CleanDelay() {
	err := DB.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("delay:%d", self.update.Message.Chat.ID)
		err := txn.Delete([]byte(key))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// GENTLE STORAGE

func (self *Request) IsGentle() bool {
	var isGentle bool
	err := DB.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("gentle:%d", self.update.Message.Chat.ID)
		item, err := txn.Get([]byte(key))
		if err != nil {
			isGentle = true
			return nil
		}
		var value int
		err = item.Value(func(val []byte) error {
			value, err = strconv.Atoi(string(val))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		isGentle = value == 1
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return isGentle
}

func (self *Request) SetGentle(gentle bool) error {
	err := DB.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("gentle:%d", self.update.Message.Chat.ID)
		var err error
		if gentle {
			err = txn.Set([]byte(key), []byte("1"))
		} else {
			err = txn.Set([]byte(key), []byte("0"))
		}
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	return err
}


// STOPPED STORAGE

func (self *Request) IsStopped() bool {
	var isStopped bool
	err := DB.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("stopped:%d", self.update.Message.Chat.ID)
		item, err := txn.Get([]byte(key))
		if err != nil {
			isStopped = false
			return nil
		}
		var value int
		err = item.Value(func(val []byte) error {
			value, err = strconv.Atoi(string(val))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		isStopped = value == 1
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return isStopped
}

func (self *Request) SetStopped(stopped bool) error {
	err := DB.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("stopped:%d", self.update.Message.Chat.ID)
		var err error
		if stopped {
			err = txn.Set([]byte(key), []byte("1"))
		} else {
			err = txn.Set([]byte(key), []byte("0"))
		}
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	return err
}

// WORDS AMOUNT STORAGE

func (self *Request) GetWordsAmount() int {
	var amount int
	err := DB.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("amount:%d", self.update.Message.Chat.ID)
		item, err := txn.Get([]byte(key))
		if err != nil {
			amount = 1
			return nil
		}
		err = item.Value(func(val []byte) error {
			amount, err = strconv.Atoi(string(val))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return amount
}

func (self *Request) SetWordsAmount(value int) error {
	err := DB.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("amount:%d", self.update.Message.Chat.ID)
		err := txn.Set([]byte(key),[]byte(strconv.Itoa(value)))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	return err
}

func (self *Request) Huify() string {
	return Huify(self.update.Message.Text, self.IsGentle(), rand.Intn(self.GetWordsAmount()+1))
}
