<template>
  <q-toolbar style="--wails-draggable:drag;" class="toolbar">
    <q-toolbar-title @dblclick="windowMaximise()">
      <text class="logo-text" v-if="config.SoftName">{{ config.SoftName.Name }}</text>

      <q-btn unelevated square size="sm" icon="format_list_bulleted">
        <q-menu>
          <q-list>
            <router-link to="/config">
              <q-item clickable>
                <q-item-section>{{ lang.systemConfig }}</q-item-section>
              </q-item>
            </router-link>
            <router-link to="/playnite/ui">
              <q-item clickable>
                <q-item-section>{{ lang.platformUI }}</q-item-section>
              </q-item>
            </router-link>
            <q-separator/>
            <q-item>
              <q-item-section>{{ lang.createCache }}</q-item-section>
              <q-item-section side>
                <q-icon name="keyboard_arrow_right"/>
              </q-item-section>
              <q-menu anchor="top end">
                <q-list>
                  <q-item clickable @click="createCache('platform')">
                    <q-item-section>{{ lang.createCachePlatform }}</q-item-section>
                  </q-item>
                  <q-item clickable @click="createCache('all')">
                    <q-item-section>{{ lang.createCacheAll }}</q-item-section>
                  </q-item>

                </q-list>
              </q-menu>
            </q-item>
            <q-separator/>
            <q-item clickable v-close-popup @click="upgrade()">
              <q-item-section>{{ lang.upgrade }}</q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="help()">
              <q-item-section>{{ lang.help }}</q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="about()">
              <q-item-section>{{ lang.about }}</q-item-section>
            </q-item>
          </q-list>

        </q-menu>
      </q-btn>
      <q-btn square size="sm" icon="checkroom">
        <q-tooltip>{{ lang.theme }}</q-tooltip>
        <q-menu>
          <q-list>
            <q-item clickable v-close-popup @click="gotoTheme('Default')">
              <q-item-section>Classic</q-item-section>
            </q-item>
            <q-item clickable v-close-popup @click="gotoTheme('Tiny')">
              <q-item-section>Tiny</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
      <q-btn unelevated square size="sm" label="P" @click="openPlatform()"/>
    </q-toolbar-title>


    <q-list class="titlebar-control row">
      <q-item></q-item>
      <q-item square clickable v-ripple @mousedown="windowMinimise()">
        <q-icon name="minimize"/>
      </q-item>
      <q-item square clickable v-ripple @mousedown="windowMaximise()">
        <q-icon name="crop_square"/>
      </q-item>
      <q-item square clickable v-ripple @mousedown="windowFullscreen()">
        <q-icon name="fullscreen"/>
      </q-item>
      <q-item square clickable v-ripple @mousedown="quit()">
        <q-icon name="close"/>
      </q-item>
    </q-list>
  </q-toolbar>


  <s-platform/>

</template>

<script lang="ts">

import UpgradeComponent from "src/components/UpgradeComponent.vue";
import SPlatform, {openPlatformDialog} from "src/pages/playnite/Platform.vue";
import AboutComponent from "src/components/AboutComponent.vue";
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
import {CheckUpgrade, CreateRomCache, UpdateOneConfig} from "app/wailsjs/go/controller/Controller";
import {ref} from "vue";

const global = useGlobalStore();
const {config, lang, activePlatform, activeFocus} = storeToRefs(global);
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

//更新
function upgrade() {

  CheckUpgrade().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", lang.value.upgradeCheckFail)
      return
    }

    let ver = resp.data
    if (ver.Upgrade == 0) {
      notify("suc", lang.value.upgradeLastest)
      return
    }

    //发现新版本
    Dialog.create({
      component: UpgradeComponent,
      componentProps: {
        data: ver,
        auto: 0,
      }
    })
  })

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

function openPlatform() {
  openPlatformDialog()
}

//打开帮助
function help() {
  BrowserOpenURL("http://www.simui.net/")
}

//关于窗口
function about() {
  Dialog.create({
    component: AboutComponent,
  })
}

export default {
  methods: {gotoTheme},
  components: {SPlatform},
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
      global,
      createCache,
      windowMinimise,
      windowMaximise,
      windowFullscreen,
      quit,
      switchFrameShow,
      upgrade,
      openPlatform,
      about,
      help,
    };
  }
}

</script>


<style scoped>

.logo-text {
  font-size: 14px;
  line-height: 1em;
  font-weight: bold;
  padding: 0 10px;
}

.q-toolbar {
  padding: 0;
  min-height: auto;
  height: auto;
  line-height: 1em;
  overflow: hidden;
}

.q-btn {
  margin: 0;
  padding: 0;
  width: 30px;
  box-shadow: none;
}

.q-item {
  font-size: 12px;
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
