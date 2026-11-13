<script setup>
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import TextAlign from '@tiptap/extension-text-align'

const model = defineModel({ type: String, default: '' })
const props = defineProps({
    fullHeight: { type: Boolean, default: false }
})

const editor = useEditor({
    content: model.value,
    extensions: [
        StarterKit.configure({
            link: false
        }),
        Link.configure({
            openOnClick: false,
            autolink: true,
            linkOnPaste: true
        }),
        TextAlign.configure({
            types: ['heading', 'paragraph']
        })
    ],
    editorProps: {
        attributes: {
            class: props.fullHeight
                ? 'prose prose-sm focus:outline-none min-h-[520px] max-h-[calc(100vh-220px)] h-full overflow-y-auto p-4 bg-white dark:bg-slate-900 rounded-b-lg text-slate-800 dark:text-slate-200'
                : 'prose prose-sm focus:outline-none min-h-[150px] max-h-[300px] overflow-y-auto p-4 bg-white dark:bg-slate-900 rounded-b-lg text-slate-800 dark:text-slate-200',
        },
    },
    onUpdate: ({ editor }) => {
        model.value = editor.getHTML()
    },
})

const linkModalOpen = ref(false)
const linkUrl = ref('')

const setLink = () => {
    if (!editor.value) return

    linkUrl.value = editor.value.getAttributes('link').href || ''
    linkModalOpen.value = true
}

const applyLink = () => {
    if (!editor.value) return

    const trimmedUrl = linkUrl.value.trim()
    if (!trimmedUrl) {
        editor.value.chain().focus().unsetLink().run()
    } else {
        const href = /^https?:\/\//i.test(trimmedUrl) ? trimmedUrl : `https://${trimmedUrl}`
        editor.value.chain().focus().setLink({ href }).run()
    }

    linkModalOpen.value = false
}

const unsetLink = () => {
    editor.value?.chain().focus().unsetLink().run()
}

// Watch jika model diubah dari luar (misal saat pilih template)
watch(model, (newValue) => {
    if (editor.value && editor.value.getHTML() !== newValue) {
        editor.value.commands.setContent(newValue || '')
    }
})

onBeforeUnmount(() => {
    editor.value?.destroy()
})
</script>

<template>
    <div class="border border-slate-200 dark:border-slate-700 rounded-lg overflow-hidden"
        :class="{ 'h-full': fullHeight }">
        <!-- Toolbar -->
        <div v-if="editor"
            class="bg-slate-50 dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700 p-2 flex flex-wrap gap-1">
            <UButton icon="i-tabler-bold" color="neutral" variant="ghost" size="xs"
                @click="editor.chain().focus().toggleBold().run()"
                :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive('bold') }" />
            <UButton icon="i-tabler-italic" color="neutral" variant="ghost" size="xs"
                @click="editor.chain().focus().toggleItalic().run()"
                :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive('italic') }" />
            <UButton icon="i-tabler-list" color="neutral" variant="ghost" size="xs"
                @click="editor.chain().focus().toggleBulletList().run()"
                :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive('bulletList') }" />
            <UButton icon="i-tabler-list-numbers" color="neutral" variant="ghost" size="xs"
                @click="editor.chain().focus().toggleOrderedList().run()"
                :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive('orderedList') }" />
            <span class="w-px h-5 bg-slate-200 dark:bg-slate-700 mx-1" />
            <UButton icon="i-tabler-align-left" color="neutral" variant="ghost" size="xs" title="Rata kiri"
                @click="editor.chain().focus().setTextAlign('left').run()"
                :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive({ textAlign: 'left' }) }" />
            <UButton icon="i-tabler-align-center" color="neutral" variant="ghost" size="xs" title="Rata tengah"
                @click="editor.chain().focus().setTextAlign('center').run()"
                :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive({ textAlign: 'center' }) }" />
            <UButton icon="i-tabler-align-right" color="neutral" variant="ghost" size="xs" title="Rata kanan"
                @click="editor.chain().focus().setTextAlign('right').run()"
                :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive({ textAlign: 'right' }) }" />
            <UButton icon="i-tabler-align-justified" color="neutral" variant="ghost" size="xs" title="Justify"
                @click="editor.chain().focus().setTextAlign('justify').run()"
                :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive({ textAlign: 'justify' }) }" />
            <span class="w-px h-5 bg-slate-200 dark:bg-slate-700 mx-1" />
            <UButton icon="i-tabler-link" color="neutral" variant="ghost" size="xs" title="Tambah atau ubah hyperlink"
                @click="setLink" :class="{ 'bg-slate-200 dark:bg-slate-700': editor.isActive('link') }" />
            <UButton icon="i-tabler-link-off" color="neutral" variant="ghost" size="xs" title="Lepas hyperlink"
                :disabled="!editor.isActive('link')" @click="unsetLink" />
            <UButton icon="i-tabler-arrow-back-up" color="neutral" variant="ghost" size="xs"
                @click="editor.chain().focus().undo().run()" />
            <UButton icon="i-tabler-arrow-forward-up" color="neutral" variant="ghost" size="xs"
                @click="editor.chain().focus().redo().run()" />
        </div>

        <!-- Area Editor -->
        <EditorContent :editor="editor" />

        <UModal v-model:open="linkModalOpen" title="Tambah Hyperlink">
            <template #body>
                <form id="hyperlink-form" class="space-y-3" @submit.prevent="applyLink">
                    <UFormField label="URL Hyperlink" hint="Contoh: https://tulungagung.go.id">
                        <UInput v-model="linkUrl" type="url" placeholder="https://..." autofocus class="w-full" />
                    </UFormField>
                </form>
            </template>
            <template #footer>
                <div class="flex justify-end gap-2">
                    <UButton label="Batal" color="neutral" variant="outline" @click="linkModalOpen = false" />
                    <UButton form="hyperlink-form" type="submit" label="Simpan Link" icon="i-tabler-link" />
                </div>
            </template>
        </UModal>
    </div>
</template>