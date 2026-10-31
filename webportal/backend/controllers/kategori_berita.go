package controllers

import (
	"net/http"
	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

func GetAdminKategori(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, nama, COALESCE(icon, ''), COALESCE(deskripsi, '') FROM kategori_berita ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id int
		var nama, icon, deskripsi string
		rows.Scan(&id, &nama, &icon, &deskripsi)
		list = append(list, map[string]interface{}{"id": id, "nama": nama, "icon": icon, "deskripsi": deskripsi})
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan saat membaca data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"kategori": list})
}

func CreateKategori(c *gin.Context) {
	var req struct {
		Nama      string `json:"nama" binding:"required"`
		Icon      string `json:"icon"`
		Deskripsi string `json:"deskripsi"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama wajib diisi"})
		return
	}
	_, err := config.DB.Exec("INSERT INTO kategori_berita (nama, icon, deskripsi) VALUES ($1, $2, $3)", req.Nama, req.Icon, req.Deskripsi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func UpdateKategori(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Nama      string `json:"nama" binding:"required"`
		Icon      string `json:"icon"`
		Deskripsi string `json:"deskripsi"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	_, err := config.DB.Exec("UPDATE kategori_berita SET nama=$1, icon=$2, deskripsi=$3 WHERE id=$4", req.Nama, req.Icon, req.Deskripsi, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteKategori(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM kategori_berita WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}