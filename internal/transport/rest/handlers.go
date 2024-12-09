package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"test/internal/models"
)

type ServiceInterface interface {
	AddSong(group string, song string) (err error, status int)
	DeleteSong(group string, song string) (err error, status int)
	EditSong(group string, song string, date string, text string, link string) (err error, status int)
	GetSongs(page int64, items int64, group string, song string, date string, text string, link string) (result models.AnswerData, err error, status int)
	GetSongText(couplet int64, group string, song string) (result models.AnswerCoupletData, err error, status int)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// AddSong godoc
// @Summary Add song
// @Description Add song based on group and song provided as json.
// @Tags song
// @Accept json
// @Produce  json
// @Param data body models.AddDeleteRequestData true "JSON with group and song"
// @Success 200 {object} nil "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /addsong [post]
func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Received request to add song")
	var respdata models.AddDeleteRequestData
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Request body: %s\n", string(body))
	if err = json.Unmarshal(body, &respdata); err != nil {
		log.Printf("ERROR: Failed to unmarshal request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Request data: group=%s, song=%s\n", respdata.Group, respdata.Song)
	err, status := h.service.AddSong(respdata.Group, respdata.Song)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	log.Printf("INFO: Added song to the database\n")
}

// DeleteSong godoc
// @Summary Delete song
// @Description Delete song based on group and song provided as json.
// @Tags song
// @Accept json
// @Produce  json
// @Param data body models.AddDeleteRequestData true "JSON with group and song"
// @Success 200 {object} nil "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /deletesong [post]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Received request to delete song")
	var respdata models.AddDeleteRequestData
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Request body: %s\n", string(body))
	if err = json.Unmarshal(body, &respdata); err != nil {
		log.Printf("ERROR: Failed to unmarshal request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Request data: group=%s, song=%s\n", respdata.Group, respdata.Song)
	err, status := h.service.DeleteSong(respdata.Group, respdata.Song)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	log.Printf("INFO: Deleted song from the database\n")
}

// EditSong godoc
// @Summary Edit song text
// @Description Edit song releaseDate, text and link based on group and song provided as json.
// @Tags song
// @Accept json
// @Produce  json
// @Param data body models.EditRequestData true "JSON with group, song, releaseDate, text, and link"
// @Success 200 {object} nil "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /editsong [post]
func (h *Handler) EditSong(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Received request to edit song")
	var respdata models.EditRequestData
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Request body: %s\n", string(body))
	if err = json.Unmarshal(body, &respdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Request data: group=%s, song=%s, releaseDate=%s, text=%s, link=%s\n", respdata.Group, respdata.Song, respdata.Date, respdata.Text, respdata.Link)
	err, status := h.service.EditSong(respdata.Group, respdata.Song, respdata.Date, respdata.Text, respdata.Link)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	log.Printf("INFO: Edited song in the database\n")
}

// GetSongs godoc
// @Summary Get all songs and their information with pagination
// @Description Retrieve songs and their details with pagination based on the page and items and filtration based on group, song, releaseDate, text and link provided as query parameters.
// @Tags songs
// @Produce  json
// @Param page query integer true "Current page" example(1)
// @Param items query integer true "Number of elements on the page" example(10)
// @Param group query string false "Group" example("Muse")
// @Param song query string false "Song name" example("Supermassive Black Hole")
// @Param releaseDate query string false "Release date in format DD.MM.YYYY" example("16.07.2006")
// @Param text query string false "Song text (multiline allowed)" example("Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight")
// @Param link query string false "Song link" example("https://www.youtube.com/watch?v=Xsp3_a-PMTw")
// @Success 200 {object} models.AnswerData "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /getdata [get]
func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Received request to get songs")
	query := r.URL.Query()
	page, err := strconv.ParseInt(query.Get("page"), 10, 64)
	if err != nil {
		log.Printf("ERROR: Failed to parse page to int %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	items, err := strconv.ParseInt(query.Get("items"), 10, 64)
	if err != nil {
		log.Printf("ERROR: Failed to parse items to int %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	group := query.Get("group")
	song := query.Get("song")
	releaseDate := query.Get("releaseDate")
	text := query.Get("text")
	link := query.Get("link")
	log.Printf("INFO: Request data: page=%d, items=%d, group=%s, song=%s, releaseDate=%s, text=%s, link=%s\n", page, items, group, song, releaseDate, text, link)
	result, err, status := h.service.GetSongs(page, items, group, song, releaseDate, text, link)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	log.Printf("INFO: Response data: items=%s\n", result.Items)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("ERROR: Failed to encode response: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: Responded\n")
}

// GetSongText godoc
// @Summary Get songs text with pagination
// @Description Retrieve song text with pagination based on the group, song and couplet provided as query parameters.
// @Tags song
// @Produce  json
// @Param group query string true "Group" example("Muse")
// @Param song query string true "Song name" example("Supermassive Black Hole")
// @Param couplet query integer true "Couplet" example(1)
// @Success 200 {object} models.AnswerCoupletData "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /getsongtext [get]
func (h *Handler) GetSongText(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	group := query.Get("group")
	song := query.Get("song")
	couplet, err := strconv.ParseInt(query.Get("couplet"), 10, 64)
	if err != nil {
		log.Printf("ERROR: Failed to parse couplet to int %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Request data: group=%s, song=%s, couplet=%d\n", group, song, couplet)
	result, err, status := h.service.GetSongText(couplet, group, song)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	log.Printf("INFO: Response data: text=%s\n", result.Text)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("ERROR: Failed to encode response: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: Responded\n")
}

// Can be used for /getinfo requests
/*
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	group := query.Get("group")
	song := query.Get("song")
	log.Printf("INFO: Request data: group=%s, song=%s\n", group, song)
	var result models.AddResponseData
	result = models.AddResponseData{
		Date: "16.07.2006",
		Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	log.Printf("INFO: Response data: releaseDate=%s, text=%s, link=%s\n", result.Date, result.Text, result.Link)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("ERROR: Failed to encode response: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
*/
