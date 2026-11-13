package controllers

import (
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/galeri
func GetAdminGaleri(c *gin.Context) {
	query := `SELECT id, judul, COALESCE(deskripsi, ''), COALESCE(tipe, 'foto'), TO_CHAR(tanggal, 'YYYY-MM-DD'), COALESCE(status, 'Draft'), COALESCE(thumbnail, ''), COALESCE(file_url, '') FROM galeri ORDER BY tanggal DESC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Galeri{}
	for rows.Next() {
		var g models.Galeri
		if err := rows.Scan(&g.ID, &g.Judul, &g.Deskripsi, &g.Tipe, &g.Tanggal, &g.Status, &g.Thumbnail, &g.FileUrl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, g)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"galeri": list})
}

// GET /api/galeri (Publik)
func GetPublicGaleri(c *gin.Context) {
	query := `SELECT id, judul, COALESCE(deskripsi, ''), COALESCE(tipe, 'foto'), TO_CHAR(tanggal, 'YYYY-MM-DD'), COALESCE(thumbnail, '') FROM galeri WHERE status = 'Publish' ORDER BY tanggal DESC LIMIT 20`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Galeri{}
	for rows.Next() {
		var g models.Galeri
		rows.Scan(&g.ID, &g.Judul, &g.Deskripsi, &g.Tipe, &g.Tanggal, &g.Thumbnail)
		list = append(list, g)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"galeri": list})
}

// POST /api/admin/galeri
func CreateGaleri(c *gin.Context) {
	var req models.Galeri
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Judul wajib diisi"})
		return
	}
	if req.Status == "" { req.Status = "Draft" }

	err := config.DB.QueryRow(
		"INSERT INTO galeri (judul, deskripsi, tipe, tanggal, status, thumbnail, file_url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		req.Judul, req.Deskripsi, req.Tipe, req.Tanggal, req.Status, req.Thumbnail, req.FileUrl,
	).Scan(&req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	LogActivity(c, "create", "galeri", req.ID, "Menambah galeri: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true, "galeri": req})
}

// PUT /api/admin/galeri/:id
func UpdateGaleri(c *gin.Context) {
	id := c.Param("id")
	var req models.Galeri
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	_, err := config.DB.Exec(
		"UPDATE galeri SET judul=$1, deskripsi=$2, tipe=$3, tanggal=$4, status=$5, thumbnail=$6, file_url=$7, updated_at=NOW() WHERE id=$8",
		req.Judul, req.Deskripsi, req.Tipe, req.Tanggal, req.Status, req.Thumbnail, req.FileUrl, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	LogActivity(c, "update", "galeri", 0, "Mengubah galeri: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/galeri/:id
func DeleteGaleri(c *gin.Context) {
	id := c.Param("id")
	res, err := config.DB.Exec("DELETE FROM galeri WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus data"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Data tidak ditemukan"})
		return
	}
	LogActivity(c, "delete", "galeri", 0, "Menghapus galeri ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}