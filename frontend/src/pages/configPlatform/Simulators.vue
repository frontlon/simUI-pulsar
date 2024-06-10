<template>
  <div class="row">
    <div class="col-3">
      <div class="sim-list-wrapper">
        <q-btn-group spread stretch>
          <q-btn dense unelevated :label="lang.addSim" @click="addSimulator()"/>
          <q-btn dense unelevated :label="lang.delSim" @click="deleteSimulator()"/>
        </q-btn-group>
        <q-list class="sim-list">
          <q-item v-for="(item,index) in simulatorList" clickable v-ripple active-class="sim-active"
                  :key="item.Id" :active="activeSimulatorId === item.Id" @click="changeSimulator(item.Id,index)">
            <q-item-section>{{ item.Name }}</q-item-section>
          </q-item>
        </q-list>
      </div>
    </div>
    <div class="col" v-if="activeSimulatorInfo">
      <h6>模拟器配置</h6>
      <q-list>
        <q-item>
          <q-item-section>
            <q-input filled square dense :label="lang.simName" style="width: 100%" v-model="activeSimulatorInfo.Name">
            </q-input>
          </q-item-section>
        </q-item>
        <q-item>
          <q-item-section>
            <q-input filled square dense :label="lang.simApp" style="width: 100%" v-model="activeSimulatorInfo.Path">
              <template v-slot:append>
                <q-icon name="launch" class="open-dialog" @click="openFileDialog('Path','app')"/>
              </template>
            </q-input>
          </q-item-section>
        </q-item>
        <q-item>
          <q-item-section>
            <q-input filled square dense :label="lang.runParam" v-model="activeSimulatorInfo.Cmd"/>
          </q-item-section>
          <q-item-section side>
            <q-btn square flat icon="data_object">
              <q-tooltip>{{ lang.tipRunParam }}</q-tooltip>
              <q-menu>
                <q-list dense>
                  <q-item clickable v-close-popup v-for="item in SIMULATOR_BOOT_PARAM_OPTIONS"
                          @click="changeCmd(item.value)">
                    <q-item-section class="cmd-label col-2">{{ item.label }}</q-item-section>
                    <q-item-section class="cmd-value col">{{ item.value }}</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </q-item-section>
        </q-item>
      </q-list>
      <h6>{{ lang.runtime }}</h6>
      <q-list>
        <q-item>
          <q-item-section>
            <q-input filled square dense :label="lang.runBefore" v-model="activeSimulatorInfo.RunBefore">
              <template v-slot:append>
                <q-icon name="first_page" class="open-dialog" @click="openFileDialog('RunBefore','app')"/>
              </template>
            </q-input>
          </q-item-section>
          <q-item-section>
            <q-input filled square dense :label="lang.runAfter" v-model="activeSimulatorInfo.RunAfter">
              <template v-slot:append>
                <q-icon name="last_page" class="open-dialog" @click="openFileDialog('RunAfter','app')"/>
              </template>
            </q-input>
          </q-item-section>
        </q-item>
      </q-list>
      <div class="q-gutter-sm bottom-bar">
        <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" @click="updateSimulator()"/>
      </div>
    </div>

    <div v-else class="col relative-position">
      <div class="absolute-full flex flex-center">{{ lang.selectASim }}</div>
    </div>
  </div>


</template>

<script setup lang="ts">

import {defineProps, onMounted, ref, watchEffect} from 'vue'
import {decodeApiData, deepClone, isEmpty, notify} from 'components/utils'
import {SIMULATOR_BOOT_PARAM_OPTIONS} from "boot/constant";
import {getPromptOpts} from "components/dialog";
import {Dialog} from "quasar";
import Sortable from 'sortablejs'
import {
  AddSimulator,
  DelSimulator,
  GetSimulatorByPlatform,
  OpenFileDialog,
  UpdateSimulator,
  UpdateSimulatorSort
} from 'app/wailsjs/go/controller/Controller'
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);
const props = defineProps(['platform'])
const activePlatformId = ref(0)
const activeSimulatorId = ref(0)
const activeSimulatorIndex = ref(0)
const activeSimulatorInfo: any = ref(null)
const simulatorList: any = ref([])

const boolOptions = [
  {label: lang.value.no, value: 0},
  {label: lang.value.yes, value: 1},
]

onMounted(() => {
  //模拟器拖拽排序
  simulatorDropSort()
})

watchEffect(() => {
  activePlatformId.value = props.platform
  GetSimulatorByPlatform(props.platform).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      simulatorList.value = resp.data
    }
  })

});

//点击模拟器
function changeSimulator(id: number, index: number) {
  activeSimulatorId.value = id
  activeSimulatorIndex.value = index
  activeSimulatorInfo.value = deepClone(simulatorList.value[index]);
}

//选择文件
function openFileDialog(dom: string, opt: string) {
  OpenFileDialog(opt).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {
      activeSimulatorInfo.value[dom] = resp.data
    }
  })
}

//更新模拟器信息
function updateSimulator() {
  console.log(activeSimulatorInfo.value)

  UpdateSimulator(JSON.stringify(activeSimulatorInfo.value)).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      activeSimulatorInfo.value = resp.data
      simulatorList.value[activeSimulatorIndex.value] = activeSimulatorInfo.value
      notify("suc", lang.value.updateSuc)
    }
  })
}

//添加模拟器
function addSimulator() {
  let opt = getPromptOpts(lang.value.tipSimName, "", lang.value.ok, false, "")
  Dialog.create(opt).onOk(resp => {
    if (resp.input == "") {
      return
    }
    console.log("addSimulator", resp.input)
    let data = resp.input
    AddSimulator(activePlatformId.value, data).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err == "") {
        let jsonObj = resp.data
        simulatorList.value.push(jsonObj)
      }
    })
  })
}

//删除模拟器
function deleteSimulator() {

  if (activeSimulatorId.value == 0) {
    notify("warn", lang.value.selectSim)
    return;
  }

  let opt = getPromptOpts(lang.value.tipDelSim, "", lang.value.ok, false)
  Dialog.create(opt).onOk(() => {

    DelSimulator(activeSimulatorId.value).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err == "") {
        simulatorList.value.splice(activeSimulatorIndex.value, 1)
        activeSimulatorInfo.value = null
        activeSimulatorId.value = 0
        activeSimulatorIndex.value = 0
        notify("suc", lang.value.delSuc)
      }
    })
  })
}

//改变模拟器参数
function changeCmd(cmd: string) {
  activeSimulatorInfo.value.Cmd = cmd
}

//模拟器拖拽排序
function simulatorDropSort() {
  const el = document.querySelectorAll('.sim-list')[0] as HTMLElement;
  Sortable.create(el, {
    animation: 150,
    sort: true,
    disabled: false,
    onEnd: (evt: any) => {
      let currentRow = simulatorList.value.splice(evt.oldIndex, 1)[0];
      simulatorList.value.splice(evt.newIndex, 0, currentRow);
      //更新数据
      let req: any = [];
      simulatorList.value.forEach((item, index) => {
        simulatorList.value[index].Sort = index + 1
        req.push(item.Id);
        if (!isEmpty(activeSimulatorInfo.value) && activeSimulatorInfo.value.Id == item.Id) {
          activeSimulatorInfo.Sort = item.Sort
          activeSimulatorIndex.value = index
        }
      });

      UpdateSimulatorSort(JSON.stringify(req)).then((result: string) => {
      })

    },
  });
}
</script>

<style scoped>
@import "src/css/manage.css";

.sim-list-wrapper {
  margin-top: 8px;
  background: var(--color-3);
  margin-right: 10px;
  min-height: 400px;
}

.sim-param-item {
  padding: 0 10px;
  margin: 0
}

.sim-list .q-item {
  padding-left: 10px;
}

.sim-active {
  color: var(--color-15);
  background: var(--q-primary);
}

.cmd-label {
  padding-left: 10px;
  font-size: 12px
}

.cmd-value {
  padding-right: 10px;
  opacity: 0.6;
  font-size: 12px
}

</style>
