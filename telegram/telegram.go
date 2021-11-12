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
	tags          []formatsPkg.Tag
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

func (t *teleBot) Start() error {
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

		_, ok := t.sessions[update.Message.Chat.ID]
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
		case "reload":
			t.reload(update)
		case "start":
			t.reload(update)
			t.askQuestion(update.Message.Chat.ID, msg)
		case "tags":
			s := t.sessions[update.Message.Chat.ID]
			msg.Text = fmt.Sprintf("%v", s.tags)
			_, err := t.bot.Send(msg)
			if err != nil {
				t.logger.Errorw("error sending message %s", err)
			}
		default:
			s := t.sessions[update.Message.Chat.ID]
			// if waiting for answer make tags
			if s.waitingAnswer {
				q := t.formats.GetQuestion(s.nextQuestion - 1)
				for _, answer := range q.Answers {
					if answer.Name == update.Message.Text {
						s.tags = append(s.tags, answer.Tags...)
					}
				}
				t.sessions[s.id] = s
				s.waitingAnswer = false
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				msg.Text = "Спасибо, следующий вопрос"
				_, err := t.bot.Send(msg)
				if err != nil {
					t.logger.Errorw("error sending message %s", err)
				}
				t.askQuestion(update.Message.Chat.ID, msg)
			}
		}
	}
	return nil
}

func (t *teleBot) reload(update tgbotapi.Update) {
	s := t.sessions[update.Message.Chat.ID]
	s.waitingAnswer = false
	s.nextQuestion = 0
}

func (t *teleBot) askQuestion(id int64, msg tgbotapi.MessageConfig) {
	s := t.sessions[id]
	var m string
	q := t.formats.GetQuestion(s.nextQuestion)
	if q == nil {
		m = "no questions left"
		s.waitingAnswer = false
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		t.showFormats(s, msg)
		return
	}
	m = q.Question
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

func (t *teleBot) showFormats(s session, msg tgbotapi.MessageConfig) {
	msg.Text = "Вопросы закончились, показываю подходящие форматы"
	_, err := t.bot.Send(msg)
	if err != nil {
		t.logger.Errorw("error sending message %s", err)
	}
	f, err := t.formats.GetFormats(s.tags)
	if err != nil {
		t.logger.Errorw("error getting formats %s", err)
		msg.Text = fmt.Sprintf("error getting formats %s", err)
		_, err = t.bot.Send(msg)
		if err != nil {
			t.logger.Errorw("error sending message %s", err)
		}
		return
	}
	if len(f) == 0 {
		msg.Text = fmt.Sprintf("Нет подходящих форматов для тегов %v", s.tags)
		_, err = t.bot.Send(msg)
		if err != nil {
			t.logger.Errorw("error sending message %s", err)
		}

	}
	for _, format := range f {
		msg.Text = t.makeFormatMsg(format)
		_, err = t.bot.Send(msg)
		if err != nil {
			t.logger.Errorw("error sending message %s", err)
		}
	}
}

func (t *teleBot) makeFormatMsg(format formatsPkg.Format) string {
	return fmt.Sprintf("Формат:%s\n Описание: %s\n Теги: %s\n", format.Name, format.Description, format.Tags)
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
