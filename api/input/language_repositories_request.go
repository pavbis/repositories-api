package input

import (
	"net/http"

	"github.com/gorilla/mux"
)

type LanguageRepositoriesRequest struct {
	LanguageName string `validate:"required,supportedLanguage"`
}

// NewLanguageRepositoriesRequest creates valid receive event input
func NewLanguageRepositoriesRequest(r *http.Request) *LanguageRepositoriesRequest {
	vars := mux.Vars(r)

	return &LanguageRepositoriesRequest{LanguageName: vars["languageName"]}
}
