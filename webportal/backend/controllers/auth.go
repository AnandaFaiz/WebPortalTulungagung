package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Fungsi Generate Token
func generateToken(userID int, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// POST /api/auth/login
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Email dan password wajib diisi"})
		return
	}

	var user models.User
	var hashedPassword string
	// Cari user berdasarkan email
	err := config.DB.QueryRow(
		"SELECT id, nama, email, password, role, COALESCE(opd, ''), is_active FROM users WHERE LOWER(email) = LOWER($1)",
		req.Email,
	).Scan(&user.ID, &user.Nama, &user.Email, &hashedPassword, &user.Role, &user.OPD, &user.IsActive)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Email atau password salah"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan server"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"statusMessage": "Akun Anda tidak aktif. Hubungi administrator."})
		return
	}

	// Cek Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Email atau password salah"})
		return
	}

	// Buat Token
	token, err := generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membuat token"})
		return
	}

	var userName string
	config.DB.QueryRow("SELECT nama FROM users WHERE id = $1", user.ID).Scan(&userName)
	ip := c.ClientIP()
	config.DB.Exec(
		`INSERT INTO activity_logs (user_id, user_nama, user_role, action, entity, description, ip_address) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`, 
		user.ID, userName, user.Role, "login", "auth", "Login berhasil sebagai "+user.Role, ip,
	)

	// Set Cookie HttpOnly
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, 0, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// GET /api/auth/me (Butuh Middleware)
func GetMe(c *gin.Context) {
	// Ambil ID dari context yang di-set oleh middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Tidak ada sesi"})
		return
	}

	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, nama, email, role, COALESCE(opd, ''), is_active FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Nama, &user.Email, &user.Role, &user.OPD, &user.IsActive)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// POST /api/auth/logout
func Logout(c *gin.Context) {
	LogActivity(c, "logout", "auth", 0, "Logout dari sistem")
	c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GET /api/auth/permissions
func GetPermissions(c *gin.Context) {
	// Ambil role user dari context (yang sudah di-set oleh AuthMiddleware)
	role, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Tidak ada sesi"})
		return
	}

	userRole := role.(string)

	// Ambil data permission dari database berdasarkan role
	query := `SELECT permission_key, allowed FROM role_permissions WHERE role = $1`
	rows, err := config.DB.Query(query, userRole)
	if err != nil {
		log.Println("❌ ERROR QUERY PERMISSIONS:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data permission"})
		return
	}
	defer rows.Close()

	// Buat map untuk menampung hasil (contoh: {"cms.berita": true})
	permissions := make(map[string]bool)
	for rows.Next() {
		var key string
		var allowed bool
		if err := rows.Scan(&key, &allowed); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		permissions[key] = allowed
	}
	
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"role":        userRole,
		"permissions": permissions,
	})
}

// PUT /api/auth/profile (Butuh Middleware)
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Tidak ada sesi"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama dan email wajib diisi dengan format yang valid"})
		return
	}

	var user models.User
	err := config.DB.QueryRow(
		`UPDATE users SET nama = $1, email = LOWER($2), updated_at = NOW()
		 WHERE id = $3
		 RETURNING id, nama, email, role, COALESCE(opd, ''), is_active, created_at, last_login_at`,
		req.Nama, req.Email, userID,
	).Scan(&user.ID, &user.Nama, &user.Email, &user.Role, &user.OPD, &user.IsActive, &user.CreatedAt, &user.LastLoginAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"statusMessage": "Email sudah digunakan"})
			return
		}
		log.Println("ERROR UPDATE PROFILE:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memperbarui profil"})
		return
	}

	LogActivity(c, "update", "user_profile", user.ID, "Memperbarui profil akun")
	c.JSON(http.StatusOK, gin.H{"success": true, "user": user})
}

// PUT /api/auth/password (Butuh Middleware)
func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Tidak ada sesi"})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Password baru minimal 8 karakter"})
		return
	}

	var hashedPassword string
	if err := config.DB.QueryRow("SELECT password FROM users WHERE id = $1", userID).Scan(&hashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "User tidak ditemukan"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Password saat ini tidak sesuai"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memproses password"})
		return
	}
	if _, err := config.DB.Exec("UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2", string(newHash), userID); err != nil {
		log.Println("ERROR CHANGE PASSWORD:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan password"})
		return
	}

	LogActivity(c, "update", "user_password", userID.(int), "Mengubah password akun")
	c.JSON(http.StatusOK, gin.H{"success": true})
}