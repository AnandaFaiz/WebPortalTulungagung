<script setup>
import { VueDraggable } from 'vue-draggable-plus'
definePageMeta({ layout: 'admin', title: 'Layanan Darurat' })
const toast = useToast()
const API_BASE = 'http://localhost:8080/api'

const list = ref([])
const loading = ref(true)
const saving = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

const emptyForm = () => ({ id: null, nama: '', nomor: '', is_active: true })
const form = ref(emptyForm())

const fetchData = async () => {
    loading.value = true
    try {
        const res = await $fetch(`${API_BASE}/admin/layanan-darurat`, { credentials: 'include' })
        list.value = res.darurat || []
    } catch (err) { toast.add({ title: 'Gagal memuat', color: 'error' }) }
    finally { loading.value = false }
}

const openCreateModal = () => { form.value = emptyForm(); editingItem.value = null; showModal.value = true }
const openEditModal = (item) => { form.value = { ...item }; editingItem.value = item; showModal.value = true }
const openDeleteModal = (item) => { deletingItem.value = item; showDeleteModal.value = true }

const saveData = async () => {
    if (!form.value.nama || !form.value.nomor) return
    saving.value = true
    try {
        const payload = { ...form.value, icon: 'auto' }
        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/layanan-darurat/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
        } else {
            await $fetch(`${API_BASE}/admin/layanan-darurat`, { method: 'POST', body: payload, credentials: 'include' })
        }
        showModal.value = false; await fetchData(); toast.add({ title: 'Berhasil', color: 'success' })
    } catch (err) { toast.add({ title: 'Gagal simpan', color: 'error' }) }
    finally { saving.value = false }
}

const confirmDelete = async () => {
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/layanan-darurat/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        showDeleteModal.value = false; await fetchData(); toast.add({ title: 'Berhasil dihapus', color: 'success' })
    } catch (err) { toast.add({ title: 'Gagal hapus', color: 'error' }) }
    finally { saving.value = false }
}

const onReorder = async () => {
    try {
        const ids = list.value.map(i => i.id)
        await $fetch(`${API_BASE}/admin/layanan-darurat/reorder`, { method: 'PUT', body: { ids }, credentials: 'include' })
    } catch (err) { await fetchData() }
}

onMounted(() => fetchData())
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Layanan Darurat" description="Kelola nomor penting darurat." action-label="Tambah Nomor"
            @action="openCreateModal" />

        <UCard :ui="{ body: 'p-0 sm:p-0' }">
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
            </div>

            <div v-else-if="list.length === 0" class="text-center py-12">
                <UIcon name="i-tabler-phone-off" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                <p class="text-slate-500">Belum ada data layanan darurat.</p>
            </div>

            <div v-else class="overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr class="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-700">
                            <th class="w-12 py-4 px-6"></th>
                            <th class="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">Nama Layanan
                            </th>
                            <th class="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">Nomor
                                Telepon</th>
                            <th class="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider">Status</th>
                            <th class="py-4 px-6 text-xs font-bold text-slate-500 uppercase tracking-wider text-right">
                                Aksi</th>
                        </tr>
                    </thead>

                    <VueDraggable v-model="list" tag="tbody" handle=".drag-handle" :animation="200" @end="onReorder"
                        class="divide-y divide-slate-100 dark:divide-slate-800">
                        <tr v-for="item in list" :key="item.id"
                            class="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">

                            <!-- Kolom 1: Drag Handle -->
                            <td class="py-4 px-6">
                                <div
                                    class="drag-handle cursor-grab active:cursor-grabbing p-1.5 rounded-lg hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors w-max">
                                    <UIcon name="i-tabler-grip-vertical" class="w-5 h-5 text-slate-400" />
                                </div>
                            </td>

                            <!-- Kolom 2: Nama & Ikon -->
                            <td class="py-4 px-6">
                                <div class="flex items-center gap-4">
                                    <div
                                        class="w-12 h-12 rounded-xl bg-rose-50 dark:bg-rose-950/50 flex items-center justify-center shrink-0">
                                        <AutoIcon :text="item.nama" class="w-6 h-6 text-rose-500" />
                                    </div>
                                    <p class="font-bold text-slate-900 dark:text-white">{{ item.nama }}</p>
                                </div>
                            </td>

                            <!-- Kolom 3: Nomor (Dibuat lebih menonjol) -->
                            <td class="py-4 px-6">
                                <div
                                    class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-red-50 dark:bg-red-950/30 text-red-600 dark:text-red-400 border border-red-100 dark:border-red-900">
                                    <UIcon name="i-tabler-phone-call" class="w-4 h-4" />
                                    <span class="font-bold font-mono tracking-wider text-base">{{ item.nomor }}</span>
                                </div>
                            </td>

                            <!-- Kolom 4: Status -->
                            <td class="py-4 px-6">
                                <AdminStatusBadge :status="item.is_active ? 'Aktif' : 'Nonaktif'" />
                            </td>

                            <!-- Kolom 5: Aksi -->
                            <td class="py-4 px-6">
                                <div class="flex items-center justify-end gap-2">
                                    <AdminButtonEdit size="sm" @click="openEditModal(item)" />
                                    <AdminButtonDelete size="sm" @click="openDeleteModal(item)" />
                                </div>
                            </td>

                        </tr>
                    </VueDraggable>
                </table>
            </div>
        </UCard>

        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Nomor' : 'Tambah Nomor'"
            :is-edit="!!editingItem" :loading="saving" @submit="saveData" @cancel="showModal = false">
            <div class="space-y-4">
                <UFormField label="Nama Layanan" required>
                    <UInput v-model="form.nama" placeholder="Contoh: Pemadam Kebakaran" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Nomor Telepon" required>
                    <UInput v-model="form.nomor" placeholder="113" size="lg" class="w-full" />
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