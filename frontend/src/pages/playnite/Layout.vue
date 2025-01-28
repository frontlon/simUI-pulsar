<template>
  <div class="border-out" @resize="windowResize">
    <div :class="'border-in bg fuzzy' + platformUi.BackgroundFuzzy" :style="{ fontFamily: currentFont }">
      <div class="bg-mask"></div>

      <header-bar-tool/>

      <div class="row">
        <!--        <div class="col" tabindex="0" style="width:8%;max-width:100px;"-->
        <!--             v-if="config.FrameShowDefault && config.FrameShowDefault.PlatformShow">-->
        <!--          <left-bar-platform/>-->
        <!--        </div>-->
        <!--        <div class="col" tabindex="1" style="width:9%;max-width: 130px"-->
        <!--             v-if="config.FrameShowDefault && config.FrameShowDefault.MenuShow">-->
        <!--          <left-bar-menu/>-->
        <!--        </div>-->
        <div class="col-7" style="position: relative">
          <rom-list-bar/>
        </div>
        <div class="col" style="min-width: 250px" v-if="config.FrameShowDefault && config.FrameShowDefault.RightShow">
          <content-bar/>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">

import {onBeforeUnmount, onMounted, ref} from "vue"
import {useGlobalStore} from 'src/stores/globalData'
import {callSrv, decodeApiData, loadFont, notify} from 'components/utils'
import {
  CheckUpgrade,
  GetAllSimulator,
  GetConfig,
  GetPlatformUi,
  UpdateOneConfig
} from 'app/wailsjs/go/controller/Controller'
import {storeToRefs} from "pinia";
// import  {initShortcuts, initThemeColor} from 'src/pages/playnite/HeaderBarTool.vue'
import {initPlatform} from 'src/pages/playnite/Platform.vue'
import {addKeyboardEvent, removeKeyboardEvent} from 'pages/playnite/modules/EventKeyboard.vue'
import {initEventGamePad} from 'pages/playnite/modules/EventGamepad.vue'
import RomListBar, {changeContentBackground} from 'src/pages/playnite/RomListBar.vue'
import ContentBar from 'src/pages/playnite/ContentBar.vue'
import HeaderBarTool, {initThemeColor} from 'src/pages/playnite/HeaderBarTool.vue'
import {debounce, Dialog} from "quasar";
import UpgradeComponent from "components/UpgradeComponent.vue";
import axios from "axios";

const global = useGlobalStore();
const {lang,config, theme, activeRom, platformUi, simulatorMap} = storeToRefs(global);
const currentFont = ref("Arial")

onMounted(async () => {

  theme.value = "Playnite"

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

  //检查更新
  checkUpgrade()
})


onBeforeUnmount(() => {
  // 解绑键盘事件
  console.log("onBeforeMount playnite")
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
@import "src/css/playnite/layout.css";
@import "src/css/transitions.css";

</style>
