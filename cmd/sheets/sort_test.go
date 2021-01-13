package sheets

import (
	"os"
	"testing"
)

func TestSheetsSort_Success(t *testing.T) {
	_, err := runCommand(NewCmdSort, os.Getenv("TEST_SPREAD_SHEET_ID"))
	if err != nil {
		t.Errorf("error running command `sheets sort`: %v", err)
	}
}
func TestSheetsSort_Fail(t *testing.T) {
	_, err := runCommand(NewCmdSort, "12345")
	if err == nil || err.Error() != "googleapi: Error 404: Requested entity was not found., notFound" {
		t.Errorf("expected error, got %v", err)
	}
}
