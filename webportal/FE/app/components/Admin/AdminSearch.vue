<script setup>
const props = defineProps({
    placeholder: { type: String, default: 'Cari data...' }
})
const model = defineModel({ type: String, default: '' })
</script>

<template>
    <UInput v-model="model" icon="i-tabler-search" :placeholder="placeholder" size="lg" class="w-full" />
</template>