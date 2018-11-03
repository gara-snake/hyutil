package hyutil

import (
	"strings"
)

const underscore = byte('_')

//CamelToSnake はcamelCaseの文字列をsnake_caseに変換します
func CamelToSnake(str string) string {

	retbuf := make([]byte, 0)

	if len(str) <= 2 {
		return strings.ToLower(str)
	}

	for i, b := range []byte(str) {

		// 大文字を探す
		if isUpper(b) {

			if 0 < i {
				retbuf = append(retbuf, underscore)
			}
			retbuf = append(retbuf, toLower(b))

		} else {
			retbuf = append(retbuf, b)
		}

	}

	return string(retbuf)
}

//SnakeToUcamel はsnake_caseの文字列をUpperCamelCaseに変換します
func SnakeToUcamel(str string) string {
	retbuf := make([]byte, 0)

	hit := false

	for i, b := range []byte(str) {

		if i == 0 {
			retbuf = append(retbuf, toUpper(b))
		} else {

			// アンダースコアを探す
			if b == underscore {
				hit = true
			} else {
				if hit {
					retbuf = append(retbuf, toUpper(b))
					hit = false
				} else {
					retbuf = append(retbuf, b)
				}
			}

		}

	}

	return string(retbuf)
}

func isWord(c byte) bool {
	return isLetter(c) || isDigit(c)
}

func isLetter(c byte) bool {
	return isLower(c) || isUpper(c)
}

func isUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

func isLower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func toLower(c byte) byte {
	if isUpper(c) {
		return c + ('a' - 'A')
	}
	return c
}

func toUpper(c byte) byte {
	if isLower(c) {
		return c - ('a' - 'A')
	}
	return c
}
