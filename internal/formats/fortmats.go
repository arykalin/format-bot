package formats

import (
	"fmt"

	dataGetter "github.com/arykalin/format-bot/internal/formats/data_getter"
)

type formats struct {
	formats   []dataGetter.Format
	questions []dataGetter.Question
}

func (f *formats) GetQuestion(num int) *dataGetter.Question {
	if len(f.questions) <= num {
		return nil
	}
	return &f.questions[num]
}

func (f *formats) GetTags(question dataGetter.Question) ([]dataGetter.Tag, error) {
	panic("implement me")
}

func (f *formats) GetFormats(tags []dataGetter.Tag) (formats []dataGetter.Format, err error) {
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
	getter dataGetter.Getter,
) (_ Formats, err error) {
	f := &formats{}
	f.formats, f.questions, err = getter.GetData()
	if err != nil {
		return nil, fmt.Errorf("failed to get formats: %w", err)
	}
	return f, nil
}
