package telegram

import (
	"fmt"
	"log"

	formatsPkg "github.com/arykalin/format-bot/formats"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type session struct {
	tags         []string
	nextQuestion int
}

type sessions map[int64]session

type teleBot struct {
	chatID   int64
	formats  formatsPkg.Formats
	bot      *tgbotapi.BotAPI
	sessions sessions
	logger   *zap.SugaredLogger
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

type TeleBot interface {
	Start() error
}

func (t teleBot) Start() error {
	t.bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("failed to get updates %w", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "непонятная команда")

		s, ok := t.sessions[update.Message.Chat.ID]
		if !ok {
			t.sessions[update.Message.Chat.ID] = session{
				tags:         nil,
				nextQuestion: 0,
			}
		}

		switch update.Message.Text {
		case "status":
			msg.Text = "I'm ok."
		case "/open":
			msg.ReplyMarkup = numericKeyboard
		case "/close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "/start":
			var m string
			q := t.formats.GetQuestion(t.sessions[update.Message.Chat.ID].nextQuestion)
			if q == nil {
				m = "no questions left"
			} else {
				m = q.Question
			}
			msg.Text = m
			s.nextQuestion = s.nextQuestion + 1
			t.sessions[update.Message.Chat.ID] = s
		}

		_, err = t.bot.Send(msg)
		if err != nil {
			t.logger.Errorw("error sending message %s", err)
		}
	}
	return nil
}

func NewBot(
	chatID int64,
	token string,
	logger *zap.SugaredLogger,
) TeleBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	f, err := formatsPkg.NewFormats("./formats/formats.json")
	if err != nil {
		log.Panic(err)
	}
	return &teleBot{
		chatID:   chatID,
		formats:  f,
		bot:      bot,
		sessions: make(sessions),
		logger:   logger.Named("teletBot"),
	}
}
