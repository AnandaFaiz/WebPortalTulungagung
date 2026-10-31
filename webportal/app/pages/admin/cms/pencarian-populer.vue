<script setup>
definePageMeta({ layout: 'admin', title: 'Pencarian Populer' })
const toast = useToast()
const API_BASE = 'http://localhost:8080/api'

const list = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)
const form = ref({ id: null, keyword: '' })

const fetchData = async () => {
    loading.value = true
    try {
        const res = await $fetch(`${API_BASE}/admin/pencarian-populer`, { credentials: 'include' })
        list.value = res.pencarian || []
    } catch (err) {
        toast.add({ title: 'Gagal memuat', description: err?.data?.statusMessage, color: 'error' })
    } finally {
        loading.value = false
    }
}

const openCreateModal = () => { form.value = { id: null, keyword: '' }; editingItem.value = null; showModal.value = true }
const openEditModal = (item) => { form.value = { ...item }; editingItem.value = item; showModal.value = true }
const openDeleteModal = (item) => { deletingItem.value = item; showDeleteModal.value = true }

const saveData = async () => {
    if (!form.value.keyword) return
    saving.value = true
    try {
        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/pencarian-populer/${editingItem.value.id}`, { method: 'PUT', body: form.value, credentials: 'include' })
        } else {
            await $fetch(`${API_BASE}/admin/pencarian-populer`, { method: 'POST', body: form.value, credentials: 'include' })
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
        await $fetch(`${API_BASE}/admin/pencarian-populer/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
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
        <AdminPageHeader title="Pencarian Populer"
            description="Kelola kata kunci yang muncul di bawah kolom pencarian Beranda." action-label="Tambah Keyword"
            @action="openCreateModal" />

        <UCard>
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            </div>
            <div v-else-if="list.length === 0" class="text-center py-12 text-slate-500">Belum ada data.</div>
            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase w-16">No</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Keyword</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase">Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(item, idx) in list" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800">
                            <td class="py-3 px-4 text-sm text-slate-500">{{ idx + 1 }}</td>
                            <td class="py-3 px-4 font-medium text-slate-900 dark:text-white">{{ item.keyword }}</td>
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

        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Keyword' : 'Tambah Keyword'"
            :is-edit="!!editingItem" :loading="saving" submit-label="Simpan" @submit="saveData"
            @cancel="showModal = false">
            <UFormField label="Kata Kunci" required>
                <UInput v-model="form.keyword" placeholder="Contoh: KTP & KK" size="lg" class="w-full" />
            </UFormField>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.keyword" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>