package models

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
