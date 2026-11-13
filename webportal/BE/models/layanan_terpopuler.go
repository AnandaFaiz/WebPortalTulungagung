package models

type LayananTerpopuler struct {
	ID        int    `json:"id"`
	Kategori  string `json:"kategori"`
	Judul     string `json:"judul"`
	Deskripsi string `json:"deskripsi"`
	Icon      string `json:"icon"`
	Link      string `json:"link"`
	Urutan    int    `json:"urutan"`
	IsActive  bool   `json:"is_active"`
}