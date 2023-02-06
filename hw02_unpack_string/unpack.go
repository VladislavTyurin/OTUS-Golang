package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func checkRepeats(current, next rune) string {
	numOfRepeat, err := strconv.Atoi(string(next))
	if err != nil {
		return string(current)
	}
	return strings.Repeat(string(current), numOfRepeat)
}

func checkInvalidString(current, next rune, escapedCurrent bool) bool {
	return (current == '\\' && !isDigit(next) && next != current) ||
		(!escapedCurrent && isDigit(current) && isDigit(next))
}

func Unpack(src string) (string, error) {
	var result strings.Builder
	var previousRune rune
	var currentRune rune
	var nextRune rune
	runes := []rune(src)
	escapedRunes := make([]bool, len(runes))

	for i := 0; i < len(runes); i++ {
		currentRune = runes[i]
		switch i {
		case 0:
			nextRune = runes[i+1]
			escapedRunes[i] = false
			if isDigit(currentRune) || checkInvalidString(currentRune, nextRune, escapedRunes[i]) {
				return "", ErrInvalidString
			}

			if currentRune != '\\' {
				result.WriteString(checkRepeats(currentRune, nextRune))
			}

		case len(runes) - 1:
			previousRune = runes[i-1]
			escapedRunes[i] = (previousRune == '\\') && !escapedRunes[i-1]

			if !isDigit(currentRune) || escapedRunes[i] {
				result.WriteRune(currentRune)
			}

		default:
			previousRune = runes[i-1]
			nextRune = runes[i+1]

			escapedRunes[i] = (previousRune == '\\') && !escapedRunes[i-1]

			if checkInvalidString(currentRune, nextRune, escapedRunes[i]) {
				return "", ErrInvalidString
			}

			if currentRune == '\\' && !escapedRunes[i] {
				continue
			}

			if !isDigit(currentRune) || escapedRunes[i] {
				result.WriteString(checkRepeats(currentRune, nextRune))
			}
		}
	}

	return result.String(), nil
}
