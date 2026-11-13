<script setup>
const props = defineProps({
    open: { type: Boolean, default: false },
    title: { type: String, required: true },
    description: { type: String, default: '' },
    isEdit: { type: Boolean, default: false },
    loading: { type: Boolean, default: false },
    submitLabel: { type: String, default: '' },
    cancelLabel: { type: String, default: 'Batal' },
    modalClass: { type: String, default: '' }
})

const emit = defineEmits(['update:open', 'submit', 'cancel'])

const finalSubmitLabel = computed(() => {
    if (props.submitLabel) return props.submitLabel
    return props.isEdit ? 'Simpan Perubahan' : 'Simpan'
})

const handleSubmit = () => {
    emit('submit')
}

const handleCancel = () => {
    emit('update:open', false)
    emit('cancel')
}
</script>

<template>
    <UModal :open="open" :title="title" :description="description" :class="modalClass"
        @update:open="emit('update:open', $event)">
        <template #body>
            <slot />
        </template>

        <template #footer>
            <div class="flex items-center justify-end gap-2">
                <UButton :label="cancelLabel" color="neutral" variant="outline" @click="handleCancel" />
                <AdminButtonSave :label="finalSubmitLabel" :action="isEdit ? 'update' : 'create'" :loading="loading"
                    @click="handleSubmit" />
            </div>
        </template>
    </UModal>
</template>