package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/activity-logs
func GetAdminActivityLogs(c *gin.Context) {
	// Ambil parameter page dan perPage
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "15"))
	if page < 1 { page = 1 }
	if perPage < 1 || perPage > 100 { perPage = 15 }
	
	offset := (page - 1) * perPage

	// Hitung total data
	var total int
	err := config.DB.QueryRow(`SELECT COUNT(*) FROM activity_logs`).Scan(&total)
	if err != nil {
		c.JSON(500, gin.H{"statusMessage": "Gagal menghitung data"})
		return
	}

	// Ambil data dengan limit & offset
	query := `SELECT id, COALESCE(user_id, 0), COALESCE(user_nama, ''), COALESCE(user_role, ''), action, COALESCE(entity, ''), COALESCE(description, ''), COALESCE(ip_address, ''), TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS') FROM activity_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := config.DB.Query(query, perPage, offset)
	if err != nil {
		log.Println("Error query logs:", err)
		c.JSON(500, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	type Log struct {
		ID          int    `json:"id"`
		UserID      int    `json:"user_id"`
		UserNama    string `json:"user_nama"`
		UserRole    string `json:"user_role"`
		Action      string `json:"action"`
		Entity      string `json:"entity"`
		Description string `json:"description"`
		IPAddress   string `json:"ip_address"`
		CreatedAt   string `json:"created_at"`
	}

	list := []Log{}
	for rows.Next() {
		var l Log
		if err := rows.Scan(&l.ID, &l.UserID, &l.UserNama, &l.UserRole, &l.Action, &l.Entity, &l.Description, &l.IPAddress, &l.CreatedAt); err != nil {
			c.JSON(500, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, l)
	}
	
	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"statusMessage": "Kesalahan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  list,
		"total": total,
		"page":  page,
		"perPage": perPage,
	})
}