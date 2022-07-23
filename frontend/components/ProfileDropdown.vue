<template>
    <a-dropdown position="bl">
        <a-avatar :size="32" class="avatar">{{ userInfo && userInfo.NickName[0] }}</a-avatar>
        <template #content>
            <div class="panel-header">
                <div class="panel-avatar">
                    <a-avatar :size="50">{{ userInfo && userInfo.NickName[0] }}</a-avatar>
                </div>
                <div class="panel-name">
                    {{ userInfo.NickName }}
                </div>
            </div>
            <a-doption class="doption">设置</a-doption>
            <a-doption class="doption" @click="handleSignOut">退出登录</a-doption>
        </template>
    </a-dropdown>
</template>

<script setup lang="ts">
import { useUserStore } from '@/store/user'
import authApi from '@/api/auth'
import { Message } from '@arco-design/web-vue'

const userStore = useUserStore()
const router = useRouter()

const userInfo = await userStore.getUserInfo()

const handleSignOut = () => {
    authApi.signOut().then(res => {
        userStore.setLogout();
        Message.info(res.data.data);
        router.go(0)
    })
}
</script>

<style lang="scss" scoped>
.avatar {
    cursor: pointer;
}

.doption {
    height: 40px;
}

.panel-header {
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    -webkit-box-orient: vertical;
    -webkit-box-direction: normal;
    -ms-flex-direction: column;
    flex-direction: column;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    -webkit-box-pack: center;
    -ms-flex-pack: center;
    justify-content: center;
    padding: 24px 8px 16px;

    .panel-avatar {
        width: 48px;
        height: 48px;
        border-radius: 50%;
        margin-bottom: 8px;
        -o-object-fit: cover;
        object-fit: cover;
    }

    .panel-name {
        font-size: 14px;
        line-height: 20px;
        font-weight: 500;
        color: var(--text-title);
        margin-bottom: 4px;
        text-align: center;
        overflow: hidden;
        word-break: break-all;
        -o-text-overflow: ellipsis;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
    }
}
</style>