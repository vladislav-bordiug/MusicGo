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
	_, err := db.pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS groups (id SERIAL PRIMARY KEY, group_name TEXT, CONSTRAINT unique_group UNIQUE(group_name));")
	if err != nil {
		return err
	}
	_, err = db.pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS songs (id SERIAL PRIMARY KEY, song_name TEXT, releaseDate TIMESTAMP, text TEXT, link TEXT, group_id INTEGER, FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE, CONSTRAINT unique_group_song UNIQUE(group_id, song_name));")
	return err
}

func (db *PGXDatabase) InsertQuery(ctx context.Context, group_name string, song_name string, releaseDate string, text string, link string) error {
	var groupID int
	err := db.pool.QueryRow(ctx, "SELECT id FROM groups WHERE group_name = $1", group_name).Scan(&groupID)
	if err != nil {
		err = db.pool.QueryRow(ctx, "INSERT INTO groups(group_name) values($1) RETURNING id", group_name).Scan(&groupID)
		if err != nil {
			return err
		}
	}
	_, err = db.pool.Exec(ctx, "INSERT INTO songs(song_name, releaseDate, text, link, group_id) values($1, TO_TIMESTAMP($2, 'DD.MM.YYYY'), $3, $4, $5)", song_name, releaseDate, text, link, groupID)
	return err
}

func (db *PGXDatabase) DeleteQuery(ctx context.Context, group_name string, song_name string) error {
	groupID, err := db.SelectGroupIdQuery(ctx, group_name)
	if err != nil {
		return err
	}
	_, err = db.pool.Exec(ctx, "DELETE FROM songs WHERE group_id = $1 AND song_name = $2", groupID, song_name)
	return err
}

func (db *PGXDatabase) SelectGroupIdQuery(ctx context.Context, group_name string) (int, error) {
	var groupID int
	err := db.pool.QueryRow(ctx, "SELECT id FROM groups WHERE group_name = $1", group_name).Scan(&groupID)
	if err != nil {
		return 0, err
	}
	return groupID, nil
}

func (db *PGXDatabase) SelectDataQuery(ctx context.Context, page int64, items int64, group string, song string, releaseDate string, text string, link string) (AnswerData, error) {
	query := "SELECT g.group_name, s.song_name, TO_CHAR(s.releaseDate, 'DD.MM.YYYY'), s.text, s.link FROM songs s JOIN groups g ON s.group_id = g.id "
	var answer AnswerData
	paramindex := 1
	setClauses := []string{}
	params := []interface{}{}
	if group != "" || song != "" || releaseDate != "" || text != "" || link != "" {
		if group != "" {
			var groupID int
			groupID, err := db.SelectGroupIdQuery(ctx, group)
			if err != nil {
				return answer, err
			}
			setClauses = append(setClauses, fmt.Sprintf("group_id = $%d", paramindex))
			params = append(params, groupID)
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
	var groupID int
	groupID, err := db.SelectGroupIdQuery(ctx, group)
	if err != nil {
		return answer, err
	}
	err = db.pool.QueryRow(ctx, "SELECT text FROM songs WHERE group_id = $1 AND song_name = $2", groupID, song).Scan(&text)
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
	var groupID int
	groupID, err := db.SelectGroupIdQuery(ctx, group_name)
	if err != nil {
		return err
	}
	params := []interface{}{groupID, song_name}
	if releaseDate != "" {
		setClauses = append(setClauses, fmt.Sprintf("releaseDate = TO_TIMESTAMP($%d, 'DD.MM.YYYY')", paramindex))
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
	query += fmt.Sprintf(" WHERE group_id = $1 AND song_name = $2")
	log.Printf("INFO: query for the database=%s\n", query)
	log.Printf("INFO: query params for the database=%s\n", params)
	_, err = db.pool.Exec(ctx, query, params...)
	return err
}
