package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

func GetStrukturOrganisasi(c *gin.Context) {
	var s models.StrukturOrganisasi
	query := `SELECT id, COALESCE(judul, ''), COALESCE(deskripsi, ''), COALESCE(gambar, '') FROM struktur_organisasi WHERE id = 1`
	err := config.DB.QueryRow(query).Scan(&s.ID, &s.Judul, &s.Deskripsi, &s.Gambar)
	
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"struktur": models.StrukturOrganisasi{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"struktur": s})
}

func UpdateStrukturOrganisasi(c *gin.Context) {
	var req models.StrukturOrganisasi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	query := `
		INSERT INTO struktur_organisasi (id, judul, deskripsi, gambar, updated_at)
		VALUES (1, $1, $2, $3, NOW())
		ON CONFLICT (id) DO UPDATE SET
			judul = EXCLUDED.judul,
			deskripsi = EXCLUDED.deskripsi,
			gambar = EXCLUDED.gambar,
			updated_at = NOW()`

	_, err := config.DB.Exec(query, req.Judul, req.Deskripsi, req.Gambar)
	if err != nil {
		log.Println("❌ ERROR UPDATE STRUKTUR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}