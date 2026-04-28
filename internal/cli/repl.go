package cli

import (
	"strings"
)

func CleanInput(text string) []string {
	lowered := strings.ToLower(text)
	separated := strings.Fields(lowered)
	return separated
}

func ConcatCommand(commands []string) string {
	return strings.Join(commands, " ")
}
