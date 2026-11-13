package models

type Layanan struct {
	ID         int    `json:"id"`
	KategoriID int    `json:"kategori_id"`
	Nama       string `json:"nama" binding:"required"`
	Link       string `json:"link"`
	Urutan     int    `json:"urutan"`
}

type KategoriLayanan struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama" binding:"required"`
	Deskripsi string    `json:"deskripsi"`
	Icon      string    `json:"icon"`
	Urutan    int       `json:"urutan"`
	Items     []Layanan `json:"items"` 
}