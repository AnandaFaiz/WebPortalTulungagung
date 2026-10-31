package controllers

import (
	"net/http"
	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

func GetAdminPencarian(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, keyword, urutan FROM pencarian_populer ORDER BY urutan ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, urutan int
		var keyword string
		rows.Scan(&id, &keyword, &urutan)
		list = append(list, map[string]interface{}{"id": id, "keyword": keyword, "urutan": urutan})
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan saat membaca data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pencarian": list})
}

func CreatePencarian(c *gin.Context) {
	var req struct {
		Keyword string `json:"keyword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Keyword wajib diisi"})
		return
	}

	var maxUrutan int
	config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM pencarian_populer").Scan(&maxUrutan)

	_, err := config.DB.Exec("INSERT INTO pencarian_populer (keyword, urutan) VALUES ($1, $2)", req.Keyword, maxUrutan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func UpdatePencarian(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Keyword string `json:"keyword" binding:"required"`
		Urutan  int    `json:"urutan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	_, err := config.DB.Exec("UPDATE pencarian_populer SET keyword=$1, urutan=$2 WHERE id=$3", req.Keyword, req.Urutan, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeletePencarian(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM pencarian_populer WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}