import { defineStore } from 'pinia'
import type { profileInfo } from '@/models/profile'
import authApi from '@/api/auth'
import { toHandlers } from 'nuxt/dist/app/compat/capi'
import user from '@/api/auth'

interface authStore {
    isLogin: boolean
    user: profileInfo
}

const userInfo = ref<profileInfo>()
const isLogin = ref<boolean>(false)


authApi.getProfile().then(res => {
    userInfo.value = res.data.data
    isLogin.value = true
}).catch(() => {
    isLogin.value = false
})

export const useUserStore = defineStore('user', {
    state: () => (<authStore>{
        isLogin: isLogin.value,
        user: userInfo.value,
    }),
    getters: {
        loginStatus: state => state.isLogin,
        userInfo: state => state.user,
    },
    actions: {
        async getUserInfo() {
            if (!this.isLogin) {
                await authApi.getProfile().then(res => {
                    this.user = res.data.data
                    this.isLogin = true
                }).catch(() => {
                    this.isLogin = false
                })
            }
            return this.user
        },
        setUserInfo(user: profileInfo) {
            this.user = user
        },
        setLogin() {
            this.isLogin = true
        },
        setLogout() {
            this.isLogin = false
        }
    },
})