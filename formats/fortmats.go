package formats

import (
	"github.com/arykalin/format-bot/formats/data_getter"
)

type formats struct {
	formats   []data_getter.Format
	questions []data_getter.Question
}

func (f *formats) GetQuestion(num int) *data_getter.Question {
	if len(f.questions) <= num {
		return nil
	}
	return &f.questions[num]
}

func (f *formats) GetTags(question data_getter.Question) ([]data_getter.Tag, error) {
	panic("implement me")
}

func (f *formats) GetFormats(tags []data_getter.Tag) (formats []data_getter.Format, err error) {
	if tags == nil {
		return f.formats, nil
	}
	if len(tags) == 0 {
		return f.formats, nil
	}
	for _, format := range f.formats {
		if subset(tags, format.Tags) {
			formats = append(formats, format)
		}
	}
	return formats, nil
}

//"./formats/formats.json"
func NewFormats(
	getter data_getter.Getter,
) (_ Formats, err error) {
	f := &formats{}
	f.formats, f.questions, err = getter.GetData()
	return f, nil
}
