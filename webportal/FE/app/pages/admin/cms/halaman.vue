<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Halaman Statis'
})

const toast = useToast()
const { can } = useAuth()
const canUpdate = computed(() => can('cms.halaman', 'update'))
const canDelete = computed(() => can('cms.halaman', 'delete'))

// ==== State ====
const halaman = ref([])
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const selectedStatus = ref('Semua')
const page = ref(1)
const perPage = 10

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

const statusOptions = ['Semua', 'Publish', 'Draft']

const emptyForm = () => ({
    id: null,
    judul: '',
    slug: '',
    tanggal: new Date().toISOString().split('T')[0],
    status: 'Draft',
    konten: ''
})
const form = ref(emptyForm())

// ==== Auto-generate slug ====
const generateSlug = (text) => {
    return text.toLowerCase().trim().replace(/[^\w\s-]/g, '').replace(/[\s_-]+/g, '-').replace(/^-+|-+$/g, '')
}

watch(() => form.value.judul, (newVal) => {
    if (!editingItem.value) {
        form.value.slug = generateSlug(newVal)
    }
})

// ==== Computed ====
const filteredHalaman = computed(() => {
    const list = halaman.value || []
    return list.filter(item => {
        const matchSearch = item.judul.toLowerCase().includes(search.value.toLowerCase()) || item.slug.toLowerCase().includes(search.value.toLowerCase())
        const matchStatus = selectedStatus.value === 'Semua' || item.status === selectedStatus.value
        return matchSearch && matchStatus
    })
})

const paginatedHalaman = computed(() => {
    const list = filteredHalaman.value || []
    const start = (page.value - 1) * perPage
    return list.slice(start, start + perPage)
})

watch([search, selectedStatus], () => { page.value = 1 })

// ==== Fetch ====
const fetchHalaman = async () => {
    loading.value = true
    try {
        const requestFetch = useRequestFetch()
        const response = await requestFetch('http://localhost:8080/api/admin/halaman', { credentials: 'include' })
        halaman.value = response.halaman || []
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
        slug: item.slug,
        tanggal: item.tanggal ? item.tanggal.split('T')[0] : '',
        status: item.status,
        konten: item.konten || ''
    }
    editingItem.value = item
    showModal.value = true
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

// ==== Save ====
const saveHalaman = async () => {
    if (!form.value.judul.trim() || !form.value.slug.trim()) {
        toast.add({ title: 'Validasi gagal', description: 'Judul dan slug wajib diisi', icon: 'i-tabler-alert-circle', color: 'error' })
        return
    }
    saving.value = true
    try {
        const requestFetch = useRequestFetch()
        const payload = {
            judul: form.value.judul,
            slug: form.value.slug,
            tanggal: form.value.tanggal,
            status: form.value.status,
            konten: form.value.konten
        }

        if (editingItem.value) {
            await requestFetch(`http://localhost:8080/api/admin/halaman/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Halaman diperbarui', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        } else {
            await requestFetch('http://localhost:8080/api/admin/halaman', { method: 'POST', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Halaman ditambahkan', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        }
        showModal.value = false
        await fetchHalaman()
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
        await requestFetch(`http://localhost:8080/api/admin/halaman/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({ title: 'Berhasil', description: 'Halaman dihapus', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchHalaman()
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

onMounted(() => { fetchHalaman() })
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Halaman Statis"
            description="Kelola halaman statis seperti Tentang, Kebijakan Privasi, FAQ, dll."
            action-label="Tambah Halaman" @action="openCreateModal" />

        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari judul atau slug..." size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-file-description" class="text-sky-500" />
                        Daftar Halaman
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredHalaman.length || 0 }} halaman</span>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data...</p>
            </div>

            <AdminEmptyState v-else-if="!filteredHalaman || filteredHalaman.length === 0" icon="i-tabler-file-off"
                title="Belum ada halaman" description="Klik Tambah Halaman untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Judul</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Slug</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Tanggal</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Status</th>
                            <th v-if="canUpdate || canDelete"
                                class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in paginatedHalaman" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <p class="font-medium text-slate-900 dark:text-white line-clamp-1 max-w-md">{{
                                    item.judul }}</p>
                            </td>
                            <td class="py-3 px-4">
                                <code
                                    class="text-xs bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 px-2 py-1 rounded">/{{
                                        item.slug }}</code>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600 dark:text-slate-300">{{ item.tanggal ?
                                item.tanggal.split('T')[0] : '-' }}</td>
                            <td class="py-3 px-4">
                                <UBadge :color="statusColor(item.status)" variant="subtle" size="sm">{{ item.status }}
                                </UBadge>
                            </td>
                            <td v-if="canUpdate || canDelete" class="py-3 px-4">
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
                <AdminPagination v-model:page="page" :total="filteredHalaman.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Halaman' : 'Tambah Halaman Baru'"
            :description="editingItem ? 'Perbarui informasi halaman.' : 'Isi form berikut.'" :is-edit="!!editingItem"
            :loading="saving" :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan Halaman'" @submit="saveHalaman"
            @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Judul Halaman" required>
                    <UInput v-model="form.judul" placeholder="Contoh: Tentang Portal" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Slug (URL)" required>
                    <UInput v-model="form.slug" placeholder="tentang-portal" size="lg" class="w-full">
                        <template #leading>
                            <span class="text-slate-400 text-sm">/</span>
                        </template>
                    </UInput>
                </UFormField>
                <div class="grid grid-cols-2 gap-4">
                    <UFormField label="Tanggal" required>
                        <UInput v-model="form.tanggal" type="date" size="lg" class="w-full" />
                    </UFormField>
                    <UFormField label="Status" required>
                        <USelect v-model="form.status" :items="['Draft', 'Publish']" size="lg" class="w-full" />
                    </UFormField>
                </div>
                <UFormField label="Konten Halaman">
                    <UTextarea v-model="form.konten" placeholder="Tulis konten halaman..." :rows="8" size="lg"
                        class="w-full" />
                </UFormField>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.judul" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>