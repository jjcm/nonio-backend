package models

import (
	"strings"
)

func createURLFromTitle(title string) string {
	convertedString := strings.ToLower(strings.Join(splitByWords(title), "-"))
	if len(convertedString) > 255 {
		convertedString = convertedString[0:254]
	}
	return convertedString
}

func splitByWords(s string) []string {
	trimmedString := strings.TrimSpace(removeReservedChars(s))
	parts := strings.Split(trimmedString, " ")

	var output []string
	for _, part := range parts {
		if strings.TrimSpace(part) == "" { // skip parts that are an empty string
			continue
		}
		output = append(output, part)
	}
	return output
}

// removeReservedChars - strip all chars except those that are safe for URLs
// https://www.ietf.org/rfc/rfc3986.txt
// section 2.3
func removeReservedChars(s string) string {
	safe := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_~ "
	output := ""
	for _, char := range s {
		for _, safeChar := range safe {
			if char == safeChar {
				output = output + string(char)
			}
		}
	}
	return output
}
