package data_getter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type jsonData struct{}

type JsonDater interface {
	GetQuestions() (questions []Question, err error)
}

func (j *jsonData) GetQuestions() (questions []Question, err error) {
	err := j.loadJson(questionPAth)
	if err != nil {
		return nil, err
	}
	d, err := j.loadJsonData()
	if err != nil {
		return nil, fmt.Errorf("failed to load formats %w", err)
	}
	err = j.saveJsonData(d, questionPAth+"-backup")
	if err != nil {
		return nil, fmt.Errorf("failed to save formats %w", err)
	}
	return &jsonData{}
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

func NewJsonData() JsonDater {
	return &jsonData{}
}
