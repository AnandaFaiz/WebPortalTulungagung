package models

type InformasiMPP struct {
	ID             int    `json:"id"`
	Deskripsi      string `json:"deskripsi"`
	Alamat         string `json:"alamat"`
	JamOperasional string `json:"jam_operasional"`
	JumlahLayanan  int    `json:"jumlah_layanan"`
	JumlahInstansi int    `json:"jumlah_instansi"`
}