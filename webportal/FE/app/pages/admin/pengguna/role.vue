<script setup lang="ts">
import { adminPermissionGroups } from '~/utils/admin-navigation'

definePageMeta({ layout: 'admin', title: 'Role & Hak Akses' })

interface AdminAccount {
    id: number
    nama: string
    email: string
    role: 'super_admin' | 'editor'
    is_active: boolean
}

type PermissionAction = 'create' | 'read' | 'update' | 'delete'
type PermissionFlags = Record<PermissionAction, boolean>
interface ApiPermission {
    permission_key: string
    can_create: boolean
    can_read: boolean
    can_update: boolean
    can_delete: boolean
}

const toast = useToast()
const { user, fetchPermissions: refreshAuthPermissions } = useAuth()
const loading = ref(true)
const saving = ref<string | null>(null)
const admins = ref<AdminAccount[]>([])
const selectedAdminId = ref<number | null>(null)
const permissions = ref<Record<string, PermissionFlags>>({})

const permissionGroups = adminPermissionGroups.filter((group) =>
    ['Konten (CMS)', 'Profil & Pemerintahan', 'Layanan Publik', 'Interaksi', 'Aplikasi', 'Media', 'Pengaturan'].includes(group.label),
)
const actions: Array<{ key: PermissionAction, label: string }> = [
    { key: 'create', label: 'Create' },
    { key: 'read', label: 'Read' },
    { key: 'update', label: 'Update' },
    { key: 'delete', label: 'Delete' },
]
const selectedAdmin = computed(() => admins.value.find((admin) => admin.id === selectedAdminId.value) || null)

const loadPermissionMatrix = async () => {
    loading.value = true
    try {
        const response = await $fetch<{
            admins: AdminAccount[]
            permissions: Record<string, ApiPermission[]>
        }>('http://localhost:8080/api/admin/permissions', { credentials: 'include' })
        admins.value = response.admins
        if (!selectedAdminId.value || !admins.value.some((admin) => admin.id === selectedAdminId.value)) {
            selectedAdminId.value = admins.value[0]?.id || null
        }
        permissions.value = {}
        for (const permission of response.permissions[String(selectedAdminId.value)] || []) {
            permissions.value[permission.permission_key] = {
                create: permission.can_create,
                read: permission.can_read,
                update: permission.can_update,
                delete: permission.can_delete,
            }
        }
    } catch (error: any) {
        toast.add({ title: 'Gagal memuat data', description: error?.data?.statusMessage || 'Terjadi kesalahan', color: 'error', icon: 'i-tabler-alert-circle' })
    } finally {
        loading.value = false
    }
}

const selectAdmin = async () => {
    await loadPermissionMatrix()
}

const hasPermission = (key: string, action: PermissionAction) => permissions.value[key]?.[action] === true

const togglePermission = async (key: string, action: PermissionAction) => {
    if (!selectedAdminId.value) return
    const toggleKey = `${key}:${action}`
    if (saving.value === toggleKey) return
    const current = hasPermission(key, action)
    saving.value = toggleKey
    permissions.value[key] = { create: false, read: false, update: false, delete: false, ...permissions.value[key], [action]: !current }
    try {
        await $fetch(`http://localhost:8080/api/admin/permissions/${selectedAdminId.value}/${encodeURIComponent(key)}/${action}`, {
            method: 'PATCH', body: { allowed: !current }, credentials: 'include',
        })
        if (selectedAdminId.value === user.value?.id) {
            await refreshAuthPermissions()
        }
        toast.add({ title: 'Permission disimpan', description: `${action.toUpperCase()} ${key} diperbarui`, color: 'success', duration: 1500 })
    } catch (error: any) {
        permissions.value[key][action] = current
        toast.add({ title: 'Gagal menyimpan', description: error?.data?.statusMessage || 'Terjadi kesalahan', color: 'error' })
    } finally {
        saving.value = null
    }
}

onMounted(loadPermissionMatrix)
</script>

<template>
    <div class="space-y-6">
        <UCard>
            <div class="flex flex-wrap items-end justify-between gap-4">
                <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Role & Hak Akses</h2>
                    <p class="text-slate-500 mt-1">Pilih admin untuk mengatur hak akses CRUD setiap modul.</p>
                </div>
                <UButton icon="i-tabler-refresh" label="Refresh" color="neutral" variant="outline" :loading="loading"
                    @click="loadPermissionMatrix" />
            </div>
            <div class="mt-5 max-w-xl">
                <label for="admin-selector"
                    class="block text-sm font-semibold text-slate-700 dark:text-slate-200 mb-2">Admin yang
                    dipilih</label>
                <select id="admin-selector" v-model="selectedAdminId"
                    class="w-full rounded-lg border border-slate-300 dark:border-slate-700 bg-white dark:bg-slate-900 px-3 py-2.5 text-sm"
                    @change="selectAdmin">
                    <option v-for="admin in admins" :key="admin.id" :value="admin.id">{{ admin.nama }} - {{ admin.role
                        === 'super_admin' ? 'Super Admin' : 'Editor' }}</option>
                </select>
            </div>
        </UCard>

        <UCard v-if="selectedAdmin">
            <template #header>
                <div class="flex items-center gap-3">
                    <UIcon :name="selectedAdmin.role === 'super_admin' ? 'i-tabler-shield-check' : 'i-tabler-edit'"
                        class="w-6 h-6 text-sky-500" />
                    <div>
                        <h3 class="font-bold text-slate-900 dark:text-white">{{ selectedAdmin.nama }}</h3>
                        <p class="text-xs text-slate-500">{{ selectedAdmin.email }} · {{ selectedAdmin.role ===
                            'super_admin' ? 'Super Admin' : 'Editor' }}</p>
                    </div>
                </div>
            </template>

            <div v-if="loading" class="py-16 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            </div>
            <div v-else class="overflow-x-auto">
                <table class="w-full min-w-180">
                    <thead>
                        <tr class="border-b-2 border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Label / Modul
                            </th>
                            <th v-for="action in actions" :key="action.key"
                                class="text-center py-3 px-4 text-xs font-bold text-slate-500 uppercase min-w-28">{{
                                    action.label }}
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <template v-for="group in permissionGroups" :key="group.label">
                            <tr class="bg-sky-50 dark:bg-sky-950/40">
                                <td colspan="5"
                                    class="py-2 px-4 text-xs font-bold text-sky-700 dark:text-sky-400 uppercase">{{
                                        group.label }}</td>
                            </tr>
                            <tr v-for="item in group.items" :key="item.key"
                                class="border-b border-slate-100 dark:border-slate-800">
                                <td class="py-3 px-4 text-sm text-slate-700 dark:text-slate-300">{{ item.label }}</td>
                                <td v-for="action in actions" :key="action.key" class="py-2 px-4 text-center">
                                    <button type="button"
                                        class="inline-flex items-center justify-center w-9 h-9 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 disabled:opacity-50"
                                        :disabled="saving === `${item.key}:${action.key}`"
                                        :aria-label="`${action.label} ${item.label}`"
                                        @click="togglePermission(item.key, action.key)">
                                        <UIcon v-if="saving === `${item.key}:${action.key}`" name="i-tabler-loader-2"
                                            class="w-5 h-5 text-sky-500 animate-spin" />
                                        <UIcon v-else
                                            :name="hasPermission(item.key, action.key) ? 'i-tabler-circle-check-filled' : 'i-tabler-circle-x-filled'"
                                            class="w-5 h-5"
                                            :class="hasPermission(item.key, action.key) ? 'text-emerald-500' : 'text-red-400/60'" />
                                    </button>
                                </td>
                            </tr>
                        </template>
                    </tbody>
                </table>
            </div>
        </UCard>
    </div>
</template>
