package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

type NotificationItem struct {
	ID        int    `json:"id"`
	UserID    *int   `json:"user_id"`
	Role      string `json:"role"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Link      string `json:"link"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

func EnsureNotificationsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id SERIAL PRIMARY KEY,
		user_id INT,
		role VARCHAR(50) DEFAULT '',
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		type VARCHAR(50) DEFAULT 'info',
		link VARCHAR(255) DEFAULT '',
		is_read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := config.DB.Exec(query); err != nil {
		log.Println("Error ensuring notifications table:", err)
	} else {
		// Seed default notifications if table is empty
		var count int
		_ = config.DB.QueryRow("SELECT COUNT(*) FROM notifications").Scan(&count)
		if count == 0 {
			seedInitialNotifications()
		}
	}
}

func seedInitialNotifications() {
	now := time.Now().Format("2006-01-02 15:04:05")
	initials := []struct {
		role    string
		title   string
		message string
		nType   string
		link    string
	}{
		{
			role:    "super_admin",
			title:   "Pengaduan Baru",
			message: "Terdapat pengaduan masyarakat baru yang perlu ditinjau.",
			nType:   "pengaduan",
			link:    "/admin/layanan-publik/pengaduan",
		},
		{
			role:    "editor",
			title:   "Artikel Diterbitkan",
			message: "Berita terbaru Anda berhasil ditayangkan di portal utama.",
			nType:   "berita",
			link:    "/admin/informasi-publik/berita",
		},
		{
			role:    "",
			title:   "Pembaruan Sistem",
			message: "Selamat datang di Panel Admin Web Portal Kabupaten Tulungagung.",
			nType:   "info",
			link:    "/admin/dashboard",
		},
	}

	for _, item := range initials {
		_, err := config.DB.Exec(
			`INSERT INTO notifications (role, title, message, type, link, is_read, created_at) VALUES ($1, $2, $3, $4, $5, false, $6)`,
			item.role, item.title, item.message, item.nType, item.link, now,
		)
		if err != nil {
			log.Println("Error seeding initial notification:", err)
		}
	}
}

// Helper to create notifications from other controllers
func CreateNotificationHelper(userID *int, targetRole string, title, message, link, notifType string) {
	EnsureNotificationsTable()
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err := config.DB.Exec(
		`INSERT INTO notifications (user_id, role, title, message, link, type, is_read, created_at) VALUES ($1, $2, $3, $4, $5, $6, false, $7)`,
		userID, targetRole, title, message, link, notifType, now,
	)
	if err != nil {
		log.Println("Error creating notification:", err)
	}
}

// GET /api/admin/notifications
func GetNotifications(c *gin.Context) {
	EnsureNotificationsTable()

	userIDVal, existsUser := c.Get("user_id")
	userRoleVal, existsRole := c.Get("user_role")

	var currentUserID int
	var currentRole string

	if existsUser {
		if uid, ok := userIDVal.(int); ok {
			currentUserID = uid
		}
	}
	if existsRole {
		if r, ok := userRoleVal.(string); ok {
			currentRole = r
		}
	}

	rows, err := config.DB.Query(`
		SELECT id, user_id, COALESCE(role, ''), title, message, COALESCE(type, 'info'), COALESCE(link, ''), is_read, 
		       TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
		FROM notifications
		WHERE (user_id IS NULL OR user_id = $1)
		  AND (role = '' OR role = $2)
		ORDER BY created_at DESC, id DESC
		LIMIT 20
	`, currentUserID, currentRole)

	if err != nil {
		log.Println("Error fetching notifications:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil notifikasi"})
		return
	}
	defer rows.Close()

	notifications := []NotificationItem{}
	unreadCount := 0

	for rows.Next() {
		var n NotificationItem
		var uid sql.NullInt64
		if err := rows.Scan(&n.ID, &uid, &n.Role, &n.Title, &n.Message, &n.Type, &n.Link, &n.IsRead, &n.CreatedAt); err != nil {
			log.Println("Error scanning notification:", err)
			continue
		}
		if uid.Valid {
			userIDInt := int(uid.Int64)
			n.UserID = &userIDInt
		}
		if !n.IsRead {
			unreadCount++
		}
		notifications = append(notifications, n)
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"unread_count":  unreadCount,
	})
}

// PUT /api/admin/notifications/:id/read
func MarkNotificationAsRead(c *gin.Context) {
	EnsureNotificationsTable()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "ID notifikasi tidak valid"})
		return
	}

	_, err = config.DB.Exec("UPDATE notifications SET is_read = true WHERE id = $1", id)
	if err != nil {
		log.Println("Error marking notification read:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memperbarui status notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusMessage": "Notifikasi ditandai dibaca"})
}

// PUT /api/admin/notifications/read-all
func MarkAllNotificationsAsRead(c *gin.Context) {
	EnsureNotificationsTable()

	userIDVal, existsUser := c.Get("user_id")
	userRoleVal, existsRole := c.Get("user_role")

	var currentUserID int
	var currentRole string

	if existsUser {
		if uid, ok := userIDVal.(int); ok {
			currentUserID = uid
		}
	}
	if existsRole {
		if r, ok := userRoleVal.(string); ok {
			currentRole = r
		}
	}

	_, err := config.DB.Exec(`
		UPDATE notifications 
		SET is_read = true 
		WHERE (user_id IS NULL OR user_id = $1)
		  AND (role = '' OR role = $2)
		  AND is_read = false
	`, currentUserID, currentRole)

	if err != nil {
		log.Println("Error marking all notifications read:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memperbarui notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusMessage": "Semua notifikasi ditandai dibaca"})
}
