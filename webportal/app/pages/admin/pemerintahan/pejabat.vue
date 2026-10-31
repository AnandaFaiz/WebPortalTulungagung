<script setup>
import { VueDraggable } from 'vue-draggable-plus'

definePageMeta({
    layout: 'admin',
    title: 'Kelola Pejabat'
})

const toast = useToast()

// ==== State ====
const pejabatList = ref([])
const loading = ref(true)
const saving = ref(false)
const reordering = ref(false)
const search = ref('')
const selectedKategori = ref('all')
const perPage = 10
const page = ref(1)

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

// ==== Opsi Kategori ====
const kategoriOptions = [
    { label: 'Semua Kategori', value: 'all' },
    { label: 'Kepala Daerah', value: 'Kepala Daerah' },
    { label: 'Sekretaris Daerah', value: 'Sekda' },
    { label: 'Asisten', value: 'Asisten' },
    { label: 'Staf Ahli', value: 'Staf Ahli' },
    { label: 'Kepala OPD', value: 'Kepala OPD' },
    { label: 'Camat', value: 'Camat' }
]

const kategoriFormOptions = [
    'Kepala Daerah',
    'Sekda',
    'Asisten',
    'Staf Ahli',
    'Kepala OPD',
    'Camat'
]

// ==== Form ====
const emptyForm = () => ({
    id: null,
    nama: '',
    jabatan: '',
    kategori: 'Kepala OPD',
    foto: '',
    periode: '2024-2029',
    is_active: true
})
const form = ref(emptyForm())

// ==== Computed ====
const filteredPejabat = computed(() => {
    return pejabatList.value.filter(item => {
        const matchSearch =
            item.nama?.toLowerCase().includes(search.value.toLowerCase()) ||
            item.jabatan?.toLowerCase().includes(search.value.toLowerCase())
        const matchKategori =
            selectedKategori.value === 'all' || item.kategori === selectedKategori.value
        return matchSearch && matchKategori
    })
})

// ==== Paginated data ====
const paginatedPejabat = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredPejabat.value.slice(start, start + perPage)
})

// ==== Reset ke halaman 1 saat filter berubah ====
watch([search, selectedKategori], () => {
    page.value = 1
})

const isDragEnabled = computed(() => {
    return !search.value && selectedKategori.value === 'all'
})

// ==== Fetch ====
const fetchPejabat = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/pejabat',{
            credentials: 'include'
        })
        pejabatList.value = response.pejabat
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
const savePejabat = async () => {
    if (!form.value.nama.trim() || !form.value.jabatan.trim()) {
        toast.add({
            title: 'Validasi gagal',
            description: 'Nama dan jabatan wajib diisi',
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
            jabatan: form.value.jabatan,
            kategori: form.value.kategori,
            foto: form.value.foto,
            periode: form.value.periode,
            is_active: form.value.is_active
        }

        if (editingItem.value) {
            await $fetch(`http://localhost:8080/api/admin/pejabat/${editingItem.value.id}`, {
                method: 'PUT',
                body: payload,
                credentials: 'include'
            })
            toast.add({
                title: 'Berhasil',
                description: 'Data pejabat berhasil diperbarui',
                icon: 'i-tabler-circle-check',
                color: 'success',
                duration: 2500
            })
        } else {
            await $fetch('http://localhost:8080/api/admin/pejabat', {
                method: 'POST',
                body: payload,
                credentials: 'include'
            })
            toast.add({
                title: 'Berhasil',
                description: 'Pejabat baru berhasil ditambahkan',
                icon: 'i-tabler-circle-check',
                color: 'success',
                duration: 2500
            })
        }

        showModal.value = false
        await fetchPejabat()
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
        await $fetch(`http://localhost:8080/api/admin/pejabat/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({
            title: 'Berhasil',
            description: 'Pejabat berhasil dihapus',
            icon: 'i-tabler-circle-check',
            color: 'success',
            duration: 2500
        })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchPejabat()
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
        const ids = pejabatList.value.map(item => item.id)
        await $fetch('http://localhost:8080/api/admin/pejabat/reorder', {
            method: 'PUT',
            body: { ids },
            credentials: 'include'
        })
        toast.add({
            title: 'Urutan disimpan',
            description: 'Urutan pejabat berhasil diperbarui',
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
        await fetchPejabat()
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

const getInitials = (nama) => {
    if (!nama) return '?'
    return nama
        .split(' ')
        .filter(w => !['h', 'hj', 'drs', 'dr', 'ir'].includes(w.toLowerCase()))
        .slice(0, 2)
        .map(w => w.charAt(0).toUpperCase())
        .join('')
}

const getKategoriLabel = (kategori) => {
    return kategori || '-'
}

const getKategoriColor = (kategori) => {
    const map = {
        'Kepala Daerah': 'primary',
        'Sekda': 'warning',
        'Asisten': 'info',
        'Staf Ahli': 'success',
        'Kepala OPD': 'neutral',
        'Camat': 'neutral'
    }
    return map[kategori] || 'neutral'
}

// ==== Init ====
onMounted(() => {
    fetchPejabat()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Kelola Pejabat"
            description="Kelola data pejabat daerah: Bupati, Wakil, Sekda, Asisten, dan Kepala OPD."
            action-label="Tambah Pejabat" @action="openCreateModal" />

        <!-- Filter -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari nama atau jabatan..." size="lg" />
                <USelect v-model="selectedKategori" :items="kategoriOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <!-- Tabel -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-user-star" class="text-sky-500" />
                        Daftar Pejabat
                    </h3>
                    <div class="flex items-center gap-3">
                        <span v-if="reordering" class="text-xs text-sky-600 flex items-center gap-1">
                            <UIcon name="i-tabler-loader-2" class="w-3.5 h-3.5 animate-spin" />
                            Menyimpan urutan...
                        </span>
                        <span class="text-sm text-slate-500">{{ filteredPejabat.length }} pejabat</span>
                    </div>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data pejabat...</p>
            </div>

            <AdminEmptyState v-else-if="filteredPejabat.length === 0" icon="i-tabler-user-off"
                title="Belum ada data pejabat" description="Klik tombol Tambah Pejabat untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="w-12"></th>
                            <th
                                class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider w-16">
                                No</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Pejabat</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Kategori</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Periode</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                        </tr>
                    </thead>

                    <VueDraggable v-if="isDragEnabled" v-model="pejabatList" tag="tbody" :animation="200"
                        handle=".drag-handle" ghost-class="opacity-50" @end="onReorder">
                        <tr v-for="(item, idx) in paginatedPejabat" :key="item.id"
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
                                        class="w-10 h-10 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center overflow-hidden shrink-0">
                                        <img v-if="item.foto" :src="item.foto" :alt="item.nama"
                                            class="w-full h-full object-cover" />
                                        <span v-else class="text-sky-600 dark:text-sky-400 font-bold text-xs">
                                            {{ getInitials(item.nama) }}
                                        </span>
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                        <p class="text-xs text-slate-400">{{ item.jabatan }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge :color="getKategoriColor(item.kategori)" variant="subtle" size="sm">
                                    {{ getKategoriLabel(item.kategori) }}
                                </UBadge>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600 dark:text-slate-300">
                                {{ item.periode || '-' }}
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
                        <tr v-for="(item, idx) in paginatedPejabat" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-200" />
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-500">{{ idx + 1 }}</td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div
                                        class="w-10 h-10 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center overflow-hidden shrink-0">
                                        <img v-if="item.foto" :src="item.foto" :alt="item.nama"
                                            class="w-full h-full object-cover" />
                                        <span v-else class="text-sky-600 dark:text-sky-400 font-bold text-xs">
                                            {{ getInitials(item.nama) }}
                                        </span>
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                        <p class="text-xs text-slate-400">{{ item.jabatan }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge :color="getKategoriColor(item.kategori)" variant="subtle" size="sm">
                                    {{ getKategoriLabel(item.kategori) }}
                                </UBadge>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600 dark:text-slate-300">
                                {{ item.periode || '-' }}
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
                <AdminPagination v-model:page="page" :total="filteredPejabat.length" :items-per-page="perPage"/>
            </template>
        </UCard>

        <!-- ==== Modal Form ==== -->
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Pejabat' : 'Tambah Pejabat Baru'"
            :description="editingItem ? 'Perbarui informasi pejabat di bawah ini.' : 'Isi form berikut untuk menambahkan pejabat baru.'"
            :is-edit="!!editingItem" :loading="saving"
            :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan Pejabat'" @submit="savePejabat"
            @cancel="closeModal">
            <div class="space-y-4">
                <!-- Preview Foto -->
                <div class="flex items-center gap-4">
                    <div
                        class="w-20 h-20 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center overflow-hidden shrink-0 ring-2 ring-white dark:ring-slate-800 shadow-md">
                        <img v-if="form.foto" :src="form.foto" class="w-full h-full object-cover" />
                        <span v-else class="text-sky-600 dark:text-sky-400 font-bold text-lg">
                            {{ form.nama ? getInitials(form.nama) : '?' }}
                        </span>
                    </div>
                    <div class="flex-1">
                        <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1.5">Foto
                            Pejabat</label>
                        <div
                            class="border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-3 text-center hover:border-sky-500 transition-colors cursor-pointer">
                            <p class="text-xs text-slate-500">Klik untuk upload foto</p>
                            <p class="text-[10px] text-slate-400 mt-0.5">PNG, JPG (Max 2MB) - Rasio 1:1</p>
                        </div>
                    </div>
                </div>

                <UFormField label="Nama Lengkap" required>
                    <UInput v-model="form.nama" placeholder="Contoh: H. Gatut Sunu Wibowo, S.E., M.E."
                        icon="i-tabler-user" size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Jabatan" required>
                    <UInput v-model="form.jabatan" placeholder="Contoh: Bupati Tulungagung" icon="i-tabler-briefcase"
                        size="lg" class="w-full" />
                </UFormField>

                <div class="grid grid-cols-2 gap-4">
                    <UFormField label="Kategori" required>
                        <USelect v-model="form.kategori" :items="kategoriFormOptions" size="lg" class="w-full" />
                    </UFormField>
                    <UFormField label="Periode">
                        <UInput v-model="form.periode" placeholder="2024-2029" icon="i-tabler-calendar" size="lg"
                            class="w-full" />
                    </UFormField>
                </div>

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