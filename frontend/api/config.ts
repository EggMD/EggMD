import request from '@/utils/requests'

const config = {
    getGlobal: () => request.GET('/config/global')
}

export default config