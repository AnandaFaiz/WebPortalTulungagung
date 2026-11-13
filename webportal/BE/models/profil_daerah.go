package models

type Demografi struct {
	Luas      string `json:"luas"`
	Penduduk  string `json:"penduduk"`
	Kecamatan string `json:"kecamatan"`
	Desa      string `json:"desa"`
}

type ProfilDaerah struct {
	ID        int       `json:"id"`
	Sejarah   string    `json:"sejarah"`
	Lambang   string    `json:"lambang"`
	Visi      string    `json:"visi"`
	Misi      []string  `json:"misi"`
	Geografis string    `json:"geografis"`
	Demografi Demografi `json:"demografi"`
}