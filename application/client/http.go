package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pavbis/repositories-api/application/types"
)

type HTTPClient interface {
	FetchData(language string) (*types.GitHubJSONResponse, error)
}

// realHTTPClient is a wrapper to make real HTTP requests.
type realHTTPClient struct {
	client http.Client
}

// NewRealHTTPClient creates a RealHttpClient.
func NewRealHTTPClient() HTTPClient {
	return &realHTTPClient{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *realHTTPClient) FetchData(language string) (*types.GitHubJSONResponse, error) {
	gitHubURL := fmt.Sprintf(
		"https://api.github.com/search/repositories?q=stars:>=10000+language:%s&sort=stars&order=desc&per_page=100",
		language)

	resp, err := c.client.Get(gitHubURL)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fetchData: external service status code is not 200")
	}

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if err != nil {
		return nil, err
	}

	ghr := types.GitHubJSONResponse{}
	ghr.ProgrammingLanguage.Name = language

	err = json.Unmarshal(body, &ghr)

	if err != nil {
		return nil, err
	}

	return &ghr, nil
}
