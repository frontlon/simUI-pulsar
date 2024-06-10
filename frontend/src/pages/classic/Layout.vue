<template>
  <div class="border-out" @resize="windowResize">
    <div class="border-in">

      <div class="row">
        <div class="col-2" style="max-width:240px"
             v-if="config.FrameShowDefault && config.FrameShowDefault.LogoShow">
          <header-bar-logo :titleInfo="config.SoftName"/>
        </div>
        <div class="col">
          <header-bar-tool/>
          <header-bar-filter/>
        </div>
      </div>

      <div class="row">
        <div class="col" tabindex="0" style="width:8%;max-width:100px;"
             v-if="config.FrameShowDefault && config.FrameShowDefault.PlatformShow">
          <left-bar-platform/>
        </div>
        <div class="col" tabindex="1" style="width:9%;max-width: 130px"
             v-if="config.FrameShowDefault && config.FrameShowDefault.MenuShow">
          <left-bar-menu/>
        </div>
        <div class="col" tabindex="2" style="position: relative;overflow: hidden">
          <content-bar/>
        </div>
        <div class="col-2" tabindex="3" style="min-width: 250px"
             v-if="config.FrameShowDefault && config.FrameShowDefault.RightShow">
          <right-bar/>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">

import {onActivated, onBeforeUnmount} from "vue"
import {useGlobalStore} from 'src/stores/globalData'
import {callSrv, decodeApiData, firstLetterToLower} from 'components/utils'
import {
  CheckUpgrade,
  GetAllSimulator,
  GetConfig,
  GetPlatformUi, GetTheme,
  UpdateOneConfig
} from 'app/wailsjs/go/controller/Controller'
import {storeToRefs} from "pinia";
import HeaderBarTool, {initShortcuts, initThemeColor} from 'src/pages/classic/HeaderBarTool.vue'
import HeaderBarFilter from 'src/pages/classic/HeaderBarFilter.vue'
import HeaderBarLogo from 'src/pages/classic/HeaderBarLogo.vue'
import LeftBarPlatform, {initPlatform} from 'src/pages/classic/LeftBarPlatform.vue'
import LeftBarMenu from 'src/pages/classic/LeftBarMenu.vue'
import RightBar from 'src/pages/classic/RightBar.vue'
import ContentBar, {changeContentBackground} from 'src/pages/classic/ContentBar.vue'
import {addKeyboardEvent, removeKeyboardEvent} from 'pages/classic/modules/EventKeyboard.vue'
import {initEventGamePad} from 'pages/classic/modules/EventGamepad.vue'
import {debounce, Dialog} from "quasar";
import UpgradeComponent from "components/UpgradeComponent.vue";
import axios from 'axios';

const global = useGlobalStore();
const {config,theme, platformUi, simulatorMap} = storeToRefs(global);

theme.value = "Default"

//启动时检查主题，跳转到默认主题
GetTheme().then((result: string) => {
  let resp = decodeApiData(result)
  if (resp.err == "") {
    let the = resp.data
    if (the != "" && the != theme.value) {
      let path = "/" + firstLetterToLower(the)
      global.goto(path, {}, 1)
    }
  }
})

onActivated(async () => {

  let conf = await callSrv("GetConfig")

  let ui = await callSrv("GetPlatformUi", conf.Config.Platform, "Default")

  //初始化UI
  platformUi.value = ui;

  //初始化配置
  global.initData(conf);

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

  //初始化快捷工具
  initShortcuts()

  //初始化快捷键
  addKeyboardEvent()

  //初始化手柄
  initEventGamePad()

  //窗口缩放事件
  window.addEventListener('resize', windowResize);

  //检查更新
  checkUpgrade()

})

onBeforeUnmount(() => {
  // 解绑键盘事件
  console.log("onBeforeMount")
  removeKeyboardEvent()
  window.removeEventListener('resize', windowResize);

})


function checkUpgrade() {
  //禁用自动更新
  if (config.value.EnableUpgrade == 0) {
    return
  }
  //请求新版本检查
  axios.get('https://www.simui.net/checkVersion.html')
      .then(httpRes => {
        let data = httpRes.data
        CheckUpgrade(data.Version).then((result: string) => {
          let resp = decodeApiData(result)
          if (resp.data == 0) {
            return
          }
          //发现新版本
          Dialog.create({
            component: UpgradeComponent,
            componentProps: {
              data: data,
              auto: 1,
            }
          })
        })
      });
}

//窗口缩放
const windowResize = debounce(() => {
  UpdateOneConfig("WindowWidth", window.innerWidth.toString())
  UpdateOneConfig("WindowHeight", window.innerHeight.toString())
}, 500);

</script>

<style scoped>
@import "src/css/classic/common.css";
@import "src/css/classic/layout.css";
</style>
