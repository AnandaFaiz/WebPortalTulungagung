package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GANTI Fungsi GetAdminAgenda
func GetAdminAgenda(c *gin.Context) {
	query := `SELECT id, judul, COALESCE(deskripsi, ''), TO_CHAR(tanggal, 'YYYY-MM-DD"T"HH24:MI:SS'), COALESCE(lokasi, ''), COALESCE(gambar, ''), COALESCE(kategori, ''), COALESCE(status, 'Draft'), COALESCE(link_event, '') FROM agenda ORDER BY tanggal ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Agenda{}
	for rows.Next() {
		var a models.Agenda
		if err := rows.Scan(&a.ID, &a.Judul, &a.Deskripsi, &a.Tanggal, &a.Lokasi, &a.Gambar, &a.Kategori, &a.Status, &a.LinkEvent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agenda": list})
}

// GANTI Fungsi GetPublicAgenda
func GetPublicAgenda(c *gin.Context) {
	query := `SELECT id, judul, COALESCE(deskripsi, ''), TO_CHAR(tanggal, 'YYYY-MM-DD"T"HH24:MI:SS'), COALESCE(lokasi, ''), COALESCE(gambar, ''), COALESCE(kategori, ''), COALESCE(link_event, '') FROM agenda WHERE status = 'Publish' AND tanggal >= NOW() ORDER BY tanggal ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Agenda{}
	for rows.Next() {
		var a models.Agenda
		if err := rows.Scan(&a.ID, &a.Judul, &a.Deskripsi, &a.Tanggal, &a.Lokasi, &a.Gambar, &a.Kategori, &a.LinkEvent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agenda": list})
}

// GANTI Fungsi CreateAgenda
func CreateAgenda(c *gin.Context) {
	var req models.Agenda
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Judul dan Tanggal wajib diisi"})
		return
	}
	if req.Status == "" { req.Status = "Draft" }

	err := config.DB.QueryRow(
		"INSERT INTO agenda (judul, deskripsi, tanggal, lokasi, gambar, kategori, status, link_event) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		req.Judul, req.Deskripsi, req.Tanggal, req.Lokasi, req.Gambar, req.Kategori, req.Status, req.LinkEvent,
	).Scan(&req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	LogActivity(c, "create", "event", req.ID, "Menambah event: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true, "agenda": req})
}

// GANTI Fungsi UpdateAgenda
func UpdateAgenda(c *gin.Context) {
	id := c.Param("id")
	var req models.Agenda
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	_, err := config.DB.Exec(
		"UPDATE agenda SET judul=$1, deskripsi=$2, tanggal=$3, lokasi=$4, gambar=$5, kategori=$6, status=$7, link_event=$8, updated_at=NOW() WHERE id=$9",
		req.Judul, req.Deskripsi, req.Tanggal, req.Lokasi, req.Gambar, req.Kategori, req.Status, req.LinkEvent, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	LogActivity(c, "update", "event", 0, "Mengubah event: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/agenda/:id
func DeleteAgenda(c *gin.Context) {
	id := c.Param("id")
	res, err := config.DB.Exec("DELETE FROM agenda WHERE id = $1", id)
	if err != nil {
		log.Println("Error Delete Agenda:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus data"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Agenda tidak ditemukan"})
		return
	}
	LogActivity(c, "delete", "event", 0, "Menghapus event ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}