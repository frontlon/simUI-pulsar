<template>

  <div>
    <q-banner class="bg-primary text-white">
      {{ lang.tipUnownedClear1 }}
      <q-btn dense unelevated filled no-caps label="/cache/unowned" icon="open_in_new" color="positive"
             @click="openCacheDir('unowned')"/>
      {{ lang.tipUnownedClear2 }}
    </q-banner>
    <q-list>
      <q-item>
        <q-item-section>
          <q-select dense filled square emit-value map-options :label="lang.selectPlatform" class="mt"
                    v-if="!isEmpty(platformList)" :options="platformList" v-model="activePlatform"
                    popup-content-class="q-select-content-min"/>
        </q-item-section>
        <q-item-section side>
          <q-btn filled square :label="lang.startCheck" class="mt" color="primary"
                 v-if="activePlatform !=null" @click="getUnownedList"/>
        </q-item-section>
      </q-item>
    </q-list>

    <q-table flat standout square dense row-key="Id" binary-state-sort separator="cell" no-data-label="没有找到无效资源"
             v-if="unownedList" :pagination="initialPagination" :rows="unownedList" :columns="listColumns"
             :rows-per-page-options="[20,50,100,200,500,900,0]"
             selection="multiple" :selected-rows-label="Id" v-model:selected="selected">

      <template v-slot:top="props">
        <div class="col-2 q-table__title">{{ lang.unownedFileList }}</div>
        <q-btn size="md" color="primary" :label="lang.delSelectFile" :disable="isEmpty(selected)"
               @click="deleteUnownedFiles('selected')" style="margin-right: 10px"></q-btn>
        <q-btn size="md" color="negative" :label="lang.delAllUnownedFiles" :disable="isEmpty(unownedList)"
               @click="deleteUnownedFiles('all')"></q-btn>
      </template>

      <template v-slot:body-cell-Path="props">
        <q-td key="Path" :props="props">
          <span class="link" @click="openDir(props.row)">{{ props.row.Path }}</span>
        </q-td>
      </template>
      <template v-slot:body-cell-ResName="props">
        <q-td key="ResName" :props="props">
          {{ lang[props.row.ResName] ? lang[props.row.ResName] : props.row.ResName }}
        </q-td>
      </template>

      <template v-slot:body-cell-Opt="props">
        <q-td key="Opt" :props="props">
          <q-btn dense unelevated style="font-size: 10px" color="primary" label="查看"
                 @click="viewFile(props.row)"/>
        </q-td>
      </template>

    </q-table>
  </div>

</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue';

import {decodeApiData, isEmpty, loading, notify} from "components/utils";
import {
  CheckUnownedRes,
  DeleteUnownedFile,
  GetPlatform, OpenCacheFolder,
  OpenFolderByPath,
  OpenUnownedBakFolder,
  RunProgram
} from "app/wailsjs/go/controller/Controller";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";
import {getPromptOpts} from "components/dialog";
import {Dialog} from "quasar";

const global = useGlobalStore();
const {config, lang} = storeToRefs(global);
const platformList: any = ref([{value: 0, label: lang.value.all}])
const activePlatform: any = ref(null)
const unownedList: any = ref(null)
const selected = ref([])
const listColumns: any = ref([
  {name: "Path", label: lang.value.filePath, field: "Path", sortable: true},
  {name: "Platform", label: lang.value.belongPlatform, field: "Platform", sortable: true},
  {name: "ResName", label: lang.value.resType, field: "ResName", sortable: true},
  {name: "Opt", label: lang.value.operation, field: "Id"},
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
function getUnownedList() {
  loading("show", lang.value.checking)
  CheckUnownedRes(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    unownedList.value = resp.data;
    loading("hide")
  })
}

//点击操作按钮
function deleteUnownedFiles(type: string) {
  let opt = {}
  if (type == "all") {
    opt = getPromptOpts(lang.value.tipDelAllUnownedFiles, "", lang.value.ok, false)
  } else {
    opt = getPromptOpts(lang.value.tipDelSelectFile, "", lang.value.ok, false)
  }
  Dialog.create(opt).onOk(() => {

    let paths: string[] = []

    if (type == "all") {
      opt = getPromptOpts(lang.value.tipDelSelectFile, "", lang.value.ok, false)
    } else {
      selected.value.forEach((item, index) => {
        paths.push(item.Path);
      });
      if (isEmpty(paths)) {
        notify("err", lang.value.noSelectFile)
        return
      }
    }

    loading("show", lang.value.operating)
    DeleteUnownedFile(activePlatform.value, paths).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err != "") {
        notify("err", resp.err)
        return
      }
      notify("suc", lang.value.operSuc)
      loading("hide")
      getUnownedList()
    })
  })

}


//查看文件
function viewFile(row: any) {
  RunProgram(row.Path).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
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
