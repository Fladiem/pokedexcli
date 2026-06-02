package main

import (
	"strings"
)

func cleanInput(text string) []string {
	if len(text) == 0 {
		var empty []string
		return empty
	}

	var low string
	var output []string

	low = strings.ToLower(text)

	output = strings.Fields(low)
	
	return output
}