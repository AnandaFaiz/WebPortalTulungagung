<script setup>
defineProps({
    label: { type: String, default: 'Simpan' },
    size: { type: String, default: 'md' },
    block: { type: Boolean, default: false },
    loading: { type: Boolean, default: false }
})
</script>

<template>
    <UButton icon="i-tabler-check" :label="label" :size="size" :block="block" :loading="loading" color="primary"
        variant="solid" />
</template>