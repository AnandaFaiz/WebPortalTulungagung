<script setup>
definePageMeta({ layout: 'admin', title: 'Manajemen Media' })
const toast = useToast()
const API_BASE = 'http://localhost:8080/api'
const { can } = useAuth()
const canCreate = computed(() => can('media', 'create'))

const list = ref([])
const loading = ref(true)
const uploadedMedia = ref('')
const search = ref('')

const showDeleteModal = ref(false)
const deletingItem = ref(null)
const saving = ref(false)

const fetchData = async () => {
    loading.value = true
    try {
        const res = await $fetch(`${API_BASE}/admin/media`, { credentials: 'include' })
        list.value = res.media || []
    } catch (err) {
        toast.add({ title: 'Gagal memuat', color: 'error' })
    } finally {
        loading.value = false
    }
}

const filteredMedia = computed(() => {
    if (!search.value) return list.value
    return list.value.filter(m => m.nama_file.toLowerCase().includes(search.value.toLowerCase()))
})

watch(uploadedMedia, async (url) => {
    if (!url) return
    await fetchData()
    uploadedMedia.value = ''
})

const copyUrl = (url) => {
    navigator.clipboard.writeText(url)
    toast.add({ title: 'Tersalin', description: 'URL gambar berhasil disalin', icon: 'i-tabler-copy', color: 'success' })
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

const confirmDelete = async () => {
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/media/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        showDeleteModal.value = false
        await fetchData()
        toast.add({ title: 'Berhasil dihapus', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menghapus', color: 'error' })
    } finally {
        saving.value = false
    }
}

onMounted(() => fetchData())
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Manajemen Media"
            description="Gudang penyimpanan gambar dan file untuk keperluan konten website.">
            <template #actions>
                <AdminImageUpload v-if="canCreate" v-model="uploadedMedia" label="file media" :aspect-ratio="1"
                    preview-class="w-12 aspect-square" />
            </template>
        </AdminPageHeader>

        <UCard>
            <div class="mb-6 max-w-md">
                <AdminSearch v-model="search" placeholder="Cari nama file..." />
            </div>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
            </div>

            <AdminEmptyState v-else-if="filteredMedia.length === 0" icon="i-tabler-photo-off"
                title="Belum ada file media" description="Klik Upload File untuk mulai mengunggah." />

            <!-- Grid Media -->
            <div v-else class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-4">
                <div v-for="item in filteredMedia" :key="item.id"
                    class="group relative bg-slate-100 dark:bg-slate-800 rounded-xl overflow-hidden border border-slate-200 dark:border-slate-700 aspect-square">
                    <img :src="item.url" :alt="item.nama_file" class="w-full h-full object-cover" />

                    <!-- Overlay Aksi -->
                    <div
                        class="absolute inset-0 bg-slate-900/70 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-2 p-2">
                        <p class="text-white text-xs text-center line-clamp-2 mb-2">{{ item.nama_file }}</p>
                        <div class="flex gap-2">
                            <UButton icon="i-tabler-copy" color="primary" variant="solid" size="sm"
                                @click="copyUrl(item.url)" title="Salin URL" />
                            <AdminButtonDelete size="sm" @click="openDeleteModal(item)" />
                        </div>
                    </div>
                </div>
            </div>
        </UCard>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama_file" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>