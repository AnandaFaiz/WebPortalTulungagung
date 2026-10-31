<script setup lang="">
defineProps({
    icon: { type: String, default: 'i-tabler-inbox' },
    title: { type: String, default: 'Tidak ada data ditemukan' },
    description: { type: String, default: 'Coba ubah filter atau tambahkan data baru.' }
})
</script>

<template>
    <div class="text-center py-12">
        <UIcon :name="icon" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
        <p class="text-slate-500 font-medium">{{ title }}</p>
        <p v-if="description" class="text-xs text-slate-400 mt-1">{{ description }}</p>
        <slot />
    </div>
</template>