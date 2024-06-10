<template>
  <q-banner class="bg-primary text-white">
    {{ lang.tipRombaseManage }}
  </q-banner>
  <q-list>
    <q-item>
      <q-item-section>
        <q-select dense filled square emit-value map-options :label="lang.selectPlatform" class="mt"
                  v-if="!isEmpty(platformList)" :options="platformList" v-model="activePlatform"
                  popup-content-class="q-select-content-min"
                  @update:model-value="changePlatform()"/>
      </q-item-section>
      <q-item-section>
        <q-select dense filled square emit-value map-options v-if="activePlatform" :options="menuList"
                  v-model="activeMenu" class="mt" popup-content-class="q-select-content-min"
                  @update:model-value="changeMenu()"
                  :label="lang.selectMenu"/>
      </q-item-section>
      <q-item-section>
        <q-input dense filled square emit-value map-options v-if="activePlatform"
                 v-model="activeKeyword" class="mt"
                 @update:model-value="searchKeyword()"
                 :label="lang.gameKeyword"/>
      </q-item-section>
    </q-item>
  </q-list>
  <div v-if="rombaseList && !isEmpty(listColumns)">
    <q-table flat standout square dense row-key="Name" binary-state-sort separator="cell" :no-data-label="lang.noData"
             :rows-per-page-options="[20,50,100,200,500,900,0]"
             :pagination="initialPagination" :rows="rombaseList" :columns="listColumns">
      <template v-slot:body="props">
        <q-tr :props="props">
          <q-td key="Name" :props="props">
            {{ props.row.Name }}
            <q-popup-edit v-model="props.row.Name" v-slot="scope">
              <q-input v-model="props.row.Name" dense autofocus :label="lang.alias" @keyup.enter="scope.set"
                       :placeholder="props.row.RomName" :error-message="lang.tipAliasIsNotEmpty"
                       :error="props.row.Name== ''"
                       @change="rombaseAddUpdate(props.row.RomName,props.row)"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseNameEn" :props="props">
            {{ props.row.BaseNameEn }}
            <q-popup-edit v-model="props.row.BaseNameEn" v-slot="scope">
              <q-input v-model="props.row.BaseNameEn" dense autofocus :label="lang.baseNameEn"
                       :placeholder="props.row.RomName"
                       @keyup.enter="scope.set"
                       @change="rombaseAddUpdate(props.row.RomName,props.row)"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseNameJp" :props="props">
            {{ props.row.BaseNameJp }}
            <q-popup-edit v-model="props.row.BaseNameJp" v-slot="scope">
              <q-input v-model="props.row.BaseNameJp" dense autofocus :label="lang.baseNameJp" @keyup.enter="scope.set"
                       :placeholder="props.row.RomName" @change="rombaseAddUpdate(props.row.RomName,props.row)"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseType" :props="props">
            {{ props.row.BaseType }}
            <q-popup-edit v-model="props.row.BaseType" v-slot="scope">
              <q-select filled square dense options-dense use-input fill-input hide-selected autofocus
                        :label="lang.baseType" :placeholder="props.row.RomName"
                        :options="enumMap.type" v-model="props.row.BaseType" @keyup.enter="scope.set"
                        popup-content-class="q-select-content"
                        @input-value="(value) => {props.row.BaseType = value;rombaseAddUpdate(props.row.RomName,props.row)}"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseYear" :props="props">
            {{ props.row.BaseYear }}
            <q-popup-edit v-model="props.row.BaseYear" v-slot="scope">
              <q-select filled square dense options-dense use-input fill-input hide-selected autofocus
                        :label="lang.baseYear" :placeholder="props.row.RomName"
                        :options="enumMap.year" v-model="props.row.BaseYear" @keyup.enter="scope.set"
                        popup-content-class="q-select-content"
                        @input-value="(value) => {props.row.BaseYear = value;rombaseAddUpdate(props.row.RomName,props.row)}"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BasePublisher" :props="props">
            {{ props.row.BasePublisher }}
            <q-popup-edit v-model="props.row.BasePublisher" v-slot="scope">
              <q-select filled square dense options-dense use-input fill-input hide-selected autofocus
                        :label="lang.basePublisher" :placeholder="props.row.RomName"
                        :options="enumMap.publisher" v-model="props.row.BasePublisher" @keyup.enter="scope.set"
                        popup-content-class="q-select-content"
                        @input-value="(value) => {props.row.BasePublisher = value;rombaseAddUpdate(props.row.RomName,props.row)}"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseProducer" :props="props">
            {{ props.row.BaseProducer }}
            <q-popup-edit v-model="props.row.BaseProducer" v-slot="scope">
              <q-select filled square dense options-dense use-input fill-input hide-selected autofocus
                        :label="lang.baseProducer" :placeholder="props.row.RomName"
                        :options="enumMap.producer" v-model="props.row.BaseProducer" @keyup.enter="scope.set"
                        popup-content-class="q-select-content"
                        @input-value="(value) => {props.row.BaseProducer = value;rombaseAddUpdate(props.row.RomName,props.row)}"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseCountry" :props="props">
            {{ props.row.BaseCountry }}
            <q-popup-edit v-model="props.row.BaseCountry" v-slot="scope">
              <q-select filled square dense options-dense use-input fill-input hide-selected autofocus
                        :label="lang.baseCountry" :placeholder="props.row.RomName"
                        :options="enumMap.country" v-model="props.row.BaseCountry" @keyup.enter="scope.set"
                        popup-content-class="q-select-content"
                        @input-value="(value) => {props.row.BaseCountry = value;rombaseAddUpdate(props.row.RomName,props.row)}"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseTranslate" :props="props">
            {{ props.row.BaseTranslate }}
            <q-popup-edit v-model="props.row.BaseTranslate" v-slot="scope">
              <q-select filled square dense options-dense use-input fill-input hide-selected autofocus
                        :label="lang.baseTranslate" :placeholder="props.row.RomName"
                        :options="enumMap.translate" v-model="props.row.BaseTranslate" @keyup.enter="scope.set"
                        popup-content-class="q-select-content"
                        @input-value="(value) => {props.row.BaseTranslate = value;rombaseAddUpdate(props.row.RomName,props.row)}"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseVersion" :props="props">
            {{ props.row.BaseVersion }}
            <q-popup-edit v-model="props.row.BaseVersion" v-slot="scope">
              <q-select filled square dense options-dense use-input fill-input hide-selected autofocus
                        :label="lang.baseVersion" :placeholder="props.row.RomName"
                        :options="enumMap.version" v-model="props.row.BaseVersion" @keyup.enter="scope.set"
                        popup-content-class="q-select-content"
                        @input-value="(value) => {props.row.BaseVersion = value;rombaseAddUpdate(props.row.RomName,props.row)}"/>
            </q-popup-edit>
          </q-td>
          <q-td key="Score" :props="props">
            {{ props.row.Score }}
            <q-popup-edit v-model="props.row.Score" v-slot="scope">
              <q-input type="number" v-model="props.row.Score" dense autofocus @keyup.enter="scope.set"
                       :label="lang.rating"
                       :placeholder="props.row.RomName" @change="rombaseAddUpdate(props.row.RomName,props.row)"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseOtherA" :props="props">
            {{ props.row.BaseOtherA }}
            <q-popup-edit v-model="props.row.BaseOtherA" v-slot="scope">
              <q-input v-model="props.row.BaseOtherA" dense autofocus @keyup.enter="scope.set"
                       :label="rombaseAlias.OtherA"
                       :placeholder="props.row.RomName" @change="rombaseAddUpdate(props.row.RomName,props.row)"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseOtherB" :props="props">
            {{ props.row.BaseOtherB }}
            <q-popup-edit v-model="props.row.BaseOtherB" v-slot="scope">
              <q-input v-model="props.row.BaseOtherB" dense autofocus @keyup.enter="scope.set"
                       :label="rombaseAlias.OtherB"
                       :placeholder="props.row.RomName" @change="rombaseAddUpdate(props.row.RomName,props.row)"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseOtherC" :props="props">
            {{ props.row.BaseOtherC }}
            <q-popup-edit v-model="props.row.BaseOtherC" v-slot="scope">
              <q-input v-model="props.row.BaseOtherC" dense autofocus @keyup.enter="scope.set"
                       :label="rombaseAlias.OtherC"
                       :placeholder="props.row.RomName" @change="rombaseAddUpdate(props.row.RomName,props.row)"/>
            </q-popup-edit>
          </q-td>
          <q-td key="BaseOtherD" :props="props">
            {{ props.row.BaseOtherD }}
            <q-popup-edit v-model="props.row.BaseOtherD" v-slot="scope">
              <q-input v-model="props.row.BaseOtherD" dense autofocus @keyup.enter="scope.set"
                       :label="rombaseAlias.OtherD"
                       :placeholder="props.row.RomName" @change="rombaseAddUpdate(props.row.RomName,props.row)"/>
            </q-popup-edit>
          </q-td>
          <q-td key="Id" :props="props">
            <q-btn dense unelevated size="xs" style="font-size: 10px" color="primary" :label="lang.run"
                   @click="runGame(props.row.Id)"/>
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </div>
  <div class="q-gutter-sm bottom-bar text-right">
    <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" :disable="isEmpty(editData)"
           @click="rombaseUpdate()"/>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {decodeApiData, isEmpty, notify} from "components/utils";
import {
  BatchSetRomBase,
  GetGameList,
  GetMenuList,
  GetPlatform,
  GetRomBaseAlias,
  GetRomBaseEnum,
  RunGame
} from "app/wailsjs/go/controller/Controller";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);
const platformList: any = ref([])
const menuList: any = ref([{value: "", label: lang.value.all}, {value: "/", label: lang.value.notCate}])
const rombaseList: any = ref(null)
const rombaseAlias: any = ref({})
const activePlatform: any = ref(null)
const activeMenu: any = ref("")
const activeKeyword: any = ref("")
const enumMap: any = ref(null)
const editData: any = ref({});
const listColumns: any = ref([])
let initialPagination = {descending: false, page: 0, rowsPerPage: 20}

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

  //选项枚举
  GetRomBaseEnum().then((result: string) => {
    let resp = decodeApiData(result)
    enumMap.value = resp.data
  })
})

//点击切换平台
function changePlatform() {
  console.log("changePlatform")
  //清空修改数据
  editData.value = {}
  activeKeyword.value = ""
  //读取菜单
  GetMenuList(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      let jsonArr = resp.data
      jsonArr.forEach(item => {
        menuList.value.push({value: item.Path, label: item.Name});
      });
    }
  })

  //读取资料项别名
  GetRomBaseAlias(activePlatform.value).then((result: string) => {
    console.log("GetRomBaseAlias", result)
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", lang.value.tipGetRomBaseAlias + resp.err)
      return
    }
    rombaseAlias.value = resp.data
    createListColumns()
  })

  //读取rom列表
  getRomList()
}

//构建资料列表项
function createListColumns() {
  listColumns.value = [
    {name: "Name", label: lang.value.alias, field: "Name", sortable: true},
    {name: "BaseNameEn", label: lang.value.baseNameEn, field: "BaseNameEn", sortable: true},
    {name: "BaseNameJp", label: lang.value.baseNameJp, field: "BaseNameJp", sortable: true},
    {name: "BaseType", label: lang.value.baseType, field: "BaseType", sortable: true},
    {name: "BaseYear", label: lang.value.baseYear, field: "BaseYear", sortable: true},
    {name: "BasePublisher", label: lang.value.basePublisher, field: "BasePublisher", sortable: true},
    {name: "BaseProducer", label: lang.value.baseProducer, field: "BaseProducer", sortable: true},
    {name: "BaseCountry", label: lang.value.baseCountry, field: "BaseCountry", sortable: true},
    {name: "BaseTranslate", label: lang.value.baseTranslate, field: "BaseTranslate", sortable: true},
    {name: "BaseVersion", label: lang.value.baseVersion, field: "BaseVersion", sortable: true},
    {name: "Score", label: lang.value.rating, field: "Score", sortable: true},
    {name: "BaseOtherA", label: rombaseAlias.value["OtherA"], field: "BaseOtherA", sortable: true},
    {name: "BaseOtherB", label: rombaseAlias.value["OtherB"], field: "BaseOtherB", sortable: true},
    {name: "BaseOtherC", label: rombaseAlias.value["OtherC"], field: "BaseOtherC", sortable: true},
    {name: "BaseOtherD", label: rombaseAlias.value["OtherD"], field: "BaseOtherD", sortable: true},
    {name: "Id", label: lang.value.runGame, field: "Id"},
  ];
}

//点击切换菜单
function changeMenu() {
  getRomList()
}

//读取rom列表
function getRomList() {
  let isLike = activeMenu.value == "/" ? 0 : 1
  var req = {
    "platform": activePlatform.value,
    "catname": activeMenu.value,
    "showSubGame": 0,
    "page": -1, //加载全部数据
    "simpleModel": "simple",
    "catnameLike": isLike,
    "keyword": activeKeyword.value,
  };
  var request = JSON.stringify(req);
  GetGameList(request).then((result: string) => {
    let resp = decodeApiData(result)
    rombaseList.value = resp.data;
  })
}

//资料修改记录
function rombaseAddUpdate(romName: string, data: any) {

  if (data.Name == "") {
    notify("err", lang.value.tipAliasIsNotEmpty)
    return
  }

  editData.value[romName] = data
  let req = JSON.stringify(editData.value)
  console.log("handleUpdate", req)
}

//资料修改提交
function rombaseUpdate() {

  if (isEmpty(editData.value)) {
    notify("err", lang.value.noModify)
    return
  }
  //清空修改记录
  let req = JSON.stringify(editData.value)
  console.log("handleUpdate", req)
  editData.value = {}
  BatchSetRomBase(activePlatform.value, req).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    notify("suc", lang.value.operSuc)
  })

}

//关键字查询
function searchKeyword() {
  console.log(activeKeyword.value)
  getRomList()
}

//运行游戏
function runGame(id: number) {
  console.log("runGame", id)
  RunGame(id, 0).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
  })
}

</script>
<style scoped>
@import "src/css/romManage.css";
</style>