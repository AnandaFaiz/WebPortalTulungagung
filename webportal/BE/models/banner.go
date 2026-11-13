package models

type Banner struct {
	ID     int    `json:"id"`
	Judul  string `json:"judul" `
	Gambar string `json:"gambar" binding:"required"`
	Link   string `json:"link"`
	Urutan int    `json:"urutan"`
	Status string `json:"status"`
	Tipe   string `json:"tipe" binding:"required"`
}