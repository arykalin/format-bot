package data_getter

import (
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"gopkg.in/Iwark/spreadsheet.v2"
)

type sheet struct {
}

type sheetData struct {
	sheetID string
	logger  *zap.SugaredLogger
	service *spreadsheet.Service
}

func (j *sheetData) getSheet(id string) (s spreadsheet.Spreadsheet, err error) {
	data, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return s, fmt.Errorf("failed to read client secret: %w", err)
	}
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		return s, fmt.Errorf("failed to get jwt config: %w", err)
	}
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	s, err = service.FetchSpreadsheet(id)
	if err != nil {
		return s, fmt.Errorf("failed to fetch spreadsheet: %w", err)
	}
	return s, nil
}

func (j *sheetData) getFormats() (formats []Format, err error) {
	jsData, err := j.loadSheet(j.sheetID)
	if err != nil {
		return nil, err
	}
	formats, err = j.loadSheetData(jsData)
	if err != nil {
		return nil, fmt.Errorf("failed to load formats %w", err)
	}

	return formats, nil
}

func (j *sheetData) loadSheet(path string) ([]byte, error) {
	js, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	return js, nil
}

func (j *sheetData) loadSheetData(js []byte) (formats []Format, err error) {
	return nil, err
}

func newSheetData(questionPath string) sheetData {
	return sheetData{
		sheetID: questionPath,
	}
}
