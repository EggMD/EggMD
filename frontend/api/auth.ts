import request from '@/utils/requests'

const user = {
    signUp: (data: requestBody) => request.POST('/user/sign-up', data),
    signIn: (data: requestBody) => request.POST('/user/sign-in', data),
    signOut: () => request.POST('/user/sign-out'),
    getProfile: () => request.GET('/user/profile')
}

export default user