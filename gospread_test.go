package gospread

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

var TEST_SHEET_API_CREDENTIAL_JSON = os.Getenv("SHEET_API_CREDENTIAL_JSON")
var TEST_SPREAD_SHEET_ID = os.Getenv("TEST_SPREAD_SHEET_ID")
var TEST_SHEET_NAME = "Test-sheet"
var TEST_STATIC_SHEET = "STATIC"

func newTestGoSpread(t *testing.T) *GoSpread {
	g, err := NewGoSpreadWithCredentialFile(TEST_SHEET_API_CREDENTIAL_JSON)
	if err != nil {
		t.Fatalf("create gospread error")
	}
	return g
}

func setupTestSheet(t *testing.T, g *GoSpread) func() {
	t.Helper()
	createTestSheet(t, g)
	return func() {
		dropTestSheet(t, g)
	}
}

func createTestSheet(t *testing.T, g *GoSpread) {
	t.Helper()
	if err := g.CreateSheet(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME); err != nil {
		t.Errorf("createTestSheet error=%s", err)
	}
}

func dropTestSheet(t *testing.T, g *GoSpread) {
	t.Helper()
	if err := g.DeleteSheetByName(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME); err != nil {
		t.Errorf("dropTestSheet error=%s", err)
	}
}

func TestNewGoSpreadWithCredentialFile(t *testing.T) {
	_, err := NewGoSpreadWithCredentialFile(TEST_SHEET_API_CREDENTIAL_JSON)
	if err != nil {
		t.Error(err)
	}
}

func TestNewGoSpreadWithCredentialFileFail(t *testing.T) {
	credential := "fail-credential.json"
	_, err := NewGoSpreadWithCredentialFile(credential)
	if err == nil {
		t.Errorf("err should be no such file or directory")
	}
}

func TestGoSpread_CreateSheet(t *testing.T) {
	g := newTestGoSpread(t)
	err := g.CreateSheet(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME)
	if err != nil {
		t.Error(err)
	}
	dropTestSheet(t, g)
}

func TestGoSpread_DeleteSheetByName(t *testing.T) {
	g := newTestGoSpread(t)
	createTestSheet(t, g)
	err := g.DeleteSheetByName(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME)
	if err != nil {
		t.Error(err)
	}
}

func TestGoSpread_DeleteAllSheet(t *testing.T) {
	g := newTestGoSpread(t)

	var err error
	createTestSheet(t, g)
	sheet := time.Now().Format("2006-01-02 15:04:05")
	err = g.CreateSheet(TEST_SPREAD_SHEET_ID, sheet)
	if err != nil {
		t.Fatal(err)
	}

	err = g.DeleteAllSheet(TEST_SPREAD_SHEET_ID, []string{sheet, TEST_STATIC_SHEET})
	if err != nil {
		t.Error(err)
	}
}

func TestGoSpread_GetSpreadSheet(t *testing.T) {
	g := newTestGoSpread(t)
	teardown := setupTestSheet(t, g)
	defer teardown()

	got, err := g.GetSpreadSheet(TEST_SPREAD_SHEET_ID)
	if err != nil {
		t.Error(err)
	}
	if got.HTTPStatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", got.HTTPStatusCode, http.StatusOK)
	}
}

func TestGoSpread_SheetExists(t *testing.T) {
	g := newTestGoSpread(t)
	var res bool
	res = g.SheetExists(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME)
	if res != false {
		t.Errorf("got exists %t want false", res)
	}

	createTestSheet(t, g)
	res = g.SheetExists(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME)
	if res != true {
		t.Errorf("got exists %t want true", res)
	}
	dropTestSheet(t, g)
}

func TestGoSpread_SheetId(t *testing.T) {
	g := newTestGoSpread(t)
	teardown := setupTestSheet(t, g)
	defer teardown()

	s, err := g.GetSpreadSheet(TEST_SPREAD_SHEET_ID)
	if err != nil {
		t.Fatal(err)
	}

	got, err := g.sheetId(s, TEST_SHEET_NAME)
	if err != nil {
		t.Error(err)
	}
	if got == 0 {
		t.Errorf("sheet id should be > 0, got %d", got)
	}
}

func TestGoSpread_AppendRows(t *testing.T) {
	g := newTestGoSpread(t)
	teardown := setupTestSheet(t, g)
	defer teardown()

	err := g.AppendRows(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME, [][]interface{}{{"a", "b", "c"}, {"1", "2", "3"}})
	if err != nil {
		t.Error(err)
	}
}

func TestGoSpread_GetValues(t *testing.T) {
	tests := []struct {
		name   string
		_range string
		want   [][]interface{}
	}{
		{
			name: "A1:C2", _range: "A1:C2",
			want: [][]interface{}{
				{"a", "b", "c"},
				{"1", "2", "3"},
			},
		},
		{
			name: "A1:A", _range: "A1:A",
			want: [][]interface{}{
				{"a"}, {"1"}, {"a3"}, {"a4"},
			},
		},
	}

	fixture := [][]interface{}{
		{"a", "b", "c", "d", "e"},
		{1, 2, 3, 4, 5},
		{"a3", "b3", "c3", "d3", "e3"},
		{"a4", "b4", "c4", "d4", "e4"},
	}

	g := newTestGoSpread(t)
	teardown := setupTestSheet(t, g)
	defer teardown()
	if err := g.AppendRows(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME, fixture); err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.GetValues(TEST_SPREAD_SHEET_ID, TEST_SHEET_NAME, tt._range)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got.Values, tt.want) {
				t.Errorf("got values %q, want %q", got.Values, tt.want)
			}
		})
	}
}

func TestGoSpread_SortSheets(t *testing.T) {
	g := newTestGoSpread(t)
	if err := g.SortSheets(TEST_SPREAD_SHEET_ID); err != nil {
		t.Error(err)
	}
}

func TestGoSpread_MoveSheet(t *testing.T) {
	name := "move-test"
	g := newTestGoSpread(t)
	if err := g.CreateSheet(TEST_SPREAD_SHEET_ID, name); err != nil {
		t.Fatal(err)
	}
	id, err := g.SheetId(TEST_SPREAD_SHEET_ID, name)
	if err != nil {
		t.Fatal(err)
	}

	if err := g.MoveSheet(id, 0, TEST_SPREAD_SHEET_ID); err != nil {
		t.Error(err)
	}
	if err := g.DeleteSheetByName(TEST_SPREAD_SHEET_ID, name); err != nil {
		t.Error(err)
	}
}
