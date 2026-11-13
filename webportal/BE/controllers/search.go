package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"pemkab-tulungagung/config"

	"github.com/gin-gonic/gin"
)

// GET /api/search?q=...
func GlobalSearch(c *gin.Context) {
	terms := normalizeSearchTerms(c.Query("q"))
	if len(terms) == 0 {
		c.JSON(http.StatusOK, emptySearchResult())
		return
	}

	condition, args := buildSearchCondition([]string{"judul", "deskripsi", "konten", "kategori", "penulis"}, terms)
	query := `
		SELECT id, judul, deskripsi, tipe, slug, url FROM (
			SELECT id, judul, COALESCE(deskripsi, '') AS deskripsi, 'berita' AS tipe, '' AS slug, '' AS url, tanggal::text AS sort_date
			FROM berita WHERE status = 'Publish' AND ` + condition + `
			UNION ALL
			SELECT id, judul, COALESCE(konten, ''), 'halaman', slug, '', tanggal::text
			FROM halaman_statis WHERE status = 'Publish' AND ` + buildSearchConditionFor([]string{"judul", "slug", "konten"}, terms) + `
			UNION ALL
			SELECT id, nama, COALESCE(deskripsi, ''), 'opd', '', COALESCE(website, ''), created_at::text
			FROM opds WHERE is_active = TRUE AND ` + buildSearchConditionFor([]string{"nama", "deskripsi", "website"}, terms) + `
			UNION ALL
			SELECT id, nama, nomor, 'darurat', '', '', id::text
			FROM layanan_darurat WHERE is_active = TRUE AND ` + buildSearchConditionFor([]string{"nama", "nomor"}, terms) + `
			UNION ALL
			SELECT id, nama, COALESCE(deskripsi, ''), 'aplikasi', '', COALESCE(url, ''), created_at::text
			FROM aplikasi WHERE is_active = TRUE AND ` + buildSearchConditionFor([]string{"nama", "deskripsi", "url"}, terms) + `
			UNION ALL
			SELECT id, judul, COALESCE(konten, ''), 'pengumuman', '', '', tanggal::text
			FROM pengumuman WHERE status = 'Publish' AND ` + buildSearchConditionFor([]string{"judul", "konten"}, terms) + `
			UNION ALL
			SELECT id, judul, COALESCE(deskripsi, ''), 'galeri', '', COALESCE(file_url, ''), tanggal::text
			FROM galeri WHERE status = 'Publish' AND ` + buildSearchConditionFor([]string{"judul", "deskripsi", "file_url"}, terms) + `
			UNION ALL
			SELECT id, judul, COALESCE(deskripsi, ''), 'agenda', '', COALESCE(link_event, ''), tanggal::text
			FROM agenda WHERE status = 'Publish' AND ` + buildSearchConditionFor([]string{"judul", "deskripsi", "lokasi", "kategori"}, terms) + `
			UNION ALL
			SELECT id, nama, CONCAT_WS(' - ', jabatan, kategori, periode), 'pejabat', '', '', id::text
			FROM pejabat WHERE is_active = TRUE AND ` + buildSearchConditionFor([]string{"nama", "jabatan", "kategori", "periode"}, terms) + `
			UNION ALL
			SELECT id, nama, COALESCE(deskripsi, ''), 'tautan', '', url, urutan::text
			FROM tautan_eksternal WHERE is_active = TRUE AND ` + buildSearchConditionFor([]string{"nama", "url", "deskripsi"}, terms) + `
			UNION ALL
			SELECT id, nama, COALESCE(deskripsi, ''), 'kategori_layanan', '', '', urutan::text
			FROM kategori_layanan WHERE ` + buildSearchConditionFor([]string{"nama", "deskripsi"}, terms) + `
			UNION ALL
			SELECT l.id, l.nama, CONCAT_WS(' - ', k.nama, l.link), 'layanan', '', COALESCE(l.link, ''), l.urutan::text
			FROM layanan l LEFT JOIN kategori_layanan k ON k.id = l.kategori_id
			WHERE ` + buildSearchConditionFor([]string{"l.nama", "l.link", "k.nama"}, terms) + `
			UNION ALL
			SELECT id, judul, COALESCE(deskripsi, ''), 'layanan_populer', '', COALESCE(link, ''), urutan::text
			FROM layanan_terpopuler WHERE is_active = TRUE AND ` + buildSearchConditionFor([]string{"judul", "deskripsi", "kategori", "link"}, terms) + `
			UNION ALL
			SELECT id, label, COALESCE(deskripsi, ''), 'menu', '', COALESCE(url, ''), urutan::text
			FROM menu_navigasi WHERE is_active = TRUE AND ` + buildSearchConditionFor([]string{"label", "url", "deskripsi", "grup"}, terms) + `
			UNION ALL
			SELECT id, 'Profil Daerah', CONCAT_WS(' ', sejarah, visi, misi::text, geografis, demografi::text), 'profil_daerah', '', '', id::text
			FROM profil_daerah WHERE ` + buildSearchConditionFor([]string{"sejarah", "visi", "misi::text", "geografis", "demografi::text"}, terms) + `
			UNION ALL
			SELECT id, COALESCE(judul, 'Struktur Organisasi'), COALESCE(deskripsi, ''), 'struktur', '', COALESCE(gambar, ''), id::text
			FROM struktur_organisasi WHERE ` + buildSearchConditionFor([]string{"judul", "deskripsi"}, terms) + `
			UNION ALL
			SELECT id, COALESCE(nama_situs, 'Portal Pemkab Tulungagung'), CONCAT_WS(' ', tagline, deskripsi, email, telepon, alamat), 'identitas', '', COALESCE(logo, ''), id::text
			FROM identitas_situs WHERE ` + buildSearchConditionFor([]string{"nama_situs", "tagline", "deskripsi", "email", "telepon", "alamat"}, terms) + `
			UNION ALL
			SELECT id, 'Kontak Pemkab Tulungagung', CONCAT_WS(' ', alamat, telepon, email, jam_operasional), 'kontak', '', '', id::text
			FROM kontak WHERE ` + buildSearchConditionFor([]string{"alamat", "telepon", "email", "jam_operasional"}, terms) + `
		) results ORDER BY sort_date DESC NULLS LAST, id DESC LIMIT 100`

	// The first condition is built separately so all UNION branches use the same
	// parameter positions and remain fully parameterized.
	query, args = buildUnionSearchQuery(query, terms)
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Println("Error Global Search:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal melakukan pencarian"})
		return
	}
	defer rows.Close()
	result := emptySearchResult()
	for rows.Next() {
		var id int
		var judul, deskripsi, tipe, slug, url string
		if err := rows.Scan(&id, &judul, &deskripsi, &tipe, &slug, &url); err != nil {
			log.Println("Error Scan Global Search:", err)
			continue
		}
		item := map[string]interface{}{"id": id, "judul": judul, "deskripsi": deskripsi, "tipe": tipe, "slug": slug, "url": url}
		if list, ok := result[tipe].([]map[string]interface{}); ok {
			result[tipe] = append(list, item)
		} else {
			result["lainnya"] = append(result["lainnya"].([]map[string]interface{}), item)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("Error Loop Global Search:", err)
	}
	appendStaticSearchResults(result, terms)
	c.JSON(http.StatusOK, result)
}

func appendStaticSearchResults(result gin.H, terms []string) {
	staticPages := []map[string]interface{}{
		{"id": "static-profil", "judul": "Profil Kabupaten Tulungagung", "deskripsi": "Profil daerah, sejarah, visi, misi, geografis, demografi, struktur organisasi, dan pejabat daerah.", "tipe": "static", "slug": "profil", "url": "/profil", "search": "profil kabupaten tulungagung sejarah visi misi geografis demografi pemerintahan"},
		{"id": "static-layanan", "judul": "Layanan Publik", "deskripsi": "Kumpulan layanan publik Kabupaten Tulungagung.", "tipe": "static", "slug": "layanan", "url": "/layanan", "search": "layanan publik administrasi kependudukan ktp kk sip sim perizinan kesehatan"},
		{"id": "static-berita", "judul": "Informasi dan Berita", "deskripsi": "Berita dan informasi resmi Pemerintah Kabupaten Tulungagung.", "tipe": "static", "slug": "berita", "url": "/berita", "search": "berita informasi kabar kegiatan pemerintah"},
		{"id": "static-pemerintahan", "judul": "Pemerintahan Kabupaten Tulungagung", "deskripsi": "Informasi OPD dan pemerintahan Kabupaten Tulungagung.", "tipe": "static", "slug": "pemerintahan", "url": "/pemerintahan", "search": "pemerintahan opd organisasi perangkat daerah"},
		{"id": "static-kontak", "judul": "Kontak", "deskripsi": "Alamat, telepon, email, dan informasi kontak Pemerintah Kabupaten Tulungagung.", "tipe": "static", "slug": "kontak", "url": "/kontak", "search": "kontak alamat telepon email"},
	}
	for _, page := range staticPages {
		if !matchesSearchTerms(page["search"].(string), terms) {
			continue
		}
		delete(page, "search")
		result["lainnya"] = append(result["lainnya"].([]map[string]interface{}), page)
	}
}

func matchesSearchTerms(value string, terms []string) bool {
	words := normalizeSearchTerms(value)
	for _, term := range terms {
		matched := false
		for _, word := range words {
			if strings.Contains(word, term) || (len(term) >= 3 && levenshteinDistance(term, word) <= 1) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

func levenshteinDistance(left, right string) int {
	previous := make([]int, len(right)+1)
	for index := range previous {
		previous[index] = index
	}
	for leftIndex, leftRune := range left {
		current := make([]int, len(right)+1)
		current[0] = leftIndex + 1
		for rightIndex, rightRune := range right {
			cost := 0
			if leftRune != rightRune {
				cost = 1
			}
			current[rightIndex+1] = minInt(current[rightIndex]+1, previous[rightIndex+1]+1, previous[rightIndex]+cost)
		}
		previous = current
	}
	return previous[len(right)]
}

func minInt(values ...int) int {
	minimum := values[0]
	for _, value := range values[1:] {
		if value < minimum {
			minimum = value
		}
	}
	return minimum
}

func emptySearchResult() gin.H {
	keys := []string{"berita", "halaman", "opd", "darurat", "aplikasi", "pengumuman", "galeri", "agenda", "pejabat", "tautan", "kategori_layanan", "layanan", "layanan_populer", "menu", "profil_daerah", "struktur", "identitas", "kontak", "lainnya"}
	result := gin.H{}
	for _, key := range keys {
		result[key] = []map[string]interface{}{}
	}
	return result
}

func normalizeSearchTerms(value string) []string {
	value = strings.ToLower(strings.TrimSpace(value))
	var normalized strings.Builder
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			normalized.WriteRune(r)
		} else {
			normalized.WriteRune(' ')
		}
	}
	seen := map[string]bool{}
	terms := []string{}
	for _, term := range strings.Fields(normalized.String()) {
		if !seen[term] {
			seen[term] = true
			terms = append(terms, term)
		}
	}
	return terms
}

func buildSearchCondition(fields []string, terms []string) (string, []interface{}) {
	parts := make([]string, 0, len(terms))
	args := make([]interface{}, 0, len(terms))
	for index, term := range terms {
		placeholders := make([]string, 0, len(fields))
		for _, field := range fields {
			placeholders = append(placeholders, "COALESCE("+field+", '')")
		}
		parts = append(parts, "LOWER(CONCAT_WS(' ', "+strings.Join(placeholders, ", ")+")) LIKE $"+itoa(index+1))
		args = append(args, "%"+escapeLikeTerm(term)+"%")
	}
	return "(" + strings.Join(parts, " AND ") + ")", args
}

func buildSearchConditionFor(fields []string, terms []string) string {
	condition, _ := buildSearchCondition(fields, terms)
	return condition
}

func buildUnionSearchQuery(query string, terms []string) (string, []interface{}) {
	// Every UNION branch uses $1..$N. PostgreSQL accepts the same parameters
	// repeatedly, so one argument list serves the whole query.
	_, args := buildSearchCondition([]string{"judul", "deskripsi", "konten", "kategori", "penulis"}, terms)
	return query, args
}

func escapeLikeTerm(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "%", `\%`)
	return strings.ReplaceAll(value, "_", `\_`)
}

func itoa(value int) string {
	return fmt.Sprintf("%d", value)
}

// GET /api/pencarian-populer
func GetPencarianPopuler(c *gin.Context) {
	rows, err := config.DB.Query("SELECT keyword FROM pencarian_populer ORDER BY urutan ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	keywords := []string{}
	for rows.Next() {
		var kw string
		rows.Scan(&kw)
		keywords = append(keywords, kw)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error Loop Pencarian Populer", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan saat membaca data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"keywords": keywords})
}

// GET /api/kategori-berita
func GetKategoriBerita(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, nama, COALESCE(icon, 'i-tabler-news'), COALESCE(deskripsi, '') FROM kategori_berita ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	type Kategori struct {
		ID        int    `json:"id"`
		Nama      string `json:"nama"`
		Icon      string `json:"icon"`
		Deskripsi string `json:"deskripsi"`
	}

	list := []Kategori{}
	for rows.Next() {
		var k Kategori
		rows.Scan(&k.ID, &k.Nama, &k.Icon, &k.Deskripsi)
		list = append(list, k)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error Loop Pencarian Berita", err)
		c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "Terjadi kesalahan saat membaca data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"kategori": list})
}
