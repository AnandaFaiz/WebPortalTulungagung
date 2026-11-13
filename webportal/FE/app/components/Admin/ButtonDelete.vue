<script setup lang="ts">
import { getResourcePermission } from '~/utils/permissions'

const route = useRoute()
const { can } = useAuth()
defineProps({
    label: { type: String, default: '' },
    size: { type: String, default: 'sm' },
    block: { type: Boolean, default: false }
})

const allowed = computed(() => {
    const resource = getResourcePermission(route.path)
    return !resource || can(resource, 'delete')
})
</script>

<template>
    <UButton v-if="allowed" icon="i-tabler-trash" :label="label || undefined" :size="size" :block="block" color="error"
        variant="solid" class="bg-red-500! hover:bg-red-600! text-white! border! border-red-600!" />
</template>