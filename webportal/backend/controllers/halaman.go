package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAdminHalaman(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, judul, slug, COALESCE(konten, ''), TO_CHAR(tanggal, 'YYYY-MM-DD'), COALESCE(status, 'Draft') FROM halaman_statis ORDER BY tanggal DESC`)
	if err != nil { c.JSON(500, gin.H{"statusMessage": "Gagal data"}); return }
	defer rows.Close()
	
	list := []models.HalamanStatis{}
	for rows.Next() {
		var h models.HalamanStatis
		rows.Scan(&h.ID, &h.Judul, &h.Slug, &h.Konten, &h.Tanggal, &h.Status)
		list = append(list, h)
	}
	if err := rows.Err(); err != nil { log.Println(err); c.JSON(500, gin.H{"statusMessage": "Kesalahan data"}); return }
	c.JSON(http.StatusOK, gin.H{"halaman": list})
}

func CreateHalaman(c *gin.Context) {
	var req models.HalamanStatis
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"statusMessage": "Data tidak valid"}); return }
	if req.Status == "" { req.Status = "Draft" }
	
	err := config.DB.QueryRow(`INSERT INTO halaman_statis (judul, slug, konten, tanggal, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`, req.Judul, req.Slug, req.Konten, req.Tanggal, req.Status).Scan(&req.ID)
	if err != nil { c.JSON(500, gin.H{"statusMessage": "Gagal simpan (Cek duplikasi slug)"}); return }
	LogActivity(c, "create", "halaman_statis", req.ID, "Menambah halaman: "+req.Judul)
	c.JSON(200, gin.H{"success": true, "halaman": req})
}

func UpdateHalaman(c *gin.Context) {
	id := c.Param("id")
	entityID, err := strconv.Atoi(id)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "ID halaman tidak valid"}); return }
	var req models.HalamanStatis
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"statusMessage": "Data invalid"}); return }
	_, err = config.DB.Exec(`UPDATE halaman_statis SET judul=$1, slug=$2, konten=$3, tanggal=$4, status=$5, updated_at=NOW() WHERE id=$6`, req.Judul, req.Slug, req.Konten, req.Tanggal, req.Status, entityID)
	if err != nil { c.JSON(500, gin.H{"statusMessage": "Gagal update"}); return }
	LogActivity(c, "update", "halaman_statis", entityID, "Mengubah halaman: "+req.Judul)
	c.JSON(200, gin.H{"success": true})
}

func DeleteHalaman(c *gin.Context) {
	id := c.Param("id")
	entityID, err := strconv.Atoi(id)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "ID halaman tidak valid"}); return }
	if _, err := config.DB.Exec(`DELETE FROM halaman_statis WHERE id=$1`, entityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus halaman"})
		return
	}
	LogActivity(c, "delete", "halaman_statis", entityID, "Menghapus halaman ID: "+id)
	c.JSON(200, gin.H{"success": true})
}