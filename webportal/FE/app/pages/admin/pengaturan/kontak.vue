<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Kontak'
})

const toast = useToast()

const loading = ref(true)
const saving = ref(false)

const form = ref({
    alamat: '',
    telepon: '',
    email: '',
    jam_operasional: '',
    embed_maps: '',
    facebook: '',
    instagram: '',
    twitter: '',
    youtube: ''
})

// ==== Fetch ====
const fetchKontak = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/kontak',{
            credentials: 'include'
        })
        const k = response.kontak || {}
        form.value = {
            alamat: k.alamat || '',
            telepon: k.telepon || '',
            email: k.email || '',
            jam_operasional: k.jam_operasional || '',
            embed_maps: k.embed_maps || '',
            facebook: k.facebook || '',
            instagram: k.instagram || '',
            twitter: k.twitter || '',
            youtube: k.youtube || ''
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
const saveKontak = async () => {
    saving.value = true
    try {
        await $fetch('http://localhost:8080/api/admin/kontak', {
            method: 'PUT',
            body: { ...form.value },
            credentials: 'include'
        })

        toast.add({
            title: 'Berhasil disimpan',
            description: 'Data Kontak berhasil diperbarui',
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
    fetchKontak()
    toast.add({
        title: 'Form di-reset',
        description: 'Perubahan yang belum disimpan telah dibatalkan',
        icon: 'i-tabler-refresh',
        color: 'neutral',
        duration: 2000
    })
}

// ==== Helper: Preview embed maps ====
const mapsPreviewSrc = computed(() => {
    const raw = form.value.embed_maps || ''
    if (!raw) return ''
    // Kalau user paste full iframe HTML, extract src="..."
    const match = raw.match(/src=["']([^"']+)["']/)
    if (match) return match[1]
    // Kalau sudah URL langsung, pakai apa adanya
    return raw
})

onMounted(() => {
    fetchKontak()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <AdminPageHeader title="Kontak"
            description="Kelola informasi kontak, lokasi kantor, dan media sosial yang tampil di halaman Kontak.">
            <template #actions>
                <div class="flex items-center gap-2">
                    <UButton icon="i-tabler-refresh" label="Reset" color="neutral" variant="outline"
                        @click="resetForm" />
                    <AdminButtonSave label="Simpan Perubahan" :loading="saving" @click="saveKontak" />
                </div>
            </template>
        </AdminPageHeader>

        <!-- Loading -->
        <div v-if="loading" class="text-center py-12">
            <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
            <p class="text-sm text-slate-500 mt-3">Memuat data...</p>
        </div>

        <template v-else>

            <!-- Section 1: Informasi Kontak -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-address-book" class="text-sky-500" />
                        Informasi Kontak
                    </h3>
                </template>

                <div class="space-y-4">
                    <UFormField label="Alamat Kantor">
                        <UTextarea v-model="form.alamat"
                            placeholder="Jl. A. Yani Timur No. 37, Tulungagung, Jawa Timur 66216" :rows="2" size="lg"
                            class="w-full" />
                    </UFormField>

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <UFormField label="Telepon">
                            <UInput v-model="form.telepon" placeholder="(0355) 321-260" icon="i-tabler-phone" size="lg"
                                class="w-full" />
                        </UFormField>

                        <UFormField label="Email">
                            <UInput v-model="form.email" type="email" placeholder="info@tulungagung.go.id"
                                icon="i-tabler-mail" size="lg" class="w-full" />
                        </UFormField>
                    </div>

                    <UFormField label="Jam Operasional">
                        <UInput v-model="form.jam_operasional" placeholder="Senin - Jumat, 08.00 - 15.00 WIB"
                            icon="i-tabler-clock" size="lg" class="w-full" />
                    </UFormField>
                </div>
            </UCard>

            <!-- Section 2: Lokasi / Embed Maps -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-map-pin" class="text-sky-500" />
                        Lokasi Kantor (Google Maps)
                    </h3>
                </template>

                <div class="space-y-4">
                    <UAlert color="info" variant="subtle" icon="i-tabler-info-circle"
                        title="Cara Mendapatkan Embed Maps"
                        description="Buka Google Maps → cari lokasi Kantor Bupati → Share → Embed a map → Copy HTML → paste di kolom di bawah." />

                    <UFormField label="Embed Google Maps">
                        <UTextarea v-model="form.embed_maps"
                            placeholder='Paste HTML embed Google Maps atau langsung URL-nya...' :rows="4" size="lg"
                            class="w-full font-mono text-xs" />
                    </UFormField>

                    <!-- Preview Maps -->
                    <div v-if="mapsPreviewSrc">
                        <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
                            Preview Peta
                        </label>
                        <div class="rounded-lg overflow-hidden border border-slate-200 dark:border-slate-700">
                            <iframe :src="mapsPreviewSrc" width="100%" height="300" style="border:0" allowfullscreen
                                loading="lazy" referrerpolicy="no-referrer-when-downgrade" />
                        </div>
                    </div>
                </div>
            </UCard>

            <!-- Section 3: Media Sosial -->
            <UCard>
                <template #header>
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-share-3" class="text-sky-500" />
                        Media Sosial
                    </h3>
                </template>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <UFormField label="Facebook">
                        <UInput v-model="form.facebook" placeholder="https://facebook.com/..." icon="i-mdi-facebook"
                            size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Instagram">
                        <UInput v-model="form.instagram" placeholder="https://instagram.com/..." icon="i-mdi-instagram"
                            size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Twitter / X">
                        <UInput v-model="form.twitter" placeholder="https://twitter.com/..." icon="i-mdi-twitter"
                            size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="YouTube">
                        <UInput v-model="form.youtube" placeholder="https://youtube.com/@..." icon="i-mdi-youtube"
                            size="lg" class="w-full" />
                    </UFormField>
                </div>

                <template #footer>
                    <p class="text-xs text-slate-400 flex items-center gap-1">
                        <UIcon name="i-tabler-info-circle" class="w-3.5 h-3.5 shrink-0" />
                        Kosongkan jika tidak ada akun di platform tersebut.
                    </p>
                </template>
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

                    <!-- Preview Info Kontak -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div class="flex items-start gap-3 p-3 bg-white dark:bg-slate-800 rounded-lg">
                            <div
                                class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                <UIcon name="i-tabler-map-pin" class="w-5 h-5 text-sky-600 dark:text-sky-400" />
                            </div>
                            <div class="flex-1 min-w-0">
                                <p class="text-xs uppercase text-slate-400 tracking-wide font-semibold">Alamat</p>
                                <p class="text-sm text-slate-700 dark:text-slate-300 mt-0.5">{{ form.alamat || '-' }}
                                </p>
                            </div>
                        </div>

                        <div class="flex items-start gap-3 p-3 bg-white dark:bg-slate-800 rounded-lg">
                            <div
                                class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                <UIcon name="i-tabler-phone" class="w-5 h-5 text-sky-600 dark:text-sky-400" />
                            </div>
                            <div class="flex-1 min-w-0">
                                <p class="text-xs uppercase text-slate-400 tracking-wide font-semibold">Telepon</p>
                                <p class="text-sm text-slate-700 dark:text-slate-300 mt-0.5">{{ form.telepon || '-' }}
                                </p>
                            </div>
                        </div>

                        <div class="flex items-start gap-3 p-3 bg-white dark:bg-slate-800 rounded-lg">
                            <div
                                class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                <UIcon name="i-tabler-mail" class="w-5 h-5 text-sky-600 dark:text-sky-400" />
                            </div>
                            <div class="flex-1 min-w-0">
                                <p class="text-xs uppercase text-slate-400 tracking-wide font-semibold">Email</p>
                                <p class="text-sm text-slate-700 dark:text-slate-300 mt-0.5">{{ form.email || '-' }}</p>
                            </div>
                        </div>

                        <div class="flex items-start gap-3 p-3 bg-white dark:bg-slate-800 rounded-lg">
                            <div
                                class="w-10 h-10 rounded-lg bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                <UIcon name="i-tabler-clock" class="w-5 h-5 text-sky-600 dark:text-sky-400" />
                            </div>
                            <div class="flex-1 min-w-0">
                                <p class="text-xs uppercase text-slate-400 tracking-wide font-semibold">Jam Operasional
                                </p>
                                <p class="text-sm text-slate-700 dark:text-slate-300 mt-0.5">{{ form.jam_operasional ||
                                    '-' }}</p>
                            </div>
                        </div>
                    </div>

                    <!-- Preview Maps -->
                    <div v-if="mapsPreviewSrc" class="bg-white dark:bg-slate-800 rounded-lg p-2">
                        <iframe :src="mapsPreviewSrc" width="100%" height="200" style="border:0; border-radius: 6px;"
                            loading="lazy" />
                    </div>

                    <!-- Preview Medsos -->
                    <div class="flex items-center gap-3 flex-wrap">
                        <span class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Medsos:</span>
                        <a v-if="form.facebook" :href="form.facebook" target="_blank"
                            class="w-9 h-9 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 hover:bg-sky-500 hover:text-white transition-colors">
                            <UIcon name="i-mdi-facebook" class="w-4 h-4" />
                        </a>
                        <a v-if="form.instagram" :href="form.instagram" target="_blank"
                            class="w-9 h-9 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 hover:bg-sky-500 hover:text-white transition-colors">
                            <UIcon name="i-mdi-instagram" class="w-4 h-4" />
                        </a>
                        <a v-if="form.twitter" :href="form.twitter" target="_blank"
                            class="w-9 h-9 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 hover:bg-sky-500 hover:text-white transition-colors">
                            <UIcon name="i-mdi-twitter" class="w-4 h-4" />
                        </a>
                        <a v-if="form.youtube" :href="form.youtube" target="_blank"
                            class="w-9 h-9 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 hover:bg-sky-500 hover:text-white transition-colors">
                            <UIcon name="i-mdi-youtube" class="w-4 h-4" />
                        </a>
                        <span v-if="!form.facebook && !form.instagram && !form.twitter && !form.youtube"
                            class="text-xs text-slate-400 italic">Belum ada media sosial</span>
                    </div>

                </div>
            </UCard>

        </template>
    </div>
</template>