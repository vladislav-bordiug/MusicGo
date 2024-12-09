package app

import (
	"context"
	"net/http"
	"test/internal/database"
	"test/internal/services"
	"test/internal/transport/rest"
)

type App struct {
	pool   database.DBPool
	ip     string
	port   string
	apiurl string
}

func NewApp(pool database.DBPool, ip string, port string, apiurl string) *App {
	return &App{pool: pool, ip: ip, port: port, apiurl: apiurl}
}
func (a *App) Run() error {
	db := database.NewPGXDatabase(a.pool)
	err := db.CreateTableQuery(context.Background())
	if err != nil {
		return err
	}
	client := &http.Client{}
	tokenservice := services.NewService(db, a.apiurl, client)
	handler := rest.NewHandler(tokenservice)
	http.HandleFunc("/addsong", handler.AddSong)
	http.HandleFunc("/deletesong", handler.DeleteSong)
	http.HandleFunc("/editsong", handler.EditSong)
	http.HandleFunc("/getdata", handler.GetSongs)
	http.HandleFunc("/getsongtext", handler.GetSongText)
	// http.HandleFunc("/info", handler.Info)
	err = http.ListenAndServe(a.ip+":"+a.port, nil)
	return err
}
