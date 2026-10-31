package controllers

import (
	"log"
	"net/http"
	"os"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// POST /api/pengaduan (Publik - Warga mengirim aduan)
func CreatePengaduan(c *gin.Context) {
	var req models.Pengaduan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak lengkap"})
		return
	}
	
	req.Status = "Baru" // Default status

	err := config.DB.QueryRow(
		`INSERT INTO pengaduan (nama, kontak, kategori, isi_aduan, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Nama, req.Kontak, req.Kategori, req.IsiAduan, req.Status,
	).Scan(&req.ID)

	if err != nil {
		log.Println("❌ ERROR INSERT PENGADUAN:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengirim aduan"})
		return
	}

	// Kirim balik ID sebagai nomor tiket
	c.JSON(http.StatusOK, gin.H{
		"success": true, 
		"message": "Aduan berhasil dikirim", 
		"tiket_id": req.ID,
	})
}

// GET /api/admin/pengaduan (Admin melihat daftar)
func GetAdminPengaduan(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, nama, kontak, kategori, isi_aduan, COALESCE(status, 'Baru'), COALESCE(balasan_admin, ''), TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI') FROM pengaduan ORDER BY created_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	list := []models.Pengaduan{}
	for rows.Next() {
		var p models.Pengaduan
		if err := rows.Scan(&p.ID, &p.Nama, &p.Kontak, &p.Kategori, &p.IsiAduan, &p.Status, &p.BalasanAdmin, &p.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca data"})
			return
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pengaduan": list})
}

// PUT /api/admin/pengaduan/:id (Admin membalas/mengubah status)
func UpdatePengaduan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status       string `json:"status"`
		BalasanAdmin string `json:"balasan_admin"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	_, err := config.DB.Exec(`UPDATE pengaduan SET status=$1, balasan_admin=$2, updated_at=NOW() WHERE id=$3`, req.Status, req.BalasanAdmin, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal memperbarui aduan"})
		return
	}

	LogActivity(c, "update", "pengaduan", 0, "Membalas aduan ID: "+id)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/pengaduan/:id
func DeletePengaduan(c *gin.Context) {
	id := c.Param("id")
	config.DB.Exec(`DELETE FROM pengaduan WHERE id=$1`, id)
	LogActivity(c, "delete", "pengaduan", 0, "Menghapus aduan ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// POST /api/admin/pengaduan/:id/send
func SendReplyPengaduan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status       string `json:"status"`
		BalasanAdmin string `json:"balasan_admin"`
	}
	c.ShouldBindJSON(&req)

	// 1. Ambil data pengadu dari DB
	var kontak string
	err := config.DB.QueryRow("SELECT kontak FROM pengaduan WHERE id = $1", id).Scan(&kontak)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Aduan tidak ditemukan"})
		return
	}

	// 2. Update status dan balasan di DB
	_, err = config.DB.Exec(`UPDATE pengaduan SET status=$1, balasan_admin=$2, updated_at=NOW() WHERE id=$3`, req.Status, req.BalasanAdmin, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan ke database"})
		return
	}

	// 3. Konfigurasi SMTP (Ganti dengan email Pemkab yang asli nanti)
	smtpHost := os.Getenv("SMTP_HOST") // contoh: smtp.gmail.com
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT")) // contoh: 587
	smtpUser := os.Getenv("SMTP_USER") // email pemkab
	smtpPass := os.Getenv("SMTP_PASS") // app password

	// 4. Siapkan Pesan
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", kontak)
	m.SetHeader("Subject", "Balasan Aduan Anda - Pemkab Tulungagung")
	
	// Body email menggunakan HTML
	body := `
	<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; border: 1px solid #ddd; padding: 20px; border-radius: 8px;">
		<h2 style="color: #0284c7;">Balasan Aduan Anda</h2>
		<p>Terima kasih telah mengirimkan aduan kepada Pemerintah Kabupaten Tulungagung.</p>
		<div style="background-color: #f8fafc; padding: 15px; border-radius: 5px; margin: 20px 0;">
			` + req.BalasanAdmin + `
		</div>
		<p>Status aduan Anda saat ini: <strong>` + req.Status + `</strong></p>
		<br>
		<p>Hormat kami,<br>Pemkab Tulungagung</p>
	</div>`
	
	m.SetBody("text/html", body)

	// 5. Kirim Email (Bungkus dengan pengecekan agar tidak crash jika SMTP belum diset)
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	
	if smtpHost != "" && smtpUser != "" {
		if err := d.DialAndSend(m); err != nil {
			log.Println("❌ GAGAL KIRIM EMAIL:", err)
			// Tetap kembalikan sukses karena data sudah masuk DB, tapi beri catatan
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tersimpan, namun gagal kirim email. Cek konfigurasi SMTP."})
			return
		}
	} else {
		log.Println("⚠️ SMTP belum dikonfigurasi di .env, email tidak terkirim.")
	}

	LogActivity(c, "update", "pengaduan", 0, "Membalas & mengirim email aduan ID: "+id)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Balasan berhasil dikirim"})
}