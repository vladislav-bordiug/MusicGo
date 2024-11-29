package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	return nil
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
	database.On("SelectDataQuery", mock.Anything, group, song, date, text, link).
		Return(AnswerData{Items: []RowDbData{{Group: "Muse", Song: "Supermassive Black Hole", Date: "16.07.2006", Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight", Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw"}}}, nil).
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
	database.AssertExpectations(t)
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
	database.On("SelectCoupletQuery", mock.Anything, group, song, couplet).
		Return(AnswerCoupletData{Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?"}, nil).
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
	database.AssertExpectations(t)
}
