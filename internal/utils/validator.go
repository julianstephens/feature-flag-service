package utils

import "fmt"

type Validator struct{}

func (v Validator) ValidatePassword(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("password cannot be empty")
	}
	if len(input) < 4 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}
