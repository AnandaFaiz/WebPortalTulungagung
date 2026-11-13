package controllers

import (
	"fmt"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func readMenu(c *gin.Context, publicOnly bool) {
	query := `SELECT id, label, url, COALESCE(deskripsi, ''), COALESCE(icon, ''), COALESCE(color, ''), grup, urutan, is_active FROM menu_navigasi`
	if publicOnly {
		query += ` WHERE is_active = TRUE`
	}
	query += ` ORDER BY grup, urutan, id`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil menu"})
		return
	}
	defer rows.Close()
	items := []models.MenuNavigasi{}
	for rows.Next() {
		var item models.MenuNavigasi
		if err := rows.Scan(&item.ID, &item.Label, &item.URL, &item.Deskripsi, &item.Icon, &item.Color, &item.Grup, &item.Urutan, &item.IsActive); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca menu"})
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan saat membaca menu"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"menu": items})
}

func GetPublicMenu(c *gin.Context) { readMenu(c, true) }
func GetAdminMenu(c *gin.Context)  { readMenu(c, false) }

func CreateMenu(c *gin.Context) {
	var req models.MenuNavigasi
	if err := c.ShouldBindJSON(&req); err != nil || req.Label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Label menu wajib diisi"})
		return
	}
	if req.URL == "" {
		req.URL = "#"
	}
	if req.Grup == "" {
		req.Grup = "Lainnya"
	}
	err := config.DB.QueryRow(`INSERT INTO menu_navigasi (label, url, deskripsi, icon, color, grup, urutan, is_active) VALUES ($1,$2,$3,$4,$5,$6,COALESCE((SELECT MAX(urutan)+1 FROM menu_navigasi WHERE grup=$6),1),$7) RETURNING id, urutan`, req.Label, req.URL, req.Deskripsi, req.Icon, req.Color, req.Grup, req.IsActive).Scan(&req.ID, &req.Urutan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan menu"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "menu": req})
}

func UpdateMenu(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.MenuNavigasi
	if err := c.ShouldBindJSON(&req); err != nil || req.Label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Label menu wajib diisi"})
		return
	}
	_, err := config.DB.Exec(`UPDATE menu_navigasi SET label=$1,url=$2,deskripsi=$3,icon=$4,color=$5,grup=$6,is_active=$7,updated_at=NOW() WHERE id=$8`, req.Label, req.URL, req.Deskripsi, req.Icon, req.Color, req.Grup, req.IsActive, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memperbarui menu"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteMenu(c *gin.Context) {
	if _, err := config.DB.Exec(`DELETE FROM menu_navigasi WHERE id=$1`, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus menu"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ReorderMenu(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data urutan tidak valid"})
		return
	}
	for i, id := range req.IDs {
		config.DB.Exec(`UPDATE menu_navigasi SET urutan = $1, updated_at = NOW() WHERE id = $2`, i+1, id)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func queryPengaduanConfig(table string, activeOnly bool) ([]models.PengaduanKonfigurasi, error) {
	allowedTables := map[string]bool{"pengaduan_kategori": true, "pengaduan_status": true}
	if !allowedTables[table] {
		return nil, fmt.Errorf("tabel konfigurasi tidak valid")
	}
	query := `SELECT id, nama, COALESCE(deskripsi, ''), urutan, is_active FROM ` + table
	if activeOnly {
		query += ` WHERE is_active = TRUE`
	}
	query += ` ORDER BY urutan, id`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.PengaduanKonfigurasi{}
	for rows.Next() {
		var item models.PengaduanKonfigurasi
		if err := rows.Scan(&item.ID, &item.Nama, &item.Deskripsi, &item.Urutan, &item.IsActive); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func GetPublicPengaduanConfig(c *gin.Context) {
	kategori, err := queryPengaduanConfig("pengaduan_kategori", true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil konfigurasi aduan"})
		return
	}
	status, err := queryPengaduanConfig("pengaduan_status", true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil konfigurasi aduan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"kategori": kategori, "status": status})
}

func getAdminPengaduanConfig(c *gin.Context, table string) {
	items, err := queryPengaduanConfig(table, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func GetAdminPengaduanKategori(c *gin.Context) { getAdminPengaduanConfig(c, "pengaduan_kategori") }
func GetAdminPengaduanStatus(c *gin.Context)   { getAdminPengaduanConfig(c, "pengaduan_status") }

func createPengaduanConfig(c *gin.Context, table string) {
	var req models.PengaduanKonfigurasi
	if err := c.ShouldBindJSON(&req); err != nil || req.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama wajib diisi"})
		return
	}
	err := config.DB.QueryRow(`INSERT INTO `+table+` (nama, deskripsi, urutan, is_active) VALUES ($1,$2,COALESCE((SELECT MAX(urutan)+1 FROM `+table+`),1),$3) RETURNING id, urutan`, req.Nama, req.Deskripsi, req.IsActive).Scan(&req.ID, &req.Urutan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama sudah digunakan atau data tidak valid"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "item": req})
}
func updatePengaduanConfig(c *gin.Context, table string) {
	var req models.PengaduanKonfigurasi
	if err := c.ShouldBindJSON(&req); err != nil || req.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama wajib diisi"})
		return
	}
	if _, err := config.DB.Exec(`UPDATE `+table+` SET nama=$1,deskripsi=$2,is_active=$3,updated_at=NOW() WHERE id=$4`, req.Nama, req.Deskripsi, req.IsActive, c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Gagal memperbarui data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
func deletePengaduanConfig(c *gin.Context, table string) {
	if _, err := config.DB.Exec(`DELETE FROM `+table+` WHERE id=$1`, c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
func CreatePengaduanKategori(c *gin.Context) { createPengaduanConfig(c, "pengaduan_kategori") }
func UpdatePengaduanKategori(c *gin.Context) { updatePengaduanConfig(c, "pengaduan_kategori") }
func DeletePengaduanKategori(c *gin.Context) { deletePengaduanConfig(c, "pengaduan_kategori") }
func CreatePengaduanStatus(c *gin.Context)   { createPengaduanConfig(c, "pengaduan_status") }
func UpdatePengaduanStatus(c *gin.Context)   { updatePengaduanConfig(c, "pengaduan_status") }
func DeletePengaduanStatus(c *gin.Context)   { deletePengaduanConfig(c, "pengaduan_status") }
