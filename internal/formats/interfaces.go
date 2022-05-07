package formats

import (
	"github.com/arykalin/format-bot/internal/formats/data_getter"
)

type Formats interface {
	GetFormats([]data_getter.Tag) ([]data_getter.Format, error)
	GetQuestion(num int) (question *data_getter.Question)
}
