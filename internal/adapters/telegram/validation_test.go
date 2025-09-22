package telegram

import (
	"testing"
	"time"
)

var loc, _ = time.LoadLocation("Asia/Yekaterinburg")

func TestParseAndValidate(t *testing.T) {
	dt1 := time.Now().Add(time.Minute)
	err := validateDate(dt1)
	if err != nil {
		t.Error(err)
	}
}

func TestParseAndValidateNegative(t *testing.T) {
	dt1 := time.Now().Add(-time.Minute)
	err := validateDate(dt1)
	if err == nil {
		t.Error()
	}

	dt1 = time.Now()
	err = validateDate(dt1)
	if err == nil {
		t.Error()
	}
}
