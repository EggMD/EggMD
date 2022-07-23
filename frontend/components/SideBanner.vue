<template>
    <div class="banner">
        <div class="banner-inner">
            <a-carousel class="carousel" animation-name="fade">
                <a-carousel-item>
                    <div class="carousel-item">
                        <div class="carousel-image">
                            <img src="~/assets/img/undraw_typewriter.svg" width="350" />
                        </div>
                        <div class="carousel-title">蛋文档</div>
                        <div class="carousel-sub-title">Documents in eggs.</div>
                    </div>
                </a-carousel-item>
            </a-carousel>
        </div>
    </div>
</template>


<style lang="less" scoped>
.banner {
    display: flex;
    align-items: center;
    justify-content: center;

    &-inner {
        flex: 1;
        height: 100%;
    }
}

.carousel {
    height: 100%;

    &-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100%;
    }

    &-title {
        margin-top: 50px;
        color: var(--color-fill-1);
        font-weight: 500;
        font-size: 20px;
        line-height: 28px;
    }

    &-sub-title {
        margin-top: 8px;
        color: var(--color-text-3);
        font-size: 14px;
        line-height: 22px;
    }

    &-image {
        // width: 320px;
        margin-top: 30px;
    }
}
</style>
