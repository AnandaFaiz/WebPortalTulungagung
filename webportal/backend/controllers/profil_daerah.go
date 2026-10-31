package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"pemkab-tulungagung/config"
	"pemkab-tulungagung/models"

	"github.com/gin-gonic/gin"
)

func GetProfilDaerah(c *gin.Context) {
	var p models.ProfilDaerah
	var misiJSON, demografiJSON []byte

	query := `SELECT id, COALESCE(sejarah, ''), COALESCE(lambang, ''), COALESCE(visi, ''), COALESCE(misi, '[]'), COALESCE(geografis, ''), COALESCE(demografi, '{}') FROM profil_daerah WHERE id = 1`
	err := config.DB.QueryRow(query).Scan(&p.ID, &p.Sejarah, &p.Lambang, &p.Visi, &misiJSON, &p.Geografis, &demografiJSON)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"statusMessage": "Data profil belum dibuat"})
		return
	}

	json.Unmarshal(misiJSON, &p.Misi)
	json.Unmarshal(demografiJSON, &p.Demografi)
	
	c.JSON(http.StatusOK, gin.H{"profil": p})
}

func UpdateProfilDaerah(c *gin.Context) {
	var req models.ProfilDaerah
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "Data tidak valid"})
		return
	}

	misiJSON, _ := json.Marshal(req.Misi)
	demografiJSON, _ := json.Marshal(req.Demografi)

	query := `
		INSERT INTO profil_daerah (id, sejarah, lambang, visi, misi, geografis, demografi, updated_at)
		VALUES (1, $1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (id) DO UPDATE SET
			sejarah = EXCLUDED.sejarah,
			lambang = EXCLUDED.lambang,
			visi = EXCLUDED.visi,
			misi = EXCLUDED.misi,
			geografis = EXCLUDED.geografis,
			demografi = EXCLUDED.demografi,
			updated_at = NOW()`

	_, err := config.DB.Exec(query, req.Sejarah, req.Lambang, req.Visi, string(misiJSON), req.Geografis, string(demografiJSON))
	if err != nil {
		log.Println("❌ ERROR UPDATE PROFIL:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal menyimpan data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}