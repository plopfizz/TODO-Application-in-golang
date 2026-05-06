package util

import (
	"errors"
	"strings"
	"time"
)

const dueDateLayout = "2006-01-02"

// ParseDueDate validates and parses the due date in YYYY-MM-DD format.
func ParseDueDate(input string) (time.Time, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return time.Time{}, errors.New("due_date is required and must use YYYY-MM-DD")
	}

	parsed, err := time.Parse(dueDateLayout, trimmed)
	if err != nil {
		return time.Time{}, errors.New("due_date must use YYYY-MM-DD format")
	}
	return parsed.UTC(), nil
}

// ValidateText checks basic constraints for task text.
func ValidateText(text string) (string, error) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "", errors.New("text is required")
	}
	if len(trimmed) > 250 {
		return "", errors.New("text length must be 250 characters or less")
	}
	return trimmed, nil
}
