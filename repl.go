package repl

import (
	"strings"
	"fmt"
)

func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string	
	}
	
	var low string
	var output []string

	low = strings.ToLower(text)

	output = strings.Fields(low)
	
	return output
}