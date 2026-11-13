package controllers

import (
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/pengumuman
func GetAdminPengumuman(c *gin.Context) {
	query := `SELECT id, judul, COALESCE(konten, ''), COALESCE(gambar, ''), COALESCE(TO_CHAR(tanggal, 'YYYY-MM-DD'), ''), COALESCE(prioritas, 'Normal'), COALESCE(status, 'Draft') FROM pengumuman ORDER BY tanggal DESC NULLS LAST, id DESC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Pengumuman{}
	for rows.Next() {
		var p models.Pengumuman
		if err := rows.Scan(&p.ID, &p.Judul, &p.Konten, &p.Gambar, &p.Tanggal, &p.Prioritas, &p.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pengumuman": list})
}

// GET /api/pengumuman (Publik)
func GetPublicPengumuman(c *gin.Context) {
	query := `SELECT id, judul, COALESCE(konten, ''), COALESCE(gambar, ''), COALESCE(TO_CHAR(tanggal, 'YYYY-MM-DD'), ''), COALESCE(prioritas, 'Normal') FROM pengumuman WHERE status = 'Publish' ORDER BY tanggal DESC LIMIT 5`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Pengumuman{}
	for rows.Next() {
		var p models.Pengumuman
		rows.Scan(&p.ID, &p.Judul, &p.Konten, &p.Gambar, &p.Tanggal, &p.Prioritas)
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pengumuman": list})
}

// POST /api/admin/pengumuman
func CreatePengumuman(c *gin.Context) {
	var req models.Pengumuman
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Judul wajib diisi"})
		return
	}
	if req.Status == "" {
		req.Status = "Draft"
	}
	if req.Prioritas == "" {
		req.Prioritas = "Normal"
	}
	var publishDate interface{}
	if req.Status == "Publish" {
		publishDate = time.Now().Format("2006-01-02")
	}

	err := config.DB.QueryRow(
		"INSERT INTO pengumuman (judul, konten, gambar, tanggal, prioritas, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Judul, req.Konten, req.Gambar, publishDate, req.Prioritas, req.Status,
	).Scan(&req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "pengumuman": req})
}

// PUT /api/admin/pengumuman/:id
func UpdatePengumuman(c *gin.Context) {
	id := c.Param("id")
	var req models.Pengumuman
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	var currentStatus string
	var currentDate string
	if err := config.DB.QueryRow("SELECT COALESCE(status, 'Draft'), COALESCE(TO_CHAR(tanggal, 'YYYY-MM-DD'), '') FROM pengumuman WHERE id = $1", id).Scan(&currentStatus, &currentDate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Pengumuman tidak ditemukan"})
		return
	}
	var publishDate interface{}
	if currentStatus == "Publish" && req.Status == "Publish" && currentDate != "" {
		publishDate = currentDate
	} else if req.Status == "Publish" {
		publishDate = time.Now().Format("2006-01-02")
	}
	_, err := config.DB.Exec(
		"UPDATE pengumuman SET judul=$1, konten=$2, gambar=$3, tanggal=$4, prioritas=$5, status=$6, updated_at=NOW() WHERE id=$7",
		req.Judul, req.Konten, req.Gambar, publishDate, req.Prioritas, req.Status, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/pengumuman/:id
func DeletePengumuman(c *gin.Context) {
	id := c.Param("id")
	res, err := config.DB.Exec("DELETE FROM pengumuman WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus data"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
