<template>

  <div class="wrapper">
    <q-banner class="bg-primary text-white">
      {{ lang.tipRomfile }}
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
          <q-select dense filled square emit-value map-options v-if="activePlatform" :options="menuOptions"
                    v-model="activeMenu" class="mt" popup-content-class="q-select-content-min"
                    @update:model-value="changeMenu()"
                    :label="lang.selectMenu"/>
        </q-item-section>
      </q-item>
    </q-list>

    <q-table
        v-if="romList && !isEmpty(listColumns)"
        flat bordered dense separator="cell" row-key="Id" selection="multiple" :pagination="initialPagination"
        :rows-per-page-options="[20,50,100,200,500,900,0]"
        :selected-rows-label="Id" v-model:selected="selected" :rows="romList" :columns="listColumns">
      <template v-slot:top="props">
        <div class="col-2 q-table__title">{{ lang.romList }}</div>
        <q-space/>
        <q-select filled square dense map-options emit-value :options="simulatorList" v-model="activeSimId"
                  style="width: 150px" :label="lang.selectSim" @update:model-value="batchChangeSim()"/>
      </template>

      <template v-slot:body-cell-Simulator="props">
        <q-td key="Simulator" :props="props">
          <q-select filled square dense map-options emit-value :options="simulatorList" v-model="props.row.SimId"
                    @update:model-value="changeSim(props.row.Id,props.row.SimId)"/>
        </q-td>
      </template>

    </q-table>
  </div>

  <!-- 目录选择框 -->
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 250px">
      <q-card-section>
        <div class="text-h6">{{ currentOpt == 'linkCopy' ? lang.copyModule : lang.moveModule }}</div>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <q-list padding class="move-list">
          <q-item clickable dense active-class="menu-active" @click="activeMenuDialog = '/'"
                  :active="activeMenuDialog == '/'">
            <q-item-section side>
              <q-avatar icon="folder" color="primary" size="xs"/>
            </q-item-section>
            <q-item-section>
              <q-item-label lines="1">{{ lang.notCate }}</q-item-label>
            </q-item-section>
          </q-item>
          <div v-for="item in menuList">
            <q-item clickable dense active-class="menu-active" @click="activeMenuDialog = item.Path"
                    :active="activeMenuDialog == item.Path">
              <q-item-section side>
                <q-avatar icon="folder" color="primary" size="xs"/>
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ item.Name }}</q-item-label>
              </q-item-section>
            </q-item>
            <div v-if="item.SubMenu" v-for="sub in item.SubMenu">
              <q-item class="sub-tree" dense clickable active-class="menu-active" @click="activeMenuDialog = sub.Path"
                      :active="activeMenuDialog == sub.Path">
                <q-item-section side>
                  <q-avatar icon="folder" color="primary" size="xs"/>
                </q-item-section>
                <q-item-section>
                  <q-item-label lines="1">{{ sub.Name }}</q-item-label>
                </q-item-section>
              </q-item>
            </div>
          </div>
        </q-list>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue';

import {decodeApiData, isEmpty, notify} from "components/utils";
import {
  BatchSetRomSimId,
  GetAllSimulator,
  GetGameList,
  GetMenuList,
  GetPlatform,
  SetRomSimId
} from "app/wailsjs/go/controller/Controller";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {config, lang} = storeToRefs(global);
const menuOptions: any = ref([{value: "", label: lang.value.all}, {value: "/", label: lang.value.notCate}])
const menuList: any = ref(null)
const showDialog: any = ref(false)
const platformList: any = ref([])
const romList: any = ref(null)
const simulatorList: any = ref([])
const activePlatform: any = ref(null)
const activeMenu: any = ref("")
const activeMenuDialog: any = ref("")
const selected = ref([])
const activeSimId: any = ref(null)
const simulatorMap: any = ref(null)
const currentOpt = ref("")
const listColumns: any = ref([
  //{name: "Id", label: "ID", field: "Id", sortable: true},
  {name: "RomName", label: lang.value.filename, field: "RomName", sortable: true},
  {name: "Name", label: lang.value.alias, field: "Name", sortable: true},
  {name: "Menu", label: lang.value.belongMenu, field: "Menu", sortable: true},
  {name: "Simulator", label: lang.value.simulator, field: "Simulator", sortable: false},
])

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
})

//加载模拟器
function loadSimulators() {
  //读取全部模拟器
  GetAllSimulator().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      if (!simulatorMap.value) {
        simulatorMap.value = {}
      }
      let jsonArr = resp.data
      Object.keys(jsonArr).forEach(platform => {

        if (!simulatorMap.value[platform]) {
          simulatorMap.value[platform] = [{value: 0, label: lang.value.defaultSimulator}]
        }
        jsonArr[platform].forEach(sim => {
          simulatorMap.value[platform].push({value: sim.Id, label: sim.Name});
        });
      });
      simulatorList.value = simulatorMap.value[activePlatform.value]
    }
  })
}


//点击切换平台
function changePlatform() {

  //创建菜单
  GetMenuList(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      menuList.value = resp.data
      let jsonArr = resp.data
      jsonArr.forEach(item => {
        menuOptions.value.push({value: item.Path, label: item.Name});
      });
    }
  })


  if (simulatorMap.value == null) {
    loadSimulators()
  } else {
    //读取模拟器列表
    simulatorList.value = simulatorMap.value[activePlatform.value]
  }

  //读取rom列表
  getRomList()
}

//点击切换菜单
function changeMenu() {
  getRomList()
}

//读取rom列表
function getRomList() {
  selected.value = [];
  let isLike = activeMenu.value == "/" ? 0 : 1
  var req = {
    "platform": activePlatform.value,
    "catname": activeMenu.value,
    "showSubGame": 0,
    "page": -1, //加载全部数据
    "simpleModel": "simple",
    "catnameLike": isLike,
  };
  var request = JSON.stringify(req);
  GetGameList(request).then((result: string) => {
    let resp = decodeApiData(result)
    romList.value = resp.data;
  })
}

//更改单个模拟器
function changeSim(romId, simId) {
  SetRomSimId(romId, simId).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    notify("suc", lang.value.operSuc)
  })
}

//批量更改模拟器
function batchChangeSim() {
  let romIds: number[] = []
  selected.value.forEach((item, index) => {
    romIds.push(item.Id);
  });

  if (romIds.length == 0) {
    notify("err", lang.value.noSelectFile)
    return
  }

  BatchSetRomSimId(romIds, activeSimId.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }

    selected.value.forEach((item, index) => {
      item.SimId = activeSimId.value;
    });
    activeSimId.value = null

    notify("suc", lang.value.operSuc)
  })
}

</script>
<style scoped>

.wrapper {
  width: 60%;
  margin: 0 auto;
}

.move-list .q-item {
  padding: 0 10px;
}

.sub-tree {
  padding-left: 30px !important;
}

.menu-active {
  background: var(--q-primary);
}
</style>
