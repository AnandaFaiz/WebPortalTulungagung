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
    return !resource || can(resource, 'update')
})
</script>

<template>
    <UButton v-if="allowed" icon="i-tabler-edit" :label="label || undefined" :size="size" :block="block" color="warning"
        variant="outline"
        class="bg-amber-50! hover:bg-amber-100! text-amber-700! border-amber-300! hover:border-amber-400!" />
</template>