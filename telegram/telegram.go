package telegram

import (
	"fmt"
	"log"

	formatsPkg "github.com/arykalin/format-bot/formats"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type session struct {
	tags []string
}

type sessions map[string]session
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "help":
			msg.Text = "type /sayhi or /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		case "/open":
			msg.ReplyMarkup = numericKeyboard
		case "/close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		t.bot.Send(msg)
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
