<template>
  <div class="app">
    <a-layout class="layout">
      <SideBar></SideBar>

      <a-layout>
        <a-layout-header>
          <HomeHeader></HomeHeader>
        </a-layout-header>
        <a-layout style="padding: 0 24px;">
          <a-layout-content>
            <a-typography>
              <div class="header-bar">
                <a-typography-title :heading="5">
                  主页
                </a-typography-title>
                <div class="header-options">
                  <a-button type="primary" @click="handleNewDoc">
                    <template #icon>
                      <icon-plus />
                    </template>
                    新建文档
                  </a-button>
                </div>
              </div>
              <a-typography-paragraph>
              </a-typography-paragraph>
            </a-typography>
          </a-layout-content>
        </a-layout>
      </a-layout>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import documentApi from '@/api/document'

const router = useRouter()

definePageMeta({
  middleware: ['login-route']
})

const handleNewDoc = () => {
  documentApi.new().then(res => {
    const uid = res.data.data
    router.push(`/docs/${uid}`)
  })
}
</script>

<style scoped lang="scss">
.app {
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  background-color: var(--bg-body);
}

.layout {
  height: 100%;
  background: var(--bg-body);
  border: 1px solid var(--color-border);
}

.layout :deep(.arco-layout-footer) {
  height: 48px;
  color: var(--color-text-2);
  font-weight: 400;
  font-size: 14px;
  line-height: 48px;
}

.layout :deep(.arco-layout-content) {
  color: var(--color-text-2);
  font-weight: 400;
  font-size: 14px;
  background: var(--color-bg-3);
}

.header-bar {
  align-items: center;
  display: flex;
  flex-shrink: 0;
  position: relative;
}

.header-options {
  align-items: center;
  display: flex;
  flex: 0 0 auto;
  list-style: none;
  margin-left: auto;
}
</style>