<template>
  <div class="absolute-full title-image">
    <q-img :src="titleInfo.image" loading="lazy" :class="{'mask': titleInfo.hideText == 0}"/>
    <h1 v-if="titleInfo.hideText == 0">{{ titleInfo.title }}</h1>
    <h2 v-if="titleInfo.hideText == 0" >{{ titleInfo.subTitle }}</h2>
  </div>
</template>

<script setup lang="ts">

import {onMounted, ref, watch} from "vue";
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';

const global = useGlobalStore();
const {config} = storeToRefs(global);

let titleInfo = ref({
      image: "",
      title: "",
      subTitle: "",
      hideText: 0,
    }
)

onMounted(() => {
  titleInfo.value = {
    image: config.value.SoftName.Image,
    title: config.value.SoftName.Name,
    subTitle: config.value.SoftName.SubName,
    hideText: config.value.SoftName.HideText,
  }
})

watch(config, (newValue, oldValue) => {
  titleInfo.value = {
    image: newValue.SoftName.Image,
    title: newValue.SoftName.Name,
    subTitle: newValue.SoftName.SubName,
    hideText: newValue.SoftName.HideText,
  }
});

</script>

<style scoped>
@import "src/css/classic/common.css";
/* 标题 */
.title-image {
  display: block;
  height: 100%;
  position: relative;
  overflow: hidden;
  padding: 0;
  background: var(--color-1);
}

.title-image .q-img {
  width: 100%;
  height: 100%;
  position: absolute;
}

.title-image .mask {
  opacity: 0.5;
}


.title-image h1, .title-image h2 {
  line-height: 1em;
  position: absolute;
  margin: 0;
  left: 10px;
}

.title-image h1 {
  font-size: 24px;
  font-weight: bold;
  bottom: 25px;
}

.title-image h2 {
  font-size: 14px;
  line-height: 1em;
  font-weight: normal;
  bottom: 10px;
}

</style>
