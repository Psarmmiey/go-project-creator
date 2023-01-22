package utils

import (
	"strings"
	"unicode"
)

func ToLowerCamel(s string) string {
	var result string
	var words []string
	l := 0
	for i, r := range s {
		if i != 0 && unicode.IsUpper(r) {
			words = append(words, s[l:i])
			l = i
		}
	}
	words = append(words, s[l:])
	for k, word := range words {
		if k == 0 {
			result += strings.ToLower(word)
		} else {
			result += ToTitle(strings.ToLower(word))
		}
	}
	return result
}

func ToUnderscore(s string) string {
	var result string
	for i, r := range s {
		if i != 0 && unicode.IsUpper(r) {
			result += "_"
		}
		result += string(unicode.ToLower(r))
	}
	return result
}

func ToTitle(s string) string {
	var result string
	for i, v := range s {
		if i == 0 {
			result += string(unicode.ToUpper(v))
		} else if unicode.IsUpper(v) {
			result += "" + string(v)
		} else {
			result += string(v)
		}
	}
	return result
}
