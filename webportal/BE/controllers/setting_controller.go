package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Struct Model
type Setting struct {
	Key   string `gorm:"primaryKey;type:varchar(50)" json:"key"`
	Value string `gorm:"type:jsonb" json:"value"` // PostgreSQL JSONB
}

// GET /api/settings/:key
func GetSetting(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	key := c.Param("key")

	var setting Setting
	if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
		// Jika belum ada, kembalikan objek kosong (bukan error 404)
		c.JSON(http.StatusOK, gin.H{"key": key, "value": "{}"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// PUT /api/settings/:key
func UpdateSetting(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	key := c.Param("key")

	var input struct {
		Value string `json:"value" binding:"required"` // Dikirim sebagai string JSON dari Nuxt
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting := Setting{Key: key, Value: input.Value}
	
	// Upsert (Update jika ada, Insert jika belum ada)
	if err := db.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan berhasil disimpan", "data": setting})
}