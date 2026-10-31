<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Daftar Pengguna'
})

const toast = useToast()
const { user: currentUser } = useAuth()

// ==== State ====
const users = ref([])
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const selectedRole = ref('all')
const selectedStatus = ref('all')
const perPage = 10
const page = ref(1)

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

// ==== Opsi Filter ====
const roleOptions = [
    { label: 'Semua Role', value: 'all' },
    { label: 'Super Admin', value: 'super_admin' },
    { label: 'Admin OPD', value: 'admin_opd' },
    { label: 'Editor', value: 'editor' }
]

const statusOptions = [
    { label: 'Semua Status', value: 'all' },
    { label: 'Aktif', value: 'aktif' },
    { label: 'Nonaktif', value: 'nonaktif' }
]

const opdOptions = [
    'Dinas Pendidikan',
    'Dinas Kesehatan',
    'Dinas Kependudukan dan Pencatatan Sipil',
    'Dinas Pekerjaan Umum dan Penataan Ruang',
    'Dinas Penanaman Modal dan PTSP',
    'Dinas Sosial',
    'Dinas Pertanian',
    'Dinas Pariwisata, Pemuda dan Olahraga',
    'Badan Perencanaan Pembangunan Daerah',
    'Badan Pengelola Keuangan dan Aset Daerah',
    'Sekretariat Daerah',
    'Sekretariat DPRD'
]

// ==== Form ====
const emptyForm = () => ({
    id: null,
    nama: '',
    email: '',
    password: '',
    role: 'editor',
    opd: '',
    is_active: true
})
const form = ref(emptyForm())

// ==== Computed: Filter ====
const filteredUsers = computed(() => {
    return users.value.filter(item => {
        const matchSearch =
            item.nama?.toLowerCase().includes(search.value.toLowerCase()) ||
            item.email?.toLowerCase().includes(search.value.toLowerCase())

        const matchRole = selectedRole.value === 'all' || item.role === selectedRole.value

        const matchStatus =
            selectedStatus.value === 'all' ||
            (selectedStatus.value === 'aktif' && item.is_active) ||
            (selectedStatus.value === 'nonaktif' && !item.is_active)

        return matchSearch && matchRole && matchStatus
    })
})

const paginatedUsers = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredUsers.value.slice(start, start + perPage)
})

watch([search, selectedRole, selectedStatus], () => {
    page.value = 1
})

// ==== Fetch ====
const fetchUsers = async () => {
    loading.value = true
    try {
        const requestFetch = useRequestFetch()
        const response = await requestFetch('http://localhost:8080/api/admin/users', { credentials: 'include' })
        users.value = response.users || []
    } catch (err) {
        toast.add({
            title: 'Gagal memuat data',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
    } finally {
        loading.value = false
    }
}

// ==== Modal Handlers ====
const openCreateModal = () => {
    form.value = emptyForm()
    editingItem.value = null
    showModal.value = true
}

const openEditModal = (item) => {
    form.value = {
        id: item.id,
        nama: item.nama,
        email: item.email,
        password: '',
        role: item.role,
        opd: item.opd || '',
        is_active: item.is_active
    }
    editingItem.value = item
    showModal.value = true
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

// ==== Save (Create / Update) ====
const saveUser = async () => {
    if (!form.value.nama.trim() || !form.value.email.trim()) {
        toast.add({
            title: 'Validasi gagal',
            description: 'Nama dan email wajib diisi',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
        return
    }

    if (!editingItem.value && !form.value.password) {
        toast.add({
            title: 'Validasi gagal',
            description: 'Password wajib diisi untuk user baru',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
        return
    }

    if (form.value.role === 'admin_opd' && !form.value.opd) {
        toast.add({
            title: 'Validasi gagal',
            description: 'Admin OPD wajib memilih OPD',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
        return
    }

    saving.value = true
    try {
        if (editingItem.value) {
            const payload = {
                nama: form.value.nama,
                email: form.value.email,
                role: form.value.role,
                opd: form.value.opd,
                is_active: form.value.is_active
            }
            if (form.value.password && form.value.password.trim()) {
                payload.password = form.value.password
            }

            await $fetch(`http://localhost:8080/api/admin/users/${editingItem.value.id}`, {
                method: 'PUT',
                body: payload,
                credentials: 'include'
            })

            toast.add({
                title: 'Berhasil',
                description: 'Data pengguna berhasil diperbarui',
                icon: 'i-tabler-circle-check',
                color: 'success',
                duration: 2500
            })
        } else {
            await $fetch('http://localhost:8080/api/admin/users', {
                method: 'POST',
                body: {
                    nama: form.value.nama,
                    email: form.value.email,
                    password: form.value.password,
                    role: form.value.role,
                    opd: form.value.opd
                },
                credentials: 'include'
            })

            toast.add({
                title: 'Berhasil',
                description: 'Pengguna baru berhasil ditambahkan',
                icon: 'i-tabler-circle-check',
                color: 'success',
                duration: 2500
            })
        }

        showModal.value = false
        await fetchUsers()
    } catch (err) {
        toast.add({
            title: 'Gagal menyimpan',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
    } finally {
        saving.value = false
    }
}

// ==== Delete ====
const confirmDelete = async () => {
    if (!deletingItem.value) return

    saving.value = true
    try {
        await $fetch(`http://localhost:8080/api/admin/users/${deletingItem.value.id}`, {
            method: 'DELETE',
            credentials: 'include'
        })

        toast.add({
            title: 'Berhasil',
            description: 'Pengguna berhasil dihapus',
            icon: 'i-tabler-circle-check',
            color: 'success',
            duration: 2500
        })

        showDeleteModal.value = false
        deletingItem.value = null
        await fetchUsers()
    } catch (err) {
        toast.add({
            title: 'Gagal menghapus',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
    } finally {
        saving.value = false
    }
}

// ==== Helpers ====
const isCurrentUser = (id) => currentUser.value?.id === id

const closeModal = () => {
    showModal.value = false
    editingItem.value = null
    form.value = emptyForm()
}

// ==== Init ====
onMounted(() => {
    fetchUsers()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Daftar Pengguna" description="Kelola akun pengguna admin, role, dan hak aksesnya."
            action-label="Tambah Pengguna" @action="openCreateModal" />

        <!-- Filter -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari nama atau email..." size="lg" />
                <USelect v-model="selectedRole" :items="roleOptions" icon="i-tabler-shield-lock" size="lg" />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <!-- Tabel -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-users" class="text-sky-500" />
                        Semua Pengguna
                    </h3>
                    <span class="text-sm text-slate-500">{{ filteredUsers.length }} pengguna</span>
                </div>
            </template>

            <!-- Loading -->
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data pengguna...</p>
            </div>

            <!-- Empty -->
            <AdminEmptyState v-else-if="filteredUsers.length === 0" icon="i-tabler-users-off"
                title="Tidak ada pengguna ditemukan" description="Coba ubah filter atau tambahkan pengguna baru." />

            <!-- Tabel -->
            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Pengguna</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Role</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                OPD</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in paginatedUsers" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-3">
                                    <div
                                        class="w-10 h-10 rounded-full bg-sky-500 flex items-center justify-center text-white font-bold shrink-0">
                                        {{ item.nama.charAt(0).toUpperCase() }}
                                    </div>
                                    <div>
                                        <div class="flex items-center gap-2">
                                            <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                            <UBadge v-if="isCurrentUser(item.id)" color="primary" variant="subtle"
                                                size="xs">
                                                Anda
                                            </UBadge>
                                        </div>
                                        <p class="text-xs text-slate-400">{{ item.email }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <AdminRoleBadge :role="item.role" />
                            </td>
                            <td class="py-3 px-4 text-sm text-slate-600 dark:text-slate-300">
                                {{ item.opd || '-' }}
                            </td>
                            <td class="py-3 px-4">
                                <AdminStatusBadge :status="item.is_active ? 'Aktif' : 'Nonaktif'" />
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center justify-end gap-2">
                                    <AdminButtonEdit size="xs" aria-label="Edit" @click="openEditModal(item)" />
                                    <AdminButtonDelete size="xs" aria-label="Hapus" :disabled="isCurrentUser(item.id)"
                                        @click="!isCurrentUser(item.id) && openDeleteModal(item)" />
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <template #footer>
                <AdminPagination v-model:page="page" :total="filteredUsers.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <!-- ==== Modal Form (Create / Edit) ==== -->
        <AdminModalForm v-model:open="showModal" :title="editingItem ? 'Edit Pengguna' : 'Tambah Pengguna Baru'"
            :description="editingItem ? 'Perbarui informasi pengguna di bawah ini.' : 'Isi form berikut untuk menambahkan pengguna baru.'"
            :is-edit="!!editingItem" :loading="saving"
            :submit-label="editingItem ? 'Simpan Perubahan' : 'Simpan Pengguna'" @submit="saveUser"
            @cancel="closeModal">
            <div class="space-y-4">
                <UFormField label="Nama Lengkap" required>
                    <UInput v-model="form.nama" placeholder="Contoh: Admin Dinas Kesehatan" icon="i-tabler-user"
                        size="lg" class="w-full" />
                </UFormField>

                <UFormField label="Email" required>
                    <UInput v-model="form.email" type="email" placeholder="nama@tulungagung.go.id" icon="i-tabler-mail"
                        size="lg" class="w-full" />
                </UFormField>

                <UFormField :label="editingItem ? 'Password Baru (opsional)' : 'Password'" :required="!editingItem"
                    :hint="editingItem ? 'Kosongkan jika tidak ingin mengubah password.' : 'Minimal 8 karakter.'">
                    <UInput v-model="form.password" type="password"
                        :placeholder="editingItem ? '••••••••' : 'Minimal 8 karakter'" icon="i-tabler-lock" size="lg"
                        class="w-full" />
                </UFormField>

                <UFormField label="Role" required>
                    <USelect v-model="form.role" :items="[
                        { label: 'Super Admin', value: 'super_admin' },
                        { label: 'Admin OPD', value: 'admin_opd' },
                        { label: 'Editor', value: 'editor' }
                    ]" size="lg" class="w-full" />
                </UFormField>

                <UFormField v-if="form.role === 'admin_opd'" label="OPD" required
                    hint="Pilih OPD tempat pengguna bertugas.">
                    <USelect v-model="form.opd" :items="opdOptions" icon="i-tabler-building-community" size="lg"
                        class="w-full" placeholder="Pilih OPD..." />
                </UFormField>

                <UFormField v-if="editingItem" label="Status Akun" hint="User nonaktif tidak dapat login.">
                    <USwitch v-model="form.is_active" :label="form.is_active ? 'Aktif' : 'Nonaktif'" />
                </UFormField>
            </div>
        </AdminModalForm>

        <!-- ==== Modal Konfirmasi Hapus ==== -->
        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="deletingItem?.nama" :loading="saving"
            @confirm="confirmDelete" />

    </div>
</template>