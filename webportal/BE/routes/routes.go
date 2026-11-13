package routes

import (
	"pemkab-tulungagung/controllers"
	"pemkab-tulungagung/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	api := r.Group("/api")

	// Auth Routes
	auth := api.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/logout", controllers.Logout)
		auth.GET("/me", middleware.AuthMiddleware(), controllers.GetMe)
		auth.PUT("/profile", middleware.AuthMiddleware(), controllers.UpdateProfile)
		auth.PUT("/password", middleware.AuthMiddleware(), controllers.ChangePassword)
		auth.GET("/permissions", middleware.AuthMiddleware(), controllers.GetPermissions)
	}

	// Public Routes
	api.GET("/berita", controllers.GetPublicBerita)
	api.GET("/berita/:id", controllers.GetPublicBeritaDetail)
	api.POST("/berita/:id/view", controllers.AddBeritaView)
	api.GET("/banner", controllers.GetPublicBanner)
	api.GET("/search", controllers.GlobalSearch)
	api.GET("/pencarian-populer", controllers.GetPencarianPopuler)
	api.GET("/kategori-berita", controllers.GetKategoriBerita)
	api.GET("/agenda", controllers.GetPublicAgenda)
	api.GET("/layanan-terpopuler", controllers.GetPublicLayanan)
	api.GET("/pengumuman", controllers.GetPublicPengumuman)
	api.GET("/galeri", controllers.GetPublicGaleri)
	api.GET("/opd", controllers.GetPublicOPD)
	api.GET("/pejabat", controllers.GetPublicPejabat)
	api.GET("/profil-daerah", controllers.GetProfilDaerah)
	api.GET("/struktur", controllers.GetStrukturOrganisasi)
	api.GET("/aplikasi", controllers.GetPublicAplikasi)
	api.GET("/layanan", controllers.GetPublicLayananKategori)
	api.GET("/identitas", controllers.GetIdentitas)
	api.GET("/kontak", controllers.GetKontak)
	api.GET("/menu-navigasi", controllers.GetPublicMenu)
	api.GET("/pengaduan-konfigurasi", controllers.GetPublicPengaduanConfig)
	api.GET("/layanan-darurat", controllers.GetPublicDarurat)
	api.POST("/pengaduan", controllers.CreatePengaduan)
	api.GET("/mpp", controllers.GetMPP)
	api.GET("/settings/:key", controllers.GetSetting)
    api.PUT("/settings/:key", controllers.UpdateSetting)

	// Admin Routes (Protected)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminPermissionMiddleware())
	{
		//Role Permissions
		admin.GET("/permissions", controllers.GetAdminPermissions)
		admin.PATCH("/permissions/:userID/:key/:action", controllers.UpdatePermission)

		//Activity logs
		admin.GET("/activity-logs", controllers.GetAdminActivityLogs)
		admin.GET("/dashboard", controllers.GetAdminDashboard)

		// Berita Routes
		admin.GET("/berita", controllers.GetAdminBerita)
		admin.POST("/berita", controllers.CreateBerita)
		admin.PUT("/berita/:id", controllers.UpdateBerita)
		admin.DELETE("/berita/:id", controllers.DeleteBerita)

		// Banner Routes
		admin.GET("/banner", controllers.GetAdminBanner)
		admin.POST("/banner", controllers.CreateBanner)
		admin.PUT("/banner/:id", controllers.UpdateBanner)
		admin.DELETE("/banner/:id", controllers.DeleteBanner)
		admin.PUT("/banner/reorder", controllers.ReorderBanner)

		// Agenda Routes
		admin.GET("/agenda", controllers.GetAdminAgenda)
		admin.POST("/agenda", controllers.CreateAgenda)
		admin.PUT("/agenda/:id", controllers.UpdateAgenda)
		admin.DELETE("/agenda/:id", controllers.DeleteAgenda)

		admin.POST("/upload", controllers.UploadFile)

		//Layanan Populer Routes
		admin.GET("/layanan-terpopuler", controllers.GetAdminLayanan)
		admin.POST("/layanan-terpopuler", controllers.CreateLayanan)
		admin.PUT("/layanan-terpopuler/:id", controllers.UpdateLayanan)
		admin.DELETE("/layanan-terpopuler/:id", controllers.DeleteLayanan)
		admin.PUT("/layanan-terpopuler/reorder", controllers.ReorderLayanan)

		//Pengumuman
		admin.GET("/pengumuman", controllers.GetAdminPengumuman)
		admin.POST("/pengumuman", controllers.CreatePengumuman)
		admin.PUT("/pengumuman/:id", controllers.UpdatePengumuman)
		admin.DELETE("/pengumuman/:id", controllers.DeletePengumuman)

		//Galeri
		admin.GET("/galeri", controllers.GetAdminGaleri)
		admin.POST("/galeri", controllers.CreateGaleri)
		admin.PUT("/galeri/:id", controllers.UpdateGaleri)
		admin.DELETE("/galeri/:id", controllers.DeleteGaleri)

		//pencarian fav
		admin.GET("/pencarian-populer", controllers.GetAdminPencarian)
		admin.POST("/pencarian-populer", controllers.CreatePencarian)
		admin.PUT("/pencarian-populer/:id", controllers.UpdatePencarian)
		admin.DELETE("/pencarian-populer/:id", controllers.DeletePencarian)

		//kategori berita
		admin.GET("/kategori-berita", controllers.GetAdminKategori)
		admin.POST("/kategori-berita", controllers.CreateKategori)
		admin.PUT("/kategori-berita/:id", controllers.UpdateKategori)
		admin.DELETE("/kategori-berita/:id", controllers.DeleteKategori)

		//OPD
		admin.GET("/opd", controllers.GetAdminOPD)
		admin.POST("/opd", controllers.CreateOPD)
		admin.PUT("/opd/:id", controllers.UpdateOPD)
		admin.DELETE("/opd/:id", controllers.DeleteOPD)
		admin.PUT("/opd/reorder", controllers.ReorderOPD)

		//Pejabat
		admin.GET("/pejabat", controllers.GetAdminPejabat)
		admin.POST("/pejabat", controllers.CreatePejabat)
		admin.PUT("/pejabat/:id", controllers.UpdatePejabat)
		admin.DELETE("/pejabat/:id", controllers.DeletePejabat)
		admin.PUT("/pejabat/reorder", controllers.ReorderPejabat)

		//Profil Daerah dan Struktur Organisasi
		admin.GET("/profil-daerah", controllers.GetProfilDaerah)
		admin.PUT("/profil-daerah", controllers.UpdateProfilDaerah)
		admin.GET("/struktur", controllers.GetStrukturOrganisasi)
		admin.PUT("/struktur", controllers.UpdateStrukturOrganisasi)

		//Aplikasi
		admin.GET("/aplikasi", controllers.GetAdminAplikasi)
		admin.POST("/aplikasi", controllers.CreateAplikasi)
		admin.PUT("/aplikasi/:id", controllers.UpdateAplikasi)
		admin.DELETE("/aplikasi/:id", controllers.DeleteAplikasi)
		admin.PUT("/aplikasi/reorder", controllers.ReorderAplikasi)

		//layanan
		admin.GET("/kategori-layanan", controllers.GetAdminKategoriLayanan)
		admin.POST("/kategori-layanan", controllers.CreateKategoriLayanan)
		admin.PUT("/kategori-layanan/:id", controllers.UpdateKategoriLayanan)
		admin.DELETE("/kategori-layanan/:id", controllers.DeleteKategoriLayanan)

		admin.GET("/layanan", controllers.GetAdminSubLayanan)
		admin.POST("/layanan", controllers.CreateSubLayanan)
		admin.PUT("/layanan/:id", controllers.UpdateSubLayanan)
		admin.DELETE("/layanan/:id", controllers.DeleteSubLayanan)

		//layanan darurat
		admin.GET("/layanan-darurat", controllers.GetAdminDarurat)
		admin.POST("/layanan-darurat", controllers.CreateDarurat)
		admin.PUT("/layanan-darurat/:id", controllers.UpdateDarurat)
		admin.DELETE("/layanan-darurat/:id", controllers.DeleteDarurat)
		admin.PUT("/layanan-darurat/reorder", controllers.ReorderDarurat)

		//Pengaduan/aduan warga
		admin.GET("/pengaduan", controllers.GetAdminPengaduan)
		admin.PUT("/pengaduan/:id", controllers.UpdatePengaduan)
		admin.DELETE("/pengaduan/:id", controllers.DeletePengaduan)
		admin.POST("/pengaduan/:id/send", controllers.SendReplyPengaduan)

		//identitas situs
		admin.GET("/identitas", controllers.GetIdentitas)
		admin.PUT("/identitas", controllers.UpdateIdentitas)

		//kontak
		admin.GET("/kontak", controllers.GetKontak)
		admin.PUT("/kontak", controllers.UpdateKontak)
		admin.GET("/media-sosial", controllers.GetKontak)
		admin.PUT("/media-sosial", controllers.UpdateKontak)

		// Menu navigasi
		admin.GET("/menu-navigasi", controllers.GetAdminMenu)
		admin.POST("/menu-navigasi", controllers.CreateMenu)
		admin.PUT("/menu-navigasi/reorder", controllers.ReorderMenu)
		admin.PUT("/menu-navigasi/:id", controllers.UpdateMenu)
		admin.DELETE("/menu-navigasi/:id", controllers.DeleteMenu)

		// Kategori dan status pengaduan
		admin.GET("/pengaduan-kategori", controllers.GetAdminPengaduanKategori)
		admin.POST("/pengaduan-kategori", controllers.CreatePengaduanKategori)
		admin.PUT("/pengaduan-kategori/:id", controllers.UpdatePengaduanKategori)
		admin.DELETE("/pengaduan-kategori/:id", controllers.DeletePengaduanKategori)
		admin.GET("/pengaduan-status", controllers.GetAdminPengaduanStatus)
		admin.POST("/pengaduan-status", controllers.CreatePengaduanStatus)
		admin.PUT("/pengaduan-status/:id", controllers.UpdatePengaduanStatus)
		admin.DELETE("/pengaduan-status/:id", controllers.DeletePengaduanStatus)

		// Halaman Statis
		admin.GET("/halaman", controllers.GetAdminHalaman)
		admin.POST("/halaman", controllers.CreateHalaman)
		admin.PUT("/halaman/:id", controllers.UpdateHalaman)
		admin.DELETE("/halaman/:id", controllers.DeleteHalaman)

		// Tautan Eksternal
		admin.GET("/tautan-eksternal", controllers.GetAdminTautan)
		admin.POST("/tautan-eksternal", controllers.CreateTautan)
		admin.PUT("/tautan-eksternal/:id", controllers.UpdateTautan)
		admin.DELETE("/tautan-eksternal/:id", controllers.DeleteTautan)
		admin.PUT("/tautan-eksternal/reorder", controllers.ReorderTautan)

		//MPP
		admin.GET("/mpp", controllers.GetMPP)
		admin.PUT("/mpp", controllers.UpdateMPP)

		//Media
		admin.GET("/media", controllers.GetAdminMedia)
		admin.DELETE("/media/:id", controllers.DeleteMedia)

		//Notifikasi
		admin.GET("/notifications", controllers.GetNotifications)
		admin.PUT("/notifications/:id/read", controllers.MarkNotificationAsRead)
		admin.PUT("/notifications/read-all", controllers.MarkAllNotificationsAsRead)
	}
}
