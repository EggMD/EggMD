<template>
    <DocHeader :uid="uid" :read-only="isReadOnly"></DocHeader>
    <DocEditor :uid="uid" :read-only="isReadOnly"></DocEditor>
</template>

<script setup lang="ts">
import documentApi from '@/api/document'
import { documentMeta } from '@/models/document'
import { useUserStore } from '@/store/user';

const route = useRoute()
const userStore = useUserStore()

const uid = ref<string>(route.params.uid.toString())

const currentUser = await userStore.getUserInfo()
const meta = ref<documentMeta>();

await documentApi.meta(uid.value).then(res => {
    meta.value = res.data.data
})
const isReadOnly = !userStore.isLogin || (meta.value.Owner.UID !== currentUser.UID)

</script>