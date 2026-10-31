<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Kelola Berita'
})

const toast = useToast()
const API_BASE = 'http://localhost:8080/api'

// ==== State ====
const berita = ref([])
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const selectedKategori = ref('Semua')
const selectedStatus = ref('Semua')
const page = ref(1)
const perPage = 10

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

const kategoriOptions = ['Semua', 'Pemerintahan', 'Ekonomi', 'Budaya', 'Pendidikan', 'Pembangunan', 'Lingkungan']
const statusOptions = ['Semua', 'Publish', 'Draft']

const emptyForm = () => ({
    id: null,
    judul: '',
    kategori: 'Pemerintahan',
    penulis: 'Admin Utama',
    tanggal: new Date().toISOString().split('T')[0],
    status: 'Draft',
    deskripsi: '',
    konten: '',
    thumbnail: ''
})
const form = ref(emptyForm())

// ==== Computed ====
const filteredBerita = computed(() => {
    return berita.value.filter(item => {
        const matchSearch = item.judul.toLowerCase().includes(search.value.toLowerCase())
        const matchKategori = selectedKategori.value === 'Semua' || item.kategori === selectedKategori.value
        const matchStatus = selectedStatus.value === 'Semua' || item.status === selectedStatus.value
        return matchSearch && matchKategori && matchStatus
    })
})

const paginatedBerita = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredBerita.value.slice(start, start + perPage)
})

watch([search, selectedKategori, selectedStatus], () => { page.value = 1 })

// ==== Fetch (Golang API) ====
const fetchBerita = async () => {
    loading.value = true
    try {
        const response = await $fetch(`${API_BASE}/admin/berita`, {
            credentials: 'include'
        })
        // ✅ PERBAIKAN: Tambahkan || [] agar tidak error jika data null
        berita.value = response.berita || []
    } catch (err) {
        toast.add({ title: 'Gagal memuat data', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        loading.value = false
    }
}

// ==== Modal Handlers ====
const openCreateModal = () => {
    form.value = emptyForm()
    editingItem.value = null
    showModal.value = true
}

const openEditModal = (item) => {
    form.value = {
        id: item.id,
        judul: item.judul,
        kategori: item.kategori,
        penulis: item.penulis || '',
        tanggal: item.tanggal ? item.tanggal.split('T')[0] : '',
        status: item.status,
        deskripsi: item.deskripsi || '',
        konten: item.konten || '',
        thumbnail: item.thumbnail || ''
    }
    editingItem.value = item
    showModal.value = true
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

// ==== Save (Golang API) ====
const saveBerita = async () => {
    if (!form.value.judul.trim()) {
        toast.add({ title: 'Validasi gagal', description: 'Judul wajib diisi', icon: 'i-tabler-alert-circle', color: 'error' })
        return
    }
    saving.value = true
    try {
        const payload = {
            judul: form.value.judul,
            kategori: form.value.kategori,
            penulis: form.value.penulis,
            tanggal: form.value.tanggal,
            status: form.value.status,
            deskripsi: form.value.deskripsi,
            konten: form.value.konten,
            thumbnail: form.value.thumbnail
        }

        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/berita/${editingItem.value.id}`, {
                method: 'PUT',
                body: payload,
                credentials: 'include'
            })
            toast.add({ title: 'Berhasil', description: 'Berita diperbarui', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        } else {
            await $fetch(`${API_BASE}/admin/berita`, {
                method: 'POST',
                body: payload,
                credentials: 'include'
            })
            toast.add({ title: 'Berhasil', description: 'Berita ditambahkan', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        }
        showModal.value = false
        await fetchBerita()
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        saving.value = false
    }
}

// ==== Delete (Golang API) ====
const confirmDelete = async () => {
    if (!deletingItem.value) return
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/berita/${deletingItem.value.id}`, {
            method: 'DELETE',
            credentials: 'include'
        })
        toast.add({ title: 'Berhasil', description: 'Berita dihapus', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchBerita()
    } catch (err) {
        toast.add({ title: 'Gagal menghapus', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        saving.value = false
    }
}

const closeModal = () => {
    showModal.value = false
    editingItem.value = null
    form.value = emptyForm()
}

const statusColor = (status) => status === 'Publish' ? 'success' : 'warning'

onMounted(() => { fetchBerita() })
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Kelola Berita" description="Kelola semua berita dan informasi resmi Pemkab Tulungagung."
            action-label="Tambah Berita" @action="openCreateModal" />

        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari judul berita..." size="lg" />
                <USelect v-model="selectedKategori" :items="kategoriOptions" icon="i-tabler-category" size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-news" class="text-sky-500" />
                        Daftar Berita
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredBerita.length }} berita</span>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat berita...</p>
            </div>

            <AdminEmptyState v-else-if="filteredBerita.length === 0" icon="i-tabler-inbox" title="Belum ada berita"
                description="Klik Tambah Berita untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Judul</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Kategori</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Penulis</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Tanggal</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in paginatedBerita" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <p class="font-medium text-slate-900 dark:text-white line-clamp-1 max-w-md">{{
                                    item.judul }}</p>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge color="primary" variant="subtle" size="sm">{{ item.kategori }}</UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-2">
                                    <div
                                        class="w-7 h-7 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 dark:text-sky-400 text-xs font-bold">
                                        {{ (item.penulis || 'A').charAt(0).toUpperCase() }}
                                    </div>
                                    <span class="text-sm text-slate-600 dark:text-slate-300">{{ item.penulis || '-'
                                        }}</span>
                                </div>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600 dark:text-slate-300">
                                {{ item.tanggal ? item.tanggal.split('T')[0] : '-' }}
                            </td>
                            <td class="py-3 px-4">
                                <UBadge :color="statusColor(item.status)" variant="subtle" size="sm">{{ item.status }}
                                </UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center justify-end gap-2">
                                    <AdminButtonEdit size="xs" @click="openEditModal(item)" />
                                    <AdminButtonDelete size="xs" @click="openDeleteModal(item)" />
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <template #footer>
                <AdminPagination v-model:page="page" :total="filteredBerita.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Berita' : 'Tambah Berita Baru'"
            :description="editingItem ? 'Perbarui informasi berita.' : 'Isi form berikut untuk menambahkan berita baru.'"
            :is-edit="!!editingItem" :loading="saving"
            :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan Berita'" @submit="saveBerita"
            @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Judul Berita" required>
                    <UInput v-model="form.judul" placeholder="Masukkan judul..." size="lg" class="w-full" />
                </UFormField>
                <div class="grid grid-cols-2 gap-4">
                    <UFormField label="Kategori" required>
                        <USelect v-model="form.kategori" :items="kategoriOptions.filter(k => k !== 'Semua')" size="lg"
                            class="w-full" />
                    </UFormField>
                    <UFormField label="Status" required>
                        <USelect v-model="form.status" :items="['Draft', 'Publish']" size="lg" class="w-full" />
                    </UFormField>
                </div>
                <div class="grid grid-cols-2 gap-4">
                    <UFormField label="Penulis">
                        <UInput v-model="form.penulis" placeholder="Nama penulis" size="lg" class="w-full" />
                    </UFormField>
                    <UFormField label="Tanggal">
                        <UInput v-model="form.tanggal" type="date" size="lg" class="w-full" />
                    </UFormField>
                </div>
                <UFormField label="Deskripsi Singkat">
                    <UTextarea v-model="form.deskripsi" placeholder="Deskripsi singkat berita..." :rows="2" size="lg"
                        class="w-full" />
                </UFormField>
                <UFormField label="Konten Berita">
                    <UTextarea v-model="form.konten" placeholder="Tulis konten berita..." :rows="6" size="lg"
                        class="w-full" />
                </UFormField>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.judul" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>