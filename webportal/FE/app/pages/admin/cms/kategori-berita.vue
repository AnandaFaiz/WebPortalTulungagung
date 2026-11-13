<script setup>
definePageMeta({ layout: 'admin', title: 'Kategori Berita' })
const toast = useToast()
const { can } = useAuth()
const canUpdate = computed(() => can('cms.kategori', 'update'))
const canDelete = computed(() => can('cms.kategori', 'delete'))
const API_BASE = 'http://localhost:8080/api'

const list = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)
const form = ref({ id: null, nama: '', deskripsi: '' })

const fetchData = async () => {
    loading.value = true
    try {
        const res = await $fetch(`${API_BASE}/admin/kategori-berita`, { credentials: 'include' })
        list.value = res.kategori || []
    } catch (err) {
        toast.add({ title: 'Gagal memuat', description: err?.data?.statusMessage, color: 'error' })
    } finally {
        loading.value = false
    }
}

const openCreateModal = () => { form.value = { id: null, nama: '', deskripsi: '' }; editingItem.value = null; showModal.value = true }
const openEditModal = (item) => { form.value = { ...item }; editingItem.value = item; showModal.value = true }
const openDeleteModal = (item) => { deletingItem.value = item; showDeleteModal.value = true }

const saveData = async () => {
    if (!form.value.nama) return
    saving.value = true
    try {
        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/kategori-berita/${editingItem.value.id}`, { method: 'PUT', body: form.value, credentials: 'include' })
        } else {
            await $fetch(`${API_BASE}/admin/kategori-berita`, { method: 'POST', body: form.value, credentials: 'include' })
        }
        showModal.value = false
        await fetchData()
        toast.add({ title: 'Berhasil', description: 'Data disimpan', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage, color: 'error' })
    } finally { saving.value = false }
}

const confirmDelete = async () => {
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/kategori-berita/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        showDeleteModal.value = false
        await fetchData()
        toast.add({ title: 'Berhasil', description: 'Data dihapus', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menghapus', description: err?.data?.statusMessage, color: 'error' })
    } finally { saving.value = false }
}

onMounted(() => fetchData())
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Kategori Berita" description="Kelola kategori yang muncul di Mega Menu Header."
            action-label="Tambah Kategori" @action="openCreateModal" />

        <UCard>
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            </div>
            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Nama Kategori
                            </th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Deskripsi</th>
                            <th v-if="canUpdate || canDelete"
                                class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase">Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in list" :key="item.id" class="border-b border-slate-100 dark:border-slate-800">
                            <td class="py-3 px-4 font-medium text-slate-900 dark:text-white">{{ item.nama }}</td>
                            <td class="py-3 px-4 text-sm text-slate-500 max-w-xs truncate" :title="item.deskripsi">
                                {{ item.deskripsi || '-' }}
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
        </UCard>

        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Kategori' : 'Tambah Kategori'"
            :is-edit="!!editingItem" :loading="saving" submit-label="Simpan" @submit="saveData"
            @cancel="showModal = false">
            <div class="space-y-4">
                <UFormField label="Nama Kategori" required>
                    <UInput v-model="form.nama" placeholder="Contoh: Pemerintahan" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Deskripsi Singkat">
                    <UInput v-model="form.deskripsi" placeholder="Kabar seputar pemda" size="lg" class="w-full" />
                </UFormField>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>