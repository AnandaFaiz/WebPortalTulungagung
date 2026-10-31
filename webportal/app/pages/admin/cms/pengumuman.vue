<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Kelola Pengumuman'
})

const toast = useToast()

// ==== State ====
const pengumuman = ref([])
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
const prioritasOptions = ['Normal', 'Tinggi']

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

const paginatedPengumuman = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredPengumuman.value.slice(start, start + perPage)
})

watch([search, selectedStatus], () => { page.value = 1 })

// ==== Fetch ====
const fetchPengumuman = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/pengumuman',
            { credentials: 'include' }
        )
        pengumuman.value = response.pengumuman
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
        tanggal: item.tanggal ? item.tanggal.split('T')[0] : '',
        prioritas: item.prioritas,
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
const savePengumuman = async () => {
    if (!form.value.judul.trim()) {
        toast.add({ title: 'Validasi gagal', description: 'Judul wajib diisi', icon: 'i-tabler-alert-circle', color: 'error' })
        return
    }
    saving.value = true
    try {
        const requestFetch = useRequestFetch()
        const payload = {
            judul: form.value.judul,
            tanggal: form.value.tanggal,
            prioritas: form.value.prioritas,
            status: form.value.status,
            konten: form.value.konten
        }

        if (editingItem.value) {
            await $fetch(`http://localhost:8080/api/admin/pengumuman/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Pengumuman diperbarui', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        } else {
            await requestFetch('http://localhost:8080/api/admin/pengumuman', { method: 'POST', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Pengumuman ditambahkan', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        }
        showModal.value = false
        await fetchPengumuman()
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
        await $fetch(`http://localhost:8080/api/admin/pengumuman/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({ title: 'Berhasil', description: 'Pengumuman dihapus', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchPengumuman()
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
const prioritasColor = (prioritas) => prioritas === 'Tinggi' ? 'error' : 'neutral'

onMounted(() => { fetchPengumuman() })
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Kelola Pengumuman" description="Kelola pengumuman resmi Pemkab Tulungagung."
            action-label="Tambah Pengumuman" @action="openCreateModal" />

        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari judul pengumuman..." size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-speakerphone" class="text-sky-500" />
                        Daftar Pengumuman
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredPengumuman.length }} pengumuman</span>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data...</p>
            </div>

            <AdminEmptyState v-else-if="filteredPengumuman.length === 0" icon="i-tabler-inbox"
                title="Belum ada pengumuman" description="Klik Tambah Pengumuman untuk memulai." />

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
                        <tr v-for="item in paginatedPengumuman" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <p class="font-medium text-slate-900 dark:text-white line-clamp-1 max-w-md">{{
                                    item.judul }}</p>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600 dark:text-slate-300">{{ item.tanggal ?
                                item.tanggal.split('T')[0] : '-' }}</td>
                            <td class="py-3 px-4">
                                <UBadge :color="prioritasColor(item.prioritas)" variant="subtle" size="sm">{{
                                    item.prioritas }}
                                </UBadge>
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
                <AdminPagination v-model:page="page" :total="filteredPengumuman.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Pengumuman' : 'Tambah Pengumuman Baru'"
            :description="editingItem ? 'Perbarui informasi pengumuman.' : 'Isi form berikut.'" :is-edit="!!editingItem"
            :loading="saving" :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan Pengumuman'"
            @submit="savePengumuman" @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Judul Pengumuman" required>
                    <UInput v-model="form.judul" placeholder="Masukkan judul..." size="lg" class="w-full" />
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
                    <UTextarea v-model="form.konten" placeholder="Tulis isi pengumuman..." :rows="5" size="lg"
                        class="w-full" />
                </UFormField>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.judul" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>