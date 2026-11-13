package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/opd
func GetAdminOPD(c *gin.Context) {
	query := `SELECT id, nama, COALESCE(deskripsi, ''), COALESCE(website, ''), COALESCE(urutan, 0), COALESCE(is_active, true) FROM opds ORDER BY urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("Error Query OPD:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.OPD{}
	for rows.Next() {
		var o models.OPD
		if err := rows.Scan(&o.ID, &o.Nama, &o.Deskripsi, &o.Website, &o.Urutan, &o.IsActive); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, o)
	}
	// Cek Error Loop
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"opds": list})
}

// GET /api/opd (Publik)
func GetPublicOPD(c *gin.Context) {
	query := `SELECT id, nama, COALESCE(deskripsi, ''), COALESCE(website, ''), COALESCE(urutan, 0) FROM opds WHERE is_active = TRUE ORDER BY urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.OPD{}
	for rows.Next() {
		var o models.OPD
		if err := rows.Scan(&o.ID, &o.Nama, &o.Deskripsi, &o.Website, &o.Urutan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, o)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"opds": list})
}

// POST /api/admin/opd
func CreateOPD(c *gin.Context) {
	var req models.OPD
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama OPD wajib diisi"})
		return
	}

	// Auto-assign urutan
	var maxUrutan int
	config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM opds").Scan(&maxUrutan)
	req.Urutan = maxUrutan

	err := config.DB.QueryRow(
		"INSERT INTO opds (nama, deskripsi, website, urutan, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		req.Nama, req.Deskripsi, req.Website, req.Urutan, req.IsActive,
	).Scan(&req.ID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	LogActivity(c, "create", "opd", req.ID, "Menambah OPD: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true, "opd": req})
}

// PUT /api/admin/opd/:id
func UpdateOPD(c *gin.Context) {
	id := c.Param("id")
	var req models.OPD
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	_, err := config.DB.Exec(
		"UPDATE opds SET nama=$1, deskripsi=$2, website=$3, is_active=$4, updated_at=NOW() WHERE id=$5",
		req.Nama, req.Deskripsi, req.Website, req.IsActive, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	LogActivity(c, "update", "opd", 0, "Mengubah OPD: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/opd/:id
func DeleteOPD(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM opds WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus data"})
		return
	}
	LogActivity(c, "delete", "opd", 0, "Menghapus OPD ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// PUT /api/admin/opd/reorder
func ReorderOPD(c *gin.Context) {
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
		_, err = tx.Exec("UPDATE opds SET urutan = $1 WHERE id = $2", i+1, id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan urutan"})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}