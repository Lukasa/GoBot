package irc

import (
	"github.com/Lukasa/GoBot/struc"
	"testing"
)

// Test that the applyFilters function works as expected. Use an IRC message that matches every filter.
// Confirm that the output matches what we expect.
func TestApplyFiltersGood(t *testing.T) {
	goodMsg := struc.NewIRCMessage()
	goodMsg.Trailing = "This is a test."
	regexFilter, _ := RegexFilterFromRegex("^This is a (.*)")

	filters := []Filter{
		YesFilter,
		regexFilter,
	}

	goodArgs := []string{
		"This is a test.",
		"test.",
	}

	goodKwargs := map[string]string{}

	success, args, kwargs := applyFilters(goodMsg, filters)
	if !success {
		t.Errorf("Successful filters failed.")
	}

	for i, arg := range goodArgs {
		if arg != args[i] {
			t.Errorf("Unexpected argument: expected %v, got %v", arg, args[i])
		}
	}

	for key, value := range goodKwargs {
		if val, ok := kwargs[key]; ok {
			if val != value {
				t.Errorf("Incorrect value: expected %v, got %v", value, val)
			}
			delete(kwargs, key)
		} else {
			t.Errorf("Missing kwarg %v", key)
		}

		if len(kwargs) > 0 {
			t.Errorf("Too many kwargs: %v", kwargs)
		}
	}
}

func TestApplyFiltersBad(t *testing.T) {
	badMsg := struc.NewIRCMessage()
	badMsg.Trailing = "This is not a test."
	regexFilter, _ := RegexFilterFromRegex("^This is a (.*)$")

	filters := []Filter{
		regexFilter,
	}

	badArgs := []string{}

	badKwargs := map[string]string{}

	success, args, kwargs := applyFilters(badMsg, filters)
	if success {
		t.Errorf("Unsuccessful filters succeeded.")
	}

	for i, arg := range badArgs {
		if arg != args[i] {
			t.Errorf("Unexpected argument: expected %v, got %v", arg, args[i])
		}
	}

	for key, value := range badKwargs {
		if val, ok := kwargs[key]; ok {
			if val != value {
				t.Errorf("Incorrect value: expected %v, got %v", value, val)
			}
			delete(kwargs, key)
		} else {
			t.Errorf("Missing kwarg %v", key)
		}

		if len(kwargs) > 0 {
			t.Errorf("Too many kwargs: %v", kwargs)
		}
	}
}
