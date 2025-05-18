package models

type Stats struct {
	TotalRead       int `json:"total_read"`
	ReadingProgress int `json:"reading_progress"`
	PopularAuthor   int `json:"popular_author"`
}
