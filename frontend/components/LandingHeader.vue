<template>
    <div class="navbar">
        <div class="left-side">
            <a-space>
                <img class="logo" alt="logo" src="~/assets/img/eggmd.svg" />
                <span :style="{ margin: 0, fontSize: '18px' }">
                    EggMD
                </span>
            </a-space>
        </div>
        <ul class="right-side" v-if="!isLogin">
            <nuxt-link to="/user/sign-in">
                <li class="nav-btn">
                    登录
                </li>
            </nuxt-link>
            <nuxt-link to="/user/sign-up">
                <li class="square-btn">
                    立即注册
                </li>
            </nuxt-link>
        </ul>
        <ul class="right-side" v-else>
            <nuxt-link to="/dashboard">
                <li class="nav-btn">
                    我的文档
                </li>
            </nuxt-link>
            <li>
                <profile-dropdown></profile-dropdown>
            </li>
        </ul>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

await userStore.getUserInfo()
const isLogin = ref<boolean>(userStore.isLogin)
</script>


<style scoped lang="scss">
.navbar {
    display: flex;
    justify-content: space-between;
    height: 100%;
    background-color: var(--color-bg-2);
    border-bottom: 1px solid var(--color-border);
    height: 60px;
}

.left-side {
    display: flex;
    align-items: center;
    padding-left: 60px;

    .logo {
        width: 40px;
    }
}

.right-side {
    margin: 0%;
    display: flex;
    padding-right: 60px;
    list-style: none;

    :deep(.locale-select) {
        border-radius: 20px;
    }

    li {
        display: flex;
        align-items: center;
        padding: 0 10px;
    }

    a {
        color: var(--color-text-1);
        text-decoration: none;
    }

    .nav-btn {
        color: #fa6167;
        align-items: center;
        box-sizing: border-box;
        display: flex;
        font-size: 16px;
        font-weight: 500;
        height: 60px;
        line-height: 24px;
        padding: 8px 24px;
        position: relative;
    }

    .square-btn {
        align-items: center;
        background: #fa8487;
        box-sizing: border-box;
        color: #fff;
        display: flex;
        font-size: 16px;
        font-weight: 500;
        height: 60px;
        line-height: 24px;
        padding: 8px 24px;
        position: relative;
    }
}
</style>