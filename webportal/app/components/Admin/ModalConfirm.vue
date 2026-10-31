<script setup>
const props = defineProps({
    open: { type: Boolean, default: false },
    title: { type: String, default: 'Konfirmasi Hapus' },
    description: { type: String, default: 'Tindakan ini tidak dapat dibatalkan.' },
    itemName: { type: String, default: '' },
    confirmLabel: { type: String, default: 'Hapus' },
    cancelLabel: { type: String, default: 'Batal' },
    loading: { type: Boolean, default: false }
})

const emit = defineEmits(['update:open', 'confirm', 'cancel'])

const handleCancel = () => {
    emit('update:open', false)
    emit('cancel')
}

const handleConfirm = () => {
    emit('confirm')
}
</script>

<template>
    <UModal :open="open" :title="title" :description="description" @update:open="emit('update:open', $event)">
        <template #body>
            <div class="text-center py-4">
                <div
                    class="w-16 h-16 rounded-full bg-red-100 dark:bg-red-950 flex items-center justify-center mx-auto mb-4">
                    <UIcon name="i-tabler-alert-triangle" class="w-8 h-8 text-red-500" />
                </div>
                <slot>
                    <p class="text-slate-700 dark:text-slate-300">Apakah Anda yakin?</p>
                    <p v-if="itemName" class="font-bold text-slate-900 dark:text-white mt-2">
                        "{{ itemName }}"
                    </p>
                </slot>
            </div>
        </template>

        <template #footer>
            <div class="flex items-center justify-end gap-2">
                <UButton :label="cancelLabel" color="neutral" variant="outline" @click="handleCancel" />
                <UButton :label="confirmLabel" color="error" icon="i-tabler-trash" :loading="loading"
                    @click="handleConfirm" />
            </div>
        </template>
    </UModal>
</template>