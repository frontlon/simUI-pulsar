<template>
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 70%;min-height:500px;" v-if="rombase && detail">
      <q-card-section class="row items-center q-pb-none q-mb-md">
        <div class="text-h6">{{ lang.editRombase }} {{ detail.Info.Name }}</div>
        <q-space/>
        <q-btn icon="close" flat round v-close-popup/>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <q-tabs v-model="tab" class="text-grey" active-color="primary" indicator-color="primary" align="justify">
          <q-tab name="base" :label="lang.baseInfo"/>
          <q-tab name="desc" :label="lang.desc"/>
          <q-tab name="strategy" :label="lang.strategy"/>
          <q-tab name="files" :label="lang.strategyFiles"/>
          <!--          <q-tab name="audio" label="音乐原声"/>-->
        </q-tabs>

        <q-separator/>
        <q-tab-panels v-model="tab" keep-alive animated>
          <!-- 基本信息 -->
          <q-tab-panel name="base">
            <q-list v-if="rombase && enumMap">
              <q-item>
                <q-item-section>
                  <q-input v-if="detail" filled square dense :label="lang.alias" v-model="detail.Info.Name"/>
                </q-item-section>
                <q-item-section>
                  <q-input filled square dense :label="lang.baseNameEn" v-model="rombase.NameEN"/>
                </q-item-section>
                <q-item-section>
                  <q-input filled square dense :label="lang.baseNameJp" v-model="rombase.NameJP"/>
                </q-item-section>
              </q-item>
              <q-item>
                <q-item-section>
                  <q-select filled square dense options-dense use-input fill-input hide-selected :label="lang.baseType"
                            :options="enumMap.type"
                            v-model="rombase.Type"
                            popup-content-style="height:40vh"
                            @input-value="(value) => {rombase.Type = value}"/>
                </q-item-section>
                <q-item-section>
                  <q-select filled square dense options-dense use-input fill-input hide-selected :label="lang.baseYear"
                            popup-content-style="height:40vh"
                            :options="enumMap.year" v-model="rombase.Year"
                            @input-value="(value) => {rombase.Year = value}"/>
                </q-item-section>
              </q-item>
              <q-item>
                <q-item-section>
                  <q-select filled square dense options-dense use-input fill-input hide-selected
                            :label="lang.baseProducer"
                            popup-content-style="height:40vh"
                            :options="enumMap.producer" v-model="rombase.Producer"
                            @input-value="(value) => {rombase.Producer = value}"/>
                </q-item-section>
                <q-item-section>
                  <q-select filled square dense options-dense use-input fill-input hide-selected
                            :label="lang.basePublisher"
                            popup-content-style="height:40vh"
                            :options="enumMap.publisher" v-model="rombase.Publisher"
                            @input-value="(value) => {rombase.Publisher = value}"/>
                </q-item-section>
              </q-item>
              <q-item>
                <q-item-section>
                  <q-select filled square dense options-dense use-input fill-input hide-selected
                            :label="lang.baseCountry"
                            popup-content-style="height:40vh"
                            :options="enumMap.country" v-model="rombase.Country"
                            @input-value="(value) => {rombase.Country = value}"/>
                </q-item-section>
                <q-item-section>
                  <q-select filled square dense options-dense use-input fill-input hide-selected
                            :label="lang.baseVersion"
                            popup-content-style="height:40vh"
                            :options="enumMap.version" v-model="rombase.Version"
                            @input-value="(value) => {rombase.Version = value}"/>
                </q-item-section>
              </q-item>
              <q-item>
                <q-item-section>
                  <q-select filled square dense options-dense use-input fill-input hide-selected
                            :label="lang.baseTranslate"
                            popup-content-style="height:40vh"
                            :options="enumMap.translate" v-model="rombase.Translate"
                            @input-value="(value) => {rombase.Translate = value}"/>
                </q-item-section>
                <q-item-section>
                  <q-select filled square dense options-dense :label="lang.rating"
                            :options="optionScore" v-model="rombase.Score"
                            popup-content-style="height:40vh"
                            @input-value="(value) => {rombase.Score = value}"/>
                </q-item-section>
              </q-item>
              <div class="separator"></div>
              <q-item>
                <q-item-section>
                  <q-input filled square dense v-if="aliasMap" :label="aliasMap.OtherA" v-model="rombase.OtherA"/>
                </q-item-section>
                <q-item-section>
                  <q-input filled square dense v-if="aliasMap" :label="aliasMap.OtherB" v-model="rombase.OtherB"/>
                </q-item-section>
              </q-item>
              <q-item>
                <q-item-section>
                  <q-input filled square dense v-if="aliasMap" :label="aliasMap.OtherC" v-model="rombase.OtherC"/>
                </q-item-section>
                <q-item-section>
                  <q-input filled square dense v-if="aliasMap" :label="aliasMap.OtherD" v-model="rombase.OtherD"/>
                </q-item-section>
              </q-item>
            </q-list>
            <q-card-actions align="right" class="text-primary">
              <q-btn color="primary" square size="md" :label="lang.update" class="update-btn" @click="updateRombase()"/>
            </q-card-actions>
          </q-tab-panel>

          <!-- 游戏简介 -->
          <q-tab-panel name="desc">
            <q-input square filled autogrow v-model="detail.DocContent" type="textarea" class="mt"/>
            <q-card-actions align="right" class="text-primary">
              <q-btn color="primary" square size="md" :label="lang.update" class="update-btn" @click="setDoc()"/>
            </q-card-actions>
          </q-tab-panel>

          <!-- 游戏攻略 -->
          <q-tab-panel name="strategy">
            <q-banner dense class="bg-grey-9 text-white">
              <div class="editor-desc">
                <q-badge outline align="middle" color="grey" label="ctrl" class="badge"/>
                <span>+</span>
                <q-badge outline align="middle" color="grey" label="shift" class="badge"/>
                <span>+</span>
                <q-badge outline align="middle" color="grey" label="v" class="badge"/>
                <span>{{ lang.pastePlainText }}</span>
              </div>
            </q-banner>
            <q-editor class="editor" square v-model="strategy" :definitions="editorDefinitions"
                      :toolbar="EDITOR_TOOLBAR"/>
            <q-card-actions align="right" class="text-primary">
              <q-btn color="primary" square size="md" :label="lang.update" class="update-btn" @click="setStrategy()"/>
            </q-card-actions>
          </q-tab-panel>

          <!-- 攻略文件 -->
          <q-tab-panel name="files">
            <q-list>
              <q-item v-for="(item,index) in strategyFiles">
                <q-item-section>
                  <q-input filled square dense :label="lang.filename" v-model="item.Name"/>
                </q-item-section>
                <q-item-section>
                  <q-input readonly filled square dense :label="lang.filePath" v-model="item.Path">
                    <template v-slot:append>
                      <q-icon name="upload_file" class="open-dialog" @click="openFileDialog(index)"/>
                    </template>
                  </q-input>
                </q-item-section>
                <q-item-section side>
                  <q-btn flat square dense icon="delete" @click="delFile(index)">
                    <q-tooltip>{{ lang.delete }}</q-tooltip>
                  </q-btn>
                </q-item-section>
              </q-item>
            </q-list>
            <div class="text-center">
              <q-btn :label="lang.add" unelevated color="secondary" class="q-mt-md" @click="addFile()"></q-btn>
            </div>
            <q-card-actions align="right" class="text-primary">
              <q-btn color="primary" square size="md" :label="lang.update" class="update-btn" @click="updateFile()"/>
            </q-card-actions>

          </q-tab-panel>

          <!-- 音乐原声 -->
          <!--          <q-tab-panel name="audio">
                      <q-list>
                        <q-item>
                          <q-item-section>
                            <q-input filled square dense label="音乐名称"/>
                          </q-item-section>
                          <q-item-section>
                            <q-file filled square dense label="音乐文件">
                              <template v-slot:append>
                                <q-icon name="music_note"/>
                              </template>
                            </q-file>
                          </q-item-section>
                          <q-item-section side>
                            <q-btn flat square dense icon="delete"/>
                          </q-item-section>
                        </q-item>
                        <q-item>
                          <q-item-section>
                            <q-input filled square dense label="音乐名称"/>
                          </q-item-section>
                          <q-item-section>
                            <q-file filled square dense label="音乐文件">
                              <template v-slot:append>
                                <q-icon name="music_note"/>
                              </template>
                            </q-file>
                          </q-item-section>
                          <q-item-section side>
                            <q-btn flat square dense icon="delete"/>
                          </q-item-section>
                        </q-item>
                      </q-list>
                      <q-card-actions align="right" class="text-primary">
                        <q-btn flat label="取消" v-close-popup/>
                        <q-btn flat label="更新"/>
                      </q-card-actions>
                    </q-tab-panel>-->

        </q-tab-panels>

      </q-card-section>
    </q-card>

  </q-dialog>
</template>

<script lang="ts">

import {ref} from 'vue'
import {decodeApiData, getFileName, isEmpty, notify} from "components/utils";
import {useRoute} from "vue-router";
import {
  GetGameDetail,
  GetGameStrategy,
  GetGameStrategyFiles,
  GetRomBase,
  GetRomBaseAlias,
  GetRomBaseEnum,
  OpenFileDialog,
  OpenFileDialogForEditor,
  SetRomBase,
  UpdateGameStrategy,
  UpdateStrategyFiles,
} from 'app/wailsjs/go/controller/Controller'
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";
import {EDITOR_TOOLBAR} from "boot/constant";
import {callback} from 'pages/classic/context/Context.vue';

const global = useGlobalStore();
const {lang} = storeToRefs(global);
const route = useRoute();
const tab = ref('base')
const romId: any = ref(null)
const romIndex: any = ref(0)
const detail: any = ref(null)
const rombase: any = ref(null)
const strategy: any = ref("")
const enumMap: any = ref(null)
const showDialog: any = ref(null)
const aliasMap: any = ref(null)
const optionScore: any = ref(["1", "1.5", "2", "2.5", "3", "3.5", "4", "4.5", "5"])
const strategyFiles: any = ref([])

//编辑器插入图片控件
const editorDefinitions = {
  image: {
    tip: 'insert image',
    icon: 'image',
    handler: editorUpdateImage
  }
}

//编辑器上传图片
function editorUpdateImage() {
  OpenFileDialogForEditor().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.data != "") {
      strategy.value += `<img src="` + resp.data + `" alt="">`
    }
  })
}

export function openBaseInfoDialog(id: any, index: number) {
  romId.value = parseInt(id);
  romIndex.value = index;
  showDialog.value = true

  //读取资料
  GetRomBase(romId.value).then((result: string) => {
    let resp = decodeApiData(result)
    rombase.value = resp.data
    if (rombase.value == null) {
      rombase.value = {};
      rombase.value.NameEN = "";
      rombase.value.NameJP = "";
      rombase.value.Type = "";
      rombase.value.Year = "";
      rombase.value.Producer = "";
      rombase.value.Publisher = "";
      rombase.value.Country = "";
      rombase.value.Translate = "";
      rombase.value.Version = "";
      rombase.value.OtherA = "";
      rombase.value.OtherB = "";
      rombase.value.OtherC = "";
      rombase.value.OtherD = "";
    }
  })

  //读取rom信息
  GetGameDetail(romId.value).then((result: string) => {
    let resp = decodeApiData(result)
    detail.value = resp.data

    //读取资料项别名
    GetRomBaseAlias(detail.value.Info.Platform).then((result: string) => {
      let resp = decodeApiData(result)
      aliasMap.value = resp.data
    })
  })

  //读取攻略
  GetGameStrategy(romId.value).then((result: string) => {
    let resp = decodeApiData(result)
    strategy.value = resp.data
  })

  //读取攻略文件
  GetGameStrategyFiles(romId.value).then((result: string) => {
    let resp = decodeApiData(result)
    strategyFiles.value = resp.data
  })


  //读取相册
  /*GetGameThumbs( romId.value).then((result: string) => {
    let resp = DecodeApiData(result)
    thumbListCarousel.value = resp.data
  })*/

  //选项枚举
  GetRomBaseEnum().then((result: string) => {
    let resp = decodeApiData(result)
    enumMap.value = resp.data
  })

}

//更新资料
function updateRombase() {
  let data = JSON.stringify(rombase.value)

  SetRomBase(romId.value, data, detail.value.Info.Name).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
      callback("editBaseInfo", romIndex.value, detail.value.Info.Name)
    } else {
      notify("err", lang.updateFail + ":" + resp.err)
    }
  })
}

function setDoc() {
  UpdateGameStrategy(romId.value, "doc", detail.value.DocContent).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
    } else {
      notify("err", lang.updateFail + ":" + resp.err)
    }
  })
}

function setStrategy() {
  UpdateGameStrategy(romId.value, "strategy", strategy.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
    } else {
      notify("err", lang.updateFail + ":" + resp.err)
    }
  })
}

//选择文件
function openFileDialog(index: number) {
  OpenFileDialog("").then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {
      strategyFiles.value[index].Path = resp.data;
      strategyFiles.value[index].Name = getFileName(resp.data)
    }
  })
}

//删除攻略文件
function delFile(index: number) {
  strategyFiles.value.splice(index, 1)
}

//添加攻略文件
function addFile() {
  strategyFiles.value.push({})
}

//更新攻略文件
function updateFile() {
  UpdateStrategyFiles(romId.value, JSON.stringify(strategyFiles.value)).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      strategyFiles.value = resp.data
      notify("suc", lang.value.updateSuc)
    } else {
      notify("err", resp.err)
    }
  })
}

export default {
  methods: {isEmpty},
  setup() {
    return {
      route,
      tab,
      romId,
      detail,
      rombase,
      strategy,
      enumMap,
      showDialog,
      aliasMap,
      optionScore,
      editorDefinitions,
      lang,
      strategyFiles,
      EDITOR_TOOLBAR,
      editorUpdateImage,
      updateRombase,
      setDoc,
      setStrategy,
      openFileDialog,
      delFile,
      addFile,
      updateFile,
    };
  }
}

</script>

<style scoped>
@import "src/css/manage.css";
</style>
