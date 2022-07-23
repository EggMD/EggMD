<template>
    <header class="header">
        <div class="search-box">
            <a-input :style="{ width: '540px' }" placeholder="搜索" allow-clear>
                <template #prefix>
                    <icon-search />
                </template>
            </a-input>
        </div>
        <ol class="options">
            <li>
                <profile-dropdown></profile-dropdown>
            </li>
        </ol>
    </header>
</template>

<style lang="scss">
.arco-dropdown-list-wrapper {
    max-height: 400px;
    width: 224px;
}
</style>

<style scoped lang="scss">
.header {
    position: relative;
    height: 64px;
    padding: 0px 24px;
    display: flex;
    -webkit-box-align: center;
    align-items: center;
    flex-shrink: 0;
}

.search-box {
    margin: 0px;
    height: auto;
}

.options {
    display: flex;
    -webkit-box-align: center;
    align-items: center;
    list-style: none;
    margin-left: auto;
    flex: 0 0 auto;
}
</style>