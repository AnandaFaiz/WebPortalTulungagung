package models

type LayananDarurat struct {
	ID       int    `json:"id"`
	Nama     string `json:"nama" binding:"required"`
	Nomor    string `json:"nomor" binding:"required"`
	Icon     string `json:"icon"`
	Urutan   int    `json:"urutan"`
	IsActive bool   `json:"is_active"`
}