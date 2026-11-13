package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/berita
func GetAdminBerita(c *gin.Context) {
	// Tambahkan COALESCE untuk setiap kolom yang mungkin NULL
	query := `
		SELECT id, judul, COALESCE(deskripsi, ''), COALESCE(konten, ''), 
		       COALESCE(kategori, 'Umum'), COALESCE(penulis, 'Admin'), 
		       COALESCE(TO_CHAR(tanggal, 'YYYY-MM-DD'), ''), COALESCE(status, 'Draft'), 
		       COALESCE(thumbnail, ''), COALESCE(views, 0) 
		FROM berita 
		ORDER BY tanggal DESC, id DESC`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("❌ Error Query Berita:", err) // Log error ke terminal
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal query ke database"})
		return
	}
	defer rows.Close()

	list := []models.Berita{}
	for rows.Next() {
		var b models.Berita
		err := rows.Scan(&b.ID, &b.Judul, &b.Deskripsi, &b.Konten, &b.Kategori, &b.Penulis, &b.Tanggal, &b.Status, &b.Thumbnail, &b.Views)
		if err != nil {
			log.Println("❌ Error Scan Berita:", err) // Log error ke terminal
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data dari database"})
			return
		}
		list = append(list, b)
	}

	if err := rows.Err(); err != nil {
		log.Println("❌ Error Loop Berita:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan saat membaca data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"berita": list})
}

// GET /api/berita (Publik)
func GetPublicBerita(c *gin.Context) {
	sort := c.DefaultQuery("sort", "terbaru")
	orderBy := "tanggal DESC"
	if sort == "populer" {
		orderBy = "views DESC, tanggal DESC"
	}

	query := `
		SELECT id, judul, COALESCE(deskripsi, ''), COALESCE(konten, ''), COALESCE(kategori, 'Umum'), 
		       COALESCE(penulis, 'Admin'), COALESCE(TO_CHAR(tanggal, 'YYYY-MM-DD'), ''), 
		       COALESCE(thumbnail, ''), COALESCE(views, 0) 
		FROM berita 
		WHERE status = 'Publish' 
		ORDER BY ` + orderBy + ` LIMIT 50`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("❌ Error Query Public Berita:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal query ke database"})
		return
	}
	defer rows.Close()

	list := []models.Berita{}
	for rows.Next() {
		var b models.Berita
		err := rows.Scan(&b.ID, &b.Judul, &b.Deskripsi, &b.Konten, &b.Kategori, &b.Penulis, &b.Tanggal, &b.Thumbnail, &b.Views)
		if err != nil {
			log.Println("❌ Error Scan Public Berita:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data dari database"})
			return
		}
		list = append(list, b)
	}

	if err := rows.Err(); err != nil {
		log.Println("❌ Error Loop Public Berita:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan saat membaca data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"berita": list})
}

// GET /api/berita/:id
func GetPublicBeritaDetail(c *gin.Context) {
	id := c.Param("id")
	var berita models.Berita

	err := config.DB.QueryRow(`
		SELECT id, judul, COALESCE(deskripsi, ''), COALESCE(konten, ''),
		       COALESCE(kategori, 'Umum'), COALESCE(penulis, 'Admin'),
		       COALESCE(TO_CHAR(tanggal, 'YYYY-MM-DD'), ''), COALESCE(status, 'Draft'),
		       COALESCE(thumbnail, ''), COALESCE(views, 0)
		FROM berita
		WHERE id = $1 AND status = 'Publish'`, id).Scan(
		&berita.ID, &berita.Judul, &berita.Deskripsi, &berita.Konten,
		&berita.Kategori, &berita.Penulis, &berita.Tanggal, &berita.Status,
		&berita.Thumbnail, &berita.Views,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Berita tidak ditemukan"})
		return
	}

	config.DB.Exec("UPDATE berita SET views = views + 1 WHERE id = $1", id)
	berita.Views++
	c.JSON(http.StatusOK, gin.H{"berita": berita})
}

// POST /api/admin/berita
func CreateBerita(c *gin.Context) {
	var req models.Berita
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Judul wajib diisi"})
		return
	}

	if req.Status == "" {
		req.Status = "Draft"
	}
	if req.Status == "Publish" {
		req.Tanggal = time.Now().Format("2006-01-02")
	} else {
		req.Tanggal = ""
	}
	var tanggal interface{}
	if req.Tanggal != "" {
		tanggal = req.Tanggal
	}
	userName, err := getCurrentUserName(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Admin login tidak ditemukan"})
		return
	}
	req.Penulis = userName

	err = config.DB.QueryRow(
		"INSERT INTO berita (judul, deskripsi, konten, kategori, penulis, tanggal, status, thumbnail) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		req.Judul, req.Deskripsi, req.Konten, req.Kategori, req.Penulis, tanggal, req.Status, req.Thumbnail,
	).Scan(&req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	LogActivity(c, "create", "berita", req.ID, "Menambah berita: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true, "berita": req})
}

// PUT /api/admin/berita/:id
func UpdateBerita(c *gin.Context) {
	id := c.Param("id")
	var req models.Berita
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	userName, err := getCurrentUserName(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Admin login tidak ditemukan"})
		return
	}
	var currentStatus string
	if err := config.DB.QueryRow("SELECT COALESCE(status, 'Draft') FROM berita WHERE id = $1", id).Scan(&currentStatus); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Berita tidak ditemukan"})
		return
	}
	var publishDate interface{}
	if req.Status == "Publish" {
		if currentStatus == "Publish" {
			publishDate = nil
		} else {
			publishDate = time.Now().Format("2006-01-02")
		}
	}
	if currentStatus == "Publish" && req.Status == "Publish" {
		_, err = config.DB.Exec(
			"UPDATE berita SET judul=$1, deskripsi=$2, konten=$3, kategori=$4, penulis=$5, status=$6, thumbnail=$7, updated_at=NOW() WHERE id=$8",
			req.Judul, req.Deskripsi, req.Konten, req.Kategori, userName, req.Status, req.Thumbnail, id,
		)
	} else {
		_, err = config.DB.Exec(
			"UPDATE berita SET judul=$1, deskripsi=$2, konten=$3, kategori=$4, penulis=$5, tanggal=$6, status=$7, thumbnail=$8, updated_at=NOW() WHERE id=$9",
			req.Judul, req.Deskripsi, req.Konten, req.Kategori, userName, publishDate, req.Status, req.Thumbnail, id,
		)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	LogActivity(c, "update", "berita", 0, "Mengubah berita: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func getCurrentUserName(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", fmt.Errorf("user tidak ditemukan")
	}
	var userName string
	err := config.DB.QueryRow("SELECT nama FROM users WHERE id = $1", userID).Scan(&userName)
	return userName, err
}

// DELETE /api/admin/berita/:id
func DeleteBerita(c *gin.Context) {
	id := c.Param("id")
	res, err := config.DB.Exec("DELETE FROM berita WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus data"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Berita tidak ditemukan"})
		return
	}
	LogActivity(c, "delete", "berita", 0, "Menghapus berita: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// POST /api/berita/:id/view
func AddBeritaView(c *gin.Context) {
	id := c.Param("id")
	config.DB.Exec("UPDATE berita SET views = views + 1 WHERE id = $1", id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
