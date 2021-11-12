package telegram

import (
	"fmt"
	"log"

	formatsPkg "github.com/arykalin/format-bot/formats"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type session struct {
	id            int64
	tags          []string
	nextQuestion  int
	waitingAnswer bool
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

func makeAnswerKeyboard(answers []formatsPkg.Answer) tgbotapi.ReplyKeyboardMarkup {
	var a [][]tgbotapi.KeyboardButton
	for _, answer := range answers {
		a = append(a, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(answer.Name)))
	}
	return tgbotapi.NewReplyKeyboard(a...)
}

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
				id:           update.Message.Chat.ID,
				tags:         nil,
				nextQuestion: 0,
			}
		}

		switch update.Message.Text {
		case "status":
			msg.Text = "I'm ok."
			_, err := t.bot.Send(msg)
			if err != nil {
				t.logger.Errorw("error sending message %s", err)
			}
		case "/open":
			msg.ReplyMarkup = numericKeyboard
			_, err := t.bot.Send(msg)
			if err != nil {
				t.logger.Errorw("error sending message %s", err)
			}
		case "/close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err := t.bot.Send(msg)
			if err != nil {
				t.logger.Errorw("error sending message %s", err)
			}
		case "/reload":
			s.nextQuestion = 0
		case "/start":
			t.askQuestion(s, msg)
		default:
			// if waiting for answer make tags
			if s.waitingAnswer {

			}
		}
	}
	return nil
}

func (t teleBot) askQuestion(s session, msg tgbotapi.MessageConfig) {
	var m string
	q := t.formats.GetQuestion(s.nextQuestion)
	if q == nil {
		m = "no questions left"
	} else {
		m = q.Question
	}
	msg.Text = m
	msg.ReplyMarkup = makeAnswerKeyboard(q.Answers)
	s.nextQuestion++
	s.waitingAnswer = true
	_, err := t.bot.Send(msg)
	if err != nil {
		t.logger.Errorw("error sending message %s", err)
	}
	t.sessions[s.id] = s
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
