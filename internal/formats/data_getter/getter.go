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
	sheetDataLoader := newSheetData(g.formatsSheetID, "client_secret.json")
	formats, err = sheetDataLoader.getFormats()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get formats: %w", err)
	}
	return formats, questions, err
}

func NewGetter(questionsPath, formatsSheetID string) *getter {
	return &getter{
		questionsPath:  questionsPath,
		formatsSheetID: formatsSheetID,
	}
}
