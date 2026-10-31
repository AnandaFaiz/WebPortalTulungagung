<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Profil Daerah'
})

const toast = useToast()

// ==== State ====
const loading = ref(true)
const saving = ref(false)

// ==== Form Data ====
const form = ref({
    sejarah: '',
    lambang: '',
    visi: '',
    misi: [],
    geografis: '',
    demografi: {
        luas: '',
        penduduk: '',
        kecamatan: '',
        desa: ''
    }
})

// ==== Fetch ====
const fetchProfil = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/profil-daerah', {
            credentials: 'include'
        })

        // Pastikan misi adalah array
        const profil = response.profil || {}
        form.value = {
            sejarah: profil.sejarah || '',
            lambang: profil.lambang || '',
            visi: profil.visi || '',
            misi: Array.isArray(profil.misi) ? [...profil.misi] : [],
            geografis: profil.geografis || '',
            demografi: profil.demografi && typeof profil.demografi === 'object'
                ? { ...profil.demografi }
                : { luas: '', penduduk: '', kecamatan: '', desa: '' }
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

// ==== Misi Handlers ====
const addMisi = () => {
    form.value.misi.push('')
}

const removeMisi = (index) => {
    form.value.misi.splice(index, 1)
}

// ==== Save ====
const saveProfil = async () => {
    saving.value = true
    try {
        await $fetch('http://localhost:8080/api/admin/profil-daerah', {
            method: 'PUT',
            body: {
                sejarah: form.value.sejarah,
                lambang: form.value.lambang,
                visi: form.value.visi,
                misi: form.value.misi.filter(m => m && m.trim()),
                geografis: form.value.geografis,
                demografi: form.value.demografi
            },
            credentials: 'include'
        })

        toast.add({
            title: 'Berhasil disimpan',
            description: 'Profil Daerah berhasil diperbarui',
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
    fetchProfil()
    toast.add({
        title: 'Form di-reset',
        description: 'Perubahan yang belum disimpan telah dibatalkan',
        icon: 'i-tabler-refresh',
        color: 'neutral',
        duration: 2000
    })
}

onMounted(() => {
    fetchProfil()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Profil Daerah"
            description="Kelola konten profil Kabupaten Tulungagung: sejarah, lambang, visi misi, dan data wilayah.">
            <template #actions>
                <div class="flex items-center gap-2">
                    <UButton icon="i-tabler-refresh" label="Reset" color="neutral" variant="outline"
                        @click="resetForm" />
                    <AdminButtonSave label="Simpan Perubahan" :loading="saving" @click="saveProfil" />
                </div>
            </template>
        </AdminPageHeader>

        <!-- Loading -->
        <div v-if="loading" class="text-center py-12">
            <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            <p class="text-sm text-slate-500 mt-3">Memuat data profil...</p>
        </div>

        <template v-else>

            <!-- Section 1: Sejarah & Lambang -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-history" class="text-sky-500" />
                        Sejarah & Lambang
                    </h3>
                </template>

                <div class="space-y-4">
                    <UFormField label="Sejarah Kabupaten Tulungagung">
                        <UTextarea v-model="form.sejarah" placeholder="Tulis sejarah singkat Kabupaten Tulungagung..."
                            :rows="6" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Makna Lambang">
                        <UTextarea v-model="form.lambang"
                            placeholder="Tulis makna filosofis lambang Kabupaten Tulungagung..." :rows="5" size="lg"
                            class="w-full" />
                    </UFormField>
                </div>
            </UCard>

            <!-- Section 2: Visi & Misi -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-target" class="text-sky-500" />
                        Visi & Misi
                    </h3>
                </template>

                <div class="space-y-4">
                    <UFormField label="Visi">
                        <UTextarea v-model="form.visi"
                            placeholder="Contoh: Terwujudnya Tulungagung yang Maju, Mandiri, dan Sejahtera..." :rows="3"
                            size="lg" class="w-full" />
                    </UFormField>

                    <div>
                        <div class="flex items-center justify-between mb-3">
                            <label class="text-sm font-medium text-slate-700 dark:text-slate-300">
                                Misi
                            </label>
                            <UButton icon="i-tabler-plus" label="Tambah Misi" color="primary" variant="soft" size="sm"
                                @click="addMisi" />
                        </div>

                        <!-- Empty state -->
                        <div v-if="form.misi.length === 0"
                            class="text-center py-8 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-lg">
                            <UIcon name="i-tabler-list-check" class="w-10 h-10 text-slate-300 mx-auto mb-2" />
                            <p class="text-sm text-slate-400">Belum ada misi</p>
                            <UButton label="Tambah Misi Pertama" color="primary" variant="link" size="sm"
                                @click="addMisi" />
                        </div>

                        <!-- List Misi -->
                        <div v-else class="space-y-2">
                            <div v-for="(misi, index) in form.misi" :key="index" class="flex items-start gap-2">
                                <div class="w-8 h-10 flex items-center justify-center shrink-0">
                                    <div
                                        class="w-7 h-7 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 dark:text-sky-400 font-bold text-xs">
                                        {{ index + 1 }}
                                    </div>
                                </div>
                                <UTextarea v-model="form.misi[index]" :placeholder="`Misi ke-${index + 1}`" :rows="2"
                                    size="lg" class="flex-1" />
                                <UButton icon="i-tabler-trash" color="error" variant="ghost" size="md" class="mt-1"
                                    aria-label="Hapus misi" @click="removeMisi(index)" />
                            </div>
                        </div>
                    </div>
                </div>
            </UCard>

            <!-- Section 3: Geografis -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-map-2" class="text-sky-500" />
                        Geografis
                    </h3>
                </template>

                <UFormField label="Deskripsi Geografis">
                    <UTextarea v-model="form.geografis"
                        placeholder="Deskripsi letak, batas wilayah, dan topografi Kabupaten Tulungagung..." :rows="5"
                        size="lg" class="w-full" />
                </UFormField>
            </UCard>

            <!-- Section 4: Demografi -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-users" class="text-sky-500" />
                        Demografi
                    </h3>
                </template>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <UFormField label="Luas Wilayah">
                        <UInput v-model="form.demografi.luas" placeholder="Contoh: 1.055,65 km²"
                            icon="i-tabler-ruler-measure" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Jumlah Penduduk">
                        <UInput v-model="form.demografi.penduduk" placeholder="Contoh: 1.100.000 jiwa"
                            icon="i-tabler-users" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Jumlah Kecamatan">
                        <UInput v-model="form.demografi.kecamatan" placeholder="Contoh: 19 kecamatan"
                            icon="i-tabler-map-pin" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Jumlah Desa/Kelurahan">
                        <UInput v-model="form.demografi.desa" placeholder="Contoh: 257 desa, 14 kelurahan"
                            icon="i-tabler-home" size="lg" class="w-full" />
                    </UFormField>
                </div>
            </UCard>

            <!-- Preview Section -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-eye" class="text-sky-500" />
                        Preview di Landing Page
                    </h3>
                </template>

                <div class="space-y-6 p-4 bg-slate-50 dark:bg-slate-800/50 rounded-lg">
                    <!-- Preview Visi -->
                    <div v-if="form.visi">
                        <p class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Visi</p>
                        <p class="text-lg font-semibold text-sky-600">{{ form.visi }}</p>
                    </div>

                    <!-- Preview Misi -->
                    <div v-if="form.misi.filter(m => m).length > 0">
                        <p class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Misi</p>
                        <div class="space-y-2">
                            <div v-for="(misi, index) in form.misi.filter(m => m)" :key="index"
                                class="flex items-start gap-3">
                                <div
                                    class="w-7 h-7 rounded-full bg-sky-500 text-white flex items-center justify-center shrink-0 font-bold text-xs">
                                    {{ index + 1 }}
                                </div>
                                <p class="text-sm text-slate-700 dark:text-slate-300 pt-1">{{ misi }}</p>
                            </div>
                        </div>
                    </div>

                    <!-- Preview Demografi -->
                    <div
                        v-if="form.demografi.luas || form.demografi.penduduk || form.demografi.kecamatan || form.demografi.desa">
                        <p class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">Demografi</p>
                        <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                            <div v-if="form.demografi.luas"
                                class="bg-white dark:bg-slate-800 p-3 rounded-lg text-center">
                                <p class="text-lg font-bold text-sky-600">{{ form.demografi.luas }}</p>
                                <p class="text-xs text-slate-500">Luas Wilayah</p>
                            </div>
                            <div v-if="form.demografi.penduduk"
                                class="bg-white dark:bg-slate-800 p-3 rounded-lg text-center">
                                <p class="text-lg font-bold text-sky-600">{{ form.demografi.penduduk }}</p>
                                <p class="text-xs text-slate-500">Penduduk</p>
                            </div>
                            <div v-if="form.demografi.kecamatan"
                                class="bg-white dark:bg-slate-800 p-3 rounded-lg text-center">
                                <p class="text-lg font-bold text-sky-600">{{ form.demografi.kecamatan }}</p>
                                <p class="text-xs text-slate-500">Kecamatan</p>
                            </div>
                            <div v-if="form.demografi.desa"
                                class="bg-white dark:bg-slate-800 p-3 rounded-lg text-center">
                                <p class="text-lg font-bold text-sky-600">{{ form.demografi.desa }}</p>
                                <p class="text-xs text-slate-500">Desa/Kelurahan</p>
                            </div>
                        </div>
                    </div>
                </div>
            </UCard>

        </template>
    </div>
</template>