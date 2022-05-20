package pass

import (
	"errors"
	"fmt"
	"unicode"
)

const (
	passMin = 8
	passMax = 128
)

var (
	ErrWrongLength  = fmt.Errorf("wrong length: min %d max %d", passMin, passMax)
	ErrWrongSymbols = errors.New("invalid symbol: only numbers/letters are allowed")
	ErrWrongFormat  = errors.New("invalid format: at least one number, uppercase and lowercase letter")
)

func Validate(password string) error {

	// password check Why not regexp? Because re2 does not support lookaheads '?= '
	count, low, up, num := 0, false, false, false
	for _, s := range password {
		if !unicode.IsLetter(s) && !unicode.IsNumber(s) {
			return ErrWrongSymbols
		}
		switch {
		case unicode.IsLower(s):
			low = true
		case unicode.IsUpper(s):
			up = true
		case unicode.IsNumber(s):
			num = true
		}
		count++
	}

	if passMin > count || count > passMax {
		return ErrWrongLength
	}
	if !(low && up && num) {
		return ErrWrongFormat
	}

	return nil
}
