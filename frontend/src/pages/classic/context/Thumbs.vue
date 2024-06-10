<template>
  <q-dialog v-model="showDialog" persistent transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 70%">
      <q-card-section class="row items-center q-pb-none q-mb-md">
        <div class="text-h6">{{ lang.editThumbs }}</div>
        <q-space/>
        <q-btn icon="close" flat round v-close-popup/>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <div class="row">
          <!-- 左窗口 -->
          <div class="col-3 left-bar">
            <q-scroll-area class="scroll-area">
              <!--图集-->

              <q-list v-if="thumbListAlumb">
                <q-expansion-item expand-icon-toggle expand-separator default-opened aria-disabled="true"
                                  v-for="(typ,index) in config.ThumbOrders">
                  <template v-slot:header>
                    <q-item-section>{{ albumTitle[typ] }}</q-item-section>
                    <q-item-section side>
                      <div class="text-grey-8 q-gutter-xs">
                        <q-btn size="sm" flat dense round icon="add" @click="addLocationImage(typ)"/>
                        <q-btn size="sm" flat dense round icon="keyboard_arrow_up" :disable="index == 0"
                               @click="cateSort('up',index)"/>
                        <q-btn size="sm" flat dense round icon="keyboard_arrow_down"
                               :disable="index == config.ThumbOrders.length-1"
                               @click="cateSort('down',index)"/>
                      </div>
                    </q-item-section>
                  </template>
                  <q-card v-if="thumbListAlumb && thumbListAlumb[typ].length > 0"
                          :class="{ 'drop-highlight': isDragOver[typ] }"
                          @dragover.prevent="handleDragOver($event,'over',typ)"
                          @dragleave="handleDragOver($event,'leave',typ)"
                          @drop.prevent="handleDrop($event,typ)">
                    <q-card-section>
                      <q-list>
                        <q-item v-for="(sub,sIndex) in thumbListAlumb[typ]" style="padding:10px" clickable
                                @contextmenu="showAlbumContextMenu($event, typ,sIndex,sub.Path)"
                                @click="showViewer(sub.Path)">
                          <q-img :src="sub.Path"/>
                        </q-item>
                      </q-list>
                    </q-card-section>
                  </q-card>
                  <q-card v-else class="empty"
                          :class="{ 'drop-highlight': isDragOver[typ] }"
                          @dragover.prevent="handleDragOver($event,'over',typ)"
                          @dragleave="handleDragOver($event,'leave',typ)"
                          @drop.prevent="handleDrop($event,typ)">
                    <q-card-section>
                      {{ lang.noThumb }}
                    </q-card-section>
                  </q-card>
                </q-expansion-item>
              </q-list>
            </q-scroll-area>
          </div>
          <!-- 右窗口 -->
          <div class="col q-ml-lg">
            <div class="row">
              <div class="col-3">
                <q-select filled square dense map-options :label="lang.selectDataResource" v-model="engine"
                          :options="dbOptions"/>
              </div>
              <div class="col q-pl-xs q-pr-xs">
                <q-input filled square dense :label="lang.gameKeyword" v-model="searchKey"/>
              </div>
              <div class="col-1">
                <q-btn class="search-btn" size="md" square color="primary" icon="search"
                       @click="searchDownloadThumb()"/>
              </div>
            </div>
            <!--图片搜索列表-->
            <div v-if="!isEmpty(downloadMap)" class="col">
              <q-scroll-area class="scroll-area">
                <q-infinite-scroll @load="nextPage" :offset="100" debounce="500">
                  <div v-for="(items, type) in downloadMap">
                    <h6>{{ type == "" ? "default" : type }}</h6>
                    <div class="thumb-list">
                      <div v-for="(item, index) in items" style="width:20%;height: auto">
                        <q-item clickable class="q-pa-md"
                                @contextmenu="showWebContextMenu($event, item.Ext, item.ImgUrl)"
                                @click="showViewer(item.ImgUrl)">
                          <q-card flat style="width:100%;" class="thumb-list-item">
                            <q-img :src="item.ImgUrl" loading="lazy" class="img-direction-1">
                              <template v-slot:error>
                                <div class="absolute-full flex flex-center">{{ lang.loadErr }}</div>
                              </template>
                              <template v-slot:loading>
                                <div class="absolute-full flex flex-center">{{ lang.loading }}</div>
                              </template>
                            </q-img>
                            <q-card-section>
                              <div class="module-title" :style="'font-size:'+ platformUi.BaseFontsize">
                                {{ item.Width }} x {{ item.Height }}<br>.{{ item.Ext }}
                              </div>
                            </q-card-section>
                          </q-card>
                        </q-item>
                      </div>
                    </div>
                  </div>
                  <template v-slot:loading>
                    <div class="row justify-center q-my-md">
                      <q-spinner-dots color="primary" size="40px"/>
                    </div>
                  </template>
                </q-infinite-scroll>
              </q-scroll-area>
            </div>
            <div v-else-if="listEmpty" class="empty">{{ lang.noSearchResult }}</div>
            <div v-else class="empty" style="padding-top: 100px">{{ lang.tipSelectDataResource }}</div>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script lang="ts">

import {h, ref, watch} from 'vue'
import {
  AddGameThumb,
  DelGameThumb,
  DownloadThumb,
  GetGameThumbs,
  LoadWebThumbs,
  OpenMultiFileDialog,
  OutputOneImage,
  SaveFileDialog,
  SortGameThumb,
  UpdateThumbsOrders
} from "app/wailsjs/go/controller/Controller";
import {
  callSrv,
  decodeApiData,
  getFileExt,
  getFileName,
  isEmpty,
  loading,
  notify,
  wailsPathDecode
} from "components/utils";
import {storeToRefs} from "pinia";
import {useGlobalStore} from "stores/globalData";
import ContextMenu from "@imengyu/vue3-context-menu";
import {api} from "v-viewer";
import 'viewerjs/dist/viewer.css'
import {CONTEXT_ICON_SIZE} from "boot/constant";
import {Dialog} from 'quasar';
import {getPromptOpts} from 'components/dialog'
import axios from "axios";

const global = useGlobalStore();
const thumbListAlumb: any = ref(null);
const {activePlatform, activeRom, platformUi, config, lang, callbackOpts, uploadServer} = storeToRefs(global);
const showDialog: any = ref(false)
const engine: any = ref(null)
const searchKey: any = ref("")
const romId: any = ref(null)
const romIndex: any = ref(null)
const romInfo: any = ref(null)
const downloadMap: any = ref({})
let pageNum = 0; //翻页数
let pageEnd = false; //是否翻到最后一页
const listEmpty = ref(false); //图片搜索列表是否为空
const albumTitle = ref({})
const dbOptions = ref([])
const isDragOver: any = ref({})

watch([romId], (newValue, oldValue) => {
  downloadMap.value = {}
  pageNum = 0;
  pageEnd = false;
  listEmpty.value = false
  engine.value = ""
})

export function openThumbsDialog(index: number, detail: any) {
  console.log("openThumbsDialog", detail)
  romInfo.value = detail
  romId.value = detail.Id
  romIndex.value = index
  searchKey.value = detail.Name
  showDialog.value = true
  albumTitle.value = {
    "thumb": lang.value.thumb,
    "snap": lang.value.snap,
    "poster": lang.value.poster,
    "packing": lang.value.packing,
    "title": lang.value.title,
    "cassette": lang.value.cassette,
    "icon": lang.value.icon,
    "gif": lang.value.gif,
    "background": lang.value.background,
    "video": lang.value.video,
  }
  dbOptions.value = [
    {label: lang.value.baidu, value: "baidu"},
    {label: 'HfsDB', value: "hfsdb"},
  ]

  //读取相册
  GetGameThumbs(romId.value, "").then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      //整理图片数据
      handleCarouselAndAlumb(resp.data);
    }
  })
}

//图片数据整理
function handleCarouselAndAlumb(lists: []) {
  thumbListAlumb.value = {}
  lists.forEach((val: any) => {
    if (thumbListAlumb.value[val.Type] == undefined) {
      thumbListAlumb.value[val.Type] = [];
    }
    thumbListAlumb.value[val.Type].push(val);
  })

  //追加上为空的类型
  Object.entries(albumTitle.value).forEach(function ([type, desc]) {
    if (thumbListAlumb.value[type] == undefined) {
      thumbListAlumb.value[type] = []
    }
  })

  console.log("thumbListAlumb", thumbListAlumb.value)
}

//点击搜索
function searchDownloadThumb() {

  if (isEmpty(engine.value.value)) {
    notify("err", lang.value.selectDataResource)
    return
  }

  pageEnd = false
  pageNum = 0
  downloadMap.value = {}
  listEmpty.value = false;
  loading("show", lang.value.loading + "...")
  LoadWebThumbs(engine.value.value, searchKey.value, pageNum).then((result: string) => {
    loading("hide")
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    if (!isEmpty(resp.data)) {
      downloadMap.value = resp.data
      setTimeout(function () {
        pageNum++;
      }, 1000);
    } else {
      console.log("pageEnd")
      pageEnd = true;
      listEmpty.value = true;
    }
  })
}

function nextPage(index: number, done: any) {

  if (pageEnd) {
    done();
    return;
  }

  //hfsdb只加载一页
  if (pageNum > 0 && engine.value.value != "baidu") {
    done();
    return;
  }

  //启动时第一页不加载
  if (pageNum == 0) {
    setTimeout(function () {
      done();
    }, 1000);
    return;
  }

  if (isEmpty(engine.value.value)) {
    notify("err", lang.value.selectDataResource)
    return
  }

  loading("show", lang.value.loading + "...")
  LoadWebThumbs(engine.value.value, searchKey.value, pageNum).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }

    if (!isEmpty(resp.data)) {
      Object.entries(resp.data).forEach(function ([type, items]) {
        if (downloadMap.value.hasOwnProperty(type)) {
          downloadMap.value[type] = downloadMap.value[type].concat(items)
        } else {
          downloadMap.value[type] = items
        }
      });
      pageNum++;
    } else {
      console.log("pageEnd")
      pageEnd = true;
    }
    loading("hide")
    setTimeout(function () {
      done()
    }, 1000);
  })
}

//网络图片右键菜单
function showWebContextMenu(event: any, ext: string, path: string) {

  event.preventDefault();
  let items: any = [];

  Object.entries(albumTitle.value).forEach(function ([type, desc]) {
    let item = {
      label: lang.value.setting + " " + desc,
      onClick: () => setThumbImage(type, ext, path)

    }
    items.push(item)
  })

  let contextData = {
    x: event.x,
    y: event.y,
    customClass: "default dark",
    items: items,
    zIndex: 9999,
  }

  ContextMenu.showContextMenu(contextData)
}

//图集图片右键菜单
function showAlbumContextMenu(event: any, type: string, index: number, path: string) {

  event.preventDefault();
  let items: any = [
    {
      label: lang.value.moveUp,
      icon: h('img', {src: '/images/context/up.png', style: CONTEXT_ICON_SIZE}),
      onClick: () => AlumbSort("up", type, index)
    },
    {
      label: lang.value.moveDown,
      divided: true,
      icon: h('img', {src: '/images/context/down.png', style: CONTEXT_ICON_SIZE}),
      onClick: () => AlumbSort("down", type, index)
    },
    {
      label: lang.value.output,
      divided: true,
      icon: h('img', {src: '/images/context/output.png', style: CONTEXT_ICON_SIZE}),
      onClick: () => outputImage(type, index, path)
    },

    {
      label: lang.value.delete,
      icon: h('img', {src: '/images/context/delete.png', style: CONTEXT_ICON_SIZE}),
      onClick: () => delAlbumImage(type, index, path)
    }
  ];

  let contextData = {
    x: event.x,
    y: event.y,
    customClass: "default dark",
    items: items,
    zIndex: 9999,
  }

  ContextMenu.showContextMenu(contextData)
}

//显示相册
function showViewer(path: string) {
  api({
    images: [path],
    options: {
      initialViewIndex: 0,
      zIndex: 10000,
    }
  })
}

//设置展示图
function setThumbImage(type: string, ext: string, path: string) {

  //下载图片
  let master = isEmpty(thumbListAlumb.value[type]) ? 1 : 0
  DownloadThumb(romId.value, type, master, path, ext).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }

    //图片插入到左侧
    if (thumbListAlumb.value[type] == undefined) {
      thumbListAlumb.value[type] = [];
    }
    let item = {
      Type: type,
      Path: resp.data,
      Mster: thumbListAlumb.value[type].length == 0,
      index: thumbListAlumb.value[type].length
    }
    thumbListAlumb.value[type].push(item);
    picCallback(type, thumbListAlumb.value[type][0].Path)
  })
}

//删除图集中的图片
function delAlbumImage(type: string, index: number, path: string) {
  let opt = getPromptOpts(lang.value.delPic, lang.value.tipDelPic, lang.value.ok, true)
  Dialog.create(opt).onOk(() => {
    let master = index == 0 ? 1 : 0
    DelGameThumb(romId.value, type, master, path).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err != "") {
        notify("err", resp.err)
        return
      }
      thumbListAlumb.value[type].splice(index, 1);

      if (thumbListAlumb.value[type].length > 0) {
        picCallback(type, thumbListAlumb.value[type][0].Path)
      } else {
        picCallback(type, "")
      }
      console.log(resp.data)
    })
  })
}

//图片排序
function AlumbSort(opt: string, type: string, currIndex: number) {

  if (opt == "up" && currIndex == 0) {
    return;
  }

  if (opt == "down" && currIndex == thumbListAlumb.value[type].length - 1) {
    return;
  }

  let req: any = [];
  let oldIndex = opt == "up" ? currIndex - 1 : currIndex + 1
  let currentRow = thumbListAlumb.value[type].splice(oldIndex, 1)[0];
  thumbListAlumb.value[type].splice(currIndex, 0, currentRow);
  thumbListAlumb.value[type].forEach((item: any, index: number) => {
    req.push(item.Path);
  });

  SortGameThumb(romId.value, type, req).then((result: string) => {
    let resp = decodeApiData(result)
    console.log(resp.data)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    thumbListAlumb.value[type] = resp.data

    picCallback(type, resp.data[0].Path)
  })
}

//导出一张图片
function outputImage(type: string, currIndex: number, path: string) {

  let de = wailsPathDecode(path)
  let name = getFileName(de)
  let ext = getFileExt(de)
  let newName = name + "_" + type + "_" + currIndex + ext

  SaveFileDialog(newName).then((result: string) => {
    let resp = decodeApiData(result)
    console.log(resp.data)
    if (resp.err != "" || resp.data == "") {
      return
    }

    //导出图片
    OutputOneImage(de, resp.data).then((result: string) => {
      let resp = decodeApiData(result)
      console.log(resp.data)
      if (resp.err != "") {
        notify("err", resp.err)
        return
      }
      notify("suc", lang.value.operSuc)
    })
  })

}

//分类排序
function cateSort(opt: string, currIndex: number) {

  if (opt == "up" && currIndex == 0) {
    return;
  }
  if (opt == "down" && currIndex == config.value.ThumbOrders.length - 1) {
    return;
  }
  let oldIndex = opt == "up" ? currIndex - 1 : currIndex + 1
  let currentRow = config.value.ThumbOrders.splice(oldIndex, 1)[0];
  config.value.ThumbOrders.splice(currIndex, 0, currentRow);
  UpdateThumbsOrders(config.value.ThumbOrders).then((result: string) => {
    let resp = decodeApiData(result)
    console.log(resp.data)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
  })
}

//添加本地图片
function addLocationImage(type: string) {

  OpenMultiFileDialog("media").then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    if (isEmpty(resp.data)) {
      return
    }

    let images = resp.data

    AddGameThumb(romId.value, type, images).then((result: string) => {
      let resp = decodeApiData(result)
      console.log(resp.data)
      if (resp.err != "") {
        notify("err", resp.err)
        return
      }

      let encodeImages = resp.data

      //追加数据
      encodeImages.forEach((img: any, index: number) => {
        let master = thumbListAlumb.value[type].length == 0 ? 1 : 0
        let s = {
          Type: type,
          Path: img,
          Master: master,
          Index: thumbListAlumb.value[type].length + 1
        }
        thumbListAlumb.value[type].push(s)

        picCallback(type, thumbListAlumb.value[type][0].Path)

      })

    })

  })


}


//图集拖拽上传
function handleDragOver(event: any, opt: string, type: string) {
  isDragOver.value[type] = opt == "over";
}

async function handleDrop(event: any, type: string) {
  isDragOver.value[type] = false;
  const files = event.dataTransfer.files;
  if (files.length > 0) {
    const file = files[0];
    // 检查文件类型
    if (!file.type.startsWith('image/') && !file.type.startsWith('video/')) {
      notify("err", lang.uploadTypeError)
      return;
    }

    // 读取文件
    let timeout = 0
    if (uploadServer.value == "") {
      //如果没有
      timeout = 1000
      uploadServer.value = await callSrv("StartUploadServer")
    }

    setTimeout(function () {
      uploadFile(files, type)
    }, timeout);

  }
}

//上传图片
function uploadFile(files, type) {
  const file = files[0];
  const formData = new FormData();
  formData.append('file', file);
  formData.append('id', romId.value);
  formData.append('type', type);

  axios.post(uploadServer.value + '/uploadFile', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  }).then(response => {
    console.log('上传成功', response.data.data);

    let item = {
      Index: thumbListAlumb.value[type].length,
      Master: thumbListAlumb.value[type].length <= 1 ? 1 : 0,
      Path: response.data.data,
      Type: type,
    }

    thumbListAlumb.value[type].push(item)

    picCallback(type, thumbListAlumb.value[type][0].Path)


  }).catch(error => {
    notify("err", error.toString())
  });
}

//更改图片回调首页列表更新
function picCallback(type: string, path: string) {
  if (type == platformUi.value.BlockThumbType) {
    //更新首页图片
    callbackOpts.value = {
      "index": romIndex.value,
      "opt": "changeThumb",
      "data": path,
    }
  }
}


export default {
  setup() {
    return {
      showDialog,
      albumTitle,
      thumbListAlumb,
      activeRom,
      platformUi,
      lang,
      config,
      dbOptions,
      engine,
      searchKey,
      downloadMap,
      listEmpty,
      isDragOver,
      searchDownloadThumb,
      nextPage,
      isEmpty,
      showWebContextMenu,
      showAlbumContextMenu,
      showViewer,
      setThumbImage,
      cateSort,
      addLocationImage,
      handleDragOver,
      handleDrop,
    }
  }
}

</script>

<style scoped>
.left-title {
  margin: 10px 0;
}

.left-bar {
  background: var(--color-2);
  padding: 5px;
}

h6 {
  width: 100%;
  margin: 10px 0;
}

.search-btn {
  height: 39px
}

.empty {
  text-align: center;
  color: var(--color-6);
}

.thumb-list {
  display: flex;
  flex-wrap: wrap;
}

.thumb-list-item {
  background: var(--color-1)
}

.module-title {
  text-align: center;
  word-break: break-all;
  white-space: pre-wrap;
}

.scroll-area {
  height: calc(100vh - 200px) !important;
}

.drop-highlight {
  background: yellow;
}
</style>
