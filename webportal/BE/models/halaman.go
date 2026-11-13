package models

type HalamanStatis struct {
	ID      int    `json:"id"`
	Judul   string `json:"judul" binding:"required"`
	Slug    string `json:"slug" binding:"required"`
	Konten  string `json:"konten"`
	Tanggal string `json:"tanggal"`
	Status  string `json:"status"`
}