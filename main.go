package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type AddDeleteRequestData struct {
	Group string `json:"group" binding:"required" example:"Muse"`
	Song  string `json:"song" binding:"required" example:"Supermassive Black Hole"`
}

type AddResponseData struct {
	Date string `json:"releaseDate"`
	Text string `json:"text"`
	Link string `json:"link"`
}

type EditRequestData struct {
	Group string `json:"group" binding:"required" example:"Muse"`
	Song  string `json:"song" binding:"required" example:"Supermassive Black Hole"`
	Date  string `json:"releaseDate" example:"16.07.2006"`
	Text  string `json:"text" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"`
	Link  string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

type RowDbData struct {
	Group string `db:"group_name" json:"group" binding:"required" example:"Muse"`
	Song  string `db:"song_name" json:"song" binding:"required" example:"Supermassive Black Hole"`
	Date  string `db:"releaseDate" json:"releaseDate" binding:"required" example:"16.07.2006"`
	Text  string `db:"text" json:"text" binding:"required" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"`
	Link  string `db:"link" json:"link" binding:"required" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

type AnswerData struct {
	Items []RowDbData `json:"items" binding:"required"`
}

type AnswerCoupletData struct {
	Text string `json:"text" binding:"required" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?"`
}

func Config() *pgxpool.Config {
	// const defaultMaxConns = int32(4)
	// const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	// dbConfig.MaxConns = defaultMaxConns
	// dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}

func CreateTableQuery(p *pgxpool.Pool) error {
	_, err := p.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS songs (group_name TEXT, song_name TEXT, releaseDate TEXT, text TEXT, link TEXT);")
	if err != nil {
		return err
	}
	return nil
}

func InsertQuery(p *pgxpool.Pool, group_name string, song_name string, releaseDate string, text string, link string) error {
	_, err := p.Exec(context.Background(), "INSERT INTO songs(group_name, song_name, releaseDate, text, link) values($1, $2, $3, $4, $5)", group_name, song_name, releaseDate, text, link)
	if err != nil {
		return err
	}
	return nil
}

func DeleteQuery(p *pgxpool.Pool, group_name string, song_name string) error {
	_, err := p.Exec(context.Background(), "DELETE FROM songs WHERE group_name = $1 AND song_name = $2", group_name, song_name)
	if err != nil {
		return err
	}
	return nil
}

func SelectDataQuery(p *pgxpool.Pool, page int64, items int64, group string, song string, releaseDate string, text string, link string) (AnswerData, error) {
	query := "SELECT * FROM songs "
	var answer AnswerData
	paramindex := 1
	setClauses := []string{}
	params := []interface{}{}
	if group != "" || song != "" || releaseDate != "" || text != "" || link != "" {
		if group != "" {
			setClauses = append(setClauses, fmt.Sprintf("group_name = $%d", paramindex))
			params = append(params, group)
			paramindex++
		}
		if song != "" {
			setClauses = append(setClauses, fmt.Sprintf("song_name = $%d", paramindex))
			params = append(params, song)
			paramindex++
		}
		if releaseDate != "" {
			setClauses = append(setClauses, fmt.Sprintf("releaseDate = $%d", paramindex))
			params = append(params, releaseDate)
			paramindex++
		}
		if text != "" {
			setClauses = append(setClauses, fmt.Sprintf("text = $%d", paramindex))
			params = append(params, text)
			paramindex++
		}
		if link != "" {
			setClauses = append(setClauses, fmt.Sprintf("link = $%d", paramindex))
			params = append(params, link)
			paramindex++
		}
	}
	if len(setClauses) > 0 {
		query += "WHERE " + strings.Join(setClauses, " AND ") + " "
	}
	query += fmt.Sprintf("LIMIT $%d OFFSET $%d", paramindex, paramindex+1)
	params = append(params, items)
	params = append(params, (page-1)*items)
	log.Printf("INFO: query for the database=%s\n", query)
	log.Printf("INFO: query params for the database=%s\n", params)
	rows, err := p.Query(context.Background(), query, params...)
	if err != nil {
		return answer, err
	}
	defer rows.Close()
	for rows.Next() {
		var result RowDbData
		if err := rows.Scan(&result.Group, &result.Song, &result.Date, &result.Text, &result.Link); err != nil {
			return answer, err
		}
		answer.Items = append(answer.Items, result)
	}
	return answer, nil
}

func SelectCoupletQuery(p *pgxpool.Pool, group string, song string, couplet int64) (AnswerCoupletData, error) {
	var text string
	var answer AnswerCoupletData
	err := p.QueryRow(context.Background(), "SELECT text FROM songs WHERE group_name = $1 AND song_name = $2", group, song).Scan(&text)
	if err != nil {
		return answer, err
	}
	result := strings.Split(text, "\n\n")
	if couplet > int64(len(result)) {
		return answer, fmt.Errorf("There is no such couplet")
	}
	answer.Text = result[couplet-1]
	return answer, nil
}

func EditQuery(p *pgxpool.Pool, group_name string, song_name string, releaseDate string, text string, link string) error {
	query := "UPDATE songs SET "
	paramindex := 3
	setClauses := []string{}
	params := []interface{}{group_name, song_name}
	if releaseDate != "" {
		setClauses = append(setClauses, fmt.Sprintf("releaseDate = $%d", paramindex))
		params = append(params, releaseDate)
		paramindex++
	}
	if text != "" {
		setClauses = append(setClauses, fmt.Sprintf("text = $%d", paramindex))
		params = append(params, text)
		paramindex++
	}
	if link != "" {
		setClauses = append(setClauses, fmt.Sprintf("link = $%d", paramindex))
		params = append(params, link)
		paramindex++
	}
	if len(setClauses) == 0 {
		return nil
	}
	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE group_name = $1 AND song_name = $2")
	log.Printf("INFO: query for the database=%s\n", query)
	log.Printf("INFO: query params for the database=%s\n", params)
	_, err := p.Exec(context.Background(), query, params...)
	if err != nil {
		return err
	}
	return nil
}

var connPool *pgxpool.Pool

func main() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: %s", err)
	}
	connPool, err = pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!", err)
	}
	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!", err)
	}
	defer connection.Release()
	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database", err)
	}
	err = CreateTableQuery(connPool)
	if err != nil {
		log.Fatal("Error while creating table in the database", err)
	}
	log.Printf("INFO: created table\n")
	defer connPool.Close()
	http.HandleFunc("/addsong", addsong)
	http.HandleFunc("/deletesong", deletesong)
	http.HandleFunc("/editsong", editsong)
	http.HandleFunc("/getdata", getdata)
	http.HandleFunc("/getsongtext", getsongtext)
	// http.HandleFunc("/info", info)
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_IP")+":"+os.Getenv("PORT"), nil))
}

// AddSong godoc
// @Summary Add song
// @Description Add song based on group and song provided as json.
// @Tags song
// @Accept json
// @Produce  json
// @Param data body AddDeleteRequestData true "JSON with group and song"
// @Success 200 {object} nil "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /addsong [post]
func addsong(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Received request to add song")
	var respdata AddDeleteRequestData
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Request body: %s\n", string(body))
	if err := json.Unmarshal(body, &respdata); err != nil {
		log.Printf("ERROR: Failed to unmarshal request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Request data: group=%s, song=%s\n", respdata.Group, respdata.Song)
	client := &http.Client{}
	encodedGroup := url.QueryEscape(respdata.Group)
	encodedSong := url.QueryEscape(respdata.Song)
	urlStr := fmt.Sprintf("%s/info?group=%s&song=%s",
		os.Getenv("API_URL"), encodedGroup, encodedSong)
	log.Printf("INFO: Url for request: %s\n", urlStr)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Printf("ERROR: Failed to create API request: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to get additional song data: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Request body: %s\n", string(body))
	var reqdata AddResponseData
	if err := json.Unmarshal(body, &reqdata); err != nil {
		log.Printf("ERROR: Failed to unmarshal request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = InsertQuery(connPool, respdata.Group, respdata.Song, reqdata.Date, reqdata.Text, reqdata.Link)
	if err != nil {
		log.Printf("ERROR: Failed to add song to the database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("INFO: Added song to the database\n")
}

// DeleteSong godoc
// @Summary Delete song
// @Description Delete song based on group and song provided as json.
// @Tags song
// @Accept json
// @Produce  json
// @Param data body AddDeleteRequestData true "JSON with group and song"
// @Success 200 {object} nil "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /deletesong [post]
func deletesong(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Received request to delete song")
	var respdata AddDeleteRequestData
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Request body: %s\n", string(body))
	if err := json.Unmarshal(body, &respdata); err != nil {
		log.Printf("ERROR: Failed to unmarshal request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Request data: group=%s, song=%s\n", respdata.Group, respdata.Song)
	err = DeleteQuery(connPool, respdata.Group, respdata.Song)
	if err != nil {
		log.Printf("ERROR: Failed to delete song from the database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("INFO: Deleted song from the database\n")
}

// EditSong godoc
// @Summary Edit song text
// @Description Edit song releaseDate, text and link based on group and song provided as json.
// @Tags song
// @Accept json
// @Produce  json
// @Param data body EditRequestData true "JSON with group, song, releaseDate, text, and link"
// @Success 200 {object} nil "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /editsong [post]
func editsong(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Received request to edit song")
	var respdata EditRequestData
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Request body: %s\n", string(body))
	if err := json.Unmarshal(body, &respdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("INFO: Request data: group=%s, song=%s, releaseDate=%s, text=%s, link=%s\n", respdata.Group, respdata.Song, respdata.Date, respdata.Text, respdata.Link)
	err = EditQuery(connPool, respdata.Group, respdata.Song, respdata.Date, respdata.Text, respdata.Link)
	if err != nil {
		log.Printf("ERROR: Failed to edit song in the database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
// @Success 200 {object} AnswerData "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /getdata [get]
func getdata(w http.ResponseWriter, r *http.Request) {
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
	result, err := SelectDataQuery(connPool, page, items, group, song, releaseDate, text, link)
	if err != nil {
		log.Printf("ERROR: Failed to get data from the database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
// @Success 200 {object} AnswerCoupletData "OK"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /getsongtext [get]
func getsongtext(w http.ResponseWriter, r *http.Request) {
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
	result, err := SelectCoupletQuery(connPool, group, song, couplet)
	if err != nil {
		log.Printf("ERROR: Failed to get data from the database: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

/*
func info(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	group := query.Get("group")
	song := query.Get("song")
	log.Printf("INFO: Request data: group=%s, song=%s\n", group, song)
	var result AddResponseData
	result = AddResponseData{
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
