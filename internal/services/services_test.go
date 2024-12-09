package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http"
	"test/internal/models"
	"testing"
)

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

func (m *MockDatabase) SelectDataQuery(ctx context.Context, page int64, items int64, group string, song string, releaseDate string, text string, link string) (models.AnswerData, error) {
	args := m.Called(ctx, page, items, group, song, releaseDate, text, link)
	return args.Get(0).(models.AnswerData), args.Error(1)
}

func (m *MockDatabase) SelectCoupletQuery(ctx context.Context, group string, song string, couplet int64) (models.AnswerCoupletData, error) {
	args := m.Called(ctx, group, song, couplet)
	return args.Get(0).(models.AnswerCoupletData), args.Error(1)
}

func (m *MockDatabase) EditQuery(ctx context.Context, group_name string, song_name string, releaseDate string, text string, link string) error {
	args := m.Called(ctx, group_name, song_name, releaseDate, text, link)
	return args.Error(0)
}

type MockHttpClient struct {
	mock.Mock
}

func NewMockHttpClient() *MockHttpClient {
	return &MockHttpClient{}
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestAddSong(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	responseData := models.AddResponseData{
		Date: "16.07.2006",
		Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		t.Fatal(err)
	}
	client.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/info" && req.URL.RawQuery == "group=Muse&song=Supermassive+Black+Hole"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(jsonData)),
	}, nil).
		Once()
	database.On("InsertQuery", context.Background(), group, song, responseData.Date, responseData.Text, responseData.Link).
		Return(nil).
		Once()
	err, status := service.AddSong(group, song)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, status)
	database.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestAddSong_NewRequestError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "://", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	err, status := service.AddSong(group, song)
	log.Println(err)
	assert.EqualError(t, err, `parse ":///info?group=Muse&song=Supermassive+Black+Hole": missing protocol scheme`)
	assert.Equal(t, http.StatusBadRequest, status)
}

func TestAddSong_DoRequestError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	client.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/info" && req.URL.RawQuery == "group=Muse&song=Supermassive+Black+Hole"
	})).Return(&http.Response{
		StatusCode: http.StatusInternalServerError,
	}, errors.New("Error doing request")).
		Once()
	err, status := service.AddSong(group, song)
	assert.Equal(t, errors.New("Error doing request"), err)
	assert.Equal(t, http.StatusInternalServerError, status)
	client.AssertExpectations(t)
}

type errReader struct{}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading body")
}

func (e *errReader) Close() error {
	return nil
}

func TestAddSong_ReadRespBodyError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	client.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/info" && req.URL.RawQuery == "group=Muse&song=Supermassive+Black+Hole"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       &errReader{},
	}, nil).
		Once()
	err, status := service.AddSong(group, song)
	assert.Equal(t, errors.New("error reading body"), err)
	assert.Equal(t, http.StatusInternalServerError, status)
	client.AssertExpectations(t)
}

func TestAddSong_UnmarshalRespBodyError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	invalidJSON := `{"group": "Muse", "song":`
	client.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/info" && req.URL.RawQuery == "group=Muse&song=Supermassive+Black+Hole"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(invalidJSON))),
	}, nil).
		Once()
	err, status := service.AddSong(group, song)
	assert.EqualError(t, err, "unexpected end of JSON input")
	assert.Equal(t, http.StatusInternalServerError, status)
	client.AssertExpectations(t)
}

func TestAddSong_InsertQuerryError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	responseData := models.AddResponseData{
		Date: "16.07.2006",
		Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		t.Fatal(err)
	}
	client.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/info" && req.URL.RawQuery == "group=Muse&song=Supermassive+Black+Hole"
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(jsonData)),
	}, nil).
		Once()
	database.On("InsertQuery", context.Background(), group, song, responseData.Date, responseData.Text, responseData.Link).
		Return(errors.New("Error inserting song")).
		Once()
	err, status := service.AddSong(group, song)
	assert.Equal(t, errors.New("Error inserting song"), err)
	assert.Equal(t, http.StatusInternalServerError, status)
	database.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestDeleteSong(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	database.On("DeleteQuery", context.Background(), group, song).
		Return(nil).
		Once()
	err, status := service.DeleteSong(group, song)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, status)
	database.AssertExpectations(t)
}

func TestDeleteSong_DeleteQueryError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	database.On("DeleteQuery", context.Background(), group, song).
		Return(errors.New("Error deleting song")).
		Once()
	err, status := service.DeleteSong(group, song)
	assert.Equal(t, errors.New("Error deleting song"), err)
	assert.Equal(t, http.StatusInternalServerError, status)
	database.AssertExpectations(t)
}

func TestEditSong(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	database.On("EditQuery", context.Background(), group, song, date, text, link).
		Return(nil).
		Once()
	err, status := service.EditSong(group, song, date, text, link)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, status)
	database.AssertExpectations(t)
}

func TestEditSong_EditQueryError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	database.On("EditQuery", context.Background(), group, song, date, text, link).
		Return(errors.New("Error editing song")).
		Once()
	err, status := service.EditSong(group, song, date, text, link)
	assert.Equal(t, errors.New("Error editing song"), err)
	assert.Equal(t, http.StatusInternalServerError, status)
	database.AssertExpectations(t)
}

func TestGetSongs(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	page := int64(1)
	items := int64(1)
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	database.On("SelectDataQuery", context.Background(), page, items, group, song, date, text, link).
		Return(models.AnswerData{}, nil).
		Once()
	_, err, status := service.GetSongs(page, items, group, song, date, text, link)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, status)
	database.AssertExpectations(t)
}

func TestGetSongs_SelectDataQueryError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	page := int64(1)
	items := int64(1)
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	database.On("SelectDataQuery", context.Background(), page, items, group, song, date, text, link).
		Return(models.AnswerData{}, errors.New("Error selecting data")).
		Once()
	_, err, status := service.GetSongs(page, items, group, song, date, text, link)
	assert.Equal(t, errors.New("Error selecting data"), err)
	assert.Equal(t, http.StatusInternalServerError, status)
	database.AssertExpectations(t)
}

func TestGetSongText(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	couplet := int64(1)
	group := "Muse"
	song := "Supermassive Black Hole"
	database.On("SelectCoupletQuery", context.Background(), group, song, couplet).
		Return(models.AnswerCoupletData{}, nil).
		Once()
	_, err, status := service.GetSongText(couplet, group, song)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, status)
	database.AssertExpectations(t)
}

func TestGetSongText_SelectCoupletQueryError(t *testing.T) {
	database := NewMockDatabase()
	client := NewMockHttpClient()
	service := NewService(database, "http://localhost:8080", client)
	couplet := int64(1)
	group := "Muse"
	song := "Supermassive Black Hole"
	database.On("SelectCoupletQuery", context.Background(), group, song, couplet).
		Return(models.AnswerCoupletData{}, errors.New("Error selecting data")).
		Once()
	_, err, status := service.GetSongText(couplet, group, song)
	assert.Equal(t, errors.New("Error selecting data"), err)
	assert.Equal(t, http.StatusInternalServerError, status)
	database.AssertExpectations(t)
}
