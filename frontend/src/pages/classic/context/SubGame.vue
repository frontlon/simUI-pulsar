<template>
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 1050px">
      <q-card-section class="row items-center q-pb-none q-mb-md">
        <div class="text-h6">{{ lang.editSubGame }}</div>
        <q-space/>
        <q-btn icon="close" flat round v-close-popup/>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <div class="row">
          <!-- 左窗口 -->
          <div class="col q-mr-sm">
            <div class="text-subtitle2 text-weight-bold q-mb-sm">{{ lang.bindedSubGame }}</div>
            <q-scroll-area class="scroll-area">
              <q-markup-table dense bordered flat class="game-table">
                <thead>
                <tr>
                  <th class="text-left">{{ lang.gameName }}</th>
                  <th class="text-left">{{ lang.romName }}</th>
                  <th class="text-left"></th>
                </tr>
                </thead>
                <tbody>
                <tr v-if="bindList && bindList.length > 0" v-for="(item,index) in bindList">
                  <td class="text-left">{{ item.Name }}</td>
                  <td class="text-left">{{ item.RomName }}</td>
                  <td class="text-right">
                    <q-btn unelevated size="xs" class="opt-btn" label="解除" color="primary" @click="unbind(item.Id,index)"/>
                  </td>
                </tr>
                <tr v-else>
                  <td colspan="3" class="text-center">{{ lang.noData }}</td>
                </tr>
                </tbody>
              </q-markup-table>
            </q-scroll-area>
          </div>
          <!-- 右窗口 -->
          <div class="col">
            <div class="text-subtitle2 text-weight-bold q-ml-sm q-mb-sm">{{ lang.notBindSubGame }}</div>
            <div class="q-gutter-xs row filter-wrapper q-ml-sm">
              <q-select standout square :options="menuList" v-model="activeMenu" dense="dense">
                <template v-slot:option="scope">
                  <q-item v-bind="scope.itemProps" dense="dense" class="option">
                    <q-item-section>
                      <q-item-label>{{ scope.opt.label }}</q-item-label>
                    </q-item-section>
                  </q-item>
                </template>
                <template v-slot:no-option>
                  <q-item>
                    <q-item-section class="select-no-option">{{ lang.noCate }}</q-item-section>
                  </q-item>
                </template>
              </q-select>

              <q-input dense filled square rounded :label="lang.gameKeyword" v-model="activeKeyword"
                       class="q-ml-sm q-mb-sm filter-input"/>
            </div>
            <q-scroll-area class="scroll-area">
              <q-markup-table dense bordered flat class="game-table">
                <thead>
                <tr>
                  <th class="text-left">{{ lang.gameName }}</th>
                  <th class="text-left">{{ lang.romName }}</th>
                  <th class="text-left"></th>
                </tr>
                </thead>
                <tbody>
                <tr v-if="unboundList && unboundList.length > 0" v-for="(item,index) in unboundList">
                  <td class="text-left">{{ item.Name }}</td>
                  <td class="text-left">{{ item.RomName }}</td>
                  <td class="text-right">
                    <q-btn unelevated size="xs" class="opt-btn" :label="lang.bind" color="primary" @click="bind(item.Id,index)"/>
                  </td>
                </tr>
                <tr v-else>
                  <td colspan="3" class="text-center">{{ lang.noData }}</td>
                </tr>
                </tbody>
              </q-markup-table>
            </q-scroll-area>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script lang="ts">

import {ref, watch} from 'vue'
import {useRoute} from "vue-router";
import {
  BindSubGame,
  GetGameListNotSubGame,
  GetMenuList,
  GetSubGameList,
  UnBindSubGame
} from "app/wailsjs/go/controller/Controller";
import {callback} from 'pages/classic/context/Context.vue';
import {decodeApiData, deepClone, notify} from "components/utils";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);
const route = useRoute();
const tab = ref('base')
const romInfo: any = ref(null)
const romIndex: any = ref(null)
const showDialog: any = ref(false)
const menuList: any = ref([{label: '全部', value: ''}, {label: '未分类', value: '/'}])
const bindList: any = ref(null)
const unboundList: any = ref(null)
const activeMenu: any = ref({"label": "全部", "value": ""});
const activeKeyword: any = ref("");


export function openSubGameDialog(index: number, detail: any) {
  console.log("openSubGame", index, detail)
  showDialog.value = true
  romInfo.value = detail;
  romIndex.value = index;
  activeKeyword.value = "";

  //菜单列表
  GetMenuList(detail.Platform).then((result: string) => {
        let resp = decodeApiData(result)
        console.log(resp.data)
        if (resp.err != "") {
          notify("err", resp.err)
          return;
        }

        if (resp.data.length > 0) {
          resp.data.forEach((item: any) => {
            let r = {label: item.Name, value: item.Path}
            menuList.value.push(r)
            if (item.SubMenu && item.SubMenu.length > 0) {
              item.SubMenu.forEach((sub: any) => {
                let s = {label: sub.Name, value: sub.Path}
                menuList.value.push(s)
              })
            }
          })
        }
      }
  )

//已绑定的子游戏列表
  GetSubGameList(detail.Id).then((result: string) => {
    let resp = decodeApiData(result)
    console.log("GetSubGameList", resp.data, resp.err)
    bindList.value = resp.data
  })

  //未绑定的子游戏列表
  getUnboundList()

}

//改变菜单
function changeMenu(menu: string, like: number) {
  activeMenu.value = menu;
}

//绑定
function bind(id: number, index: number) {
  BindSubGame(romInfo.value.Id, id).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      let item = deepClone(unboundList.value[index])
      unboundList.value.splice(index, 1)
      bindList.value.push(item)
      callback("bindSubGame", 0, "")
    } else {
      notify("err", resp.err)
    }
  })
}

//解绑
function unbind(id: number, index: number) {
  UnBindSubGame(id).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      let item = deepClone(bindList.value[index])
      bindList.value.splice(index, 1)
      unboundList.value.push(item)
      callback("bindSubGame", 0, "")
    } else {
      notify("err", resp.err)
    }
  })
}

watch(activeMenu, (newValue, oldValue) => {
  getUnboundList()
});

watch(activeKeyword, (newValue, oldValue) => {
  getUnboundList()
});

//读取游戏列表
function getUnboundList() {

  let menu = ""
  let like = 0
  if (activeMenu.value.value != "") {
    menu = activeMenu.value.value
  }
  if (activeMenu.value.value == "") {
    like = 1
  }

  console.log(romInfo.value, menu, like, activeKeyword.value)

  var req = {
    "platform": romInfo.value.Platform,
    "catname": menu,
    "catnameLike": like,
    "keyword": activeKeyword.value,
    "page": 0,
    "simpleModel": "simple",
  };
  var request = JSON.stringify(req);
  unboundList.value = []
  GetGameListNotSubGame(request).then((result: string) => {
    let resp = decodeApiData(result)
    console.log(resp.data)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }

    resp.data.forEach((item: any, index: number) => {
      if (item.Id == romInfo.value.Id) {
        return
      }
      unboundList.value.push(item)
    })

  })
}

export default {
  setup() {
    return {
      route,
      tab,
      romInfo,
      romIndex,
      showDialog,
      menuList,
      bindList,
      unboundList,
      activeMenu,
      activeKeyword,
      lang,
      bind,
      unbind,
      changeMenu,
    };
  }
}

</script>

<style scoped>
@import "src/css/manage.css";

.option {
  padding: 5px 10px
}

.game-table {
  width: 100%;
}

.scroll-area {
  height: 450px;
  width: 100%;
  paading: 0;
  margin: 0;
}

.game-table td {
  word-break: break-all;
  word-wrap: break-word;
  white-space: pre-wrap;
}

.game-table .opt-btn{
  min-width: 30px;
  padding: 4px 0!important;

}

</style>
