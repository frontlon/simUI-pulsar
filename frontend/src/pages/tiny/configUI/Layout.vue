<template>

  <div class="wrapper">
    <div>

      <q-item style="--wails-draggable:drag">
        <q-item-section avatar>
          <router-link to="/tiny">
            <q-btn flat size="lg" icon="arrow_back" class="vertical-middle"/>
          </router-link>
        </q-item-section>
        <q-item-section>
          <h3>{{ lang.uiConfig }}</h3>
        </q-item-section>
      </q-item>
    </div>

    <q-list v-if="platformList">
      <q-item>
        <q-item-section>
          <q-select filled square dense options-dense emit-value map-options class="mt" :label="lang.selectPlatform"
                    popup-content-class="q-select-content-min"
                    :options="platformList" v-model="activePlatform" @update:model-value="changePlatform()"/>
        </q-item-section>
        <q-item-section v-if="activeUi">
          <q-select filled square dense options-dense emit-value map-options :options="sortOptions"
                    v-model="activeUi.RomSort" :label="lang.sortMethod"/>
        </q-item-section>
      </q-item>
    </q-list>
    <div v-if="activeUi">

      <h6>{{ lang.displayConfig }}</h6>
      <q-list>
        <q-item>
          <q-item-section>
            <q-select filled square dense options-dense emit-value map-options :options="thumbOptions"
                      v-model="activeUi.BlockThumbType" :label="lang.blockThumbType"/>
          </q-item-section>
          <q-item-section>
            <q-select filled square dense options-dense emit-value map-options :options="titleOptions"
                      v-model="activeUi.NameType" :label="lang.showNameType"/>
          </q-item-section>
        </q-item>

        <q-item>
          <q-item-section>
            <q-select filled square dense options-dense emit-value map-options :options="fontsizeOptions"
                      v-model="activeUi.BaseFontsize"
                      :label="lang.baseFontsize"/>
          </q-item-section>
          <q-item-section>
            <q-select filled square dense options-dense emit-value map-options :options="boolOptions"
                      v-model="activeUi.HideCarousel" :label="lang.hideCarousel"/>
          </q-item-section>
        </q-item>

      </q-list>


      <h6>{{ lang.backgroundConfig }}</h6>
      <q-list>
        <q-item>
          <q-item-section>
            <q-input filled square dense :label="lang.backgroundImage" v-model="activeUi.BackgroundImage">
              <template v-slot:append>
                <q-icon name="upload_file" @click="openFileDialog('image','BackgroundImage')"/>
              </template>
            </q-input>
          </q-item-section>
          <q-item-section>
            <q-input filled square dense :label="lang.backgroundMask" v-model="activeUi.BackgroundMask">
              <template v-slot:append>
                <q-icon name="upload_file" @click="openFileDialog('image','BackgroundMask')"/>
              </template>
            </q-input>
          </q-item-section>
        </q-item>
        <q-item>
          <q-item-section>
            <q-select filled square dense options-dense emit-value map-options :options="backgroundRepeatOptions"
                      v-model="activeUi.BackgroundRepeat"
                      :label="lang.backgroundRepeat"/>
          </q-item-section>
          <q-item-section>
            <q-select filled square dense options-dense emit-value map-options :options="backgroundFuzzyOptions"
                      v-model="activeUi.BackgroundFuzzy"
                      :label="lang.backgroundFuzzy"/>
          </q-item-section>
        </q-item>
      </q-list>
    </div>

    <div class="q-gutter-sm bottom-bar">
      <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" @click="updatePlatfirmUi()"/>
    </div>

  </div>

</template>

<script setup lang="ts">

import {onMounted, ref} from 'vue'
import 'animate.css'
import {decodeApiData, notify, wailsPathDecode} from "components/utils";
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {
  GetConfig,
  GetPlatform,
  GetPlatformUi,
  OpenFileDialog,
  UpdatePlatformUi
} from 'app/wailsjs/go/controller/Controller'
import {
  backgroundFuzzyOptions,
  backgroundRepeatOptions,
  boolOptions,
  fontsizeOptions,
  sortOptions,
  thumbOptions,
  titleOptions,
} from 'src/pages/classic/configUI/UiConst.vue';

const global = useGlobalStore();
const {config, lang, rootPath, theme} = storeToRefs(global);
const cfg = ref(null);
const uiMap: any = ref([]);
const platformList: any = ref(null);
const activePlatform: any = ref(null)
const activeUi: any = ref(null)
onMounted(() => {

  GetConfig().then((result: string) => {
    let resp = decodeApiData(result)
    cfg.value = resp.data.Config
  })

  //加载平台
  GetPlatform().then((result: string) => {
    let resp = decodeApiData(result)

    if (resp.err == "") {
      platformList.value = [{value: 0, label: lang.value.uiConfigSelectAllPlatform}];

      let jsonArr = resp.data
      jsonArr.forEach(item => {
        uiMap.value[item.Id] = item.Ui;
        platformList.value.push({value: item.Id, label: item.Name});
      });
    }
  })
})

function changePlatform() {
  console.log("changePlatform", activePlatform.value, theme.value)
  GetPlatformUi(activePlatform.value, theme.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      activeUi.value = resp.data
      if (activeUi.value.BackgroundImage != "") {
        activeUi.value.BackgroundImage = wailsPathDecode(activeUi.value.BackgroundImage)
        activeUi.value.BackgroundImage = activeUi.value.BackgroundImage.replace(rootPath.value, "");
      }
      if (activeUi.value.BackgroundMask != "") {
        activeUi.value.BackgroundMask = wailsPathDecode(activeUi.value.BackgroundMask)
        activeUi.value.BackgroundMask = activeUi.value.BackgroundMask.replace(rootPath.value, "");
      }
    }
  })
}

function updatePlatfirmUi() {

  let data = JSON.stringify(activeUi.value)
  console.log("UpdatePlatformUi", activePlatform.value, theme.value, data)

  UpdatePlatformUi(activePlatform.value, theme.value, data).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
      uiMap.value[activePlatform.value] = activeUi.value
    } else {
      notify("err", resp.err)
    }
  })
}

//选择文件
function openFileDialog(opt: string, par: string) {
  OpenFileDialog(opt).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {
      activeUi.value[par] = resp.data
    }
  })
}

</script>
<style scoped>

.wrapper {
  max-width: 60%;
  margin: 0 auto;
}

.q-tab-panels {
  background: none;
  padding: 0;
}

.q-tab-panel, .q-item {
  padding: 0;
}

.bottom-bar {
  text-align: right;
  margin-top: 20px;
}

.q-item {
  margin: 8px 0
}

.mt {
  margin-top: 8px;
}

.click-style {
  height: 100px;
}

.active {
  background: var(--q-primary);
}

h6 {
  margin-top: 30px;
  margin-bottom: 0;
}

.animate-item {
  height: 70px;
  line-height: 70px;
  text-align: center;
  width: 100%;
  background: var(--color-3);
  margin: 5px
}

.list-style {
  width: 100%;
  margin: 5px;
}

.thumb-orders {
  background: #ccc;
  margin: 5px;
}

</style>
