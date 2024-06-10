<template>

  <div class="manage-wrapper">
    <div>
      <q-item style="--wails-draggable:drag">
        <q-item-section avatar>
            <q-btn flat size="lg" icon="arrow_back" class="vertical-middle" @click="$router.back()"/>
        </q-item-section>
        <q-item-section>
          <h3>{{lang.systemConfig}}</h3>
        </q-item-section>
      </q-item>
    </div>

    <q-tabs v-model="tab" class="text-grey" active-color="primary" indicator-color="primary" align="justify">
      <q-tab name="base" :label="lang.base"/>
      <q-tab name="shortcut" :label="lang.shortcut"/>
      <q-tab name="enum" :label="lang.enum"/>
    </q-tabs>

    <q-separator/>
    <q-tab-panels v-model="tab" keep-alive animated transition-prev="scale"
                  transition-next="scale">
      <!-- 基本信息 -->
      <q-tab-panel name="base" v-if="config">

        <q-list>
          <q-item>
            <q-item-section>
              <q-input filled square dense :label="lang.softName" v-model="config.SoftName.Name"/>
            </q-item-section>
            <q-item-section>
              <q-input filled square dense :label="lang.softSubName" v-model="config.SoftName.SubName"/>
            </q-item-section>
            <q-item-section>
              <q-input filled square dense :label="lang.softNameImage" v-model="config.SoftName.Image">
                <template v-slot:append>
                  <q-icon name="upload_file" class="open-dialog" @click="openFileDialog('image')"/>
                </template>
              </q-input>
            </q-item-section>
            <q-item-section>
              <q-select filled square dense map-options emit-value :options="boolOptions" :label="lang.softNameHideText"
                        v-model="config.SoftName.HideText"/>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-select filled square dense :options="langList" v-model="config.Lang" :label="lang.lang + ' Language'"/>
            </q-item-section>
            <q-item-section>
              <q-select filled square dense map-options emit-value :options="boolOptions" :label="lang.checkUpgrade"
                        v-model="config.EnableUpgrade"/>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-select filled square dense map-options emit-value :options="zoomOptions" :label="lang.windowZoomDesc"
                        v-model="config.WindowZoom">
              </q-select>

            </q-item-section>
            <q-item-section>
              <q-select filled square dense map-options emit-value :options="boolOptions" :label="lang.adminRunGame"
                        v-model="config.AdminRunGame"/>
            </q-item-section>
            <q-item-section>
              <q-select filled square dense map-options emit-value :options="boolOptions" :label="lang.gameMultiOpen"
                        v-model="config.GameMultiOpen"/>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-input filled square dense :label="lang.baiduPicSearch" v-model="config.SearchEnginesBaidu"/>
            </q-item-section>
          </q-item>
        </q-list>
        <div class="q-gutter-sm bottom-bar">
          <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" @click="updateConfig()"/>
        </div>
      </q-tab-panel>

      <!-- 快捷工具 -->
      <q-tab-panel name="shortcut">
        <shortcut/>
      </q-tab-panel>

      <!-- 资料枚举 -->
      <q-tab-panel name="enum">
        <rombase-enum/>
      </q-tab-panel>
    </q-tab-panels>

  </div>


</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {decodeApiData, notify} from "components/utils";
import RombaseEnum from 'src/pages/config/RombaseEnum.vue'
import Shortcut from "pages/config/Shortcut.vue";
import {GetBaseConfig, OpenFileDialog, UpdateBaseConfig, UpdateShortcut} from 'app/wailsjs/go/controller/Controller'

const global = useGlobalStore();
const tab = ref('base')
const config = ref(null)
const {langList,lang} = storeToRefs(global);

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

onMounted(() => {
  //加载配置
  GetBaseConfig().then((result: string) => {
    let resp = decodeApiData(result)
    config.value = resp.data;
    if (config.value.WindowZoom == 0){
      config.value.WindowZoom = 1
    }
  })

})

//更新系统配置
function updateConfig() {

  UpdateBaseConfig(JSON.stringify(config.value)).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
    } else {
      notify("err",  lang.value.updateFail + ":" + resp.err)
    }
  })
}

//更新快捷工具
function updateShortcut() {
  UpdateShortcut(JSON.stringify(shortcuts.value)).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
    } else {
      notify("err", lang.value.updateFail +  + resp.err)
    }
  })
}

//选择文件
function openFileDialog(opt: string) {
  OpenFileDialog(opt).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {
      config.value.SoftName.Image = resp.data
    }
  })
}
</script>
<style scoped>
@import "src/css/manage.css";

.base-list-wrapper .q-textarea {
  margin: 8px 4px;
}

.base-list-wrapper textarea {
  height: 300px;
}

</style>
