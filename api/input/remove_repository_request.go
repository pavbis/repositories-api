package input

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var ErrRepoID = errors.New("missing or invalid repository id provided")

type RemoveRepositoryRequest struct {
	RepositoryID uuid.UUID
}

func NewRemoveRepositoryRequest(r *http.Request) (*RemoveRepositoryRequest, error) {
	vars := mux.Vars(r)

	// parse the provided repository id and ensure it's a valid uuid
	repoID, err := uuid.Parse(vars["repositoryId"])

	if err != nil {
		return nil, ErrRepoID
	}

	return &RemoveRepositoryRequest{RepositoryID: repoID}, nil
}
