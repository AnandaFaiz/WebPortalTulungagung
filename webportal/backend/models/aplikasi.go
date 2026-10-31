package models

type Aplikasi struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	Icon      string `json:"icon"`
	Url       string `json:"url"`
	Urutan    int    `json:"urutan"`
	IsActive  bool   `json:"is_active"`
}