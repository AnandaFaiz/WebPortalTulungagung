<script setup>
definePageMeta({ layout: 'admin', title: 'Informasi MPP' })
const toast = useToast()
const API_BASE = 'http://localhost:8080/api'

const loading = ref(true)
const saving = ref(false)

const form = ref({
    deskripsi: '',
    alamat: '',
    jam_operasional: '',
    jumlah_layanan: 0,
    jumlah_instansi: 0
})

const fetchData = async () => {
    loading.value = true
    try {
        const res = await $fetch(`${API_BASE}/admin/mpp`, { credentials: 'include' })
        if (res.mpp) form.value = res.mpp
    } catch (err) {
        toast.add({ title: 'Gagal memuat', color: 'error' })
    } finally {
        loading.value = false
    }
}

const saveData = async () => {
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/mpp`, {
            method: 'PUT',
            body: form.value,
            credentials: 'include'
        })
        toast.add({ title: 'Berhasil', description: 'Informasi MPP diperbarui', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', color: 'error' })
    } finally {
        saving.value = false
    }
}

onMounted(() => fetchData())
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Informasi MPP" description="Kelola informasi Mal Pelayanan Publik.">
            <template #actions>
                <AdminButtonSave label="Simpan Perubahan" :loading="saving" @click="saveData" />
            </template>
        </AdminPageHeader>

        <div v-if="loading" class="py-12 text-center">
            <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <!-- Form Kiri -->
            <div class="lg:col-span-2 space-y-6">
                <UCard>
                    <template #header>
                        <h3 class="font-bold flex items-center gap-2">
                            <UIcon name="i-tabler-building-store" class="text-sky-500" /> Detail MPP
                        </h3>
                    </template>
                    <div class="space-y-4">
                        <UFormField label="Deskripsi Singkat">
                            <UTextarea v-model="form.deskripsi" :rows="4" size="lg" class="w-full" />
                        </UFormField>
                        <UFormField label="Alamat Lengkap">
                            <UInput v-model="form.alamat" icon="i-tabler-map-pin" size="lg" class="w-full" />
                        </UFormField>
                        <UFormField label="Jam Operasional">
                            <UInput v-model="form.jam_operasional" icon="i-tabler-clock" size="lg" class="w-full" />
                        </UFormField>
                        <div class="grid grid-cols-2 gap-4">
                            <UFormField label="Jumlah Jenis Layanan">
                                <UInput v-model.number="form.jumlah_layanan" type="number" size="lg" class="w-full" />
                            </UFormField>
                            <UFormField label="Jumlah Instansi">
                                <UInput v-model.number="form.jumlah_instansi" type="number" size="lg" class="w-full" />
                            </UFormField>
                        </div>
                    </div>
                </UCard>
            </div>

            <!-- Preview Kanan -->
            <div class="lg:col-span-1">
                <UCard class="sticky top-24">
                    <template #header>
                        <h3 class="font-bold flex items-center gap-2">
                            <UIcon name="i-tabler-eye" class="text-sky-500" /> Preview
                        </h3>
                    </template>
                    <div class="space-y-4">
                        <p class="text-sm text-slate-500">{{ form.deskripsi || 'Deskripsi belum diisi...' }}</p>
                        <div class="flex items-start gap-2 text-sm">
                            <UIcon name="i-tabler-map-pin" class="text-sky-500 mt-0.5 shrink-0" />
                            <span>{{ form.alamat || 'Alamat belum diisi' }}</span>
                        </div>
                        <div class="flex items-start gap-2 text-sm">
                            <UIcon name="i-tabler-clock" class="text-sky-500 mt-0.5 shrink-0" />
                            <span>{{ form.jam_operasional || 'Jam operasional belum diisi' }}</span>
                        </div>
                        <div class="grid grid-cols-2 gap-2 mt-4">
                            <div class="bg-slate-50 p-3 rounded-lg text-center">
                                <p class="text-xl font-bold text-sky-600">{{ form.jumlah_layanan }}</p>
                                <p class="text-[10px] text-slate-500 uppercase font-bold">Layanan</p>
                            </div>
                            <div class="bg-slate-50 p-3 rounded-lg text-center">
                                <p class="text-xl font-bold text-sky-600">{{ form.jumlah_instansi }}</p>
                                <p class="text-[10px] text-slate-500 uppercase font-bold">Instansi</p>
                            </div>
                        </div>
                    </div>
                </UCard>
            </div>
        </div>
    </div>
</template>