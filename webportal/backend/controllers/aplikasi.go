package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/aplikasi
func GetAdminAplikasi(c *gin.Context) {
	query := `SELECT id, nama, COALESCE(deskripsi, ''), COALESCE(icon, ''), COALESCE(url, ''), COALESCE(urutan, 0), COALESCE(is_active, true) FROM aplikasi ORDER BY urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("Error Query Aplikasi:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Aplikasi{}
	for rows.Next() {
		var a models.Aplikasi
		if err := rows.Scan(&a.ID, &a.Nama, &a.Deskripsi, &a.Icon, &a.Url, &a.Urutan, &a.IsActive); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"aplikasi": list})
}

// GET /api/aplikasi (Publik)
func GetPublicAplikasi(c *gin.Context) {
	query := `SELECT id, nama, COALESCE(deskripsi, ''), COALESCE(icon, ''), COALESCE(url, ''), COALESCE(urutan, 0) FROM aplikasi WHERE is_active = TRUE ORDER BY urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Aplikasi{}
	for rows.Next() {
		var a models.Aplikasi
		if err := rows.Scan(&a.ID, &a.Nama, &a.Deskripsi, &a.Icon, &a.Url, &a.Urutan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"aplikasi": list})
}

// POST /api/admin/aplikasi
func CreateAplikasi(c *gin.Context) {
	var req models.Aplikasi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama aplikasi wajib diisi"})
		return
	}

	var maxUrutan int
	config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM aplikasi").Scan(&maxUrutan)
	req.Urutan = maxUrutan

	err := config.DB.QueryRow(
		"INSERT INTO aplikasi (nama, deskripsi, icon, url, urutan, is_active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Nama, req.Deskripsi, req.Icon, req.Url, req.Urutan, req.IsActive,
	).Scan(&req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "aplikasi": req})
}

// PUT /api/admin/aplikasi/:id
func UpdateAplikasi(c *gin.Context) {
	id := c.Param("id")
	var req models.Aplikasi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	_, err := config.DB.Exec(
		"UPDATE aplikasi SET nama=$1, deskripsi=$2, icon=$3, url=$4, is_active=$5, updated_at=NOW() WHERE id=$6",
		req.Nama, req.Deskripsi, req.Icon, req.Url, req.IsActive, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/aplikasi/:id
func DeleteAplikasi(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM aplikasi WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// PUT /api/admin/aplikasi/reorder
func ReorderAplikasi(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Array ids wajib diisi"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memulai transaksi"})
		return
	}

	for i, id := range req.IDs {
		_, err = tx.Exec("UPDATE aplikasi SET urutan = $1 WHERE id = $2", i+1, id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan urutan"})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}