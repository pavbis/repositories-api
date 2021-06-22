package writemodel

import (
	"github.com/pavbis/zal-case-study/application/client"
	"github.com/pavbis/zal-case-study/application/storage"
	"github.com/pavbis/zal-case-study/application/types"
)

// WriteOperationsHandler handles data between external and internal storage
type WriteOperationsHandler interface {
	HandleRepositories(pl *types.ProgrammingLanguage) (types.LanguageID, error)
}

type writeLanguageRepositoriesCommandHandler struct {
	client  client.HTTPClient
	storage storage.PersistsProgrammingLanguage
}

// NewWriteLanguageRepositoriesCommandHandler creates new instance of writeLanguageRepositoriesCommandHandler in valid state
func NewWriteLanguageRepositoriesCommandHandler(
	c client.HTTPClient, s storage.PersistsProgrammingLanguage) WriteOperationsHandler {
	return &writeLanguageRepositoriesCommandHandler{c, s}
}

// HandleRepositories fetches the data and persists it
func (ch *writeLanguageRepositoriesCommandHandler) HandleRepositories(pl *types.ProgrammingLanguage) (types.LanguageID, error) {
	var languageID types.LanguageID

	respData, err := ch.client.FetchData(pl.Name)

	if err != nil {
		return languageID, err
	}

	languageID, err = ch.storage.PersistProgrammingLanguageRepositories(respData)

	if err != nil {
		return languageID, err
	}

	return languageID, nil
}
