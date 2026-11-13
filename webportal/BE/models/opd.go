package models

type OPD struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	Website   string `json:"website"`
	Urutan    int    `json:"urutan"`
	IsActive  bool   `json:"is_active"`
}