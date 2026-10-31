<script setup>
import { VueDraggable } from 'vue-draggable-plus'

definePageMeta({
    layout: 'admin',
    title: 'Wallpaper Hero'
})

const toast = useToast()
const API_BASE = 'http://localhost:8080/api'

// ==== State ====
const banners = ref([])
const loading = ref(true)
const saving = ref(false)
const reordering = ref(false)

const uploading = ref(false)
const fileInput = ref(null)

const search = ref('')
const selectedStatus = ref('Semua')
const selectedTipe = ref('Semua')
const page = ref(1)
const perPage = 10

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

const statusOptions = ['Semua', 'Aktif', 'Nonaktif']

const linkOptions = [
    { label: '— Tanpa Link —', value: 'none' },
    { label: 'Beranda', value: '/' },
    { label: 'Profil', value: '/profil' },
    { label: 'Pemerintahan', value: '/pemerintahan' },
    { label: 'Layanan Publik', value: '/layanan' },
    { label: 'Informasi & Berita', value: '/berita' },
    { label: 'Kontak', value: '/kontak' }
]

// Form tanpa Judul & Link
const emptyForm = () => ({
    id: null,
    gambar: '',
    status: 'Aktif',
    tipe: 'beranda'
})
const form = ref(emptyForm())

// ==== Upload Logic ====
const triggerFileInput = () => { fileInput.value?.click() }

const handleFileUpload = async (event) => {
    const file = event.target.files?.[0]
    if (!file) return

    uploading.value = true
    try {
        const formData = new FormData()
        formData.append('file', file)

        const response = await $fetch(`${API_BASE}/admin/upload`, {
            method: 'POST',
            body: formData,
            credentials: 'include'
        })

        form.value.gambar = response.url
        toast.add({ title: 'Berhasil', description: 'Gambar berhasil diunggah', icon: 'i-tabler-circle-check', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal upload', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        uploading.value = false
        if (fileInput.value) fileInput.value.value = ''
    }
}

// ==== Computed ====
const filteredBanners = computed(() => {
    return banners.value.filter(item => {
        const matchSearch = !search.value || item.judul?.toLowerCase().includes(search.value.toLowerCase())
        const matchStatus = selectedStatus.value === 'Semua' || item.status === selectedStatus.value
        const matchTipe = selectedTipe.value === 'Semua' || item.tipe === selectedTipe.value
        return matchSearch && matchStatus && matchTipe
    })
})

const paginatedBanners = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredBanners.value.slice(start, start + perPage)
})

const isDragEnabled = computed(() => {
    return !search.value && selectedStatus.value === 'Semua' && selectedTipe.value === 'Semua' && banners.value.length <= perPage
})

watch([search, selectedStatus, selectedTipe], () => { page.value = 1 })

// ==== Fetch ====
const fetchBanner = async () => {
    loading.value = true
    try {
        const response = await $fetch(`${API_BASE}/admin/banner`, { credentials: 'include' })
        banners.value = response.banner || []
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
        gambar: item.gambar || '',
        status: item.status,
        tipe: item.tipe || 'beranda'
    }
    editingItem.value = item
    showModal.value = true
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

// ==== Save ====
const saveBanner = async () => {
    if (!form.value.gambar || !form.value.tipe) {
        toast.add({ title: 'Validasi gagal', description: 'Gambar dan Halaman wajib diisi', icon: 'i-tabler-alert-circle', color: 'error' })
        return
    }
    saving.value = true
    try {
        // ✅ Sistem otomatis mengisi Judul dan Link di belakang layar
        const payload = {
            judul: `Wallpaper ${form.value.tipe}`,
            gambar: form.value.gambar,
            link: '',
            status: form.value.status,
            tipe: form.value.tipe
        }

        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/banner/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Wallpaper diperbarui', icon: 'i-tabler-circle-check', color: 'success' })
        } else {
            await $fetch(`${API_BASE}/admin/banner`, { method: 'POST', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Wallpaper ditambahkan', icon: 'i-tabler-circle-check', color: 'success' })
        }
        showModal.value = false
        await fetchBanner()
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
        await $fetch(`${API_BASE}/admin/banner/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({ title: 'Berhasil', description: 'Wallpaper dihapus', icon: 'i-tabler-circle-check', color: 'success' })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchBanner()
    } catch (err) {
        toast.add({ title: 'Gagal menghapus', description: 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        saving.value = false
    }
}

// ==== Reorder ====
const onReorder = async () => {
    reordering.value = true
    try {
        const ids = banners.value.map(item => item.id)
        await $fetch(`${API_BASE}/admin/banner/reorder`, { method: 'PUT', body: { ids }, credentials: 'include' })
        toast.add({ title: 'Urutan disimpan', description: 'Urutan wallpaper diperbarui', icon: 'i-tabler-circle-check', color: 'success', duration: 2000 })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan urutan', description: 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
        await fetchBanner()
    } finally {
        reordering.value = false
    }
}

const closeModal = () => {
    showModal.value = false
    editingItem.value = null
    form.value = emptyForm()
}

onMounted(() => { fetchBanner() })
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Wallpaper Hero" description="Kelola gambar latar (wallpaper) untuk setiap halaman."
            action-label="Tambah Wallpaper" @action="openCreateModal" />

        <!-- Filter Section -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <AdminSearch v-model="search" placeholder="Cari judul banner..." />
                <USelect v-model="selectedTipe" :items="[
                    { label: 'Semua Halaman', value: 'Semua' },
                    { label: 'Beranda', value: 'beranda' },
                    { label: 'Informasi & Berita', value: 'berita' },
                    { label: 'Pemerintahan', value: 'pemerintahan' },
                    { label: 'Layanan Publik', value: 'layanan' },
                    { label: 'Kontak', value: 'kontak' }
                ]" icon="i-tabler-layout" size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-photo" class="text-sky-500" />
                        Daftar Wallpaper
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredBanners.length }} gambar</span>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            </div>

            <AdminEmptyState v-else-if="filteredBanners.length === 0" icon="i-tabler-photo-off"
                title="Belum ada wallpaper" description="Klik Tambah Wallpaper untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="w-12"></th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Preview Gambar
                            </th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Halaman Tujuan
                            </th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase">Aksi</th>
                        </tr>
                    </thead>

                    <VueDraggable v-if="isDragEnabled" v-model="banners" tag="tbody" :animation="200"
                        handle=".drag-handle" ghost-class="opacity-50" @end="onReorder">
                        <tr v-for="item in banners" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50">
                            <td class="py-3 px-4">
                                <div
                                    class="drag-handle cursor-grab active:cursor-grabbing p-1 rounded hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors">
                                    <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-400" />
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <div
                                    class="w-32 h-16 rounded overflow-hidden bg-slate-200 shrink-0 border border-slate-200">
                                    <img v-if="item.gambar" :src="item.gambar" class="w-full h-full object-cover" />
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge color="primary" variant="subtle" size="sm" class="capitalize">
                                    {{ item.tipe }}
                                </UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <AdminStatusBadge :status="item.status" />
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
                        <tr v-for="item in paginatedBanners" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800">
                            <td class="py-3 px-4">
                                <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-200" />
                            </td>
                            <td class="py-3 px-4">
                                <div
                                    class="w-32 h-16 rounded overflow-hidden bg-slate-200 shrink-0 border border-slate-200">
                                    <img v-if="item.gambar" :src="item.gambar" class="w-full h-full object-cover" />
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge color="primary" variant="subtle" size="sm" class="capitalize">
                                    {{ item.tipe }}
                                </UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <AdminStatusBadge :status="item.status" />
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
                <AdminPagination v-model:page="page" :total="filteredBanners.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <!-- Modal Form (Sangat Sederhana) -->
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Wallpaper' : 'Tambah Wallpaper'"
            :description="editingItem ? 'Perbarui gambar wallpaper.' : 'Upload gambar wallpaper baru.'"
            :is-edit="!!editingItem" :loading="saving"
            :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan Wallpaper'" @submit="saveBanner"
            @cancel="closeModal">
            <div class="space-y-5">
                <!-- 1. Pilih Halaman -->
                <UFormField label="Tampil di Halaman" required>
                    <USelect v-model="form.tipe" :items="[
                        { label: 'Beranda', value: 'beranda' },
                        { label: 'Informasi & Berita', value: 'berita' },
                        { label: 'Pemerintahan', value: 'pemerintahan' },
                        { label: 'Layanan Publik', value: 'layanan' },
                        { label: 'Kontak', value: 'kontak' }
                    ]" value-key="value" size="lg" class="w-full" />
                </UFormField>

                <!-- 2. Upload Gambar -->
                <UFormField label="Gambar Wallpaper" required>
                    <div class="space-y-3">
                        <div v-if="form.gambar"
                            class="relative w-full aspect-video rounded-lg overflow-hidden border border-slate-200 dark:border-slate-700 bg-slate-100">
                            <img :src="form.gambar" class="w-full h-full object-cover" />
                            <div class="absolute top-2 right-2">
                                <UButton icon="i-tabler-trash" color="error" variant="solid" size="sm"
                                    @click="form.gambar = ''" />
                            </div>
                        </div>

                        <div v-else
                            class="relative border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-8 text-center hover:border-sky-500 transition-colors bg-slate-50 dark:bg-slate-900/50 cursor-pointer"
                            @click="triggerFileInput">
                            <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/webp" class="hidden"
                                @change="handleFileUpload" />
                            <UIcon :name="uploading ? 'i-tabler-loader-2' : 'i-tabler-cloud-upload'"
                                class="w-10 h-10 text-slate-400 mx-auto mb-3" :class="{ 'animate-spin': uploading }" />
                            <p class="text-sm font-medium text-slate-700 dark:text-slate-300">
                                {{ uploading ? 'Sedang mengunggah...' : 'Klik atau tarik gambar ke sini' }}
                            </p>
                            <p class="text-xs text-slate-500 mt-1">PNG, JPG, WebP (Maks. 5MB)</p>
                        </div>
                    </div>
                </UFormField>

                <!-- 3. Status -->
                <UFormField label="Status">
                    <USwitch v-model="form.status" :label="form.status === 'Aktif' ? 'Aktif' : 'Nonaktif'"
                        color="success" true-value="Aktif" false-value="Nonaktif" />
                </UFormField>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="'Wallpaper ini'" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>