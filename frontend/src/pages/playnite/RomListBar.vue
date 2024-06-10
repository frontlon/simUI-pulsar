<template>
  <!-- ROM列表 -->
  <div class="rom-wrapper">
    <q-scroll-area ref="scrollAreaRef" class="rom-scroll" :visible="false">
      <q-infinite-scroll @load="nextPage" :offset="100" debounce="500" style="width:100%;max-width: 100%">

        <div class="rom-block">
          <div v-for="(item, index) in items" :style="'width:'+ 100 / platformUi.BlockSize + '%;'">
            <q-item clickable :class="'q-pa-xs q-ma-' + platformUi.BlockMargin"
                    :active-class="'active animate__animated animate__faster ' + platformUi.BlockClickAnimate"
                    :id="'rom-ele-'+ index" :active="activeRom === item.Id"
                    @click="clickGame(item.Id,index)">
              <q-card flat style="width: 100%" :style="platformUi.BlockHideBackground ?'background:none' : ''">
                <q-img v-if="item.ThumbPic" :src="item.ThumbPic" loading="lazy" fit="fill"
                       :class="'img-direction-'+platformUi.BlockDirection">
                  <template v-slot:error>
                    <div class="absolute-full flex flex-center image-error">{{ lang.thumbError }}</div>
                  </template>
                  <template v-slot:loading>
                    <div class="absolute-full flex flex-center image-error">{{ lang.loading }}</div>
                  </template>
                </q-img>
                <q-img v-else src="" :class="'img-direction-'+platformUi.BlockDirection">
                  <template v-slot:default>
                    <div class="absolute-full flex flex-center image-error">{{ lang.noThumb }}</div>
                  </template>
                </q-img>
                <q-card-section v-if="!platformUi.BlockHideTitle" class="module-title"
                                :style="{'background':platformUi.BlockHideBackground ?'none' : ''}">
                  {{ platformUi.NameType == 1 ? item.Name : item.RomName }}
                </q-card-section>
              </q-card>
            </q-item>
          </div>
        </div>


        <!--        <q-list class="rom-list">
                  <div v-for="(item, index) in items">
                    <q-item clickable dense
                            :active-class="'active animate__animated animate__faster ' + platformUi.BlockClickAnimate"
                            :id="'ele'+ index" :active="activeRom === item.Id"
                            @click="clickGame(item.Id,index)">
                      <q-item-section side v-if="item.ThumbPic">
                        <q-img :src="item.ThumbPic" loading="lazy" fit="fill" class="rom-img"/>
                      </q-item-section>
                      <q-item-section>
                        {{ platformUi.NameType == 1 ? item.Name : item.RomName }}
                      </q-item-section>
                    </q-item>
                  </div>
                </q-list>-->

      </q-infinite-scroll>

    </q-scroll-area>
  </div>


</template>
<script lang="ts">
import 'animate.css'
import {ref, watch} from 'vue';
import {decodeApiData, getTimestamp, isEmpty} from "components/utils";
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {GetGameCount, GetGameList} from "app/wailsjs/go/controller/Controller";
import {closePlatformDialog, getMenuListByPlatform, getPlatformList} from 'pages/playnite/Platform.vue'
import {closeContentDialog, openRunGameDialogByKeyboard} from 'pages/playnite/ContentBar.vue'

const global = useGlobalStore();
const {
  activeMenu, activePlatform, activeRom, activeFocus, romCount, config, lang, romState, activeKeyword,
  activeType, activeYear, activeProducer, activePublisher, activeCountry, activeTranslate, scrollAreaRef,
  activeVersion, activeScore, activeComplete, activeLetter, menuLike, platformUi, callbackOpts, theme
} = storeToRefs(global);

const items: any = ref([])
const rating = ref(2)
const listColumns: any = ref([])
const letterList: any = ref(["ALL", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "#"])
const root = document.documentElement;
let pageNum = 0; //翻页数
let pageEnd = false; //是否翻到最后一页
let initialPagination = {
  descending: false,
  page: 0,
  rowsPerPage: 0
}
const qItems = ref([]);
let clickExpire = 0

//监听rom列表刷新
watch([romState], (newValue, oldValue) => {
  pageEnd = false;

  //滚动条回到最顶端
  if (scrollAreaRef.value != null) {
    scrollAreaRef.value.setScrollPosition("vertical", 0)
  }
  var req = {
    "theme": theme.value,
    "platform": activePlatform.value,
    "catname": activeMenu.value,
    "catnameLike": menuLike.value,
    "baseType": activeType.value.value,
    "basePublisher": activePublisher.value.value,
    "baseYear": activeYear.value.value,
    "baseCountry": activeCountry.value.value,
    "baseTranslate": activeTranslate.value.value,
    "baseVersion": activeVersion.value.value,
    "baseProducer": activeProducer.value.value,
    "score": activeScore.value.value,
    "complete": activeComplete.value.value,
    "keyword": activeKeyword.value,
    "page": 0,
    "simpleModel": "simple",
    "letter": activeLetter.value == "ALL" ? "" : activeLetter.value,
  };
  var request = JSON.stringify(req);
  GetGameList(request).then((result: string) => {

    let resp = decodeApiData(result)
    items.value = resp.data;
    console.log(resp.data)
  })
  GetGameCount(request).then((result: string) => {
    let resp = decodeApiData(result)
    romCount.value = resp.data;
  })
});

//点击游戏
function clickGame(romId: number, index: number) {
  activeRom.value = romId
  activeFocus.value = [3, index]
}


function nextPage(index: number, done: any) {
  if (pageEnd) {
    done();
    return;
  }

  //启动时第一页不加载
  if (pageNum == 0) {
    pageNum++;
    done();
    return;
  }
  let req = {
    "theme": theme.value,
    "platform": activePlatform.value,
    "catname": activeMenu.value,
    "catnameLike": menuLike.value,
    "baseType": activeType.value.value,
    "basePublisher": activeProducer.value.value,
    "baseYear": activeYear.value.value,
    "baseCountry": activeCountry.value.value,
    "baseTranslate": activeTranslate.value.value,
    "baseVersion": activeVersion.value.value,
    "baseProducer": activePublisher.value.value,
    "score": activeScore.value.value,
    "complete": activeComplete.value.value,
    "keyword": activeKeyword.value,
    "page": pageNum,
    "simpleModel": "simple",
    "letter": activeLetter.value == "ALL" ? "" : activeLetter.value,
  };
  var request = JSON.stringify(req);
  GetGameList(request).then((result: string) => {
    let resp = decodeApiData(result)
    console.log(resp.data)
    if (resp.data.length > 0) {
      items.value = items.value.concat(resp.data)
      setTimeout(function () {
        pageNum++;
        done()
      }, 1000);
    } else {
      console.log("pageEnd")
      pageEnd = true;
      done();
    }

  })
  GetGameCount(request).then((result: string) => {
    let resp = decodeApiData(result)
    romCount.value = resp.data;
  })

}

//键盘方向键
export function romEventKeyboard(direction: string) {

  let focusType = 3 //默认rom列表
  let eleName = "rom-ele-"
  if (document.getElementById('platform-ele-0') != null) {
    focusType = 1
    eleName = "platform-ele-"
  } else if (document.getElementById('game-ele-0') != null) {
    focusType = 4
    eleName = "game-ele-"
  } else if (document.getElementById('sim-ele-0') != null) {
    focusType = 5
    eleName = "sim-ele-"
  }

  if (focusType == 1) {
    if (document.getElementById('menu-ele-0') != null) {
      if (direction == "ArrowUp") { //平台替换为菜单
        focusType = 2
        eleName = "menu-ele-"
      } else if (direction == "ArrowDown") { //平台替换为菜单
        focusType = 2
        eleName = "menu-ele-"
      }
    }
  }

  let act = activeFocus.value[1]
  switch (focusType) {
    case 1: //平台
      if (activeFocus.value[0] != focusType) {
        let platforms = getPlatformList()
        platforms.forEach((item: any, index: number) => {
          if (item.Id == activePlatform.value) {
            act = index
          }
        })
      }
      break;
    case 2: //菜单
      if (activeFocus.value[0] != focusType) {
        let menus = getMenuListByPlatform(activePlatform.value)
        menus.forEach((item: any, index: number) => {
          if (item.Path == activeMenu.value) {
            act = index
          }
        })
        if (act == -1) {
          act = 0
        }
      }
      break
    default:
      if (activeFocus.value[0] != focusType) {
        act = -1
      }
  }

  //焦点不在ROM里
  activeFocus.value = [focusType, act];

  let currFocus = activeFocus.value[1];
  let newIndex = 0;
  let blockSize = platformUi.value.BlockSize == 0 ? 1 : platformUi.value.BlockSize;

  switch (direction) {
    case "ArrowUp":
      if (currFocus == 0) {
        return
      }
      if (focusType == 3) {
        newIndex = currFocus - blockSize
      } else {
        newIndex = currFocus - 1
      }
      newIndex = newIndex < 0 ? 0 : newIndex;
      break;
    case "ArrowDown":
      if (focusType == 3) {
        newIndex = currFocus + blockSize
      } else {
        newIndex = currFocus + 1
      }
      break;
    case "ArrowLeft":
      if (currFocus == 0) {
        return
      }
      newIndex = currFocus - 1
      newIndex = newIndex < 0 ? 0 : newIndex;
      break;
    case "ArrowRight":
      newIndex = currFocus + 1
      break;
  }


  if ((focusType == 4 || focusType == 5) && direction == "ArrowDown" && newIndex == 0) {
    newIndex = 1
  }

  let ele = document.getElementById(eleName + newIndex);

  if (ele == null) {
    newIndex = 0
    ele = document.getElementById(eleName + newIndex);
  }

  if (ele == null) {
    return
  }

  activeFocus.value = [focusType, newIndex]

  if (focusType == 2 || focusType == 4 || focusType == 5) {
    ele?.focus();
  } else {
    ele?.focus()
    ele?.click()
  }
}


//手柄确定键
export function clickEventKeyboard() {

  let focusType = 3 //默认rom列表
  let eleName = "rom-ele-"
  if (document.getElementById('menu-ele-0') != null) {
    focusType = 2
    eleName = "menu-ele-"
  } else if (document.getElementById('game-ele-0') != null) {
    focusType = 4
    eleName = "game-ele-"
  } else if (document.getElementById('sim-ele-0') != null) {
    focusType = 5
    eleName = "sim-ele-"
  }

  console.log("click ele", eleName + activeFocus.value[1])

  //如果点击菜单，加长时间防止连击
  if (focusType == 2) {
    let t = getTimestamp()
    if (clickExpire > 0 && t - clickExpire < 2) {
      return
    }
    clickExpire = t
  }

  let ele = document.getElementById(eleName + activeFocus.value[1]);
  if (ele == null) {
    ele = document.getElementById(eleName + "0");
    if (ele == null) {
      return
    }
  }

  if (focusType == 3) {
    openRunGameDialogByKeyboard()
  } else {
    ele?.focus()
    ele?.click()
  }

}


//手柄取消键
export function cancelEventKeyboard() {

  let focusType = 3 //默认rom列表
  if (document.getElementById('menu-ele-0') != null) {
    focusType = 2
  } else if (document.getElementById('game-ele-0') != null) {
    focusType = 4
  } else if (document.getElementById('sim-ele-0') != null) {
    focusType = 5
  }

  if (focusType == 2) { //关闭平台
    closePlatformDialog()
  } else if (focusType == 4 || focusType == 5) {  //关闭子游戏
    closeContentDialog()
  }

}


//更换背景图
export function changeContentBackground(bg: string) {

  root.style.setProperty('--base-font-size', platformUi.value.BaseFontsize);

  if (isEmpty(bg)) {
    root.style.setProperty('--content-backgorund-image', 'url(' + platformUi.value.BackgroundImage + ')');
  } else {
    root.style.setProperty('--content-backgorund-image', 'url(' + bg + ')');
  }

  if (platformUi.value.BackgroundMask != "") {
    root.style.setProperty('--content-backgorund-mask', 'url(' + platformUi.value.BackgroundMask + ')');
  }else{
    root.style.setProperty('--content-backgorund-mask', '');
  }

  let r = platformUi.value.BackgroundRepeat;
  if (r == "" || r == "cover") {
    root.style.setProperty('--content-backgorund-repeat', "no-repeat");
    root.style.setProperty('--content-backgorund-size', "100% 100%");
  } else if (r == "repeat") {
    root.style.setProperty('--content-backgorund-repeat', r);
    root.style.setProperty('--content-backgorund-size', "auto auto");
  } else {
    root.style.setProperty('--content-backgorund-repeat', r);
    root.style.setProperty('--content-backgorund-size', "contain");
  }
}

export default {
  components: {},
  setup() {
    return {
      nextPage,
      clickGame,
      activeMenu,
      activePlatform,
      activeRom,
      romCount,
      config,
      lang,
      romState,
      activeKeyword,
      activeType,
      activeYear,
      activeProducer,
      activePublisher,
      activeCountry,
      activeTranslate,
      activeVersion,
      activeScore,
      activeComplete,
      activeLetter,
      scrollAreaRef,
      menuLike,
      platformUi,
      callbackOpts,
      items,
      rating,
      listColumns,
      letterList,
      initialPagination,
      qItems,
      isEmpty,
    }
  }
}

</script>

<style scoped>
@import "src/css/playnite/romlistBar.css";

</style>
