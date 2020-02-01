package models

import (
	"fmt"
	"strconv"
	"strings"
)

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

// getUniqueCommentParentIDs
// this will take in a slice of comments and return a unique list of parent IDs
// as a comma separated string. it's used to make SQL queries faster
func getUniqueCommentParentIDs(comments []Comment) string {
	fmt.Println("comments:", comments)
	var parentIDs []string

	if len(comments) == 0 {
		return "0"
	}

	for _, c := range comments {
		currentParentID := strconv.Itoa(c.ParentID)
		if !stringInSlice(currentParentID, parentIDs) {
			parentIDs = append(parentIDs, strconv.Itoa(c.ParentID))
		}
	}
	return strings.Join(parentIDs, ",")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
