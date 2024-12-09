package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"test/internal/database"
	"test/internal/models"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Service struct {
	database database.Database
	apiurl   string
	client   httpClient
}

func NewService(db database.Database, apiurl string, client httpClient) *Service {
	return &Service{database: db, apiurl: apiurl, client: client}
}

func (s *Service) AddSong(group string, song string) (err error, status int) {
	encodedGroup := url.QueryEscape(group)
	encodedSong := url.QueryEscape(song)
	urlStr := fmt.Sprintf("%s/info?group=%s&song=%s",
		s.apiurl, encodedGroup, encodedSong)
	log.Printf("INFO: Url for request: %s\n", urlStr)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create API request: %v\n", err)
		return err, http.StatusBadRequest
	}
	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to get additional song data: %v\n", err)
		return err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response body: %v\n", err)
		return err, http.StatusInternalServerError
	}
	log.Printf("DEBUG: Request body: %s\n", string(body))
	var reqdata models.AddResponseData
	if err = json.Unmarshal(body, &reqdata); err != nil {
		log.Printf("ERROR: Failed to unmarshal response body: %v\n", err)
		return err, http.StatusInternalServerError
	}
	err = s.database.InsertQuery(context.Background(), group, song, reqdata.Date, reqdata.Text, reqdata.Link)
	if err != nil {
		log.Printf("ERROR: Failed to add song to the database: %v\n", err)
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (s *Service) DeleteSong(group string, song string) (err error, status int) {
	err = s.database.DeleteQuery(context.Background(), group, song)
	if err != nil {
		log.Printf("ERROR: Failed to delete song from the database: %v\n", err)
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (s *Service) EditSong(group string, song string, date string, text string, link string) (err error, status int) {
	err = s.database.EditQuery(context.Background(), group, song, date, text, link)
	if err != nil {
		log.Printf("ERROR: Failed to edit song in the database: %v\n", err)
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (s *Service) GetSongs(page int64, items int64, group string, song string, date string, text string, link string) (result models.AnswerData, err error, status int) {
	result, err = s.database.SelectDataQuery(context.Background(), page, items, group, song, date, text, link)
	if err != nil {
		log.Printf("ERROR: Failed to get data from the database: %v\n", err)
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}

func (s *Service) GetSongText(couplet int64, group string, song string) (result models.AnswerCoupletData, err error, status int) {
	result, err = s.database.SelectCoupletQuery(context.Background(), group, song, couplet)
	if err != nil {
		log.Printf("ERROR: Failed to get data from the database: %v\n", err)
		return result, err, http.StatusInternalServerError
	}
	return result, nil, http.StatusOK
}
