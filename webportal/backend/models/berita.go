package models

type Berita struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	Konten    string `json:"konten"`
	Kategori  string `json:"kategori"`
	Penulis   string `json:"penulis"`
	Tanggal   string `json:"tanggal"`
	Status    string `json:"status"`
	Thumbnail string `json:"thumbnail"`
	Views     int    `json:"views"`
}