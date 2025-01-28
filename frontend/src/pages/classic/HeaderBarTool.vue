<template>
  <q-toolbar style="--wails-draggable:drag;background: var(--color-2);overflow: hidden">
    <q-toolbar-title @dblclick="windowMaximise()">
      <q-btn square size="sm" label="L" class="text-weight-bold"
             :class="{'bar-hide':config.FrameShowDefault && config.FrameShowDefault.LogoShow == false}"
             @click="switchFrameShow('LogoShow')">
        <q-tooltip>{{ lang.showHideLogo }}</q-tooltip>
      </q-btn>
      <q-btn square size="sm" label="P" class="text-weight-bold"
             :class="{'bar-hide':config.FrameShowDefault && config.FrameShowDefault.PlatformShow == false}"
             @click="switchFrameShow('PlatformShow')">
        <q-tooltip>{{ lang.showHidePlatform }}</q-tooltip>
      </q-btn>
      <q-btn square size="sm" label="M" class="text-weight-bold"
             :class="{'bar-hide':config.FrameShowDefault && config.FrameShowDefault.MenuShow == false}"
             @click="switchFrameShow('MenuShow')">
        <q-tooltip>{{ lang.showHideMenu }}</q-tooltip>
      </q-btn>
      <q-btn split square size="sm" label="R" class="text-weight-bold"
             :class="{'bar-hide':config.FrameShowDefault && config.FrameShowDefault.RightShow == false}"
             @click="switchFrameShow('RightShow')">
        <q-tooltip>{{ lang.showHideSidebar }}</q-tooltip>
      </q-btn>
      <router-link to="/classic/ui">
        <q-btn square size="sm" icon="auto_fix_high">
          <q-tooltip>{{ lang.platformUI }}</q-tooltip>
        </q-btn>
      </router-link>

      <q-btn square size="sm" icon="palette">
        <q-tooltip>{{ lang.themeColor }}</q-tooltip>
        <q-menu style="width: 200px">
          <q-list>
            <q-item>
              <q-item-section>
                <q-item-label>{{ lang.skinColor }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <div>
                  <q-btn square flat size="sm" icon="invert_colors" style="border: 0;background: none">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-color @change="changeColor" v-model="themeColor" no-header/>
                    </q-popup-proxy>
                  </q-btn>
                  <q-btn square flat size="sm" icon="title" style="border: 0;background: none">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-color @change="changeColor" v-model="themeTextColor" no-header/>
                    </q-popup-proxy>
                  </q-btn>
                </div>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label>{{ lang.brandColor }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <div>
                  <q-btn square flat size="sm" icon="invert_colors" style="border: 0;background: none">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-color @change="changeColor" v-model="themeBrandColor" no-header/>
                    </q-popup-proxy>
                  </q-btn>
                  <q-btn square flat size="sm" icon="title" style="border: 0;background: none">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-color @change="changeColor" v-model="themeBrandTextColor" no-header/>
                    </q-popup-proxy>
                  </q-btn>
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
      <q-btn square size="sm" icon="checkroom">
        <q-tooltip>{{ lang.theme }}</q-tooltip>
        <q-menu>
          <q-list>
            <q-item clickable v-close-popup @click="gotoTheme('Playnite')">
              <q-item-section>Playnite</q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="gotoTheme('Tiny')">
              <q-item-section>Tiny</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
      <q-btn square size="sm" icon="flash_on">
        <q-tooltip>{{ lang.shortcutTools }}</q-tooltip>
        <q-menu>
          <q-list>
            <q-item clickable v-close-popup v-for="(item, index) in shortCutList" @click="runShortcut(item.Path)">
              <q-item-section>{{ item.Name }}</q-item-section>
            </q-item>
          </q-list>
        </q-menu>

      </q-btn>
      <q-btn square size="sm" icon="refresh">
        <q-tooltip>{{ lang.createCache }}</q-tooltip>
        <q-menu>
          <q-list>
            <q-item clickable v-close-popup @click="createCache('platform')">
              <q-item-section>{{ lang.createCachePlatform }}</q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="createCache('all')">
              <q-item-section>{{ lang.createCacheAll }}</q-item-section>
            </q-item>
            <q-separator/>
            <q-item clickable>
              <q-item-section>{{ lang.clearData }}</q-item-section>
              <q-item-section side>
                <q-icon name="keyboard_arrow_right"/>
              </q-item-section>
              <q-menu anchor="top end" self="top start">
                <q-list>
                  <q-item clickable @click="clearGameStat()">
                    <q-item-section>{{ lang.clearGameStat }}</q-item-section>
                  </q-item>
                  <q-item clickable @click="clearNotExistConfig()">
                    <q-item-section>{{ lang.clearNotExistConfig }}</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>

      <q-btn square size="sm" icon="settings">
        <q-tooltip>{{ lang.config }}</q-tooltip>
        <q-menu>
          <q-list>
            <router-link to="/platform">
              <q-item clickable v-close-popup>
                <q-item-section>{{ lang.platformConfig }}</q-item-section>
              </q-item>
            </router-link>
            <router-link to="/config">
              <q-item clickable v-close-popup>
                <q-item-section>{{ lang.systemConfig }}</q-item-section>
              </q-item>
            </router-link>
            <router-link to="/romManage">
              <q-item clickable v-close-popup>
                <q-item-section>{{ lang.romManage }}</q-item-section>
              </q-item>
            </router-link>
            <q-separator/>

            <!--            <q-item clickable v-close-popup>
                          <q-item-section>{{lang.output}}</q-item-section>
                        </q-item>-->
<!--            <q-item clickable v-close-popup @click="copyright()">
              <q-item-section>{{ lang.copyright }}</q-item-section>
            </q-item>-->
            <q-item clickable v-close-popup @click="about()">
              <q-item-section>{{ lang.about }}</q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="help()">
              <q-item-section>{{ lang.help }}</q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="upgrade()">
              <q-item-section>{{ lang.upgrade }}</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </q-toolbar-title>

    <q-list class="titlebar-control row">
      <q-item></q-item>
      <q-item clickable v-ripple @mousedown="windowMinimise()">
        <q-icon name="minimize"/>
      </q-item>
      <q-item clickable v-ripple @mousedown="windowMaximise()">
        <q-icon name="crop_square"/>
      </q-item>
      <q-item clickable v-ripple @mousedown="windowFullscreen()">
        <q-icon name="fullscreen"/>
      </q-item>
      <q-item clickable v-ripple @mousedown="quit()">
        <q-icon name="close"/>
      </q-item>
    </q-list>
  </q-toolbar>

</template>

<script lang="ts">

import AboutComponent from "src/components/AboutComponent.vue";
import CopyrightComponent from "components/CopyrightComponent.vue";
import UpgradeComponent from "src/components/UpgradeComponent.vue";
import {Dialog, setCssVar} from "quasar";
import {
  BrowserOpenURL,
  Quit,
  WindowFullscreen,
  WindowIsFullscreen,
  WindowMinimise,
  WindowReload,
  WindowToggleMaximise,
  WindowUnfullscreen,
} from "app/wailsjs/runtime/runtime";
import {decodeApiData, getGradientColors, gotoTheme, loading, notify} from "components/utils";
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {
  CheckUpgrade,
  ClearGameStat,
  ClearNotExistGameConfig,
  CreateRomCache,
  GetShortcuts,
  RunProgram,
  UpdateColorsConfig,
  UpdateOneConfig
} from "app/wailsjs/go/controller/Controller";
import {ref, watch} from "vue";
import {getPromptOpts} from "components/dialog";
import axios from "axios";

const global = useGlobalStore();
const {config, lang, activePlatform} = storeToRefs(global);
const shortCutList = ref(null)
const themeColor = ref("rgb(0,0,0)");
const themeTextColor = ref("");
const themeBrandColor = ref("");
const themeBrandTextColor = ref("");
const root = document.documentElement;

//初始化主题颜色
export function initThemeColor() {

  themeColor.value = config.value.Colors.ThemeColor
  themeTextColor.value = config.value.Colors.ThemeTextColor
  themeBrandColor.value = config.value.Colors.ThemeBrandColor
  themeBrandTextColor.value = config.value.Colors.ThemeBrandTextColor

  const root = document.documentElement;

  //主题颜色
  if (config.value.Colors.ThemeColor != "") {
    let colors = getGradientColors(config.value.Colors.ThemeColor)
    colors.forEach((color, index) => {
      root.style.setProperty('--color-' + index, color);
    });
  }

  //主题字体颜色
  if (config.value.Colors.ThemeTextColor != "") {
    root.style.setProperty('--color-text', config.value.Colors.ThemeTextColor);
  }

  //品牌颜色
  if (config.value.Colors.ThemeBrandColor != "") {
    setCssVar('primary', config.value.Colors.ThemeBrandColor)
  }

  //品牌字体颜色
  if (config.value.Colors.ThemeBrandTextColor != "") {
    root.style.setProperty('--color-text-brand', config.value.Colors.ThemeBrandTextColor);
  }
}

//读取快捷工具
export function initShortcuts() {

  GetShortcuts(true).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      shortCutList.value = resp.data
    }

  })
}

//运行快捷工具
function runShortcut(path: string) {
  RunProgram(path).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return;
    }
  })
}

//监控主题颜色变更
watch([themeColor, themeTextColor, themeBrandColor, themeBrandTextColor], ([newThemeColor, newThemeTextColor, newThemeBrandColor, newThemeBrandTextColor], [oldThemeColor, oldThemeTextColor, oldThemeBrandColor, oldThemeBrandTextColor]) => {

  //主题颜色
  if (newThemeColor != oldThemeColor) {
    let colors = getGradientColors(newThemeColor)
    colors.forEach((color, index) => {
      root.style.setProperty('--color-' + index, color);
    });
  }

  //主题字体颜色
  if (newThemeTextColor != oldThemeTextColor) {
    root.style.setProperty('--color-text', newThemeTextColor);
  }

  //品牌颜色
  if (newThemeBrandColor != oldThemeBrandColor) {
    setCssVar('primary', newThemeBrandColor)
  }

  //品牌字体颜色
  if (newThemeBrandTextColor != oldThemeBrandTextColor) {
    root.style.setProperty('--color-text-brand', newThemeBrandTextColor);
  }

})


//改变品牌色
function changeColor() {

  config.value.Colors = {
    "ThemeColor": themeColor.value,
    "ThemeTextColor": themeTextColor.value,
    "ThemeBrandColor": themeBrandColor.value,
    "ThemeBrandTextColor": themeBrandTextColor.value,
  }

  var req = JSON.stringify(config.value.Colors);

  UpdateColorsConfig(req).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return;
    }
  })
}

//更新缓存
export function createCache(opt: string) {
  let platform = activePlatform.value
  if (opt == "all") {
    platform = 0;
  }

  loading("show", lang.createCacheing)

  CreateRomCache(platform).then((result: string) => {
    WindowReload()
  })
}

//关于窗口
function about() {
  Dialog.create({
    component: AboutComponent,
  })
}

//版权窗口
function copyright() {
  Dialog.create({
    component: CopyrightComponent,
  })
}


function upgrade() {
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
            notify("suc", lang.value.upgradeLastest)
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


//最小化
function windowMinimise() {
  WindowMinimise();
}

//最大化
function windowMaximise() {
  if (config.value.WindowState == 1) {
    config.value.WindowState = 0
  } else {
    config.value.WindowState = 1
  }
  UpdateOneConfig("WindowState", config.value.WindowState.toString())
  WindowToggleMaximise()
}

//全屏
function windowFullscreen() {
  const isFull = WindowIsFullscreen();
  isFull.then(value => {
    if (value) {
      config.value.WindowState = 0
      WindowUnfullscreen();
    } else {
      config.value.WindowState = 3
      WindowFullscreen();
    }
    UpdateOneConfig("WindowState", config.value.WindowState.toString())
  });
}

//关闭
function quit() {
  Quit();
}

//切换框架显示状态
function switchFrameShow(type: string) {
  let status = !config.value.FrameShowDefault[type];
  config.value.FrameShowDefault[type] = status
  UpdateOneConfig("FrameShowDefault", JSON.stringify(config.value.FrameShowDefault)).then((result: string) => {
  })

}

//打开帮助
function help() {
  BrowserOpenURL("http://www.simui.net/")
}

//清理游戏统计
function clearGameStat() {
  let opt = getPromptOpts(lang.value.clearGameStat, lang.value.tipClearGameStat, lang.value.ok, true)
  Dialog.create(opt).onOk(() => {
    ClearGameStat().then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err != "") {
        notify("err", resp.err)
        return
      }
      notify("suc", lang.value.operSuc)
    })
  })
}

//清理不存在的配置
function clearNotExistConfig() {
  let opt = getPromptOpts(lang.value.clearNotExistConfig, lang.value.tipClearNotExistConfig, lang.value.ok, true)
  Dialog.create(opt).onOk(() => {
    loading("show", lang.loading)
    ClearNotExistGameConfig().then((result: string) => {
      loading("hide")
      let resp = decodeApiData(result)
      if (resp.err != "") {
        notify("err", resp.err)
        return
      }
      notify("suc", lang.value.operSuc)
    })
  })
}


export default {
  methods: {gotoTheme},
  setup() {
    return {
      config,
      lang,
      activePlatform,
      shortCutList,
      themeColor,
      themeTextColor,
      themeBrandColor,
      themeBrandTextColor,
      changeColor,
      createCache,
      about,
      copyright,
      windowMinimise,
      windowMaximise,
      windowFullscreen,
      quit,
      switchFrameShow,
      runShortcut,
      help,
      upgrade,
      clearGameStat,
      clearNotExistConfig,
      global,
    };
  }
}

</script>


<style scoped>
@import "src/css/classic/common.css";

.q-toolbar {
  padding: 0;
  min-height: auto;
  background: var(--color-2);
  height: auto;
  line-height: 23px;
}

.q-btn {
  margin: 0;
  padding: 0;
  width: 30px;
  background: var(--color-2);
  box-shadow: none;
  border-left: 1px solid var(--color-0);
  border-right: 1px solid var(--color-4);
}

.q-btn:first-child {
  border-left: 1px solid var(--color-4);
}

.q-item {
}

.bar-hide {
  color: var(--color-5)
}

.titlebar-control {
}

.titlebar-control .q-item {
  width: 30px;
  height: 26px;
  min-height: 26px;
  line-height: 26px;
  padding: 0;
  text-align: center;
  border-left: 1px solid var(--color-0);
  border-right: 1px solid var(--color-4);
}

.titlebar-control .q-item .q-icon {
  margin: 0 auto;
  height: 100%;
}

.titlebar-control .q-item:first-child {
  border-left: 0;
}

.titlebar-control .q-item:last-child {
  border-right: 0;
}

</style>
