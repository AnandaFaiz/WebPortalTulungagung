package controllers

import (
	"log"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

func GetAdminTautan(c *gin.Context) {
	rows, _ := config.DB.Query(`SELECT id, nama, url, COALESCE(deskripsi, ''), COALESCE(icon, ''), COALESCE(urutan, 0), COALESCE(is_active, true) FROM tautan_eksternal ORDER BY urutan ASC`)
	defer rows.Close()
	list := []models.TautanEksternal{}
	for rows.Next() {
		var t models.TautanEksternal
		rows.Scan(&t.ID, &t.Nama, &t.Url, &t.Deskripsi, &t.Icon, &t.Urutan, &t.IsActive)
		list = append(list, t)
	}
	if err := rows.Err(); err != nil { log.Println(err); c.JSON(500, gin.H{"statusMessage": "Kesalahan data"}); return }
	c.JSON(200, gin.H{"tautan": list})
}

func CreateTautan(c *gin.Context) {
	var req models.TautanEksternal
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"statusMessage": "Data invalid"}); return }
	var max int
	config.DB.QueryRow(`SELECT COALESCE(MAX(urutan), 0) + 1 FROM tautan_eksternal`).Scan(&max)
	req.Urutan = max
	config.DB.QueryRow(`INSERT INTO tautan_eksternal (nama, url, deskripsi, icon, urutan, is_active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, req.Nama, req.Url, req.Deskripsi, req.Icon, req.Urutan, req.IsActive).Scan(&req.ID)
	LogActivity(c, "create", "tautan_eksternal", req.ID, "Menambah tautan: "+req.Nama)
	c.JSON(200, gin.H{"success": true, "tautan": req})
}

func UpdateTautan(c *gin.Context) {
	id := c.Param("id")
	var req models.TautanEksternal
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"statusMessage": "Data invalid"}); return }
	config.DB.Exec(`UPDATE tautan_eksternal SET nama=$1, url=$2, deskripsi=$3, icon=$4, is_active=$5 WHERE id=$6`, req.Nama, req.Url, req.Deskripsi, req.Icon, req.IsActive, id)
	LogActivity(c, "update", "tautan_eksternal", 0, "Mengubah tautan: "+req.Nama)
	c.JSON(200, gin.H{"success": true})
}

func DeleteTautan(c *gin.Context) {
	id := c.Param("id")
	config.DB.Exec(`DELETE FROM tautan_eksternal WHERE id=$1`, id)
	LogActivity(c, "delete", "tautan_eksternal", 0, "Menghapus tautan ID: "+id)
	c.JSON(200, gin.H{"success": true})
}

func ReorderTautan(c *gin.Context) {
	var req struct{ IDs []int `json:"ids"` }
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 { c.JSON(400, gin.H{"statusMessage": "Data invalid"}); return }
	tx, _ := config.DB.Begin()
	for i, id := range req.IDs { tx.Exec(`UPDATE tautan_eksternal SET urutan=$1 WHERE id=$2`, i+1, id) }
	tx.Commit()
	c.JSON(200, gin.H{"success": true})
}