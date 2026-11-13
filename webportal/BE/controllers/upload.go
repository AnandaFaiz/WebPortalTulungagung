package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "File tidak ditemukan"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Format harus JPG, PNG, atau WebP"})
		return
	}

	filename := fmt.Sprintf("img-%d%s", time.Now().Unix(), ext)
	uploadDir := "./public/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	filePath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan file"})
		return
	}

	fileUrl := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)

	_, err = config.DB.Exec(
		`INSERT INTO media (nama_file, url, tipe_file, ukuran_file) VALUES ($1, $2, $3, $4)`,
		file.Filename, fileUrl, file.Header.Get("Content-Type"), file.Size,
	)
	if err != nil {
		fmt.Println("❌ Gagal mencatat media ke DB:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"url":     fileUrl,
	})
}