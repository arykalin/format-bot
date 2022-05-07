package data_getter

import "fmt"

type getter struct {
	questionsPath  string
	formatsSheetID string
}

type Getter interface {
	GetData() (formats []Format, questions []Question, err error)
}

func (g *getter) GetData() (formats []Format, questions []Question, err error) {
	jsonDataLoader := newJsonData(g.questionsPath)
	questions, err = jsonDataLoader.getQuestions()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get questions: %w", err)
	}
	return nil, questions, err
}

func NewGetter() *getter {
	return &getter{}
}
