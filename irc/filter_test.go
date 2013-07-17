package irc

import (
	"github.com/Lukasa/GoBot/struc"
	"testing"
)

// Test the building and function of the regex filter.
func TestRegexpFilterGenerator(t *testing.T) {
	testStrings := []string{
		"!m Lukasa",
		"!m the power of love",
		"!m I love !m don't you?",
		"No",
	}
	testMatches := []bool{true, true, true, false}
	testArgs := [][]string{
		[]string{"!m Lukasa", "Lukasa"},
		[]string{"!m the power of love", "the power of love"},
		[]string{"!m I love !m don't you?", "I love !m don't you?"},
		[]string{},
	}
	testRegex := `^!m (?P<target>.*)$`

	filter, err := RegexFilterFromRegex(testRegex)
	if err != nil {
		t.Errorf("Failed to build filter: %v", err)
	}

	for i, str := range testStrings {
		msg := struc.NewIRCMessage()
		msg.Trailing = str

		match, args, _ := filter(msg)

		if match != testMatches[i] {
			t.Errorf("Inconsistent match: expected %v, got %v", testMatches[i], match)
		}

		for j, testArg := range testArgs[i] {
			if testArg != args[j] {
				t.Errorf("Inconsistent arg: expected %v, got %v", testArg, args[j])
			}
		}
	}
}
