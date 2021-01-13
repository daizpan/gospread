package sheets

import (
	"os"
	"testing"
)

func TestSheetsMove_Success(t *testing.T) {
	_, err := runCommand(NewCmdMove, "113591128", "0", os.Getenv("TEST_SPREAD_SHEET_ID"))
	if err != nil {
		t.Error(err)
	}
}
func TestSheetsMove_Fail(t *testing.T) {
	_, err := runCommand(NewCmdMove, "12345")
	if err == nil || err.Error() != "accepts 3 arg(s), received 1" {
		t.Errorf("expected error, got %v", err)
	}
}
