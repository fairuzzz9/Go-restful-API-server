package utilities

import (
	"unicode"
)

// Rot18 = Rot13(alphabets) + rot5(numeric)
// ROT13 basically rotates a character by 13 places, ‘A’ to ‘N’, ‘B’ to ‘M’ and so on.
// The ROT5, rotates the digits: ‘0’ to ‘5’, ‘1’ to ‘6’ and so on.

func Rot18(input string) string {

	var result []rune
	rot5map := map[rune]rune{'0': '5', '1': '6', '2': '7', '3': '8', '4': '9', '5': '0', '6': '1', '7': '2', '8': '3', '9': '4'}

	for _, i := range input {
		switch {
		case !unicode.IsLetter(i) && !unicode.IsNumber(i):
			result = append(result, i)
		case i >= 'A' && i <= 'Z':
			result = append(result, 'A'+(i-'A'+13)%26)
		case i >= 'a' && i <= 'z':
			result = append(result, 'a'+(i-'a'+13)%26)
		case i >= '0' && i <= '9':
			result = append(result, rot5map[i])
		case unicode.IsSpace(i):
			result = append(result, ' ')
		}
	}
	return string(result[:])
}
