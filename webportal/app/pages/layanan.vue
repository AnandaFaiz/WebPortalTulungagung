<script setup>
import layananData from '~/data/layanan.json'

const data = layananData
</script>

<template>
    <div>
        <UPageHero title="Layanan Publik"
            description="Akses berbagai layanan publik Pemkab Tulungagung dalam satu halaman." />

        <!-- Section Kategori Layanan -->
        <UPageSection title="Kategori Layanan" description="Pilih kategori layanan yang Anda butuhkan.">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <UCard v-for="item in data.kategori" :key="item.id"
                    class="hover:ring-2 hover:ring-primary transition-all cursor-pointer">
                    <div class="flex flex-col items-center text-center gap-3 py-4">
                        <div class="w-16 h-16 rounded-2xl bg-primary/10 flex items-center justify-center">
                            <UIcon :name="item.icon" class="w-8 h-8 text-primary" />
                        </div>
                        <h3 class="font-bold text-lg">{{ item.nama }}</h3>
                        <p class="text-sm text-gray-500">{{ item.deskripsi }}</p>
                        <UButton label="Selengkapnya" variant="link" trailing-icon="i-tabler-arrow-right" />
                    </div>
                </UCard>
            </div>
        </UPageSection>

        <!-- Section MPP -->
        <UPageSection title="Mal Pelayanan Publik (MPP)"
            description="Layanan terpadu dari berbagai instansi dalam satu gedung.">
            <UCard>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <div>
                        <p class="text-gray-600 dark:text-gray-400 leading-relaxed mb-6">
                            {{ data.mpp.deskripsi }}
                        </p>
                        <div class="space-y-3 text-sm">
                            <div class="flex items-start gap-3">
                                <UIcon name="i-tabler-map-pin" class="text-primary mt-0.5 shrink-0" />
                                <span>{{ data.mpp.alamat }}</span>
                            </div>
                            <div class="flex items-start gap-3">
                                <UIcon name="i-tabler-clock" class="text-primary mt-0.5 shrink-0" />
                                <span>{{ data.mpp.jam_operasional }}</span>
                            </div>
                        </div>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <UCard variant="subtle" class="text-center">
                            <p class="text-3xl font-bold text-primary">{{ data.mpp.jumlah_layanan }}+</p>
                            <p class="text-xs text-gray-500 mt-1">Jenis Layanan</p>
                        </UCard>
                        <UCard variant="subtle" class="text-center">
                            <p class="text-3xl font-bold text-primary">{{ data.mpp.instansi }}</p>
                            <p class="text-xs text-gray-500 mt-1">Instansi</p>
                        </UCard>
                    </div>
                </div>
            </UCard>
        </UPageSection>

        <!-- Section Pengaduan -->
        <UPageSection title="Pengaduan Masyarakat"
            description="Sampaikan keluhan, saran, atau aspirasi Anda kepada Pemkab Tulungagung.">
            <UCard>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                        <h3 class="font-bold mb-3">Kanal Pengaduan</h3>
                        <div class="space-y-3">
                            <div class="flex items-center gap-3 p-3 rounded-lg bg-gray-50 dark:bg-gray-800">
                                <UIcon name="i-tabler-world" class="text-primary shrink-0" />
                                <div>
                                    <p class="font-medium text-sm">SP4N LAPOR</p>
                                    <p class="text-xs text-gray-500">www.lapor.go.id</p>
                                </div>
                            </div>
                            <div class="flex items-center gap-3 p-3 rounded-lg bg-gray-50 dark:bg-gray-800">
                                <UIcon name="i-tabler-phone" class="text-primary shrink-0" />
                                <div>
                                    <p class="font-medium text-sm">Telepon</p>
                                    <p class="text-xs text-gray-500">(0355) 321-260</p>
                                </div>
                            </div>
                            <div class="flex items-center gap-3 p-3 rounded-lg bg-gray-50 dark:bg-gray-800">
                                <UIcon name="i-tabler-mail" class="text-primary shrink-0" />
                                <div>
                                    <p class="font-medium text-sm">Email</p>
                                    <p class="text-xs text-gray-500">pengaduan@tulungagung.go.id</p>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div>
                        <UButton label="Buat Pengaduan Baru" icon="i-tabler-edit" color="primary" size="lg"
                            block />
                        <p class="text-xs text-gray-500 mt-3 text-center">
                            Pengaduan Anda akan ditindaklanjuti dalam 5 hari kerja.
                        </p>
                    </div>
                </div>
            </UCard>
        </UPageSection>

        <!-- Section Darurat -->
        <UPageSection title="Layanan Darurat" description="Nomor penting yang dapat dihubungi dalam situasi darurat.">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                <UCard v-for="(item, i) in data.darurat" :key="i" class="hover:ring-2 hover:ring-error transition-all">
                    <div class="flex items-center gap-4">
                        <div class="w-12 h-12 rounded-full bg-error/10 flex items-center justify-center shrink-0">
                            <UIcon :name="item.icon" class="w-6 h-6 text-error" />
                        </div>
                        <div>
                            <p class="font-bold text-lg">{{ item.nomor }}</p>
                            <p class="text-sm text-gray-500">{{ item.nama }}</p>
                        </div>
                    </div>
                </UCard>
            </div>
        </UPageSection>
    </div>
</template>