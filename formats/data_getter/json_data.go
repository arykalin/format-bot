package data_getter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type jsonData struct {
	questionPath string
}

func (j *jsonData) backupQuestions(questions []Question) error {
	var d data
	d.Questions = questions

	return j.saveJsonData(d, j.questionPath+"-backup")
}

func (j *jsonData) getQuestions() (questions []Question, err error) {
	jsData, err := j.loadJson(j.questionPath)
	if err != nil {
		return nil, err
	}
	questions, err = j.loadJsonData(jsData)
	if err != nil {
		return nil, fmt.Errorf("failed to load formats %w", err)
	}

	return questions, nil
}

func (j *jsonData) loadJson(path string) ([]byte, error) {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	return js, nil
}

func (j *jsonData) loadJsonData(js []byte) (questions []Question, err error) {
	var d data
	err = json.Unmarshal(js, &d)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal formats %w", err)
	}
	return d.Questions, nil
}

func (j *jsonData) saveJsonData(d data, path string) error {
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal formats %w", err)
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write formats %w", err)
	}
	return nil
}

func newJsonData(questionPath string) jsonData {
	return jsonData{
		questionPath: questionPath,
	}
}
