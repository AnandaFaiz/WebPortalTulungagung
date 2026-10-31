package models

type Pejabat struct {
	ID       int    `json:"id"`
	Nama     string `json:"nama" binding:"required"`
	Jabatan  string `json:"jabatan" binding:"required"`
	Kategori string `json:"kategori"`
	Foto     string `json:"foto"`
	Periode  string `json:"periode"`
	Urutan   int    `json:"urutan"`
	IsActive bool   `json:"is_active"`
}