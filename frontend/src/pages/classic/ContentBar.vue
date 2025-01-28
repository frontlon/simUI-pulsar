<template>
  <!-- ROM列表 -->
  <div :class="'rom-wrapper fuzzy' + platformUi.BackgroundFuzzy">
    <div class="bg-mask"></div>
    <q-scroll-area ref="scrollAreaRef" class="rom-scroll" :visible="false">
      <q-infinite-scroll @load="nextPage" :offset="100" debounce="500" style="width:100%;max-width: 100%">

        <!-- 模块样式1 -->
        <div v-if="platformUi.RomListStyle == 1" class="rom-block">
          <div v-for="(item, index) in items" :style="'width:'+ 100 / platformUi.BlockSize + '%;'">
            <q-item clickable :class="'q-pa-sm q-ma-' + platformUi.BlockMargin"
                    :active-class="'active animate__animated animate__faster ' + platformUi.BlockClickAnimate"
                    :id="'ele'+ index" :active="activeRom === item.Id"
                    @contextmenu="showContextMenu($event, index,item.Id, item)"
                    @click="clickGame(item.Id,index)">
              <q-card flat style="width: 100%" :style="platformUi.BlockHideBackground ?'background:none' : ''"
                      @dblclick="dbClickRunGame(item.Id)">

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

        <!-- 模块样式2 -->
        <div v-else-if=" platformUi.RomListStyle== 2
                " class="rom-block">
          <div v-for="(item, index) in items" :style="'width:'+ 100 / platformUi.BlockSize + '%;'">
            <q-item clickable :class="'q-pa-sm q-ma-' + platformUi.BlockMargin"
                    :active-class="'active animate__animated animate__faster ' + platformUi.BlockClickAnimate"
                    :id="'ele'+ index" :active="activeRom === item.Id"
                    @contextmenu="showContextMenu($event, index,item.Id, item)"
                    @click="clickGame(item.Id,index)">
              <q-card flat style="width: 100%" @dblclick="dbClickRunGame(item.Id)"
                      :style="platformUi.BlockHideBackground ?'background:none' : ''">
                <q-img v-if="item.ThumbPic" :src="item.ThumbPic" loading="lazy" fit="contain"
                       :class="'img-direction-'+platformUi.BlockDirection">
                  <div class="absolute-bottom text-center" v-if="!platformUi.BlockHideTitle"
                       :style="{'padding':'10px 0','color':'var(--color-text)','background':platformUi.BlockHideBackground ?'none' : ''}">
                    {{ platformUi.NameType == 1 ? item.Name : item.RomName }}
                  </div>
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
              </q-card>
            </q-item>
          </div>
        </div>

        <!-- 模块样式3 -->
        <div v-else-if="platformUi.RomListStyle == 3" class="rom-block">
          <div v-for="(item, index) in items" :style="'width:'+ 100 / platformUi.BlockSize + '%;'">
            <q-item clickable :class="'q-pa-sm q-ma-' + platformUi.BlockMargin"
                    :active-class="'active animate__animated animate__faster ' + platformUi.BlockClickAnimate"
                    :id="'ele'+ index" :active="activeRom === item.Id"
                    @contextmenu="showContextMenu($event, index,item.Id, item)"
                    @click="clickGame(item.Id,index)">
              <q-card flat style="width: 100%" @dblclick="dbClickRunGame(item.Id)"
                      :style="platformUi.BlockHideBackground ?'background:none' : ''">
                <q-img v-if="item.ThumbPic" :src="item.ThumbPic" loading="lazy" fit="contain"
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
                <q-card-section>
                  <div class="row no-wrap items-center">
                    <div class="col" v-if="!platformUi.BlockHideTitle"
                         :style="{'background':platformUi.BlockHideBackground ?'none' : ''}">
                      {{ platformUi.NameType == 1 ? item.Name : item.RomName }}
                    </div>
                  </div>
                  <q-rating v-model="item.Score" :max="5" size="14px" color="orange" color-half="orange"
                            icon-half="star_half" readonly/>
                </q-card-section>
              </q-card>
            </q-item>
          </div>
        </div>

        <!-- 列表模式 -->
        <div v-else-if="platformUi.RomListStyle == 4" class="q-pa-md">

          <q-markup-table flat standout square dense class="rom-list">
            <thead>
            <tr>
              <th v-for="col in listColumns">{{ col.label }}</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="(item,index) in items" :id="'ele'+ index"
                :class="activeRom === item.Id ? 'active animate__animated animate__faster ' + platformUi.BlockClickAnimate : ''"
                @contextmenu="showContextMenu($event, index,item.Id, item)"
                @click="clickGame(item.Id,index)" @dblclick="dbClickRunGame(item.Id)">
              <td v-for="col in listColumns">{{ item[col.field] }}</td>
            </tr>
            </tbody>
          </q-markup-table>
        </div>

        <template v-slot:loading>
          <div class="row justify-center q-my-md">
            <q-spinner-dots color="primary" size="40px"/>
          </div>
        </template>

      </q-infinite-scroll>

    </q-scroll-area>
  </div>

  <q-list class="letter-list">
    <q-item clickable v-ripple dense v-for="letter in letterList" :key="letter" class="letter"
            active-class="letter-active" :active="activeLetter === letter" @click="changeLetter(letter)">
      <q-item-section>
        {{ letter }}
      </q-item-section>
    </q-item>
  </q-list>

  <!-- 悬浮按钮 -->
  <div class="fab">
    <q-fab direction="up" label-position="right" label-class="bg-grey-3 text-purple" color="orange"
           icon="keyboard_arrow_up">
      <q-fab-action label-class="bg-grey-3 text-grey-8" external-label label-position="left" color="primary"
                    icon="hide_source" :label="lang.hideGame" @click="getHideGame()"/>
      <q-fab-action label-class="bg-grey-3 text-grey-8" external-label label-position="left" color="positive"
                    icon="new_label" :label="lang.inputGame" @click="openInputDialog()"/>
      <!--      <q-fab-action label-class="bg-grey-3 text-grey-8" external-label label-position="left" color="positive"-->
      <!--                    icon="new_label" :label="lang.outputGame" @click="openOutputDialog()"/>-->
      <q-fab-action label-class="bg-grey-3 text-grey-8" external-label label-position="left" color="secondary"
                    icon="menu_open" :label="lang.menuManage" @click="openManageMenuDialog()"/>
    </q-fab>
  </div>

  <!-- fab功能 -->
  <fab-menu/>
  <fab-input-rom/>
  <fab-output-rom/>
  <!-- 右键菜单 -->
  <content-context/>
</template>
<script lang="ts">

import {ref, watch} from 'vue';
import ContentContext, {keybordOpenWin, showContextMenu} from 'pages/classic/context/Context.vue';
import FabMenu, {openManageMenuDialog} from 'pages/classic/modules/FabMenu.vue';
import FabInputRom, {openInputDialog} from 'pages/classic/modules/FabInputRom.vue';
import FabOutputRom, {openOutputDialog} from 'pages/classic/modules/FabOutputRom.vue';

import 'animate.css'
import {decodeApiData, deepClone, firstLetterToLower, isEmpty, loading, notify} from "components/utils";
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {CreateRomCache, GetGameCount, GetGameList, RunGame} from "app/wailsjs/go/controller/Controller";

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
let focusType = 3; //3rom列表

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

//监听rom列表刷新
watch([romState], async(newValue, oldValue) => {
  pageEnd = false;
  pageNum = 0;
  //滚动条回到最顶端
  scrollAreaRef.value.setScrollPosition("vertical", 0)

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
    "page": pageNum,
    "simpleModel": "simple",
    "letter": activeLetter.value == "ALL" ? "" : activeLetter.value,
  };

  var request = JSON.stringify(req);
  let result = await GetGameList(request)

  let resp = decodeApiData(result)
  items.value = resp.data;

  //列表模式
  listColumns.value = [{label: lang.value.alias, field: "Name"},
    {label: lang.value.romName, field: "RomName"}
  ]

  platformUi.value.RomListColumn.forEach((item: any, index: number) => {
    listColumns.value.push({
      label: lang.value[firstLetterToLower(item)],
      field: item,
    })
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

//右键菜单操作回调
watch(callbackOpts, (newValue: any, oldValue) => {

  // let romId = newValue.id;
  let opt = newValue.opt;
  let romIndex = newValue.index;
  let data = newValue.data;

  switch (opt) {
    case "renameRomLink": //模块重命名
    case "editBaseInfo": //编辑资料
      items.value[romIndex].Name = data
      break;
    case "bindSubGame":
      global.incRomState();
      break;
    case "romHide": //rom隐藏
    case "moveRomLink": //移动ROM
    case "deleteRomLink": //删除rom链接
    case "deleteRom": //删除rom
      items.value.splice(romIndex, 1)
      break;
    case "copyRomLink": //复制模块
      if (activeMenu.value == '' || activeMenu.value == data.menu) {
        let item = deepClone(items.value[romIndex])
        item.Id = data.id
        item.Name = data.name
        item.Menu = data.menu
        items.value.splice(romIndex + 1, 0, item);
      }
      break;
    case "inputRom": //导入ROM
      loading("show", lang.loading)
      CreateRomCache(activePlatform.value).then((result: string) => {
        global.incRomState();
        loading("hide")
      })
      break;
    case "changeThumb": //更新展示图
      items.value[romIndex].ThumbPic = data
      break;
    case "setFavorite": //设置喜爱
      if(activeMenu.value == "favorite" && data == 0){
        items.value.splice(romIndex, 1)
      }
  }
})

function nextPage(index: number, done: any) {
  if (pageEnd) {
    done();
    return;
  }

  console.log("pageNum", pageNum)

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

//运行游戏
export function dbClickRunGame(id: number) {
  console.log("dbclick", id)
  RunGame(id, 0).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
  })
}

//字母搜索
function changeLetter(letter: string) {
  activeLetter.value = letter;
  global.incRomState()
}

//键盘方向键
export function romEventKeyboard(direction: string) {

  let newIndex = 0;

  let blockSize = platformUi.value.BlockSize == 0 ? 1 : platformUi.value.BlockSize;
  //列表模式
  if (platformUi.value.RomListStyle == 4) {
    blockSize = 1
  }

  //焦点不在ROM里
  if (activeFocus.value[0] != focusType) {
    activeFocus.value = [focusType, -1];
  }

  let currFocus = activeFocus.value[1];

  switch (direction) {
    case "ArrowUp":
      if (currFocus == 0) {
        return
      }
      newIndex = currFocus - blockSize
      newIndex = newIndex < 0 ? 0 : newIndex;
      break;
    case "ArrowDown":
      newIndex = currFocus + blockSize
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
  console.log("-=-=-=-", currFocus)

  let ele = document.getElementById('ele' + newIndex);
  if (isEmpty(ele)) {
    return
  }

  activeFocus.value = [focusType, newIndex]

  if (platformUi.value.RomListStyle == 4) {
    //列表模式点击
    ele?.click()
  } else {
    ele?.focus()
  }

}

//事件点击游戏
export function clickGameEvent() {
  if (activeFocus.value[0] != 3) {
    return
  }
  let ele = document.getElementById('ele' + activeFocus.value[1]);
  ele?.click()
}

//keyboard event
export function manageEventKeyboard(opt) {
  if (activeFocus.value[0] != focusType) {
    return
  }
  let rom = items.value[activeFocus.value[1]]
  keybordOpenWin(opt, activeFocus.value[1], rom)
}


//运行游戏 - 当前焦点
export function runActiveGame() {
  if (activeFocus.value[0] != 3) {
    return
  }

  let rom = items.value[activeFocus.value[1]]
  if (isEmpty(rom)) {
    return
  }

  RunGame(rom.Id, 0).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
  })
}

//显示隐藏的游戏
export function getHideGame() {
  if (activeMenu.value == "hide") {
    return;
  }
  activeMenu.value = "hide";
  global.incRomState();
}

export default {
  components: {ContentContext, FabMenu, FabInputRom, FabOutputRom},
  setup() {
    return {
      showContextMenu,
      openManageMenuDialog,
      getHideGame,
      openInputDialog,
      openOutputDialog,
      nextPage,
      dbClickRunGame,
      changeLetter,
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
    }
  }
}

</script>

<style scoped>
@import "src/css/classic/common.css";
@import "src/css/classic/contentBar.css";
</style>
