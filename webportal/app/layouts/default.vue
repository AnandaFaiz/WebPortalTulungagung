<script setup>
const items = [
    { label: 'Beranda', to: '/' },
    { label: 'Profil', to: '/profil' },
    { label: 'Pemerintahan', to: '/pemerintahan' },
    { label: 'Layanan Publik', to: '/layanan' },
    { label: 'Informasi & Berita', to: '/berita' },
    { label: 'Kontak', to: '/kontak' },
]
</script>

<template>
    <div>
        <!-- Header Navigasi -->
        <UHeader title="Pemkab Tulungagung" :links="items">
            <template #right>
                <UButton icon="i-tabler-search" color="neutral" variant="ghost" />
                <UColorModeButton />
            </template>
        </UHeader>

        <!-- Konten Utama -->
        <UMain>
            <slot />
        </UMain>

        <!-- Footer -->
        <UFooter>
            <template #left>
                <p class="text-sm text-gray-500">
                    © {{ new Date().getFullYear() }} Pemerintah Kabupaten Tulungagung
                </p>
            </template>
            <template #right>
                <div class="flex items-center gap-2">
                    <UButton icon="i-mdi-facebook" color="neutral" variant="ghost" to="#" />
                    <UButton icon="i-mdi-instagram" color="neutral" variant="ghost" to="#" />
                    <UButton icon="i-mdi-twitter" color="neutral" variant="ghost" to="#" />
                    <UButton icon="i-tabler-brand-youtube" color="neutral" variant="ghost" to="#" />
                </div>
            </template>
        </UFooter>
    </div>
</template>