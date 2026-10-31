<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Pengaturan'
})

const settings = [
    { title: 'Identitas Situs', description: 'Nama portal, tagline, deskripsi SEO, logo, dan favicon.', icon: 'i-tabler-id-badge-2', to: '/admin/pengaturan/identitas', color: 'text-sky-500' },
    { title: 'Kontak', description: 'Alamat kantor, telepon, jam operasional, peta, dan media sosial.', icon: 'i-tabler-address-book', to: '/admin/pengaturan/kontak', color: 'text-emerald-500' },
    { title: 'Menu Navigasi', description: 'Atur tautan dan urutan menu yang tampil di portal publik.', icon: 'i-tabler-menu-2', to: '/admin/pengaturan/menu', color: 'text-amber-500' }
]
</script>

<template>
    <div class="space-y-6 max-w-5xl">
        <AdminPageHeader title="Pengaturan"
            description="Kelola konfigurasi utama yang membentuk pengalaman portal publik." />

        <UAlert color="info" variant="subtle" icon="i-tabler-info-circle" title="Pusat pengaturan portal"
            description="Pilih kategori untuk memperbarui informasi yang tampil kepada masyarakat. Perubahan disimpan langsung ke portal." />

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <NuxtLink v-for="item in settings" :key="item.to" :to="item.to" class="block group">
                <UCard
                    class="h-full transition-shadow group-hover:shadow-md group-hover:ring-2 group-hover:ring-sky-200 dark:group-hover:ring-sky-800">
                    <div class="flex items-start gap-4">
                        <div
                            class="w-11 h-11 rounded-xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center shrink-0">
                            <UIcon :name="item.icon" class="w-6 h-6" :class="item.color" />
                        </div>
                        <div class="min-w-0 flex-1">
                            <div class="flex items-center justify-between gap-3">
                                <h2 class="font-bold text-slate-900 dark:text-white">{{ item.title }}</h2>
                                <UIcon name="i-tabler-arrow-up-right" class="w-5 h-5 text-slate-400 shrink-0" />
                            </div>
                            <p class="text-sm text-slate-500 mt-2 leading-relaxed">{{ item.description }}</p>
                        </div>
                    </div>
                </UCard>
            </NuxtLink>
        </div>

        <UCard>
            <template #header>
                <h2 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                    <UIcon name="i-tabler-adjustments" class="text-sky-500" />
                    Rekomendasi pengelolaan
                </h2>
            </template>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm text-slate-600 dark:text-slate-300">
                <div class="flex gap-3">
                    <UIcon name="i-tabler-check" class="text-emerald-500 shrink-0" /><span>Gunakan email resmi yang
                        aktif untuk
                        menerima pertanyaan warga.</span>
                </div>
                <div class="flex gap-3">
                    <UIcon name="i-tabler-check" class="text-emerald-500 shrink-0" /><span>Pastikan logo dan favicon
                        menggunakan
                        aset resmi dengan ukuran ringan.</span>
                </div>
                <div class="flex gap-3">
                    <UIcon name="i-tabler-check" class="text-emerald-500 shrink-0" /><span>Periksa tautan menu setelah
                        mengubah URL
                        halaman publik.</span>
                </div>
            </div>
        </UCard>
    </div>
</template>