package models

type Pengumuman struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul" binding:"required"`
	Konten    string `json:"konten"`
	Tanggal   string `json:"tanggal"`
	Prioritas string `json:"prioritas"`
	Status    string `json:"status"`
}