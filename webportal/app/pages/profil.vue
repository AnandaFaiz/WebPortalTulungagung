<script setup>
import profilData from '~/data/profil.json'

const profil = profilData
</script>

<template>
    <div>
        <UPageHero title="Profil Kabupaten Tulungagung"
            description="Mengenal lebih dekat sejarah, visi misi, dan jajaran pemerintahan Kabupaten Tulungagung." />

        <!-- Section Sejarah & Lambang -->
        <UPageSection title="Sejarah & Lambang" description="Jejak panjang perjalanan Kabupaten Tulungagung.">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <UCard>
                    <h3 class="font-bold text-lg mb-3">Sejarah Singkat</h3>
                    <p class="text-gray-600 dark:text-gray-400 leading-relaxed">{{ profil.sejarah }}</p>
                </UCard>
                <UCard>
                    <h3 class="font-bold text-lg mb-3">Makna Lambang</h3>
                    <p class="text-gray-600 dark:text-gray-400 leading-relaxed">{{ profil.lambang }}</p>
                </UCard>
            </div>
        </UPageSection>

        <!-- Section Visi & Misi -->
        <UPageSection title="Visi & Misi" description="Arah pembangunan Kabupaten Tulungagung.">
            <UCard class="mb-6">
                <h3 class="font-bold text-lg mb-3">Visi</h3>
                <p class="text-xl font-medium text-primary">{{ profil.visi }}</p>
            </UCard>

            <h3 class="font-bold text-lg mb-4">Misi</h3>
            <div class="space-y-3">
                <UCard v-for="(misi, index) in profil.misi" :key="index" variant="subtle">
                    <div class="flex items-start gap-4">
                        <div
                            class="w-8 h-8 rounded-full bg-primary text-white flex items-center justify-center shrink-0 font-bold">
                            {{ index + 1 }}
                        </div>
                        <p>{{ misi }}</p>
                    </div>
                </UCard>
            </div>
        </UPageSection>

        <!-- Section Profil Pejabat -->
        <UPageSection title="Pejabat Daerah"
            description="Bupati, Wakil Bupati, dan jajaran Pemerintah Kabupaten Tulungagung.">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <UCard v-for="item in profil.pejabat" :key="item.id" class="text-center">
                    <div
                        class="w-40 h-40 mx-auto rounded-full bg-gray-200 dark:bg-gray-800 flex items-center justify-center mb-4">
                        <UIcon name="i-tabler-user" class="w-20 h-20 text-gray-400" />
                    </div>
                    <h3 class="font-bold text-lg">{{ item.nama }}</h3>
                    <p class="text-primary">{{ item.jabatan }}</p>
                </UCard>
            </div>
        </UPageSection>
    </div>
</template>