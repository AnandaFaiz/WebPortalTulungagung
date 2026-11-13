<script setup>
import Cropper from 'cropperjs'
import 'cropperjs/dist/cropper.css'

const model = defineModel({ type: String, default: '' })
const props = defineProps({
    label: { type: String, default: 'Gambar' },
    aspectRatio: { type: Number, default: NaN },
    maxWidth: { type: Number, default: 1920 },
    maxHeight: { type: Number, default: 1080 },
    previewClass: { type: String, default: 'aspect-video' }
})

const API_BASE = 'http://localhost:8080/api'
const toast = useToast()
const fileInput = ref(null)
const imageElement = ref(null)
const cropper = ref(null)
const sourceUrl = ref('')
const open = ref(false)
const uploading = ref(false)
const zoom = ref(0)

const selectFile = () => fileInput.value?.click()

const handleFile = (event) => {
    const file = event.target.files?.[0]
    if (!file) return
    if (!file.type.startsWith('image/')) {
        toast.add({ title: 'Format tidak didukung', description: 'Pilih file gambar.', color: 'error' })
        return
    }

    sourceUrl.value = URL.createObjectURL(file)
    zoom.value = 0
    open.value = true
    nextTick(() => {
        cropper.value?.destroy()
        cropper.value = new Cropper(imageElement.value, {
            aspectRatio: Number.isNaN(props.aspectRatio) ? NaN : props.aspectRatio,
            viewMode: 1,
            autoCropArea: 1,
            responsive: true,
            background: false,
            movable: true,
            zoomable: true,
            cropBoxMovable: true,
            cropBoxResizable: true
        })
    })
}

const setZoom = (value) => {
    zoom.value = Math.max(-1, Math.min(1, value))
    cropper.value?.zoomTo(zoom.value)
}

const closeEditor = () => {
    cropper.value?.destroy()
    cropper.value = null
    open.value = false
    if (sourceUrl.value) URL.revokeObjectURL(sourceUrl.value)
    sourceUrl.value = ''
    if (fileInput.value) fileInput.value.value = ''
}

const uploadCrop = async () => {
    if (!cropper.value) return
    uploading.value = true
    try {
        const canvas = cropper.value.getCroppedCanvas({
            maxWidth: props.maxWidth,
            maxHeight: props.maxHeight,
            imageSmoothingEnabled: true,
            imageSmoothingQuality: 'high'
        })
        const blob = await new Promise(resolve => canvas.toBlob(resolve, 'image/jpeg', 0.9))
        if (!blob) throw new Error('Gambar hasil crop tidak tersedia')

        const formData = new FormData()
        formData.append('file', blob, `cropped-${Date.now()}.jpg`)
        const response = await $fetch(`${API_BASE}/admin/upload`, {
            method: 'POST',
            body: formData,
            credentials: 'include'
        })
        model.value = response.url
        toast.add({ title: 'Berhasil', description: 'Gambar berhasil di-crop dan diunggah', color: 'success' })
        closeEditor()
    } catch (error) {
        toast.add({ title: 'Gagal upload', description: error?.data?.statusMessage || error?.message || 'Terjadi kesalahan', color: 'error' })
    } finally {
        uploading.value = false
    }
}

const removeImage = () => { model.value = '' }

onBeforeUnmount(closeEditor)
</script>

<template>
    <div class="space-y-3">
        <div v-if="model"
            :class="['relative overflow-hidden rounded-lg border border-slate-200 dark:border-slate-700', previewClass]">
            <img :src="model" :alt="label" class="w-full h-full object-cover" />
            <div class="absolute top-2 right-2 flex gap-2">
                <UButton icon="i-tabler-replace" color="neutral" variant="solid" size="sm" title="Ganti gambar"
                    @click="selectFile" />
                <UButton icon="i-tabler-trash" color="error" variant="solid" size="sm" title="Hapus gambar"
                    @click="removeImage" />
            </div>
        </div>
        <button v-else type="button"
            class="w-full border-2 border-dashed border-slate-300 dark:border-slate-600 rounded-lg p-6 text-center hover:border-sky-500 transition-colors"
            @click="selectFile">
            <UIcon name="i-tabler-cloud-upload" class="w-8 h-8 text-slate-400 mx-auto mb-2" />
            <span class="block text-sm font-medium text-slate-700 dark:text-slate-300">Pilih {{ label.toLowerCase()
                }}</span>
            <span class="block text-xs text-slate-500 mt-1">Anda dapat crop, resize, dan zoom sebelum upload</span>
        </button>
        <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/webp" class="hidden"
            @change="handleFile" />

        <UModal v-model:open="open" :title="`Atur ${label}`" :ui="{ content: 'sm:max-w-4xl' }">
            <template #body>
                <div class="space-y-4">
                    <div
                        class="h-[55vh] min-h-64 bg-slate-950 rounded-lg overflow-hidden flex items-center justify-center">
                        <img ref="imageElement" :src="sourceUrl" :alt="label" class="block max-w-full" />
                    </div>
                    <div class="flex flex-wrap items-center justify-between gap-3">
                        <div class="flex items-center gap-2">
                            <UButton icon="i-tabler-zoom-out" color="neutral" variant="outline" title="Zoom out"
                                @click="setZoom(zoom - 0.1)" />
                            <input v-model.number="zoom" type="range" min="-1" max="1" step="0.1" class="w-32"
                                aria-label="Zoom gambar" @input="setZoom(zoom)" />
                            <UButton icon="i-tabler-zoom-in" color="neutral" variant="outline" title="Zoom in"
                                @click="setZoom(zoom + 0.1)" />
                        </div>
                        <div class="flex gap-2">
                            <UButton icon="i-tabler-rotate-clockwise" color="neutral" variant="outline"
                                title="Putar kanan" @click="cropper?.rotate(90)" />
                            <UButton icon="i-tabler-rotate-2" color="neutral" variant="outline" title="Putar kiri"
                                @click="cropper?.rotate(-90)" />
                        </div>
                    </div>
                    <p class="text-xs text-slate-500">Geser gambar dan ubah area crop sesuai kebutuhan label.</p>
                </div>
            </template>
            <template #footer>
                <div class="flex justify-end gap-2">
                    <UButton label="Batal" color="neutral" variant="outline" @click="closeEditor" />
                    <UButton label="Crop & Upload" icon="i-tabler-check" :loading="uploading" @click="uploadCrop" />
                </div>
            </template>
        </UModal>
    </div>
</template>
