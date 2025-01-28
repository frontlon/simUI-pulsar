<template>

  <!-- 平台栏 -->
  <q-dialog maximized square v-model="showDialog" @contextmenu="closePlatformDialog" @hide="callbackDialog" transition-show="fade" transition-hide="fade">
    <q-card flat class="wrapper" @click="closePlatformDialog">
      <q-btn label="X" class="close-btn" @click="closePlatformDialog"></q-btn>
      <q-scroll-area class="platform-scroll-area" ref="platformScrollRef" v-if="platformList.length > 0">
        <div class="platform-list" id="platformListDom">
          <div class="platform-item" v-for="(item,index) in platformList" :id="'platform-ele-'+index"
               @click.stop="changePlatform(item.Id,index)"
               :class="activePlatform == item.Id ? 'platform-active animate__animated animate__faster animate__flipInY' : ''">
            <div class="platform-bg"></div>
            <div class="platform-content">
              <div :class="activePlatform == item.Id ? 'platform-ico-active' : ''">
                <q-img class="platform-logo" fit="contain" :src="item.Icon" style="width: 100%"/>
                <div class="platform-title">{{ item.Name }}</div>
              </div>
              <q-scroll-area class="menu-scroll-area">
                <q-list class="menu-list"
                        v-if="activePlatform === item.Id && menuList[activePlatform] && menuList[activePlatform].length > 0">
                  <q-item clickable :active="menu.Path === activeMenu" :id="'menu-ele-'+mIndex"
                          active-class="menu-active" @click.stop="changeMenu(menu)"
                          v-for="(menu,mIndex) in menuList[activePlatform]">
                    {{ menu.Name }}
                  </q-item>
                </q-list>
              </q-scroll-area>
            </div>
          </div>
        </div>
      </q-scroll-area>
    </q-card>
  </q-dialog>
</template>

<script lang="ts">
import 'animate.css'
import {ref} from 'vue'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {decodeApiData, notify} from 'components/utils';
import {
  GetGameCountByMenu,
  GetMenuList,
  GetPlatform,
  GetPlatformById,
  GetPlatformUi,
  UpdateOneConfig
} from "app/wailsjs/go/controller/Controller";
import {changeContentBackground} from 'src/pages/playnite/RomListBar.vue'

const global = useGlobalStore();
const {
  activePlatform, activeRom, activeMenu, activeFocus, menuLike,
  activeLetter, platformUi, rombaseAlias, config, lang,activePlatformDesc,
} = storeToRefs(global);

const platformList = ref([])
const menuList = ref([])
const root = document.documentElement;
const showDialog = ref(false)
const platformScrollRef = ref(null)
let tmpRomFocus = null
const romCountMap = ref({})

export function initPlatform() {
  //读取全部平台数据
  platformList.value = []
  GetPlatform().then((result: string) => {
    let resp = decodeApiData(result)
    platformList.value = resp.data;
    getPlatformDesc()
  })
}


export function openPlatformDialog() {
  showDialog.value = true
  tmpRomFocus = activeFocus.value
  getMenuList()

  //滚动到激活项目
  platformList.value.forEach((item: any, index: number) => {
    if (item.Id == activePlatform.value) {
      setTimeout(function () {
        if (menuList.value[activePlatform.value]) {
          menuList.value[activePlatform.value].forEach(function (item: any, mIndex: number) {
            if (item.Path == activeMenu.value) {
              activeFocus.value = [2, mIndex]
            }
          })
        }
        platformScroll(index, 0)
      }, 100);
    }
  })

}

//加载平台描述
export function getPlatformDesc() {
  //滚动条回到最顶端
  GetPlatformById(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      activePlatformDesc.value = resp.data.Desc;
    }
  })
}

export function getPlatformList() {
  return platformList.value
}

function platformScroll(index, speed = 300) {

  let windowWidth = window.innerWidth;
  let outEle = document.getElementById('platformListDom');

  if (outEle == null) {
    return
  }

  let outRect = outEle.getBoundingClientRect();

  if (outRect.width < windowWidth) {
    return
  }

  let ele = document.getElementById('platform-ele-' + index);
  let offset = ele.offsetLeft - windowWidth / 2 + 100
  platformScrollRef.value.setScrollPosition('horizontal', offset, speed)
}

export function getMenuListByPlatform(platform) {
  return menuList.value[platform]
}

//点击切换平台
function changePlatform(id: number, index: number) {
  console.log("changePlatform", id)
  activeFocus.value = [1, index]
  activePlatform.value = id
  activeRom.value = 0
  activeLetter.value = "ALL"
  activeMenu.value = ""
  menuLike.value = 0
  tmpRomFocus = null //记录当前rom所在位置
  changePlatformUi(id)
  getPlatformDesc()
  UpdateOneConfig("Platform", id.toString()).then((result: string) => {
  })

  getMenuList()

  global.incRomState()

  platformScroll(index)
}

//更新平台ui
function changePlatformUi(id: number) {
  console.log("changePlatformUi", id)
  //读取平台ui
  GetPlatformUi(id, "Playnite").then((result: string) => {
    let resp = decodeApiData(result)
    console.log("changePlatformUi", resp.data)

    platformUi.value = resp.data;

    //更换背景图
    changeContentBackground(platformUi.value.BackgroundImage)
  })
}


async function getMenuList() {
  if (menuList.value[activePlatform.value]) {
    return
  }

  GetMenuList(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }

    //读取rom数量
    GetGameCountByMenu(activePlatform.value).then((countRes: string) => {
      let countResp = decodeApiData(countRes)
      if (countResp.err != "") {
        notify("err", resp.err)
        return
      }
      romCountMap.value = countResp.data
      menuList.value[activePlatform.value] = handleMenuData(resp.data);
    })

  })


}

function handleMenuData(data: any[]) {

  let index = 0;
  let create = [
    {Path: "", Name: lang.value.all, IsLike: 0, Index: index++, SubMenu: []},
  ];

  if (data.length == 0) {
    return create;
  }

  if (romCountMap.value["/"] != undefined && romCountMap.value["/"] > 0) {
    create.push({Path: "/", Name: lang.value.notCate, IsLike: 0, Index: index++, SubMenu: []})
  }


  data.forEach(function (item: any) {
    let masterMenu: any = {Path: item.Path, Name: item.Name, IsLike: 0, Index: index++, SubMenu: []}
    create.push(masterMenu)
  })
  return create;
}

export function changeMenu(item: any) {
  console.log("changeMenu")
  showDialog.value = false

  if (item == null) {
    item = menuList.value[activePlatform.value][activeFocus.value[1]]
  }

  activeFocus.value = [2, item.Index]
  activeMenu.value = item.Path
  menuLike.value = item.IsLike
  activeLetter.value = "ALL"
  let menuJson = JSON.stringify([item.Path, item.IsLike])
  UpdateOneConfig("Menu", menuJson).then((result: string) => {
  })
  global.incRomState()
  tmpRomFocus = null
}

export function closePlatformDialog() {
  showDialog.value = false
}

function callbackDialog() {
  if (tmpRomFocus) {
    activeFocus.value = tmpRomFocus
  }
}

export default {
  setup() {
    return {
      activePlatform,
      activeRom,
      activeMenu,
      menuLike,
      menuList,
      activeLetter,
      platformUi,
      rombaseAlias,
      config,
      lang,
      platformList,
      showDialog,
      closePlatformDialog,
      callbackDialog,
      changePlatform,
      changeMenu,
      platformScrollRef,
    }
  },
}
</script>
<style scoped>
@import "src/css/playnite/platform.css";
@import "src/css/playnite/layout.css";

</style>
