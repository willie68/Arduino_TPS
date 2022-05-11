package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateID() string {
	uuidStr := uuid.NewString()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	return uuidStr
}

// NOTE: this isn't multi-Unicode-codepoint aware, like specifying skintone or
//       gender of an emoji: https://unicode.org/emoji/charts/full-emoji-modifiers.html
func Substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
