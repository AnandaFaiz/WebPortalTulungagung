<script setup>
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'

const model = defineModel({ type: String, default: '' })

const editor = useEditor({
    content: model.value,
    extensions: [
        StarterKit.configure({
            link: { openOnClick: false }
        })
    ],
    editorProps: {
        attributes: {
            class: 'prose prose-sm focus:outline-none min-h-[150px] max-h-[300px] overflow-y-auto p-4 bg-white dark:bg-slate-900 rounded-b-lg text-slate-800 dark:text-slate-200',
        },
    },
    onUpdate: ({ editor }) => {
        model.value = editor.getHTML()
    },
})

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
    <div class="border border-slate-200 dark:border-slate-700 rounded-lg overflow-hidden">
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
            <UButton icon="i-tabler-arrow-back-up" color="neutral" variant="ghost" size="xs"
                @click="editor.chain().focus().undo().run()" />
            <UButton icon="i-tabler-arrow-forward-up" color="neutral" variant="ghost" size="xs"
                @click="editor.chain().focus().redo().run()" />
        </div>

        <!-- Area Editor -->
        <EditorContent :editor="editor" />
    </div>
</template>