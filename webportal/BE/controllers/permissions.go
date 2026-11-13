package controllers

import (
	"log"
	"net/http"
	"strconv"

	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

type AdminPermission struct {
	ID            int    `json:"id"`
	PermissionKey string `json:"permission_key"`
	CanCreate     bool   `json:"can_create"`
	CanRead       bool   `json:"can_read"`
	CanUpdate     bool   `json:"can_update"`
	CanDelete     bool   `json:"can_delete"`
}

func requireSuperAdmin(c *gin.Context) bool {
	role, exists := c.Get("user_role")
	if !exists || role != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"statusMessage": "Hanya Super Admin yang dapat mengatur hak akses"})
		return false
	}
	return true
}

// GET /api/admin/permissions
func GetAdminPermissions(c *gin.Context) {
	if !requireSuperAdmin(c) {
		return
	}

	rows, err := config.DB.Query(`
		SELECT id, nama, email, role, COALESCE(opd, ''), COALESCE(is_active, true)
		FROM users WHERE role IN ('super_admin', 'editor') ORDER BY role, nama`)
	if err != nil {
		log.Println("Error query admin users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil daftar admin"})
		return
	}
	defer rows.Close()

	type Admin struct {
		ID       int    `json:"id"`
		Nama     string `json:"nama"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		OPD      string `json:"opd"`
		IsActive bool   `json:"is_active"`
	}
	admins := []Admin{}
	for rows.Next() {
		var admin Admin
		if err := rows.Scan(&admin.ID, &admin.Nama, &admin.Email, &admin.Role, &admin.OPD, &admin.IsActive); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data admin"})
			return
		}
		admins = append(admins, admin)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca daftar admin"})
		return
	}

	permissionsRows, err := config.DB.Query(`
		SELECT id, user_id, permission_key, can_create, can_read, can_update, can_delete
		FROM user_permissions ORDER BY user_id, permission_key`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil matriks hak akses"})
		return
	}
	defer permissionsRows.Close()

	permissions := map[int][]AdminPermission{}
	for permissionsRows.Next() {
		var userID int
		var permission AdminPermission
		if err := permissionsRows.Scan(&permission.ID, &userID, &permission.PermissionKey, &permission.CanCreate, &permission.CanRead, &permission.CanUpdate, &permission.CanDelete); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca matriks hak akses"})
			return
		}
		permissions[userID] = append(permissions[userID], permission)
	}
	if err := permissionsRows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca matriks hak akses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"admins": admins, "permissions": permissions})
}

// PATCH /api/admin/permissions/:userID/:key/:action
func UpdatePermission(c *gin.Context) {
	if !requireSuperAdmin(c) {
		return
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Admin tidak valid"})
		return
	}
	key := c.Param("key")
	action := c.Param("action")
	if action != "create" && action != "read" && action != "update" && action != "delete" {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Aksi permission tidak valid"})
		return
	}

	var req struct {
		Allowed bool `json:"allowed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nilai permission tidak valid"})
		return
	}

	var role string
	if err := config.DB.QueryRow(`SELECT role FROM users WHERE id = $1`, userID).Scan(&role); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Admin tidak ditemukan"})
		return
	}
	if role != "super_admin" && role != "editor" {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Role admin tidak valid"})
		return
	}

	column := "can_" + action
	query := `
		INSERT INTO user_permissions (user_id, permission_key, ` + column + `, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, permission_key) DO UPDATE
		SET ` + column + ` = EXCLUDED.` + column + `, updated_at = NOW()`
	_, err = config.DB.Exec(query, userID, key, req.Allowed)
	if err != nil {
		log.Println("Error update permission:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan hak akses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "user_id": userID, "permission_key": key, "action": action, "allowed": req.Allowed})
}
