package storage

import "github.com/pavbis/zal-case-study/application/types"

type postgresReadStorage struct {
	sqlExecutor Executor
}

// NewPostgresReadStore creates new read store instance in valid state
func NewPostgresReadStore(e Executor) RepresentsReadStorage {
	return &postgresReadStorage{sqlExecutor: e}
}

// ReadRepositoriesForLanguage reads repositories from database by provided programming language
func (p *postgresReadStorage) ReadRepositoriesForLanguage(l *types.ProgrammingLanguage) ([]byte, error) {
	row := p.sqlExecutor.QueryRow(
		`SELECT COALESCE((SELECT json_strip_nulls(json_agg(r))
                 FROM (
                          WITH language AS (
                              SELECT "languageId"
                              FROM programming_languages
                              WHERE language_name = $1
                          )
                          SELECT r.full_name,
                                 r.stars,
                                 r.description
                          FROM repositories r
                          WHERE "languageId" = (SELECT "languageId" FROM language)
                          ORDER BY stars DESC
                          LIMIT 1000
                      ) r), '[]')`,
		l.Name)

	return scanOrFail(row)
}

// ReadTopRepositoriesPerLanguage reads repositories from database by provided programming language
func (p *postgresReadStorage) ReadTopRepositoriesPerLanguage() ([]byte, error) {
	row := p.sqlExecutor.QueryRow(
		`SELECT COALESCE((SELECT json_strip_nulls(json_agg(tl))
                 FROM (
                          WITH ranked_repos AS (
                              SELECT pl.language_name,
                                     r.full_name,
                                     r.stars,
                                     RANK() OVER (PARTITION BY "languageId" ORDER BY stars DESC) AS rank
                              FROM repositories r
                                       JOIN programming_languages pl USING ("languageId")
                          )
                          SELECT rr.language_name,
                                 rr.full_name,
                                 rr.stars
                          FROM ranked_repos rr
                          WHERE rank = 1
                          ORDER BY rr.stars DESC
                      ) tl), '[]')`)

	return scanOrFail(row)
}

func (p *postgresReadStorage) ReadRepositoriesSumsForLanguages() ([]byte, error) {
	row := p.sqlExecutor.QueryRow(
		`SELECT COALESCE((SELECT json_strip_nulls(json_agg(tl))
                 FROM (
                          WITH language_sum_starts AS (
                              SELECT pl."languageId"  AS language_id,
                                     pl.language_name AS language_name,
                                     SUM(r.stars)     AS stars_sum
                              FROM programming_languages pl
                                       JOIN repositories r USING ("languageId")
                              GROUP BY language_id, language_name
                              ORDER BY stars_sum DESC
                          )
                          SELECT language_name,
                                 stars_sum,
                                 stars_sum - COALESCE(LAG(stars_sum) OVER (), stars_sum) AS stars_difference
                          FROM language_sum_starts
                      ) tl), '[]');`)

	return scanOrFail(row)
}

func (p *postgresReadStorage) ReadRepositoriesAndLanguages() ([]byte, error) {
	row := p.sqlExecutor.QueryRow(
		`SELECT COALESCE(
    (SELECT json_agg(json_build_object(
         'languageId', "languageId",
         'language_name', language_name,
         'repositories', (
             SELECT COALESCE(
                (SELECT json_agg(json_build_object(
                     'repository_id', "repositoryId",
                     'repository_name', full_name,
                     'stars', stars
                 ) ORDER BY "full_name")
                 FROM repositories r
                 WHERE r."languageId" = pl."languageId"),
                '[]'
            )
         )
     ) ORDER BY pl.language_name)
    FROM programming_languages pl),
'[]')`)

	return scanOrFail(row)
}
