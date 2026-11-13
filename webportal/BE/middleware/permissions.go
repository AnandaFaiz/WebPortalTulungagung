package middleware

import (
	"net/http"
	"strings"

	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

type permissionRoute struct {
	prefix string
	key    string
}

var adminPermissionRoutes = []permissionRoute{
	{prefix: "/admin/berita", key: "cms.berita"},
	{prefix: "/admin/kategori-berita", key: "cms.kategori"},
	{prefix: "/admin/pengumuman", key: "cms.pengumuman"},
	{prefix: "/admin/galeri", key: "cms.galeri"},
	{prefix: "/admin/banner", key: "cms.banner"},
	{prefix: "/admin/agenda", key: "cms.event"},
	{prefix: "/admin/layanan-terpopuler", key: "cms.layanan"},
	{prefix: "/admin/pencarian-populer", key: "cms.pencarian"},
	{prefix: "/admin/halaman", key: "cms.halaman"},
	{prefix: "/admin/profil-daerah", key: "pemerintahan.profil"},
	{prefix: "/admin/struktur", key: "pemerintahan.struktur"},
	{prefix: "/admin/pejabat", key: "pemerintahan.pejabat"},
	{prefix: "/admin/opd", key: "pemerintahan.opd"},
	{prefix: "/admin/tautan-eksternal", key: "pemerintahan.tautan"},
	{prefix: "/admin/kategori-layanan", key: "layanan.kategori"},
	{prefix: "/admin/layanan-darurat", key: "layanan.darurat"},
	{prefix: "/admin/layanan", key: "layanan.daftar"},
	{prefix: "/admin/mpp", key: "layanan.mpp"},
	{prefix: "/admin/pengaduan-kategori", key: "interaksi.kategori"},
	{prefix: "/admin/pengaduan-status", key: "interaksi.kategori"},
	{prefix: "/admin/pengaduan", key: "interaksi.tiket"},
	{prefix: "/admin/aplikasi", key: "aplikasi"},
	{prefix: "/admin/media", key: "media"},
	{prefix: "/admin/upload", key: "media"},
	{prefix: "/admin/identitas", key: "pengaturan.identitas"},
	{prefix: "/admin/kontak", key: "pengaturan.kontak"},
	{prefix: "/admin/media-sosial", key: "interaksi.medsos"},
	{prefix: "/admin/menu-navigasi", key: "pengaturan.menu"},
}

func permissionAction(method string) string {
	if method == http.MethodGet {
		return "read"
	}
	if method == http.MethodPost {
		return "create"
	}
	if method == http.MethodDelete {
		return "delete"
	}
	if method == http.MethodPut || method == http.MethodPatch {
		return "update"
	}
	return "read"
}

func findPermissionKey(path string) string {
	for _, route := range adminPermissionRoutes {
		if path == route.prefix || strings.HasPrefix(path, route.prefix+"/") {
			return route.key
		}
	}
	return ""
}

// AdminPermissionMiddleware applies CRUD permissions to admin API routes.
func AdminPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/admin/permissions") || strings.HasPrefix(path, "/api/admin/users") {
			if role, _ := c.Get("user_role"); role != "super_admin" {
				c.JSON(http.StatusForbidden, gin.H{"statusMessage": "Akses khusus Super Admin"})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		if strings.HasPrefix(path, "/api/admin/activity-logs") {
			if role, _ := c.Get("user_role"); role != "super_admin" {
				c.JSON(http.StatusForbidden, gin.H{"statusMessage": "Akses khusus Super Admin"})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// Super Admin bypass: memiliki akses penuh ke seluruh rute admin
		if role, _ := c.Get("user_role"); role == "super_admin" {
			c.Next()
			return
		}

		key := findPermissionKey(strings.TrimPrefix(path, "/api"))
		if key == "" {
			c.Next()
			return
		}

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Sesi tidak valid"})
			c.Abort()
			return
		}
		var allowed bool
		column := "can_" + permissionAction(c.Request.Method)
		err := config.DB.QueryRow(`SELECT `+column+` FROM user_permissions WHERE user_id = $1 AND permission_key = $2`, userID, key).Scan(&allowed)
		if err != nil || !allowed {
			c.JSON(http.StatusForbidden, gin.H{"statusMessage": "Anda tidak memiliki izin untuk aksi ini"})
			c.Abort()
			return
		}
		c.Next()
	}
}
