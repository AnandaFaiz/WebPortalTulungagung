<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Kelola Pengumuman'
})

// ==== Data Mock ====
const pengumuman = ref([
    {
        id: 1,
        judul: 'Pengumuman Lelang Proyek Pembangunan Jembatan',
        tanggal: '2026-10-09',
        prioritas: 'Tinggi',
        status: 'Publish'
    },
    {
        id: 2,
        judul: 'Jadwal Pelayanan MPP Selama Bulan Oktober',
        tanggal: '2026-10-07',
        prioritas: 'Normal',
        status: 'Publish'
    },
    {
        id: 3,
        judul: 'Penerimaan CPNS Pemkab Tulungagung Tahun 2026',
        tanggal: '2026-10-05',
        prioritas: 'Tinggi',
        status: 'Publish'
    },
    {
        id: 4,
        judul: 'Draft: Sosialisasi Program Vaksinasi Booster',
        tanggal: '2026-10-03',
        prioritas: 'Normal',
        status: 'Draft'
    }
])

// ==== State UI ====
const search = ref('')
const selectedStatus = ref('Semua')
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)
const page = ref(1)

// ==== Opsi ====
const statusOptions = ['Semua', 'Publish', 'Draft']
const prioritasOptions = ['Normal', 'Tinggi']

// ==== Form ====
const emptyForm = () => ({
    id: null,
    judul: '',
    tanggal: new Date().toISOString().split('T')[0],
    prioritas: 'Normal',
    status: 'Draft',
    konten: ''
})
const form = ref(emptyForm())

// ==== Computed ====
const filteredPengumuman = computed(() => {
    return pengumuman.value.filter(item => {
        const matchSearch = item.judul.toLowerCase().includes(search.value.toLowerCase())
        const matchStatus = selectedStatus.value === 'Semua' || item.status === selectedStatus.value
        return matchSearch && matchStatus
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

const savePengumuman = () => {
    if (!form.value.judul.trim()) return

    if (editingItem.value) {
        const index = pengumuman.value.findIndex(p => p.id === editingItem.value.id)
        if (index !== -1) {
            pengumuman.value[index] = { ...pengumuman.value[index], ...form.value }
        }
    } else {
        const newId = Math.max(...pengumuman.value.map(p => p.id), 0) + 1
        pengumuman.value.unshift({ ...form.value, id: newId })
    }
    closeModal()
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

const confirmDelete = () => {
    if (deletingItem.value) {
        pengumuman.value = pengumuman.value.filter(p => p.id !== deletingItem.value.id)
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
const prioritasColor = (prioritas) => prioritas === 'Tinggi' ? 'error' : 'neutral'
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <UCard>
            <div class="flex items-center justify-between flex-wrap gap-4">
                <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Kelola Pengumuman</h2>
                    <p class="text-slate-500 mt-1">Kelola pengumuman resmi Pemkab Tulungagung.</p>
                </div>
                <UButton icon="i-tabler-plus" label="Tambah Pengumuman" color="primary" size="lg"
                    @click="openCreateModal" />
            </div>
        </UCard>

        <!-- Filter -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari judul pengumuman..." size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg"
                    placeholder="Status" />
            </div>
        </UCard>

        <!-- Tabel -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-speakerphone" class="text-sky-500" />
                        Daftar Pengumuman
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredPengumuman.length }} pengumuman ditemukan</span>
                </div>
            </template>

            <div v-if="filteredPengumuman.length === 0" class="text-center py-12">
                <UIcon name="i-tabler-inbox" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                <p class="text-slate-500 font-medium">Tidak ada pengumuman ditemukan</p>
                <p class="text-xs text-slate-400 mt-1">Coba ubah filter atau tambahkan pengumuman baru.</p>
            </div>

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Judul</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Tanggal</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Prioritas</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in filteredPengumuman" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <p class="font-medium text-slate-900 dark:text-white line-clamp-1 max-w-md">{{
                                    item.judul }}</p>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600 dark:text-slate-300">{{ item.tanggal }}</td>
                            <td class="py-3 px-4">
                                <UBadge :color="prioritasColor(item.prioritas)" variant="subtle" size="sm">
                                    {{ item.prioritas }}
                                </UBadge>
                            </td>
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

            <template #footer>
                <div class="flex items-center justify-between flex-wrap gap-2">
                    <p class="text-xs text-slate-400">
                        Menampilkan {{ filteredPengumuman.length }} dari {{ filteredPengumuman.length }} data
                    </p>
                    <UPagination v-model="page" :total="filteredPengumuman.length" :page-count="10" />
                </div>
            </template>
        </UCard>

        <!-- Modal Form -->
        <UModal v-model:open="showModal" :title="editingItem ? 'Edit Pengumuman' : 'Tambah Pengumuman Baru'"
            :description="editingItem ? 'Perbarui informasi pengumuman di bawah ini.' : 'Isi form berikut untuk menambahkan pengumuman baru.'">
            <template #body>
                <div class="space-y-4">
                    <UFormField label="Judul Pengumuman" required>
                        <UInput v-model="form.judul" placeholder="Masukkan judul pengumuman..." size="lg"
                            class="w-full" />
                    </UFormField>

                    <div class="grid grid-cols-2 gap-4">
                        <UFormField label="Tanggal" required>
                            <UInput v-model="form.tanggal" type="date" size="lg" class="w-full" />
                        </UFormField>
                        <UFormField label="Prioritas" required>
                            <USelect v-model="form.prioritas" :items="prioritasOptions" size="lg" class="w-full" />
                        </UFormField>
                    </div>

                    <UFormField label="Status" required>
                        <USelect v-model="form.status" :items="['Draft', 'Publish']" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Isi Pengumuman">
                        <UTextarea v-model="form.konten" placeholder="Tulis isi pengumuman di sini..." :rows="5"
                            size="lg" class="w-full" />
                    </UFormField>
                </div>
            </template>

            <template #footer>
                <div class="flex items-center justify-end gap-2">
                    <UButton label="Batal" color="neutral" variant="outline" @click="closeModal" />
                    <AdminButtonSave :label="editingItem ? 'Simpan Perubahan' : 'Simpan Pengumuman'"
                        @click="savePengumuman" />
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
                    <p class="text-slate-700 dark:text-slate-300">Apakah Anda yakin ingin menghapus pengumuman:</p>
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