package writemodel

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pavbis/zal-case-study/application/types"
)

const languageID = "34ffdec9-26e4-4c2f-b9ae-4dc9cb647dc5"

// ErrorWhileFetchingData simulates client's fetch data error response
var ErrorWhileFetchingData = errors.New("error while fetching data")

// ErrStorage simulates storage's error
var ErrStorage = errors.New("storage error")

// FakeHTTPClientWithError is fake client which provokes ErrorWhileFetchingData
type FakeHTTPClientWithError struct{}

func (f *FakeHTTPClientWithError) FetchData(language string) (*types.GitHubJSONResponse, error) {
	return nil, ErrorWhileFetchingData
}

// FakeHTTPClientWithoutError simulates valid client response
type FakeHTTPClientWithoutError struct{}

func (f *FakeHTTPClientWithoutError) FetchData(language string) (*types.GitHubJSONResponse, error) {
	return &types.GitHubJSONResponse{}, nil
}

// StorageWhichReturnsLanguageID simulates valid storage result
type StorageWhichReturnsLanguageID struct{}

func (s *StorageWhichReturnsLanguageID) PersistProgrammingLanguageRepositories(gh *types.GitHubJSONResponse) (types.LanguageID, error) {
	newUUID := uuid.MustParse(languageID)

	return types.LanguageID{UUID: newUUID}, nil
}

// StorageWhichReturnsError simulates storage error
type StorageWhichReturnsError struct{}

func (s *StorageWhichReturnsError) PersistProgrammingLanguageRepositories(gh *types.GitHubJSONResponse) (types.LanguageID, error) {
	return types.LanguageID{}, ErrStorage
}

func Test_WithClientError(t *testing.T) {
	pl := &types.ProgrammingLanguage{Name: "test"}
	client := &FakeHTTPClientWithError{}
	storage := &StorageWhichReturnsLanguageID{}
	commandHandler := NewWriteLanguageRepositoriesCommandHandler(client, storage)

	_, err := commandHandler.HandleRepositories(pl)

	if !errors.Is(err, ErrorWhileFetchingData) {
		t.Errorf("got result %d but expected %d", err, ErrorWhileFetchingData)
	}

	if err.Error() != "error while fetching data" {
		t.Errorf("got result %q but expected %q", err.Error(), "error while fetching data")
	}
}

func Test_WithStorageError(t *testing.T) {
	pl := &types.ProgrammingLanguage{Name: "test"}
	client := &FakeHTTPClientWithoutError{}
	storage := &StorageWhichReturnsError{}
	commandHandler := NewWriteLanguageRepositoriesCommandHandler(client, storage)

	_, err := commandHandler.HandleRepositories(pl)

	if !errors.Is(err, ErrStorage) {
		t.Errorf("got result %d but expected %d", err, ErrorWhileFetchingData)
	}

	if err.Error() != "storage error" {
		t.Errorf("got result %q but expected %q", err.Error(), "storage error")
	}
}

func Test_WithoutError(t *testing.T) {
	pl := &types.ProgrammingLanguage{Name: "test"}
	client := &FakeHTTPClientWithoutError{}
	storage := &StorageWhichReturnsLanguageID{}
	commandHandler := NewWriteLanguageRepositoriesCommandHandler(client, storage)

	result, _ := commandHandler.HandleRepositories(pl)

	if result.UUID.String() != languageID {
		t.Errorf("got result %s but expected %s", result, languageID)
	}
}
