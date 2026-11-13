package models

type StrukturOrganisasi struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul"`
	Deskripsi string `json:"deskripsi"`
	Gambar    string `json:"gambar"`
}