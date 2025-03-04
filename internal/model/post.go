package model

type Post struct {
	PostNumber int `json:"post_number"`
	Title string `json:"title"`
	Content string `json:"content"`
	Writer string `json:"writer"`
	WriteDate string `json:"write_date"`
	DataType string `json:"data_type"`
}
