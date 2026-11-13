package controllers

import (
	"net/http"
	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

type dashboardActivity struct {
	ID     int    `json:"id"`
	User   string `json:"user"`
	Action string `json:"action"`
	Target string `json:"target"`
	Time   string `json:"time"`
}

type dashboardComplaint struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Date   string `json:"date"`
}

func dashboardCount(query string) (int, error) {
	var value int
	err := config.DB.QueryRow(query).Scan(&value)
	return value, err
}

func GetAdminDashboard(c *gin.Context) {
	newsCount, err := dashboardCount(`SELECT COUNT(*) FROM berita`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil statistik berita"})
		return
	}
	complaintCount, err := dashboardCount(`SELECT COUNT(*) FROM pengaduan`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil statistik pengaduan"})
		return
	}
	viewsCount, err := dashboardCount(`SELECT COALESCE(SUM(COALESCE(views, 0)), 0) FROM berita`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil statistik tayangan"})
		return
	}
	activeApps, err := dashboardCount(`SELECT COUNT(*) FROM aplikasi WHERE is_active = TRUE`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil statistik aplikasi"})
		return
	}

	activities := []dashboardActivity{}
	activityRows, err := config.DB.Query(`
		SELECT id, COALESCE(user_nama, 'Sistem'), COALESCE(action, ''),
		       COALESCE(description, ''), TO_CHAR(created_at, 'HH24:MI')
		FROM activity_logs
		ORDER BY created_at DESC
		LIMIT 4`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil aktivitas"})
		return
	}
	for activityRows.Next() {
		var activity dashboardActivity
		if err := activityRows.Scan(&activity.ID, &activity.User, &activity.Action, &activity.Target, &activity.Time); err != nil {
			activityRows.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca aktivitas"})
			return
		}
		activities = append(activities, activity)
	}
	if err := activityRows.Err(); err != nil {
		activityRows.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan saat membaca aktivitas"})
		return
	}
	activityRows.Close()

	complaints := []dashboardComplaint{}
	complaintRows, err := config.DB.Query(`
		SELECT id, COALESCE(isi_aduan, ''), COALESCE(status, 'Baru'),
		       TO_CHAR(created_at, 'YYYY-MM-DD')
		FROM pengaduan
		ORDER BY created_at DESC
		LIMIT 3`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil pengaduan terbaru"})
		return
	}
	for complaintRows.Next() {
		var complaint dashboardComplaint
		if err := complaintRows.Scan(&complaint.ID, &complaint.Title, &complaint.Status, &complaint.Date); err != nil {
			complaintRows.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal membaca pengaduan terbaru"})
			return
		}
		complaints = append(complaints, complaint)
	}
	if err := complaintRows.Err(); err != nil {
		complaintRows.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Kesalahan saat membaca pengaduan"})
		return
	}
	complaintRows.Close()

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"news":        newsCount,
			"complaints":  complaintCount,
			"news_views":  viewsCount,
			"active_apps": activeApps,
		},
		"activities": activities,
		"complaints": complaints,
	})
}
