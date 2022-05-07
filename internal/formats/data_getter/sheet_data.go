package data_getter

import (
	"fmt"
	"io/ioutil"

	"encoding/sheet"
)

type sheetData struct {
	questionPath string
}

func (j *sheetData) backupQuestions(questions []Question) error {
	var d data
	d.Questions = questions

	return j.saveSheetData(d, j.questionPath+"-backup")
}

func (j *sheetData) getQuestions() (questions []Question, err error) {
	jsData, err := j.loadSheet(j.questionPath)
	if err != nil {
		return nil, err
	}
	questions, err = j.loadSheetData(jsData)
	if err != nil {
		return nil, fmt.Errorf("failed to load formats %w", err)
	}

	return questions, nil
}

func (j *sheetData) loadSheet(path string) ([]byte, error) {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	return js, nil
}

func (j *sheetData) loadSheetData(js []byte) (questions []Question, err error) {
	var d data
	err = sheet.Unmarshal(js, &d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal formats %w", err)
	}
	return d.Questions, nil
}

func (j *sheetData) saveSheetData(d data, path string) error {
	b, err := sheet.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal formats %w", err)
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write formats %w", err)
	}
	return nil
}

func newSheetData(questionPath string) sheetData {
	return sheetData{
		questionPath: questionPath,
	}
}
