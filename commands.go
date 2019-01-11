package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CommandStart Command

func (self *CommandStart) Handle() error {
	message := "Привет! Я бот-хуебот.\n" +
		"Я буду хуифицировать некоторые из ваших фраз.\n" +
		"Сейчас режим вежливости %s\n" +
		"За подробностями в /help"
	err := self.request.SetStopped(false)
	if err != nil {
		self.request.AnswerErrorWithLog("Не удалось сохранить", err)
		return nil
	}
	if self.request.IsGentle() {
		message = fmt.Sprintf(message, "включен")
	} else {
		message = fmt.Sprintf(message, "отключен")
	}
	self.request.Answer(message)
	return nil
}

type CommandStop Command

func (self *CommandStop) Handle() error {
	err := self.request.SetStopped(true)
	if err != nil {
		self.request.AnswerErrorWithLog("Не удалось сохранить", err)
		return nil
	}
	self.request.Answer("Выключаюсь")
	return nil
}

type CommandHelp Command

func (self *CommandHelp) Handle() error {
	self.request.Answer(
		"Вежливый режим:\n" +
			"  Для включения используйте команду /gentle\n" +
			"  Для отключения - /hardcore\n" +
			"Частота ответов: /delay N, где N - любое любое натуральное число\n" +
			"Число хуифицируемых слов: /amount N, где N - от 1 до 10\n" +
			"Для остановки используйте /stop")
	return nil
}

type CommandDelay Command

func (self *CommandDelay) Handle() error {
	command := strings.Fields(self.request.update.Message.Text)
	if len(command) < 2 {
		currentDelayMessage := "Сейчас я пропускаю случайное число сообщений от 0 до "
		currentDelayMessage += strconv.Itoa(self.request.GetMaxDelay())
		self.request.Answer(currentDelayMessage)
		return nil
	}
	commandArg := command[len(command)-1]
	tryDelay, err := strconv.Atoi(commandArg)
	if err != nil || tryDelay < 1 || tryDelay > 1000000 {
		self.request.Answer("Неправильный аргумент, отправьте `/delay N`, где N любое натуральное число меньше 1000000")
		return nil
	}
	err = self.request.SetMaxDelay(tryDelay)
	if err != nil {
		self.request.AnswerErrorWithLog("Не удалось сохранить, отправьте еще раз `/delay N`, где N любое натуральное число меньше 1000000", err)
		return nil
	}
	self.request.Answer("Я буду пропускать случайное число сообщений от 0 до " + commandArg)
	self.request.CleanDelay()
	return nil
}

type CommandHardcore Command

func (self *CommandHardcore) Handle() error {
	err := self.request.SetGentle(false)
	if err != nil {
		self.request.AnswerErrorWithLog("Не удалось сохранить", err)
		return nil
	}
	self.request.Answer("Вежливый режим отключен.\nЧтобы включить его, используйте команду /gentle")
	return nil
}

type CommandGentle Command

func (self *CommandGentle) Handle() error {
	err := self.request.SetGentle(true)
	if err != nil {
		self.request.AnswerErrorWithLog("Не удалось сохранить", err)
		return nil
	}
	self.request.Answer("Вежливый режим включен.\nЧтобы отключить его, используйте команду /hardcore")
	return nil
}

type CommandAmount Command

func (self *CommandAmount) Handle() error {
	command := strings.Fields(self.request.update.Message.Text)
	if len(command) < 2 {
		currentWordsAmount := self.request.GetWordsAmount()
		self.request.Answer("Сейчас я хуифицирую случайное число слов от 1 до " + strconv.Itoa(currentWordsAmount))
		return nil
	}
	commandArg := command[len(command)-1]
	tryWordsAmount, err := strconv.Atoi(commandArg)
	if err != nil || tryWordsAmount < 1 || tryWordsAmount > 10 {
		self.request.Answer("Неправильный аргумент, отправьте `/amount N`, где N любое натуральное число не больше 10")
		return nil
	}
	err = self.request.SetWordsAmount(tryWordsAmount)
	if err != nil {
		self.request.AnswerErrorWithLog("Не удалось сохранить, отправьте еще раз `/amount N`, где N любое натуральное число не больше 10", err)
		return nil
	}
	self.request.Answer("Я буду хуифицировать случайное число слов от 1 до " + commandArg)
	return nil
}

type CommandNotFound Command

func (self *CommandNotFound) Handle() error {
	return errors.New("Команда не найдена")
}
