<template>
    <div id="editor" class="editor"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, getCurrentInstance } from 'vue'
import documentApi from '@/api/document'
import DragDrop from 'editorjs-drag-drop';
import EditorJS from "@editorjs/editorjs";
import Undo from 'editorjs-undo';

const { default: Header } = await import('@editorjs/header');
const { default: Table } = await import('@editorjs/table');
const { default: NestedList } = await import('@editorjs/nested-list');
const { default: CheckList } = await import('@editorjs/checklist');
const { default: ImageTool } = await import('@editorjs/image');
const { default: Underline } = await import('@editorjs/underline');
const { default: Marker } = await import('@editorjs/marker');
const { default: Warning } = await import('@editorjs/warning');
const { default: InlineCode } = await import('@editorjs/inline-code');
const { default: Delimiter } = await import('@editorjs/delimiter');
const { default: CodeTool } = await import('@editorjs/code');

const tools = {
    header: { class: Header },
    table: { class: Table, inlineToolbar: true, config: { rows: 2, cols: 3, }, },
    list: { class: NestedList, inlineToolbar: true, },
    checklist: { class: CheckList, inlineToolbar: true, },
    image: {
        class: ImageTool, config: {
            endpoints: {
                byFile: 'http://localhost:8008/uploadFile', // Your backend file uploader endpoint
                byUrl: 'http://localhost:8008/fetchUrl', // Your endpoint that provides uploading by Url
            }
        }
    },
    code: CodeTool,
    inlineCode: { class: InlineCode, shortcut: 'CMD+SHIFT+M', },
    underline: Underline,
    marker: { class: Marker, shortcut: 'CMD+SHIFT+M', },
    warning: { class: Warning, inlineToolbar: true, shortcut: 'CMD+SHIFT+W', config: { titlePlaceholder: 'Title', messagePlaceholder: 'Message', }, },
    delimiter: Delimiter,
}


const props = defineProps({
    uid: { default: '' },
    readOnly: { default: true },
});

const content = ref()
const change_ready = ref(true);
const onChangeTimeout = ref(0);
const onChange = (api, event) => {
    api.saver.save().then(async (data) => {
        data.value = data;
        if (change_ready.value) {

            delete (data.value.value)
            documentApi.save(uid, data)

            if (onChangeTimeout.value > 0) {
                change_ready.value = false;

                setTimeout(() => {
                    change_ready.value = true;
                }, onChangeTimeout.value);
            }
        }
    });
};

const uid = props.uid
await documentApi.content(uid).then(res => {
    content.value = res.data.data
})

onMounted(() => {
    const editor = new EditorJS({
        holder: "editor",
        minHeight: 0,
        onChange: onChange,
        data: content.value,
        logLevel: "ERROR",
        autofocus: true,
        placeholder: '写点什么吧~',
        tools: tools,
        readOnly: props.readOnly,
        onReady: () => {
            new Undo({ editor });
            new DragDrop(editor);
        },
    });
});
</script>

<style lang="scss" scoped>
.editor {
    font-size: 16px;
    padding-left: 15px;
    padding-right: 15px;
}
</style>

