package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestCreateTableQuery(t *testing.T) {
	mockk, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	database := NewPGXDatabase(mockk)
	defer mockk.Close()
	mockk.ExpectExec("CREATE TABLE IF NOT EXISTS songs").WillReturnResult(pgxmock.NewResult("CREATE", 1))
	err = database.CreateTableQuery(context.Background())
	assert.NoError(t, err)
	if err := mockk.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertQuery(t *testing.T) {
	mockk, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	database := NewPGXDatabase(mockk)
	defer mockk.Close()
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	mockk.ExpectExec("INSERT INTO songs").
		WithArgs(group, song, date, text, link).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	err = database.InsertQuery(context.Background(), group, song, date, text, link)
	assert.NoError(t, err)
	if err := mockk.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteQuery(t *testing.T) {
	mockk, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	database := NewPGXDatabase(mockk)
	defer mockk.Close()
	group := "Muse"
	song := "Supermassive Black Hole"
	mockk.ExpectExec("DELETE FROM songs").
		WithArgs(group, song).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	err = database.DeleteQuery(context.Background(), group, song)
	assert.NoError(t, err)
	if err := mockk.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectDataQuery(t *testing.T) {
	mockk, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	database := NewPGXDatabase(mockk)
	defer mockk.Close()
	page := int64(1)
	items := int64(10)
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	mockk.ExpectQuery("SELECT \\* FROM").
		WithArgs(group, song, date, text, link, items, (page-1)*items).
		WillReturnRows(pgxmock.NewRows([]string{"group_name", "song_name", "releaseDate", "text", "link"}).
			AddRow(group, song, date, text, link))
	_, err = database.SelectDataQuery(context.Background(), page, items, group, song, date, text, link)
	assert.NoError(t, err)
	if err := mockk.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSelectCoupletQuery(t *testing.T) {
	mockk, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	database := NewPGXDatabase(mockk)
	defer mockk.Close()
	couplet := int64(1)
	group := "Muse"
	song := "Supermassive Black Hole"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	mockk.ExpectQuery("SELECT text FROM songs").
		WithArgs(group, song).
		WillReturnRows(pgxmock.NewRows([]string{"text"}).
			AddRow(text))
	_, err = database.SelectCoupletQuery(context.Background(), group, song, couplet)
	assert.NoError(t, err)
	if err := mockk.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEditQuery(t *testing.T) {
	mockk, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	database := NewPGXDatabase(mockk)
	defer mockk.Close()
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	mockk.ExpectExec("UPDATE songs SET").
		WithArgs(group, song, date, text, link).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	err = database.EditQuery(context.Background(), group, song, date, text, link)
	assert.NoError(t, err)
	if err := mockk.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

type MockDatabase struct {
	mock.Mock
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}

func (m *MockDatabase) InsertQuery(ctx context.Context, group_name string, song_name string, releaseDate string, text string, link string) error {
	args := m.Called(ctx, group_name, song_name, releaseDate, text, link)
	return args.Error(0)
}

func (m *MockDatabase) CreateTableQuery(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDatabase) DeleteQuery(ctx context.Context, group_name string, song_name string) error {
	args := m.Called(ctx, group_name, song_name)
	return args.Error(0)
}

func (m *MockDatabase) SelectDataQuery(ctx context.Context, page int64, items int64, group string, song string, releaseDate string, text string, link string) (AnswerData, error) {
	args := m.Called(ctx, group, song, releaseDate, text, link)
	return args.Get(0).(AnswerData), args.Error(1)
}

func (m *MockDatabase) SelectCoupletQuery(ctx context.Context, group string, song string, couplet int64) (AnswerCoupletData, error) {
	args := m.Called(ctx, group, song, couplet)
	return args.Get(0).(AnswerCoupletData), args.Error(1)
}

func (m *MockDatabase) EditQuery(ctx context.Context, group_name string, song_name string, releaseDate string, text string, link string) error {
	args := m.Called(ctx, group_name, song_name, releaseDate, text, link)
	return args.Error(0)
}

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestAddSong(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	requestData := AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	httpclient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/info" && req.URL.RawQuery == "group=Muse&song=Supermassive+Black+Hole"
	})).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"releaseDate": "16.07.2006", "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight", "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"}`))),
	}, nil).
		Once()
	database.On("InsertQuery", mock.Anything, requestData.Group, requestData.Song, "16.07.2006", "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight", "https://www.youtube.com/watch?v=Xsp3_a-PMTw").
		Return(nil).
		Once()
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.addsong(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	httpclient.AssertExpectations(t)
	database.AssertExpectations(t)
}

func TestAddSong_ReadBodyError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}
	req, err := http.NewRequest("POST", "/addsong", io.NopCloser(&errorReader{}))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.addsong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error reading body")
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error reading body")
}

func TestAddSong_UnmarshalError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}

	invalidJSON := `{"Group": "Muse", "Song":`
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader([]byte(invalidJSON)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.addsong(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "unexpected end of JSON input")
}

func TestAddSong_CreateRequestError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}
	t.Setenv("API_URL", "::invalid-url")
	requestData := AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.addsong(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "parse \"::invalid-url/info?group=Muse&song=Supermassive+Black+Hole\": missing protocol scheme")
}

func TestAddSong_ClientDoError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}
	requestData := AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	httpclient.On("Do", mock.Anything).
		Return(&http.Response{
			StatusCode: 500,
		}, fmt.Errorf("mocked error")).
		Once()
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.addsong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "mocked error")
}

func TestAddSong_ReadBodyError2(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}
	requestData := AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	httpclient.On("Do", mock.Anything).
		Return(&http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(&errorReader{}),
		}, nil).
		Once()
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.addsong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error reading body")
}

func TestAddSong_UnmarshalError2(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	requestData := AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	httpclient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/info" && req.URL.RawQuery == "group=Muse&song=Supermassive+Black+Hole"
	})).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"releaseDate": "16.07.2006", "text":`))),
	}, nil).
		Once()
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.addsong(rr, req)
	httpclient.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "unexpected end of JSON input")
}

func TestAddSong_InsertQueryError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	requestData := AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	httpclient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/info" && req.URL.RawQuery == "group=Muse&song=Supermassive+Black+Hole"
	})).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"releaseDate": "16.07.2006", "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight", "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"}`))),
	}, nil).
		Once()
	database.On("InsertQuery", mock.Anything, requestData.Group, requestData.Song, "16.07.2006", "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight", "https://www.youtube.com/watch?v=Xsp3_a-PMTw").
		Return(fmt.Errorf("database error")).
		Once()
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.addsong(rr, req)
	httpclient.AssertExpectations(t)
	database.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "database error")
}

func TestDeleteSong(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	requestData := AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	database.On("DeleteQuery", mock.Anything, requestData.Group, requestData.Song).
		Return(nil).
		Once()
	req, err := http.NewRequest("POST", "/deletesong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.deletesong(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	database.AssertExpectations(t)
}

func TestDeleteSong_ReadBodyError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}
	req, err := http.NewRequest("POST", "/deletesong", io.NopCloser(&errorReader{}))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.deletesong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error reading body")
}

func TestDeleteSong_UnmarshalError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}
	invalidJSON := `{"Group": "Muse", "Song":`
	req, err := http.NewRequest("POST", "/deletesong", bytes.NewReader([]byte(invalidJSON)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.deletesong(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "unexpected end of JSON input")
}

func TestDeleteSong_DeleteQueryError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	requestData := AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	database.On("DeleteQuery", mock.Anything, requestData.Group, requestData.Song).
		Return(fmt.Errorf("database error")).
		Once()
	req, err := http.NewRequest("POST", "/deletesong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.deletesong(rr, req)
	database.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "database error")
}

func TestEditSong(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	requestData := EditRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
		Date:  "16.07.2006",
		Text:  "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:  "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	requestBody, _ := json.Marshal(requestData)
	database.On("EditQuery", mock.Anything, requestData.Group, requestData.Song, requestData.Date, requestData.Text, requestData.Link).Return(nil).
		Return(nil).
		Once()
	req, err := http.NewRequest("POST", "/editsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.editsong(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	database.AssertExpectations(t)
}

func TestEditSong_ReadBodyError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}
	req, err := http.NewRequest("POST", "/editsong", io.NopCloser(&errorReader{}))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.editsong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error reading body")
}

func TestEditSong_UnmarshalError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		Database: database,
		Client:   &httpclient,
	}
	invalidJSON := `{"Group": "Muse", "Song":`
	req, err := http.NewRequest("POST", "/editsong", bytes.NewReader([]byte(invalidJSON)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.editsong(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "unexpected end of JSON input")
}

func TestEditSong_EditQueryError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	requestData := EditRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
		Date:  "16.07.2006",
		Text:  "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:  "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	requestBody, _ := json.Marshal(requestData)
	database.On("EditQuery", mock.Anything, requestData.Group, requestData.Song, requestData.Date, requestData.Text, requestData.Link).Return(nil).
		Return(fmt.Errorf("database error")).
		Once()
	req, err := http.NewRequest("POST", "/editsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.editsong(rr, req)
	database.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "database error")
}

func TestGetData(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	page := 1
	items := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	expectedResponse := AnswerData{
		Items: []RowDbData{
			{
				Group: group,
				Song:  song,
				Date:  date,
				Text:  text,
				Link:  link,
			},
		},
	}
	database.On("SelectDataQuery", mock.Anything, group, song, date, text, link).
		Return(expectedResponse, nil).
		Once()
	urlStr := fmt.Sprintf("/getdata?page=%d&items=%d&group=%s&song=%s&releaseDate=%s&text=%s&link=%s",
		page, items, url.QueryEscape(group), url.QueryEscape(song), url.QueryEscape(date), url.QueryEscape(text), url.QueryEscape(link))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.getdata(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualResponse AnswerData
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
	database.AssertExpectations(t)
}

func TestGetData_ParseIntPageError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	page := "1d"
	items := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	urlStr := fmt.Sprintf("/getdata?page=%s&items=%d&group=%s&song=%s&releaseDate=%s&text=%s&link=%s",
		page, items, url.QueryEscape(group), url.QueryEscape(song), url.QueryEscape(date), url.QueryEscape(text), url.QueryEscape(link))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.getdata(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "strconv.ParseInt: parsing \"1d\": invalid syntax")
}

func TestGetData_ParseIntItemsError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	page := 1
	items := "1d"
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	urlStr := fmt.Sprintf("/getdata?page=%d&items=%s&group=%s&song=%s&releaseDate=%s&text=%s&link=%s",
		page, items, url.QueryEscape(group), url.QueryEscape(song), url.QueryEscape(date), url.QueryEscape(text), url.QueryEscape(link))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.getdata(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "strconv.ParseInt: parsing \"1d\": invalid syntax")
}

func TestGetData_SelectDataQueryError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	page := 1
	items := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	database.On("SelectDataQuery", mock.Anything, group, song, date, text, link).
		Return(AnswerData{}, fmt.Errorf("database error")).
		Once()
	urlStr := fmt.Sprintf("/getdata?page=%d&items=%d&group=%s&song=%s&releaseDate=%s&text=%s&link=%s",
		page, items, url.QueryEscape(group), url.QueryEscape(song), url.QueryEscape(date), url.QueryEscape(text), url.QueryEscape(link))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.getdata(rr, req)
	database.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "database error")
}

func TestGetData_EncodeError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	page := 1
	items := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	expectedResponse := AnswerData{
		Items: []RowDbData{
			{
				Group: group,
				Song:  song,
				Date:  date,
				Text:  text,
				Link:  link,
			},
		},
	}
	database.On("SelectDataQuery", mock.Anything, group, song, date, text, link).
		Return(expectedResponse, nil).
		Once()
	urlStr := fmt.Sprintf("/getdata?page=%d&items=%d&group=%s&song=%s&releaseDate=%s&text=%s&link=%s",
		page, items, url.QueryEscape(group), url.QueryEscape(song), url.QueryEscape(date), url.QueryEscape(text), url.QueryEscape(link))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)
	rr := httptest.NewRecorder()
	handler.getdata(&errorWriter{ResponseRecorder: rr}, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, logBuffer.String(), "mocked write error")
}

type errorWriter struct {
	*httptest.ResponseRecorder
}

func (ew *errorWriter) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("mocked write error")
}

func TestGetSongText(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	couplet := int64(1)
	group := "Muse"
	song := "Supermassive Black Hole"
	expectedResponse := AnswerCoupletData{
		Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?",
	}
	database.On("SelectCoupletQuery", mock.Anything, group, song, couplet).
		Return(expectedResponse, nil).
		Once()
	urlStr := fmt.Sprintf("/getsongtext?couplet=%d&group=%s&song=%s",
		couplet, url.QueryEscape(group), url.QueryEscape(song))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.getsongtext(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualResponse AnswerCoupletData
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
	database.AssertExpectations(t)
}

func TestGetSongText_ParseIntCoupletError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	couplet := "1d"
	group := "Muse"
	song := "Supermassive Black Hole"
	urlStr := fmt.Sprintf("/getsongtext?couplet=%s&group=%s&song=%s",
		couplet, url.QueryEscape(group), url.QueryEscape(song))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.getsongtext(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "strconv.ParseInt: parsing \"1d\": invalid syntax")
}

func TestGetSongText_SelectCoupletQueryError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	couplet := int64(1)
	group := "Muse"
	song := "Supermassive Black Hole"
	database.On("SelectCoupletQuery", mock.Anything, group, song, couplet).
		Return(AnswerCoupletData{}, fmt.Errorf("database error")).
		Once()
	urlStr := fmt.Sprintf("/getsongtext?couplet=%d&group=%s&song=%s",
		couplet, url.QueryEscape(group), url.QueryEscape(song))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.getsongtext(rr, req)
	database.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "database error")
}

func TestGetSongText_EncodeError(t *testing.T) {
	database := NewMockDatabase()
	httpclient := MockHttpClient{}
	handler := &Handler{
		database,
		&httpclient,
	}
	couplet := int64(1)
	group := "Muse"
	song := "Supermassive Black Hole"
	expectedResponse := AnswerCoupletData{
		Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?",
	}
	database.On("SelectCoupletQuery", mock.Anything, group, song, couplet).
		Return(expectedResponse, nil).
		Once()
	urlStr := fmt.Sprintf("/getsongtext?couplet=%d&group=%s&song=%s",
		couplet, url.QueryEscape(group), url.QueryEscape(song))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)
	rr := httptest.NewRecorder()
	handler.getsongtext(&errorWriter{ResponseRecorder: rr}, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, logBuffer.String(), "mocked write error")
}
