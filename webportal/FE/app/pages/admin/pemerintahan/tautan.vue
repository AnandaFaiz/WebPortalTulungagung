<script setup>
import { VueDraggable } from 'vue-draggable-plus'

definePageMeta({
    layout: 'admin',
    title: 'Tautan Eksternal'
})

const toast = useToast()
const { can } = useAuth()
const canUpdate = computed(() => can('pemerintahan.tautan', 'update'))
const canDelete = computed(() => can('pemerintahan.tautan', 'delete'))

// ==== State ====
const tautanList = ref([])
const loading = ref(true)
const saving = ref(false)
const reordering = ref(false)
const search = ref('')
const selectedStatus = ref('all')
const page = ref(1)
const perPage = 10

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

// ==== Opsi ====
const statusOptions = [
    { label: 'Semua Status', value: 'all' },
    { label: 'Aktif', value: 'aktif' },
    { label: 'Nonaktif', value: 'nonaktif' }
]

// ==== Preset Ikon (untuk memudahkan editor) ====
const iconPresets = [
    { label: 'External Link (default)', value: 'i-tabler-external-link' },
    { label: 'File Search (PPID)', value: 'i-tabler-file-search' },
    { label: 'Scale (JDIH)', value: 'i-tabler-scale' },
    { label: 'Gavel (DPRD)', value: 'i-tabler-gavel' },
    { label: 'Building Bank', value: 'i-tabler-building-bank' },
    { label: 'World / Website', value: 'i-tabler-world' },
    { label: 'News', value: 'i-tabler-news' },
    { label: 'File Text', value: 'i-tabler-file-text' },
    { label: 'Shield (Keamanan)', value: 'i-tabler-shield' },
    { label: 'Users (Masyarakat)', value: 'i-tabler-users' }
]

// ==== Form ====
const emptyForm = () => ({
    id: null,
    nama: '',
    url: '',
    deskripsi: '',
    icon: 'i-tabler-external-link',
    is_active: true
})
const form = ref(emptyForm())

// ==== Computed ====
const filteredTautan = computed(() => {
    const list = tautanList.value || []
    return list.filter(item => {
        const matchSearch =
            item.nama?.toLowerCase().includes(search.value.toLowerCase()) ||
            item.deskripsi?.toLowerCase().includes(search.value.toLowerCase())
        const matchStatus =
            selectedStatus.value === 'all' ||
            (selectedStatus.value === 'aktif' && item.is_active) ||
            (selectedStatus.value === 'nonaktif' && !item.is_active)
        return matchSearch && matchStatus
    })
})

const paginatedTautan = computed(() => {
    const list = filteredTautan.value || []
    const start = (page.value - 1) * perPage
    return list.slice(start, start + perPage)
})

const isDragEnabled = computed(() => {
    return !search.value && selectedStatus.value === 'all' && (tautanList.value?.length || 0) <= perPage
})

watch([search, selectedStatus], () => {
    page.value = 1
})

// ==== Fetch ====
const fetchTautan = async () => {
    loading.value = true
    try {
        const requestFetch = useRequestFetch()
        const response = await requestFetch('http://localhost:8080/api/admin/tautan-eksternal', { credentials: 'include' })
        tautanList.value = Array.isArray(response.tautan) ? response.tautan : []
    } catch (err) {
        toast.add({
            title: 'Gagal memuat data',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
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
    form.value = { ...item }
    editingItem.value = item
    showModal.value = true
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

// ==== Save ====
const saveTautan = async () => {
    if (!form.value.nama.trim() || !form.value.url.trim()) {
        toast.add({
            title: 'Validasi gagal',
            description: 'Nama dan URL wajib diisi',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
        return
    }

    saving.value = true
    try {
        const requestFetch = useRequestFetch()
        const payload = {
            nama: form.value.nama,
            url: form.value.url,
            deskripsi: form.value.deskripsi,
            icon: form.value.icon,
            is_active: form.value.is_active
        }

        if (editingItem.value) {
            await requestFetch(`http://localhost:8080/api/admin/tautan-eksternal/${editingItem.value.id}`, {
                method: 'PUT',
                body: payload,
                credentials: 'include'
            })
            toast.add({
                title: 'Berhasil',
                description: 'Tautan berhasil diperbarui',
                icon: 'i-tabler-circle-check',
                color: 'success',
                duration: 2500
            })
        } else {
            await requestFetch('http://localhost:8080/api/admin/tautan-eksternal', {
                method: 'POST',
                body: payload,
                credentials: 'include'
            })
            toast.add({
                title: 'Berhasil',
                description: 'Tautan baru berhasil ditambahkan',
                icon: 'i-tabler-circle-check',
                color: 'success',
                duration: 2500
            })
        }

        showModal.value = false
        await fetchTautan()
    } catch (err) {
        toast.add({
            title: 'Gagal menyimpan',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
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
        await requestFetch(`http://localhost:8080/api/admin/tautan-eksternal/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({
            title: 'Berhasil',
            description: 'Tautan berhasil dihapus',
            icon: 'i-tabler-circle-check',
            color: 'success',
            duration: 2500
        })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchTautan()
    } catch (err) {
        toast.add({
            title: 'Gagal menghapus',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
    } finally {
        saving.value = false
    }
}

// ==== Reorder ====
const onReorder = async () => {
    reordering.value = true
    try {
        const requestFetch = useRequestFetch()
        const ids = tautanList.value.map(item => item.id)
        await requestFetch('http://localhost:8080/api/admin/tautan-eksternal/reorder', {
            method: 'PUT',
            body: { ids },
            credentials: 'include'
        })
        toast.add({
            title: 'Urutan disimpan',
            description: 'Urutan tautan berhasil diperbarui',
            icon: 'i-tabler-circle-check',
            color: 'success',
            duration: 2000
        })
    } catch (err) {
        toast.add({
            title: 'Gagal menyimpan urutan',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
        await fetchTautan()
    } finally {
        reordering.value = false
    }
}

// ==== Helpers ====
const closeModal = () => {
    showModal.value = false
    editingItem.value = null
    form.value = emptyForm()
}

const extractDomain = (url) => {
    if (!url) return null
    try {
        return new URL(url).hostname.replace('www.', '')
    } catch {
        return url
    }
}

// ==== Init ====
onMounted(() => {
    fetchTautan()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Tautan Eksternal"
            description="Kelola tautan ke portal/website pemerintah lainnya (PPID, JDIH, DPRD, dll)."
            action-label="Tambah Tautan" @action="openCreateModal" />

        <!-- Filter -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari nama atau deskripsi..." size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <!-- Tabel -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-external-link" class="text-sky-500" />
                        Daftar Tautan
                    </h3>
                    <div class="flex items-center gap-3">
                        <span v-if="reordering" class="text-xs text-sky-600 flex items-center gap-1">
                            <UIcon name="i-tabler-loader-2" class="w-3.5 h-3.5 animate-spin" />
                            Menyimpan urutan...
                        </span>
                        <span class="text-sm text-slate-500">{{ filteredTautan?.length || 0 }} tautan</span>
                    </div>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data tautan...</p>
            </div>

            <AdminEmptyState v-else-if="!filteredTautan?.length" icon="i-tabler-link-off"
                title="Belum ada tautan eksternal" description="Klik tombol Tambah Tautan untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="w-12"></th>
                            <th
                                class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider w-16">
                                No</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Nama</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                URL</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Status</th>
                            <th v-if="canUpdate || canDelete"
                                class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                        </tr>
                    </thead>

                    <VueDraggable v-if="isDragEnabled && tautanList && tautanList.length > 0" v-model="tautanList"
                        tag="tbody" :animation="200" handle=".drag-handle" ghost-class="opacity-50" @end="onReorder">
                        <tr v-for="(item, idx) in tautanList" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <div
                                    class="drag-handle cursor-grab active:cursor-grabbing p-1 rounded hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors">
                                    <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-400" />
                                </div>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-500">{{ idx + 1 }}</td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div
                                        class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                        <UIcon :name="item.icon || 'i-tabler-external-link'"
                                            class="w-5 h-5 text-sky-600 dark:text-sky-400" />
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                        <p class="text-xs text-slate-400 line-clamp-1 max-w-md">{{ item.deskripsi || '-'
                                        }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <a :href="item.url" target="_blank" rel="noopener"
                                    class="text-xs text-sky-600 hover:text-sky-700 hover:underline flex items-center gap-1 max-w-50 truncate">
                                    <UIcon name="i-tabler-external-link" class="w-3.5 h-3.5 shrink-0" />
                                    {{ extractDomain(item.url) }}
                                </a>
                            </td>
                            <td class="py-3 px-4">
                                <AdminStatusBadge :status="item.is_active ? 'Aktif' : 'Nonaktif'" />
                            </td>
                            <td v-if="canUpdate || canDelete" class="py-3 px-4">
                                <div class="flex items-center justify-end gap-2">
                                    <AdminButtonEdit size="xs" @click="openEditModal(item)" />
                                    <AdminButtonDelete size="xs" @click="openDeleteModal(item)" />
                                </div>
                            </td>
                        </tr>
                    </VueDraggable>

                    <tbody v-else>
                        <tr v-for="(item, idx) in paginatedTautan" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-200" />
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-500">{{ (page - 1) * perPage + idx + 1 }}</td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div
                                        class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                        <UIcon :name="item.icon || 'i-tabler-external-link'"
                                            class="w-5 h-5 text-sky-600 dark:text-sky-400" />
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                        <p class="text-xs text-slate-400 line-clamp-1 max-w-md">{{ item.deskripsi || '-'
                                        }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <a :href="item.url" target="_blank" rel="noopener"
                                    class="text-xs text-sky-600 hover:text-sky-700 hover:underline flex items-center gap-1 max-w-50 truncate">
                                    <UIcon name="i-tabler-external-link" class="w-3.5 h-3.5 shrink-0" />
                                    {{ extractDomain(item.url) }}
                                </a>
                            </td>
                            <td class="py-3 px-4">
                                <AdminStatusBadge :status="item.is_active ? 'Aktif' : 'Nonaktif'" />
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
                <AdminPagination v-model:page="page" :total="filteredTautan?.length || 0" :items-per-page="perPage" />
            </template>
        </UCard>

        <!-- ==== Modal Form ==== -->
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Tautan' : 'Tambah Tautan Baru'"
            :description="editingItem ? 'Perbarui informasi tautan di bawah ini.' : 'Isi form berikut untuk menambahkan tautan baru.'"
            :is-edit="!!editingItem" :loading="saving"
            :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan Tautan'" @submit="saveTautan"
            @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Nama Tautan" required>
                    <UInput v-model="form.nama" placeholder="Contoh: PPID, JDIH, DPRD" icon="i-tabler-link" size="lg"
                        class="w-full" />
                </UFormField>

                <UFormField label="URL" required>
                    <UInput v-model="form.url" type="url" placeholder="https://ppid.tulungagung.go.id"
                        icon="i-tabler-world" size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Deskripsi Singkat">
                    <UTextarea v-model="form.deskripsi" placeholder="Deskripsi singkat tentang tautan ini..." :rows="2"
                        size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Ikon">
                    <div class="space-y-2">
                        <USelect v-model="form.icon" :items="iconPresets" value-key="value" size="lg" class="w-full" />
                        <!-- Preview ikon -->
                        <div class="flex items-center gap-2 text-xs text-slate-500">
                            <span>Preview:</span>
                            <div class="w-8 h-8 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center">
                                <UIcon :name="form.icon" class="w-4 h-4 text-sky-600 dark:text-sky-400" />
                            </div>
                            <code
                                class="text-[10px] bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded">{{ form.icon }}</code>
                        </div>
                    </div>
                </UFormField>

                <UFormField label="Status">
                    <USwitch v-model="form.is_active" :label="form.is_active ? 'Aktif' : 'Nonaktif'" color="success" />
                </UFormField>
            </div>
        </AdminModalForm>

        <!-- ==== Modal Konfirmasi Hapus ==== -->
        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama" :loading="saving"
            @confirm="confirmDelete" />

    </div>
</template>