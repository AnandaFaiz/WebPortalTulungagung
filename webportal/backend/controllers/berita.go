package controllers

import (
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
		       TO_CHAR(tanggal, 'YYYY-MM-DD'), COALESCE(status, 'Draft'), 
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
		SELECT id, judul, COALESCE(deskripsi, ''), COALESCE(kategori, 'Umum'), 
		       COALESCE(penulis, 'Admin'), TO_CHAR(tanggal, 'YYYY-MM-DD'), 
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
		err := rows.Scan(&b.ID, &b.Judul, &b.Deskripsi, &b.Kategori, &b.Penulis, &b.Tanggal, &b.Thumbnail, &b.Views)
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

// POST /api/admin/berita
func CreateBerita(c *gin.Context) {
	var req models.Berita
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Judul wajib diisi"})
		return
	}
	
	if req.Tanggal == "" {
		req.Tanggal = time.Now().Format("2006-01-02")
	}
	if req.Status == "" {
		req.Status = "Draft"
	}

	err := config.DB.QueryRow(
		"INSERT INTO berita (judul, deskripsi, konten, kategori, penulis, tanggal, status, thumbnail) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		req.Judul, req.Deskripsi, req.Konten, req.Kategori, req.Penulis, req.Tanggal, req.Status, req.Thumbnail,
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

	_, err := config.DB.Exec(
		"UPDATE berita SET judul=$1, deskripsi=$2, konten=$3, kategori=$4, penulis=$5, tanggal=$6, status=$7, thumbnail=$8, updated_at=NOW() WHERE id=$9",
		req.Judul, req.Deskripsi, req.Konten, req.Kategori, req.Penulis, req.Tanggal, req.Status, req.Thumbnail, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update data"})
		return
	}
	LogActivity(c, "update", "berita", 0, "Mengubah berita: "+req.Judul)
	c.JSON(http.StatusOK, gin.H{"success": true})
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