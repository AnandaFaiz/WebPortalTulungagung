package controllers

import (
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/banner
func GetAdminBanner(c *gin.Context) {
	query := `SELECT id, COALESCE(judul, ''), COALESCE(gambar, ''), COALESCE(link, ''), COALESCE(urutan, 0), COALESCE(status, 'Aktif'), COALESCE(tipe, 'beranda') FROM banner ORDER BY urutan ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Banner{}
	for rows.Next() {
		var b models.Banner
		if err := rows.Scan(&b.ID, &b.Judul, &b.Gambar, &b.Link, &b.Urutan, &b.Status, &b.Tipe); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, b)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"banner": list})
}

// GET /api/banner (Publik)
func GetPublicBanner(c *gin.Context) {
	tipe := c.DefaultQuery("tipe", "beranda")
	query := `SELECT id, COALESCE(judul, ''), COALESCE(gambar, ''), COALESCE(link, ''), COALESCE(urutan, 0) FROM banner WHERE status = 'Aktif' AND tipe = $1 ORDER BY urutan ASC`
	rows, err := config.DB.Query(query, tipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Banner{}
	for rows.Next() {
		var b models.Banner
		rows.Scan(&b.ID, &b.Judul, &b.Gambar, &b.Link, &b.Urutan)
		list = append(list, b)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"banner": list})
}

// POST /api/admin/banner
func CreateBanner(c *gin.Context) {
	var req models.Banner
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	if req.Status == "" { req.Status = "Aktif" }
	if req.Tipe == "" { req.Tipe = "beranda" }

	var maxUrutan int
	config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM banner").Scan(&maxUrutan)
	req.Urutan = maxUrutan

	err := config.DB.QueryRow(
		"INSERT INTO banner (judul, gambar, link, urutan, status, tipe) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Judul, req.Gambar, req.Link, req.Urutan, req.Status, req.Tipe,
	).Scan(&req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	LogActivity(c, "create", "banner", req.ID, "Menambah banner: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true, "banner": req})
}

// PUT /api/admin/banner/:id
func UpdateBanner(c *gin.Context) {
	id := c.Param("id")
	var req models.Banner
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	_, err := config.DB.Exec(
		"UPDATE banner SET judul=$1, gambar=$2, link=$3, urutan=$4, status=$5, tipe=$6, updated_at=NOW() WHERE id=$7",
		req.Judul, req.Gambar, req.Link, req.Urutan, req.Status, req.Tipe, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	LogActivity(c, "update", "banner", 0, "Mengubah banner: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/banner/:id
func DeleteBanner(c *gin.Context) {
	id := c.Param("id")
	config.DB.Exec("DELETE FROM banner WHERE id = $1", id)
	LogActivity(c, "delete", "banner", 0, "Menghapus banner ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// PUT /api/admin/banner/reorder
func ReorderBanner(c *gin.Context) {
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
		_, err = tx.Exec("UPDATE banner SET urutan = $1 WHERE id = $2", i+1, id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan urutan"})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}