package controllers

import (
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/layanan-terpopuler
func GetAdminLayanan(c *gin.Context) {
	query := `SELECT id, kategori, judul, COALESCE(deskripsi, ''), COALESCE(icon, ''), COALESCE(link, ''), urutan, is_active FROM layanan_terpopuler ORDER BY kategori ASC, urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.LayananTerpopuler{}
	for rows.Next() {
		var l models.LayananTerpopuler
		if err := rows.Scan(&l.ID, &l.Kategori, &l.Judul, &l.Deskripsi, &l.Icon, &l.Link, &l.Urutan, &l.IsActive); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, l)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"layanan": list})
}

// GET /api/layanan-terpopuler (Publik)
func GetPublicLayanan(c *gin.Context) {
	query := `SELECT id, kategori, judul, COALESCE(deskripsi, ''), COALESCE(icon, ''), COALESCE(link, ''), urutan FROM layanan_terpopuler WHERE is_active = TRUE ORDER BY kategori ASC, urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.LayananTerpopuler{}
	for rows.Next() {
		var l models.LayananTerpopuler
		rows.Scan(&l.ID, &l.Kategori, &l.Judul, &l.Deskripsi, &l.Icon, &l.Link, &l.Urutan)
		list = append(list, l)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"layanan": list})
}

// POST /api/admin/layanan-terpopuler
func CreateLayanan(c *gin.Context) {
	var req models.LayananTerpopuler
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	
	// Otomatis set urutan paling akhir berdasarkan kategori
	var maxUrutan int
	config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM layanan_terpopuler WHERE kategori = $1", req.Kategori).Scan(&maxUrutan)
	req.Urutan = maxUrutan

	err := config.DB.QueryRow(
		"INSERT INTO layanan_terpopuler (kategori, judul, deskripsi, icon, link, urutan, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		req.Kategori, req.Judul, req.Deskripsi, req.Icon, req.Link, req.Urutan, req.IsActive,
	).Scan(&req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "layanan": req})
}

// PUT /api/admin/layanan-terpopuler/:id
func UpdateLayanan(c *gin.Context) {
	id := c.Param("id")
	var req models.LayananTerpopuler
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	_, err := config.DB.Exec(
		"UPDATE layanan_terpopuler SET kategori=$1, judul=$2, deskripsi=$3, icon=$4, link=$5, urutan=$6, is_active=$7, updated_at=NOW() WHERE id=$8",
		req.Kategori, req.Judul, req.Deskripsi, req.Icon, req.Link, req.Urutan, req.IsActive, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/layanan-terpopuler/:id
func DeleteLayanan(c *gin.Context) {
	id := c.Param("id")
	res, err := config.DB.Exec("DELETE FROM layanan_terpopuler WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ReorderLayanan(c *gin.Context) {
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
		_, err = tx.Exec("UPDATE layanan_terpopuler SET urutan = $1 WHERE id = $2", i+1, id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan urutan"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}