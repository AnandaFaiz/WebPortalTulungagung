package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GET /api/mpp
func GetMPP(c *gin.Context) {
	var m models.InformasiMPP
	query := `SELECT id, COALESCE(deskripsi, ''), COALESCE(alamat, ''), COALESCE(jam_operasional, ''), COALESCE(jumlah_layanan, 0), COALESCE(jumlah_instansi, 0) FROM informasi_mpp WHERE id = 1`
	
	err := config.DB.QueryRow(query).Scan(&m.ID, &m.Deskripsi, &m.Alamat, &m.JamOperasional, &m.JumlahLayanan, &m.JumlahInstansi)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"mpp": models.InformasiMPP{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mpp": m})
}

// PUT /api/admin/mpp
func UpdateMPP(c *gin.Context) {
	var req models.InformasiMPP
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	query := `
		INSERT INTO informasi_mpp (id, deskripsi, alamat, jam_operasional, jumlah_layanan, jumlah_instansi, updated_at)
		VALUES (1, $1, $2, $3, $4, $5, NOW())
		ON CONFLICT (id) DO UPDATE SET
			deskripsi = EXCLUDED.deskripsi,
			alamat = EXCLUDED.alamat,
			jam_operasional = EXCLUDED.jam_operasional,
			jumlah_layanan = EXCLUDED.jumlah_layanan,
			jumlah_instansi = EXCLUDED.jumlah_instansi,
			updated_at = NOW()`

	_, err := config.DB.Exec(query, req.Deskripsi, req.Alamat, req.JamOperasional, req.JumlahLayanan, req.JumlahInstansi)
	if err != nil {
		log.Println("❌ ERROR UPDATE MPP:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}

	// Catat Log
	LogActivity(c, "update", "informasi_mpp", 1, "Memperbarui Informasi MPP")

	c.JSON(http.StatusOK, gin.H{"success": true})
}