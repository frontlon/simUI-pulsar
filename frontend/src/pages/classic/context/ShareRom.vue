<template>
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 500px">
      <q-card-section class="row items-center q-pb-none q-mb-md">
        <div class="text-h6">{{ lang.outputShare }}</div>
        <q-space/>
        <q-btn icon="close" flat round v-close-popup/>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <q-input filled square dense :label="lang.outputDir" v-model="outputPath" @click="openSaveDialog()">
          <template v-slot:append>
            <q-icon style="cursor: pointer" name="output" class="open-dialog" @click="openSaveDialog()"/>
          </template>
        </q-input>
      </q-card-section>
      <q-card-section tag="label">
        <q-list>
          <q-item disable>
            <q-item-section avatar top>
              <q-checkbox v-model="checkMaster" disable/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ lang.romFile }}</q-item-label>
              <q-item-label caption>
                {{ lang.romMasterFile }}
              </q-item-label>
            </q-item-section>
          </q-item>
          <q-item tag="label">
            <q-item-section avatar top>
              <q-checkbox v-model="checkSlave"/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ lang.romSubGameFile }}</q-item-label>
              <q-item-label caption>
                {{ lang.tipDelSubGame }}
              </q-item-label>
            </q-item-section>
          </q-item>
          <q-item tag="label">
            <q-item-section avatar top>
              <q-checkbox v-model="checkRes"/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ lang.gameRes }}</q-item-label>
              <q-item-label caption>
                {{ lang.tipGameRes }}
              </q-item-label>
            </q-item-section>
          </q-item>
          <q-item tag="label">
            <q-item-section avatar top>
              <q-checkbox v-model="checkRombase"/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ lang.baseInfo }}</q-item-label>
              <q-item-label caption>
                {{ lang.tipOutputRombase }}
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>


      <q-card-actions align="right" class="text-primary">
        <q-btn color="primary" square size="md" :label="lang.ok" :disable="outputPath == ''"
               v-close-popup class="update-btn" @click="outputRom()"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script lang="ts">

import {ref} from 'vue'
import {decodeApiData, isEmpty, notify} from "components/utils";
import {storeToRefs} from "pinia";
import {useGlobalStore} from "stores/globalData";
import {OutputOneRom, SaveFileDialog} from "app/wailsjs/go/controller/Controller";

const global = useGlobalStore();
const {activePlatform, lang, activeRom, platformUi, config} = storeToRefs(global);
const showDialog: any = ref(false)
const romId: any = ref(null)
const romInfo: any = ref(null)
const outputPath: any = ref("")
const checkMaster: any = ref(true)
const checkSlave: any = ref(true)
const checkRes: any = ref(true)
const checkRombase: any = ref(true)


//导出分享ROM
export function openShareRomDialog(id: number, detail: any) {
  outputPath.value = ""
  romId.value = id
  romInfo.value = detail
  showDialog.value = true
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

  OutputOneRom(romId.value, outputPath.value, checkSlave.value, checkRes.value, checkRombase.value).then((result: string) => {
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
      checkMaster,
      checkSlave,
      checkRes,
      checkRombase,
      lang,
      isEmpty,
      outputRom,
      openSaveDialog,
    }
  }
}

</script>

<style scoped>
</style>
