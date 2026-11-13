package models

type Agenda struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	Tanggal   string `json:"tanggal" binding:"required"`
	Lokasi    string `json:"lokasi"`
	Gambar    string `json:"gambar"`
	Kategori  string `json:"kategori"`
	Status    string `json:"status"`
	LinkEvent string `json:"link_event"`
}