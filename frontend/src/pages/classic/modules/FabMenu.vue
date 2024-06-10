<template>

  <!-- 添加菜单对话框 -->
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 600px">
      <q-card-section>
        <div class="text-h6">{{ lang.menuManage }}</div>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <q-list padding class="menu-tree-list">
          <q-item dense>
            <q-item-section side>
              <q-avatar icon="folder" color="primary" size="xs"/>
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ lang.rootMenu }}</q-item-label>
            </q-item-section>
            <q-item-section side>
              <div class="text-grey-8 q-gutter-xs">
                <q-btn size="sm" flat dense round icon="add" @click="manageAddMenu('/')"/>
              </div>
            </q-item-section>
          </q-item>
        </q-list>

        <q-list>
          <q-expansion-item expand-icon-toggle expand-separator :default-opened="item.SubMenu.length > 0"
                            aria-disabled="true"
                            v-for="(item,index) in menuList">
            <template v-slot:header>
              <q-item-section side>
                <q-avatar icon="folder" color="primary" size="xs"/>
              </q-item-section>
              <q-item-section>{{ item.Name }}</q-item-section>
              <q-item-section side>
                <div class="text-grey-8 q-gutter-xs">
                  <q-btn size="sm" flat dense round icon="add" @click="manageAddMenu(item.Path)"/>
                  <q-btn size="sm" flat dense round icon="edit" @click="manageEditMenu(item.Path,item.Name)"/>
                  <q-btn size="sm" flat dense round icon="delete" @click="manageDelMenu(item.Path)"/>

                  <q-btn size="sm" :disable="index == 0" flat dense round icon="keyboard_arrow_up"
                         @click="manageSortMenu('up',-1,index)"/>
                  <q-btn size="sm" :disable="index == menuList.length-1" flat dense round icon="keyboard_arrow_down"
                         @click="manageSortMenu('down',-1,index)"/>
                </div>
              </q-item-section>
            </template>
            <q-card v-if="item.SubMenu.length > 0">
              <q-card-section>
                <q-list>
                  <q-item v-for="(sub,sIndex) in item.SubMenu">
                    <q-item-section side>
                      <q-avatar icon="folder" color="positive" size="xs"/>
                    </q-item-section>
                    <q-item-section>
                      <q-item-label>{{ sub.Name }}</q-item-label>
                    </q-item-section>
                    <q-item-section top side>
                      <div class="text-grey-8 q-gutter-xs">
                        <q-btn size="sm" flat dense round icon="edit" @click="manageEditMenu(sub.Path,sub.Name)"/>
                        <q-btn size="sm" flat dense round icon="delete" @click="manageDelMenu(sub.Path)"/>
                        <q-btn size="sm" flat dense round icon="keyboard_arrow_up" :disable="sIndex == 0"
                               @click="manageSortMenu('up',index,sIndex)"/>
                        <q-btn size="sm" flat dense round icon="keyboard_arrow_down"
                               :disable="sIndex == item.SubMenu.length-1"
                               @click="manageSortMenu('down',index,sIndex)"/>
                      </div>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-card-section>
            </q-card>
            <q-card v-else class="empty">
              无子分类
            </q-card>
          </q-expansion-item>

        </q-list>

      </q-card-section>
    </q-card>
  </q-dialog>

</template>
<script lang="ts">

import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {ref} from "vue";
import {AddMenu, DeleteMenu, GetMenuList, RenameMenu, SortMenu} from "app/wailsjs/go/controller/Controller";
import {decodeApiData, notify} from "components/utils";
import {getPromptOpts} from "components/dialog";
import {Dialog} from "quasar";
import {callbackManageMenu} from 'pages/classic/LeftBarMenu.vue';

const global = useGlobalStore();
const {activePlatform, activeMenu, config, lang} = storeToRefs(global);
const menuList: any = ref([])
const showDialog: any = ref(false)

//打开菜单管理
export function openManageMenuDialog() {

  if (activePlatform.value < 1) {
    notify("err", "请选择一个平台")
    return;
  }

  showDialog.value = true;
  createMenuList();
}

function createMenuList() {
  GetMenuList(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    menuList.value = resp.data;
  })
}

//添加菜单
function manageAddMenu(path: string) {
  let opt = getPromptOpts(lang.value.createMenu, "", lang.value.ok, false, "")

  Dialog.create(opt).onOk(resp => {
    if (resp.input == "") {
      notify("err", lang.value.tipMenuNameIsNotEmpty)
      return
    }
    let data = resp.input

    AddMenu(activePlatform.value, path.toString(), data.toString()).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err != "") {
        notify("err", resp.err)
        return
      }
      menuList.value = resp.data;
      createMenuList();
      callbackManageMenu();
    })
  })
}

//编辑菜单
function manageEditMenu(path: string, name: string) {
  let opt = getPromptOpts(lang.value.renameMenu, "", lang.value.ok, false, name)

  Dialog.create(opt).onOk(resp => {
    if (resp.input == "") {
      notify("err", lang.value.tipMenuNameIsNotEmpty)
      return
    }
    if (resp.input == name) {
      return;
    }

    let data = resp.input
    RenameMenu(activePlatform.value, path, data).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err != "") {
        notify("err", resp.err)
        return;
      }
      menuList.value = resp.data;
      createMenuList();
      callbackManageMenu();
    })
  })
}

//删除菜单
function manageDelMenu(path: string) {
  let opt = getPromptOpts(lang.value.delMenu, lang.value.tipDelMenu, lang.value.ok, false)
  Dialog.create(opt).onOk(() => {
    DeleteMenu(activePlatform.value, path).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err != "") {
        notify("err", resp.err)
        return;
      }
      createMenuList();
      callbackManageMenu();
    })
  })
}

//菜单排序
function manageSortMenu(opt: string, pIndex: number, currIndex: number) {

  let req: any = [];
  let oldIndex = opt == "up" ? currIndex - 1 : currIndex + 1

  if (pIndex == -1) {
    //主目录
    let currentRow = menuList.value.splice(oldIndex, 1)[0];
    menuList.value.splice(currIndex, 0, currentRow);
    menuList.value.forEach((item: any, index: number) => {
      req.push(item.Path);
    });
  } else {
    //子目录
    let currentRow = menuList.value[pIndex].SubMenu.splice(oldIndex, 1)[0];
    menuList.value[pIndex].SubMenu.splice(currIndex, 0, currentRow);
    menuList.value[pIndex].SubMenu.forEach((item: any, index: number) => {
      req.push(item.Path);
    });
  }

  SortMenu(activePlatform.value, req).then((result: string) => {
    let resp = decodeApiData(result)
    console.log(resp.data)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    callbackManageMenu();
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
      manageAddMenu,
      manageEditMenu,
      manageDelMenu,
      manageSortMenu,
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
