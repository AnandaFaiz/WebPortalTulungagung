<script setup>
const props = defineProps({
    status: { type: String, required: true }
})

// Mapping warna otomatis berdasarkan nilai status
const colorMap = {
    // Konten
    'Publish': 'success',
    'Draft': 'warning',
    'Archived': 'neutral',

    // User aktif
    'Aktif': 'success',
    'Nonaktif': 'error',

    // Pengaduan
    'Baru': 'error',
    'Diproses': 'warning',
    'Selesai': 'success',

    // Prioritas pengumuman
    'Tinggi': 'error',
    'Normal': 'neutral'
}

const color = computed(() => colorMap[props.status] || 'neutral')
</script>

<template>
    <UBadge :color="color" variant="subtle" size="sm">
        {{ status }}
    </UBadge>
</template>