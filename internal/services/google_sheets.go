package services

import (
	"google.golang.org/api/sheets/v4"
	"wb_logistic_assistant/external/google_sheets_api"
	"wb_logistic_assistant/external/google_sheets_api/auth"
	"wb_logistic_assistant/internal/errors"
)

type GoogleSheetsService interface {
	GetSheets(id string) ([]*sheets.Sheet, error)
	GetSheetIDByName(id, name string) (int64, error)
	GetSheetNameByID(id string, sheetID int64) (string, error)
	UpdateValues(id, name, cellRange string, values [][]interface{}, isRawInput bool) error
	AppendValues(id, name, cellRange string, values [][]interface{}, isRawInput, isInsertOverwrite bool) error
	ClearValues(id, sheetID, cellRange string) error
	ClearFormat(id string, sheetID int64, row1, column1, row2, column2 int64) error
	ClearFilters(id string, sheetID int64) error
}

type BaseGoogleSheetsService struct {
	client *google_sheets_api.Client
	actor  auth.Actor
}

func NewBaseGoogleSheetsService(client *google_sheets_api.Client, actor auth.Actor) *BaseGoogleSheetsService {
	return &BaseGoogleSheetsService{
		client: client,
		actor:  actor,
	}
}

func (s *BaseGoogleSheetsService) GetSheets(id string) ([]*sheets.Sheet, error) {
	res, err := s.client.GetSheets(s.actor, id)
	if err != nil {
		return nil, errors.Wrap(err, "BaseGoogleSheetsService.GetSheets()", "")
	}
	return res, nil
}

func (s *BaseGoogleSheetsService) GetSheetIDByName(id, name string) (int64, error) {
	sheets, err := s.GetSheets(id)
	if err != nil {
		return 0, err
	}

	for _, sheet := range sheets {
		if sheet.Properties.Title == name {
			return sheet.Properties.SheetId, nil
		}
	}

	return 0, errors.New("BaseGoogleSheetsService.GetSheetIDByName()", "sheet not found")
}

func (s *BaseGoogleSheetsService) GetSheetNameByID(id string, sheetID int64) (string, error) {
	sheets, err := s.GetSheets(id)
	if err != nil {
		return "", err
	}

	for _, sheet := range sheets {
		if sheet.Properties.SheetId == sheetID {
			return sheet.Properties.Title, nil
		}
	}

	return "", errors.New("BaseGoogleSheetsService.GetSheetNameByID()", "sheet not found")
}

func (s *BaseGoogleSheetsService) UpdateValues(id, name, cellRange string, values [][]interface{}, isRawInput bool) error {
	inputOption := google_sheets_api.ValueInputOptionUserEntered
	if isRawInput {
		inputOption = google_sheets_api.ValueInputOptionRaw
	}

	err := s.client.UpdateValues(s.actor, id, name, cellRange, inputOption, values)
	if err != nil {
		return errors.Wrapf(err, "BaseGoogleSheetsService.UpdateValues()", "error to update text %s in sheet on range [%s]", inputOption, cellRange)
	}

	return nil
}

func (s *BaseGoogleSheetsService) AppendValues(id, name, cellRange string, values [][]interface{}, isRawInput, isInsertOverwrite bool) error {
	inputOption := google_sheets_api.ValueInputOptionUserEntered
	if isRawInput {
		inputOption = google_sheets_api.ValueInputOptionRaw
	}

	insertDataOption := google_sheets_api.InsertDataOptionInsertRows
	if isInsertOverwrite {
		insertDataOption = google_sheets_api.InsertDataOptionOverwrite
	}

	err := s.client.AppendValues(s.actor, id, name, cellRange, inputOption, insertDataOption, values)
	if err != nil {
		return errors.Wrapf(err, "BaseGoogleSheetsService.AppendValues()", "Error to append %s text in sheet on range [%s] with insert option %s", inputOption, cellRange, insertDataOption)
	}
	return nil
}

// ClearValues Clears the values of cells in the specified range
//
//	Range format: A1:D1 or A:B, etc.
//	If need clearing all values cellRange: A:Z
func (s *BaseGoogleSheetsService) ClearValues(id, sheetName, cellRange string) error {
	err := s.client.ClearValues(s.actor, id, sheetName, cellRange)
	if err != nil {
		return errors.Wrapf(err, "BaseGoogleSheetsService.ClearValues()", "failed to clear values sheet to range [%s]", cellRange)
	}
	return nil
}

// ClearFormat Clears the formatting of table styles in the specified range
//
//	If need clear all row and column: 0
func (s *BaseGoogleSheetsService) ClearFormat(id string, sheetID int64, row1, column1, row2, column2 int64) error {
	err := s.client.ClearFormat(s.actor, id, sheetID, row1, column1, row2, column2)
	if err != nil {
		return errors.Wrap(err, "BaseGoogleSheetsService.ClearFormat()", "")
	}
	return nil
}

// ClearFilters Clears the formatting of all table styles
func (s *BaseGoogleSheetsService) ClearFilters(id string, sheetID int64) error {
	err := s.client.ClearFilters(s.actor, id, sheetID)
	if err != nil {
		return errors.Wrap(err, "BaseGoogleSheetsService.ClearFilters()", "")
	}
	return nil
}
