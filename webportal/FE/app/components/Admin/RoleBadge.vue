<script setup>
const props = defineProps({
    role: { type: String, required: true }
})

const roleMap = {
    super_admin: { label: 'Super Admin', color: 'primary', icon: 'i-tabler-shield-check' },
    editor: { label: 'Editor', color: 'success', icon: 'i-tabler-edit' }
}

const config = computed(() => roleMap[props.role] || {
    label: props.role,
    color: 'neutral',
    icon: 'i-tabler-user'
})
</script>

<template>
    <UBadge :color="config.color" variant="subtle" size="sm">
        <UIcon :name="config.icon" class="w-3.5 h-3.5 mr-1" />
        {{ config.label }}
    </UBadge>
</template>