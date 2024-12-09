package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"test/internal/models"
	"testing"
)

type MockInterface struct {
	mock.Mock
}

func NewMockInterface() *MockInterface {
	return &MockInterface{}
}

func (m *MockInterface) AddSong(group string, song string) (err error, status int) {
	args := m.Called(group, song)
	return args.Error(0), args.Get(1).(int)
}

func (m *MockInterface) DeleteSong(group string, song string) (err error, status int) {
	args := m.Called(group, song)
	return args.Error(0), args.Get(1).(int)
}

func (m *MockInterface) EditSong(group string, song string, date string, text string, link string) (err error, status int) {
	args := m.Called(group, song, date, text, link)
	return args.Error(0), args.Get(1).(int)
}

func (m *MockInterface) GetSongs(page int64, items int64, group string, song string, date string, text string, link string) (result models.AnswerData, err error, status int) {
	args := m.Called(page, items, group, song, date, text, link)
	return args.Get(0).(models.AnswerData), args.Error(1), args.Get(2).(int)
}

func (m *MockInterface) GetSongText(couplet int64, group string, song string) (result models.AnswerCoupletData, err error, status int) {
	args := m.Called(couplet, group, song)
	return args.Get(0).(models.AnswerCoupletData), args.Error(1), args.Get(2).(int)
}

func TestAddSong(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	requestData := models.AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	mockinterface.On("AddSong", requestData.Group, requestData.Song).
		Return(nil, http.StatusOK).
		Once()
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.AddSong(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockinterface.AssertExpectations(t)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading body")
}

func TestAddSong_ReadBodyError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	req, err := http.NewRequest("POST", "/addsong", io.NopCloser(errReader(0)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.AddSong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error reading body")
}

func TestAddSong_UnmarshalError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	invalidJSON := `{"group": "Muse", "song":`
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader([]byte(invalidJSON)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.AddSong(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "unexpected end of JSON input")
}

func TestAddSong_AddSongError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	requestData := models.AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	mockinterface.On("AddSong", requestData.Group, requestData.Song).
		Return(errors.New("error adding song"), http.StatusInternalServerError).
		Once()
	req, err := http.NewRequest("POST", "/addsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.AddSong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	log.Println(rr.Code)
	mockinterface.AssertExpectations(t)
	assert.Contains(t, rr.Body.String(), "error adding song")
}

func TestDeleteSong(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	requestData := models.AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	mockinterface.On("DeleteSong", requestData.Group, requestData.Song).
		Return(nil, http.StatusOK).
		Once()
	req, err := http.NewRequest("POST", "/deletesong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.DeleteSong(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockinterface.AssertExpectations(t)
}

func TestDeleteSong_ReadBodyError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	req, err := http.NewRequest("POST", "/deletesong", io.NopCloser(errReader(0)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.DeleteSong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error reading body")
}

func TestDeleteSong_UnmarshalError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	invalidJSON := `{"group": "Muse", "song":`
	req, err := http.NewRequest("POST", "/deletesong", bytes.NewReader([]byte(invalidJSON)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.DeleteSong(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "unexpected end of JSON input")
}

func TestDeleteSong_AddSongError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	requestData := models.AddDeleteRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
	}
	requestBody, _ := json.Marshal(requestData)
	mockinterface.On("DeleteSong", requestData.Group, requestData.Song).
		Return(errors.New("error deleting song"), http.StatusInternalServerError).
		Once()
	req, err := http.NewRequest("POST", "/deletesong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.DeleteSong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	log.Println(rr.Code)
	mockinterface.AssertExpectations(t)
	assert.Contains(t, rr.Body.String(), "error deleting song")
}

func TestEditSong(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	requestData := models.EditRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
		Date:  "16.07.2006",
		Text:  "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:  "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	requestBody, _ := json.Marshal(requestData)
	mockinterface.On("EditSong", requestData.Group, requestData.Song, requestData.Date, requestData.Text, requestData.Link).
		Return(nil, http.StatusOK).
		Once()
	req, err := http.NewRequest("POST", "/editsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.EditSong(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockinterface.AssertExpectations(t)
}

func TestEditSong_ReadBodyError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	req, err := http.NewRequest("POST", "/editsong", io.NopCloser(errReader(0)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.EditSong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error reading body")
}

func TestEditSong_UnmarshalError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	invalidJSON := `{"group": "Muse", "song":`
	req, err := http.NewRequest("POST", "/editsong", bytes.NewReader([]byte(invalidJSON)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.EditSong(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "unexpected end of JSON input")
}

func TestEditSong_AddSongError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	requestData := models.EditRequestData{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
		Date:  "16.07.2006",
		Text:  "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:  "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	requestBody, _ := json.Marshal(requestData)
	mockinterface.On("EditSong", requestData.Group, requestData.Song, requestData.Date, requestData.Text, requestData.Link).
		Return(errors.New("error deleting song"), http.StatusInternalServerError).
		Once()
	req, err := http.NewRequest("POST", "/editsong", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.EditSong(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	log.Println(rr.Code)
	mockinterface.AssertExpectations(t)
	assert.Contains(t, rr.Body.String(), "error deleting song")
}

func TestGetSongs(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	page := 1
	items := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	expectedResponse := models.AnswerData{
		Items: []models.RowDbData{
			{
				Group: group,
				Song:  song,
				Date:  date,
				Text:  text,
				Link:  link,
			},
		},
	}
	mockinterface.On("GetSongs", int64(page), int64(items), group, song, date, text, link).
		Return(expectedResponse, nil, http.StatusOK).
		Once()
	urlStr := fmt.Sprintf("/getdata?page=%d&items=%d&group=%s&song=%s&releaseDate=%s&text=%s&link=%s",
		page, items, url.QueryEscape(group), url.QueryEscape(song), url.QueryEscape(date), url.QueryEscape(text), url.QueryEscape(link))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.GetSongs(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualResponse models.AnswerData
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
	mockinterface.AssertExpectations(t)
}

func TestGetSongs_ParseIntPageError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
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
	handler.GetSongs(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "strconv.ParseInt: parsing \"1d\": invalid syntax")
}

func TestGetSongs_ParseIntItemsError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
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
	handler.GetSongs(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "strconv.ParseInt: parsing \"1d\": invalid syntax")
}

func TestGetSongs_GetSongsError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	page := 1
	items := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	mockinterface.On("GetSongs", int64(page), int64(items), group, song, date, text, link).
		Return(models.AnswerData{}, errors.New("error getting songs"), http.StatusInternalServerError).
		Once()
	urlStr := fmt.Sprintf("/getdata?page=%d&items=%d&group=%s&song=%s&releaseDate=%s&text=%s&link=%s",
		page, items, url.QueryEscape(group), url.QueryEscape(song), url.QueryEscape(date), url.QueryEscape(text), url.QueryEscape(link))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.GetSongs(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, rr.Body.String(), "error getting songs\n")
	mockinterface.AssertExpectations(t)
}

type errorWriter struct {
	*httptest.ResponseRecorder
}

func (ew *errorWriter) Write(data []byte) (int, error) {
	ew.Body.Write([]byte("forced encoding error"))
	return 0, fmt.Errorf("forced encoding error")
}

func TestGetSongs_NewEncoderEncodeError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	page := 1
	items := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	date := "16.07.2006"
	text := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"
	link := "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
	expectedResponse := models.AnswerData{
		Items: []models.RowDbData{
			{
				Group: group,
				Song:  song,
				Date:  date,
				Text:  text,
				Link:  link,
			},
		},
	}
	mockinterface.On("GetSongs", int64(page), int64(items), group, song, date, text, link).
		Return(expectedResponse, nil, http.StatusOK).
		Once()
	urlStr := fmt.Sprintf("/getdata?page=%d&items=%d&group=%s&song=%s&releaseDate=%s&text=%s&link=%s",
		page, items, url.QueryEscape(group), url.QueryEscape(song), url.QueryEscape(date), url.QueryEscape(text), url.QueryEscape(link))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.GetSongs(&errorWriter{ResponseRecorder: rr}, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "forced encoding error")
	mockinterface.AssertExpectations(t)
}

func TestGetSongText(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	couplet := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	expectedResponse := models.AnswerCoupletData{
		Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?",
	}
	mockinterface.On("GetSongText", int64(couplet), group, song).
		Return(expectedResponse, nil, http.StatusOK).
		Once()
	urlStr := fmt.Sprintf("/getsongtext?couplet=%d&group=%s&song=%s",
		couplet, url.QueryEscape(group), url.QueryEscape(song))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.GetSongText(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualResponse models.AnswerCoupletData
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
	mockinterface.AssertExpectations(t)
}

func TestGetSongText_ParseIntPageError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
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
	handler.GetSongText(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "strconv.ParseInt: parsing \"1d\": invalid syntax")
}

func TestGetSongText_GetSongsError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	couplet := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	mockinterface.On("GetSongText", int64(couplet), group, song).
		Return(models.AnswerCoupletData{}, errors.New("error getting song text"), http.StatusInternalServerError).
		Once()
	urlStr := fmt.Sprintf("/getsongtext?couplet=%d&group=%s&song=%s",
		couplet, url.QueryEscape(group), url.QueryEscape(song))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.GetSongText(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, rr.Body.String(), "error getting song text\n")
	mockinterface.AssertExpectations(t)
}

func TestGetSongText_NewEncoderEncodeError(t *testing.T) {
	mockinterface := NewMockInterface()
	handler := &Handler{
		mockinterface,
	}
	couplet := 1
	group := "Muse"
	song := "Supermassive Black Hole"
	expectedResponse := models.AnswerCoupletData{
		Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?",
	}
	mockinterface.On("GetSongText", int64(couplet), group, song).
		Return(expectedResponse, nil, http.StatusOK).
		Once()
	urlStr := fmt.Sprintf("/getsongtext?couplet=%d&group=%s&song=%s",
		couplet, url.QueryEscape(group), url.QueryEscape(song))
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.GetSongText(&errorWriter{ResponseRecorder: rr}, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "forced encoding error")
	mockinterface.AssertExpectations(t)
}
