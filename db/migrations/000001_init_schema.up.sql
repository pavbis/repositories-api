CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE DOMAIN non_empty AS VARCHAR(100) NOT NULL CHECK ( length(value) > 0 );

CREATE TABLE IF NOT EXISTS "programming_languages"
(
    "languageId"    CHAR(36)        NOT NULL PRIMARY KEY,
    "language_name" non_empty       UNIQUE,
    "created_at"    timestamptz     NOT NULL DEFAULT (NOW()),
    "updated_at"    timestamptz     NOT NULL DEFAULT (NOW())
);

CREATE TABLE IF NOT EXISTS "repositories"
(
    "repositoryId"    CHAR(36)        NOT NULL PRIMARY KEY,
    "languageId"      CHAR(36)        NOT NULL REFERENCES programming_languages("languageId") ON DELETE CASCADE,
    "full_name"       non_empty,
    "stars"           BIGINT          NOT NULL,
    "createdAt"       timestamptz     NOT NULL,
    "owner"           VARCHAR(100)    NOT NULL,
    "description"     TEXT            NOT NULL
);

CREATE UNIQUE INDEX repositories_language_id_full_name_id on repositories ("languageId", full_name);
CREATE INDEX repositories_stars_idx ON repositories (stars DESC);