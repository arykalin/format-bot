package formats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Formats interface {
	GetTags(Question) ([]Tag, error)
	GetFormats([]Tag) ([]Format, error)
	GetQuestion(num int) (question *Question)
}

type formats struct {
	j         []byte
	formats   []Format
	questions []Question
}

func (f *formats) GetQuestion(num int) *Question {
	if len(f.questions) <= num {
		return nil
	}
	return &f.questions[num]
}

func (f *formats) GetTags(question Question) ([]Tag, error) {
	panic("implement me")
}

func (f *formats) GetFormats(tags []Tag) (formats []Format, err error) {
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

func (f *formats) loadJson(path string) error {
	j, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}
	f.j = j
	return nil
}
func (f *formats) loadData() (err error) {
	d := data{}
	err = json.Unmarshal(f.j, &d)
	if err != nil {
		return fmt.Errorf("failed to unmarshal formats %w", err)
	}
	f.formats = d.Formats
	f.questions = d.Questions
	return nil
}

func NewFormats() (Formats, error) {
	f := &formats{}
	err := f.loadJson("./formats/formats.json")
	if err != nil {
		return nil, err
	}
	err = f.loadData()
	if err != nil {
		return nil, err
	}
	return f, nil
}
