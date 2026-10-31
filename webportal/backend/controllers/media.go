package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/media
func GetAdminMedia(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, nama_file, url, COALESCE(tipe_file, ''), COALESCE(ukuran_file, 0), TO_CHAR(created_at, 'YYYY-MM-DD') FROM media ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	type Media struct {
		ID         int    `json:"id"`
		NamaFile   string `json:"nama_file"`
		Url        string `json:"url"`
		TipeFile   string `json:"tipe_file"`
		UkuranFile int    `json:"ukuran_file"`
		CreatedAt  string `json:"created_at"`
	}

	list := []Media{}
	for rows.Next() {
		var m Media
		if err := rows.Scan(&m.ID, &m.NamaFile, &m.Url, &m.TipeFile, &m.UkuranFile, &m.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, m)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"media": list})
}

// DELETE /api/admin/media/:id
func DeleteMedia(c *gin.Context) {
	id := c.Param("id")

	// 1. Ambil nama file dari database
	var fileUrl string
	err := config.DB.QueryRow(`SELECT url FROM media WHERE id = $1`, id).Scan(&fileUrl)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Media tidak ditemukan"})
		return
	}

	// 2. Hapus file fisik dari folder public/uploads
	filename := filepath.Base(fileUrl)
	filePath := filepath.Join("./public/uploads", filename)
	os.Remove(filePath)

	// 3. Hapus record dari database
	config.DB.Exec(`DELETE FROM media WHERE id = $1`, id)

	LogActivity(c, "delete", "media", 0, "Menghapus file media: "+filename)

	c.JSON(http.StatusOK, gin.H{"success": true})
}