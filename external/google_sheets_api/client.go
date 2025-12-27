package google_sheets_api

import (
	"fmt"
	"google.golang.org/api/sheets/v4"
	"regexp"
	"wb_logistic_assistant/external/google_sheets_api/auth"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

// validateActor checks if the actor is authorized and has the necessary services
func validateActor(actor auth.Actor) (*sheets.SpreadsheetsService, error) {
	if !actor.IsAuth() {
		return nil, fmt.Errorf("actor is not authorized")
	}
	spreadsheets := actor.Spreadsheets()
	if spreadsheets == nil {
		return nil, fmt.Errorf("spreadsheets service not initialized")
	}
	return spreadsheets, nil
}

// GetSheets retrieves all sheets from a spreadsheet
func (c *Client) GetSheets(actor auth.Actor, spreadsheetID string) ([]*sheets.Sheet, error) {
	svc, err := validateActor(actor)
	if err != nil {
		return nil, err
	}
	res, err := svc.Get(spreadsheetID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get sheets: %w", err)
	}
	return res.Sheets, nil
}

// GetSpreadsheetInfo returns spreadsheet metadata
func (c *Client) GetSpreadsheetInfo(actor auth.Actor, spreadsheetID string) (*sheets.Spreadsheet, error) {
	svc, err := validateActor(actor)
	if err != nil {
		return nil, err
	}
	return svc.Get(spreadsheetID).Do()
}

// GetValue retrieves cell data in the specified range
func (c *Client) GetValue(actor auth.Actor, spreadsheetID, sheetName, range_ string) (*sheets.ValueRange, error) {
	svc, err := validateActor(actor)
	if err != nil {
		return nil, err
	}
	if !isValidRange(range_) {
		return nil, fmt.Errorf("invalid range format: %s", range_)
	}
	return svc.Values.Get(spreadsheetID, sheetName+"!"+range_).Do()
}

// UpdateValues modifies cell values
func (c *Client) UpdateValues(actor auth.Actor, spreadsheetID, sheetName, range_, inputOption string, values [][]interface{}) error {
	svc, err := validateActor(actor)
	if err != nil {
		return err
	}
	if !isValidRange(range_) {
		return fmt.Errorf("invalid range format: %s", range_)
	}
	if inputOption == "" {
		inputOption = ValueInputOptionRaw
	}
	_, err = svc.Values.Update(spreadsheetID, sheetName+"!"+range_, &sheets.ValueRange{Values: values}).ValueInputOption(inputOption).Do()
	return err
}

// UpdateValuesExtended is like Update but allows full ValueRange object
func (c *Client) UpdateValuesExtended(actor auth.Actor, spreadsheetID, sheetName, range_, inputOption string, valueRange *sheets.ValueRange) error {
	svc, err := validateActor(actor)
	if err != nil {
		return err
	}
	if !isValidRange(range_) {
		return fmt.Errorf("invalid range format: %s", range_)
	}
	if inputOption == "" {
		inputOption = ValueInputOptionRaw
	}
	_, err = svc.Values.Update(spreadsheetID, sheetName+"!"+range_, valueRange).ValueInputOption(inputOption).Do()
	return err
}

// AppendValues adds new rows or columns
func (c *Client) AppendValues(actor auth.Actor, spreadsheetID, sheetName, range_, inputOption, insertDataOption string, values [][]interface{}) error {
	svc, err := validateActor(actor)
	if err != nil {
		return err
	}
	if !isValidRange(range_) {
		return fmt.Errorf("invalid range format: %s", range_)
	}
	if inputOption == "" {
		inputOption = ValueInputOptionRaw
	}
	if insertDataOption == "" {
		insertDataOption = InsertDataOptionInsertRows
	}
	_, err = svc.Values.Append(spreadsheetID, sheetName+"!"+range_, &sheets.ValueRange{Values: values}).ValueInputOption(inputOption).InsertDataOption(insertDataOption).Do()
	return err
}

// AppendValuesExtended adds rows/columns with full ValueRange
func (c *Client) AppendValuesExtended(actor auth.Actor, spreadsheetID, sheetName, range_, inputOption, insertDataOption string, valueRange *sheets.ValueRange) error {
	svc, err := validateActor(actor)
	if err != nil {
		return err
	}
	if !isValidRange(range_) {
		return fmt.Errorf("invalid range format: %s", range_)
	}
	if inputOption == "" {
		inputOption = ValueInputOptionRaw
	}
	if insertDataOption == "" {
		insertDataOption = InsertDataOptionInsertRows
	}
	_, err = svc.Values.Append(spreadsheetID, sheetName+"!"+range_, valueRange).ValueInputOption(inputOption).InsertDataOption(insertDataOption).Do()
	return err
}

// BatchUpdate performs multiple sheet operations at once
func (c *Client) BatchUpdate(actor auth.Actor, spreadsheetID string, batchUpdateRequest *sheets.BatchUpdateSpreadsheetRequest) error {
	svc, err := validateActor(actor)
	if err != nil {
		return err
	}
	_, err = svc.BatchUpdate(spreadsheetID, batchUpdateRequest).Do()
	return err
}

// ClearValues removes cell contents but keeps formatting
func (c *Client) ClearValues(actor auth.Actor, spreadsheetID, sheetName, range_ string) error {
	svc, err := validateActor(actor)
	if err != nil {
		return err
	}
	if !isValidRange(range_) {
		return fmt.Errorf("invalid range format: %s", range_)
	}
	r := sheetName
	if range_ != "" {
		r = r + "!" + range_
	}
	_, err = svc.Values.Clear(spreadsheetID, r, &sheets.ClearValuesRequest{}).Do()
	return err
}

// ClearFormat resets style formatting in the specified range
func (c *Client) ClearFormat(actor auth.Actor, spreadsheetID string, sheetID int64, startRow, startColumn, endRow, endColumn int64) error {
	requests := []*sheets.Request{
		{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					SheetId: sheetID,
					GridProperties: &sheets.GridProperties{
						FrozenRowCount:    0,
						FrozenColumnCount: 0,
					},
				},
				Fields: "gridProperties(frozenRowCount, frozenColumnCount)",
			},
		},
		{
			UnmergeCells: &sheets.UnmergeCellsRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartRowIndex:    startRow,
					EndRowIndex:      endRow,
					StartColumnIndex: startColumn,
					EndColumnIndex:   endColumn,
				},
			},
		},
		{
			UpdateCells: &sheets.UpdateCellsRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartRowIndex:    startRow,
					EndRowIndex:      endRow,
					StartColumnIndex: startColumn,
					EndColumnIndex:   endColumn,
				},
				Fields: "userEnteredFormat",
			},
		},
	}
	return c.BatchUpdate(actor, spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{Requests: requests})
}

// ClearFilters removes all active filters on the sheet
func (c *Client) ClearFilters(actor auth.Actor, spreadsheetID string, sheetID int64) error {
	request := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				ClearBasicFilter: &sheets.ClearBasicFilterRequest{
					SheetId: sheetID,
				},
			},
		},
	}
	return c.BatchUpdate(actor, spreadsheetID, request)
}

// isValidRange validates range in A1 notation
func isValidRange(range_ string) bool {
	matched, err := regexp.MatchString(`^([A-Za-z]+[0-9]*)(:[A-Za-z]+[0-9]*)?$`, range_)
	return err == nil && matched
}
