<template>

  <q-dialog v-model="showDialog" persistent transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 600px">
      <q-card-section>
        <div class="text-h6">编辑ROM模拟器设置</div>
      </q-card-section>
      <q-card-section class="q-pt-none" v-if="activeSimulatorInfo">
        <q-list>
          <q-item>
            <q-item-section>
              <q-input filled square dense :label="lang.runParam" v-model="activeSimulatorInfo.Cmd"/>
            </q-item-section>
            <q-item-section>
              <q-input filled square dense :label="lang.unzipRun" v-model="activeSimulatorInfo.Unzip">
                <q-tooltip>{{ lang.tipUnzipRun }}</q-tooltip>
              </q-input>
            </q-item-section>
          </q-item>
          <q-item>
            {{ lang.runtime }}
          </q-item>
          <q-item>
            <q-item-section>
              <q-input filled square dense :label="lang.runBefore" v-model="activeSimulatorInfo.RunBefore">
                <template v-slot:append>
                  <q-icon name="first_page" class="open-dialog" @click="openFileDialog('RunBefore')"/>
                </template>
              </q-input>
            </q-item-section>
            <q-item-section>
              <q-input filled square dense :label="lang.runAfter" v-model="activeSimulatorInfo.RunAfter">
                <template v-slot:append>
                  <q-icon name="last_page" class="open-dialog" @click="openFileDialog('RunAfter')"/>
                </template>
              </q-input>
            </q-item-section>
          </q-item>
        </q-list>

        <div class="q-gutter-sm bottom-bar text-right">
          <q-btn flat square v-close-popup :label="lang.cancel" size="lg"/>
          <q-btn flat square :label="lang.update" size="lg" color="primary" @click="updateRomSimulator()"/>
        </div>

      </q-card-section>
    </q-card>
  </q-dialog>

</template>
<script lang="ts">

import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {ref} from "vue";
import {GetRomSimSettingById, OpenFileDialog, UpdateRomSimSetting} from "app/wailsjs/go/controller/Controller";
import {decodeApiData, notify} from "components/utils";

const global = useGlobalStore();
const {activePlatform, activeMenu, config, lang} = storeToRefs(global);
const menuList: any = ref([])
const showDialog: any = ref(false)
const romId: any = ref(null)
const simId: any = ref(null)
const activeSimulatorInfo: any = ref(null)

//打开菜单管理
export function openEditRomSimulatorDialog(rid: number, sid: number) {
  romId.value = rid
  simId.value = sid
  showDialog.value = true;

  GetRomSimSettingById(rid, sid).then((result: string) => {
    let resp = decodeApiData(result)
    activeSimulatorInfo.value = resp.data;
  })

}

//更新数据
function updateRomSimulator() {
  let data = JSON.stringify(activeSimulatorInfo.value)
  UpdateRomSimSetting(romId.value, simId.value, data).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return;
    }
    notify("suc", lang.value.operSuc)
  })

}


//选择文件
function openFileDialog(dom: string) {
  OpenFileDialog("app").then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {
      activeSimulatorInfo.value[dom] = resp.data
    }
  })
}

export default {

  setup() {
    return {
      activePlatform,
      activeMenu,
      config,
      lang,
      menuList,
      showDialog,
      activeSimulatorInfo,
      updateRomSimulator,
      openFileDialog,
    };
  }
}

</script>

<style scoped>

.menu-tree-list .q-item {
  border-bottom: 1px solid var(--color-4);
}

.empty {
  text-align: center;
  font-size: 12px;
  color: var(--color-6);
  padding: 10px 0;
}
</style>