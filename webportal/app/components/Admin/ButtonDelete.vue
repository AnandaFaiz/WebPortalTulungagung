<script setup>
defineProps({
    label: { type: String, default: '' },
    size: { type: String, default: 'sm' },
    block: { type: Boolean, default: false }
})
</script>

<template>
    <UButton icon="i-tabler-trash" :label="label || undefined" :size="size" :block="block" color="error" variant="solid"
        class="bg-red-500! hover:bg-red-600! text-white! border! border-red-600!" />
</template>