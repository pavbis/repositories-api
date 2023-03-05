package input

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type LanguageRepositoriesRequest struct {
	LanguageName string `validate:"required,supportedLanguage"`
}

// NewLanguageRepositoriesRequest creates valid receive event input
func NewLanguageRepositoriesRequest(r *http.Request) *LanguageRepositoriesRequest {
	language := chi.URLParam(r, "languageName")

	return &LanguageRepositoriesRequest{LanguageName: language}
}
