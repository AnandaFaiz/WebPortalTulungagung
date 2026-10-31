package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// POST /api/admin/permissions/sync
func SyncAdminPermissions(c *gin.Context) {
	if role, exists := c.Get("user_role"); !exists || role != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"statusMessage": "Hanya Super Admin yang dapat menyinkronkan permission"})
		return
	}

	var req struct {
		Keys []string `json:"keys" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Daftar fitur tidak valid"})
		return
	}

	keys := make([]string, 0, len(req.Keys))
	seen := map[string]bool{}
	for _, key := range req.Keys {
		key = strings.TrimSpace(key)
		if key != "" && !seen[key] {
			keys = append(keys, key)
			seen[key] = true
		}
	}
	if len(keys) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Daftar fitur kosong"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memulai sinkronisasi"})
		return
	}
	defer tx.Rollback()

	for _, role := range []string{"super_admin", "admin_opd", "editor"} {
		for _, key := range keys {
			if _, err := tx.Exec(`
				INSERT INTO role_permissions (role, permission_key, allowed)
				VALUES ($1, $2, $3)
				ON CONFLICT (role, permission_key) DO UPDATE SET
					allowed = CASE WHEN $1 = 'super_admin' THEN TRUE ELSE role_permissions.allowed END,
					updated_at = NOW()`, role, key, role == "super_admin"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menambahkan permission baru"})
				return
			}
		}

		placeholders := make([]string, len(keys))
		args := []any{role}
		for i, key := range keys {
			placeholders[i] = "$" + strconv.Itoa(i+2)
			args = append(args, key)
		}
		if _, err := tx.Exec(`DELETE FROM role_permissions WHERE role = $1 AND permission_key NOT IN (`+strings.Join(placeholders, ",")+`)`, args...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menghapus permission usang"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan sinkronisasi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "keys": keys})
}

// GET /api/admin/permissions
func GetAdminPermissions(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, role, permission_key, allowed FROM role_permissions ORDER BY role, permission_key`)
	if err != nil {
		log.Println("Error query permissions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	// Struct sederhana untuk response
	type Perm struct {
		ID            int    `json:"id"`
		Role          string `json:"role"`
		PermissionKey string `json:"permission_key"`
		Allowed       bool   `json:"allowed"`
	}

	list := []Perm{}
	for rows.Next() {
		var p Perm
		if err := rows.Scan(&p.ID, &p.Role, &p.PermissionKey, &p.Allowed); err != nil {
			c.JSON(500, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"statusMessage": "Kesalahan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": list})
}

// PATCH /api/admin/permissions/:role/:key
func UpdatePermission(c *gin.Context) {
	role := c.Param("role")
	key := c.Param("key")

	var req struct {
		Allowed bool `json:"allowed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	// Cegah perubahan untuk super_admin
	if role == "super_admin" {
		c.JSON(403, gin.H{"statusMessage": "Permission Super Admin tidak dapat diubah"})
		return
	}

	_, err := config.DB.Exec(`UPDATE role_permissions SET allowed=$1, updated_at=NOW() WHERE role=$2 AND permission_key=$3`, req.Allowed, role, key)
	if err != nil {
		log.Println("Error update permission:", err)
		c.JSON(500, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "role": role, "permission_key": key, "allowed": req.Allowed})
}
