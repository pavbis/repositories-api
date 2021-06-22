package storage

import (
	"database/sql"

	"github.com/pavbis/zal-case-study/application/types"
)

// Executor is the interface for sql operations
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type (
	// RepresentsWriteStorage is a combined interface which holds the both interfaces (PersistsProgrammingLanguage, ProgrammingLanguageDeleter)
	RepresentsWriteStorage interface {
		PersistsProgrammingLanguage
		ProgrammingLanguageRepositoryDeleter
	}

	// PersistsProgrammingLanguage is interface which represents the programming language persist operation
	PersistsProgrammingLanguage interface {
		PersistProgrammingLanguageRepositories(gh *types.GitHubJSONResponse) (types.LanguageID, error)
	}

	// ProgrammingLanguageRepositoryDeleter is interface which represents the repository delete operation
	ProgrammingLanguageRepositoryDeleter interface {
		RemoveRepository(rn *types.RepositoryID) (*types.RepositoryID, error)
	}
)

type (
	// RepresentsReadStorage is a combined interface of all read storage interfaces
	RepresentsReadStorage interface {
		ProvidesRepositoriesForLanguage
		ProvidesTopRepositoriesPerLanguage
		ProvidesRepositoryStarsSum
		ProvidesRepositoriesAndCorrespondingLanguages
	}

	// ProvidesRepositoriesForLanguage represents the read operation by programming language
	ProvidesRepositoriesForLanguage interface {
		ReadRepositoriesForLanguage(l *types.ProgrammingLanguage) ([]byte, error)
	}

	// ProvidesTopRepositoriesPerLanguage represents the read top repositories operation
	ProvidesTopRepositoriesPerLanguage interface {
		ReadTopRepositoriesPerLanguage() ([]byte, error)
	}

	// ProvidesRepositoryStarsSum represents the stars sum read operation
	ProvidesRepositoryStarsSum interface {
		ReadRepositoriesSumsForLanguages() ([]byte, error)
	}

	// ProvidesRepositoriesAndCorrespondingLanguages represents the languages and repositories
	ProvidesRepositoriesAndCorrespondingLanguages interface {
		ReadRepositoriesAndLanguages() ([]byte, error)
	}
)
