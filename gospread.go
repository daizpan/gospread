package gospread

import (
	"context"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/mattn/natural"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoSpread struct {
	*sheets.Service
}

// NewGoSpreadWithCredentialFile is credential fileを引数にGoSpreadを生成
func NewGoSpreadWithCredentialFile(credential string) (*GoSpread, error) {
	b, err := ioutil.ReadFile(credential)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return NewGoSpreadWithCredentialsJSON(ctx, b)
}

func NewGoSpreadWithCredentialsJSON(ctx context.Context, jsonData []byte) (*GoSpread, error) {
	cred, err := google.CredentialsFromJSON(ctx, jsonData, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	srv, err := sheets.NewService(ctx, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}

	return &GoSpread{srv}, nil
}

func (g *GoSpread) CreateSheet(spreadsheetId string, sheetName string) error {
	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: sheetName,
					},
				},
			},
		},
	}
	_, err := g.Service.Spreadsheets.BatchUpdate(spreadsheetId, req).Do()
	if err != nil {
		return err
	}

	return nil
}

func (g *GoSpread) DeleteSheet(spreadsheetId string, sheetId int64) error {
	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteSheet: &sheets.DeleteSheetRequest{
					SheetId: sheetId,
				},
			},
		},
	}
	_, err := g.Service.Spreadsheets.BatchUpdate(spreadsheetId, req).Do()
	if err != nil {
		return err
	}

	return nil
}

func (g *GoSpread) DeleteSheetByName(spreadsheetId string, sheetName string) error {
	spreadSheet, err := g.GetSpreadSheet(spreadsheetId)
	if err != nil {
		return err
	}
	sheetId, err := g.sheetId(spreadSheet, sheetName)
	if err != nil {
		return err
	}
	if err := g.DeleteSheet(spreadsheetId, sheetId); err != nil {
		return err
	}
	return nil
}

func (g *GoSpread) DeleteAllSheet(spreadsheetId string, exclusionSheet []string) error {
	spreadSheet, err := g.GetSpreadSheet(spreadsheetId)
	if err != nil {
		return err
	}

	isSheetExists := func(e []string, s string) bool {
		for _, es := range e {
			if es == s {
				return true
			}
		}
		return false
	}

	deleteSheets := []int64{}
	for _, sheet := range spreadSheet.Sheets {
		if !isSheetExists(exclusionSheet, sheet.Properties.Title) {
			deleteSheets = append(deleteSheets, sheet.Properties.SheetId)
		}
	}

	sheetsRequests := []*sheets.Request{}
	for _, sheetId := range deleteSheets {
		r := &sheets.Request{
			DeleteSheet: &sheets.DeleteSheetRequest{
				SheetId: sheetId,
			},
		}
		sheetsRequests = append(sheetsRequests, r)
	}

	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: sheetsRequests,
	}
	_, err = g.BatchUpdate(spreadsheetId, req)
	if err != nil {
		return err
	}
	return nil
}

func (g *GoSpread) GetSpreadSheet(spreadsheetId string) (*sheets.Spreadsheet, error) {
	return g.Spreadsheets.Get(spreadsheetId).Do()
}

// SheetExists: Check sheet exists.
func (g *GoSpread) SheetExists(spreadsheetId string, sheetName string) bool {
	spreadSheet, err := g.Service.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		return false
	}

	if _, err := g.sheetId(spreadSheet, sheetName); err != nil {
		return false
	}
	return true
}

func (g *GoSpread) SheetId(spreadsheetId string, sheetName string) (int64, error) {
	spread, err := g.GetSpreadSheet(spreadsheetId)
	if err != nil {
		return 0, err
	}
	for _, sheet := range spread.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId, nil
		}
	}
	return 0, fmt.Errorf("%s is not exists.", sheetName)
}

func (g *GoSpread) sheetId(spreadSheet *sheets.Spreadsheet, sheetName string) (int64, error) {
	for _, sheet := range spreadSheet.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId, nil
		}
	}
	return 0, fmt.Errorf("%s is not exists.", sheetName)
}

// GetValues: Get sheet data
func (g *GoSpread) GetValues(spreadSheetId string, sheetName string, _range string) (*sheets.ValueRange, error) {
	readRange := sheetName + "!" + _range
	return g.Service.Spreadsheets.Values.Get(spreadSheetId, readRange).Do()
}

func (g *GoSpread) AppendRows(spreadsheetId string, sheetName string, values [][]interface{}) error {
	value := &sheets.ValueRange{
		Values: values,
	}
	_, err := g.Service.Spreadsheets.Values.Append(spreadsheetId, sheetName+"!A1", value).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}
	return nil
}

func (g *GoSpread) UpdateValues(spreadsheetId string, updateRange string, value *sheets.ValueRange) (*sheets.UpdateValuesResponse, error) {
	return g.Spreadsheets.Values.Update(spreadsheetId, updateRange, value).ValueInputOption("USER_ENTERED").Do()
}

func (g *GoSpread) BatchUpdate(spreadsheetId string, req *sheets.BatchUpdateSpreadsheetRequest) (*sheets.BatchUpdateSpreadsheetResponse, error) {
	return g.Spreadsheets.BatchUpdate(spreadsheetId, req).Do()
}

func (g *GoSpread) SortSheets(spreadsheetId string) error {
	spread, err := g.GetSpreadSheet(spreadsheetId)
	if err != nil {
		return err
	}

	sort.Slice(spread.Sheets, func(i, j int) bool {
		return natural.NaturalComp(spread.Sheets[i].Properties.Title, spread.Sheets[j].Properties.Title) < 0
	})

	req := []*sheets.Request{}
	for i, s := range spread.Sheets {
		r := &sheets.Request{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					SheetId: s.Properties.SheetId,
					Index:   int64(i),
				},
				Fields: "SheetId,Index",
			},
		}
		req = append(req, r)
	}

	br := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: req,
	}
	_, err = g.BatchUpdate(spreadsheetId, br)
	if err != nil {
		return err
	}
	return nil
}

func (g *GoSpread) MoveSheet(sheetId int64, index int64, spreadsheetId string) error {
	br := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
					Properties: &sheets.SheetProperties{
						SheetId: sheetId,
						Index:   index,
					},
					Fields: "SheetId,Index",
				},
			},
		},
	}
	_, err := g.BatchUpdate(spreadsheetId, br)
	if err != nil {
		return err
	}
	return nil
}
