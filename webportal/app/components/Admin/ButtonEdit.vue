<script setup>
defineProps({
    label: { type: String, default: '' },
    size: { type: String, default: 'sm' },
    block: { type: Boolean, default: false }
})
</script>

<template>
    <UButton icon="i-tabler-edit" :label="label || undefined" :size="size" :block="block" color="warning"
        variant="outline"
        class="bg-amber-50! hover:bg-amber-100! text-amber-700! border-amber-300! hover:border-amber-400!" />
</template>