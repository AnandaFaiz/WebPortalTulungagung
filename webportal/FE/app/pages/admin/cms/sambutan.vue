<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Sambutan Bupati'
})

const toast = useToast()
const saving = ref(false)
const loading = ref(true)

// ==== State (single record) ====
const sambutan = ref({
    judul: 'Membangun Tulungagung yang Lebih Maju',
    nama_bupati: 'H. Gatut Sunu Wibowo, S.E., M.E.',
    jabatan: 'Bupati Tulungagung',
    foto: '',
    isi: ''
})

// ==== Fetch data ====
const fetchSambutan = async () => {
    loading.value = true
    try {
        // Placeholder: baca dari mock dulu
        // Nanti kalau ada backend: useRequestFetch('/api/admin/sambutan')
        await new Promise(r => setTimeout(r, 300))
        sambutan.value = {
            judul: 'Membangun Tulungagung yang Lebih Maju',
            nama_bupati: 'H. Gatut Sunu Wibowo, S.E., M.E.',
            jabatan: 'Bupati Tulungagung',
            foto: '',
            isi: 'Assalamu\'alaikum Warahmatullahi Wabarakatuh. Puji syukur kita panjatkan ke hadirat Tuhan Yang Maha Esa atas segala rahmat dan karunia-Nya.\n\nSelamat datang di Portal Resmi Pemerintah Kabupaten Tulungagung. Portal ini kami hadirkan sebagai wujud komitmen kami dalam memberikan pelayanan publik yang transparan, akuntabel, dan mudah diakses oleh seluruh masyarakat.\n\nSemoga portal ini bermanfaat dan menjadi jembatan komunikasi antara pemerintah dan masyarakat. Terima kasih.'
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
const saveSambutan = async () => {
    if (!sambutan.value.judul.trim() || !sambutan.value.nama_bupati.trim() || !sambutan.value.isi.trim()) {
        toast.add({
            title: 'Validasi gagal',
            description: 'Judul, nama bupati, dan isi sambutan wajib diisi',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
        return
    }

    saving.value = true
    try {
        // Placeholder: nanti POST ke backend
        await new Promise(r => setTimeout(r, 800))

        toast.add({
            title: 'Berhasil disimpan',
            description: 'Sambutan Bupati berhasil diperbarui',
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
    fetchSambutan()
    toast.add({
        title: 'Form di-reset',
        description: 'Perubahan yang belum disimpan telah dibatalkan',
        icon: 'i-tabler-refresh',
        color: 'neutral',
        duration: 2000
    })
}

onMounted(() => {
    fetchSambutan()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Sambutan Bupati"
            description="Kelola konten sambutan Bupati yang tampil di halaman Beranda.">
            <template #actions>
                <div class="flex items-center gap-2">
                    <UButton icon="i-tabler-refresh" label="Reset" color="neutral" variant="outline"
                        @click="resetForm" />
                    <AdminButtonSave label="Simpan Perubahan" :loading="saving" @click="saveSambutan" />
                </div>
            </template>
        </AdminPageHeader>

        <!-- Loading -->
        <div v-if="loading" class="text-center py-12">
            <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            <p class="text-sm text-slate-500 mt-3">Memuat data...</p>
        </div>

        <!-- Form -->
        <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">

            <!-- Kolom Kiri: Preview Foto -->
            <div class="lg:col-span-1">
                <UCard>
                    <template #header>
                        <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <UIcon name="i-tabler-photo" class="text-sky-500" />
                            Foto Bupati
                        </h3>
                    </template>

                    <div class="space-y-4">
                        <AdminImageUpload v-model="sambutan.foto" label="foto bupati" :aspect-ratio="3 / 4"
                            preview-class="aspect-3/4" />
                    </div>
                </UCard>
            </div>

            <!-- Kolom Kanan: Form Konten -->
            <div class="lg:col-span-2 space-y-6">
                <UCard>
                    <template #header>
                        <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <UIcon name="i-tabler-edit" class="text-sky-500" />
                            Informasi Sambutan
                        </h3>
                    </template>

                    <div class="space-y-4">
                        <UFormField label="Judul Sambutan" required>
                            <UInput v-model="sambutan.judul" placeholder="Contoh: Membangun Tulungagung yang Lebih Maju"
                                size="lg" class="w-full" />
                        </UFormField>

                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <UFormField label="Nama Bupati" required>
                                <UInput v-model="sambutan.nama_bupati"
                                    placeholder="Contoh: H. Gatut Sunu Wibowo, S.E., M.E." size="lg" class="w-full" />
                            </UFormField>

                            <UFormField label="Jabatan">
                                <UInput v-model="sambutan.jabatan" placeholder="Contoh: Bupati Tulungagung" size="lg"
                                    class="w-full" />
                            </UFormField>
                        </div>
                    </div>
                </UCard>

                <UCard>
                    <template #header>
                        <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                            <UIcon name="i-tabler-file-text" class="text-sky-500" />
                            Isi Sambutan
                        </h3>
                    </template>

                    <UFormField label="Teks Sambutan" required
                        hint="Gunakan enter dua kali untuk membuat paragraf baru.">
                        <UTextarea v-model="sambutan.isi" placeholder="Tulis isi sambutan Bupati di sini..." :rows="14"
                            size="lg" class="w-full" />
                    </UFormField>
                </UCard>
            </div>
        </div>

        <!-- Preview Section -->
        <UCard v-if="!loading">
            <template #header>
                <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                    <UIcon name="i-tabler-eye" class="text-sky-500" />
                    Preview Tampilan di Beranda
                </h3>
            </template>

            <div class="bg-slate-50 dark:bg-slate-800/50 rounded-lg p-6 md:p-8">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-8 items-center max-w-4xl mx-auto">
                    <!-- Preview Foto -->
                    <div class="flex justify-center">
                        <div
                            class="w-56 h-72 rounded-2xl overflow-hidden shadow-lg bg-white dark:bg-slate-700 ring-1 ring-slate-200 dark:ring-slate-700">
                            <img v-if="sambutan.foto" :src="sambutan.foto" alt="Bupati"
                                class="w-full h-full object-cover" />
                            <div v-else class="w-full h-full flex items-center justify-center">
                                <UIcon name="i-tabler-user" class="w-20 h-20 text-slate-300" />
                            </div>
                        </div>
                    </div>

                    <!-- Preview Teks -->
                    <div>
                        <UBadge color="primary" variant="subtle" class="mb-4">Sambutan Bupati</UBadge>
                        <h2 class="text-2xl font-bold text-slate-900 dark:text-white mb-3">{{ sambutan.judul }}</h2>
                        <p class="text-slate-600 dark:text-slate-400 text-sm leading-relaxed whitespace-pre-line mb-4">
                            {{
                                sambutan.isi }}</p>
                        <div>
                            <p class="font-bold text-slate-900 dark:text-white">{{ sambutan.nama_bupati }}</p>
                            <p class="text-sm text-slate-500">{{ sambutan.jabatan }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </UCard>

    </div>
</template>