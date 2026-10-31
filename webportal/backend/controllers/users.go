package controllers

import (
	"log"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAdminUsers(c *gin.Context) {
	rows, _ := config.DB.Query(`SELECT id, nama, email, role, COALESCE(opd, ''), COALESCE(is_active, true) FROM users ORDER BY id ASC`)
	defer rows.Close()
	list := []models.User{}
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Nama, &u.Email, &u.Role, &u.OPD, &u.IsActive)
		list = append(list, u)
	}
	if err := rows.Err(); err != nil { log.Println(err); c.JSON(500, gin.H{"statusMessage": "Kesalahan data"}); return }
	c.JSON(200, gin.H{"users": list})
}

func CreateUser(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"statusMessage": "Data invalid"}); return }
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	err := config.DB.QueryRow(`INSERT INTO users (nama, email, password, role, opd, is_active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, req.Nama, req.Email, string(hash), req.Role, req.OPD, req.IsActive).Scan(&req.ID)
	if err != nil { c.JSON(500, gin.H{"statusMessage": "Gagal simpan (Cek duplikasi email)"}); return }
	req.Password = ""
	LogActivity(c, "create", "user", req.ID, "Menambah user baru: "+req.Nama+" ("+req.Role+")")
	c.JSON(200, gin.H{"success": true, "user": req})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(400, gin.H{"statusMessage": "Data invalid"}); return }
	
	if req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		config.DB.Exec(`UPDATE users SET nama=$1, email=$2, password=$3, role=$4, opd=$5, is_active=$6, updated_at=NOW() WHERE id=$7`, req.Nama, req.Email, string(hash), req.Role, req.OPD, req.IsActive, id)
	} else {
		config.DB.Exec(`UPDATE users SET nama=$1, email=$2, role=$3, opd=$4, is_active=$5, updated_at=NOW() WHERE id=$6`, req.Nama, req.Email, req.Role, req.OPD, req.IsActive, id)
	}
	LogActivity(c, "update", "user", 0, "Mengubah data user: "+req.Nama)
	c.JSON(200, gin.H{"success": true})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	config.DB.Exec(`DELETE FROM users WHERE id=$1`, id)
	LogActivity(c, "delete", "user", 0, "Menghapus user ID: "+id)
	c.JSON(200, gin.H{"success": true})
}