import request from '@/utils/requests'

const document = {
    new: () => request.POST("/doc"),
    meta: (uid: string) => request.GET(`/doc/${uid}/meta`),
    content: (uid: string) => request.GET(`/doc/${uid}/content`),
    save: (uid: string, data: requestBody) => request.POST(`/doc/${uid}/save`, data),
    getSetting: (uid: string) => request.GET(`/doc/${uid}/setting`),
    updateSetting: (uid: string, data: requestBody) => request.POST(`/doc/${uid}/setting`, data)
}

export default document