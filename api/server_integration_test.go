package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/pavbis/repositories-api/application/storage"
	"github.com/pavbis/repositories-api/application/types"
	"github.com/pavbis/repositories-api/application/writemodel"
)

var s Server

func TestMain(m *testing.M) {
	initializeServer()

	code := m.Run()
	os.Exit(code)
}

func TestHealthStatus(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/api/health", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	checkMessageValue(t, response.Body.Bytes(), "status", "OK")
}

func TestPostLanguageWithInvalidLanguageName(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/api/languages/rust", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	checkMessageValue(
		t,
		response.Body.Bytes(),
		"error",
		"Key: 'LanguageRepositoriesRequest.LanguageName' Error:Field validation for 'LanguageName' failed on the 'supportedLanguage' tag")
}

func TestGetRepositoriesWithInvalidLanguageName(t *testing.T) {
	req := authRequest(http.MethodGet, "/api/languages/rust", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	checkMessageValue(
		t,
		response.Body.Bytes(),
		"error",
		"Key: 'LanguageRepositoriesRequest.LanguageName' Error:Field validation for 'LanguageName' failed on the 'supportedLanguage' tag")
}

func TestGetRepositoriesWithValidLanguageNameAndEmptyResultSet(t *testing.T) {
	req := authRequest(http.MethodGet, "/api/languages/java", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBodyIsEmptyArray(t, response.Body)
}

func TestStatisticsListWithEmptyRDBMS(t *testing.T) {
	if err := truncateProgrammingLanguagesTable(); err != nil {
		t.Error(err)
	}

	req := authRequest(http.MethodGet, "/api/stats/top-list", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBodyIsEmptyArray(t, response.Body)
}

func TestGetRepositoriesWithValidLanguageName(t *testing.T) {
	if err := truncateProgrammingLanguagesTable(); err != nil {
		t.Error(err)
	}

	// persist data for go programming language to check the response body
	store := storage.NewPostgresWriteStore(s.db)
	// we are using another client here to prevent real HTTP call to external api
	client := &FakeJSONFileReadingClient{}
	commandHandler := writemodel.NewWriteLanguageRepositoriesCommandHandler(client, store)
	pl := &types.ProgrammingLanguage{Name: "go"}
	// write data to database
	_, _ = commandHandler.HandleRepositories(pl)

	req := authRequest(http.MethodGet, "/api/languages/go", nil)
	response := executeRequest(req)
	expected, _ := readFileContent("testdata/internal_response_data.json")

	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBody(t, response.Body.Bytes(), expected)
}

func TestDeleteRepositoryWithInvalidRepositoryId(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/api/repositories/invalidUuid", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	checkMessageValue(t, response.Body.Bytes(), "error", "missing or invalid repository id provided")
}

func TestDeleteRepositoryWithValidRepositoryIdButInconsistentRepo(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/api/repositories/34ffdec9-26e4-4c2f-b9ae-4dc9cb647dc5", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	checkMessageValue(t, response.Body.Bytes(), "error", "repository not found")
}

func TestDeleteRepositoryWithValidExistingRepositoryId(t *testing.T) {
	var repositoryUUIDAsString string
	_ = s.db.QueryRow(`SELECT "repositoryId" FROM repositories LIMIT 1`).Scan(&repositoryUUIDAsString)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/repositories/%s", repositoryUUIDAsString), nil)
	response := executeRequest(req)
	expected := bytes.NewBufferString(fmt.Sprintf("successfully deleted repository %s", repositoryUUIDAsString))

	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBody(t, response.Body.Bytes(), expected.Bytes())
}

// the repos for golang are persisted ATM
func TestStatisticsTopListWithRecordsInRDBMS(t *testing.T) {
	req := authRequest(http.MethodGet, "/api/stats/top-list", nil)
	response := executeRequest(req)
	expected, _ := readFileContent("testdata/top_list_data.json")

	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBody(t, response.Body.Bytes(), expected)
}

// the repos for golang are persisted ATM
func TestStatisticsCountReposWithRecordsInRDBMS(t *testing.T) {
	req := authRequest(http.MethodGet, "/api/stats/count-repositories", nil)
	response := executeRequest(req)
	expected, _ := readFileContent("testdata/count_repositories_data.json")

	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBody(t, response.Body.Bytes(), expected)
}

func TestStatisticsCountReposWithEmptyRDBMS(t *testing.T) {
	if err := truncateProgrammingLanguagesTable(); err != nil {
		t.Error(err)
	}

	req := authRequest(http.MethodGet, "/api/stats/count-repositories", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	checkResponseBodyIsEmptyArray(t, response.Body)
}

// helper functions start here
// initializes the server, there is no need to execute s.Run(":1111")
// the http test recorder just collects the request/response information
func initializeServer() {
	s = Server{}
	s.Initialize()
}

// creates new instance of HTTP recorder
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)

	return rr
}

// checks the response code returned by defined handler
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code is %d. Got %d", expected, actual)
	}
}

// checks message value in json response for specific key
func checkMessageValue(t *testing.T, body []byte, fieldName string, expected string) {
	var m map[string]interface{}
	_ = json.Unmarshal(body, &m)

	fieldValue := m[fieldName]
	if fieldValue != expected {
		t.Errorf("Expected %v. Got %v", expected, fieldValue)
	}
}

// checks response body
func checkResponseBody(t *testing.T, body []byte, expected []byte) {
	var m1 []interface{}
	_ = json.Unmarshal(body, &m1)

	var m2 []interface{}
	_ = json.Unmarshal(expected, &m2)

	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("\n %v. \n %v", m2, m1)
	}
}

// checks response body is empty array, which is mostly the default COALESCE value
func checkResponseBodyIsEmptyArray(t *testing.T, respBody *bytes.Buffer) {
	expected := bytes.NewBufferString("[]")

	if !bytes.Equal(respBody.Bytes(), expected.Bytes()) {
		t.Errorf("expected %v got %v", expected.Bytes(), respBody.Bytes())
	}
}

// removes all existing programming languages from "programming_languages" table.
func truncateProgrammingLanguagesTable() error {
	if _, err := s.db.Exec(`DELETE FROM programming_languages WHERE "languageId" IS NOT NULL`); err != nil {
		return err
	}

	return nil
}

// reads a content of a file
func readFileContent(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// creates valid auth request
func authRequest(method string, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json; charset=utf-8")
	req.Header.Add("Authorization", basicAuthValue())

	return req
}

// return the auth header value
func basicAuthValue() string {
	auth := os.Getenv("AUTH_USER") + ":" + os.Getenv("AUTH_PASS")
	return "Basic " + base64.URLEncoding.EncodeToString([]byte(auth))
}

// FakeJsonFileReadingClient is a struct which simulates external api response while reading defined json file
type FakeJSONFileReadingClient struct{}

// FetchData fetches data from defined json file and fills the GitHubJSONResponse struct
func (c *FakeJSONFileReadingClient) FetchData(language string) (*types.GitHubJSONResponse, error) {
	fileContent, _ := readFileContent("testdata/external_response_data.json")

	ghr := types.GitHubJSONResponse{}
	ghr.ProgrammingLanguage.Name = language

	err := json.Unmarshal(fileContent, &ghr)

	if err != nil {
		return nil, err
	}

	return &ghr, nil
}
