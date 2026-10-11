<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Kelola Galeri'
})

// ==== Data Mock ====
const galeri = ref([
    {
        id: 1,
        judul: 'Peresmian Jembatan Desa',
        tipe: 'foto',
        tanggal: '2026-10-05',
        thumbnail: 'https://picsum.photos/seed/1/400/300',
        status: 'Publish'
    },
    {
        id: 2,
        judul: 'Festival Budaya 2026',
        tipe: 'foto',
        tanggal: '2026-10-07',
        thumbnail: 'https://picsum.photos/seed/2/400/300',
        status: 'Publish'
    },
    {
        id: 3,
        judul: 'Pelatihan UMKM Digital',
        tipe: 'video',
        tanggal: '2026-10-04',
        thumbnail: 'https://picsum.photos/seed/3/400/300',
        status: 'Publish'
    },
    {
        id: 4,
        judul: 'Penghijauan Pesisir',
        tipe: 'foto',
        tanggal: '2026-10-03',
        thumbnail: 'https://picsum.photos/seed/4/400/300',
        status: 'Publish'
    },
    {
        id: 5,
        judul: 'Rapat Paripurna DPRD',
        tipe: 'foto',
        tanggal: '2026-10-02',
        thumbnail: 'https://picsum.photos/seed/5/400/300',
        status: 'Draft'
    },
    {
        id: 6,
        judul: 'Kunjungan Kerja Bupati',
        tipe: 'video',
        tanggal: '2026-10-01',
        thumbnail: 'https://picsum.photos/seed/6/400/300',
        status: 'Publish'
    }
])

// ==== State UI ====
const search = ref('')
const selectedTipe = ref('Semua')
const selectedStatus = ref('Semua')
const viewMode = ref('grid') // 'grid' | 'list'
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)
const page = ref(1)

// ==== Opsi ====
const tipeOptions = ['Semua', 'foto', 'video']
const statusOptions = ['Semua', 'Publish', 'Draft']

// ==== Form ====
const emptyForm = () => ({
    id: null,
    judul: '',
    tipe: 'foto',
    tanggal: new Date().toISOString().split('T')[0],
    status: 'Draft',
    thumbnail: '',
    fileUrl: ''
})
const form = ref(emptyForm())

// ==== Computed ====
const filteredGaleri = computed(() => {
    return galeri.value.filter(item => {
        const matchSearch = item.judul.toLowerCase().includes(search.value.toLowerCase())
        const matchTipe = selectedTipe.value === 'Semua' || item.tipe === selectedTipe.value
        const matchStatus = selectedStatus.value === 'Semua' || item.status === selectedStatus.value
        return matchSearch && matchTipe && matchStatus
    })
})

// ==== CRUD ====
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

const saveGaleri = () => {
    if (!form.value.judul.trim()) return

    if (editingItem.value) {
        const index = galeri.value.findIndex(g => g.id === editingItem.value.id)
        if (index !== -1) {
            galeri.value[index] = { ...galeri.value[index], ...form.value }
        }
    } else {
        const newId = Math.max(...galeri.value.map(g => g.id), 0) + 1
        galeri.value.unshift({
            ...form.value,
            id: newId,
            thumbnail: form.value.thumbnail || `https://picsum.photos/seed/${newId}/400/300`
        })
    }
    closeModal()
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

const confirmDelete = () => {
    if (deletingItem.value) {
        galeri.value = galeri.value.filter(g => g.id !== deletingItem.value.id)
    }
    showDeleteModal.value = false
    deletingItem.value = null
}

const closeModal = () => {
    showModal.value = false
    editingItem.value = null
    form.value = emptyForm()
}

// ==== Helper ====
const statusColor = (status) => status === 'Publish' ? 'success' : 'warning'
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <UCard>
            <div class="flex items-center justify-between flex-wrap gap-4">
                <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Kelola Galeri</h2>
                    <p class="text-slate-500 mt-1">Kelola dokumentasi foto dan video kegiatan Pemkab Tulungagung.</p>
                </div>
                <UButton icon="i-tabler-plus" label="Upload Media" color="primary" size="lg" @click="openCreateModal" />
            </div>
        </UCard>

        <!-- Filter -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari judul..." size="lg"
                    class="md:col-span-2" />
                <USelect v-model="selectedTipe" :items="tipeOptions" icon="i-tabler-photo" size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <!-- Grid Galeri -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-photo" class="text-sky-500" />
                        Daftar Media
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredGaleri.length }} media ditemukan</span>
                </div>
            </template>

            <!-- Empty State -->
            <div v-if="filteredGaleri.length === 0" class="text-center py-12">
                <UIcon name="i-tabler-photo-off" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                <p class="text-slate-500 font-medium">Tidak ada media ditemukan</p>
                <p class="text-xs text-slate-400 mt-1">Coba ubah filter atau upload media baru.</p>
            </div>

            <!-- Grid -->
            <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                <div v-for="item in filteredGaleri" :key="item.id"
                    class="group relative bg-white dark:bg-slate-800 rounded-xl overflow-hidden border border-slate-200 dark:border-slate-700 hover:border-sky-400 hover:shadow-lg transition-all">
                    <!-- Thumbnail -->
                    <div class="relative aspect-video bg-slate-100 dark:bg-slate-900 overflow-hidden">
                        <img :src="item.thumbnail" :alt="item.judul"
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />

                        <!-- Overlay saat hover -->
                        <div
                            class="absolute inset-0 bg-slate-900/70 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                            <AdminButtonEdit size="sm" @click="openEditModal(item)" />
                            <AdminButtonDelete size="sm" @click="openDeleteModal(item)" />
                        </div>

                        <!-- Badge Tipe -->
                        <div class="absolute top-2 left-2 flex gap-1">
                            <UBadge :color="item.tipe === 'video' ? 'error' : 'primary'" variant="solid" size="xs">
                                <UIcon :name="item.tipe === 'video' ? 'i-tabler-video' : 'i-tabler-photo'"
                                    class="w-3 h-3 mr-1" />
                                {{ item.tipe.toUpperCase() }}
                            </UBadge>
                        </div>

                        <!-- Badge Status -->
                        <div class="absolute top-2 right-2">
                            <UBadge :color="statusColor(item.status)" variant="solid" size="xs">
                                {{ item.status }}
                            </UBadge>
                        </div>
                    </div>

                    <!-- Info -->
                    <div class="p-3">
                        <h4
                            class="font-semibold text-sm text-slate-900 dark:text-white line-clamp-2 group-hover:text-sky-600 transition-colors">
                            {{ item.judul }}
                        </h4>
                        <p class="text-xs text-slate-400 mt-1 flex items-center gap-1">
                            <UIcon name="i-tabler-calendar" class="w-3 h-3" />
                            {{ item.tanggal }}
                        </p>
                    </div>
                </div>
            </div>

            <!-- Pagination -->
            <template #footer>
                <div class="flex items-center justify-between flex-wrap gap-2">
                    <p class="text-xs text-slate-400">
                        Menampilkan {{ filteredGaleri.length }} dari {{ filteredGaleri.length }} media
                    </p>
                    <UPagination v-model="page" :total="filteredGaleri.length" :page-count="12" />
                </div>
            </template>
        </UCard>

        <!-- Modal Form -->
        <UModal v-model:open="showModal" :title="editingItem ? 'Edit Media' : 'Upload Media Baru'"
            :description="editingItem ? 'Perbarui informasi media di bawah ini.' : 'Upload foto atau video kegiatan baru.'">
            <template #body>
                <div class="space-y-4">
                    <UFormField label="Judul Media" required>
                        <UInput v-model="form.judul" placeholder="Masukkan judul media..." size="lg" class="w-full" />
                    </UFormField>

                    <div class="grid grid-cols-2 gap-4">
                        <UFormField label="Tipe Media" required>
                            <USelect v-model="form.tipe" :items="[
                                { label: 'Foto', value: 'foto' },
                                { label: 'Video', value: 'video' }
                            ]" size="lg" class="w-full" />
                        </UFormField>
                        <UFormField label="Tanggal" required>
                            <UInput v-model="form.tanggal" type="date" size="lg" class="w-full" />
                        </UFormField>
                    </div>

                    <UFormField label="Status" required>
                        <USelect v-model="form.status" :items="['Draft', 'Publish']" size="lg" class="w-full" />
                    </UFormField>

                    <!-- Upload Area -->
                    <UFormField :label="form.tipe === 'video' ? 'Upload Video' : 'Upload Foto'">
                        <div
                            class="border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-6 text-center hover:border-sky-500 transition-colors cursor-pointer">
                            <UIcon name="i-tabler-cloud-upload" class="w-10 h-10 text-slate-400 mx-auto mb-2" />
                            <p class="text-sm text-slate-500 font-medium">Klik untuk memilih file</p>
                            <p class="text-xs text-slate-400 mt-1">
                                {{ form.tipe === 'video' ? 'MP4, WebM (Max 50MB)' : 'PNG, JPG, WebP (Max 5MB)' }}
                            </p>
                        </div>
                    </UFormField>

                    <!-- Preview Thumbnail -->
                    <UFormField v-if="form.thumbnail" label="Preview">
                        <div class="rounded-lg overflow-hidden border border-slate-200 dark:border-slate-700">
                            <img :src="form.thumbnail" alt="Preview" class="w-full aspect-video object-cover" />
                        </div>
                    </UFormField>

                    <UFormField label="Deskripsi (opsional)">
                        <UTextarea v-model="form.fileUrl" placeholder="Tulis deskripsi singkat media ini..." :rows="3"
                            size="lg" class="w-full" />
                    </UFormField>
                </div>
            </template>

            <template #footer>
                <div class="flex items-center justify-end gap-2">
                    <UButton label="Batal" color="neutral" variant="outline" @click="closeModal" />
                    <AdminButtonSave :label="editingItem ? 'Simpan Perubahan' : 'Upload Media'" @click="saveGaleri" />
                </div>
            </template>
        </UModal>

        <!-- Modal Konfirmasi Hapus -->
        <UModal v-model:open="showDeleteModal" title="Konfirmasi Hapus"
            description="Tindakan ini tidak dapat dibatalkan.">
            <template #body>
                <div class="text-center py-4">
                    <div
                        class="w-16 h-16 rounded-full bg-red-100 dark:bg-red-950 flex items-center justify-center mx-auto mb-4">
                        <UIcon name="i-tabler-alert-triangle" class="w-8 h-8 text-red-500" />
                    </div>
                    <p class="text-slate-700 dark:text-slate-300">Apakah Anda yakin ingin menghapus media:</p>
                    <p class="font-bold text-slate-900 dark:text-white mt-2">"{{ deletingItem?.judul }}"</p>
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