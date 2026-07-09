<template>

  <div class="manage-wrapper system-config">
    <q-item class="page-header" style="--wails-draggable:drag">
      <q-item-section avatar>
        <q-btn flat size="lg" icon="arrow_back" class="vertical-middle" @click="$router.back()"/>
      </q-item-section>
      <q-item-section>
        <h3>{{ lang.systemConfig }}</h3>
      </q-item-section>
    </q-item>

    <div class="config-content">
      <q-tabs v-model="tab" class="config-tabs text-grey" active-color="primary" indicator-color="primary"
              align="left">
        <q-tab name="base" :label="lang.base"/>
        <q-tab name="shortcut" :label="lang.shortcut"/>
        <q-tab name="enum" :label="lang.enum"/>
      </q-tabs>

      <q-separator/>
      <q-tab-panels v-model="tab" keep-alive animated transition-prev="scale" transition-next="scale">
        <!-- 基本信息 -->
        <q-tab-panel name="base" class="config-panel" v-if="config">

          <q-card flat bordered class="config-section">
            <q-card-section>
              <div class="section-title">{{ lang.softName }}</div>
              <div class="row q-col-gutter-md">
                <div class="col-12 col-sm-6">
                  <q-input filled square :label="lang.softName" v-model="config.SoftName.Name"/>
                </div>
                <div class="col-12 col-sm-6">
                  <q-input filled square :label="lang.softSubName" v-model="config.SoftName.SubName"/>
                </div>
                <div class="col-12 col-sm-6">
                  <q-input filled square :label="lang.softNameImage" v-model="config.SoftName.Image">
                    <template v-slot:append>
                      <q-icon name="upload_file" class="open-dialog" @click="openFileDialog('image','SoftNameImage')"/>
                    </template>
                  </q-input>
                </div>
                <div class="col-12 col-sm-6">
                  <q-select filled square map-options emit-value :options="boolOptions" :label="lang.softNameHideText"
                            v-model="config.SoftName.HideText"/>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <q-card flat bordered class="config-section">
            <q-card-section>
              <div class="section-title">{{ lang.base }}</div>
              <div class="row q-col-gutter-md">
                <div class="col-12 col-sm-6">
                  <q-select filled square :options="langList" v-model="config.Lang" :label="lang.lang + ' Language'"/>
                </div>
                <div class="col-12 col-sm-6">
                  <q-select filled square map-options emit-value :options="boolOptions" :label="lang.checkUpgrade"
                            v-model="config.EnableUpgrade"/>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <q-card flat bordered class="config-section">
            <q-card-section>
              <div class="section-title">{{ lang.windowZoomDesc }}</div>
              <div class="row q-col-gutter-md">
                <div class="col-12 col-sm-4">
                  <q-select filled square map-options emit-value :options="zoomOptions" :label="lang.windowZoomDesc"
                            v-model="config.WindowZoom"/>
                </div>
                <div class="col-12 col-sm-4">
                  <q-select filled square map-options emit-value :options="boolOptions" :label="lang.adminRunGame"
                            v-model="config.AdminRunGame"/>
                </div>
                <div class="col-12 col-sm-4">
                  <q-select filled square map-options emit-value :options="boolOptions" :label="lang.gameMultiOpen"
                            v-model="config.GameMultiOpen"/>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <q-card flat bordered class="config-section">
            <q-card-section>
              <div class="section-title">{{ lang.splashScreenSrc }}</div>
              <div class="row q-col-gutter-md">
                <div class="col-12">
                  <q-input filled square :label="lang.splashScreenSrc" v-model="config.SplashScreen.Src">
                    <template v-slot:append>
                      <q-icon name="upload_file" class="open-dialog" @click="openFileDialog('image','SplashScreenSrc')"/>
                    </template>
                  </q-input>
                </div>
                <div class="col-12 col-sm-6">
                  <q-select filled square map-options emit-value :options="splashScreenSizeOptions"
                            :label="lang.splashScreenSize" v-model="config.SplashScreen.Size"/>
                </div>
                <div class="col-12 col-sm-6">
                  <q-select filled square map-options emit-value :options="splashScreenTimeOptions"
                            :label="lang.splashScreenTime" v-model="config.SplashScreen.Time"/>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <q-card flat bordered class="config-section">
            <q-card-section>
              <div class="section-title">{{ lang.baiduPicSearch }}</div>
              <div class="row q-col-gutter-md">
                <div class="col-12">
                  <q-input filled square :label="lang.baiduPicSearch" v-model="config.SearchEnginesBaidu"/>
                </div>
              </div>
            </q-card-section>
          </q-card>

          <div class="bottom-bar">
            <q-btn unelevated :label="lang.update" color="primary" class="update-btn" @click="updateConfig()"/>
          </div>
        </q-tab-panel>

        <!-- 快捷工具 -->
        <q-tab-panel name="shortcut" class="config-panel">
          <shortcut/>
        </q-tab-panel>

        <!-- 资料枚举 -->
        <q-tab-panel name="enum" class="config-panel">
          <rombase-enum/>
        </q-tab-panel>
      </q-tab-panels>
    </div>
  </div>

</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {decodeApiData, notify} from "components/utils";
import RombaseEnum from 'src/pages/config/RombaseEnum.vue'
import Shortcut from "pages/config/Shortcut.vue";
import {GetBaseConfig, OpenFileDialog, UpdateBaseConfig} from 'app/wailsjs/go/controller/Controller'

const global = useGlobalStore();
const tab = ref('base')
const config = ref(null)
const {langList, lang} = storeToRefs(global);

const boolOptions = [
  {label: '否', value: 0},
  {label: '是', value: 1},
]

const zoomOptions = [
  {label: '100%', value: 1},
  {label: '125%', value: 1.25},
  {label: '150%', value: 1.50},
  {label: '175%', value: 1.75},
]

const splashScreenSizeOptions = [
  {label: '25%', value: "25%"},
  {label: '50%', value: "50%"},
  {label: '75%', value: "75%"},
  {label: '100%', value: "100%"},
]

const splashScreenTimeOptions = [
  {label: lang.value.close, value: 0},
  {label: '1' + lang.value.second, value: 1},
  {label: '2' + lang.value.second, value: 2},
  {label: '3' + lang.value.second, value: 3},
  {label: '4' + lang.value.second, value: 4},
  {label: '5' + lang.value.second, value: 5},
]


onMounted(() => {
  GetBaseConfig().then((result: string) => {
    let resp = decodeApiData(result)
    config.value = resp.data;

    if (config.value.WindowZoom == 0) {
      config.value.WindowZoom = 1
    }
  })
})

function updateConfig() {
  UpdateBaseConfig(JSON.stringify(config.value)).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
    } else {
      notify("err", lang.value.updateFail + ":" + resp.err)
    }
  })
}

function openFileDialog(opt: string, type: string) {
  OpenFileDialog(opt).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {

      switch (type) {
        case "SoftNameImage":
          config.value.SoftName.Image = resp.data;
          break;
        case "SplashScreenSrc":
          config.value.SplashScreen.Src = resp.data
          break;
      }
    }
  })
}
</script>
<style scoped>
@import "src/css/manage.css";
@import "src/css/page/config.css";
</style>
