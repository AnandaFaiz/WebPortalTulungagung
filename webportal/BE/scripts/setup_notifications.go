package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment OS")
	}

	// Database connection
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	fmt.Println("Berhasil terkoneksi ke PostgreSQL!")

	// 1. Create notifications table
	createTableSQL := `
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

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Gagal membuat tabel notifications:", err)
	}
	fmt.Println("✓ Tabel notifications dibuat/terverifikasi")

	// 2. Create indexes
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_role ON notifications(role);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user_role ON notifications(user_id, role);`,
	}

	for _, idx := range indexes {
		_, err = db.Exec(idx)
		if err != nil {
			log.Printf("Peringatan: Gagal membuat index: %v", err)
		}
	}
	fmt.Println("✓ Index dibuat/terverifikasi")

	// 3. Check if table is empty and seed initial notifications
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notifications").Scan(&count)
	if err != nil {
		log.Fatal("Gagal menghitung notifikasi:", err)
	}

	if count == 0 {
		now := time.Now().Format("2006-01-02 15:04:05")

		// Notifications for super_admin
		superAdminNotifs := []struct {
			title   string
			message string
			nType   string
			link    string
		}{
			{"Pengaduan Baru", "Terdapat pengaduan masyarakat baru yang perlu ditinjau.", "pengaduan", "/admin/layanan-publik/pengaduan"},
			{"User Baru Terdaftar", "Seorang user baru telah mendaftar di portal.", "user", "/admin/pengguna"},
			{"Pengaturan Sistem Diubah", "Pengaturan identitas/kontak situs telah diperbarui.", "settings", "/admin/pengaturan"},
		}

		for _, n := range superAdminNotifs {
			_, err = db.Exec(
				`INSERT INTO notifications (role, title, message, type, link, is_read, created_at) VALUES ($1, $2, $3, $4, $5, false, $6)`,
				"super_admin", n.title, n.message, n.nType, n.link, now,
			)
			if err != nil {
				log.Printf("Gagal menambahkan notifikasi super_admin: %v", err)
			}
		}

		// Notifications for editor
		editorNotifs := []struct {
			title   string
			message string
			nType   string
			link    string
		}{
			{"Artikel Diterbitkan", "Berita terbaru Anda berhasil ditayangkan di portal utama.", "berita", "/admin/informasi-publik/berita"},
			{"Artikel Perlu Review", "Terdapat artikel menunggu persetujuan publikasi.", "review", "/admin/informasi-publik/berita"},
			{"Komentar Baru", "Pengunjung meninggalkan komentar pada artikel Anda.", "komentar", "/admin/informasi-publik/berita"},
		}

		for _, n := range editorNotifs {
			_, err = db.Exec(
				`INSERT INTO notifications (role, title, message, type, link, is_read, created_at) VALUES ($1, $2, $3, $4, $5, false, $6)`,
				"editor", n.title, n.message, n.nType, n.link, now,
			)
			if err != nil {
				log.Printf("Gagal menambahkan notifikasi editor: %v", err)
			}
		}

		// Global notifications (for all roles)
		globalNotifs := []struct {
			title   string
			message string
			nType   string
			link    string
		}{
			{"Selamat Datang di Panel Admin", "Selamat datang di Panel Admin Web Portal Kabupaten Tulungagung.", "info", "/admin/dashboard"},
			{"Pembaruan Sistem", "Sistem telah diperbarui ke versi terbaru. Silakan cek fitur-fitur baru.", "update", "/admin/dashboard"},
		}

		for _, n := range globalNotifs {
			_, err = db.Exec(
				`INSERT INTO notifications (role, title, message, type, link, is_read, created_at) VALUES ($1, $2, $3, $4, $5, false, $6)`,
				"", n.title, n.message, n.nType, n.link, now,
			)
			if err != nil {
				log.Printf("Gagal menambahkan notifikasi global: %v", err)
			}
		}

		fmt.Println("✓ Notifikasi awal berhasil di-seed (3 super_admin, 3 editor, 2 global)")
	} else {
		fmt.Printf("✓ Tabel sudah berisi %d notifikasi, skip seeding\n", count)
	}

	// 4. Verify the setup
	rows, err := db.Query(`
		SELECT id, role, title, message, type, link, is_read, created_at
		FROM notifications
		ORDER BY created_at DESC, id DESC
	`)
	if err != nil {
		log.Fatal("Gagal query verifikasi:", err)
	}
	defer rows.Close()

	fmt.Println("\n=== DAFTAR NOTIFIKASI ===")
	fmt.Printf("%-4s %-14s %-30s %-50s %-12s %-30s %-6s %s\n", "ID", "ROLE", "TITLE", "MESSAGE", "TYPE", "LINK", "READ", "CREATED_AT")
	for rows.Next() {
		var id int
		var role, title, message, nType, link string
		var isRead bool
		var createdAt time.Time

		err := rows.Scan(&id, &role, &title, &message, &nType, &link, &isRead, &createdAt)
		if err != nil {
			log.Printf("Error scanning: %v", err)
			continue
		}

		roleDisplay := role
		if role == "" {
			roleDisplay = "(global)"
		}
		readStr := "No"
		if isRead {
			readStr = "Yes"
		}

		// Truncate long fields for display
		if len(title) > 28 {
			title = title[:28] + ".."
		}
		if len(message) > 48 {
			message = message[:48] + ".."
		}
		if len(link) > 28 {
			link = link[:28] + ".."
		}

		fmt.Printf("%-4d %-14s %-30s %-50s %-12s %-30s %-6s %s\n",
			id, roleDisplay, title, message, nType, link, readStr, createdAt.Format("2006-01-02 15:04:05"))
	}

	// Count unread for each role
	for _, role := range []string{"super_admin", "editor"} {
		var unreadCount int
		err = db.QueryRow(`
			SELECT COUNT(*) FROM notifications 
			WHERE (user_id IS NULL OR user_id = 1) 
			  AND (role = '' OR role = $1)
			  AND is_read = false
		`, role).Scan(&unreadCount)
		if err != nil {
			log.Printf("Error counting unread for %s: %v", role, err)
		} else {
			fmt.Printf("\nUnread count for %s: %d\n", role, unreadCount)
		}
	}

	fmt.Println("\n✅ Setup notifikasi selesai!")
}