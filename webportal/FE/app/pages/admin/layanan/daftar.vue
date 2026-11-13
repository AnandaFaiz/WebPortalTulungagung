<script setup>
definePageMeta({ layout: 'admin', title: 'Daftar Layanan' })
const toast = useToast()
const { can } = useAuth()
const canUpdate = computed(() => can('layanan.daftar', 'update'))
const canDelete = computed(() => can('layanan.daftar', 'delete'))
const API_BASE = 'http://localhost:8080/api'

const list = ref([])
const kategoriList = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

const form = ref({ id: null, kategori_id: null, nama: '', link: '' })

const fetchData = async () => {
    loading.value = true
    try {
        const [resLayanan, resKat] = await Promise.all([
            $fetch(`${API_BASE}/admin/layanan`, { credentials: 'include' }),
            $fetch(`${API_BASE}/admin/kategori-layanan`, { credentials: 'include' })
        ])
        list.value = resLayanan.layanan || []
        kategoriList.value = resKat.kategori || []
    } catch (err) { toast.add({ title: 'Gagal memuat', color: 'error' }) }
    finally { loading.value = false }
}

const openCreateModal = () => { form.value = { id: null, kategori_id: kategoriList.value[0]?.id, nama: '', link: '' }; editingItem.value = null; showModal.value = true }
const openEditModal = (item) => { form.value = { ...item }; editingItem.value = item; showModal.value = true }
const openDeleteModal = (item) => { deletingItem.value = item; showDeleteModal.value = true }

const saveData = async () => {
    if (!form.value.nama || !form.value.kategori_id) return
    saving.value = true
    try {
        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/layanan/${editingItem.value.id}`, { method: 'PUT', body: form.value, credentials: 'include' })
        } else {
            await $fetch(`${API_BASE}/admin/layanan`, { method: 'POST', body: form.value, credentials: 'include' })
        }
        showModal.value = false; await fetchData()
        toast.add({ title: 'Berhasil', color: 'success' })
    } catch (err) { toast.add({ title: 'Gagal menyimpan', color: 'error' }) }
    finally { saving.value = false }
}

const confirmDelete = async () => {
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/layanan/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        showDeleteModal.value = false; await fetchData()
        toast.add({ title: 'Berhasil dihapus', color: 'success' })
    } catch (err) { toast.add({ title: 'Gagal menghapus', color: 'error' }) }
    finally { saving.value = false }
}
onMounted(() => fetchData())
</script>
<template>
    <div class="space-y-6">
        <AdminPageHeader title="Daftar Layanan" description="Kelola isi dropdown layanan untuk tiap kategori."
            action-label="Tambah Layanan" @action="openCreateModal" />
        <!-- Ganti bagian <UCard> tabel dengan ini -->
        <UCard :ui="{ body: 'p-0' }">
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
            </div>

            <div v-else-if="list.length === 0" class="text-center py-12">
                <UIcon name="i-tabler-list-details" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                <p class="text-slate-500">Belum ada daftar layanan.</p>
            </div>

            <div v-else class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
                            <th class="py-3 px-6 text-xs font-bold text-slate-500 uppercase">Nama Layanan</th>
                            <th class="py-3 px-6 text-xs font-bold text-slate-500 uppercase">Kategori Induk</th>
                            <th class="py-3 px-6 text-xs font-bold text-slate-500 uppercase">Link Tujuan</th>
                            <th v-if="canUpdate || canDelete"
                                class="py-3 px-6 text-xs font-bold text-slate-500 uppercase text-right">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                        <tr v-for="item in list" :key="item.id"
                            class="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-4 px-6 font-semibold text-slate-900 dark:text-white">{{ item.nama }}</td>
                            <td class="py-4 px-6">
                                <UBadge color="primary" variant="subtle" size="sm"
                                    class="flex items-center gap-1.5 w-max">
                                    <UIcon name="i-tabler-folder" class="w-3.5 h-3.5" />
                                    {{kategoriList.find(k => k.id === item.kategori_id)?.nama || 'Tanpa Kategori'}}
                                </UBadge>
                            </td>
                            <td class="py-4 px-6">
                                <a v-if="item.link" :href="item.link" target="_blank"
                                    class="text-sm text-sky-600 hover:underline flex items-center gap-1 w-max">
                                    <UIcon name="i-tabler-external-link" class="w-3.5 h-3.5" />
                                    {{ item.link }}
                                </a>
                                <span v-else class="text-sm text-slate-400 italic">-</span>
                            </td>
                            <td v-if="canUpdate || canDelete" class="py-4 px-6">
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
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Layanan' : 'Tambah Layanan'"
            :is-edit="!!editingItem" :loading="saving" @submit="saveData" @cancel="showModal = false">
            <div class="space-y-4">
                <UFormField label="Pilih Kategori" required>
                    <USelect v-model="form.kategori_id" :items="kategoriList.map(k => ({ label: k.nama, value: k.id }))"
                        value-key="value" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Nama Layanan" required>
                    <UInput v-model="form.nama" placeholder="Contoh: Akta Kelahiran" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Link Tujuan">
                    <UInput v-model="form.link" placeholder="/layanan#akta" size="lg" class="w-full" />
                </UFormField>
            </div>
        </AdminModalForm>
        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>