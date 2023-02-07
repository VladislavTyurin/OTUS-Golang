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

func checkRepeats(prev rune, current rune) string {
	numOfRepeat, err := strconv.Atoi(string(current))
	if err != nil {
		return string(prev)
	}
	return strings.Repeat(string(prev), numOfRepeat)
}

func escaped(prev rune, escapedPrev bool) bool {
	return prev == '\\' && !escapedPrev
}

func checkInvalid(prev, current rune, escapedPrev bool) bool {
	if escaped(prev, escapedPrev) {
		return !isDigit(current) && current != '\\'
	}
	return !escapedPrev && isDigit(prev) && isDigit(current)
}

func Unpack(src string) (string, error) {
	var result strings.Builder
	var currentRune rune
	var previousRune rune
	var escapedPrevious bool
	var i int

	for i, currentRune = range src {
		if i == 0 {
			if isDigit(currentRune) {
				return "", ErrInvalidString
			}
		} else {
			if checkInvalid(previousRune, currentRune, escapedPrevious) {
				return "", ErrInvalidString
			}

			if !isDigit(previousRune) && previousRune != '\\' || escapedPrevious {
				result.WriteString(checkRepeats(previousRune, currentRune))
			}
		}

		escapedPrevious = escaped(previousRune, escapedPrevious)
		previousRune = currentRune
	}

	if !escapedPrevious && previousRune == '\\' {
		return "", ErrInvalidString
	}

	if i > 0 && (!isDigit(previousRune) && previousRune != '\\' || escapedPrevious) {
		result.WriteRune(previousRune)
	}

	return result.String(), nil
}
