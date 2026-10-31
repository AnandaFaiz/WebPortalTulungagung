<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Identitas Situs'
})

const toast = useToast()

const loading = ref(true)
const saving = ref(false)

const form = ref({
    nama_situs: '',
    tagline: '',
    deskripsi: '',
    logo: '/Logo-Kabupaten-Tulungagung.png',
    favicon: '',
    email: '',
    telepon: '',
    alamat: ''
})

// ==== Fetch ====
const fetchIdentitas = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/identitas', {
            credentials: 'include'
        })
        const s = response.identitas || {}
        form.value = {
            nama_situs: s.nama_situs || 'Portal Pemkab Tulungagung',
            tagline: s.tagline || '',
            deskripsi: s.deskripsi || '',
            logo: s.logo || '/Logo-Kabupaten-Tulungagung.png',
            favicon: s.favicon || '',
            email: s.email || '',
            telepon: s.telepon || '',
            alamat: s.alamat || ''
        }
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

// ==== Save ====
const saveIdentitas = async () => {
    if (!form.value.nama_situs.trim()) {
        toast.add({
            title: 'Validasi gagal',
            description: 'Nama situs wajib diisi',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
        return
    }

    saving.value = true
    try {
        const requestFetch = useRequestFetch()
        await requestFetch('http://localhost:8080/api/admin/identitas', {
            method: 'PUT',
            body: { ...form.value },
            credentials: 'include'
        })

        toast.add({
            title: 'Berhasil disimpan',
            description: 'Identitas Situs berhasil diperbarui',
            icon: 'i-tabler-circle-check',
            color: 'success',
            duration: 2500
        })
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

// ==== Reset ====
const resetForm = () => {
    fetchIdentitas()
    toast.add({
        title: 'Form di-reset',
        description: 'Perubahan yang belum disimpan telah dibatalkan',
        icon: 'i-tabler-refresh',
        color: 'neutral',
        duration: 2000
    })
}

onMounted(() => {
    fetchIdentitas()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Identitas Situs"
            description="Kelola identitas dasar portal: nama, tagline, logo, dan kontak utama.">
            <template #actions>
                <div class="flex items-center gap-2">
                    <UButton icon="i-tabler-refresh" label="Reset" color="neutral" variant="outline"
                        @click="resetForm" />
                    <AdminButtonSave label="Simpan Perubahan" :loading="saving" @click="saveIdentitas" />
                </div>
            </template>
        </AdminPageHeader>

        <!-- Loading -->
        <div v-if="loading" class="text-center py-12">
            <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            <p class="text-sm text-slate-500 mt-3">Memuat data...</p>
        </div>

        <template v-else>
            <!-- Section 1: Identitas Dasar -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-id-badge-2" class="text-sky-500" />
                        Identitas Dasar
                    </h3>
                </template>

                <div class="space-y-4">
                    <UFormField label="Nama Situs" required>
                        <UInput v-model="form.nama_situs" placeholder="Portal Pemkab Tulungagung"
                            icon="i-tabler-building" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Tagline" hint="Kalimat singkat yang muncul di bawah nama situs.">
                        <UInput v-model="form.tagline" placeholder="Portal Resmi Pemerintah Daerah"
                            icon="i-tabler-quote" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Deskripsi Situs" hint="Deskripsi ini akan muncul di meta description untuk SEO.">
                        <UTextarea v-model="form.deskripsi"
                            placeholder="Portal resmi Pemerintah Kabupaten Tulungagung yang menyediakan..." :rows="3"
                            size="lg" class="w-full" />
                    </UFormField>
                </div>
            </UCard>

            <!-- Section 2: Logo & Favicon -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-photo" class="text-sky-500" />
                        Logo & Favicon
                    </h3>
                </template>

                <div class="space-y-8">

                    <!-- LOGO UTAMA — Layout Horizontal -->
                    <div>
                        <label class="block text-sm font-semibold text-slate-700 dark:text-slate-300 mb-3">
                            Logo Utama
                            <span class="text-xs font-normal text-slate-400 ml-2">PNG, transparan, max 500KB</span>
                        </label>

                        <div class="flex items-start gap-4">
                            <!-- Preview Logo -->
                            <div
                                class="w-24 h-24 rounded-xl bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center justify-center overflow-hidden shrink-0">
                                <img v-if="form.logo" :src="form.logo" alt="Logo"
                                    class="w-full h-full object-contain p-2" />
                                <UIcon v-else name="i-tabler-photo" class="w-8 h-8 text-slate-400" />
                            </div>

                            <!-- Upload + Input -->
                            <div class="flex-1 space-y-2">
                                <div
                                    class="border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-4 text-center hover:border-sky-500 hover:bg-sky-50/50 dark:hover:bg-sky-950/20 transition-colors cursor-pointer">
                                    <div class="flex items-center justify-center gap-2 text-slate-500">
                                        <UIcon name="i-tabler-cloud-upload" class="w-5 h-5" />
                                        <span class="text-sm font-medium">Klik untuk upload logo</span>
                                    </div>
                                    <p class="text-xs text-slate-400 mt-1">Atau paste URL di bawah</p>
                                </div>
                                <UInput v-model="form.logo" placeholder="/Logo-Kabupaten-Tulungagung.png"
                                    icon="i-tabler-link" size="md" class="w-full" />
                            </div>
                        </div>
                    </div>

                    <!-- Divider -->
                    <div class="border-t border-slate-200 dark:border-slate-700"></div>

                    <!-- FAVICON — Layout Horizontal -->
                    <div>
                        <label class="block text-sm font-semibold text-slate-700 dark:text-slate-300 mb-3">
                            Favicon
                            <span class="text-xs font-normal text-slate-400 ml-2">ICO/PNG 32×32px</span>
                        </label>

                        <div class="flex items-start gap-4">
                            <!-- Preview Favicon -->
                            <div
                                class="w-16 h-16 rounded-lg bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center justify-center overflow-hidden shrink-0">
                                <img v-if="form.favicon" :src="form.favicon" alt="Favicon"
                                    class="w-full h-full object-contain p-1.5" />
                                <UIcon v-else name="i-tabler-favicon" class="w-5 h-5 text-slate-400" />
                            </div>

                            <!-- Upload + Input -->
                            <div class="flex-1 space-y-2">
                                <div
                                    class="border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-4 text-center hover:border-sky-500 hover:bg-sky-50/50 dark:hover:bg-sky-950/20 transition-colors cursor-pointer">
                                    <div class="flex items-center justify-center gap-2 text-slate-500">
                                        <UIcon name="i-tabler-cloud-upload" class="w-5 h-5" />
                                        <span class="text-sm font-medium">Klik untuk upload favicon</span>
                                    </div>
                                    <p class="text-xs text-slate-400 mt-1">Atau paste URL di bawah</p>
                                </div>
                                <UInput v-model="form.favicon" placeholder="/favicon.ico" icon="i-tabler-link" size="md"
                                    class="w-full" />
                            </div>
                        </div>
                    </div>

                </div>
            </UCard>

            <!-- Section 3: Kontak Utama -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-address-book" class="text-sky-500" />
                        Kontak Utama
                    </h3>
                </template>

                <div class="space-y-4">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <UFormField label="Email Resmi">
                            <UInput v-model="form.email" type="email" placeholder="info@tulungagung.go.id"
                                icon="i-tabler-mail" size="lg" class="w-full" />
                        </UFormField>

                        <UFormField label="Telepon">
                            <UInput v-model="form.telepon" placeholder="(0355) 321-260" icon="i-tabler-phone" size="lg"
                                class="w-full" />
                        </UFormField>
                    </div>

                    <UFormField label="Alamat Kantor">
                        <UTextarea v-model="form.alamat"
                            placeholder="Jl. A. Yani Timur No. 37, Tulungagung, Jawa Timur 66216" :rows="2" size="lg"
                            class="w-full" />
                    </UFormField>
                </div>
            </UCard>

            <!-- Preview Section -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-eye" class="text-sky-500" />
                        Preview Identitas
                    </h3>
                </template>

                <div class="p-6 bg-slate-50 dark:bg-slate-800/50 rounded-lg">
                    <div class="flex items-center gap-3 mb-4">
                        <img v-if="form.logo" :src="form.logo" alt="Logo" class="h-12 w-auto" />
                        <div class="leading-tight">
                            <p class="text-xs text-slate-500">{{ form.tagline }}</p>
                            <p class="font-extrabold text-slate-900 dark:text-white">{{ form.nama_situs }}</p>
                        </div>
                    </div>
                    <p v-if="form.deskripsi" class="text-sm text-slate-600 dark:text-slate-400 leading-relaxed">
                        {{ form.deskripsi }}
                    </p>
                    <div
                        class="mt-4 pt-4 border-t border-slate-200 dark:border-slate-700 grid grid-cols-1 md:grid-cols-3 gap-3 text-xs text-slate-500">
                        <div v-if="form.email" class="flex items-center gap-1.5">
                            <UIcon name="i-tabler-mail" class="w-3.5 h-3.5" />
                            {{ form.email }}
                        </div>
                        <div v-if="form.telepon" class="flex items-center gap-1.5">
                            <UIcon name="i-tabler-phone" class="w-3.5 h-3.5" />
                            {{ form.telepon }}
                        </div>
                        <div v-if="form.alamat" class="flex items-center gap-1.5 md:col-span-3">
                            <UIcon name="i-tabler-map-pin" class="w-3.5 h-3.5 shrink-0" />
                            {{ form.alamat }}
                        </div>
                    </div>
                </div>
            </UCard>
        </template>

    </div>
</template>