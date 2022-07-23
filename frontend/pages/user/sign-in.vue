<template>
    <div class="container">
        <div class="logo">
            <img class="logo" alt="logo" src="~/assets/img/eggmd.svg" />
            <div class="logo-text">EggMD</div>
        </div>
        <SideBanner></SideBanner>
        <div class="content">
            <div class="content-inner">
                <div class="login-form-wrapper">
                    <div class="login-form-title">欢迎回来</div>
                    <div class="login-form-sub-title">登录你的 EggMD 账号</div>

                    <a-form ref="loginForm" :model="signInForm" class="login-form" layout="vertical"
                        @submit="handleSubmit">
                        <a-form-item field="Email" label="电子邮箱" :rules="[{ required: true, message: '电子邮箱不能为空' }]"
                            :validate-trigger="['change', 'blur', 'submit']">
                            <a-input v-model="signInForm.Email"></a-input>
                        </a-form-item>

                        <a-form-item field="Password" label="密码" :rules="[{ required: true, message: '密码不能为空' }]"
                            :validate-trigger="['change', 'blur', 'submit']">
                            <a-input-password v-model="signInForm.Password" allow-clear></a-input-password>
                        </a-form-item>

                        <a-space :size="16" direction="vertical">
                            <a-button type="primary" html-type="submit" long>
                                登录
                            </a-button>
                            <NuxtLink to="/user/sign-up" style="text-decoration: none;">
                                <a-button type="text" long class="login-form-register-btn">
                                    没有账号？前往注册
                                </a-button>
                            </NuxtLink>
                        </a-space>
                    </a-form>
                </div>
            </div>
            <div class="footer">

            </div>
        </div>
    </div>
</template>


<script setup lang="ts">
import { ref } from 'vue';
import authApi from '@/api/auth'
import { useUserStore } from '@/store/user'

definePageMeta({
    middleware: ['no-login-route']
})

const router = useRouter()
const userStore = useUserStore()

const signInForm = ref<Object>({
    Email: '',
    Password: '',
})

const handleSubmit = () => {
    authApi.signIn(signInForm.value).then(() => {
        userStore.getUserInfo().then(() => {
            router.push({ path: '/dashboard' })
        })
    })
}
</script>


<style lang="less" scoped>
.container {
    display: flex;
    height: 100vh;

    .banner {
        width: 550px;
        background: linear-gradient(163.85deg, #1d2129 0%, #00308f 100%);
    }

    .content {
        position: relative;
        display: flex;
        flex: 1;
        align-items: center;
        justify-content: center;
        padding-bottom: 40px;
    }

    .footer {
        position: absolute;
        right: 0;
        bottom: 0;
        width: 100%;
    }
}

.logo {
    position: fixed;
    top: 24px;
    left: 22px;
    z-index: 1;
    display: inline-flex;
    align-items: center;
    width: 40px;

    &-text {
        margin-right: 4px;
        margin-left: 50px;
        margin-top: 10px;
        color: var(--color-fill-1);
        font-size: 20px;
    }
}

.login-form {
    &-wrapper {
        width: 320px;
    }

    &-title {
        color: var(--color-text-1);
        font-weight: 500;
        font-size: 24px;
        line-height: 32px;
    }

    &-sub-title {
        color: var(--color-text-3);
        font-size: 16px;
        line-height: 24px;
    }

    &-error-msg {
        height: 32px;
        color: rgb(var(--red-6));
        line-height: 32px;
    }

    &-password-actions {
        display: flex;
        justify-content: space-between;
    }

    &-register-btn {
        text-decoration: none;
        color: var(--color-text-3) !important;
    }
}
</style>
