package models

type Galeri struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	Tipe      string `json:"tipe" binding:"required"` // 'foto' atau 'video'
	Tanggal   string `json:"tanggal"`
	Status    string `json:"status"`
	Thumbnail string `json:"thumbnail"`
	FileUrl   string `json:"file_url"`
}