<script setup>
import beritaData from '~/data/berita.json'
import potensiData from '~/data/potensi.json'
import pengumumanData from '~/data/pengumuman.json'
import aplikasiData from '~/data/aplikasi.json'

const layananCepat = [
    { title: 'Kependudukan', description: 'KTP, KK, Akta Kelahiran', icon: 'i-tabler-users', to: '/layanan/kependudukan' },
    { title: 'Perizinan', description: 'Izin Usaha, IMB, SIUP', icon: 'i-tabler-file-text', to: '/layanan/perizinan' },
    { title: 'Pengaduan', description: 'Laporkan keluhan & aspirasi', icon: 'i-tabler-speakerphone', to: '/layanan/pengaduan' },
    { title: 'Darurat', description: 'Ambulans, Damkar, Polisi', icon: 'i-tabler-phone-call', to: '/layanan/darurat' }
]

const berita = beritaData
const potensi = potensiData
const pengumuman = pengumumanData
const aplikasi = aplikasiData
</script>

<template>
    <div>
        <!-- Hero Section -->
        <UPageHero title="Selamat Datang di Portal Pemkab Tulungagung"
            description="Portal resmi Pemerintah Kabupaten Tulungagung. Dapatkan informasi terbaru, layanan publik, dan berita kegiatan pemerintah daerah dalam satu pintu."
            :links="[
                { label: 'Lihat Layanan', color: 'primary', size: 'xl', to: '/layanan' },
                { label: 'Baca Berita', color: 'neutral', variant: 'outline', size: 'xl', to: '/berita' }
            ]">
            <template #headline>
                <UBadge color="primary" variant="subtle" size="lg">
                    Portal Resmi Pemerintah Daerah
                </UBadge>
            </template>
        </UPageHero>

        <!-- Section Layanan Cepat -->
        <UPageSection title="Layanan Cepat"
            description="Akses layanan publik yang paling sering dibutuhkan masyarakat.">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <UPageCard v-for="(item, index) in layananCepat" :key="index" v-bind="item" variant="subtle"
                    class="hover:ring-2 hover:ring-primary transition-all cursor-pointer" />
            </div>
        </UPageSection>

        <!-- Section Sambutan Bupati -->
        <UPageSection>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
                <!-- Kolom Kiri: Foto Bupati -->
                <div class="flex justify-center">
                    <div
                        class="w-64 h-80 bg-gray-200 dark:bg-gray-800 rounded-2xl flex items-center justify-center shadow-lg">
                        <!-- Placeholder Foto -->
                        <UIcon name="i-tabler-user" class="w-24 h-24 text-gray-400" />
                    </div>
                </div>

                <!-- Kolom Kanan: Teks Sambutan -->
                <div>
                    <UBadge color="primary" variant="subtle" class="mb-4">Sambutan Bupati</UBadge>
                    <h2 class="text-3xl font-bold mb-4">Membangun Tulungagung yang Lebih Maju</h2>
                    <p class="text-gray-600 dark:text-gray-400 leading-relaxed mb-4">
                        Assalamu'alaikum Warahmatullahi Wabarakatuh. Puji syukur kita panjatkan ke hadirat Tuhan Yang
                        Maha Esa
                        atas segala rahmat dan karunia-Nya.
                    </p>
                    <p class="text-gray-600 dark:text-gray-400 leading-relaxed mb-6">
                        Selamat datang di Portal Resmi Pemerintah Kabupaten Tulungagung. Portal ini kami hadirkan
                        sebagai wujud
                        komitmen kami dalam memberikan pelayanan publik yang transparan, akuntabel, dan mudah diakses
                        oleh
                        seluruh masyarakat.
                    </p>
                    <div>
                        <p class="font-bold text-lg">H. Gatut Sunoto, S.Pd., M.M.</p>
                        <p class="text-sm text-gray-500">Bupati Tulungagung</p>
                    </div>
                </div>
            </div>
        </UPageSection>

        <!-- Section Potensi Daerah: Pariwisata -->
        <UPageSection title="Jelajahi Tulungagung"
            description="Destinasi wisata alam dan budaya yang menakjubkan di Kabupaten Tulungagung.">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <UCard v-for="item in potensi.pariwisata" :key="item.id" class="overflow-hidden">
                    <template #header>
                        <div
                            class="h-56 bg-linear-to-br from-emerald-400 to-emerald-600 rounded-lg flex items-center justify-center">
                            <UIcon name="i-tabler-beach" class="w-16 h-16 text-white/80" />
                        </div>
                    </template>
                    <h3 class="font-bold text-lg">{{ item.nama }}</h3>
                    <p class="text-sm text-gray-500 mt-2">{{ item.deskripsi }}</p>
                    <template #footer>
                        <div class="flex items-center gap-2 text-sm text-gray-400">
                            <UIcon name="i-tabler-map-pin" />
                            <span>{{ item.lokasi }}</span>
                        </div>
                    </template>
                </UCard>
            </div>
        </UPageSection>

        <!-- Section Potensi Daerah: UMKM -->
        <UPageSection title="Produk Unggulan Daerah" description="UMKM dan produk unggulan khas Tulungagung.">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <UCard v-for="item in potensi.umkm" :key="item.id">
                    <template #header>
                        <div
                            class="h-40 bg-linear-to-br from-amber-400 to-amber-600 rounded-lg flex items-center justify-center">
                            <UIcon name="i-tabler-shopping-bag" class="w-14 h-14 text-white/80" />
                        </div>
                    </template>
                    <UBadge color="warning" variant="subtle" size="sm" class="mb-2">
                        {{ item.kategori }}
                    </UBadge>
                    <h3 class="font-bold text-lg">{{ item.nama }}</h3>
                    <p class="text-sm text-gray-500 mt-2">{{ item.deskripsi }}</p>
                </UCard>
            </div>
        </UPageSection>

        <!-- Section Investasi (CTA) -->
        <UPageCTA title="Investasi di Tulungagung"
            description="Peluang investasi dengan potensi besar dan kemudahan perizinan." :links="[
                { label: 'Lihat Peluang Investasi', color: 'primary', size: 'lg', to: '/potensi/investasi' },
                { label: 'Hubungi Kami', color: 'neutral', variant: 'outline', size: 'lg', to: '/kontak' }
            ]" />

        <!-- Section Berita & Pengumuman -->
        <UPageSection title="Informasi Terkini" description="Berita dan pengumuman resmi dari Pemkab Tulungagung.">
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                <!-- Kolom Kiri: Berita (2/3 lebar) -->
                <div class="lg:col-span-2">
                    <h3 class="font-semibold text-lg mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-news" class="text-primary" />
                        Berita Terbaru
                    </h3>
                    <div class="space-y-4">
                        <UCard v-for="item in berita" :key="item.id">
                            <div class="flex gap-4">
                                <div class="w-24 h-24 bg-gray-200 dark:bg-gray-800 rounded-lg shrink-0"></div>
                                <div>
                                    <UBadge color="primary" variant="subtle" size="sm" class="mb-1">
                                        {{ item.kategori }}
                                    </UBadge>
                                    <h4 class="font-bold">{{ item.judul }}</h4>
                                    <p class="text-sm text-gray-500 mt-1 line-clamp-2">{{ item.deskripsi }}</p>
                                    <p class="text-xs text-gray-400 mt-2">{{ item.tanggal }}</p>
                                </div>
                            </div>
                        </UCard>
                    </div>
                </div>

                <!-- Kolom Kanan: Pengumuman (1/3 lebar) -->
                <div>
                    <h3 class="font-semibold text-lg mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-bell" class="text-primary" />
                        Pengumuman
                    </h3>
                    <div class="space-y-3">
                        <UCard v-for="item in pengumuman" :key="item.id" variant="subtle">
                            <p class="text-sm font-medium">{{ item.judul }}</p>
                            <p class="text-xs text-gray-400 mt-2 flex items-center gap-1">
                                <UIcon name="i-mdi-calendar-blank" />
                                {{ item.tanggal }}
                            </p>
                        </UCard>
                    </div>
                </div>
            </div>
        </UPageSection>

        <!-- Section Aplikasi Daerah -->
        <UPageSection title="Aplikasi Daerah"
            description="Akses layanan digital Pemkab Tulungagung melalui aplikasi berikut.">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <UCard v-for="item in aplikasi" :key="item.id"
                    class="text-center hover:ring-2 hover:ring-primary transition-all cursor-pointer">
                    <div class="flex flex-col items-center gap-3 py-4">
                        <div class="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center">
                            <UIcon :name="item.icon" class="w-8 h-8 text-primary" />
                        </div>
                        <h3 class="font-bold">{{ item.nama }}</h3>
                        <p class="text-sm text-gray-500">{{ item.deskripsi }}</p>
                    </div>
                </UCard>
            </div>
        </UPageSection>
    </div>
</template>