<template>
    <a-layout-sider hide-trigger :style="{ width: '260px' }" class="sider">
        <div class="brand">
            <img class="logo" alt="logo" src="~/assets/img/eggmd.svg" />
            <div class="logo-text">EggMD</div>
        </div>

        <a-menu class="menu" :defaultSelectedKeys="['home']" :style="{ width: '100%' }">
            <a-menu-item key="home">
                <icon-home size="20" />主页
            </a-menu-item>

            <a-sub-menu key="space">
                <template #title>
                    <icon-folder size="20" />我的空间
                </template>
                <a-menu-item key="1_1">Menu 1</a-menu-item>
                <a-menu-item key="1_2">Menu 2</a-menu-item>
                <a-sub-menu key="2" title="Navigation 2">
                    <a-menu-item key="2_1">Menu 1</a-menu-item>
                    <a-menu-item key="2_2">Menu 2</a-menu-item>
                </a-sub-menu>
                <a-sub-menu key="3" title="Navigation 3">
                    <a-menu-item key="3_1">Menu 1</a-menu-item>
                    <a-menu-item key="3_2">Menu 2</a-menu-item>
                    <a-menu-item key="3_3">Menu 3</a-menu-item>
                </a-sub-menu>
            </a-sub-menu>

            <a-menu-item key="favorites">
                <icon-star size="20" />收藏
            </a-menu-item>

            <a-menu-item key="trash">
                <icon-delete size="20" />回收站
            </a-menu-item>
        </a-menu>
    </a-layout-sider>
</template>

<style lang="scss">
.sider {
    height: 100%;
    background-color: var(--color-fill-2);
}

.brand {
    margin: 0 12px;
    padding: 5px 8px;
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    overflow: hidden;
    padding-top: 16px;
    padding-bottom: 40px;
}

.logo-text {
    font-size: 14px;
    line-height: 30px;
    padding-left: 5px;
    color: var(--text-title);
    display: flex;
    font-weight: 700;
}

.arco-layout-sider .logo {
    height: 32px;
    background: rgba(255, 255, 255, 0.2);
}

.arco-layout-sider-light .logo {
    background: var(--color-fill-2);
}

.arco-menu-light {
    background: var(--color-fill-2);

    .arco-menu-item {
        font-weight: 800;
        background: var(--color-fill-2);
        display: flex;
        align-items: center;

        &:hover {
            background: var(--color-fill-3);
        }
    }

    .arco-menu-item.arco-menu-selected {
        font-weight: 800;
        --B500: #3370ff;
        --B500-raw: 51, 112, 255;
        --primary-pri-500: var(--B500);
        --primary-pri-500-raw: var(--B500-raw);
        background-color: rgba(var(--primary-pri-500-raw), 0.08);
    }

    .arco-menu-inline-header {
        display: flex;
        align-items: center;

        font-weight: 800;
        background: var(--color-fill-2);

        &:hover {
            background: var(--color-fill-3);
        }
    }
}
</style>