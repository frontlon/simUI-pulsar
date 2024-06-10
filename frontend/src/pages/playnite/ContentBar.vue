<template>
  <div class="rightbar-wrapper q-gutter-y-md">
    <div v-if="activeRom && !isEmpty(detail)">

      <div class="bg-wrapper">
        <div class="bg-mask"></div>

        <q-slider v-if="thumbListAlumb['video'] && thumbListAlumb['video'].length > 0"
                  v-model="videoVolumeSlider" class="bg-video-slider" :min="0" :max="100" @change="changeVideoVolume"/>
        <video v-if="thumbListAlumb['video'] && thumbListAlumb['video'].length > 0"
               ref="videoRef"
               loop autoplay :volume="config.VideoVolume" key="0" class="bg-video"
               :src="thumbListAlumb['video'][0].Path"/>
      </div>

      <div class="content-wrapper">

        <div class="title-wrapper">
          <h1>{{ platformUi.NameType == 1 ? detail.Info.Name : detail.Info.RomName }}</h1>

          <q-carousel v-if="thumbListCarousel.length > 0 && platformUi.HideCarousel == 0"
                      animated infinite arrows :autoplay="issetVideo ? 0 : 2500" v-model="slide"
                      transition-prev="jump-left" transition-next="jump-right"
                      @mouseenter="autoplay = false" @mouseleave="autoplay = true"
                      class="right-carousel">
            <q-carousel-slide v-for="(item,index) in thumbListCarousel" :name="index" style="padding: 0"
                              :img-src="item.Type == 'video' ? '' : item.Path"
                              @click="item.Type != 'video' ?showViewer(item.Index) : null">
            </q-carousel-slide>
          </q-carousel>

        </div>

        <div class="row content">
          <div class="col-4" style="padding-left: 30px">

            <!--启动按钮-->
            <div class="run-btn-list">
              <q-btn class="btn-primary" no-caps :label="lang.runGame" @click="openRunGameDialog()"/>
            </div>

            <div class="rating">
              <star-rating v-model:rating="detail.Info.Score" :increment="0.5" :read-only="true" :star-size="18"
                           :show-rating="false" active-color="#ff9800" inactive-color="#784f1f" :clearable="true"
                           :star-points="[23,2, 14,17, 0,19, 10,34, 7,50, 23,43, 38,50, 36,34, 46,19, 31,17]"/>
            </div>

            <div class="rombase">
              <div class="rombase-item">
                <p class="title">{{ lang.complete }}</p>
                <p v-if="detail.Info.Complete == 1">{{ lang.passed }}</p>
                <p v-else-if="detail.Info.Complete == 2">{{ lang.perfectClear }}</p>
                <p v-else>{{ lang.notCleared }}</p>
              </div>

              <div class="rombase-item">
                <p class="title">{{ lang.runNum }}</p>
                <p>{{ detail.Info.RunNum }}</p>
              </div>

              <div class="rombase-item" style="width: 100%">
                <p class="title">{{ lang.lastRunTime }}</p>
                <p>{{ !isEmpty(detail.Info.RunLasttime) ? getDateTime(detail.Info.RunLasttime) : lang.notRun }}</p>
              </div>
            </div>
          </div>
          <div class="col-4" style="padding: 0 30px;">
            <!--资料-->
            <div class="right-title">{{ lang.baseInfo }}</div>
            <div class="rombase">
              <div class="rombase-item">
                <p class="title">{{ lang.baseType }}</p>
                <p v-if="!isEmpty(detail.Info.BaseType)">{{ detail.Info.BaseType }}</p>
                <p class="no-set" v-else>{{ lang.noSet }}</p>
              </div>
              <div class="rombase-item">
                <p class="title">{{ lang.baseYear }}</p>
                <p v-if="!isEmpty(detail.Info.BaseYear)">{{ detail.Info.BaseYear }}</p>
                <p class="no-set" v-else>{{ lang.noSet }}</p>
              </div>
              <div class="rombase-item">
                <p class="title">{{ lang.baseProducer }}</p>
                <p v-if="!isEmpty(detail.Info.BaseProducer)">{{ detail.Info.BaseProducer }}</p>
                <p class="no-set" v-else>{{ lang.noSet }}</p>
              </div>
              <div class="rombase-item">
                <p class="title">{{ lang.basePublisher }}</p>
                <p v-if="!isEmpty(detail.Info.BasePublisher)">{{ detail.Info.BasePublisher }}</p>
                <p class="no-set" v-else>{{ lang.noSet }}</p>
              </div>
              <div class="rombase-item">
                <p class="title">{{ lang.baseCountry }}</p>
                <p v-if="!isEmpty(detail.Info.BaseCountry)">{{ detail.Info.BaseCountry }}</p>
                <p class="no-set" v-else>{{ lang.noSet }}</p>
              </div>
              <div class="rombase-item">
                <p class="title">{{ lang.baseTranslate }}</p>
                <p v-if="!isEmpty(detail.Info.BaseTranslate)">{{ detail.Info.BaseTranslate }}</p>
                <p class="no-set" v-else>{{ lang.noSet }}</p>
              </div>
              <div class="rombase-item">
                <p class="title">{{ lang.baseVersion }}</p>
                <p v-if="!isEmpty(detail.Info.BaseVersion)">{{ detail.Info.BaseVersion }}</p>
                <p class="no-set" v-else>{{ lang.noSet }}</p>
              </div>

            </div>

          </div>
          <div class="col-4" style="padding-right: 30px">
            <div>
              <div class="right-title">{{ lang.desc }}</div>
              <q-scroll-area class="doc" ref="scrollAreaDocRef">
                <div v-if="detail.DocContent" class="selectable" v-html="detail.DocContent">
                </div>
                <div v-else class="doc-empty">
                  {{ lang.GameDocIsNotExists }}
                </div>
              </q-scroll-area>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else-if="activePlatform > 0 && activePlatformDesc" v-html="activePlatformDesc"
         ref="platformDescRef" class="panel-desc">
    </div>
    <div v-else class="panel-default">
    </div>
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
import {ref, watch} from 'vue'
import 'viewerjs/dist/viewer.css'
import {api} from 'v-viewer'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import StarRating from 'vue-star-rating'
import {decodeApiData, getDateTime, isEmpty} from "components/utils";
import {
  GetGameDetail,
  GetGameThumbs,
  RunGame,
  SetRomSimId,
  UpdateOneConfig,
} from 'app/wailsjs/go/controller/Controller'
import {changeContentBackground} from 'src/pages/playnite/RomListBar.vue'

const global = useGlobalStore();
const {
  activePlatform,
  activeRom,
  platformUi,
  rombaseAlias,
  config,
  lang,
  simulatorMap,
  activePlatformDesc,
  activeFocus,
} = storeToRefs(global);
const slide = ref(0)
const detail: any = ref(null)
const simulators: any = ref([])
const thumbListCarousel: any = ref([]);
const thumbListAlumb: any = ref({});
const issetVideo: any = ref(false) //是否存在视频资料
const scrollAreaDocRef = ref(null)
const platformDescRef = ref(null)
const videoRef = ref(null)
const showGameDialog: any = ref(false)
const showSimDialog: any = ref(false)
const selectGameId: any = ref(null)
const videoVolumeSlider = ref(null)
let tmpRomFocus = null

//监听rom变更
watch(activeRom, (newValue, oldValue) => {
  if (activeRom.value > 0) {

    //读取模拟器信息
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

    //读取相册
    GetGameThumbs(activeRom.value, "").then((result: string) => {
      let resp = decodeApiData(result)

      //初始化视频音量
      videoVolumeSlider.value = config.value.VideoVolume * 100

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

  config.value.ThumbOrders.forEach((typ: string) => {
    if (typ == "video" || typ == "icon" || thumbListAlumb.value[typ] == undefined) {
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
        initialViewIndex: index
      }
    })
  }

}

//打开游戏选择对话框
function openRunGameDialog() {

  tmpRomFocus = activeFocus.value

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


//打开模拟器选择对话框
function openSimDialog(romId) {
  selectGameId.value = romId
  simulators.value = []

  let platform = detail.value.Info.Platform
  if (!isEmpty(simulatorMap.value[platform])) {

    if (simulatorMap.value[platform].length == 1) {
      //运行游戏
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
  }
  showGameDialog.value = false
}

function runGame(simId) {
  SetRomSimId(detail.value.Info.Id, simId)
  RunGame(selectGameId.value, simId)
  showSimDialog.value = false
}

export function openRunGameDialogByKeyboard() {

  if (activeFocus.value[0] != 3) {
    return
  }

  openRunGameDialog()
}

function callbackCloseDialog() {
  activeFocus.value = tmpRomFocus
}

export function closeContentDialog() {
  showGameDialog.value = false
  showSimDialog.value = false
}

//窗口失去焦点，暂停视频
window.addEventListener('blur', function () {
  if (videoRef.value != null) {
    videoRef.value.pause()
  }
});


//修改视频音量
function changeVideoVolume() {
  let volume = videoVolumeSlider.value / 100
  config.value.VideoVolume = volume
  UpdateOneConfig("VideoVolume", volume.toString()).then((result: string) => {
    config.value.VideoVolume = volume.toString()
  })
}

export default {
  components: {StarRating},
  setup() {
    return {
      activeRom,
      detail,
      thumbListAlumb,
      platformUi,
      thumbListCarousel,
      issetVideo,
      slide,
      config,
      lang,
      scrollAreaDocRef,
      activePlatform,
      activePlatformDesc,
      platformDescRef,
      showGameDialog,
      showSimDialog,
      simulators,
      openRunGameDialog,
      openSimDialog,
      getDateTime,
      runGame,
      isEmpty,
      showViewer,
      callbackCloseDialog,
      changeVideoVolume,
      videoRef,
      videoVolumeSlider,
    }
  }
}


</script>

<style scoped>
@import "src/css/playnite/contentBar.css";
</style>
