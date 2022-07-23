/// <reference types="vite/client" />

declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/ban-types
    const component: DefineComponent<{}, {}, any>
    export default component
}

declare module 'notiwind' {
    import { Plugin } from 'vue'
    const Notifications: Plugin
    export default Notifications
    export function notify(
        notification: {
            group: string
            title: string
            text?: string
            kind?: 'error' | 'success'
        },
        time: number
    ): void
}

declare type requestBody = { [key: string]: any }
declare type requestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'
declare interface map {
    [key: string]: any
    [index: number]: any
}
