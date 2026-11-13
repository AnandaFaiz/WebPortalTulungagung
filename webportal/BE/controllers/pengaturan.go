package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

func GetIdentitas(c *gin.Context) {
	var i models.IdentitasSitus
	query := `SELECT id, COALESCE(nama_situs, ''), COALESCE(tagline, ''), COALESCE(deskripsi, ''), COALESCE(logo, ''), COALESCE(favicon, ''), COALESCE(email, ''), COALESCE(telepon, ''), COALESCE(alamat, '') FROM identitas_situs WHERE id = 1`
	
	err := config.DB.QueryRow(query).Scan(&i.ID, &i.NamaSitus, &i.Tagline, &i.Deskripsi, &i.Logo, &i.Favicon, &i.Email, &i.Telepon, &i.Alamat)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"identitas": models.IdentitasSitus{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"identitas": i})
}

func UpdateIdentitas(c *gin.Context) {
	var req models.IdentitasSitus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	query := `
		INSERT INTO identitas_situs (id, nama_situs, tagline, deskripsi, logo, favicon, email, telepon, alamat, updated_at)
		VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (id) DO UPDATE SET
			nama_situs = EXCLUDED.nama_situs, tagline = EXCLUDED.tagline,
			deskripsi = EXCLUDED.deskripsi, logo = EXCLUDED.logo,
			favicon = EXCLUDED.favicon, email = EXCLUDED.email,
			telepon = EXCLUDED.telepon, alamat = EXCLUDED.alamat, updated_at = NOW()`

	_, err := config.DB.Exec(query, req.NamaSitus, req.Tagline, req.Deskripsi, req.Logo, req.Favicon, req.Email, req.Telepon, req.Alamat)
	if err != nil {
		log.Println("❌ ERROR UPDATE IDENTITAS:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	LogActivity(c, "update", "identitas_situs", 1, "Memperbarui identitas situs")
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ================= KONTAK =================

func GetKontak(c *gin.Context) {
	var k models.Kontak
	query := `SELECT id, COALESCE(alamat, ''), COALESCE(telepon, ''), COALESCE(email, ''), COALESCE(jam_operasional, ''), COALESCE(embed_maps, ''), COALESCE(facebook, ''), COALESCE(instagram, ''), COALESCE(twitter, ''), COALESCE(youtube, '') FROM kontak WHERE id = 1`
	
	err := config.DB.QueryRow(query).Scan(&k.ID, &k.Alamat, &k.Telepon, &k.Email, &k.JamOperasional, &k.EmbedMaps, &k.Facebook, &k.Instagram, &k.Twitter, &k.Youtube)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"kontak": models.Kontak{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"kontak": k})
}

func UpdateKontak(c *gin.Context) {
	var req models.Kontak
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	query := `
		INSERT INTO kontak (id, alamat, telepon, email, jam_operasional, embed_maps, facebook, instagram, twitter, youtube, updated_at)
		VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (id) DO UPDATE SET
			alamat = EXCLUDED.alamat, telepon = EXCLUDED.telepon,
			email = EXCLUDED.email, jam_operasional = EXCLUDED.jam_operasional,
			embed_maps = EXCLUDED.embed_maps, facebook = EXCLUDED.facebook,
			instagram = EXCLUDED.instagram, twitter = EXCLUDED.twitter,
			youtube = EXCLUDED.youtube, updated_at = NOW()`

	_, err := config.DB.Exec(query, req.Alamat, req.Telepon, req.Email, req.JamOperasional, req.EmbedMaps, req.Facebook, req.Instagram, req.Twitter, req.Youtube)
	if err != nil {
		log.Println("❌ ERROR UPDATE KONTAK:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	LogActivity(c, "update", "kontak", 1, "Memperbarui informasi kontak situs")
	c.JSON(http.StatusOK, gin.H{"success": true})
}