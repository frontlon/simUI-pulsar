<template>
  <div class="border-out" @resize="windowResize">
    <div :class="'border-in bg fuzzy' + platformUi.BackgroundFuzzy" :style="{ fontFamily: currentFont }">
      <div class="bg-mask"></div>
      <header-bar-tool/>
      <rom-list-bar/>
    </div>
  </div>
</template>

<script setup lang="ts">

import {onBeforeUnmount, onMounted, ref} from "vue"
import {useGlobalStore} from 'src/stores/globalData'
import {callSrv, decodeApiData, loadFont, notify} from 'components/utils'
import {GetAllSimulator, GetConfig, GetPlatformUi, UpdateOneConfig} from 'app/wailsjs/go/controller/Controller'
import {storeToRefs} from "pinia";
import {initPlatform} from 'src/pages/tiny/Platform.vue'
import {addKeyboardEvent, removeKeyboardEvent} from 'pages/tiny/modules/EventKeyboard.vue'
import {initEventGamePad} from 'pages/tiny/modules/EventGamepad.vue'
import RomListBar, {changeContentBackground} from 'src/pages/tiny/RomListBar.vue'
import HeaderBarTool, {initThemeColor} from 'src/pages/tiny/HeaderBarTool.vue'
import {debounce} from "quasar";

const global = useGlobalStore();
const {lang,config, theme, activeRom, platformUi, simulatorMap} = storeToRefs(global);
const currentFont = ref("Arial")

onMounted(async () => {

  theme.value = "Tiny"

  let conf = await callSrv("GetConfig")

  let ui = await callSrv("GetPlatformUi", conf.Config.Platform, theme.value)

  //初始化UI
  platformUi.value = ui;

  //初始化配置
  global.initData(conf);

  //初始化字体
  initFont()

  //初始化主题颜色
  initThemeColor()

  //初始化内容区域
  changeContentBackground("")

  //通知rom刷新列表
  global.incRomState();

  //加载平台
  initPlatform()

  //读取模拟器列表，用于ROM右键菜单
  GetAllSimulator().then((result: string) => {
    let resp = decodeApiData(result)
    simulatorMap.value = resp.data
  })

  //添加键盘事件
  addKeyboardEvent()

  //添加手柄事件
  initEventGamePad()

  //窗口缩放事件
  window.addEventListener('resize', windowResize);

})

onBeforeUnmount(() => {
  // 解绑键盘事件
  console.log("onBeforeMount tiny")
  removeKeyboardEvent()
  window.removeEventListener('resize', windowResize);
})


//窗口缩放
const windowResize = debounce(() => {
  UpdateOneConfig("WindowWidth", window.innerWidth.toString())
  UpdateOneConfig("WindowHeight", window.innerHeight.toString())
}, 500);


//初始化字体
const initFont = () => {
  let ui = platformUi.value
  if (ui.Font.Type == 1) {
    //系统字体
    currentFont.value = ui.Font.Family
  } else if (ui.Font.Type == 2) {
    //用户字体
    loadFont(ui.Font.Family, ui.Font.Format, ui.Font.Src).then((res) => {
      currentFont.value = ui.Font.Family
    }).catch((err) => {
      notify("err", ui.Font.Family + lang.value.fontLoadErr)
    })
  }
};

</script>

<style scoped>
@import "src/css/tiny/layout.css";
@import "src/css/transitions.css";
</style>
