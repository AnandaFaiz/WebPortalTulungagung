package models

type MenuNavigasi struct {
	ID        int    `json:"id"`
	Label     string `json:"label" binding:"required"`
	URL       string `json:"url"`
	Deskripsi string `json:"deskripsi"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	Grup      string `json:"grup"`
	Urutan    int    `json:"urutan"`
	IsActive  bool   `json:"is_active"`
}

type PengaduanKonfigurasi struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	Urutan    int    `json:"urutan"`
	IsActive  bool   `json:"is_active"`
}
