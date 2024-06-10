<template>
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 800px">
      <q-card-section>
        <div class="text-h6">{{ lang.outputShare }}</div>
      </q-card-section>
      <q-card-section class="q-pt-none">

        <div class="row">
          <div class="col-3" style="margin-right: 10px">
            <q-select dense filled square emit-value map-options :label="lang.selectPlatform" class="mt"
                      v-if="!isEmpty(platformList)" :options="platformList" v-model="activePlatform"/>
          </div>
          <div class="col">
            <q-input filled square dense :label="lang.outputDir" v-model="outputPath" @click="openSaveDialog()">
              <template v-slot:append>
                <q-icon style="cursor: pointer" name="output" class="open-dialog" @click="openSaveDialog()"/>
              </template>
            </q-input>
          </div>
        </div>

        <q-list>
          <q-item disable>
            <q-item-section avatar top>
              <q-checkbox v-model="checkRom" disable/>
            </q-item-section>
            <q-item-section>
              <q-item-label>ROM和子ROM</q-item-label>
              <q-item-label caption>
                导出平台下的所有ROM
              </q-item-label>
            </q-item-section>
          </q-item>

          <q-item disable>
            <q-item-section avatar top>
              <q-checkbox v-model="checkRes" disable/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ lang.gameRes }}</q-item-label>
              <q-item-label caption>
                {{ lang.tipGameRes }}
              </q-item-label>
            </q-item-section>
          </q-item>
          <q-item disable>
            <q-item-section avatar top>
              <q-checkbox v-model="checkRombase" disable/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ lang.baseInfo }}</q-item-label>
              <q-item-label caption>
                {{ lang.tipOutputRombase }}
              </q-item-label>
            </q-item-section>
          </q-item>
          <q-item tag="label">
            <q-item-section avatar top>
              <q-checkbox v-model="checkSimulator"/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ lang.simulator }}</q-item-label>
              <q-item-label caption>
                {{ lang.tipOutputSimulator }}
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>

      </q-card-section>
      <q-card-actions align="right" class="text-primary">
        <q-btn flat :label="lang.cancel" v-close-popup/>
        <q-btn flat :label="lang.output" v-close-popup @click="outputRom()" :disable="outputPath == ''"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script lang="ts">

import {ref} from 'vue'
import {decodeApiData, isEmpty, notify} from "components/utils";
import {storeToRefs} from "pinia";
import {useGlobalStore} from "stores/globalData";
import {GetPlatform, OutputRomByPlatform, SaveFileDialog} from "app/wailsjs/go/controller/Controller";

const global = useGlobalStore();
const {lang, activeRom, platformUi, config} = storeToRefs(global);
const showDialog: any = ref(false)
const romId: any = ref(null)
const romInfo: any = ref(null)
const outputPath: any = ref("")
const checkRom: any = ref(true)
const checkRes: any = ref(true)
const checkRombase: any = ref(true)
const checkSimulator: any = ref(false)

const platformList: any = ref([])
const activePlatform = ref(null)
const menuList: any = ref(null)

//导出分享ROM
export function openOutputDialog() {
  outputPath.value = ""
  showDialog.value = true

  //读取全部平台数据
  GetPlatform().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      platformList.value = []
      let jsonArr = resp.data
      jsonArr.forEach(item => {
        platformList.value.push({value: item.Id, label: item.Name});
      });
    }
  })
}

function changePlatform() {
  //读取菜单
}

//文件导出选择
function openSaveDialog() {
  let file = romInfo.value.RomName + ".shr"
  SaveFileDialog(file).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return;
    }
    outputPath.value = resp.data;
  })
}

//导出ROM
function outputRom() {
  OutputRomByPlatform(activePlatform.value, outputPath.value, checkSimulator.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return;
    }
    notify("suc", lang.outputSuc)

  })
}

export default {
  setup() {
    return {
      showDialog,
      activeRom,
      platformUi,
      config,
      outputPath,
      checkRom,
      checkRes,
      checkRombase,
      checkSimulator,
      lang,
      platformList,
      activePlatform,
      menuList,
      changePlatform,
      isEmpty,
      outputRom,
      openSaveDialog,
    }
  }
}

</script>

<style scoped>

.menu-table {
  margin-top: 10px;
}

.menu-list .q-item {
  padding: 0 10px;
}

.sub-tree {
  padding-left: 30px !important;
}
</style>