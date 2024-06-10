<template>

  <div>
    <q-banner class="bg-primary text-white">
      {{ lang.tipRomRepeatCheck1 }}
      <q-btn dense unelevated filled no-caps label="/cache/repeat" icon="open_in_new" color="positive"
             @click="openCacheDir('repeat')"/>
      {{ lang.tipRomRepeatCheck2 }}
    </q-banner>
    <q-list>
      <q-item>
        <q-item-section>
          <q-select dense filled square emit-value map-options label="选择平台" class="mt"
                    popup-content-class="q-select-content-min"
                    v-if="!isEmpty(platformList)" :options="platformList" v-model="activePlatform"/>
        </q-item-section>
        <q-item-section side>
          <q-btn filled square :label="lang.startCheck" class="mt" color="primary"
                 v-if="activePlatform !=null" @click="getRepeatList"/>
        </q-item-section>
      </q-item>
    </q-list>

    <q-table flat standout square dense row-key="Path" binary-state-sort separator="cell"
             :no-data-label="lang.unownedNotFound"
             v-if="repeatList" :pagination="initialPagination" :rows="repeatList" :columns="listColumns"
             :rows-per-page-options="[20,50,100,200,500,900,0]"
             selection="multiple" :selected-rows-label="Path" v-model:selected="selected">

      <template v-slot:top="props">
        <div class="col-2 q-table__title">{{ lang.unownedResultList }}</div>
        <q-btn size="md" color="primary" :label="lang.delSelectRom" :disable="isEmpty(selected)"
               @click="deleteUnownedFiles('selected')" style="margin-right: 10px"></q-btn>
      </template>

      <template v-slot:body-cell-Path="props">
        <q-td key="Path" :props="props">
          <span class="link" @click="openDir(props.row)">{{ props.row.Path }}</span>
        </q-td>
      </template>
      <template v-slot:body-cell-Size="props">
        <q-td key="Size" :props="props">
          {{ props.row.Size }}
        </q-td>
      </template>

      <template v-slot:body-cell-Opt="props">
        <q-td key="Opt" :props="props">
          <q-btn dense unelevated style="font-size: 10px" color="primary" :label="lang.run"/>
        </q-td>
      </template>

    </q-table>
  </div>

</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue';

import {decodeApiData, isEmpty, loading, notify} from "components/utils";
import {
  CheckRomRepeat,
  CreateRomCache,
  DeleteRepeatFile,
  GetPlatform,
  OpenCacheFolder,
  OpenFolderByPath,
} from "app/wailsjs/go/controller/Controller";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";
import {getPromptOpts} from "components/dialog";
import {Dialog} from "quasar";

const global = useGlobalStore();
const {config, lang} = storeToRefs(global);
const platformList: any = ref([])
const activePlatform: any = ref(null)
const repeatList: any = ref(null)
const selected = ref([])
const listColumns: any = ref([
  {name: "Path", label: "文件路径", field: "Path"},
  {name: "Size", label: "文件大小(byte)", field: "Size", sortable: true},
  {name: "Opt", label: "操作", field: "Path"},
])
let initialPagination = {descending: false, page: 0, rowsPerPage: 50}
onMounted(() => {
  //读取全部平台数据
  GetPlatform().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      let jsonArr = resp.data
      jsonArr.forEach(item => {
        platformList.value.push({value: item.Id, label: item.Name});
      });
    }
  })

})

//查找无效文件
function getRepeatList() {
  loading("show", lang.value.checking)
  CreateRomCache(activePlatform.value).then((result: string) => {
    CheckRomRepeat(activePlatform.value).then((result: string) => {
      let resp = decodeApiData(result)
      repeatList.value = resp.data;
      loading("hide")
    })
  })
}

//点击操作按钮
function deleteUnownedFiles() {
  let opt = getPromptOpts(lang.value.tipDelSelectFile, "", lang.value.ok, false)

  Dialog.create(opt).onOk(() => {
    let paths: string[] = []

    selected.value.forEach((item, index) => {
      paths.push(item.Path);
    });
    if (isEmpty(paths)) {
      notify("err", lang.value.noSelectPlatform)
      return
    }

    loading("show", lang.value.operating)
    DeleteRepeatFile(paths).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err != "") {
        notify("err", resp.err)
        return
      }
      notify("suc", lang.value.operSuc)
      loading("hide")
      getRepeatList()
    })
  })

}

//打开目录
function openDir(row: any) {
  OpenFolderByPath(row.Path).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
  })
}

//打开cache目录
function openCacheDir(type: string) {
  OpenCacheFolder(type, 1).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
  })
}


</script>
<style scoped>

.link:hover {
  text-decoration: underline;
  cursor: pointer;
}

</style>
