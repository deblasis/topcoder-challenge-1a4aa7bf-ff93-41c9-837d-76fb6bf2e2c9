package data

import "errors"

func StringNotEmptyValidator(s string) error {
	if s == "" {
		return errors.New("Should not be empty")
	}
	return nil
}
