package formats

import (
	data_getter2 "github.com/arykalin/format-bot/internal/formats/data_getter"
)

type formats struct {
	formats   []data_getter2.Format
	questions []data_getter2.Question
}

func (f *formats) GetQuestion(num int) *data_getter2.Question {
	if len(f.questions) <= num {
		return nil
	}
	return &f.questions[num]
}

func (f *formats) GetTags(question data_getter2.Question) ([]data_getter2.Tag, error) {
	panic("implement me")
}

func (f *formats) GetFormats(tags []data_getter2.Tag) (formats []data_getter2.Format, err error) {
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
	getter data_getter2.Getter,
) (_ Formats, err error) {
	f := &formats{}
	f.formats, f.questions, err = getter.GetData()
	return f, nil
}
