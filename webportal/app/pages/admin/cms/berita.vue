<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Kelola Berita'
})

// ==== Data Mock (nanti diganti dari API) ====
const berita = ref([
    {
        id: 1,
        judul: 'Pemkab Tulungagung Luncurkan Program Digitalisasi Layanan Desa',
        kategori: 'Pemerintahan',
        penulis: 'Admin Utama',
        tanggal: '2026-10-08',
        status: 'Publish'
    },
    {
        id: 2,
        judul: 'Festival Budaya Tulungagung 2026 Resmi Dibuka',
        kategori: 'Budaya',
        penulis: 'Editor Budaya',
        tanggal: '2026-10-07',
        status: 'Publish'
    },
    {
        id: 3,
        judul: 'Realisasi Investasi Semester I Capai Rp652 Miliar',
        kategori: 'Ekonomi',
        penulis: 'Admin Utama',
        tanggal: '2026-10-06',
        status: 'Publish'
    },
    {
        id: 4,
        judul: 'Draft: Program Beasiswa Pelajar Berprestasi',
        kategori: 'Pendidikan',
        penulis: 'Editor Pendidikan',
        tanggal: '2026-10-05',
        status: 'Draft'
    }
])

// ==== State UI ====
const search = ref('')
const selectedKategori = ref('Semua')
const selectedStatus = ref('Semua')
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)
const page = ref(1)

// ==== Opsi Filter ====
const kategoriOptions = ['Semua', 'Pemerintahan', 'Ekonomi', 'Budaya', 'Pendidikan', 'Pembangunan', 'Lingkungan']
const statusOptions = ['Semua', 'Publish', 'Draft']

// ==== Form Data ====
const emptyForm = () => ({
    id: null,
    judul: '',
    kategori: 'Pemerintahan',
    penulis: 'Admin Utama',
    tanggal: new Date().toISOString().split('T')[0],
    status: 'Draft',
    konten: ''
})
const form = ref(emptyForm())

// ==== Computed: Filter & Search ====
const filteredBerita = computed(() => {
    return berita.value.filter(item => {
        const matchSearch = item.judul.toLowerCase().includes(search.value.toLowerCase())
        const matchKategori = selectedKategori.value === 'Semua' || item.kategori === selectedKategori.value
        const matchStatus = selectedStatus.value === 'Semua' || item.status === selectedStatus.value
        return matchSearch && matchKategori && matchStatus
    })
})

// ==== Aksi CRUD ====
const openCreateModal = () => {
    form.value = emptyForm()
    editingItem.value = null
    showModal.value = true
}

const openEditModal = (item) => {
    form.value = { ...item }
    editingItem.value = item
    showModal.value = true
}

const saveBerita = () => {
    if (!form.value.judul.trim()) return

    if (editingItem.value) {
        // Update
        const index = berita.value.findIndex(b => b.id === editingItem.value.id)
        if (index !== -1) {
            berita.value[index] = { ...berita.value[index], ...form.value }
        }
    } else {
        // Create
        const newId = Math.max(...berita.value.map(b => b.id), 0) + 1
        berita.value.unshift({ ...form.value, id: newId })
    }
    closeModal()
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

const confirmDelete = () => {
    if (deletingItem.value) {
        berita.value = berita.value.filter(b => b.id !== deletingItem.value.id)
    }
    showDeleteModal.value = false
    deletingItem.value = null
}

const closeModal = () => {
    showModal.value = false
    editingItem.value = null
    form.value = emptyForm()
}

// ==== Helper: Warna Badge Status ====
const statusColor = (status) => {
    return status === 'Publish' ? 'success' : 'warning'
}
</script>

<template>
    <div class="space-y-6">

        <!-- Header Card -->
        <UCard>
            <div class="flex items-center justify-between flex-wrap gap-4">
                <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Kelola Berita</h2>
                    <p class="text-slate-500 mt-1">Kelola semua berita dan informasi resmi Pemkab Tulungagung.</p>
                </div>
                <UButton icon="i-tabler-plus" label="Tambah Berita" color="primary" size="lg"
                    @click="openCreateModal" />
            </div>
        </UCard>

        <!-- Filter & Search -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari judul berita..." size="lg" />
                <USelect v-model="selectedKategori" :items="kategoriOptions" icon="i-tabler-category" size="lg"
                    placeholder="Kategori" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg"
                    placeholder="Status" />
            </div>
        </UCard>

        <!-- Tabel Berita -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-news" class="text-sky-500" />
                        Daftar Berita
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredBerita.length }} berita ditemukan</span>
                </div>
            </template>

            <!-- Empty State -->
            <div v-if="filteredBerita.length === 0" class="text-center py-12">
                <UIcon name="i-tabler-inbox" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                <p class="text-slate-500 font-medium">Tidak ada berita ditemukan</p>
                <p class="text-xs text-slate-400 mt-1">Coba ubah filter atau tambahkan berita baru.</p>
            </div>

            <!-- Tabel -->
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
                        <tr v-for="item in filteredBerita" :key="item.id"
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
                                        {{ item.penulis.charAt(0) }}
                                    </div>
                                    <span class="text-sm text-slate-600 dark:text-slate-300">{{ item.penulis }}</span>
                                </div>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600 dark:text-slate-300">{{ item.tanggal }}</td>
                            <td class="py-3 px-4">
                                <UBadge :color="statusColor(item.status)" variant="subtle" size="sm">
                                    {{ item.status }}
                                </UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center justify-end gap-2">
                                    <AdminButtonEdit size="xs" aria-label="Edit" @click="openEditModal(item)" />
                                    <AdminButtonDelete size="xs" aria-label="Hapus" @click="openDeleteModal(item)" />
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Pagination -->
            <template #footer>
                <div class="flex items-center justify-between flex-wrap gap-2">
                    <p class="text-xs text-slate-400">
                        Menampilkan {{ filteredBerita.length }} dari {{ filteredBerita.length }} data
                    </p>
                    <UPagination v-model="page" :total="filteredBerita.length" :page-count="10" />
                </div>
            </template>
        </UCard>

        <!-- ===== Modal Form (Create/Edit) ===== -->
        <UModal v-model:open="showModal" :title="editingItem ? 'Edit Berita' : 'Tambah Berita Baru'"
            :description="editingItem ? 'Perbarui informasi berita di bawah ini.' : 'Isi form berikut untuk menambahkan berita baru.'">
            <template #body>
                <div class="space-y-4">
                    <UFormField label="Judul Berita" required>
                        <UInput v-model="form.judul" placeholder="Masukkan judul berita..." size="lg" class="w-full" />
                    </UFormField>

                    <div class="grid grid-cols-2 gap-4">
                        <UFormField label="Kategori" required>
                            <USelect v-model="form.kategori" :items="kategoriOptions.filter(k => k !== 'Semua')"
                                size="lg" class="w-full" />
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

                    <UFormField label="Thumbnail">
                        <div
                            class="border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-6 text-center hover:border-sky-500 transition-colors cursor-pointer">
                            <UIcon name="i-tabler-cloud-upload" class="w-8 h-8 text-slate-400 mx-auto mb-2" />
                            <p class="text-sm text-slate-500">Klik untuk upload gambar</p>
                            <p class="text-xs text-slate-400 mt-1">PNG, JPG, WebP (Max 5MB)</p>
                        </div>
                    </UFormField>

                    <UFormField label="Konten Berita">
                        <UTextarea v-model="form.konten" placeholder="Tulis konten berita di sini..." :rows="6"
                            size="lg" class="w-full" />
                    </UFormField>
                </div>
            </template>

            <template #footer>
                <div class="flex items-center justify-end gap-2">
                    <UButton label="Batal" color="neutral" variant="outline" @click="closeModal" />
                    <AdminButtonSave :label="editingItem ? 'Simpan Perubahan' : 'Simpan Berita'" @click="saveBerita" />
                </div>
            </template>
        </UModal>

        <!-- ===== Modal Konfirmasi Hapus ===== -->
        <UModal v-model:open="showDeleteModal" title="Konfirmasi Hapus"
            description="Tindakan ini tidak dapat dibatalkan.">
            <template #body>
                <div class="text-center py-4">
                    <div
                        class="w-16 h-16 rounded-full bg-red-100 dark:bg-red-950 flex items-center justify-center mx-auto mb-4">
                        <UIcon name="i-tabler-alert-triangle" class="w-8 h-8 text-red-500" />
                    </div>
                    <p class="text-slate-700 dark:text-slate-300">
                        Apakah Anda yakin ingin menghapus berita:
                    </p>
                    <p class="font-bold text-slate-900 dark:text-white mt-2">
                        "{{ deletingItem?.judul }}"
                    </p>
                </div>
            </template>

            <template #footer>
                <div class="flex items-center justify-end gap-2">
                    <UButton label="Batal" color="neutral" variant="outline" @click="showDeleteModal = false" />
                    <AdminButtonDelete label="Hapus" @click="confirmDelete" />
                </div>
            </template>
        </UModal>

    </div>
</template>