<template>
  <div class="border-out">
    <div class="border-in">
      <transition name="fade">
        <div v-if="isVisible" class="splash fade-element" :style="splashStyle"></div>
      </transition>
    </div>
  </div>
</template>

<script setup lang="ts">

import {onMounted, ref} from "vue"
import {useGlobalStore} from 'src/stores/globalData'
import {callSrv, decodeApiData, firstLetterToLower} from 'components/utils'
import {GetConfig, GetTheme} from 'app/wailsjs/go/controller/Controller'
import {storeToRefs} from "pinia";
import {initThemeColor} from 'src/pages/classic/HeaderBarTool.vue'


const global = useGlobalStore();
const {config, theme, platformUi, simulatorMap} = storeToRefs(global);

theme.value = "Default"
const splashStyle: any = ref(null);

const isVisible = ref(false);

//启动时检查主题，跳转到默认主题
GetTheme().then((result: string) => {
  let resp = decodeApiData(result)
  if (resp.err == "") {
    theme.value = "/" + firstLetterToLower(resp.data)
  }
})

onMounted(async () => {

  let conf = await callSrv("GetConfig")

  //初始化配置
  global.initData(conf);

  if (config.value.SplashScreen.Time == 0) {
    global.goto(theme.value, {}, 1)
    return
  }

  //初始化主题颜色
  initThemeColor()


  splashStyle.value = {
    "background-image": "url(" + config.value.SplashScreen.Src + ")",
    "background-size": config.value.SplashScreen.Size,
  }

  // 淡入
  isVisible.value = true;

  // 延迟跳转
  setTimeout(() => {
    isVisible.value = false;
    global.goto(theme.value, {}, 1)
  }, config.value.SplashScreen.Time * 1000);

})

</script>

<style scoped>
.splash {
  width: 100vw;
  height: 100vh;
  background: center center no-repeat;
  transition: opacity 0.5s ease-in-out;
}

</style>
