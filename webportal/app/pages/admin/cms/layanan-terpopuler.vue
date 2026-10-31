<script setup>
import { VueDraggable } from 'vue-draggable-plus'

definePageMeta({ layout: 'admin', title: 'Layanan Terpopuler' })
const toast = useToast()
const API_BASE = 'http://localhost:8080/api'

const list = ref([])
const loading = ref(true)
const saving = ref(false)
const reordering = ref(false)
const search = ref('')

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

// Hapus urutan dari form
const emptyForm = () => ({
    id: null, kategori: 'Warga', judul: '', deskripsi: '', icon: 'i-tabler-star', link: '', is_active: true
})
const form = ref(emptyForm())

const iconOptions = [
    { label: 'Users', value: 'i-tabler-users' },
    { label: 'Heartbeat', value: 'i-tabler-heartbeat' },
    { label: 'Speakerphone', value: 'i-tabler-speakerphone' },
    { label: 'File Text', value: 'i-tabler-file-text' },
    { label: 'Certificate', value: 'i-tabler-certificate' },
    { label: 'Shopping Bag', value: 'i-tabler-shopping-bag' },
    { label: 'Briefcase', value: 'i-tabler-briefcase' },
    { label: 'School', value: 'i-tabler-school' },
    { label: 'Phone Call', value: 'i-tabler-phone-call' }
]

const filteredList = computed(() => list.value.filter(item =>
    item.judul.toLowerCase().includes(search.value.toLowerCase())
))

// Cek apakah drag diizinkan (hanya jika tidak ada filter pencarian)
const isDragEnabled = computed(() => !search.value)

const fetchData = async () => {
    loading.value = true
    try {
        const res = await $fetch(`${API_BASE}/admin/layanan-terpopuler`, { credentials: 'include' })
        list.value = res.layanan || []
    } catch (err) {
        toast.add({ title: 'Gagal memuat', description: err?.data?.statusMessage, color: 'error' })
    } finally {
        loading.value = false
    }
}

const openCreateModal = () => { form.value = emptyForm(); editingItem.value = null; showModal.value = true }
const openEditModal = (item) => { form.value = { ...item }; editingItem.value = item; showModal.value = true }
const openDeleteModal = (item) => { deletingItem.value = item; showDeleteModal.value = true }

const saveData = async () => {
    if (!form.value.judul) return
    saving.value = true
    try {
        // Hapus urutan dari payload
        const payload = {
            kategori: form.value.kategori, judul: form.value.judul, deskripsi: form.value.deskripsi,
            icon: form.value.icon, link: form.value.link, is_active: form.value.is_active
        }
        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/layanan-terpopuler/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
        } else {
            await $fetch(`${API_BASE}/admin/layanan-terpopuler`, { method: 'POST', body: payload, credentials: 'include' })
        }
        showModal.value = false
        await fetchData()
        toast.add({ title: 'Berhasil', description: 'Data disimpan', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage, color: 'error' })
    } finally {
        saving.value = false
    }
}

const confirmDelete = async () => {
    if (!deletingItem.value) return
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/layanan-terpopuler/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        showDeleteModal.value = false
        await fetchData()
        toast.add({ title: 'Berhasil', description: 'Data dihapus', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menghapus', description: err?.data?.statusMessage, color: 'error' })
    } finally {
        saving.value = false
    }
}

// Fungsi Drag and Drop
const onReorder = async () => {
    reordering.value = true
    try {
        const ids = list.value.map(item => item.id)
        await $fetch(`${API_BASE}/admin/layanan-terpopuler/reorder`, {
            method: 'PUT',
            body: { ids },
            credentials: 'include'
        })
        toast.add({ title: 'Urutan disimpan', description: 'Urutan layanan diperbarui', icon: 'i-tabler-circle-check', color: 'success', duration: 2000 })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan urutan', description: 'Terjadi kesalahan', color: 'error' })
        await fetchData()
    } finally {
        reordering.value = false
    }
}

const closeModal = () => { showModal.value = false; editingItem.value = null; form.value = emptyForm() }
onMounted(() => fetchData())
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Layanan Terpopuler"
            description="Kelola daftar layanan di Beranda (Warga, Usaha, Pekerja)." action-label="Tambah Layanan"
            @action="openCreateModal" />

        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari layanan..." size="lg" />
            </div>
        </UCard>

        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-list" class="text-sky-500" /> Daftar Layanan
                    </h3>
                    <span v-if="reordering" class="text-xs text-sky-600 flex items-center gap-1">
                        <UIcon name="i-tabler-loader-2" class="w-3.5 h-3.5 animate-spin" /> Menyimpan...
                    </span>
                </div>
            </template>
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            </div>
            <AdminEmptyState v-else-if="filteredList.length === 0" icon="i-tabler-list-off" title="Belum ada layanan"
                description="Klik Tambah Layanan untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="w-12"></th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Layanan</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Kategori</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase">Aksi</th>
                        </tr>
                    </thead>

                    <!-- Mode Drag (Tanpa Filter) -->
                    <VueDraggable v-if="isDragEnabled" v-model="list" tag="tbody" :animation="200" handle=".drag-handle"
                        ghost-class="opacity-50" @end="onReorder">
                        <tr v-for="item in list" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50">
                            <td class="py-3 px-4">
                                <div
                                    class="drag-handle cursor-grab active:cursor-grabbing p-1 rounded hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors">
                                    <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-400" />
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div class="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center">
                                        <UIcon :name="item.icon" class="w-5 h-5 text-sky-600" />
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white">{{ item.judul }}</p>
                                        <p class="text-xs text-slate-500">{{ item.deskripsi }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge color="primary" variant="subtle" size="sm">{{ item.kategori }}</UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge :color="item.is_active ? 'success' : 'neutral'" variant="subtle" size="sm">{{
                                    item.is_active
                                        ? 'Aktif' : 'Nonaktif' }}</UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center justify-end gap-2">
                                    <AdminButtonEdit size="xs" @click="openEditModal(item)" />
                                    <AdminButtonDelete size="xs" @click="openDeleteModal(item)" />
                                </div>
                            </td>
                        </tr>
                    </VueDraggable>

                    <!-- Mode Biasa (Dengan Filter) -->
                    <tbody v-else>
                        <tr v-for="item in filteredList" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800">
                            <td class="py-3 px-4">
                                <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-200" />
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div class="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center">
                                        <UIcon :name="item.icon" class="w-5 h-5 text-sky-600" />
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white">{{ item.judul }}</p>
                                        <p class="text-xs text-slate-500">{{ item.deskripsi }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge color="primary" variant="subtle" size="sm">{{ item.kategori }}</UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge :color="item.is_active ? 'success' : 'neutral'" variant="subtle" size="sm">{{
                                    item.is_active
                                        ? 'Aktif' : 'Nonaktif' }}</UBadge>
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
        </UCard>

        <!-- MODAL FORM (Diperbarui) -->
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Layanan' : 'Tambah Layanan'"
            :is-edit="!!editingItem" :loading="saving" submit-label="Simpan" @submit="saveData" @cancel="closeModal">
            <div class="space-y-4">

                <!-- Kategori & Icon (Sejajar) -->
                <div class="grid grid-cols-2 gap-4">
                    <UFormField label="Kategori" required>
                        <USelect v-model="form.kategori" :items="['Warga', 'Usaha', 'Pencari Kerja']" size="lg"
                            class="w-full" />
                    </UFormField>
                    <UFormField label="Icon">
                        <USelect v-model="form.icon" :items="iconOptions" value-key="value" size="lg" class="w-full" />
                    </UFormField>
                </div>


                <UFormField label="Judul Layanan" required>
                    <UInput v-model="form.judul" placeholder="Contoh: Layanan Kependudukan" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Deskripsi Singkat">
                    <UInput v-model="form.deskripsi" placeholder="KTP, KK, Akta" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Link Tujuan">
                    <UInput v-model="form.link" placeholder="/layanan#kependudukan" icon="i-tabler-link" size="lg"
                        class="w-full" />
                </UFormField>
                <UFormField label="Status">
                    <USwitch v-model="form.is_active" :label="form.is_active ? 'Aktif' : 'Nonaktif'" color="success" />
                </UFormField>
                <!-- Kotak Preview -->
                <div
                    class="p-4 rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50">
                    <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-3">Preview Tampilan</p>
                    <div class="flex gap-4">
                        <UIcon :name="form.icon" class="w-8 h-8 text-slate-400 shrink-0" />
                        <div>
                            <h4 class="font-bold text-slate-800 dark:text-slate-200">{{ form.judul || 'Judul Layanan' }}
                            </h4>
                            <p class="text-sm text-slate-500 mt-1">{{ form.deskripsi || 'Deskripsi singkat layanan...'
                            }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.judul" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>