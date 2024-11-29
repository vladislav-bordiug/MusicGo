package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"strings"
)

type DBPool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
}

type PGXDatabase struct {
	pool DBPool
}

func NewPGXDatabase(pool DBPool) *PGXDatabase {
	return &PGXDatabase{pool: pool}
}

func (db *PGXDatabase) CreateTableQuery(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS songs (group_name TEXT, song_name TEXT, releaseDate TEXT, text TEXT, link TEXT);")
	return err
}

func (db *PGXDatabase) InsertQuery(ctx context.Context, group_name string, song_name string, releaseDate string, text string, link string) error {
	_, err := db.pool.Exec(ctx, "INSERT INTO songs(group_name, song_name, releaseDate, text, link) values($1, $2, $3, $4, $5)", group_name, song_name, releaseDate, text, link)
	return err
}

func (db *PGXDatabase) DeleteQuery(ctx context.Context, group_name string, song_name string) error {
	_, err := db.pool.Exec(ctx, "DELETE FROM songs WHERE group_name = $1 AND song_name = $2", group_name, song_name)
	return err
}

func (db *PGXDatabase) SelectDataQuery(ctx context.Context, page int64, items int64, group string, song string, releaseDate string, text string, link string) (AnswerData, error) {
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
	rows, err := db.pool.Query(ctx, query, params...)
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

func (db *PGXDatabase) SelectCoupletQuery(ctx context.Context, group string, song string, couplet int64) (AnswerCoupletData, error) {
	var text string
	var answer AnswerCoupletData
	err := db.pool.QueryRow(ctx, "SELECT text FROM songs WHERE group_name = $1 AND song_name = $2", group, song).Scan(&text)
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

func (db *PGXDatabase) EditQuery(ctx context.Context, group_name string, song_name string, releaseDate string, text string, link string) error {
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
	_, err := db.pool.Exec(ctx, query, params...)
	return err
}
