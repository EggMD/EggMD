import axios from 'axios'
import type { AxiosRequestConfig, AxiosRequestHeaders } from 'axios'
import { HTTP_REQUEST_TIMEOUT_MILLISECONDS } from '@/consts'
import { Message } from '@arco-design/web-vue'

const instance = axios.create({
    baseURL: process.client ? '/' : process.env.SERVER_URL,
    timeout: HTTP_REQUEST_TIMEOUT_MILLISECONDS
})

instance.interceptors.request.use(
    (config: AxiosRequestConfig) => {
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

const request = (url: string, method: requestMethod, data?: requestBody) => {
    return instance.request({ url: `/api${url}`, method, data }).catch((error) => {
        if (error.code === `ECONNABORTED` && error.message.includes(`timeout`)) {
            console.log('Request timeout')
        }
        // if (error.response.data.error === 40300) {
        //     const authStore = useAuthStore()
        //     authStore.setLogout()
        //     router.push('/login')
        // }
        if (error.response.data.msg !== undefined && error.response.data.error / 100 !== 401) {
            Message.error(error.response.data.msg)
        }
        throw error
    })
}

const httpMethod = {
    GET: (url: string) => request(url, 'GET'),
    POST: (url: string, data?: requestBody) => request(url, 'POST', data),
    PUT: (url: string, data: requestBody) => request(url, 'PUT', data),
    DELETE: (url: string) => request(url, 'DELETE'),
    DOWNLOAD: (url: string) => instance.request({ url, method: 'GET', responseType: 'blob' }),
}

export default httpMethod