<script setup lang="ts">
import { getResourcePermission } from '~/utils/permissions'

const props = defineProps({
    label: { type: String, default: 'Simpan' },
    size: { type: String, default: 'md' },
    block: { type: Boolean, default: false },
    loading: { type: Boolean, default: false },
    action: { type: String, default: 'update' }
})
const route = useRoute()
const { can } = useAuth()
const allowed = computed(() => {
    const resource = getResourcePermission(route.path)
    return !resource || can(resource, props.action as 'create' | 'read' | 'update' | 'delete')
})
</script>

<template>
    <UButton v-if="allowed" icon="i-tabler-check" :label="label" :size="size" :block="block" :loading="loading"
        color="primary" variant="solid" />
</template>