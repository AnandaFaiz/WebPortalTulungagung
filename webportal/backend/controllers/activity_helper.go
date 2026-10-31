package controllers

import (
	"log"
	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

// Fungsi pembantu untuk mencatat aktivitas
func LogActivity(c *gin.Context, action string, entity string, entityID int, description string) {
	// Ambil ID user dari context (yang di-set oleh AuthMiddleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return // Jika tidak ada user (misal gagal login), jangan catat
	}
	userID := userIDVal.(int)

	userRoleVal, _ := c.Get("user_role")
	userRole := userRoleVal.(string)

	// Ambil nama user dari database
	var userName string
	config.DB.QueryRow("SELECT nama FROM users WHERE id = $1", userID).Scan(&userName)

	// Ambil IP Address
	ip := c.ClientIP()

	// Simpan ke database
	_, err := config.DB.Exec(
		`INSERT INTO activity_logs (user_id, user_nama, user_role, action, entity, entity_id, description, ip_address)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, userName, userRole, action, entity, entityID, description, ip,
	)
	
	if err != nil {
		log.Println("❌ GAGAL MENCATAT LOG:", err)
	}
}