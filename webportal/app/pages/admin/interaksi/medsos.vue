<script setup>
definePageMeta({ layout: 'admin', title: 'Media Sosial' })
const toast = useToast()
const loading = ref(true)
const saving = ref(false)
const form = ref({ alamat: '', telepon: '', email: '', jam_operasional: '', embed_maps: '', facebook: '', instagram: '', twitter: '', youtube: '' })
const fields = [
    { key: 'facebook', label: 'Facebook', icon: 'i-mdi-facebook', placeholder: 'https://facebook.com/...' },
    { key: 'instagram', label: 'Instagram', icon: 'i-mdi-instagram', placeholder: 'https://instagram.com/...' },
    { key: 'twitter', label: 'X / Twitter', icon: 'i-tabler-brand-x', placeholder: 'https://x.com/...' },
    { key: 'youtube', label: 'YouTube', icon: 'i-mdi-youtube', placeholder: 'https://youtube.com/@...' }
]

const fetchData = async () => {
    loading.value = true
    try {
        const response = await $fetch('http://localhost:8080/api/admin/kontak', { credentials: 'include' })
        const kontak = response.kontak || {}
        form.value = {
            alamat: kontak.alamat || '',
            telepon: kontak.telepon || '',
            email: kontak.email || '',
            jam_operasional: kontak.jam_operasional || '',
            embed_maps: kontak.embed_maps || '',
            ...fields.reduce((result, field) => ({ ...result, [field.key]: kontak[field.key] || '' }), {})
        }
    } catch (err) {
        toast.add({ title: 'Gagal memuat data', color: 'error' })
    } finally { loading.value = false }
}

const saveData = async () => {
    saving.value = true
    try {
        await $fetch('http://localhost:8080/api/admin/kontak', { method: 'PUT', body: form.value, credentials: 'include' })
        toast.add({ title: 'Media sosial disimpan', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', description: err?.data?.statusMessage || 'Terjadi kesalahan', color: 'error' })
    } finally { saving.value = false }
}

onMounted(fetchData)
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Media Sosial"
            description="Kelola tautan media sosial yang tampil di landing page dan halaman kontak.">
            <template #actions>
                <AdminButtonSave label="Simpan Perubahan" :loading="saving" @click="saveData" />
            </template>
        </AdminPageHeader>

        <div v-if="loading" class="py-12 text-center">
            <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
        </div>
        <UCard v-else>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <UFormField v-for="field in fields" :key="field.key" :label="field.label">
                    <UInput v-model="form[field.key]" :placeholder="field.placeholder" :icon="field.icon" size="lg"
                        class="w-full" />
                </UFormField>
            </div>
            <UAlert class="mt-6" color="info" variant="subtle" icon="i-tabler-info-circle"
                title="Tautan kosong tidak ditampilkan"
                description="Simpan URL lengkap dengan https:// agar tombol media sosial dapat dibuka dari landing page." />
        </UCard>
    </div>
</template>
