package formats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Formats interface {
	GetTags(Question) ([]Tag, error)
	GetFormats(*Tag) ([]Format, error)
}

type formats struct {
	j         []byte
	formats   []Format
	questions []Question
}

type Tag string

type data struct {
	Formats   []Format   `json:"formats"`
	Questions []Question `json:"questions"`
}

type Format struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        []Tag  `json:"tags"`
}

type Question struct {
	Number   int      `json:"number"`
	Question string   `json:"question"`
	Answers  []Answer `json:"answers"`
}

type Answer struct {
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
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

func NewFormats(path string) (*formats, error) {
	f := &formats{}
	err := f.loadJson(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
