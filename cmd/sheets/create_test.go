package sheets

import (
	"os"
	"testing"
)

func TestSheetsCreate_Success(t *testing.T) {
	_, err := runCommand(NewCmdCreate, "new-sheet-1", os.Getenv("TEST_SPREAD_SHEET_ID"))
	if err != nil {
		t.Error(err)
	}
}
func TestSheetsCreate_Fail(t *testing.T) {
	_, err := runCommand(NewCmdCreate, "12345")
	if err == nil || err.Error() != "accepts 2 arg(s), received 1" {
		t.Errorf("expected error, got %v", err)
	}
}
