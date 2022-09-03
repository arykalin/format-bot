package data_getter

import (
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"gopkg.in/Iwark/spreadsheet.v2"
)

type sheetData struct {
	sheetID    string
	secretPath string
	config     sheetConfig
}

type sheetConfig struct {
	creatorEmailIdx int
	nameIdx         int
	descriptionIdx  int
	tagsStartIdx    int
	tagsEndIdx      int
	skip            int
}

// getSheet returns the google sheet with formats.
func (s *sheetData) getSheet() (sheet *spreadsheet.Sheet, err error) {
	secretData, err := ioutil.ReadFile(s.secretPath)
	if err != nil {
		return sheet, fmt.Errorf("failed to read client secret: %w", err)
	}
	conf, err := google.JWTConfigFromJSON(secretData, spreadsheet.Scope)
	if err != nil {
		return sheet, fmt.Errorf("failed to get jwt config: %w", err)
	}
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	fetchSheet, err := service.FetchSpreadsheet(s.sheetID)
	if err != nil {
		return sheet, fmt.Errorf("failed to fetch spreadsheet: %w", err)
	}
	sheet, err = fetchSheet.SheetByIndex(0)
	if err != nil {
		return sheet, fmt.Errorf("failed to fetch sheet: %w", err)
	}
	return sheet, nil
}

// getFormats returns the formats from the sheet.
func (s *sheetData) getFormats() (formats []Format, err error) {
	sheet, err := s.getSheet()
	if err != nil {
		return nil, fmt.Errorf("failed to load sheet: %s", err)
	}
	for i := range sheet.Rows {
		if i < s.config.skip {
			// skip
			continue
		}
		format := Format{}

		var creatorEmail string
		if len(sheet.Rows[i]) > s.config.creatorEmailIdx {
			creatorEmail = sheet.Rows[i][s.config.creatorEmailIdx].Value
		}
		format.CreatorEmail = creatorEmail

		var name string
		if len(sheet.Rows[i]) > s.config.nameIdx {
			name = sheet.Rows[i][s.config.nameIdx].Value
		}
		format.Name = name

		var description string
		if len(sheet.Rows[i]) > s.config.descriptionIdx {
			description = sheet.Rows[i][s.config.descriptionIdx].Value
		}
		format.Description = description

		var tags []string
		for tagNum := s.config.tagsStartIdx; tagNum <= s.config.tagsEndIdx; tagNum++ {
			if len(sheet.Rows[i]) > tagNum {
				tags = append(tags, sheet.Rows[i][tagNum].Value)
			}
		}
		format.Tags = s.getTags(tags)
		formats = append(formats, format)
	}
	return formats, nil
}

func (s *sheetData) getTags(tagGroups []string) (tagsList []Tag) {
	for _, tagGroup := range tagGroups {
		tags := strings.Split(tagGroup, ",")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			tag = strings.ReplaceAll(tag, "\"", "")
			tagsList = append(tagsList, Tag(tag))
		}
	}
	return tagsList
}

func newSheetData(sheetID, secretPath string) sheetData {
	return sheetData{
		sheetID:    sheetID,
		secretPath: secretPath,
		config: sheetConfig{
			creatorEmailIdx: 1,
			nameIdx:         2,
			descriptionIdx:  3,
			tagsStartIdx:    4,
			tagsEndIdx:      10,
			skip:            1,
		},
	}
}
