<script setup>
import { VueDraggable } from 'vue-draggable-plus'
definePageMeta({ layout: 'admin', title: 'Menu Navigasi' })
const toast = useToast()
const { can } = useAuth()
const canCreate = computed(() => can('pengaturan.menu', 'create'))
const canUpdate = computed(() => can('pengaturan.menu', 'update'))
const canDelete = computed(() => can('pengaturan.menu', 'delete'))
const list = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)
const form = ref({ id: null, label: '', url: '', deskripsi: '', icon: 'i-tabler-link', color: 'text-sky-600', grup: 'Lainnya', is_active: true })
const groups = ['Berita', 'Layanan Publik', 'Pemerintahan', 'Lainnya']
const emptyForm = () => ({ id: null, label: '', url: '', deskripsi: '', icon: 'i-tabler-link', color: 'text-sky-600', grup: 'Lainnya', is_active: true })

const fetchData = async () => {
    loading.value = true
    try { const response = await $fetch('http://localhost:8080/api/admin/menu-navigasi', { credentials: 'include' }); list.value = response.menu || [] }
    catch (err) { toast.add({ title: 'Gagal memuat menu', color: 'error' }) }
    finally { loading.value = false }
}
const openCreate = () => { form.value = emptyForm(); editingItem.value = null; showModal.value = true }
const openEdit = item => { form.value = { ...item }; editingItem.value = item; showModal.value = true }
const openDelete = item => { deletingItem.value = item; showDeleteModal.value = true }
const saveData = async () => {
    if (!form.value.label.trim()) { toast.add({ title: 'Label menu wajib diisi', color: 'error' }); return }
    saving.value = true
    try {
        const url = editingItem.value ? `http://localhost:8080/api/admin/menu-navigasi/${editingItem.value.id}` : 'http://localhost:8080/api/admin/menu-navigasi'
        await $fetch(url, { method: editingItem.value ? 'PUT' : 'POST', body: form.value, credentials: 'include' })
        showModal.value = false; await fetchData(); toast.add({ title: 'Menu berhasil disimpan', color: 'success' })
    } catch (err) { toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage || 'Terjadi kesalahan', color: 'error' }) }
    finally { saving.value = false }
}
const deleteData = async () => {
    saving.value = true
    try { await $fetch(`http://localhost:8080/api/admin/menu-navigasi/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' }); showDeleteModal.value = false; await fetchData(); toast.add({ title: 'Menu dihapus', color: 'success' }) }
    catch (err) { toast.add({ title: 'Gagal menghapus', color: 'error' }) }
    finally { saving.value = false }
}
const onReorder = async () => {
    try {
        const ids = list.value.map(item => item.id)
        await $fetch('http://localhost:8080/api/admin/menu-navigasi/reorder', {
            method: 'PUT',
            body: { ids },
            credentials: 'include'
        })
        toast.add({ title: 'Urutan menu diperbarui', color: 'success', duration: 1500 })
    } catch {
        toast.add({ title: 'Gagal menyimpan urutan', color: 'error' })
        await fetchData()
    }
}
onMounted(fetchData)
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Menu Navigasi" description="Atur tautan yang tampil pada navigasi landing page."
            :action-label="canCreate ? 'Tambah Menu' : undefined" @action="openCreate" />
        <UCard :ui="{ body: 'p-0' }">
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
            </div>
            <AdminEmptyState v-else-if="!list.length" icon="i-tabler-menu-2" title="Belum ada menu"
                description="Tambahkan menu untuk ditampilkan di landing page." />
            <VueDraggable v-else v-model="list" :animation="150" handle=".drag-handle" tag="div" @end="onReorder">
                <div v-for="item in list" :key="item.id"
                    class="flex items-center gap-4 px-6 py-4 border-b border-slate-100 dark:border-slate-800">
                    <UIcon name="i-tabler-grip-vertical" class="drag-handle w-5 h-5 text-slate-400 cursor-grab" />
                    <UIcon :name="item.icon || 'i-tabler-link'" class="w-5 h-5 text-sky-600" />
                    <div class="min-w-0 flex-1">
                        <p class="font-semibold text-slate-900 dark:text-white">{{ item.label }}</p>
                        <p class="text-xs text-slate-500 truncate">{{ item.grup }} - {{ item.url }}</p>
                    </div>
                    <UBadge :color="item.is_active ? 'success' : 'neutral'" variant="subtle">{{ item.is_active ? 'Aktif'
                        : 'Nonaktif' }}</UBadge>
                    <div class="flex gap-2">
                        <AdminButtonEdit size="xs" @click="openEdit(item)" />
                        <AdminButtonDelete size="xs" @click="openDelete(item)" />
                    </div>
                </div>
            </VueDraggable>
        </UCard>
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Menu' : 'Tambah Menu'"
            :is-edit="!!editingItem" :loading="saving" @submit="saveData" @cancel="showModal = false">
            <div class="space-y-4">
                <UFormField label="Label Menu" required>
                    <UInput v-model="form.label" placeholder="Contoh: Profil Daerah" size="lg" class="w-full" />
                </UFormField>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <UFormField label="URL">
                        <UInput v-model="form.url" placeholder="/profil" size="lg" class="w-full" />
                    </UFormField>
                    <UFormField label="Grup">
                        <USelect v-model="form.grup" :items="groups" size="lg" class="w-full" />
                    </UFormField>
                </div>
                <UFormField label="Deskripsi">
                    <UTextarea v-model="form.deskripsi" placeholder="Keterangan singkat menu" :rows="2"
                        class="w-full" />
                </UFormField>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <UFormField label="Ikon Tabler">
                        <UInput v-model="form.icon" placeholder="i-tabler-link" size="lg" class="w-full" />
                    </UFormField>
                    <UFormField label="Warna Ikon">
                        <UInput v-model="form.color" placeholder="text-sky-600" size="lg" class="w-full" />
                    </UFormField>
                </div>
                <UCheckbox v-model="form.is_active" label="Tampilkan di landing page" />
            </div>
        </AdminModalForm>
        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.label" :loading="saving"
            @confirm="deleteData" />
    </div>
</template>
