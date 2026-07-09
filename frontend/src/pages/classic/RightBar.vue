<template>
  <div class="rightbar-wrapper scrollable">
    <div class="q-gutter-y-md">
      <div v-if="activeRom && detail">
        <q-card flat class="tab-card">
          <q-tabs v-model="panel" dense class="tabs" active-color="primary" indicator-color="primary" align="justify">
            <q-tab name="detail" :label="lang.info"/>
            <q-tab name="thumbs" :label="lang.thumbs"/>
            <q-tab name="strategy" :label="lang.strategy"/>
          </q-tabs>
          <q-separator/>

          <q-tab-panels v-model="panel" keep-alive animated class="panels">
            <q-tab-panel name="detail" class="detail-panel">
              <q-carousel v-if="thumbListCarousel.length > 0"
                          animated infinite arrows :autoplay="issetVideo ? 0 : 2000" v-model="slide"
                          transition-prev="jump-left" transition-next="jump-right"
                          @mouseenter="autoplay = false" @mouseleave="autoplay = true"
                          class="right-carousel detail-carousel">
                <q-carousel-slide v-for="(item,index) in thumbListCarousel" :name="index" style="padding: 0"
                                  :img-src="item.Type == 'video' ? '' : item.Path"
                                  @click="item.Type != 'video' ?showViewer(item.Index) : null">
                  <video v-if="item.Type == 'video'" ref="carouselVideoRef" @mouseleave="changeVideoVolume" loop
                         controls autoplay
                         :volume="config.VideoVolume" :key="item.Index" class="carousel-video" :src="item.Path">
                  </video>
                </q-carousel-slide>
              </q-carousel>
              <div class="detail-no-carousel" v-else>
                <q-icon size="lg" name="wallpaper"/>
                <div>{{ lang.noThumb }}</div>
              </div>

              <div class="detail-body">
                <div class="detail-section detail-section--overview">
                  <div class="rating">
                    <q-rating v-model="detail.Info.Score" :max="5" size="22px" color="orange" color-half="orange"
                              icon-half="star_half" @update:model-value="updateRating()"/>
                  </div>
                  <div class="state-row">
                      <div
                        class="state-stat state-stat--clickable"
                        :class="{
                          'state-stat--passed': detail.Info.Complete == 1,
                          'state-stat--perfect': detail.Info.Complete == 2,
                        }"
                        @click="updateComplete()"
                      >
                        <q-icon
                          :name="detail.Info.Complete == 1 ? 'img:/images/svg/complete1.svg'
                            : detail.Info.Complete == 2 ? 'img:/images/svg/complete2.svg'
                            : 'img:/images/svg/complete0.svg'"
                          size="22px"
                        />
                        <div
                          class="state-val"
                          :class="{
                            'text-positive': detail.Info.Complete == 1,
                            'state-perfect-text': detail.Info.Complete == 2,
                          }"
                        >
                          {{
                            detail.Info.Complete == 1 ? lang.passed
                              : detail.Info.Complete == 2 ? lang.perfectClear
                                : lang.notCleared
                          }}
                        </div>
                        <div class="state-label">{{ lang.complete }}</div>
                      </div>
                      <div class="state-stat state-stat--middle">
                        <q-icon name="img:/images/svg/number.svg" size="22px"/>
                        <div class="state-val">{{ detail.Info.RunNum }}</div>
                        <div class="state-label">{{ lang.runNum }}</div>
                      </div>
                      <div class="state-stat">
                        <q-icon name="img:/images/svg/time.svg" size="22px"/>
                        <div class="state-val state-time-text">{{
                            !isEmpty(detail.Info.RunLasttime) ? getDate(detail.Info.RunLasttime) : lang.notRun
                          }}
                        </div>
                        <div class="state-label">{{ lang.lastRunTime }}</div>
                      </div>
                    </div>
                  </div>

                  <div class="detail-divider"/>

                  <div class="detail-section">
                    <div class="strategy-section-title">
                      <q-icon name="play_circle" size="18px"/>
                      <span>{{ lang.runGame }}</span>
                    </div>

                    <q-list v-if="simulators.length > 0" class="sim-list" dense>
                      <q-item
                        v-for="item in simulators"
                        :key="item.Id"
                        class="sim-item"
                        :class="{ 'sim-item--active': detail.Info.SimId === item.Id }"
                      >
                        <q-item-section>
                          <q-item-label>
                            <q-radio size="xs" v-model="detail.Info.SimId" dense checked-icon="task_alt"
                                     unchecked-icon="panorama_fish_eye" color="primary" keep-color
                                     :val="item.Id" :label="item.Name" @click="UpdateRomSimId()"/>
                          </q-item-label>
                        </q-item-section>
                        <q-item-section side>
                          <q-btn size="xs" flat round icon="edit" class="sim-edit-btn"
                                 @click="openEditRomSimulatorDialog(detail.Info.Id,item.Id)"/>
                        </q-item-section>
                      </q-item>
                    </q-list>

                    <div class="run-btn-list">
                      <q-btn color="primary" no-caps class="run-btn-main"
                             :label="detail.Info.Name" @click="runGame(detail.Info.Id)"/>
                      <div v-if="detail.Sublist?.length" class="run-btn-sublist">
                        <q-btn flat no-caps class="run-btn-sub" v-for="(sg) in detail.Sublist"
                               :key="sg.Id" :label="sg.Name" @click="runGame(sg.Id)"/>
                      </div>
                    </div>
                  </div>

                  <div class="detail-divider"/>

                  <div class="detail-section">
                    <div class="strategy-section-title">
                      <q-icon name="info" size="18px"/>
                      <span>{{ lang.baseInfo }}</span>
                    </div>
                    <div class="rombase">
                      <p class="bg-indigo">{{ detail.Extra.PlatformName }}
                        <q-tooltip>{{ lang.belongPlatform }}</q-tooltip>
                      </p>
                      <p class="bg-indigo" v-if="detail.Info.Menu!= '' && detail.Info.Menu != '/'">{{ detail.Info.Menu }}
                        <q-tooltip>{{ lang.belongMenu }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseNameEn)" class="bg-orange">{{ detail.Info.BaseNameEn }}
                        <q-tooltip>{{ lang.baseNameEn }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseNameJp)" class="bg-red">{{ detail.Info.BaseNameJp }}
                        <q-tooltip>{{ lang.baseNameJp }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseType)" class="bg-pink">{{ detail.Info.BaseType }}
                        <q-tooltip>{{ lang.baseType }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseYear)" class="bg-purple">{{ detail.Info.BaseYear }}
                        <q-tooltip>{{ lang.baseYear }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseProducer)" class="bg-deep-purple">{{ detail.Info.BaseProducer }}
                        <q-tooltip>{{ lang.baseProducer }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BasePublisher)" class="bg-indigo">{{ detail.Info.BasePublisher }}
                        <q-tooltip>{{ lang.basePublisher }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseCountry)" class="bg-blue">{{ detail.Info.BaseCountry }}
                        <q-tooltip>{{ lang.baseCountry }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseTranslate)" class="bg-green">{{ detail.Info.BaseTranslate }}
                        <q-tooltip>{{ lang.baseTranslate }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseVersion)" class="bg-cyan-9">{{ detail.Info.BaseVersion }}
                        <q-tooltip>{{ lang.baseVersion }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseOtherA)" class="bg-blue-10">{{ detail.Info.BaseOtherA }}
                        <q-tooltip>{{ isEmpty(rombaseAlias.OtherA) ? lang.baseOtherA : rombaseAlias.OtherA }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseOtherB)" class="bg-teal">{{ detail.Info.BaseOtherB }}
                        <q-tooltip>{{ isEmpty(rombaseAlias.OtherB) ? lang.baseOtherB : rombaseAlias.OtherB }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseOtherC)" class="bg-yellow-10">{{ detail.Info.BaseOtherC }}
                        <q-tooltip>{{ isEmpty(rombaseAlias.OtherC) ? lang.baseOtherC : rombaseAlias.OtherC }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.BaseOtherD)" class="bg-brown-6">{{ detail.Info.BaseOtherD }}
                        <q-tooltip>{{ isEmpty(rombaseAlias.OtherD) ? lang.baseOtherD : rombaseAlias.OtherD }}</q-tooltip>
                      </p>
                      <p v-if="!isEmpty(detail.Info.Size)" class="bg-red-13">{{ detail.Info.Size }}
                        <q-tooltip>{{ lang.fileSize }}</q-tooltip>
                      </p>
                      <p v-if="detail.Extra.RomPath" class="rombase-path">{{ detail.Extra.RomPath }}
                        <q-tooltip>{{ lang.romPath }}</q-tooltip>
                      </p>
                    </div>
                  </div>

                  <template v-if="detail.DocContent">
                    <div class="detail-divider"/>
                    <div class="detail-section detail-section--doc">
                      <div class="strategy-section-title">
                        <q-icon name="article" size="18px"/>
                        <span>{{ lang.desc }}</span>
                      </div>
                      <q-scroll-area class="doc" ref="scrollAreaDocRef">
                        <div class="doc-content selectable" v-html="detail.DocContent"/>
                      </q-scroll-area>
                    </div>
                  </template>
              </div>
            </q-tab-panel>

            <q-tab-panel name="thumbs" class="thumbs-panel">
              <template v-if="thumbListCarousel.length > 0">
                <template v-for="typ in config.ThumbOrders" :key="typ">
                  <div
                    v-if="thumbListAlumb[typ]"
                    class="thumbs-album-box"
                  >
                    <div class="strategy-section-title">
                      <q-icon :name="albumIcon[typ] || 'image'" size="18px"/>
                      <span>{{ albumTitle[typ] }}</span>
                      <span class="thumbs-album-count">{{ thumbListAlumb[typ].length }}</span>
                    </div>

                    <div
                      v-if="typ != 'video'"
                      class="thumbs-grid"
                    >
                      <div
                        v-for="item in thumbListAlumb[typ]"
                        :key="item.Index"
                        class="thumbs-grid-item"
                        @click="showViewer(item.Index)"
                      >
                        <q-img
                          :src="item.Path"
                          class="thumbs-img"
                          fit="contain"
                        />
                        <div class="thumbs-img-overlay">
                          <q-icon name="zoom_in" size="20px"/>
                        </div>
                      </div>
                    </div>

                    <div v-else class="thumbs-video-list">
                      <video
                        v-for="item in thumbListAlumb[typ]"
                        :key="item.Index"
                        ref="albumVideoRef"
                        class="thumbs-video"
                        @mouseleave="changeVideoVolume"
                        loop
                        controls
                        autoplay
                        :volume="config.VideoVolume"
                        :src="item.Path"
                      />
                    </div>
                  </div>
                </template>
              </template>

              <div v-else class="thumbs-empty">
                <q-icon name="photo_library" size="48px"/>
                <div>{{ lang.noThumb }}</div>
              </div>
            </q-tab-panel>

            <q-tab-panel name="strategy" class="strategy-panel">
              <div v-if="strategyFiles.length > 0" class="strategy-files-box">
                <div class="strategy-section-title">
                  <q-icon name="folder_open" size="18px"/>
                  <span>{{ lang.strategyFiles }}</span>
                </div>
                <q-list class="strategy-file-list" dense>
                  <q-item
                    v-for="item in strategyFiles"
                    :key="item.Path"
                    v-ripple
                    clickable
                    class="strategy-file-item"
                    @click="runProgram(item.Path)"
                  >
                    <q-item-section avatar>
                      <q-avatar size="28px" class="strategy-file-icon">
                        <q-icon name="description" size="16px"/>
                      </q-avatar>
                    </q-item-section>
                    <q-item-section>
                      <q-item-label class="strategy-file-name">{{ item.Name }}</q-item-label>
                    </q-item-section>
                    <q-item-section side>
                      <q-icon name="launch" size="14px" class="strategy-file-action"/>
                    </q-item-section>
                  </q-item>
                </q-list>
              </div>

              <div v-if="strategy != ''" class="strategy-content-box">
                <div class="strategy-content selectable" v-html="formattedStrategy"></div>
              </div>
              <div v-else-if="strategy == '' && strategyFiles.length == 0" class="strategy-empty">
                <q-icon name="menu_book" size="48px"/>
                <div>{{ lang.noStrategy }}</div>
              </div>
            </q-tab-panel>
          </q-tab-panels>
        </q-card>
      </div>

      <div class="panel-desc" v-else-if="activePlatform > 0 && PlatformDesc" v-html="PlatformDesc"
           ref="platformDescRef">
      </div>
      <div v-else class="panel-default" :style="{'background-image':'url(\'' + logo+'\')'}">
      </div>

    </div>
  </div>
  <rom-sim/>
</template>

<script setup lang="ts">
import {computed, onMounted, ref, watch} from 'vue'
import 'viewerjs/dist/viewer.css'
import {api} from 'v-viewer'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {decodeApiData, getDate, isEmpty, notify} from "components/utils";
import {
  GetGameDetail,
  GetGameStrategy,
  GetGameStrategyFiles,
  GetGameThumbs,
  GetPlatformById,
  RunGame,
  RunProgram,
  SetRomSimId,
  UpdateGameComplete,
  UpdateGameScore,
  UpdateOneConfig,
} from 'app/wailsjs/go/controller/Controller'
import {changeContentBackground} from 'src/pages/classic/ContentBar.vue'
import RomSim, {openEditRomSimulatorDialog} from 'src/pages/classic/modules/RomSim.vue'

const global = useGlobalStore();
const {
  activePlatform,
  activeRom,
  platformUi,
  rombaseAlias,
  config,
  lang,
  simulatorMap,
  logo,
} = storeToRefs(global);
const panel = ref('detail')
const slide = ref(0)
const link = ref(1)
const autoplay = ref(true)
const detail: any = ref(null)
const simulators: any = ref([])
const strategy: any = ref("")
const formattedStrategy = computed(() => formatStrategyContent(strategy.value))
const PlatformDesc = ref("")
const thumbListCarousel: any = ref([]);
const thumbListAlumb: any = ref({});
const strategyFiles: any = ref([])
const issetVideo: any = ref(false) //是否存在视频资料
const scrollAreaDocRef = ref(null)
const platformDescRef = ref(null)
const carouselVideoRef = ref(null)
const albumVideoRef = ref(null)

const root = document.documentElement;

const albumTitle = {
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

const albumIcon: Record<string, string> = {
  "thumb": "image",
  "snap": "photo_camera",
  "poster": "style",
  "packing": "inventory_2",
  "title": "title",
  "cassette": "album",
  "icon": "widgets",
  "gif": "gif_box",
  "background": "wallpaper",
  "video": "videocam",
}

onMounted(() => {
  getPlatformDesc()
})

//监听rom变更
watch(activeRom, (newValue, oldValue) => {
  if (activeRom.value > 0) {

    //滚动条回到最顶端
    if (scrollAreaDocRef.value != null) {

      scrollAreaDocRef.value.setScrollPosition("vertical", 0)
    }

    //读取模拟器信息
    simulators.value = []
    if (!isEmpty(simulatorMap.value[activePlatform.value])) {
      simulators.value = simulatorMap.value[activePlatform.value]
    }

    //读取rom信息
    GetGameDetail(activeRom.value).then((result: string) => {

      let resp = decodeApiData(result)
      if (resp.err == "") {
        detail.value = resp.data

        //设置默认模拟器选项
        if (!isEmpty(simulators.value)) {
          let simId = 0
          simulators.value.forEach((item: any, index: number) => {
            if (item.Id == detail.value.Info.SimId) {
              simId = item.Id
            }
          })
          detail.value.Info.SimId = simId > 0 ? simId : simulators.value[0].Id
        }
      }
    })

    //读取攻略
    GetGameStrategy(activeRom.value).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err == "") {
        strategy.value = resp.data
      }
    })

    //读取攻略文件
    GetGameStrategyFiles(activeRom.value).then((result: string) => {
      let resp = decodeApiData(result)
      strategyFiles.value = resp.data
    })

    //读取相册
    GetGameThumbs(activeRom.value, "").then((result: string) => {
      let resp = decodeApiData(result)
      //整理图片数据
      handleCarouselAndAlumb(resp.data);

      //更换背景图
      if (thumbListAlumb.value["background"] != undefined && thumbListAlumb.value["background"].length > 0) {
        changeContentBackground(thumbListAlumb.value["background"][0].Path)
      } else {
        changeContentBackground("")
      }
    })
  }
})

//监听平台变更
watch(activePlatform, (newValue, oldValue) => {
  getPlatformDesc()
});

//加载平台描述
function getPlatformDesc() {
  //滚动条回到最顶端
  PlatformDesc.value = ""
  if (activePlatform.value > 0) {
    GetPlatformById(activePlatform.value).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err == "") {
        PlatformDesc.value = resp.data.Desc;
      }
    })
  }
}

//图片数据整理
function handleCarouselAndAlumb(data) {
  thumbListCarousel.value = []
  issetVideo.value = false
  slide.value = 0
  thumbListAlumb.value = {}
  data.forEach((val: any) => {
    if (thumbListAlumb.value[val.Type] == undefined) {
      thumbListAlumb.value[val.Type] = [];
    }
    thumbListAlumb.value[val.Type].push(val);
  })

  //处理滑动图
  if (thumbListAlumb.value["video"] != undefined && thumbListAlumb.value["video"].length > 0) {
    issetVideo.value = true
    thumbListCarousel.value = thumbListAlumb.value["video"]
  }

  config.value.ThumbOrders.forEach((typ: string) => {
    if (typ == "video" || thumbListAlumb.value[typ] == undefined) {
      return
    }
    thumbListCarousel.value = thumbListCarousel.value.concat(thumbListAlumb.value[typ])
  })
}

//显示相册
function showViewer(romIndex: number) {

  if (thumbListCarousel.value.length == 0) {
    return
  }

  let views: any = []
  let index = 0
  let tmpIdx = 0
  thumbListCarousel.value.forEach((val: any) => {
    if (val.Type == "video") {
      return
    }
    views.push(val.Path)
    if (romIndex == val.Index) {
      index = tmpIdx
    }
    tmpIdx++
  })

  if (views.length > 0) {
    api({
      images: views,
      options: {
        initialViewIndex: index,
        title: [1, (image, imageData) => `${imageData.naturalWidth} × ${imageData.naturalHeight}`],
      }
    })
  }

}

//更新评分
function updateRating() {
  console.log(detail.value.Info.Score)
  UpdateGameScore(detail.value.Info.Id, detail.value.Info.Score).then((result: string) => {
  })

}

//更新通关状态
function updateComplete() {
  if (detail.value.Info.Complete >= 2) {
    detail.value.Info.Complete = 0
  } else {
    detail.value.Info.Complete++
  }

  UpdateGameComplete(detail.value.Info.Id, detail.value.Info.Complete).then((result: string) => {
  })
}

function UpdateRomSimId() {
  SetRomSimId(detail.value.Info.Id, detail.value.Info.SimId).then((result: string) => {
  })
}

function runGame(id: number) {
  RunGame(id, detail.value.Info.SimId).then((result: string) => {
  })
}

//运行文件
function runProgram(path: string) {
  RunProgram(path).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return;
    }
  })
}

function formatStrategyContent(content: string) {
  if (!content) {
    return ""
  }
  return content
    .replace(/\[([^\]]+)\]/g, '<div class="strategy-block-title">$1</div>')
    .replace(/(^|<br\s*\/?>)\s*--\s*(<br\s*\/?>|$)/gi, '$1<div class="strategy-block-divider"></div>$2')
}

//修改视频音量
function changeVideoVolume(obj) {
  const video = obj.target;
  const volume = video.volume;
  if (config.value.VideoVolume == volume) {
    console.log("changeVideoVolume 没有修改")
    return
  }
  console.log("changeVideoVolume", volume)
  UpdateOneConfig("VideoVolume", volume.toString()).then((result: string) => {
    config.value.VideoVolume = volume.toString()
  })
}

//窗口失去焦点，暂停视频
window.addEventListener('blur', function () {
  if (carouselVideoRef.value != null) {
    carouselVideoRef.value.forEach((item) => {
      item.pause()
    });
  }

  if (albumVideoRef.value != null) {
    albumVideoRef.value.forEach((item) => {
      item.pause()
    });
  }
});

</script>

<style scoped>
@import "src/css/classic/rightBar.css";

</style>
