<script setup>
const props = defineProps({
    total: { type: Number, required: true },
    itemsPerPage: { type: Number, default: 10 }
})
const page = defineModel('page', { type: Number, default: 1 })
</script>

<template>
    <div class="flex items-center justify-between flex-wrap gap-3">
        <p class="text-xs text-slate-400">
            Menampilkan {{ total }} data (Halaman {{ page }})
        </p>
        <UPagination v-if="total > itemsPerPage" v-model:page="page" :total="total" :items-per-page="itemsPerPage" />
    </div>
</template>