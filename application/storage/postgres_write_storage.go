package storage

import (
	"errors"

	"github.com/pavbis/repositories-api/application/types"
)

// ErrRepoNotFound represents error in case the repository is not found
var ErrRepoNotFound = errors.New("repository not found")

type postgresWriteStorage struct {
	sqlExecutor Executor
}

// NewPostgresWriteStore creates new write store instance in valid state
func NewPostgresWriteStore(e Executor) RepresentsWriteStorage {
	return &postgresWriteStorage{sqlExecutor: e}
}

// PersistProgrammingLanguageRepositories handles the whole database write operation
func (s *postgresWriteStorage) PersistProgrammingLanguageRepositories(gh *types.GitHubJSONResponse) (types.LanguageID, error) {
	languageID, err := s.persistProgrammingLanguage(gh)

	if err != nil {
		return languageID, err
	}

	for _, repo := range gh.Items {
		_, err = s.sqlExecutor.Exec(
			`INSERT INTO repositories ("repositoryId", "languageId", full_name, stars, "createdAt", owner, description)
VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6)
ON CONFLICT ("languageId", full_name)
    DO UPDATE SET stars       = EXCLUDED.stars,
                  description = EXCLUDED.description;`,
			languageID.UUID.String(), repo.FullName, repo.StargazersCount, repo.CreatedAt, repo.Owner.Login, repo.Description)

		if err != nil {
			return languageID, err
		}
	}

	return languageID, nil
}

func (s *postgresWriteStorage) persistProgrammingLanguage(gh *types.GitHubJSONResponse) (types.LanguageID, error) {
	var languageID types.LanguageID

	err := s.sqlExecutor.QueryRow(
		`INSERT INTO programming_languages("languageId", "language_name")
		VALUES (uuid_generate_v4(), $1)
		ON CONFLICT ("language_name") DO UPDATE SET "updated_at" = NOW()
		RETURNING "languageId";`,
		gh.ProgrammingLanguage.Name).Scan(&languageID.UUID)

	if err != nil {
		return languageID, err
	}

	return languageID, nil
}

func (s *postgresWriteStorage) RemoveRepository(rn *types.RepositoryID) (*types.RepositoryID, error) {
	result, err := s.sqlExecutor.Exec(`DELETE FROM repositories r WHERE r."repositoryId" = $1`, rn.UUID.String())

	if err != nil {
		return nil, err
	}

	affectedRows, err := result.RowsAffected()

	if err != nil {
		return nil, err
	}

	if affectedRows == 0 {
		return nil, ErrRepoNotFound
	}

	return rn, nil
}
