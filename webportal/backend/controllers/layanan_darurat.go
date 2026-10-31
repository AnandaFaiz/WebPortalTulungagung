package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

func GetPublicDarurat(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, nama, nomor, COALESCE(icon, ''), urutan FROM layanan_darurat WHERE is_active = TRUE ORDER BY urutan ASC`)
	if err != nil {
		log.Println("Error Query Layanan Darurat Publik:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data layanan darurat"})
		return
	}
	defer rows.Close()
	list := []models.LayananDarurat{}
	for rows.Next() {
		var d models.LayananDarurat
		if err := rows.Scan(&d.ID, &d.Nama, &d.Nomor, &d.Icon, &d.Urutan); err != nil {
			log.Println("Error Scan Layanan Darurat Publik:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data layanan darurat"})
			return
		}
		list = append(list, d)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error Loop Layanan Darurat Publik:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data layanan darurat"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"darurat": list})
}

func GetAdminDarurat(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, nama, nomor, COALESCE(icon, ''), urutan, is_active FROM layanan_darurat ORDER BY urutan ASC`)
	if err != nil {
		log.Println("Error Query Admin Layanan Darurat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data layanan darurat"})
		return
	}
	defer rows.Close()
	list := []models.LayananDarurat{}
	for rows.Next() {
		var d models.LayananDarurat
		if err := rows.Scan(&d.ID, &d.Nama, &d.Nomor, &d.Icon, &d.Urutan, &d.IsActive); err != nil {
			log.Println("Error Scan Admin Layanan Darurat:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data layanan darurat"})
			return
		}
		list = append(list, d)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error Loop Admin Layanan Darurat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data layanan darurat"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"darurat": list})
}

func CreateDarurat(c *gin.Context) {
	var req models.LayananDarurat
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"statusMessage": "Data tidak valid"}); return }
	var max int
	if err := config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM layanan_darurat").Scan(&max); err != nil {
		log.Println("Error Menentukan Urutan Layanan Darurat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan layanan darurat"})
		return
	}
	if err := config.DB.QueryRow("INSERT INTO layanan_darurat (nama, nomor, icon, urutan, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id", req.Nama, req.Nomor, req.Icon, max, req.IsActive).Scan(&req.ID); err != nil {
		log.Println("Error Menyimpan Layanan Darurat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan layanan darurat"})
		return
	}
	LogActivity(c, "create", "layanan_darurat", req.ID, "Menambah layanan darurat: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true, "darurat": req})
}

func UpdateDarurat(c *gin.Context) {
	id := c.Param("id")
	var req models.LayananDarurat
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	if _, err := config.DB.Exec("UPDATE layanan_darurat SET nama=$1, nomor=$2, icon=$3, is_active=$4 WHERE id=$5", req.Nama, req.Nomor, req.Icon, req.IsActive, id); err != nil {
		log.Println("Error Update Layanan Darurat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memperbarui layanan darurat"})
		return
	}
	LogActivity(c, "update", "layanan_darurat", 0, "Mengubah layanan darurat: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteDarurat(c *gin.Context) {
	id := c.Param("id")
	if _, err := config.DB.Exec("DELETE FROM layanan_darurat WHERE id = $1", id); err != nil {
		log.Println("Error Delete Layanan Darurat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus layanan darurat"})
		return
	}
	LogActivity(c, "delete", "layanan_darurat", 0, "Menghapus layanan darurat ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ReorderDarurat(c *gin.Context) {
	var req struct{ IDs []int `json:"ids"` }
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Error Memulai Transaksi Layanan Darurat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengurutkan layanan darurat"})
		return
	}
	for i, id := range req.IDs {
		if _, err := tx.Exec("UPDATE layanan_darurat SET urutan=$1 WHERE id=$2", i+1, id); err != nil {
			tx.Rollback()
			log.Println("Error Mengurutkan Layanan Darurat:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengurutkan layanan darurat"})
			return
		}
	}
	if err := tx.Commit(); err != nil {
		log.Println("Error Commit Urutan Layanan Darurat:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengurutkan layanan darurat"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}