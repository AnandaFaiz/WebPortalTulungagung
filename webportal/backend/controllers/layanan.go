package controllers

import (
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

// ==== KATEGORI LAYANAN ====

// GET /api/layanan (Publik - Bersarang/Nested)
func GetPublicLayananKategori(c *gin.Context) {
	katRows, err := config.DB.Query("SELECT id, nama, COALESCE(deskripsi, ''), COALESCE(icon, ''), COALESCE(urutan, 0) FROM kategori_layanan ORDER BY urutan ASC")
	if err != nil {
		log.Println("Error Query Kategori:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer katRows.Close()

	var kategoris []models.KategoriLayanan
	for katRows.Next() {
		var k models.KategoriLayanan
		katRows.Scan(&k.ID, &k.Nama, &k.Deskripsi, &k.Icon, &k.Urutan)
		
		layananRows, _ := config.DB.Query("SELECT id, kategori_id, nama, COALESCE(link, ''), COALESCE(urutan, 0) FROM layanan WHERE kategori_id = $1 ORDER BY urutan ASC", k.ID)
		
		var items []models.Layanan
		for layananRows.Next() {
			var l models.Layanan
			layananRows.Scan(&l.ID, &l.KategoriID, &l.Nama, &l.Link, &l.Urutan)
			items = append(items, l)
		}
		if err := layananRows.Err(); err != nil {
			log.Println("Error Loop Sub-Layanan:", err)
			layananRows.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data sub-layanan"})
			return
		}
		layananRows.Close()
		if items == nil { items = []models.Layanan{} }
		
		k.Items = items
		kategoris = append(kategoris, k)
	}
	if err := katRows.Err(); err != nil {
		log.Println("Error Loop Kategori:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"kategori": kategoris})
}

// GET /api/admin/kategori-layanan
func GetAdminKategoriLayanan(c *gin.Context) {
	rows, _ := config.DB.Query("SELECT id, nama, COALESCE(deskripsi, ''), COALESCE(icon, ''), COALESCE(urutan, 0) FROM kategori_layanan ORDER BY urutan ASC")
	defer rows.Close()
	
	list := []models.KategoriLayanan{}
	for rows.Next() {
		var k models.KategoriLayanan
		rows.Scan(&k.ID, &k.Nama, &k.Deskripsi, &k.Icon, &k.Urutan)
		list = append(list, k)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error Loop Admin Kategori:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"kategori": list})
}

// POST /api/admin/kategori-layanan
func CreateKategoriLayanan(c *gin.Context) {
	var req models.KategoriLayanan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama wajib diisi"})
		return
	}
	var maxUrut int
	config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM kategori_layanan").Scan(&maxUrut)
	
	err := config.DB.QueryRow("INSERT INTO kategori_layanan (nama, deskripsi, icon, urutan) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Nama, req.Deskripsi, req.Icon, maxUrut).Scan(&req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan"})
		return
	}
	LogActivity(c, "create", "kategori_layanan", req.ID, "Menambah kategori layanan: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true, "kategori": req})
}

// PUT /api/admin/kategori-layanan/:id
func UpdateKategoriLayanan(c *gin.Context) {
	id := c.Param("id")
	var req models.KategoriLayanan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}
	config.DB.Exec("UPDATE kategori_layanan SET nama=$1, deskripsi=$2, icon=$3 WHERE id=$4", req.Nama, req.Deskripsi, req.Icon, id)
	LogActivity(c, "update", "kategori_layanan", 0, "Mengubah kategori layanan: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/kategori-layanan/:id
func DeleteKategoriLayanan(c *gin.Context) {
	id := c.Param("id")
	config.DB.Exec("DELETE FROM kategori_layanan WHERE id = $1", id)
	LogActivity(c, "delete", "kategori_layanan", 0, "Menghapus kategori layanan ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}


// ==== SUB-LAYANAN (DROPDOWN ITEMS) ====

// GET /api/admin/layanan
func GetAdminSubLayanan(c *gin.Context) {
	rows, _ := config.DB.Query("SELECT id, kategori_id, nama, COALESCE(link, ''), COALESCE(urutan, 0) FROM layanan ORDER BY kategori_id, urutan ASC")
	defer rows.Close()
	
	list := []models.Layanan{}
	for rows.Next() {
		var l models.Layanan
		rows.Scan(&l.ID, &l.KategoriID, &l.Nama, &l.Link, &l.Urutan)
		list = append(list, l)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error Loop Admin Sub-Layanan:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"layanan": list})
}

// POST /api/admin/layanan
func CreateSubLayanan(c *gin.Context) {
	var req models.Layanan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Nama wajib diisi"})
		return
	}
	var maxUrut int
	config.DB.QueryRow("SELECT COALESCE(MAX(urutan), 0) + 1 FROM layanan WHERE kategori_id = $1", req.KategoriID).Scan(&maxUrut)
	
	config.DB.QueryRow("INSERT INTO layanan (kategori_id, nama, link, urutan) VALUES ($1, $2, $3, $4) RETURNING id", req.KategoriID, req.Nama, req.Link, maxUrut).Scan(&req.ID)
	LogActivity(c, "create", "layanan", req.ID, "Menambah layanan: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true, "layanan": req})
}

// PUT /api/admin/layanan/:id
func UpdateSubLayanan(c *gin.Context) {
	id := c.Param("id")
	var req models.Layanan
	if err := c.ShouldBindJSON(&req); err != nil { return }
	config.DB.Exec("UPDATE layanan SET nama=$1, link=$2 WHERE id=$3", req.Nama, req.Link, id)
	LogActivity(c, "update", "layanan", 0, "Mengubah layanan: "+req.Nama)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DELETE /api/admin/layanan/:id
func DeleteSubLayanan(c *gin.Context) {
	id := c.Param("id")
	config.DB.Exec("DELETE FROM layanan WHERE id = $1", id)
	LogActivity(c, "delete", "layanan", 0, "Menghapus layanan ID: "+id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}