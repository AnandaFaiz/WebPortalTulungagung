<script setup lang="ts">
import { adminPermissionGroups, adminPermissionKeys } from '~/utils/admin-navigation'

definePageMeta({
    layout: 'admin',
    title: 'Role & Hak Akses'
})

const toast = useToast()

// ==== State ====
const loading = ref(true)
const permissions = ref<Record<string, Record<string, boolean>>>({})
const toggling = ref<string | null>(null)  // key: "role:permission"
const { user } = useAuth()

// ==== Definisi Role ====
const roles = [
    {
        key: 'super_admin',
        label: 'Super Admin',
        description: 'Akses penuh ke seluruh sistem. Permission-nya tidak dapat diubah.',
        color: 'primary',
        icon: 'i-tabler-shield-check',
        locked: true
    },
    {
        key: 'admin_opd',
        label: 'Admin OPD',
        description: 'Mengelola konten dan layanan untuk OPD tertentu saja.',
        color: 'warning',
        icon: 'i-tabler-building-community',
        locked: false
    },
    {
        key: 'editor',
        label: 'Editor',
        description: 'Fokus pada konten editorial di seluruh portal.',
        color: 'success',
        icon: 'i-tabler-edit',
        locked: false
    }
]

// ==== Grup permission untuk tabel ====
const legacyPermissionGroups = [
    {
        label: 'Dashboard',
        items: [{ key: 'dashboard', label: 'Dashboard' }]
    },
    {
        label: 'Konten (CMS)',
        items: [
            { key: 'cms.halaman', label: 'Halaman Statis' },
            { key: 'cms.berita', label: 'Berita' },
            { key: 'cms.pengumuman', label: 'Pengumuman' },
            { key: 'cms.galeri', label: 'Galeri' },
            { key: 'cms.sambutan', label: 'Sambutan Bupati' },
            { key: 'cms.banner', label: 'Banner' }
        ]
    },
    {
        label: 'Profil & Pemerintahan',
        items: [
            { key: 'pemerintahan.profil', label: 'Profil Daerah' },
            { key: 'pemerintahan.struktur', label: 'Struktur Organisasi' },
            { key: 'pemerintahan.pejabat', label: 'Profil Pejabat' },
            { key: 'pemerintahan.opd', label: 'OPD' },
            { key: 'pemerintahan.dprd', label: 'DPRD' },
            { key: 'pemerintahan.ppid', label: 'PPID' },
            { key: 'pemerintahan.jdih', label: 'JDIH' }
        ]
    },
    {
        label: 'Layanan Publik',
        items: [
            { key: 'layanan.shortcut', label: 'Shortcut Layanan' },
            { key: 'layanan.kategori', label: 'Kategori Layanan' },
            { key: 'layanan.daftar', label: 'Daftar Layanan' },
            { key: 'layanan.mpp', label: 'Informasi MPP' },
            { key: 'layanan.darurat', label: 'Layanan Darurat' }
        ]
    },
    {
        label: 'Interaksi',
        items: [
            { key: 'interaksi.medsos', label: 'Media Sosial' },
            { key: 'interaksi.tiket', label: 'Tiket Pengaduan' },
            { key: 'interaksi.kategori', label: 'Kategori & Status' }
        ]
    },
    {
        label: 'Potensi Daerah',
        items: [
            { key: 'potensi.pariwisata', label: 'Pariwisata' },
            { key: 'potensi.umkm', label: 'UMKM' },
            { key: 'potensi.investasi', label: 'Investasi' }
        ]
    },
    {
        label: 'Sistem',
        items: [
            { key: 'aplikasi', label: 'Aplikasi Daerah' },
            { key: 'media', label: 'Media Library' },
            { key: 'pengaturan.identitas', label: 'Identitas Situs' },
            { key: 'pengaturan.menu', label: 'Menu Navigasi' },
            { key: 'pengaturan.kontak', label: 'Kontak' },
            { key: 'pengguna.daftar', label: 'Daftar Pengguna' },
            { key: 'pengguna.role', label: 'Role & Hak Akses' },
            { key: 'pengguna.log', label: 'Log Aktivitas' }
        ]
    }
]

const permissionGroups = adminPermissionGroups

// ==== Fetch permissions dari API ====
const fetchPermissions = async () => {
    loading.value = true
    try {
        if (user.value?.role === 'super_admin') {
            await $fetch('http://localhost:8080/api/admin/permissions/sync', {
                method: 'POST',
                body: { keys: adminPermissionKeys },
                credentials: 'include'
            })
        }
        const response = await $fetch<{ permissions: any[] }>('http://localhost:8080/api/admin/permissions', { credentials: 'include' })

        // Konversi array → object nested { role: { permission_key: allowed } }
        const map: Record<string, Record<string, boolean>> = {
            super_admin: {},
            admin_opd: {},
            editor: {}
        }
        for (const p of response.permissions) {
            if (!map[p.role]) map[p.role] = {}
            map[p.role]![p.permission_key] = p.allowed
        }
        permissions.value = map
    } catch (err: any) {
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

// ==== Cek izin ====
const hasPermission = (roleKey: string, permissionKey: string): boolean => {
    return permissions.value[roleKey]?.[permissionKey] === true
}

// ==== Toggle permission ====
const togglePermission = async (role: any, permissionKey: string) => {
    // Cegah kalau role terkunci
    if (role.locked) return

    const currentValue = hasPermission(role.key, permissionKey)
    const newValue = !currentValue
    const toggleKey = `${role.key}:${permissionKey}`

    // Cegah klik ganda
    if (toggling.value === toggleKey) return

    toggling.value = toggleKey

    // Optimistic update — update UI duluan, biar terasa cepat
    // Optimistic update — update UI duluan, biar terasa cepat
    const rolePerms = permissions.value[role.key] || {}
    rolePerms[permissionKey] = newValue
    permissions.value[role.key] = rolePerms

    try {
        await $fetch(`http://localhost:8080/api/admin/permissions/${role.key}/${permissionKey}`, {
            method: 'PATCH',
            body: { allowed: newValue },
            credentials: 'include'
        })

        toast.add({
            title: 'Berhasil disimpan',
            description: `${role.label}: ${newValue ? 'Diizinkan' : 'Tidak diizinkan'}`,
            icon: 'i-tabler-circle-check',
            color: 'success',
            duration: 2000
        })
    } catch (err: any) {

        // Rollback kalau gagal
        const rolePerms = permissions.value[role.key] || {}
        rolePerms[permissionKey] = currentValue
        permissions.value[role.key] = rolePerms

        toast.add({
            title: 'Gagal menyimpan',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
    } finally {
        toggling.value = null
    }
}

// ==== Fetch saat mount ====
onMounted(() => {
    fetchPermissions()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <UCard>
            <div class="flex items-center justify-between flex-wrap gap-4">
                <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Role & Hak Akses</h2>
                    <p class="text-slate-500 mt-1">
                        Klik ikon centang/silang untuk mengubah hak akses. Perubahan langsung tersimpan ke database.
                    </p>
                </div>
                <UButton icon="i-tabler-refresh" label="Refresh" color="neutral" variant="outline" :loading="loading"
                    @click="fetchPermissions" />
            </div>
        </UCard>

        <!-- Kartu Role -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <UCard v-for="role in roles" :key="role.key" :class="role.locked ? 'opacity-90' : ''">
                <div class="flex items-start gap-4">
                    <div class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0"
                        :class="`bg-${role.color}-100 dark:bg-${role.color}-950`">
                        <UIcon :name="role.icon" class="w-6 h-6" :class="`text-${role.color}-600`" />
                    </div>
                    <div class="flex-1">
                        <div class="flex items-center gap-2 mb-1">
                            <h3 class="font-bold text-slate-900 dark:text-white">{{ role.label }}</h3>
                            <UBadge v-if="role.locked" color="neutral" variant="subtle" size="xs">
                                <UIcon name="i-tabler-lock" class="w-3 h-3 mr-0.5" />
                                Terkunci
                            </UBadge>
                        </div>
                        <p class="text-xs text-slate-500 leading-relaxed">{{ role.description }}</p>
                    </div>
                </div>
            </UCard>
        </div>

        <!-- Matrix Permission -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between flex-wrap gap-2">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-shield-lock" class="text-sky-500" />
                        Matriks Hak Akses
                    </h3>
                    <div class="flex items-center gap-4 text-xs text-slate-500">
                        <div class="flex items-center gap-1.5">
                            <UIcon name="i-tabler-circle-check-filled" class="text-emerald-500 w-4 h-4" />
                            <span>Diizinkan</span>
                        </div>
                        <div class="flex items-center gap-1.5">
                            <UIcon name="i-tabler-circle-x-filled" class="text-red-500" />
                            <span>Tidak diizinkan</span>
                        </div>
                        <div class="flex items-center gap-1.5">
                            <UIcon name="i-tabler-lock" class="text-slate-400 w-4 h-4" />
                            <span>Terkunci</span>
                        </div>
                    </div>
                </div>
            </template>

            <!-- Loading State -->
            <div v-if="loading" class="py-16 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat data permission...</p>
            </div>

            <!-- Matrix -->
            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b-2 border-slate-200 dark:border-slate-700">
                            <th
                                class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider sticky left-0 bg-white dark:bg-slate-800 z-10">
                                Fitur
                            </th>
                            <th v-for="role in roles" :key="role.key"
                                class="text-center py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider min-w-40">
                                <div class="flex items-center justify-center gap-1.5">
                                    <UIcon :name="role.icon" class="w-4 h-4" :class="`text-${role.color}-600`" />
                                    {{ role.label }}
                                </div>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <template v-for="group in permissionGroups" :key="group.label">
                            <!-- Group Header -->
                            <tr class="bg-sky-50 dark:bg-sky-950/40">
                                <td :colspan="roles.length + 1"
                                    class="py-2 px-4 text-xs font-bold text-sky-700 dark:text-sky-400 uppercase tracking-wider">
                                    {{ group.label }}
                                </td>
                            </tr>

                            <!-- Items -->
                            <tr v-for="item in group.items" :key="item.key"
                                class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                                <td
                                    class="py-2.5 px-4 text-sm text-slate-700 dark:text-slate-300 sticky left-0 bg-white dark:bg-slate-800">
                                    {{ item.label }}
                                </td>

                                <td v-for="role in roles" :key="role.key" class="py-2.5 px-4 text-center">
                                    <!-- Kalau role locked, tampilkan ikon statis -->
                                    <div v-if="role.locked" class="flex justify-center">
                                        <UIcon name="i-tabler-lock" class="w-4 h-4 text-slate-300"
                                            title="Permission terkunci" />
                                    </div>

                                    <!-- Kalau tidak locked, tampilkan tombol toggle -->
                                    <button v-else type="button" :disabled="toggling === `${role.key}:${item.key}`"
                                        :aria-label="`${hasPermission(role.key, item.key) ? 'Cabut' : 'Beri'} izin ${item.label} untuk ${role.label}`"
                                        class="relative inline-flex items-center justify-center w-9 h-9 rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:bg-slate-100 dark:hover:bg-slate-700"
                                        @click="togglePermission(role, item.key)">
                                        <!-- Loading -->
                                        <UIcon v-if="toggling === `${role.key}:${item.key}`" name="i-tabler-loader-2"
                                            class="w-5 h-5 text-sky-500 animate-spin" />

                                        <!-- Centang (diizinkan) -->
                                        <UIcon v-else-if="hasPermission(role.key, item.key)"
                                            name="i-tabler-circle-check-filled" class="w-5 h-5 text-emerald-500" />

                                        <!-- Silang (tidak diizinkan) -->
                                        <UIcon v-else name="i-tabler-circle-x-filled" class="w-5 h-5 text-red-400/60" />
                                    </button>
                                </td>
                            </tr>
                        </template>
                    </tbody>
                </table>
            </div>

            <template #footer>
                <div class="flex items-start gap-2 text-xs text-slate-400">
                    <UIcon name="i-tabler-info-circle" class="w-4 h-4 mt-0.5 shrink-0" />
                    <p>
                        Klik ikon di kolom <strong>Admin OPD</strong> atau <strong>Editor</strong> untuk mengubah hak
                        akses.
                        Perubahan langsung tersimpan ke database dan berlaku untuk semua user dengan role tersebut.
                        Kolom <strong>Super Admin</strong> terkunci untuk mencegah Anda terkunci dari sistem.
                    </p>
                </div>
            </template>
        </UCard>

    </div>
</template>