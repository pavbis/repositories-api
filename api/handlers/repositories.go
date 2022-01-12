package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pavbis/repositories-api/api/input"
	"github.com/pavbis/repositories-api/application/client"
	"github.com/pavbis/repositories-api/application/storage"
	"github.com/pavbis/repositories-api/application/types"
	"github.com/pavbis/repositories-api/application/writemodel"

	"github.com/go-playground/validator/v10"
)

// ReceiveRepositoriesRequestHandler handles incoming request and executes storage's write operation
func ReceiveRepositoriesRequestHandler(db storage.Executor, w http.ResponseWriter, r *http.Request) {
	receiveRepositoriesRequest := input.NewLanguageRepositoriesRequest(r)

	v := validator.New()
	_ = v.RegisterValidation("supportedLanguage", func(fl validator.FieldLevel) bool {
		return types.SupportedProgrammingLanguageEnum(fl.Field().String()).IsValid()
	})

	if err := v.Struct(receiveRepositoriesRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpClient := client.NewRealHTTPClient()
	writeStorage := storage.NewPostgresWriteStore(db)
	pl := &types.ProgrammingLanguage{Name: receiveRepositoriesRequest.LanguageName}
	commandHandler := writemodel.NewWriteLanguageRepositoriesCommandHandler(httpClient, writeStorage)

	result, err := commandHandler.HandleRepositories(pl)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		fmt.Sprintf("Language %s was successfully created with id %s", pl.Name, result))
}

// ReadRepositoriesRequestHandler handles incoming request and executes storage's read operation
func ReadRepositoriesRequestHandler(db storage.Executor, w http.ResponseWriter, r *http.Request) {
	readRepositoriesRequest := input.NewLanguageRepositoriesRequest(r)

	v := validator.New()
	// add custom validation rule
	_ = v.RegisterValidation("supportedLanguage", func(fl validator.FieldLevel) bool {
		return types.SupportedProgrammingLanguageEnum(fl.Field().String()).IsValid()
	})

	if err := v.Struct(readRepositoriesRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	readStorage := storage.NewPostgresReadStore(db)
	pl := &types.ProgrammingLanguage{Name: readRepositoriesRequest.LanguageName}
	result, err := readStorage.ReadRepositoriesForLanguage(pl)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, result)
}

// RemoveRepositoryRequestHandler handles incoming request and executes storage's remove operation
func RemoveRepositoryRequestHandler(db storage.Executor, w http.ResponseWriter, r *http.Request) {
	removeRepoRequest, err := input.NewRemoveRepositoryRequest(r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	repoName := &types.RepositoryID{UUID: removeRepoRequest.RepositoryID}
	writeStorage := storage.NewPostgresWriteStore(db)
	result, err := writeStorage.RemoveRepository(repoName)

	if err != nil {
		if errors.Is(err, storage.ErrRepoNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, fmt.Sprintf("succesfully deleted repository %s", result.UUID.String()))
}

// TopRepositoryForLanguageRequestHandler executes storage's read top list operation
func TopRepositoryForLanguageRequestHandler(db storage.Executor, w http.ResponseWriter, r *http.Request) {
	readStore := storage.NewPostgresReadStore(db)
	result, err := readStore.ReadTopRepositoriesPerLanguage()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, result)
}

// CountRepositoriesStarsForLanguagesRequestHandler executes storage's count repositories operation
func CountRepositoriesStarsForLanguagesRequestHandler(db storage.Executor, w http.ResponseWriter, r *http.Request) {
	readStore := storage.NewPostgresReadStore(db)
	result, err := readStore.ReadRepositoriesSumsForLanguages()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, result)
}

// ListLanguagesAndRepositoriesRequestHandler executes storage's count repositories operation
func ListLanguagesAndRepositoriesRequestHandler(db storage.Executor, w http.ResponseWriter, r *http.Request) {
	readStore := storage.NewPostgresReadStore(db)
	result, err := readStore.ReadRepositoriesAndLanguages()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, result)
}
