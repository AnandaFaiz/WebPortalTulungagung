<script setup>
import { VueDraggable } from 'vue-draggable-plus'

definePageMeta({
    layout: 'admin',
    title: 'Kelola OPD'
})

const toast = useToast()

// ==== State ====
const opds = ref([])
const loading = ref(true)
const saving = ref(false)
const reordering = ref(false)
const search = ref('')
const selectedStatus = ref('all')
const perPage = 10
const page = ref(1)

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

// ==== Form ====
const emptyForm = () => ({
    id: null,
    nama: '',
    deskripsi: '',
    website: '',
    is_active: true
})
const form = ref(emptyForm())

// ==== Computed ====
const filteredOpds = computed(() => {
    return opds.value.filter(item => {
        const matchSearch = item.nama?.toLowerCase().includes(search.value.toLowerCase())
        const matchStatus =
            selectedStatus.value === 'all' ||
            (selectedStatus.value === 'aktif' && item.is_active) ||
            (selectedStatus.value === 'nonaktif' && !item.is_active)
        return matchSearch && matchStatus
    })
})

const paginatedOpds = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredOpds.value.slice(start, start + perPage)
})

watch([search, selectedStatus], () => {
    page.value = 1
})

// ==== Cek apakah drag-drop aktif (hanya kalau tidak ada filter/search) ====
const isDragEnabled = computed(() => {
    return !search.value && selectedStatus.value === 'all'
})

// ==== Fetch ====
const fetchOpds = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/opd',{
            credentials: 'include'
        })
        opds.value = response.opds
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
const saveOpd = async () => {
    if (!form.value.nama.trim()) {
        toast.add({
            title: 'Validasi gagal',
            description: 'Nama OPD wajib diisi',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
        return
    }

    saving.value = true
    try {
        const payload = {
            nama: form.value.nama,
            deskripsi: form.value.deskripsi,
            website: form.value.website,
            is_active: form.value.is_active
        }

        if (editingItem.value) {
            await $fetch(`http://localhost:8080/api/admin/opd/${editingItem.value.id}`, {
                method: 'PUT',
                body: payload,
                credentials: 'include'
            })
            toast.add({
                title: 'Berhasil',
                description: 'Data OPD berhasil diperbarui',
                icon: 'i-tabler-circle-check',
                color: 'success',
                duration: 2500
            })
        } else {
            await $fetch('http://localhost:8080/api/admin/opd', {
                method: 'POST',
                body: payload,
                credentials: 'include'
            })
            toast.add({
                title: 'Berhasil',
                description: 'OPD baru berhasil ditambahkan',
                icon: 'i-tabler-circle-check',
                color: 'success',
                duration: 2500
            })
        }

        showModal.value = false
        await fetchOpds()
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
        await $fetch(`http://localhost:8080/api/admin/opd/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({
            title: 'Berhasil',
            description: 'OPD berhasil dihapus',
            icon: 'i-tabler-circle-check',
            color: 'success',
            duration: 2500
        })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchOpds()
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

// ==== Reorder (drag-drop) ====
const onReorder = async () => {
    reordering.value = true
    try {
        const ids = opds.value.map(item => item.id)
        await $fetch('http://localhost:8080/api/admin/opd/reorder', {
            method: 'PUT',
            body: { ids },
            credentials: 'include'
        })
        toast.add({
            title: 'Urutan disimpan',
            description: 'Urutan OPD berhasil diperbarui',
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
        // Refresh data untuk kembalikan urutan asli
        await fetchOpds()
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

const getInitials = (nama) => {
    if (!nama) return '?'
    return nama
        .split(' ')
        .filter(w => !['dinas', 'badan', 'sekretariat', 'kantor'].includes(w.toLowerCase()))
        .slice(0, 2)
        .map(w => w.charAt(0).toUpperCase())
        .join('')
}

// ==== Init ====
onMounted(() => {
    fetchOpds()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Kelola OPD"
            description="Kelola direktori Organisasi Perangkat Daerah. Drag baris untuk mengatur urutan."
            action-label="Tambah OPD" @action="openCreateModal" />

        <!-- Filter -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari nama OPD..." size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <!-- Tabel -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-building-community" class="text-sky-500" />
                        Daftar OPD
                    </h3>
                    <div class="flex items-center gap-3">
                        <span v-if="reordering" class="text-xs text-sky-600 flex items-center gap-1">
                            <UIcon name="i-tabler-loader-2" class="w-3.5 h-3.5 animate-spin" />
                            Menyimpan urutan...
                        </span>
                        <span class="text-sm text-slate-500">{{ filteredOpds.length }} OPD</span>
                    </div>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data OPD...</p>
            </div>

            <AdminEmptyState v-else-if="filteredOpds.length === 0" icon="i-tabler-building-off"
                title="Belum ada data OPD" description="Klik tombol Tambah OPD untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="w-12"></th>
                            <th
                                class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider w-16">
                                No</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Nama OPD</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Deskripsi</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Website</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                        </tr>
                    </thead>

                    <!-- ==== Kalau drag enabled: pakai VueDraggable ==== -->
                    <VueDraggable v-if="isDragEnabled" v-model="opds" tag="tbody" :animation="200" handle=".drag-handle"
                        ghost-class="opacity-50" @end="onReorder">
                        <tr v-for="(item, idx) in paginatedOpds" :key="item.id"
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
                                        class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 dark:text-sky-400 font-bold text-xs shrink-0">
                                        {{ getInitials(item.nama) }}
                                    </div>
                                    <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <p class="text-sm text-slate-600 dark:text-slate-300 max-w-md line-clamp-2">
                                    {{ item.deskripsi || '-' }}
                                </p>
                            </td>
                            <td class="py-3 px-4">
                                <a v-if="item.website" :href="item.website" target="_blank" rel="noopener"
                                    class="text-xs text-sky-600 hover:text-sky-700 hover:underline flex items-center gap-1 max-w-50 truncate">
                                    <UIcon name="i-tabler-external-link" class="w-3.5 h-3.5 shrink-0" />
                                    {{ extractDomain(item.website) }}
                                </a>
                                <span v-else class="text-xs text-slate-400 italic">-</span>
                            </td>
                            <td class="py-3 px-4">
                                <AdminStatusBadge :status="item.is_active ? 'Aktif' : 'Nonaktif'" />
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center justify-end gap-2">
                                    <AdminButtonEdit size="xs" @click="openEditModal(item)" />
                                    <AdminButtonDelete size="xs" @click="openDeleteModal(item)" />
                                </div>
                            </td>
                        </tr>
                    </VueDraggable>

                    <!-- ==== Kalau filter aktif: tbody biasa (drag tidak aktif) ==== -->
                    <tbody v-else>
                        <tr v-for="(item, idx) in paginatedOpds" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-200" />
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-500">{{ idx + 1 }}</td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div
                                        class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 dark:text-sky-400 font-bold text-xs shrink-0">
                                        {{ getInitials(item.nama) }}
                                    </div>
                                    <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <p class="text-sm text-slate-600 dark:text-slate-300 max-w-md line-clamp-2">
                                    {{ item.deskripsi || '-' }}
                                </p>
                            </td>
                            <td class="py-3 px-4">
                                <a v-if="item.website" :href="item.website" target="_blank" rel="noopener"
                                    class="text-xs text-sky-600 hover:text-sky-700 hover:underline flex items-center gap-1 max-w-50 truncate">
                                    <UIcon name="i-tabler-external-link" class="w-3.5 h-3.5 shrink-0" />
                                    {{ extractDomain(item.website) }}
                                </a>
                                <span v-else class="text-xs text-slate-400 italic">-</span>
                            </td>
                            <td class="py-3 px-4">
                                <AdminStatusBadge :status="item.is_active ? 'Aktif' : 'Nonaktif'" />
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
                <AdminPagination v-model:page="page" :total="filteredOpds.length" :items-per-page="perPage"/>
            </template>
        </UCard>

        <!-- ==== Modal Form ==== -->
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit OPD' : 'Tambah OPD Baru'"
            :description="editingItem ? 'Perbarui informasi OPD di bawah ini.' : 'Isi form berikut untuk menambahkan OPD baru.'"
            :is-edit="!!editingItem" :loading="saving" :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan OPD'"
            @submit="saveOpd" @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Nama OPD" required>
                    <UInput v-model="form.nama" placeholder="Contoh: Dinas Kesehatan" icon="i-tabler-building-community"
                        size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Deskripsi Singkat">
                    <UTextarea v-model="form.deskripsi"
                        placeholder="Contoh: Melayani bidang kesehatan masyarakat di Kabupaten Tulungagung." :rows="3"
                        size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Website OPD" hint="Kosongkan jika OPD belum punya website.">
                    <UInput v-model="form.website" type="url" placeholder="https://dinkes.tulungagung.go.id"
                        icon="i-tabler-world" size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Status">
                    <USwitch v-model="form.is_active" :label="form.is_active ? 'Aktif' : 'Nonaktif'" color="success"/>
                </UFormField>
            </div>
        </AdminModalForm>

        <!-- ==== Modal Konfirmasi Hapus ==== -->
        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama" :loading="saving"
            @confirm="confirmDelete" />

    </div>
</template>