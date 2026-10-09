<script setup>
import pemerintahanData from '~/data/pemerintahan.json'

const data = pemerintahanData
</script>

<template>
    <div>
        <UPageHero title="Pemerintahan Kabupaten Tulungagung"
            description="Informasi OPD, PPID, JDIH, dan DPRD Kabupaten Tulungagung." />

        <!-- Section OPD -->
        <UPageSection title="Perangkat Daerah (OPD)"
            description="Direktori Dinas, Badan, dan Sekretariat di lingkungan Pemkab Tulungagung.">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <UCard v-for="item in data.opd" :key="item.id"
                    class="hover:ring-2 hover:ring-primary transition-all cursor-pointer">
                    <div class="flex items-start gap-3">
                        <div class="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
                            <UIcon name="i-tabler-buildings" class="w-5 h-5 text-primary" />
                        </div>
                        <div>
                            <h3 class="font-bold text-sm">{{ item.singkatan }}</h3>
                            <p class="text-xs text-gray-500">{{ item.nama }}</p>
                        </div>
                    </div>
                </UCard>
            </div>
        </UPageSection>

        <!-- Section PPID -->
        <UPageSection title="PPID" description="Pejabat Pengelola Informasi dan Dokumentasi.">
            <UCard>
                <p class="text-gray-600 dark:text-gray-400 leading-relaxed mb-6">
                    {{ data.ppid.deskripsi }}
                </p>
                <h3 class="font-bold mb-3">Layanan PPID:</h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div v-for="(layanan, i) in data.ppid.layanan" :key="i"
                        class="flex items-center gap-3 p-3 rounded-lg bg-gray-50 dark:bg-gray-800">
                        <UIcon name="i-tabler-file-text" class="text-primary shrink-0" />
                        <span class="text-sm">{{ layanan }}</span>
                    </div>
                </div>
            </UCard>
        </UPageSection>

        <!-- Section JDIH -->
        <UPageSection title="JDIH" description="Jaringan Dokumentasi dan Informasi Hukum Kabupaten Tulungagung.">
            <UCard>
                <div class="space-y-3">
                    <div v-for="item in data.jdih" :key="item.id"
                        class="flex flex-col md:flex-row md:items-center justify-between p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary transition-colors">
                        <div>
                            <UBadge color="primary" variant="subtle" size="sm" class="mb-2">{{ item.jenis }}</UBadge>
                            <p class="font-bold">{{ item.nomor }}</p>
                            <p class="text-sm text-gray-500">{{ item.tentang }}</p>
                        </div>
                        <div class="text-sm text-gray-400 mt-2 md:mt-0">{{ item.tanggal }}</div>
                    </div>
                </div>
            </UCard>
        </UPageSection>

        <!-- Section DPRD -->
        <UPageSection title="DPRD Kabupaten Tulungagung" description="Dewan Perwakilan Rakyat Daerah.">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
                <UCard class="text-center">
                    <p class="text-4xl font-bold text-primary">{{ data.dprd.jumlah_anggota }}</p>
                    <p class="text-sm text-gray-500 mt-2">Jumlah Anggota</p>
                </UCard>
                <UCard class="text-center">
                    <p class="text-2xl font-bold text-primary">{{ data.dprd.periode }}</p>
                    <p class="text-sm text-gray-500 mt-2">Periode</p>
                </UCard>
                <UCard class="text-center">
                    <p class="text-4xl font-bold text-primary">{{ data.dprd.komisi.length }}</p>
                    <p class="text-sm text-gray-500 mt-2">Komisi</p>
                </UCard>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UCard v-for="(k, i) in data.dprd.komisi" :key="i" variant="subtle">
                    <h3 class="font-bold">{{ k.nama }}</h3>
                    <p class="text-sm text-gray-500 mt-1">{{ k.bidang }}</p>
                </UCard>
            </div>
        </UPageSection>
    </div>
</template>