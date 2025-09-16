package telegram

import (
	"testing"
)

func TestParseAndValidate(t *testing.T) {
	dt1 := "2025-09-17 01:25"
	_, err := parseAndValidateDate(dt1)
	if err != nil {
		t.Error(err)
	}
}
