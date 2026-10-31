package models

type TautanEksternal struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama" binding:"required"`
	Url       string `json:"url" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	Icon      string `json:"icon"`
	Urutan    int    `json:"urutan"`
	IsActive  bool   `json:"is_active"`
}