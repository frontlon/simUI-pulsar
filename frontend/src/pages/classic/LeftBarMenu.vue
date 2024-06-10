<template>
  <!-- 菜单栏 -->
  <q-scroll-area class="menu-wrapper wrap">
    <q-item class="menu-title">
      <q-item-section>{{ lang.gameNum }}</q-item-section>
      <q-item-section side>
        <q-badge :label="romCount"/>
      </q-item-section>
    </q-item>

    <q-list>
      <div v-if="menuList.length > 0" v-for="item in menuList">

        <!-- 折叠菜单 -->
        <q-expansion-item v-if="item.SubMenu.length > 0" expand-separator :label="item.Name" class="item-font"
                          :id="'menuEle' + item.Index">
          <div class="menu-sub-item-wrapper">
            <q-item clickable v-for="sub in item.SubMenu" class="menu-sub-item" active-class="btn-primary"
                    :id="'menuEle' + sub.Index" :active="sub.Path === activeMenu  && menuLike == sub.IsLike"
                    @click="changeMenu(sub.Path,sub.Index,sub.IsLike)">
              <q-item-section>
                <q-item-label>{{ sub.Name }}</q-item-label>
              </q-item-section>
            </q-item>
          </div>
        </q-expansion-item>

        <!-- 按钮 -->
        <q-item v-else clickable class="menu-item" active-class="btn-primary" :id="'menuEle' + item.Index"
                :active="item.Path === activeMenu && menuLike == item.IsLike"
                @click="changeMenu(item.Path,item.Index,item.IsLike)">
          <q-item-section>
            <q-item-label>{{ item.Name }}</q-item-label>
          </q-item-section>
        </q-item>
      </div>
    </q-list>
  </q-scroll-area>

</template>

<script lang="ts">

import {ref, watch} from 'vue'
import {callSrv, decodeApiData, isEmpty, notify} from 'components/utils'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {GetMenuList, UpdateOneConfig} from "app/wailsjs/go/controller/Controller";


const global = useGlobalStore();
const {
  activeMenu,
  activeFocus,
  menuLike,
  activeLetter,
  config,
  lang,
  activePlatform,
  romCount,
  scrollAreaRef
} = storeToRefs(global);
const menuList = ref([])
const menuIndex = ref(0)
const romCountMap = ref({})
watch(activePlatform, (newValue, oldValue) => {
  createMenuList();
});

async function createMenuList() {
  menuList.value = []

  romCountMap.value = await callSrv("GetGameCountByMenu", activePlatform.value)

  GetMenuList(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    menuList.value = handleMenuData(resp.data);
  })
}

function handleMenuData(data: any[]) {

  let index = 0;
  let create = [
    {Path: "", Name: lang.value.all, IsLike: 0, Index: index++, SubMenu: []},
  ];


  if (romCountMap.value["/"] != undefined && romCountMap.value["/"] > 0) {
    create.push({Path: "/", Name: lang.value.notCate, IsLike: 0, Index: index++, SubMenu: []})
  }

  if (data.length == 0) {
    return create;
  }


  data.forEach(function (item: any) {

    let subMenu: any[] = [];
    let masterMenu: any = {Path: item.Path, Name: item.Name, IsLike: 0, Index: index++, SubMenu: []}

    if (item.SubMenu && item.SubMenu.length > 0) {
      subMenu.push({Path: item.Path, Name: lang.value.all, IsLike: 1, Index: index++, SubMenu: []})

      if (romCountMap.value[item.Path] != undefined && romCountMap.value[item.Path] > 0) {
        subMenu.push({Path: item.Path, Name: lang.value.notCate, IsLike: 0, Index: index++, SubMenu: []})
      }

      item.SubMenu.forEach(function (sub: any) {
        subMenu.push({Path: sub.Path, Name: sub.Name, IsLike: 0, Index: index++, SubMenu: []})
      })

    }
    masterMenu.SubMenu = subMenu
    create.push(masterMenu)
  })

  return create;
}

//点击菜单
function changeMenu(menu: string, index: number, isLike: number) {
  if (menu == activeMenu.value && menuLike.value == isLike) {
    return
  }
  activeFocus.value = [2, index]
  activeMenu.value = menu
  menuLike.value = isLike
  activeLetter.value = "ALL"
  scrollAreaRef.value.setScrollPosition("vertical", 0) //滚动条回到最顶端
  global.incRomState()
  let menuJson = JSON.stringify([menu, isLike])
  UpdateOneConfig("Menu", menuJson).then((result: string) => {
  })
}


//菜单管理回调
export function callbackManageMenu() {
  createMenuList();
}


//键盘方向键
export function menuEventKeyboard(direction: string) {

  let newIndex = 0;
  let focusType = 2; //2菜单

  //焦点不在ROM里
  if (activeFocus.value[0] != focusType) {
    activeFocus.value = [focusType, -1];
  }

  let currFocus = activeFocus.value[1];
  let ele;
  switch (direction) {
    case "ArrowUp":
      if (currFocus == 0) {
        newIndex = menuList.value.length - 1
        newIndex = newIndex < 0 ? 0 : newIndex
        ele = document.getElementById('menuEle' + newIndex);
      } else {
        newIndex = currFocus - 1
        newIndex = newIndex < 0 ? 0 : newIndex;
        ele = document.getElementById('menuEle' + newIndex);
      }

      break;
    case "ArrowDown":
      newIndex = currFocus + 1
      ele = document.getElementById('menuEle' + newIndex);
      if (isEmpty(ele)) {
        ele = document.getElementById('menuEle0');
        newIndex = 0
      }
      break;
  }
  if (isEmpty(ele)) {
    return
  }

  activeFocus.value = [focusType, newIndex]
  ele.focus()

}

//事件点击菜单
export function clickMenuEvent() {
  if (activeFocus.value[0] != 2) {
    return
  }
  let ele = document.getElementById('menuEle' + activeFocus.value[1]);
  ele?.click()
}

export default {
  setup() {
    return {
      menuList,
      menuIndex,
      activeMenu,
      menuLike,
      config,
      lang,
      activePlatform,
      romCount,
      changeMenu,
    };
  }
}
</script>

<style scoped>
@import "src/css/classic/common.css";
@import "src/css/classic/menuBar.css";
</style>