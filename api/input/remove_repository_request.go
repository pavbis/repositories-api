package input

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/google/uuid"
)

var ErrRepoID = errors.New("missing or invalid repository id provided")

type RemoveRepositoryRequest struct {
	RepositoryID uuid.UUID
}

func NewRemoveRepositoryRequest(r *http.Request) (*RemoveRepositoryRequest, error) {
	repositoryID := chi.URLParam(r, "repositoryId")

	// parse the provided repository id and ensure it's a valid uuid
	repoID, err := uuid.Parse(repositoryID)

	if err != nil {
		return nil, ErrRepoID
	}

	return &RemoveRepositoryRequest{RepositoryID: repoID}, nil
}
