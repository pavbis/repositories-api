package types

import "github.com/google/uuid"

// ProgrammingLanguage represents programming language structure
type ProgrammingLanguage struct {
	Name string
}

// LanguageID represents programming language id structure
type LanguageID struct {
	UUID uuid.UUID `json:"language_id"`
}

// SupportedProgrammingLanguageEnum represents the allowed languages
type SupportedProgrammingLanguageEnum string

const (
	Golang     SupportedProgrammingLanguageEnum = "go"
	Java       SupportedProgrammingLanguageEnum = "java"
	PHP        SupportedProgrammingLanguageEnum = "php"
	Javascript SupportedProgrammingLanguageEnum = "javascript"
	Ruby       SupportedProgrammingLanguageEnum = "ruby"
)

// IsValid executes the check on provided value validity
func (s SupportedProgrammingLanguageEnum) IsValid() bool {
	switch s {
	case Golang, Java, PHP, Javascript, Ruby:
		return true
	}

	return false
}
