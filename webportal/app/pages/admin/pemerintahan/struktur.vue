<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Struktur Organisasi'
})

const toast = useToast()

const loading = ref(true)
const saving = ref(false)

const form = ref({
    judul: '',
    deskripsi: '',
    gambar: ''
})

// ==== Fetch ====
const fetchStruktur = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/struktur',{
            credentials: 'include'
        })
        const s = response.struktur || {}
        form.value = {
            judul: s.judul || 'Struktur Organisasi Pemkab Tulungagung',
            deskripsi: s.deskripsi || '',
            gambar: s.gambar || ''
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
const saveStruktur = async () => {
    saving.value = true
    try {
        await $fetch('http://localhost:8080/api/admin/struktur', {
            method: 'PUT',
            body: {
                judul: form.value.judul,
                deskripsi: form.value.deskripsi,
                gambar: form.value.gambar
            },
            credentials: 'include'
        })

        toast.add({
            title: 'Berhasil disimpan',
            description: 'Struktur Organisasi berhasil diperbarui',
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
    fetchStruktur()
    toast.add({
        title: 'Form di-reset',
        description: 'Perubahan yang belum disimpan telah dibatalkan',
        icon: 'i-tabler-refresh',
        color: 'neutral',
        duration: 2000
    })
}

onMounted(() => {
    fetchStruktur()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Struktur Organisasi" description="Kelola bagan struktur organisasi Pemkab Tulungagung.">
            <template #actions>
                <div class="flex items-center gap-2">
                    <UButton icon="i-tabler-refresh" label="Reset" color="neutral" variant="outline"
                        @click="resetForm" />
                    <AdminButtonSave label="Simpan Perubahan" :loading="saving" @click="saveStruktur" />
                </div>
            </template>
        </AdminPageHeader>

        <!-- Loading -->
        <div v-if="loading" class="text-center py-12">
            <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            <p class="text-sm text-slate-500 mt-3">Memuat data...</p>
        </div>

        <template v-else>
            <!-- Info Card -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-info-circle" class="text-sky-500" />
                        Informasi Struktur
                    </h3>
                </template>

                <div class="space-y-4">
                    <UFormField label="Judul">
                        <UInput v-model="form.judul" placeholder="Struktur Organisasi Pemkab Tulungagung"
                            icon="i-tabler-heading" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Deskripsi Singkat">
                        <UTextarea v-model="form.deskripsi"
                            placeholder="Deskripsi singkat tentang struktur organisasi..." :rows="3" size="lg"
                            class="w-full" />
                    </UFormField>
                </div>
            </UCard>

            <!-- Upload Gambar -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-photo" class="text-sky-500" />
                        Gambar Bagan
                    </h3>
                </template>

                <!-- Preview -->
                <div v-if="form.gambar" class="mb-4">
                    <div
                        class="relative rounded-lg overflow-hidden border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800">
                        <img :src="form.gambar" alt="Bagan Struktur Organisasi" class="w-full h-auto" />
                        <div class="absolute top-3 right-3 flex gap-2">
                            <UButton icon="i-tabler-refresh" color="neutral" variant="solid" size="sm" label="Ganti" />
                            <UButton icon="i-tabler-trash" color="error" variant="solid" size="sm" aria-label="Hapus"
                                @click="form.gambar = ''" />
                        </div>
                    </div>
                </div>

                <!-- Upload Area -->
                <div v-else
                    class="border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-10 text-center hover:border-sky-500 transition-colors cursor-pointer">
                    <UIcon name="i-tabler-cloud-upload" class="w-12 h-12 text-slate-400 mx-auto mb-3" />
                    <p class="text-sm text-slate-500 font-medium">Klik untuk upload gambar bagan</p>
                    <p class="text-xs text-slate-400 mt-1">PNG, JPG, WebP — Rekomendasi 2000×1200 (Max 5MB)</p>
                </div>

                <template #footer>
                    <div class="flex items-start gap-2 text-xs text-slate-400">
                        <UIcon name="i-tabler-info-circle" class="w-4 h-4 mt-0.5 shrink-0" />
                        <p>
                            Upload gambar bagan struktur organisasi. Setelah backend upload siap, gambar akan otomatis
                            disimpan di server.
                            Untuk sementara, kamu bisa paste URL gambar langsung ke field di bawah.
                        </p>
                    </div>
                </template>
            </UCard>

            <!-- Input URL Gambar (fallback manual) -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-link" class="text-sky-500" />
                        URL Gambar (Manual)
                    </h3>
                </template>

                <UFormField label="URL Gambar Bagan" hint="Kosongkan kalau belum ada gambar.">
                    <UInput v-model="form.gambar" type="url" placeholder="https://example.com/bagan.png"
                        icon="i-tabler-photo" size="lg" class="w-full" />
                </UFormField>
            </UCard>

            <!-- Preview Section -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-eye" class="text-sky-500" />
                        Preview di Landing Page
                    </h3>
                </template>

                <div class="p-4 bg-slate-50 dark:bg-slate-800/50 rounded-lg">
                    <h3 class="text-2xl font-bold text-slate-900 dark:text-white mb-3">
                        {{ form.judul || 'Struktur Organisasi' }}
                    </h3>
                    <p v-if="form.deskripsi" class="text-slate-500 mb-4">
                        {{ form.deskripsi }}
                    </p>
                    <div v-if="form.gambar" class="bg-white dark:bg-slate-800 rounded-lg p-2">
                        <img :src="form.gambar" alt="Preview" class="w-full h-auto rounded" />
                    </div>
                    <div v-else class="bg-white dark:bg-slate-800 rounded-lg p-8 text-center">
                        <UIcon name="i-tabler-photo-off" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                        <p class="text-sm text-slate-400">Belum ada gambar bagan</p>
                    </div>
                </div>
            </UCard>
        </template>
    </div>
</template>