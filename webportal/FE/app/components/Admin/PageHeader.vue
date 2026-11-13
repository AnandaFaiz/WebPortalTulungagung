<script setup>
import { getResourcePermission } from '~/utils/permissions'

const route = useRoute()
const { can } = useAuth()

defineProps({
    title: { type: String, required: true },
    description: { type: String, default: '' },
    actionLabel: { type: String, default: '' },
    actionIcon: { type: String, default: 'i-tabler-plus' }
})

const resource = computed(() => getResourcePermission(route.path))
const canCreate = computed(() => !resource.value || can(resource.value, 'create'))

const emit = defineEmits(['action'])
</script>

<template>
    <UCard>
        <div class="flex items-center justify-between flex-wrap gap-4">
            <div>
                <h2 class="text-2xl font-bold text-slate-900 dark:text-white">{{ title }}</h2>
                <p v-if="description" class="text-slate-500 mt-1">{{ description }}</p>
            </div>
            <slot name="actions">
                <UButton v-if="actionLabel && canCreate" :icon="actionIcon" :label="actionLabel" color="primary"
                    size="lg" @click="emit('action')" />
            </slot>
        </div>
    </UCard>
</template>