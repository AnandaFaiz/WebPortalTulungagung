package controllers

import (
	"log"
	"net/http"

	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

// GET /api/search?q=...
func GlobalSearch(c *gin.Context) {
	q := c.Query("q")
	if len(q) < 3 {
		c.JSON(http.StatusOK, gin.H{
			"berita":  []interface{}{},
			"halaman": []interface{}{},
			"opd":     []interface{}{},
		})
		return
	}

	searchTerm := "%" + q + "%"
	result := gin.H{
		"berita":  []map[string]interface{}{},
		"halaman": []map[string]interface{}{},
		"opd":     []map[string]interface{}{},
	}

	// 1. Cari di Berita
	beritaRows, err := config.DB.Query(`SELECT id, judul, 'berita' as tipe FROM berita WHERE status = 'Publish' AND LOWER(judul) LIKE LOWER($1) LIMIT 5`, searchTerm)
	if err == nil {
		defer beritaRows.Close()
		beritaList := []map[string]interface{}{}
		for beritaRows.Next() {
			var id int
			var judul, tipe string
			beritaRows.Scan(&id, &judul, &tipe)
			beritaList = append(beritaList, map[string]interface{}{"id": id, "judul": judul, "tipe": tipe})
		}
		if err := beritaRows.Err(); err != nil {
			log.Println("Error Loop Search Berita:", err)
		}
		result["berita"] = beritaList
	} else {
		log.Println("Error Search Berita:", err)
	}

	// 2. Cari di Halaman Statis
	halamanRows, err := config.DB.Query(`SELECT id, judul, slug, 'halaman' as tipe FROM halaman_statis WHERE status = 'Publish' AND (LOWER(judul) LIKE LOWER($1) OR LOWER(slug) LIKE LOWER($1)) LIMIT 5`, searchTerm)
	if err == nil {
		defer halamanRows.Close()
		halamanList := []map[string]interface{}{}
		for halamanRows.Next() {
			var id int
			var judul, slug, tipe string
			halamanRows.Scan(&id, &judul, &slug, &tipe)
			halamanList = append(halamanList, map[string]interface{}{"id": id, "judul": judul, "slug": slug, "tipe": tipe})
		}
		if err := halamanRows.Err(); err != nil {
			log.Println("Error Loop Search Halaman", err)
		}
		result["halaman"] = halamanList
	} else {
		log.Println("Error Search Halaman:", err)
	}

	// 3. Cari di OPD
	opdRows, err := config.DB.Query(`SELECT id, nama as judul, 'opd' as tipe FROM opds WHERE is_active = TRUE AND LOWER(nama) LIKE LOWER($1) LIMIT 5`, searchTerm)
	if err == nil {
		defer opdRows.Close()
		opdList := []map[string]interface{}{}
		for opdRows.Next() {
			var id int
			var judul, tipe string
			opdRows.Scan(&id, &judul, &tipe)
			opdList = append(opdList, map[string]interface{}{"id": id, "judul": judul, "tipe": tipe})
		}
		if err := opdRows.Err(); err != nil {
			log.Println("Error Loop Search OPD", err)
		}
		result["opd"] = opdList
	} else {
		log.Println("Error Search OPD:", err)
	}

	daruratRows, err := config.DB.Query(`SELECT id, nama as judul, nomor as deskripsi, 'darurat' as tipe FROM layanan_darurat WHERE is_active = TRUE AND LOWER(nama) LIKE LOWER($1) LIMIT 5`, searchTerm)
	if err == nil {
		defer daruratRows.Close()
		daruratList := []map[string]interface{}{}
		for daruratRows.Next() {
			var id int
			var judul, deskripsi, tipe string
			daruratRows.Scan(&id, &judul, &deskripsi, &tipe)
			daruratList = append(daruratList, map[string]interface{}{"id": id, "judul": judul, "deskripsi": deskripsi, "tipe": tipe})
		}
		if err := daruratRows.Err(); err != nil {
			log.Println("Error Loop Search Darurat:", err)
		}
		result["darurat"] = daruratList
	} else {
		log.Println("Error Query Search Darurat:", err)
	}

	aplikasiRows, err := config.DB.Query(`SELECT id, nama as judul, 'aplikasi' as tipe, url FROM aplikasi WHERE is_active = TRUE AND LOWER(nama) LIKE LOWER($1) LIMIT 5`, searchTerm)
	if err == nil {
		defer aplikasiRows.Close()
		aplikasiList := []map[string]interface{}{}
		for aplikasiRows.Next() {
			var id int
			var judul, tipe, url string
			aplikasiRows.Scan(&id, &judul, &tipe, &url)
			aplikasiList = append(aplikasiList, map[string]interface{}{"id": id, "judul": judul, "tipe": tipe, "url": url})
		}
		if err := aplikasiRows.Err(); err != nil {
			log.Println("Error Loop Search Aplikasi:", err)
		}
		result["aplikasi"] = aplikasiList
	} else {
		log.Println("Error Query Search Aplikasi:", err)
	}

	c.JSON(http.StatusOK, result)
}

// GET /api/pencarian-populer
func GetPencarianPopuler(c *gin.Context) {
	rows, err := config.DB.Query("SELECT keyword FROM pencarian_populer ORDER BY urutan ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	keywords := []string{}
	for rows.Next() {
		var kw string
		rows.Scan(&kw)
		keywords = append(keywords, kw)
	}
	if err := rows.Err(); err != nil{
			log.Println("Error Loop Pencarian Populer", err)
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan saat membaca data"})
			return
		}
	c.JSON(http.StatusOK, gin.H{"keywords": keywords})
}

// GET /api/kategori-berita
func GetKategoriBerita(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, nama, COALESCE(icon, 'i-tabler-news'), COALESCE(deskripsi, '') FROM kategori_berita ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	type Kategori struct {
		ID        int    `json:"id"`
		Nama      string `json:"nama"`
		Icon      string `json:"icon"`
		Deskripsi string `json:"deskripsi"`
	}

	list := []Kategori{}
	for rows.Next() {
		var k Kategori
		rows.Scan(&k.ID, &k.Nama, &k.Icon, &k.Deskripsi)
		list = append(list, k)
	}
	if err := rows.Err(); err != nil{
			log.Println("Error Loop Pencarian Berita", err)
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan saat membaca data"})
			return
		}
	c.JSON(http.StatusOK, gin.H{"kategori": list})
}