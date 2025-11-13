package validation

import (
	"regexp"

	"cruder/internal/model"
)

var (
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]{3,32}$`)
	emailPattern    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	uuidPattern     = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

func ValidateUsername(username string) error {
	if username == "" {
		return model.NewValidationError("username is required")
	}
	if !usernamePattern.MatchString(username) {
		return model.NewValidationError("username is invalid")
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return model.NewValidationError("email is required")
	}
	if !emailPattern.MatchString(email) {
		return model.NewValidationError("email is invalid")
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if fullName == "" {
		return model.NewValidationError("full name is required")
	}
	if len(fullName) > 100 {
		return model.NewValidationError("full name is too long")
	}
	return nil
}

func ValidateUUID(value string) error {
	if value == "" {
		return model.NewValidationError("uuid is required")
	}
	if !uuidPattern.MatchString(value) {
		return model.NewValidationError("uuid is invalid")
	}
	return nil
}

func ValidateID(id int64) error {
	if id < 1 {
		return model.NewValidationError("id must be positive")
	}
	return nil
}

func ValidateCreateUserInput(username, email, fullName string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidateFullName(fullName); err != nil {
		return err
	}
	return nil
}

func ValidateUpdateUserInput(username, email, fullName string) error {
	if username == "" && email == "" && fullName == "" {
		return model.NewValidationError("no fields to update")
	}
	if username != "" {
		if err := ValidateUsername(username); err != nil {
			return err
		}
	}
	if email != "" {
		if err := ValidateEmail(email); err != nil {
			return err
		}
	}
	if fullName != "" {
		if err := ValidateFullName(fullName); err != nil {
			return err
		}
	}
	return nil
}