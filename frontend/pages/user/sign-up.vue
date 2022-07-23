<template>
    <div class="container">
        <div class="logo">
            <img class="logo" alt="logo" src="~/assets/img/eggmd.svg" />
            <div class="logo-text">EggMD</div>
        </div>
        <SideBanner></SideBanner>
        <div class="content">
            <div class="content-inner">
                <div class="register-form-wrapper">
                    <div class="register-form-title">欢迎使用 EggMD</div>
                    <div class="register-form-sub-title">30s 注册你的 EggMD 账号</div>

                    <a-form :model="signUpForm" class="register-form" layout="vertical" @submit="handleSubmit">
                        <a-form-item field="Email" label="电子邮箱" :rules="[{ required: true, message: '电子邮箱不能为空' }]"
                            :validate-trigger="['change', 'blur', 'submit']">
                            <a-input v-model="signUpForm.Email"></a-input>
                        </a-form-item>

                        <a-form-item field="LoginName" label="用户名" :rules="[{ required: true, message: '用户名不能为空' }]"
                            :validate-trigger="['change', 'blur', 'submit']">
                            <a-input v-model="signUpForm.LoginName"></a-input>
                        </a-form-item>


                        <a-form-item field="Password" label="密码" :rules="[{ required: true, message: '密码不能为空' }]"
                            :validate-trigger="['change', 'blur', 'submit']">
                            <a-input-password v-model="signUpForm.Password" allow-clear></a-input-password>
                        </a-form-item>

                        <a-space :size="16" direction="vertical">
                            <a-button type="primary" html-type="submit" long>
                                注册新账号
                            </a-button>
                            <NuxtLink to="/user/sign-in" style="text-decoration: none;">
                                <a-button type="text" long class="register-form-register-btn">
                                    已有账号？前往登录
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
import { Message } from '@arco-design/web-vue'
const router = useRouter()

definePageMeta({
    middleware: ['no-login-route']
})

const signUpForm = ref<Object>({
    Email: '',
    LoginName: '',
    Password: '',
})

const handleSubmit = () => {
    authApi.signUp(signUpForm.value).then(res => {
        Message.info(res.data.data);
        router.push({ path: '/user/sign-in' });
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

.register-form {
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
