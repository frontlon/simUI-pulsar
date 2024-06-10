<template>
  <!-- 导入游戏对话框 -->
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 600px;">
      <q-card-section>
        <div class="text-h6">{{ lang.inputGame }}</div>
      </q-card-section>
      <q-card-section class="q-pt-none">

        <q-tabs v-model="tab" class="text-grey" active-color="primary" indicator-color="primary" align="justify">
          <q-tab name="rom" :label="lang.simRom"/>
          <q-tab name="ps3" :label="lang.folderRom"/>
          <q-tab name="pc" :label="lang.pcGame"/>
          <q-tab name="share" :label="lang.inputShare"/>
        </q-tabs>

        <q-separator/>
        <q-tab-panels v-model="tab" keep-alive animated transition-prev="scale"
                      transition-next="scale">
          <!-- 模拟器ROM -->
          <q-tab-panel name="rom">

            <q-banner inline-actions rounded dense class="bg-primary text-white">
              {{ lang.tipInputSimRom }}
            </q-banner>

            <q-list class="output-list">
              <q-item>
                <q-item-section>
                  <q-input filled square dense :label="lang.selectRomFile" class="input" v-model="modelInputFile.rom">
                    <template v-slot:append>
                      <q-icon name="upload_file" class="open-dialog" @click="inputSelectFileMulti('rom')"/>
                    </template>
                  </q-input>
                </q-item-section>
              </q-item>

              <q-item v-if="romPathOptions.length > 1">
                <q-item-section>
                  <div class="list-title">{{ lang.inputWhereRomDir }}</div>
                  <q-option-group type="radio" :options="romPathOptions" v-model="modelRomPath"/>
                </q-item-section>
              </q-item>

            </q-list>
          </q-tab-panel>

          <!-- PS3 ROM(文件夹) -->
          <q-tab-panel name="ps3">
            <q-banner inline-actions rounded dense class="bg-primary text-white">
              {{ lang.tipInputPS3 }}
            </q-banner>
            <q-list class="output-list">
              <q-item>
                <q-item-section>
                  <q-input filled square dense :label="lang.selectFolder" class="input" v-model="modelInputFile.ps3">
                    <template v-slot:append>
                      <q-icon name="upload_file" class="open-dialog" @click="inputSelectFolder('ps3')"/>
                    </template>
                  </q-input>
                </q-item-section>
              </q-item>
            </q-list>
          </q-tab-panel>

          <!-- PC游戏 -->
          <q-tab-panel name="pc">
            <q-banner inline-actions rounded dense class="bg-primary text-white">
              {{ lang.tipInputPC }}
            </q-banner>
            <q-list class="output-list">
              <q-item>
                <q-item-section>
                  <q-input filled square dense :label="lang.selectBootFile" class="input" v-model="modelInputFile.pc">
                    <template v-slot:append>
                      <q-icon name="upload_file" class="open-dialog" @click="inputSelectFileSingle('pc')"/>
                    </template>
                  </q-input>
                  <q-input filled square dense :label="lang.runParam" class="input" v-model="modelInputParam.pc"/>
                </q-item-section>

              </q-item>
              <q-item>
                <q-item-section>
                  <q-toggle
                      v-model="addBatFile"
                      label="添加为Bat文件(解决多个同名pc游戏的问题)"
                  />
                </q-item-section>
              </q-item>
            </q-list>
          </q-tab-panel>

          <!-- 导入分享 -->
          <q-tab-panel name="share">
            <q-banner inline-actions rounded dense class="bg-primary text-white">
              {{ lang.tipInputShare }}
            </q-banner>
            <q-list class="output-list">
              <q-item>
                <q-item-section>
                  <q-input filled square dense :label="lang.selectShareFile" class="input"
                           v-model="modelInputFile.share">
                    <template v-slot:append>
                      <q-icon name="upload_file" class="open-dialog" @click="inputSelectFileSingle('share')"/>
                    </template>
                  </q-input>
                </q-item-section>
              </q-item>
              <q-item v-if="romPathOptions.length > 1">
                <q-item-section>
                  <div class="list-title">{{ lang.inputWhereRomDir }}</div>
                  <q-option-group type="radio" :options="romPathOptions" v-model="modelRomPath"/>
                </q-item-section>
              </q-item>
            </q-list>
          </q-tab-panel>
        </q-tab-panels>

        <div class="q-gutter-sm bottom-bar text-right">
          <q-btn flat square v-close-popup :label="lang.cancel" size="lg"/>
          <q-btn flat square :label="lang.add" size="lg" color="primary" @click="inputGame()"
                 :disable="modelInputFile[tab] == '' || modelRomPath==''"/>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>

</template>
<script lang="ts">

import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {
  AddGame,
  GetPlatformById,
  OpenDirectoryDialog,
  OpenFileDialog,
  OpenMultiFileDialog
} from "app/wailsjs/go/controller/Controller";
import {decodeApiData, isEmpty, notify} from "components/utils";
import {ref} from "vue";

const global = useGlobalStore();
const {activePlatform, activeMenu, config, callbackOpts, rootPath, lang} = storeToRefs(global);
const showDialog: any = ref(false)
const modelInputFile: any = ref({"rom": "", "ps3": "", "pc": "", "share": ""})
const modelInputParam: any = ref({"rom": "", "ps3": "", "pc": "", "share": ""})
const modelRomPath: any = ref("")
const romPathOptions: any = ref([])

const tab = ref('rom')
const addBatFile = ref(false)

let platform = null

//打开导入对话框
export function openInputDialog() {
  if (activePlatform.value < 1) {
    notify("err", lang.value.tipSelectAPlatform)
    return;
  }
  //读取平台信息
  GetPlatformById(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }

    platform = resp.data;

    if (platform.RomPath.length == 0) {
      notify("err", lang.value.romPathNotFound)
      return
    }

    //填充rom目录radio列表
    romPathOptions.value = [];
    platform.RomPath.forEach((item: any, index: number) => {
      let rel = item.replace(rootPath.value, "")
      romPathOptions.value.push({label: rel, value: item})
      if (index == 0) {
        modelRomPath.value = item;
      }
    })
    showDialog.value = true
  })
}

//选择单文件
function inputSelectFileSingle(opt: string) {
  OpenFileDialog(opt).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    if (isEmpty(resp.data)) {
      return
    }
    modelInputFile.value[opt] = resp.data
  })
}

//选择多文件
function inputSelectFileMulti(opt: string) {
  OpenMultiFileDialog(opt).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    if (isEmpty(resp.data)) {
      return
    }
    modelInputFile.value[opt] = resp.data.join(";")
  })
}

//选择文件夹
function inputSelectFolder(opt: string) {
  OpenDirectoryDialog().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    if (isEmpty(resp.data)) {
      return
    }
    modelInputFile.value[opt] = resp.data
  })
}

//导入游戏
function inputGame() {
  let files = modelInputFile.value[tab.value].split(";")
  let param = modelInputParam.value[tab.value]
  let menu = activeMenu.value == "" ? "/" : activeMenu.value;
  let addBat = addBatFile.value ? 1 : 0
  modelInputFile.value[tab.value] = ""
  modelInputParam.value[tab.value] = ""

  AddGame(tab.value, activePlatform.value, menu, modelRomPath.value, files, param, addBat).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }

    //清空文本框
    notify("suc", lang.value.inputSuc)

    callbackOpts.value = {
      "index": 0,
      "opt": "inputRom",
      "data": "",
    }
  })
}

export default {
  setup() {
    return {
      activePlatform,
      activeMenu,
      config,
      showDialog,
      tab,
      lang,
      modelRomPath,
      modelInputFile,
      modelInputParam,
      romPathOptions,
      addBatFile,
      inputSelectFileSingle,
      inputSelectFileMulti,
      inputSelectFolder,
      inputGame,
    };
  }
}

</script>

<style scoped>
.output-list .q-item {
  padding: 0;
}

.input {
  margin: 10px 0;
}

.list-title {
  font-size: 14px;
  margin: 10px 0;
}

</style>
