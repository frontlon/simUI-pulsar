<template>
  <!-- ROM列表 -->
  <div class="platform-title">
    {{ platformName }}
  </div>

  <div class="rom-wrapper" v-if="items != null && items.length > 0">
    <div v-if="platformUi.HideCarousel == 0" ref="carouselRef" class="thumb-box animate__animated animate__faster">
      <q-carousel v-if="thumbListCarousel.length > 0" animated infinite :autoplay="2500" v-model="slide"
                  class="carousel">
        <q-carousel-slide v-for="(item,index) in thumbListCarousel" :name="index">
          <q-img class="carousel-img" :src="item.Path" fit="contain"/>
        </q-carousel-slide>
      </q-carousel>
      <div class="no-carousel" v-else>
        <q-icon size="lg" name="wallpaper"></q-icon> &nbsp;
        {{ lang.noThumb }}
      </div>
    </div>
    <div class="page-box">
      <div class="page-num top">{{ detailIndex }}</div>
      <div class="page-num bottom">{{ romCount }}</div>
    </div>
    <div class="rom-box" id="rom-box">
      <q-scroll-area ref="scrollAreaRef" class="rom-scroll">
        <q-infinite-scroll @load="nextPage" :offset="100" debounce="500" style="width:100%;max-width: 100%">

          <q-list class="rom-list">
            <q-item clickable dense v-for="(item, index) in items"
                    :active-class="'active animate__animated  animate__flipInX'"
                    :id="'rom-ele-'+ index" :active="activeRom === item.Id"
                    @click="clickGame(item.Id,index)">
              <q-item-section side v-if="item.ThumbPic">
                <q-img :src="item.ThumbPic" loading="lazy" fit="contain" class="rom-img"/>
              </q-item-section>
              <q-item-section>
                <div class="rom-title" @dblclick="openRunGameDialog()">
                  <div class="master-title">{{ platformUi.NameType == 1 ? item.Name : item.RomName }}</div>
                  <div class="sub-title">{{ item.subTitle }}</div>
                </div>
              </q-item-section>
            </q-item>
          </q-list>

        </q-infinite-scroll>

      </q-scroll-area>

    </div>


  </div>
  <div class="rom-wrapper empty-romlist" v-else-if="items != null">
    {{ lang.noGame }}
  </div>

  <!-- 选择游戏对话框 -->
  <q-dialog v-if="!isEmpty(detail)" v-model="showGameDialog" @hide="callbackCloseDialog" transition-show="scale"
            transition-hide="scale">
    <q-card style="min-width: 600px">
      <q-card-section>
        <q-list class="menu-list">
          <q-item clickable v-ripple id="game-ele-0" key="0"
                  @click.stop="openSimDialog(detail.Info.Id)">
            <q-item-section>{{ detail.Info.Name }}</q-item-section>
          </q-item>
          <q-item clickable v-ripple
                  v-for="(sg,sIndex) in detail.Sublist" :id="'game-ele-'+(sIndex+1)" :key="sIndex+1"
                  @click.stop="openSimDialog(sg.Id)">
            <q-item-section>{{ sg.Name }}</q-item-section>
          </q-item>

        </q-list>
      </q-card-section>
    </q-card>
  </q-dialog>

  <!-- 选择模拟器对话框 -->
  <q-dialog v-if="simulators" v-model="showSimDialog" @hide="callbackCloseDialog" transition-show="scale"
            transition-hide="scale">
    <q-card style="min-width: 600px">
      <q-card-section>
        <q-list class="sim-list" v-if="simulators">
          <q-item v-for="(item,index) in simulators" clickable v-ripple :id="'sim-ele-'+index" :key="index"
                  @click="runGame(item.Id)">
            <q-item-section>
              {{ item.Name }} - {{ index }}
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>
  </q-dialog>


</template>
<script lang="ts">
import 'animate.css'
import {ref, watch} from 'vue';
import {decodeApiData, getTimestamp, isEmpty} from "components/utils";
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {
  GetGameCount,
  GetGameDetail,
  GetGameList,
  GetGameThumbs,
  GetPlatformById,
  RunGame,
  SetRomSimId
} from "app/wailsjs/go/controller/Controller";
import {closePlatformDialog, getMenuListByPlatform, getPlatformList} from 'pages/tiny/Platform.vue'

const global = useGlobalStore();
const {
  activeMenu, activePlatform, activeRom, activeFocus, romCount, config, lang, romState, activeKeyword,
  activeType, activeYear, activeProducer, activePublisher, activeCountry, activeTranslate, simulatorMap,
  activeVersion, activeScore, activeComplete, activeLetter, menuLike, platformUi, callbackOpts, theme
} = storeToRefs(global);

const items: any = ref(null)
const detail: any = ref(null)
const detailIndex: any = ref(0)
const platformName: any = ref("")
const root = document.documentElement;
let pageNum = 0; //翻页数
let pageEnd = false; //是否翻到最后一页
const thumbListCarousel: any = ref([]);
const slide = ref(0)
const carouselRef = ref(null)
const scrollAreaRef = ref(null)

const showGameDialog: any = ref(false)
const showSimDialog: any = ref(false)
const selectGameId: any = ref(null)
const simulators: any = ref([])
let tmpRomFocus = null
let clickExpire = 0

//监听rom列表刷新
watch([romState], (newValue, oldValue) => {

  pageEnd = false;

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
    items.value = [];

    if (resp.data.length == 0) {
      return
    }

    resp.data.forEach((item, index) => {
      resp.data[index]['subTitle'] = createSubTitle(item)
    });
    items.value = resp.data;

    //默认选中第一个游戏
    if (resp.data.length > 0) {
      clickGame(resp.data[0].Id, 0)
    }

    //滚动条回到最顶端
    setTimeout(function () {
      if (scrollAreaRef.value != null) {
        scrollAreaRef.value.setScrollPosition("vertical", 0)
      }
    }, 50)

  })
  GetGameCount(request).then((result: string) => {
    let resp = decodeApiData(result)
    romCount.value = resp.data;
  })

  //读取平台信息
  GetPlatformById(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    platformName.value = resp.data.Name;
  })


});

function createSubTitle(item) {
  let subTitle = []
  if (item.Name != item.RomName) {
    subTitle.push(item.RomName)
  }
  if (item.BaseYear != "") {
    subTitle.push(item.BaseYear)
  }
  if (item.BaseType != "") {
    subTitle.push(item.BaseType)
  }
  return subTitle.join(" | ")
}

//点击游戏
function clickGame(romId: number, index: number) {
  activeRom.value = romId
  activeFocus.value = [3, index]

  //读取rom信息
  detail.value = items.value[index]
  detailIndex.value = index + 1

  //读取模拟器信息
  if (!isEmpty(simulatorMap.value[activePlatform.value])) {
    simulators.value = simulatorMap.value[activePlatform.value]
  }

  //读取相册
  GetGameThumbs(activeRom.value, "").then((result: string) => {
    let resp = decodeApiData(result)
    //整理图片数据
    handleCarouselAndAlumb(resp.data);

    //更新动画
    if (carouselRef.value != null) {
      if (carouselRef.value.classList.contains("animate__fadeInDown")) {
        carouselRef.value.classList.remove("animate__fadeInDown");
      }
      setTimeout(function () {
        carouselRef.value.classList.add("animate__fadeInDown");
      }, 50)
    }
  })

}

//翻页
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
    if (resp.data.length > 0) {

      resp.data.forEach((item, index) => {
        resp.data[index]['subTitle'] = createSubTitle(item)
      });

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

//图片数据整理
function handleCarouselAndAlumb(data) {
  thumbListCarousel.value = []
  let thumbMap = {}
  data.forEach((val: any) => {
    if (thumbMap[val.Type] == undefined) {
      thumbMap[val.Type] = [];
    }
    thumbMap[val.Type].push(val);
  })

  config.value.ThumbOrders.forEach((typ: string) => {

    if (thumbMap[typ] == undefined || typ == "video" || typ == "icon" || typ == "background") {
      return
    }

    thumbListCarousel.value = thumbListCarousel.value.concat(thumbMap[typ])
  })

  if (thumbMap["background"] != undefined && thumbMap["background"].length > 0) {
    changeContentBackground(thumbMap["background"][0].Path)
  } else {
    changeContentBackground("")
  }
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
  let romNum = 0
  let romEle = document.getElementById('rom-ele-' + currFocus);
  if (romEle == null) {
    romEle = document.getElementById('rom-ele-0');
  }

  if (romEle != null) {
    let boxEle = document.getElementById('rom-box');
    romNum = parseInt(boxEle.offsetHeight / romEle.offsetHeight);
  }

  switch (direction) {
    case "ArrowUp":
      if (currFocus == 0) {
        return
      }
      newIndex = currFocus - 1
      newIndex = newIndex < 0 ? 0 : newIndex;
      break;
    case "ArrowDown":
      newIndex = currFocus + 1
      break;
    case "ArrowLeft":
      if (focusType == 3) {
        newIndex = currFocus - romNum
      } else {
        newIndex = currFocus - 1
      }
      newIndex = newIndex < 0 ? 0 : newIndex;
      if (currFocus == 0) {
        return
      }
      break;
    case "ArrowRight":
      if (focusType == 3) {
        newIndex = currFocus + romNum
        newIndex = newIndex > items.value.length - 1 ? items.value.length - 1 : newIndex;
        if (currFocus == items.value.length - 1) {
          return
        }
      } else {
        newIndex = currFocus + 1
      }
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
  } else {
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

export function openRunGameDialogByKeyboard() {

  if (activeFocus.value[0] != 3) {
    return
  }

  openRunGameDialog()
}


//打开游戏选择对话框
function openRunGameDialog() {

  tmpRomFocus = activeFocus.value

  //读取rom信息
  GetGameDetail(activeRom.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      detail.value = resp.data

      //设置默认模拟器选项
      let simId = 0
      if (simulators.value && simulators.value.length > 0) {
        simulators.value.forEach((item: any, index: number) => {
          if (item.Id == detail.value.Info.SimId) {
            simId = item.Id
          }
        })
      }
      detail.value.Info.SimId = simId

      if (detail.value.Sublist.length == 0) {
        //一个游戏，直接打开模拟器对话框
        openSimDialog(detail.value.Info.Id)
      } else {
        showGameDialog.value = true

        //设置默认模拟器选项
        setTimeout(function () {
          //activeFocus.value = [4, 0]
          let ele = document.getElementById('game-ele-0');
          ele?.focus()
        }, 100);

      }

    }
  })

}


//打开模拟器选择对话框
function openSimDialog(romId) {
  selectGameId.value = romId
  simulators.value = []

  let platform = detail.value.Info.Platform
  if (!simulatorMap.value[platform] || simulatorMap.value[platform].length == 0) {
    //无模拟器运行
    runGame(0)
  } else if (simulatorMap.value[platform].length == 1) {
    //只有一个模拟器运行游戏
    let simId = simulatorMap.value[platform][0].Id
    runGame(simId)
  } else if (simulatorMap.value[platform].length > 1) {
    //选择模拟器
    simulators.value = simulatorMap.value[platform]
    showSimDialog.value = true

    //设置默认模拟器选项
    setTimeout(function () {
      simulators.value.forEach((item: any, index: number) => {
        if (item.Id == detail.value.Info.SimId) {
          let ele = document.getElementById('sim-ele-' + index);
          ele?.focus()
          //activeFocus.value = [5, index]
        }
      })
    }, 100);
  }

  showGameDialog.value = false
}

function runGame(simId) {
  SetRomSimId(detail.value.Info.Id, simId)
  RunGame(selectGameId.value, simId)
  showSimDialog.value = false
}

function callbackCloseDialog() {
  activeFocus.value = tmpRomFocus
}

function closeContentDialog() {
  showGameDialog.value = false
  showSimDialog.value = false
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
      detail,
      detailIndex,
      thumbListCarousel,
      slide,
      carouselRef,
      platformName,
      isEmpty,
      showGameDialog,
      callbackCloseDialog,
      simulators,
      showSimDialog,
      runGame,
      openSimDialog,
      openRunGameDialog,
    }
  }
}

</script>

<style scoped>
@import "src/css/tiny/romlistBar.css";

</style>
