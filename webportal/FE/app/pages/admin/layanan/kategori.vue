<script setup>
definePageMeta({ layout: 'admin', title: 'Kategori Layanan' })
const toast = useToast()
const { can } = useAuth()
const canUpdate = computed(() => can('layanan.kategori', 'update'))
const canDelete = computed(() => can('layanan.kategori', 'delete'))
const API_BASE = 'http://localhost:8080/api'

const list = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

// Form tanpa icon manual
const form = ref({ id: null, nama: '', deskripsi: '' })

const fetchData = async () => {
    loading.value = true
    try {
        const res = await $fetch(`${API_BASE}/admin/kategori-layanan`, { credentials: 'include' })
        list.value = res.kategori || []
    } catch (err) { toast.add({ title: 'Gagal memuat', color: 'error' }) }
    finally { loading.value = false }
}

const openCreateModal = () => { form.value = { id: null, nama: '', deskripsi: '' }; editingItem.value = null; showModal.value = true }
const openEditModal = (item) => { form.value = { ...item }; editingItem.value = item; showModal.value = true }
const openDeleteModal = (item) => { deletingItem.value = item; showDeleteModal.value = true }

const saveData = async () => {
    if (!form.value.nama) return
    saving.value = true
    try {
        // Masukkan dummy icon string, karena AutoIcon yang akan mengurus di tampilan
        const payload = { ...form.value, icon: 'auto' }

        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/kategori-layanan/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
        } else {
            await $fetch(`${API_BASE}/admin/kategori-layanan`, { method: 'POST', body: payload, credentials: 'include' })
        }
        showModal.value = false; await fetchData()
        toast.add({ title: 'Berhasil', color: 'success' })
    } catch (err) { toast.add({ title: 'Gagal menyimpan', color: 'error' }) }
    finally { saving.value = false }
}

const confirmDelete = async () => {
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/kategori-layanan/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        showDeleteModal.value = false; await fetchData()
        toast.add({ title: 'Berhasil dihapus', color: 'success' })
    } catch (err) { toast.add({ title: 'Gagal menghapus', color: 'error' }) }
    finally { saving.value = false }
}
onMounted(() => fetchData())
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Kategori Layanan" description="Kelola kategori utama layanan publik."
            action-label="Tambah Kategori" @action="openCreateModal" />

        <UCard :ui="{ body: 'p-0' }">
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
            </div>

            <div v-else-if="list.length === 0" class="text-center py-12">
                <UIcon name="i-tabler-category" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                <p class="text-slate-500">Belum ada kategori layanan.</p>
            </div>

            <div v-else class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
                            <th class="py-3 px-6 text-xs font-bold text-slate-500 uppercase w-16">Ikon</th>
                            <th class="py-3 px-6 text-xs font-bold text-slate-500 uppercase">Nama Kategori</th>
                            <th class="py-3 px-6 text-xs font-bold text-slate-500 uppercase">Deskripsi</th>
                            <th v-if="canUpdate || canDelete"
                                class="py-3 px-6 text-xs font-bold text-slate-500 uppercase text-right">Aksi</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
                        <tr v-for="item in list" :key="item.id"
                            class="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-4 px-6">
                                <div
                                    class="w-10 h-10 rounded-lg bg-sky-50 dark:bg-sky-950 flex items-center justify-center">
                                    <AutoIcon :text="item.nama" size="w-5 h-5 text-sky-600" />
                                </div>
                            </td>
                            <td class="py-4 px-6 font-semibold text-slate-900 dark:text-white">{{ item.nama }}</td>
                            <td class="py-4 px-6 text-sm text-slate-500">{{ item.deskripsi || '-' }}</td>
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

        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Kategori' : 'Tambah Kategori'"
            :is-edit="!!editingItem" :loading="saving" @submit="saveData" @cancel="showModal = false">
            <div class="space-y-4">
                <UFormField label="Nama Kategori" required>
                    <UInput v-model="form.nama" placeholder="Contoh: Kependudukan, Kesehatan, Pekerjaan..." size="lg"
                        class="w-full" />
                </UFormField>
                <UFormField label="Deskripsi Singkat">
                    <UInput v-model="form.deskripsi" placeholder="Penjelasan singkat kategori" size="lg"
                        class="w-full" />
                </UFormField>

                <div
                    class="p-4 rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50 flex items-center gap-4">
                    <div class="w-12 h-12 rounded-lg bg-white shadow-sm flex items-center justify-center text-sky-600">
                        <AutoIcon :text="form.nama" size="w-6 h-6" />
                    </div>
                    <div>
                        <p class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-1">Preview Ikon Otomatis
                        </p>
                        <p class="text-sm text-slate-600 dark:text-slate-400">Ikon akan menyesuaikan dengan teks nama
                            kategori.</p>
                    </div>
                </div>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>