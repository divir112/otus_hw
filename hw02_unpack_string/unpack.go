package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const shielding = '\\'

func Unpack(str string) (string, error) {
	var result strings.Builder
	prevIsDigit := true
	prevIsShielding := false
	letters := []rune(str)

	for i := 0; i < len(letters); i++ {
		if err := validate(prevIsShielding, prevIsDigit, letters[i]); err != nil {
			return "", err
		}

		if !prevIsShielding && letters[i] == shielding {
			prevIsShielding = true
			prevIsDigit = false
			continue
		}

		if unicode.IsDigit(letters[i]) && !prevIsShielding {
			prevIsDigit = true
			continue
		}

		if len(letters)-1 == i {
			if !unicode.IsDigit(letters[i]) || prevIsShielding {
				result.WriteRune(letters[i])
			}
			break
		}

		prevIsDigit = false
		if unicode.IsDigit(letters[i+1]) && prevIsShielding || unicode.IsDigit(letters[i+1]) {
			iterations, _ := strconv.Atoi(string(letters[i+1]))
			for j := 0; j < iterations; j++ {
				result.WriteRune(letters[i])
			}
		} else {
			result.WriteRune(letters[i])
		}
		prevIsShielding = false
	}
	return result.String(), nil
}

func validate(prevIsShielding, prevIsDigit bool, letter rune) error {
	if prevIsShielding && (!unicode.IsDigit(letter) && letter != shielding) {
		return ErrInvalidString
	}

	if unicode.IsDigit(letter) && prevIsDigit {
		return ErrInvalidString
	}
	return nil
}
