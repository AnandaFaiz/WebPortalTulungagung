<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Kelola Galeri'
})

const toast = useToast()

// ==== State ====
const galeri = ref([])
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const selectedTipe = ref('Semua')
const selectedStatus = ref('Semua')
const page = ref(1)
const perPage = 12

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

const tipeOptions = ['Semua', 'foto', 'video']
const statusOptions = ['Semua', 'Publish', 'Draft']

const emptyForm = () => ({
    id: null,
    judul: '',
    deskripsi: '',
    tipe: 'foto',
    tanggal: new Date().toISOString().split('T')[0],
    status: 'Draft',
    thumbnail: '',
    file_url: ''
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

const paginatedGaleri = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredGaleri.value.slice(start, start + perPage)
})

watch([search, selectedTipe, selectedStatus], () => { page.value = 1 })

// ==== Fetch ====
const fetchGaleri = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/galeri', { credentials: 'include' }
        )
        galeri.value = response.galeri || []
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
        deskripsi: item.deskripsi || '',
        tipe: item.tipe,
        tanggal: item.tanggal ? item.tanggal.split('T')[0] : '',
        status: item.status,
        thumbnail: item.thumbnail || '',
        file_url: item.file_url || ''
    }
    editingItem.value = item
    showModal.value = true
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

// ==== Save ====
const saveGaleri = async () => {
    if (!form.value.judul.trim()) {
        toast.add({ title: 'Validasi gagal', description: 'Judul wajib diisi', icon: 'i-tabler-alert-circle', color: 'error' })
        return
    }
    saving.value = true
    try {
        const requestFetch = useRequestFetch()
        const payload = {
            judul: form.value.judul,
            deskripsi: form.value.deskripsi,
            tipe: form.value.tipe,
            tanggal: form.value.tanggal,
            status: form.value.status,
            thumbnail: form.value.thumbnail,
            file_url: form.value.file_url
        }

        if (editingItem.value) {
            await requestFetch(`http://localhost:8080/api/admin/galeri/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Galeri diperbarui', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        } else {
            await requestFetch('http://localhost:8080/api/admin/galeri', { method: 'POST', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Galeri ditambahkan', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        }
        showModal.value = false
        await fetchGaleri()
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        saving.value = false
    }
}

// ==== Delete ====
const confirmDelete = async () => {
    if (!deletingItem.value) return
    saving.value = true
    try {
        const requestFetch = useRequestFetch()
        await requestFetch(`http://localhost:8080/api/admin/galeri/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({ title: 'Berhasil', description: 'Galeri dihapus', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchGaleri()
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

onMounted(() => { fetchGaleri() })
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Kelola Galeri" description="Kelola dokumentasi foto dan video kegiatan."
            action-label="Upload Media" @action="openCreateModal" />

        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari judul..." size="lg"
                    class="md:col-span-2" />
                <USelect v-model="selectedTipe" :items="tipeOptions" icon="i-tabler-photo" size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-photo" class="text-sky-500" />
                        Daftar Media
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredGaleri.length }} media</span>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data...</p>
            </div>

            <AdminEmptyState v-else-if="filteredGaleri.length === 0" icon="i-tabler-photo-off" title="Belum ada media"
                description="Klik Upload Media untuk memulai." />

            <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                <div v-for="item in paginatedGaleri" :key="item.id"
                    class="group relative bg-white dark:bg-slate-800 rounded-xl overflow-hidden border border-slate-200 dark:border-slate-700 hover:border-sky-400 hover:shadow-lg transition-all">
                    <div class="relative aspect-video bg-slate-100 dark:bg-slate-900 overflow-hidden">
                        <img v-if="item.thumbnail" :src="item.thumbnail" :alt="item.judul"
                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
                        <div v-else class="w-full h-full flex items-center justify-center">
                            <UIcon :name="item.tipe === 'video' ? 'i-tabler-video' : 'i-tabler-photo'"
                                class="w-12 h-12 text-slate-300" />
                        </div>

                        <div
                            class="absolute inset-0 bg-slate-900/70 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                            <AdminButtonEdit size="sm" @click="openEditModal(item)" />
                            <AdminButtonDelete size="sm" @click="openDeleteModal(item)" />
                        </div>

                        <div class="absolute top-2 left-2">
                            <UBadge :color="item.tipe === 'video' ? 'error' : 'primary'" variant="solid" size="xs">
                                <UIcon :name="item.tipe === 'video' ? 'i-tabler-video' : 'i-tabler-photo'"
                                    class="w-3 h-3 mr-1" />
                                {{ item.tipe.toUpperCase() }}
                            </UBadge>
                        </div>

                        <div class="absolute top-2 right-2">
                            <UBadge :color="statusColor(item.status)" variant="solid" size="xs">{{ item.status }}
                            </UBadge>
                        </div>
                    </div>

                    <div class="p-3">
                        <h4
                            class="font-semibold text-sm text-slate-900 dark:text-white line-clamp-2 group-hover:text-sky-600 transition-colors">
                            {{ item.judul }}</h4>
                        <p class="text-xs text-slate-400 mt-1 flex items-center gap-1">
                            <UIcon name="i-tabler-calendar" class="w-3 h-3" />
                            {{ item.tanggal ? item.tanggal.split('T')[0] : '-' }}
                        </p>
                    </div>
                </div>
            </div>

            <template #footer>
                <AdminPagination v-model:page="page" :total="filteredGaleri.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Media' : 'Upload Media Baru'"
            :description="editingItem ? 'Perbarui informasi media.' : 'Upload foto atau video baru.'"
            :is-edit="!!editingItem" :loading="saving" :submit-label="editingItem ? 'Simpan Perubahan' : 'Upload Media'"
            @submit="saveGaleri" @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Judul Media" required>
                    <UInput v-model="form.judul" placeholder="Masukkan judul..." size="lg" class="w-full" />
                </UFormField>
                <div class="grid grid-cols-2 gap-4">
                    <UFormField label="Tipe Media" required>
                        <USelect v-model="form.tipe"
                            :items="[{ label: 'Foto', value: 'foto' }, { label: 'Video', value: 'video' }]" size="lg"
                            class="w-full" />
                    </UFormField>
                    <UFormField label="Tanggal" required>
                        <UInput v-model="form.tanggal" type="date" size="lg" class="w-full" />
                    </UFormField>
                </div>
                <UFormField label="Status" required>
                    <USelect v-model="form.status" :items="['Draft', 'Publish']" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="URL Thumbnail / Gambar" hint="URL langsung ke gambar.">
                    <UInput v-model="form.thumbnail" placeholder="https://..." icon="i-tabler-link" size="lg"
                        class="w-full" />
                </UFormField>
                <UFormField label="URL File (Video)" hint="Kosongkan jika tipe foto.">
                    <UInput v-model="form.file_url" placeholder="https://..." icon="i-tabler-link" size="lg"
                        class="w-full" />
                </UFormField>
                <UFormField label="Deskripsi">
                    <UTextarea v-model="form.deskripsi" placeholder="Deskripsi singkat..." :rows="3" size="lg"
                        class="w-full" />
                </UFormField>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.judul" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>