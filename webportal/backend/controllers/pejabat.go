package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/pejabat
func GetAdminPejabat(c *gin.Context) {
	query := `SELECT id, nama, jabatan, COALESCE(kategori, ''), COALESCE(foto, ''), COALESCE(periode, ''), COALESCE(urutan, 0), COALESCE(is_active, true) FROM pejabat ORDER BY urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("Error Query Pejabat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Pejabat{}
	for rows.Next() {
		var p models.Pejabat
		if err := rows.Scan(&p.ID, &p.Nama, &p.Jabatan, &p.Kategori, &p.Foto, &p.Periode, &p.Urutan, &p.IsActive); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pejabat": list})
}

// GET /api/pejabat (Publik)
func GetPublicPejabat(c *gin.Context) {
	query := `SELECT id, nama, jabatan, COALESCE(kategori, ''), COALESCE(foto, ''), COALESCE(periode, ''), COALESCE(urutan, 0) FROM pejabat WHERE is_active = TRUE ORDER BY urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Pejabat{}
	for rows.Next() {
		var p models.Pejabat
		if err := rows.Scan(&p.ID, &p.Nama, &p.Jabatan, &p.Kategori, &p.Foto, &p.Periode, &p.Urutan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pejabat": list})
}

// POST /api/admin/pejabat
func CreatePejabat(c *gin.Context) {
	var req models.Pejabat
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama dan Jabatan wajib diisi"})
		return
	}

	var maxUrutan int
	config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM pejabat").Scan(&maxUrutan)
	req.Urutan = maxUrutan

	err := config.DB.QueryRow(
		"INSERT INTO pejabat (nama, jabatan, kategori, foto, periode, urutan, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		req.Nama, req.Jabatan, req.Kategori, req.Foto, req.Periode, req.Urutan, req.IsActive,
	).Scan(&req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	LogActivity(c, "create", "pejabat", req.ID, "Menambah pejabat: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true, "pejabat": req})
}

// PUT /api/admin/pejabat/:id
func UpdatePejabat(c *gin.Context) {
	id := c.Param("id")
	var req models.Pejabat
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	_, err := config.DB.Exec(
		"UPDATE pejabat SET nama=$1, jabatan=$2, kategori=$3, foto=$4, periode=$5, is_active=$6, updated_at=NOW() WHERE id=$7",
		req.Nama, req.Jabatan, req.Kategori, req.Foto, req.Periode, req.IsActive, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	LogActivity(c, "update", "pejabat", 0, "Mengubah pejabat: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/pejabat/:id
func DeletePejabat(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM pejabat WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus data"})
		return
	}
	LogActivity(c, "delete", "pejabat", 0, "Menghapus pejabat ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// PUT /api/admin/pejabat/reorder
func ReorderPejabat(c *gin.Context) {
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
		_, err = tx.Exec("UPDATE pejabat SET urutan = $1 WHERE id = $2", i+1, id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan urutan"})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}