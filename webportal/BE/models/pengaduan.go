package models

type Pengaduan struct {
	ID           int    `json:"id"`
	Nama         string `json:"nama" binding:"required"`
	Kontak       string `json:"kontak" binding:"required"`
	Kategori     string `json:"kategori" binding:"required"`
	IsiAduan     string `json:"isi_aduan" binding:"required"`
	Status       string `json:"status"`
	BalasanAdmin string `json:"balasan_admin"`
	CreatedAt    string `json:"created_at"`
}