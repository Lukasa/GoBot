package util

import (
	"bytes"
	"testing"
)

// Test the nasty-ass Bash format string conversion.
func TestBashFmtStringToGoFmtString(t *testing.T) {
	inputs := []string{
		"You're doing great work, ${1}!",
		"Test some stuff ${0}${1}.",
		"Test in${1}side a word.",
		"Test ordering ${1} is right ${0}.",
	}

	outputs := []string{
		"You're doing great work, %v!",
		"Test some stuff %v%v.",
		"Test in%vside a word.",
		"Test ordering %v is right %v.",
	}

	indices := [][]int{
		[]int{1},
		[]int{0, 1},
		[]int{1},
		[]int{1, 0},
	}

	for i, input := range inputs {
		output, outindices := BashFmtStringToGoFmtString(input)

		if bytes.Compare([]byte(output), []byte(outputs[i])) != 0 {
			t.Errorf("Incorrect output:\n\texpected: %v,\n\tgot: %v.\n", []byte(outputs[i]), []byte(output))
		}

		if len(indices[i]) != len(outindices) {
			t.Errorf("Incorrect number of indices: expected %v, got %v.", len(indices[i]), outindices)
		}

		for j, index := range outindices {
			if index != indices[i][j] {
				t.Errorf("Invalid index at %v: expected %v, got %v.", j, indices[i][j], index)
			}
		}
	}
}
