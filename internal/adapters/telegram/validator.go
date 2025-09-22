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

func validateDate(t time.Time) error {
	if t.UTC().Before(time.Now().UTC()) {
		return fmt.Errorf("%s: %s", validationErr.Error(), "incorrect date")
	}
	return nil
}
