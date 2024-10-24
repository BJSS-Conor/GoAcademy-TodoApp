package stringUtils

import "testing"

func TestIsEmptyOrWhitespace(t *testing.T) {
	var testCases = []struct {
		testName string
		input    string
		expected bool
	}{
		{"Non-empty string", "This is a string", false},
		{"String with only spaces", "   ", true},
		{"Empty string", "", true},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			if result := IsEmptyOrWhitespace(test.input); result != test.expected {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", result, test.expected)
			}
		})
	}
}
