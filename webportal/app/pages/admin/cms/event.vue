<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Kelola Event'
})

const toast = useToast()
const API_BASE = 'http://localhost:8080/api'

// ==== State ====
const agendaList = ref([])
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const page = ref(1)
const perPage = 10

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

// State untuk Upload Gambar
const uploading = ref(false)
const fileInput = ref(null)

const triggerFileInput = () => {
    fileInput.value?.click()
}

const handleFileUpload = async (event) => {
    const file = event.target.files?.[0]
    if (!file) return

    uploading.value = true
    try {
        const formData = new FormData()
        formData.append('file', file)

        const response = await $fetch('http://localhost:8080/api/admin/upload', {
            method: 'POST',
            body: formData,
            credentials: 'include'
        })

        form.value.gambar = response.url
        toast.add({ title: 'Berhasil', description: 'Gambar berhasil diunggah', icon: 'i-tabler-circle-check', color: 'success' })
    } catch (err) {
        console.error("Detail Error:", err)
        toast.add({
            title: 'Gagal upload',
            description: err?.data?.statusMessage || err.message || 'Terjadi kesalahan jaringan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
    } finally {
        uploading.value = false
        if (fileInput.value) fileInput.value.value = ''
    }
}

// ==== Form ====
const emptyForm = () => ({
    id: null,
    judul: '',
    deskripsi: '',
    tanggal: '',
    lokasi: '',
    gambar: '',
    kategori: 'Budaya',
    status: 'Draft',
    link_event: ''
})
const form = ref(emptyForm())

// ==== Computed ====
const filteredAgenda = computed(() => {
    return agendaList.value.filter(item =>
        item.judul.toLowerCase().includes(search.value.toLowerCase())
    )
})

const paginatedAgenda = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredAgenda.value.slice(start, start + perPage)
})

watch(search, () => { page.value = 1 })

// ==== Fetch ====
const fetchAgenda = async () => {
    loading.value = true
    try {
        const response = await $fetch(`${API_BASE}/admin/agenda`, { credentials: 'include' })
        agendaList.value = response.agenda || []
    } catch (err) {
        toast.add({ title: 'Gagal memuat data', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        loading.value = false
    }
}

// ==== Modal ====
const openCreateModal = () => {
    form.value = emptyForm()
    editingItem.value = null
    showModal.value = true
}

const openEditModal = (item) => {
    form.value = { ...item }
    if (form.value.tanggal) form.value.tanggal = form.value.tanggal.substring(0, 16)
    editingItem.value = item
    showModal.value = true
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

// ==== Save ====
const saveAgenda = async () => {
    if (!form.value.judul || !form.value.tanggal) {
        toast.add({ title: 'Validasi gagal', description: 'Judul dan Tanggal wajib diisi', icon: 'i-tabler-alert-circle', color: 'error' })
        return
    }
    saving.value = true
    try {
        const payload = {
            judul: form.value.judul,
            deskripsi: form.value.deskripsi,
            tanggal: form.value.tanggal + ':00',
            lokasi: form.value.lokasi,
            gambar: form.value.gambar,
            kategori: form.value.kategori,
            status: form.value.status,
            link_event: form.value.link_event
        }

        if (editingItem.value) {
            await $fetch(`${API_BASE}/admin/agenda/${editingItem.value.id}`, { method: 'PUT', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Agenda diperbarui', icon: 'i-tabler-circle-check', color: 'success' })
        } else {
            await $fetch(`${API_BASE}/admin/agenda`, { method: 'POST', body: payload, credentials: 'include' })
            toast.add({ title: 'Berhasil', description: 'Agenda ditambahkan', icon: 'i-tabler-circle-check', color: 'success' })
        }
        showModal.value = false
        await fetchAgenda()
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage || 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        saving.value = false
    }
}

// ==== Delete ====
const confirmDelete = async () => {
    if (!deletingItem.value) return
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/agenda/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        toast.add({ title: 'Berhasil', description: 'Agenda dihapus', icon: 'i-tabler-circle-check', color: 'success' })
        showDeleteModal.value = false
        deletingItem.value = null
        await fetchAgenda()
    } catch (err) {
        toast.add({ title: 'Gagal menghapus', description: 'Terjadi kesalahan', icon: 'i-tabler-alert-circle', color: 'error' })
    } finally {
        saving.value = false
    }
}

const closeModal = () => { showModal.value = false; editingItem.value = null; form.value = emptyForm() }
const formatDate = (dateStr) => {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => { fetchAgenda() })

// Cek apakah ada agenda yang sudah terlewat
const pastEvents = computed(() => {
    const now = new Date()
    return agendaList.value.filter(item => new Date(item.tanggal) < now)
})
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Kelola Event" description="Kelola jadwal kegiatan dan acara Pemkab Tulungagung."
            action-label="Tambah Event" @action="openCreateModal" />

        <!--Notif-->
        <UAlert v-if="pastEvents.length > 0" color="error" variant="subtle" icon="i-tabler-alert-triangle"
            :title="`Perhatian: Ada ${pastEvents.length} event yang sudah terlewat!`"
            description="Segera periksa daftar event di bawah dan hapus event yang sudah selesai!" />

        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari nama event..." size="lg" />
            </div>
        </UCard>


        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-calendar-event" class="text-sky-500" /> Daftar Event
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredAgenda.length }} event</span>
                </div>
            </template>

            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            </div>
            <AdminEmptyState v-else-if="filteredAgenda.length === 0" icon="i-tabler-calendar-off"
                title="Belum ada event" description="Klik Tambah Event untuk memulai." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Event</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Waktu</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase">Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in paginatedAgenda" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50">
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div class="w-16 h-10 rounded overflow-hidden bg-slate-200 shrink-0">
                                        <img v-if="item.gambar" :src="item.gambar" class="w-full h-full object-cover" />
                                    </div>
                                    <div>
                                        <p class="font-medium text-slate-900 dark:text-white line-clamp-1">{{ item.judul
                                        }}</p>
                                        <p class="text-xs text-slate-500 flex items-center gap-1">
                                            <UIcon name="i-tabler-map-pin" /> {{ item.lokasi || '-' }}
                                        </p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600">{{ formatDate(item.tanggal) }}</td>
                            <td class="py-3 px-4">
                                <UBadge :color="item.status === 'Publish' ? 'success' : 'warning'" variant="subtle"
                                    size="sm">{{
                                        item.status }}</UBadge>
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
                <AdminPagination v-model:page="page" :total="filteredAgenda.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <!-- MODAL FORM -->
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Event' : 'Tambah Event'"
            :is-edit="!!editingItem" :loading="saving" submit-label="Simpan Event" @submit="saveAgenda"
            @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Judul Event" required>
                    <UInput v-model="form.judul" placeholder="Nama acara..." size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Waktu Pelaksanaan" required>
                    <UInput v-model="form.tanggal" type="datetime-local" size="lg" class="w-full" />
                </UFormField>

                <div class="grid grid-cols-2 gap-4">
                    <UFormField label="Lokasi">
                        <UInput v-model="form.lokasi" placeholder="Tempat..." size="lg" class="w-full" />
                    </UFormField>
                    <UFormField label="Kategori">
                        <USelect v-model="form.kategori"
                            :items="['Budaya', 'Ekonomi', 'Pendidikan', 'Kesehatan', 'Pemerintahan']" size="lg"
                            class="w-full" />
                    </UFormField>
                </div>

                <!-- Upload Gambar -->
                <UFormField label="Gambar Event">
                    <div class="space-y-3">
                        <div v-if="form.gambar"
                            class="relative w-full aspect-video rounded-lg overflow-hidden border border-slate-200 dark:border-slate-700 bg-slate-100">
                            <img :src="form.gambar" class="w-full h-full object-cover" />
                            <div class="absolute top-2 right-2">
                                <UButton icon="i-tabler-trash" color="error" variant="solid" size="sm"
                                    @click="form.gambar = ''" />
                            </div>
                        </div>

                        <div v-else
                            class="relative border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-6 text-center hover:border-sky-500 transition-colors bg-slate-50 dark:bg-slate-900/50 cursor-pointer"
                            @click="triggerFileInput">
                            <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/webp" class="hidden"
                                @change="handleFileUpload" />
                            <UIcon :name="uploading ? 'i-tabler-loader-2' : 'i-tabler-cloud-upload'"
                                class="w-8 h-8 text-slate-400 mx-auto mb-2" :class="{ 'animate-spin': uploading }" />
                            <p class="text-sm font-medium text-slate-700 dark:text-slate-300">
                                {{ uploading ? 'Sedang mengunggah...' : 'Klik untuk upload gambar' }}
                            </p>
                            <p class="text-xs text-slate-500 mt-1">PNG, JPG, WebP (Maks. 5MB)</p>
                        </div>
                    </div>
                </UFormField>

                <!-- Link Event -->
                <UFormField label="Link Event">
                    <UInput v-model="form.link_event" placeholder="https://instagram.com/p/..." icon="i-tabler-link"
                        size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Deskripsi Singkat">
                    <UTextarea v-model="form.deskripsi" placeholder="Detail acara..." :rows="3" size="lg"
                        class="w-full" />
                </UFormField>

                <UFormField label="Status">
                    <USelect v-model="form.status" :items="['Draft', 'Publish']" size="lg" class="w-full" />
                </UFormField>
            </div>
        </AdminModalForm>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.judul" :loading="saving"
            @confirm="confirmDelete" />
    </div>
</template>