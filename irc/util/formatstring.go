package util

import (
	"strconv"
	"strings"
)

// BashFmtStringToGoFmtString takes a format string as used in a Botscript and converts it to a Go format string
// and a slice of indices into the expected arguments.
// This function is a prime candidate for refactoring: it's monolithic and fairly bizarrely complicated.
func BashFmtStringToGoFmtString(in string) (string, []int) {
	inRune := []rune(in)
	outRune := make([]rune, len(inRune))
	indices := make([]int, 0)

	// Loop over the string. If we get a $, it indicates a variable replacement unless prefixed by a
	// backslash.
	for i, j := 0, 0; i < len(inRune); i, j = i+1, j+1 {
		char := inRune[i]

		if (char != '$') || (inRune[i-1] == '\\') {
			outRune[j] = char
			continue
		}

		// We have a '$'. The portion in between braces is the variable name/number.
		// In this early version, it should be a number.
		i += 2
		name := make([]rune, 0, 2)
		for (i < len(inRune)) && (inRune[i] != '}') {
			name = append(name, inRune[i])
			i++
		}

		// The variable name needs to be a number. If it isn't, ignore it and keep going.
		varIndex, err := strconv.Atoi(string(name))
		if err != nil {
			continue
		}

		// In the output string, replace the variable with '%v'.
		outRune[j] = '%'
		outRune[j+1] = 'v'
		j++

		// Now, store the indices.
		indices = append(indices, varIndex)
	}

	outString := string(outRune)
	outString = strings.TrimRight(outString, "\x00") // The fact that I need to do this is crap.

	return outString, indices
}
