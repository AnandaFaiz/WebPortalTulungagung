<script setup>
const kontak = [
    {
        icon: 'i-tabler-map-pin',
        label: 'Alamat Kantor',
        nilai: 'Jl. A. Yani Timur No. 37, Tulungagung, Jawa Timur 66216',
        link: null
    },
    {
        icon: 'i-tabler-phone-call',
        label: 'Nomor Pengaduan (24 Jam)',
        nilai: '(0355) 321-260',
        link: 'tel:0355321260'
    },
    {
        icon: 'i-tabler-mail',
        label: 'Email',
        nilai: 'info@tulungagung.go.id',
        link: 'mailto:info@tulungagung.go.id'
    },
    {
        icon: 'i-tabler-clock',
        label: 'Jam Operasional',
        nilai: 'Senin - Jumat, 08.00 - 15.00 WIB',
        link: null
    }
]
</script>

<template>
    <div>
        <UPageHero title="Hubungi Kami"
            description="Sampaikan pertanyaan, saran, atau aspirasi Anda kepada Pemkab Tulungagung." />

        <!-- Section Kontak Utama -->
        <UPageSection title="Informasi Kontak" description="Berbagai kanal komunikasi resmi Pemkab Tulungagung.">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <UCard v-for="(item, i) in kontak" :key="i" class="hover:ring-2 hover:ring-primary transition-all">
                    <div class="flex items-start gap-4">
                        <div
                            class="w-12 h-12 rounded-xl bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                            <UIcon :name="item.icon" class="w-6 h-6 text-sky-500" />
                        </div>
                        <div class="flex-1">
                            <p class="text-xs uppercase text-gray-400 tracking-wide font-semibold">{{ item.label }}</p>
                            <component :is="item.link ? 'a' : 'p'" :href="item.link"
                                class="mt-1 font-medium block hover:text-primary transition-colors">
                                {{ item.nilai }}
                            </component>
                        </div>
                    </div>
                </UCard>
            </div>
        </UPageSection>

        <!-- Section Peta -->
        <UPageSection title="Lokasi Kantor" description="Kantor Bupati Tulungagung.">
            <UCard>
                <div class="h-96 bg-slate-200 dark:bg-slate-800 rounded-lg flex items-center justify-center">
                    <div class="text-center">
                        <UIcon name="i-tabler-map-pin" class="w-16 h-16 text-sky-500 mx-auto" />
                        <p class="mt-4 text-gray-500">Peta Lokasi Kantor Bupati Tulungagung</p>
                        <p class="text-xs text-gray-400 mt-1">(Embed Google Maps akan ditambahkan di sini)</p>
                    </div>
                </div>
            </UCard>
        </UPageSection>

        <!-- Section CTA Pengaduan -->
        <UPageCTA title="Punya Keluhan atau Aspirasi?"
            description="Sampaikan melalui kanal pengaduan resmi Pemkab Tulungagung." :links="[
                { label: 'Buat Pengaduan', color: 'primary', size: 'lg', to: '/layanan' },
                { label: 'Lihat FAQ', color: 'neutral', variant: 'outline', size: 'lg', to: '#' }
            ]" />
    </div>
</template>