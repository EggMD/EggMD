import { defineNuxtConfig } from 'nuxt'
import Components from 'unplugin-vue-components/vite'
import IconsResolver from 'unplugin-icons/resolver'
import { ArcoResolver } from "unplugin-vue-components/resolvers";

// https://v3.nuxtjs.org/api/configuration/nuxt.config
export default defineNuxtConfig({
    ssr: false,
    css: [
        '@arco-design/web-vue/lib/message/style/index.less',
    ],
    build: {
        transpile: ['compute-scroll-into-view'],
    },
    modules: ['@pinia/nuxt'],
    vite: {
        optimizeDeps: {
            include: ["@editorjs/editorjs"],
        },
        // example of overwriting arco-design .less vars
        css: {
            preprocessorOptions: {
                less: {
                    // modifyVars: {
                    //     'arcoblue-6': 'black',
                    //     '@menu-horizontal-padding-vertical': '0px',
                    //     '@menu-horizontal-padding-horizontal': '0px',
                    //     '@tooltip-border-radius': '5px',
                    //     '@tooltip-mini-padding-vertical': '1px',
                    //     '@tooltip-mini-padding-horizontal': '5px',
                    //     '@menu-vertical-padding-vertical': '0px',
                    //     '@menu-vertical-padding-horizontal': '0px',
                    //     '@menu-collapse-padding-vertical': '0px',
                    //     '@menu-collapse-padding-horizontal': '0px',
                    //     '@menu-item-indent-spacing': '0px',
                    //     '@menu-font-weight-item-selected': '700',
                    //     '@btn-font-weight': '500',
                    //     '@btn-border-radius': '5px'
                    // },
                    javascriptEnabled: true,
                },
            },
        },
        plugins: [
            Components({
                dts: true, // enabled by default if `typescript` is installed
                resolvers: [IconsResolver({
                    // to avoid naming conflicts
                    // a prefix can be specified for icons
                    prefix: 'i'
                }),
                ArcoResolver({
                    importStyle: 'less',
                    resolveIcons: true
                }
                )],
            })
        ],
        server: {
            proxy: {
                '/api': {
                    target: process.env.SERVER_URL
                }
            }
        },
    }
})
