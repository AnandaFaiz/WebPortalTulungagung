<script setup>
import { VueDraggable } from 'vue-draggable-plus';

definePageMeta({
    layout: 'admin',
    title: 'Kelola Aplikasi'
})

const toast = useToast()

const aplikasiList = ref([])
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

const statusOptions = [
    { label: 'Semua Status', value: 'all' },
    { label: 'Aktif', value: 'aktif' },
    { label: 'Nonaktif', value: 'nonaktif' }
]

const iconPresets = [
    { label: 'Apps (default)', value: 'i-tabler-apps' },
    { label: 'Users (Kependudukan)', value: 'i-tabler-users' },
    { label: 'Car (Samsat)', value: 'i-tabler-car' },
    { label: 'Speakerphone (Pengaduan)', value: 'i-tabler-speakerphone' },
    { label: 'File Check (Perizinan)', value: 'i-tabler-file-check' },
    { label: 'Heartbeat (Kesehatan)', value: 'i-tabler-heartbeat' },
    { label: 'School (Pendidikan)', value: 'i-tabler-school' },
    { label: 'Receipt (Pajak)', value: 'i-tabler-receipt' },
    { label: 'Building Store', value: 'i-tabler-building-store' },
    { label: 'Shopping Bag (UMKM)', value: 'i-tabler-shopping-bag' },
    { label: 'Mountain (Wisata)', value: 'i-tabler-mountain' },
    { label: 'Chart Line (Data)', value: 'i-tabler-chart-line' }
]

const emptyForm = () => ({
    id: null,
    nama: '',
    deskripsi: '',
    icon: 'i-tabler-apps',
    url: '',
    is_active: true
})
const form = ref(emptyForm())

const filteredAplikasi = computed(() => {
    return aplikasiList.value.filter(item => {
        const matchSearch =
            item.nama?.toLowerCase().includes(search.value.toLowerCase()) || item.deskripsi?.toLowerCase().includes(search.value.toLowerCase())
        const matchStatus =
            selectedStatus.value === 'all' || (selectedStatus.value === 'aktif' && item.is_active) || (selectedStatus.value === 'nonaktif' && !item.is_active)
        return matchSearch && matchStatus
    })
})

const paginatedAplikasi = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredAplikasi.value.slice(start, start + perPage)
})

const isDragEnabled = computed(() => {
    return !search.value && selectedStatus.value === 'all' && aplikasiList.value.length <= perPage
})

watch([search, selectedStatus], () => { page.value = 1 })

const fetchAplikasi = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/aplikasi', {
            credentials: 'include'
        })
        aplikasiList.value = response.aplikasi
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

const saveAplikasi = async () => {
    if (!form.value.nama.trim()) {
        toast.add({ title: 'Validasi gagal', description: 'Nama aplikasi wajib diisi', icon: 'i-tabler-alert-circle', color: 'error' })
        return
    }
    saving.value = true
    try {
        const payload = {
            nama: form.value.nama,
            deskripsi: form.value.deskripsi,
            icon: form.value.icon,
            url: form.value.url,
            is_active: form.value.is_active
        }
        if (editingItem.value) {
            await $fetch(`http://localhost:8080/api/admin/aplikasi/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Aplikasi berhasil diperbarui', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        } else {
            await $fetch('http://localhost:8080/api/admin/aplikasi', { method: 'POST', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Aplikasi baru berhasil ditambahkan', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        }
        showModal.value = false
        await fetchAplikasi()
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        saving.value = false
    }
}

const confirmDelete = async () => {
    if (!deletingItem.value) return
    saving.value = true
    try {
        await $fetch(`http://localhost:8080/api/admin/aplikasi/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({ title: 'Berhasil', description: 'Aplikasi berhasil dihapus', icon: 'i-tabler-circle-check', color: 'success', duration: 2500 })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchAplikasi()
    } catch (err) {
        toast.add({ title: 'Gagal menghapus', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        saving.value = false
    }
}

const onReorder = async () => {
    reordering.value = true
    try {
        const ids = aplikasiList.value.map(item => item.id)
        await $fetch('http://localhost:8080/api/admin/aplikasi/reorder', { method: 'PUT', body: { ids }, credentials: 'include' })
        toast.add({ title: 'Urutan disimpan', description: 'Urutan aplikasi berhasil diperbarui', icon: 'i-tabler-circle-check', color: 'success', duration: 2000 })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan urutan', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
        await fetchAplikasi()
    } finally {
        reordering.value = false
    }
}

const closeModal = () => {
    showModal.value = false
    editingItem.value = null
    form.value = emptyForm()
}

const extractDomain = (url) => {
    if (!url || url === '#') return null
    try { return new URL(url).hostname.replace('www.', '') } catch { return url }
}

onMounted(() => { fetchAplikasi() })
</script>

<template>
    <div class="space-y-6">

        <AdminPageHeader title="Kelola Aplikasi"
            description="Kelola aplikasi digital milik Pemkab Tulungagung yang tampil di Beranda."
            action-label="Tambah Aplikasi" @action="openCreateModal" />

        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari nama atau deskripsi..." size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-apps" class="text-sky-500" />
                        Daftar Aplikasi
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredAplikasi.length }} aplikasi</span>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data aplikasi...</p>
            </div>

            <AdminEmptyState v-else-if="filteredAplikasi.length === 0" icon="i-tabler-apps-off"
                title="Belum ada aplikasi" description="Klik tombol Tambah Aplikasi untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="w-12"></th>
                            <th
                                class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider w-16">
                                No</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aplikasi</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                URL</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                        </tr>
                    </thead>

                    <VueDraggable v-if="isDragEnabled" v-model="aplikasiList" tag="tbody" :animation="200"
                        handle=".drag-handle" ghost-class="opacity-50" @end="onReorder">
                        <tr v-for="(item, idx) in aplikasiList" :key="item.id"
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
                                        <UIcon :name="item.icon" class="w-5 h-5 text-sky-600 dark:text-sky-400" />
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                        <p class="text-xs text-slate-400 line-clamp-1 max-w-md">{{ item.deskripsi || '-'
                                        }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <span v-if="!item.url || item.url === '#'"
                                    class="text-xs text-slate-400 italic">-</span>
                                <a v-else :href="item.url" target="_blank" rel="noopener"
                                    class="text-xs text-sky-600 hover:text-sky-700 hover:underline flex items-center gap-1 max-w-50 truncate">
                                    <UIcon name="i-tabler-external-link" class="w-3.5 h-3.5 shrink-0" />
                                    {{ extractDomain(item.url) }}
                                </a>
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

                    <tbody v-else>
                        <tr v-for="(item, idx) in paginatedAplikasi" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-200" />
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-500">{{ (page - 1) * perPage + idx + 1 }}</td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div
                                        class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                        <UIcon :name="item.icon" class="w-5 h-5 text-sky-600 dark:text-sky-400" />
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                        <p class="text-xs text-slate-400 line-clamp-1 max-w-md">{{ item.deskripsi || '-'
                                        }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <span v-if="!item.url || item.url === '#'"
                                    class="text-xs text-slate-400 italic">-</span>
                                <a v-else :href="item.url" target="_blank" rel="noopener"
                                    class="text-xs text-sky-600 hover:text-sky-700 hover:underline flex items-center gap-1 max-w-50 truncate">
                                    <UIcon name="i-tabler-external-link" class="w-3.5 h-3.5 shrink-0" />
                                    {{ extractDomain(item.url) }}
                                </a>
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
                <AdminPagination v-model:page="page" :total="filteredAplikasi.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <!-- Modal Form -->
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Aplikasi' : 'Tambah Aplikasi Baru'"
            :description="editingItem ? 'Perbarui informasi aplikasi di bawah ini.' : 'Isi form berikut untuk menambahkan aplikasi baru.'"
            :is-edit="!!editingItem" :loading="saving"
            :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan Aplikasi'" @submit="saveAplikasi"
            @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Nama Aplikasi" required>
                    <UInput v-model="form.nama" placeholder="Contoh: Gerdu Kalimasada" icon="i-tabler-apps" size="lg"
                        class="w-full" />
                </UFormField>

                <UFormField label="Deskripsi Singkat">
                    <UTextarea v-model="form.deskripsi" placeholder="Contoh: Layanan Administrasi Kependudukan Online"
                        :rows="2" size="lg" class="w-full" />
                </UFormField>

                <UFormField label="URL Aplikasi" hint="Kosongkan kalau belum ada.">
                    <UInput v-model="form.url" type="url" placeholder="https://aplikasi.tulungagung.go.id"
                        icon="i-tabler-world" size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Ikon">
                    <div class="space-y-2">
                        <USelect v-model="form.icon" :items="iconPresets" value-key="value" size="lg" class="w-full" />
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
                    <USwitch v-model="form.is_active" :label="form.is_active ? 'Aktif' : 'Nonaktif'" />
                </UFormField>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama" :loading="saving"
            @confirm="confirmDelete" />

    </div>
</template>
