package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func GenerateID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

func ValidateComment(content string) error {
	if strings.TrimSpace(content) == "" {
		return errors.New("comment cannot be empty")
	}

	if len(content) > 2000 {
		return errors.New("comment exceeds 2000 characters limit")
	}

	return nil
}
