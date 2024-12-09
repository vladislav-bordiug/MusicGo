package database

import (
	"context"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTableQuery(t *testing.T) {
	mockk, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	database := NewPGXDatabase(mockk)
	defer mockk.Close()
	mockk.ExpectExec("CREATE TABLE IF NOT EXISTS groups").WillReturnResult(pgxmock.NewResult("CREATE", 1))
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
	mockk.ExpectQuery("INSERT INTO groups").
		WithArgs(group).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mockk.ExpectExec("INSERT INTO songs").
		WithArgs(song, date, text, link, 1).
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
	mockk.ExpectQuery("SELECT id FROM groups").
		WithArgs(group).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mockk.ExpectExec("DELETE FROM songs").
		WithArgs(1, song).
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
	mockk.ExpectQuery("SELECT id FROM groups").
		WithArgs(group).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mockk.ExpectQuery("SELECT g.group_name, s.song_name, TO_CHAR\\(s.releaseDate, \\'DD.MM.YYYY\\'\\), s.text, s.link FROM songs s JOIN groups g ON s.group_id = g.id").
		WithArgs(1, song, date, text, link, items, (page-1)*items).
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
	mockk.ExpectQuery("SELECT id FROM groups").
		WithArgs(group).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mockk.ExpectQuery("SELECT text FROM songs").
		WithArgs(1, song).
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
	mockk.ExpectQuery("SELECT id FROM groups").
		WithArgs(group).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
	mockk.ExpectExec("UPDATE songs SET").
		WithArgs(1, song, date, text, link).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	err = database.EditQuery(context.Background(), group, song, date, text, link)
	assert.NoError(t, err)
	if err := mockk.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
