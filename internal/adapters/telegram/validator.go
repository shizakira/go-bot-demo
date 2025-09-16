package telegram

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var validationErr = errors.New("validation error")
var v = validator.New()

func validateTitle(s string) error {
	if err := v.Var(s, `required,min=3`); err != nil {
		return fmt.Errorf("%s: %s", validationErr.Error(), "title required min 3")
	}
	return nil
}

func validateDesc(s string) error {
	if err := v.Var(s, `required`); err != nil {
		return fmt.Errorf("%s: %w", validationErr.Error(), err)
	}
	return nil
}

func validateDate(t1 time.Time) error {
	t2 := time.Now()
	if !t1.After(t2) {
		return fmt.Errorf("%s: %s", validationErr.Error(), "incorrect date")
	}
	return nil
}

func parseAndValidateDate(datetime string) (*time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02 15:04", datetime)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", validationErr.Error(), "incorrect date")
	}
	err = validateDate(parsedDate)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", validationErr.Error(), "incorrect date")
	}
	return &parsedDate, nil
}
