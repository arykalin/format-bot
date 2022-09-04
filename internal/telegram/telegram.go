package telegram

import (
	"fmt"
	"log"
	"strings"

	"github.com/arykalin/format-bot/internal/formats"
	"github.com/arykalin/format-bot/internal/formats/data_getter"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type session struct {
	id            int64
	tags          []data_getter.Tag
	nextQuestion  int
	waitingAnswer bool
}

type sessions map[int64]session

type teleBot struct {
	chatID   int64
	formats  formats.Formats
	bot      *tgbotapi.BotAPI
	sessions sessions
	logger   *zap.SugaredLogger
	getter   data_getter.Getter
}

func makeAnswerKeyboard(answers []data_getter.Answer) tgbotapi.ReplyKeyboardMarkup {
	var a [][]tgbotapi.KeyboardButton
	for _, answer := range answers {
		a = append(a, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(answer.Name)))
	}
	a = append(a, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Пропустить вопрос")))
	return tgbotapi.NewReplyKeyboard(a...)
}

type TeleBot interface {
	Start() error
	UpdateFormats() error
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

		s := strings.ToLower(update.Message.Text)
		s = strings.ReplaceAll(s, "/", "")
		switch s {
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
		case "list":
			t.listFormats(update.Message.Chat.ID)
		default:
			s := t.sessions[update.Message.Chat.ID]
			// if waiting for answer make tags
			if s.waitingAnswer {
				t.handleAnswer(s, update, msg)
				t.askQuestion(update.Message.Chat.ID, msg)
			}
		}
	}
	return nil
}

func (t *teleBot) UpdateFormats() error {
	f, err := formats.NewFormats(t.getter)
	if err != nil {
		return fmt.Errorf("failed to get new formats: %w", err)
	}
	t.formats = f
	return nil
}

// handleAnswer is checking answers and append tags from it
func (t *teleBot) handleAnswer(s session, update tgbotapi.Update, msg tgbotapi.MessageConfig) {
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
}

func (t *teleBot) reload(update tgbotapi.Update) {
	s := t.sessions[update.Message.Chat.ID]
	s.tags = []data_getter.Tag{}
	s.waitingAnswer = false
	s.nextQuestion = 0
	t.sessions[update.Message.Chat.ID] = s
}

func (t *teleBot) askQuestion(id int64, msg tgbotapi.MessageConfig) {
	var err error
	s := t.sessions[id]
	var m string
	//TODO: https://trello.com/c/LZI840VO
	// if only one format left don't ask next question and show formats
	gotFormats, err := t.formats.GetFormats(s.tags)
	if err != nil {
		t.logger.Errorw("error getting formats %s", err)
	}
	if len(gotFormats) == 1 && t.sessions[id].nextQuestion != 0 {
		s.waitingAnswer = false
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		msg.Text = "Остался только один формат для выбранных тегов. Нет смысла дальше задавать вопросы"
		_, err = t.bot.Send(msg)
		if err != nil {
			t.logger.Errorw("error sending message %s", err)
		}
		t.showFormats(s, msg)
		return
	}

	q := t.formats.GetQuestion(s.nextQuestion)
	// if no questions left show formats
	if q == nil {
		//m = "no questions left"
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
	_, err = t.bot.Send(msg)
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
		var tags []string
		for _, tag := range s.tags {
			tags = append(tags, fmt.Sprintf("\"%s\"", tag))
		}
		msg.Text = fmt.Sprintf("Нет подходящих форматов для тегов %s", strings.Join(tags, ", "))
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

func (t *teleBot) makeFormatMsg(format data_getter.Format) string {
	var tags []string
	for _, tag := range format.Tags {
		tags = append(tags, fmt.Sprintf("\"%s\"", tag))
	}
	return fmt.Sprintf("Формат:%s\n Описание: %s\n Теги: %s\n", format.Name, format.Description, strings.Join(tags, ", "))
}

func (t *teleBot) listFormats(id int64) {
	formats, err := t.formats.GetFormats(nil)
	if err != nil {
		t.logger.Errorw("error getting formats %s", err)
		msg := tgbotapi.NewMessage(id, fmt.Sprintf("error getting formats %s", err))
		_, err = t.bot.Send(msg)
		if err != nil {
			t.logger.Errorw("error sending message %s", err)
		}
	}
	for _, format := range formats {
		msg := tgbotapi.NewMessage(id, fmt.Sprintf(
			"Имя формата: %s\nОписание формата: %s\nТеги формата: %s",
			format.Name,
			format.Description,
			format.Tags,
		))
		_, err = t.bot.Send(msg)
		if err != nil {
			t.logger.Errorw("error sending message %s", err)
		}
	}
	return
}

func NewBot(chatID int64, token string, sheetId string, logger *zap.SugaredLogger) TeleBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	getter := data_getter.NewGetter("./internal/formats/data_getter/questions.json", sheetId)
	return &teleBot{
		getter: getter,
		chatID: chatID,
		//formats:  f,
		bot:      bot,
		sessions: make(sessions),
		logger:   logger.Named("teletBot"),
	}
}
