<template>

  <!-- 平台栏 -->
  <q-scroll-area class="platform-wrapper wrap">

    <q-select standout square :options="tagList" v-model="activeTag" emit-value map-options :label="lang.label"
              dense="dense"
              style="font-size:12px">
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">
            {{ lang.notLabel }}
          </q-item-section>
        </q-item>
      </template>
    </q-select>

    <q-list>
      <q-item clickable v-ripple class="platform-item" active-class="btn-primary platform-active"
              id="platformEle0"
              :active="activePlatform === 0" v-if="activeTag==''" @click="changePlatform(0,0)">
        <q-item-section>
          <q-item-label>{{ lang.all }}</q-item-label>
        </q-item-section>
      </q-item>
      <q-item clickable v-ripple active-class="btn-primary platform-active" class="platform-item"
              v-for="(item,index) in platformList"
              v-show="activeTag == '' || activeTag == item.Tag"
              :id="'platformEle' + (index+1)"
              :active="activePlatform === item.Id" @click="changePlatform(item.Id,index+1)">
        <q-item-section>
          <q-img fit="contain" v-if="item.Icon !== ''" :src="item.Icon" style="width: 100%"/>
          <q-item-label v-if="!item.HideName">{{ item.Name }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>

  </q-scroll-area>

</template>

<script lang="ts">
import {ref, watch} from 'vue'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {decodeApiData, isEmpty} from 'components/utils';
import {
  GetAllPlatformTag,
  GetPlatform,
  GetPlatformUi,
  GetRomBaseAlias,
  UpdateOneConfig
} from "app/wailsjs/go/controller/Controller";
import {changeContentBackground} from 'src/pages/classic/ContentBar.vue'


const global = useGlobalStore();
const {
  activePlatform,
  activeRom,
  activeMenu,
  activeFocus,
  menuLike,
  activeLetter,
  platformUi,
  rombaseAlias,
  config,
  lang,
  scrollAreaRef,
} = storeToRefs(global);

const platformList = ref([])
const tagList = ref([])
const activeTag = ref("")
const root = document.documentElement;

export function initPlatform() {

  //读取平台标签
  tagList.value = [{label: lang.value.all, value: ''}]
  GetAllPlatformTag().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      resp.data.forEach((item: any, index: number) => {
        let s = {label: item, value: item}
        tagList.value.push(s)
      })
    }
  })

  //读取全部平台数据
  platformList.value = []
  GetPlatform().then((result: string) => {
    let resp = decodeApiData(result)
    platformList.value = resp.data;
  })
}

watch(activePlatform, (newValue, oldValue) => {
  changePlatformUi(newValue)
});

//点击切换平台
function changePlatform(id: number, index: number) {
  console.log("changePlatform", id)
  changePlatformUi(id)
  activeFocus.value = [1, index]
  activePlatform.value = id
  activeRom.value = 0
  activeMenu.value = ""
  menuLike.value = 0
  activeLetter.value = "ALL"
  scrollAreaRef.value.setScrollPosition("vertical", 0) //滚动条回到最顶端

  global.clearFilter()
  global.incRomState()
  UpdateOneConfig("Platform", id.toString()).then((result: string) => {
  })
  let menu = JSON.stringify([activeMenu.value, menuLike.value]);
  UpdateOneConfig("Menu", menu).then((result: string) => {
  })
}

//更新平台ui
function changePlatformUi(id: number) {
  console.log("changePlatformUi", id)
  //读取平台ui
  GetPlatformUi(id, "Default").then((result: string) => {
    let resp = decodeApiData(result)
    console.log("changePlatformUi", resp.data)

    platformUi.value = resp.data;

    //更换背景图
    changeContentBackground(platformUi.value.BackgroundImage)
  })

  //读取平台标签别名
  GetRomBaseAlias(id).then((result: string) => {
    let resp = decodeApiData(result)
    rombaseAlias.value = resp.data;
  })
}

//键盘方向键
export function platformEventKeyboard(direction: string) {

  let newIndex = 0;
  let focusType = 1; //1平台

  //焦点不在ROM里
  if (activeFocus.value[0] != focusType) {
    activeFocus.value = [focusType, -1];
  }

  let currFocus = activeFocus.value[1];

  let ele;
  switch (direction) {
    case "ArrowUp":
      if (currFocus == 0) {
        //如果当前在第一个，则跳到最后一个
        let itemCount = document.querySelectorAll('.platform-item').length;
        itemCount--
        ele = document.getElementById('platformEle' + itemCount);
        newIndex = itemCount
      } else {
        newIndex = currFocus - 1
        newIndex = newIndex < 0 ? 0 : newIndex;
        ele = document.getElementById('platformEle' + newIndex);
      }
      break;
    case "ArrowDown":
      newIndex = currFocus + 1
      ele = document.getElementById('platformEle' + newIndex);
      if (isEmpty(ele)) {
        ele = document.getElementById('platformEle0');
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

//事件点击平台
export function clickPlatformEvent() {
  if (activeFocus.value[0] != 1) {
    return
  }
  let ele = document.getElementById('platformEle' + activeFocus.value[1]);
  ele?.click()
}

export default {
  setup() {
    return {
      activePlatform,
      activeRom,
      activeMenu,
      menuLike,
      activeLetter,
      platformUi,
      rombaseAlias,
      config,
      lang,
      platformList,
      tagList,
      activeTag,
      changePlatform,
    }
  },
}
</script>
<style scoped>
@import "src/css/classic/common.css";
@import "src/css/classic/platformBar.css";
</style>
