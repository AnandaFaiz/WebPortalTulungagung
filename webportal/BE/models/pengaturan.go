package models

type IdentitasSitus struct {
	ID        int    `json:"id"`
	NamaSitus string `json:"nama_situs"`
	Tagline   string `json:"tagline"`
	Deskripsi string `json:"deskripsi"`
	Logo      string `json:"logo"`
	Favicon   string `json:"favicon"`
	Email     string `json:"email"`
	Telepon   string `json:"telepon"`
	Alamat    string `json:"alamat"`
}

type Kontak struct {
	ID             int    `json:"id"`
	Alamat         string `json:"alamat"`
	Telepon        string `json:"telepon"`
	Email          string `json:"email"`
	JamOperasional string `json:"jam_operasional"`
	EmbedMaps      string `json:"embed_maps"`
	Facebook       string `json:"facebook"`
	Instagram      string `json:"instagram"`
	Twitter        string `json:"twitter"`
	Youtube        string `json:"youtube"`
}