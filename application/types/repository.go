package types

import "github.com/google/uuid"

// Owner represents the owner structure
type Owner struct {
	Login string
}

// GitHubRepository represents the repository structure
type GitHubRepository struct {
	FullName        string `json:"full_name"`
	Owner           Owner
	Description     string
	CreatedAt       string `json:"created_at"`
	StargazersCount int    `json:"stargazers_count"`
}

// GitHubJSONResponse represents the json structure of the external api response
type GitHubJSONResponse struct {
	ProgrammingLanguage
	Items []GitHubRepository
}

// RepositoryID represents the repository uuid
type RepositoryID struct {
	UUID uuid.UUID `json:"repository_id"`
}
