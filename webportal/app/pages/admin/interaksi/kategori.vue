<script setup>
definePageMeta({ layout: 'admin', title: 'Kategori & Status Aduan' })
const toast = useToast()
const API = 'http://localhost:8080/api/admin'
const tabs = [{ label: 'Kategori Aduan', key: 'kategori', icon: 'i-tabler-category' }, { label: 'Status Aduan', key: 'status', icon: 'i-tabler-progress-check' }]
const activeTab = ref('kategori')
const data = ref({ kategori: [], status: [] })
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)
const form = ref({ id: null, nama: '', deskripsi: '', is_active: true })
const currentItems = computed(() => data.value[activeTab.value])
const endpoint = key => `${API}/pengaduan-${key}`
const fetchData = async () => {
    loading.value = true
    try {
        const [kategori, status] = await Promise.all([['kategori', 'status'].map(key => $fetch(endpoint(key), { credentials: 'include' }))].flat())
        data.value = { kategori: kategori.items || [], status: status.items || [] }
    } catch (err) { toast.add({ title: 'Gagal memuat konfigurasi', color: 'error' }) }
    finally { loading.value = false }
}
const openCreate = () => { form.value = { id: null, nama: '', deskripsi: '', is_active: true }; editingItem.value = null; showModal.value = true }
const openEdit = item => { form.value = { ...item }; editingItem.value = item; showModal.value = true }
const openDelete = item => { deletingItem.value = item; showDeleteModal.value = true }
const saveData = async () => {
    if (!form.value.nama.trim()) { toast.add({ title: 'Nama wajib diisi', color: 'error' }); return }
    saving.value = true
    try { const url = editingItem.value ? `${endpoint(activeTab.value)}/${editingItem.value.id}` : endpoint(activeTab.value); await $fetch(url, { method: editingItem.value ? 'PUT' : 'POST', body: form.value, credentials: 'include' }); showModal.value = false; await fetchData(); toast.add({ title: 'Konfigurasi disimpan', color: 'success' }) }
    catch (err) { toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage || 'Terjadi kesalahan', color: 'error' }) }
    finally { saving.value = false }
}
const deleteData = async () => {
    saving.value = true
    try { await $fetch(`${endpoint(activeTab.value)}/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' }); showDeleteModal.value = false; await fetchData(); toast.add({ title: 'Konfigurasi dihapus', color: 'success' }) }
    catch (err) { toast.add({ title: 'Gagal menghapus', color: 'error' }) }
    finally { saving.value = false }
}
onMounted(fetchData)
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Kategori & Status Aduan"
            description="Kelola pilihan kategori dan alur status yang digunakan warga serta admin."
            action-label="Tambah Data" @action="openCreate" />
        <div class="flex gap-2 border-b border-slate-200 dark:border-slate-700"><button v-for="tab in tabs"
                :key="tab.key" class="flex items-center gap-2 px-4 py-3 text-sm font-semibold border-b-2"
                :class="activeTab === tab.key ? 'text-sky-600 border-sky-500' : 'text-slate-500 border-transparent'"
                @click="activeTab = tab.key">
                <UIcon :name="tab.icon" class="w-5 h-5" />{{ tab.label }}
            </button></div>
        <UCard :ui="{ body: 'p-0' }">
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
            </div>
            <AdminEmptyState v-else-if="!currentItems.length" icon="i-tabler-inbox" title="Belum ada data"
                description="Tambahkan konfigurasi baru untuk digunakan di frontend." />
            <div v-else class="divide-y divide-slate-100 dark:divide-slate-800">
                <div v-for="item in currentItems" :key="item.id" class="flex items-center gap-4 px-6 py-4">
                    <div class="w-10 h-10 rounded-lg bg-sky-50 dark:bg-sky-950 flex items-center justify-center">
                        <AutoIcon :text="item.nama" size="w-5 h-5 text-sky-600" />
                    </div>
                    <div class="flex-1">
                        <p class="font-semibold text-slate-900 dark:text-white">{{ item.nama }}</p>
                        <p class="text-sm text-slate-500">{{ item.deskripsi || 'Tanpa deskripsi' }}</p>
                    </div>
                    <UBadge :color="item.is_active ? 'success' : 'neutral'" variant="subtle">{{ item.is_active ? 'Aktif'
                        : 'Nonaktif' }}</UBadge>
                    <div class="flex gap-2">
                        <AdminButtonEdit size="xs" @click="openEdit(item)" />
                        <AdminButtonDelete size="xs" @click="openDelete(item)" />
                    </div>
                </div>
            </div>
        </UCard>
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Konfigurasi' : 'Tambah Konfigurasi'"
            :is-edit="!!editingItem" :loading="saving" @submit="saveData" @cancel="showModal = false">
            <div class="space-y-4">
                <UFormField label="Nama" required>
                    <UInput v-model="form.nama" placeholder="Contoh: Layanan Publik" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Deskripsi">
                    <UTextarea v-model="form.deskripsi" placeholder="Penjelasan singkat" :rows="3" class="w-full" />
                </UFormField>
                <UCheckbox v-model="form.is_active" label="Aktif dan tampil di frontend" />
            </div>
        </AdminModalForm>
        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama" :loading="saving"
            @confirm="deleteData" />
    </div>
</template>
