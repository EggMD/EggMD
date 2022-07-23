<template>
    <a-page-header>
        <template #title>
            {{ meta.Title }}
        </template>
        <template #subtitle>
            {{ dayjs(meta.UpdatedAt).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
        <template #extra v-if="!props.readOnly">
            <div class="extra-bar">
                <div class="extra-share">
                    <a-button type="primary">
                        <template #icon>
                            <icon-link />
                        </template>
                        <template #default>分享</template>
                    </a-button>
                </div>
                <div class="extra-setting">
                    <a-button type="text" shape="circle" @click="settingModalVisible = true">
                        <icon-more />
                    </a-button>
                    <a-modal v-model:visible="settingModalVisible" @ok="handleChangeSetting">
                        <template #title>文档设置 </template>
                        <a-form :model="settingForm">
                            <a-form-item field="Title" label="文档标题">
                                <a-input v-model="settingForm.Title" />
                            </a-form-item>
                        </a-form>
                    </a-modal>
                </div>
                <span class="extar-line"></span>
                <div>
                    <profile-dropdown></profile-dropdown>
                </div>
            </div>
        </template>
        <p></p>
    </a-page-header>
</template>

<script setup lang="ts">
import documentApi from '@/api/document'
import { documentMeta } from '@/models/document'
import { Message } from '@arco-design/web-vue';
import dayjs from 'dayjs'

const router = useRouter()

const props = defineProps({
    uid: { default: '' },
    readOnly: { default: true },
});
const uid = props.uid

const meta = ref<documentMeta>();
const getMeta = async () => {
    await documentApi.meta(uid).then(res => {
        meta.value = res.data.data
    }).catch(() => {
        router.push('/dashboard')
    })
}
await getMeta()


const settingModalVisible = ref<boolean>(false)
const settingForm = ref<Object>({
    Title: '',
})
const getSetting = async () => {
    documentApi.getSetting(uid).then(res => {
        settingForm.value = res.data.data
    })
}
await getSetting()

const handleChangeSetting = () => {
    documentApi.updateSetting(uid, settingForm.value).then(async res => {
        await getSetting()
        await getMeta()
        settingModalVisible.value = false
        Message.info(res.data.data)
    })
}

</script>

<style lang="scss" scoped>
.extra-bar {
    flex: 1;
    display: flex;
    justify-content: flex-end;
    align-items: center;
}

.extar-line {
    margin: 0 16px 0 14px;
    width: 1px;
    height: 32px;
    background-color: rgb(222, 224, 227);
    margin: 0 20px 0 20px;
}

.extra-setting {
    display: inline-block;
    font-style: normal;
    line-height: 0;
    text-align: center;
    font-size: 30px;
    margin-left: 20px;
}

.extra-share {}
</style>