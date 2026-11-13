<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Profil Saya'
})

const { user } = useAuth()
const toast = useToast()
const savingProfile = ref(false)
const savingPassword = ref(false)

const profileForm = reactive({ nama: '', email: '' })
const passwordForm = reactive({ current_password: '', new_password: '', confirmation: '' })

const roleLabels = {
    super_admin: 'Super Admin',
    editor: 'Editor'
}

const syncProfileForm = () => {
    profileForm.nama = user.value?.nama || ''
    profileForm.email = user.value?.email || ''
}

const formatDate = (value) => value
    ? new Intl.DateTimeFormat('id-ID', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
    : 'Belum tersedia'

const saveProfile = async () => {
    if (!profileForm.nama.trim() || !profileForm.email.trim()) {
        toast.add({ title: 'Data belum lengkap', description: 'Nama dan email wajib diisi.', color: 'error', icon: 'i-tabler-alert-circle' })
        return
    }

    savingProfile.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/auth/profile', {
            method: 'PUT', body: { nama: profileForm.nama, email: profileForm.email }, credentials: 'include'
        })
        user.value = response.user
        toast.add({ title: 'Profil diperbarui', description: 'Data akun berhasil disimpan.', color: 'success', icon: 'i-tabler-circle-check' })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan profil', description: err?.data?.statusMessage || 'Terjadi kesalahan.', color: 'error', icon: 'i-tabler-alert-circle' })
    } finally {
        savingProfile.value = false
    }
}

const savePassword = async () => {
    if (passwordForm.new_password !== passwordForm.confirmation) {
        toast.add({ title: 'Konfirmasi tidak cocok', description: 'Ulangi password baru dengan benar.', color: 'error', icon: 'i-tabler-alert-circle' })
        return
    }

    savingPassword.value = true
    try {
        await $fetch('http://localhost:8080/api/auth/password', {
            method: 'PUT', body: {
                current_password: passwordForm.current_password,
                new_password: passwordForm.new_password
            }, credentials: 'include'
        })
        passwordForm.current_password = ''
        passwordForm.new_password = ''
        passwordForm.confirmation = ''
        toast.add({ title: 'Password diperbarui', description: 'Gunakan password baru saat login berikutnya.', color: 'success', icon: 'i-tabler-circle-check' })
    } catch (err) {
        toast.add({ title: 'Gagal mengubah password', description: err?.data?.statusMessage || 'Terjadi kesalahan.', color: 'error', icon: 'i-tabler-alert-circle' })
    } finally {
        savingPassword.value = false
    }
}

onMounted(syncProfileForm)
</script>

<template>
    <div class="space-y-6 max-w-5xl">
        <AdminPageHeader title="Profil Saya" description="Kelola identitas dan keamanan akun yang sedang digunakan." />

        <div class="grid grid-cols-1 lg:grid-cols-[minmax(0,1fr)_18rem] gap-6">
            <div class="space-y-6">
                <UCard>
                    <template #header>
                        <div class="flex items-center gap-3">
                            <div
                                class="w-10 h-10 rounded-full bg-sky-500 text-white flex items-center justify-center font-bold">
                                {{ user?.nama?.charAt(0)?.toUpperCase() || '?' }}
                            </div>
                            <div>
                                <h2 class="font-bold text-slate-900 dark:text-white">Informasi akun</h2>
                                <p class="text-sm text-slate-500">Nama dan email untuk komunikasi internal.</p>
                            </div>
                        </div>
                    </template>
                    <div class="space-y-4">
                        <UFormField label="Nama lengkap" required>
                            <UInput v-model="profileForm.nama" icon="i-tabler-user" size="lg" class="w-full" />
                        </UFormField>
                        <UFormField label="Email" required>
                            <UInput v-model="profileForm.email" type="email" icon="i-tabler-mail" size="lg"
                                class="w-full" />
                        </UFormField>
                    </div>
                    <template #footer>
                        <div class="flex justify-end">
                            <UButton label="Simpan profil" icon="i-tabler-device-floppy" :loading="savingProfile"
                                @click="saveProfile" />
                        </div>
                    </template>
                </UCard>

                <UCard>
                    <template #header>
                        <div>
                            <h2 class="font-bold text-slate-900 dark:text-white">Keamanan akun</h2>
                            <p class="text-sm text-slate-500 mt-1">Ganti password secara berkala dan jangan
                                membagikannya.</p>
                        </div>
                    </template>
                    <div class="space-y-4">
                        <UFormField label="Password saat ini" required>
                            <UInput v-model="passwordForm.current_password" type="password" icon="i-tabler-lock"
                                size="lg" class="w-full" />
                        </UFormField>
                        <UFormField label="Password baru" hint="Minimal 8 karakter." required>
                            <UInput v-model="passwordForm.new_password" type="password" icon="i-tabler-key" size="lg"
                                class="w-full" />
                        </UFormField>
                        <UFormField label="Ulangi password baru" required>
                            <UInput v-model="passwordForm.confirmation" type="password" icon="i-tabler-key" size="lg"
                                class="w-full" />
                        </UFormField>
                    </div>
                    <template #footer>
                        <div class="flex justify-end">
                            <UButton label="Ubah password" icon="i-tabler-shield-lock" :loading="savingPassword"
                                @click="savePassword" />
                        </div>
                    </template>
                </UCard>
            </div>

            <UCard class="h-fit">
                <template #header>
                    <h2 class="font-bold text-slate-900 dark:text-white">Ringkasan akun</h2>
                </template>
                <dl class="space-y-4 text-sm">
                    <div>
                        <dt class="text-slate-500">Peran</dt>
                        <dd class="font-semibold text-slate-900 dark:text-white mt-1">{{ roleLabels[user?.role] ||
                            user?.role || '-'
                            }}</dd>
                    </div>
                    <div>
                        <dt class="text-slate-500">Akun dibuat</dt>
                        <dd class="font-semibold text-slate-900 dark:text-white mt-1">{{ formatDate(user?.created_at) }}
                        </dd>
                    </div>
                    <div>
                        <dt class="text-slate-500">Login terakhir</dt>
                        <dd class="font-semibold text-slate-900 dark:text-white mt-1">{{ formatDate(user?.last_login_at)
                            }}</dd>
                    </div>
                </dl>
            </UCard>
        </div>
    </div>
</template>