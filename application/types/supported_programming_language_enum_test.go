package types

import (
	"testing"
)

func Test_SupportedProgrammingLanguage(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult bool
	}{
		{
			name:           "Test with go",
			input:          "go",
			expectedResult: true,
		},
		{
			name:           "Test with java",
			input:          "java",
			expectedResult: true,
		},
		{
			name:           "Test with php",
			input:          "php",
			expectedResult: true,
		},
		{
			name:           "Test with javascript",
			input:          "javascript",
			expectedResult: true,
		},
		{
			name:           "Test with ruby",
			input:          "ruby",
			expectedResult: true,
		},
		{
			name:           "Test with invalid language",
			input:          "rust",
			expectedResult: false,
		},
	}

	for _, test := range tests {
		supportedMethod := SupportedProgrammingLanguageEnum(test.input)
		result := supportedMethod.IsValid()

		if result != test.expectedResult {
			t.Errorf(
				"for supported programming language test '%s', got result %t but expected %t",
				test.name,
				result,
				test.expectedResult,
			)
		}
	}
}
