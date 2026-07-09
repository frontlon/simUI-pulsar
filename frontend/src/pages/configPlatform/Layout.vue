<template>
  <div class="manage-wrapper platform-config">
    <q-item class="page-header" style="--wails-draggable:drag">
      <q-item-section avatar @click="$router.back()">
        <q-btn flat size="lg" icon="arrow_back" class="vertical-middle"/>
      </q-item-section>
      <q-item-section>
        <h3>{{ lang.platformConfig }}</h3>
      </q-item-section>
    </q-item>
    <div class="row platform-layout">

      <div class="col-auto left-wrapper">
        <div class="sidebar-inner">
          <q-btn-group spread class="sidebar-actions">
            <q-btn unelevated :label="lang.addPlatform" @click="addPlatform()"/>
            <q-btn unelevated :label="lang.delPlatform" @click="deletePlatform()"/>
          </q-btn-group>

          <q-scroll-area class="platform-scroll" :visible="false">
            <q-list class="platform-list">
              <q-item v-for="(item,index) in platformList" clickable v-ripple active-class="platform-active"
                      :key="item.Id" :active="activePlatformId === item.Id" @click="changePlatform(item.Id,index)">
                <q-item-section>{{ item.Name }}</q-item-section>
              </q-item>
            </q-list>
          </q-scroll-area>
        </div>
      </div>
      <div class="col right-wrapper" v-if="activePlatformInfo">
        <q-tabs v-model="tab" class="config-tabs text-grey" active-color="primary" indicator-color="primary"
                align="left">
          <q-tab name="base" :label="lang.platformConfig" tabindex="1"/>
          <q-tab name="simulator" :label="lang.simConfig" tabindex="2"/>
          <q-tab name="desc" :label="lang.platformDesc" tabindex="3"/>
          <q-tab name="rombase" :label="lang.rombaseAlias" tabindex="4"/>
        </q-tabs>
        <q-separator/>
        <q-tab-panels v-model="tab" keep-alive animated>
          <!-- 基本信息 -->
          <q-tab-panel name="base" class="config-panel">
            <q-card flat bordered class="config-section">
              <q-card-section>
                <div class="section-title">{{ lang.platformBaseInfo }}</div>
                <div class="row q-col-gutter-md">
                  <div class="col-12 col-sm-6">
                    <q-input filled square :label="lang.platformName" v-model="activePlatformInfo.Name"/>
                  </div>
                  <div class="col-12 col-sm-6">
                    <q-input filled square :label="lang.label" v-model="activePlatformInfo.Tag"/>
                  </div>
                  <div class="col-12 col-sm-6">
                    <q-input filled square :label="lang.platformIcon" class="open-dialog"
                             v-model="activePlatformInfo.Icon">
                      <template v-slot:append>
                        <q-icon name="upload_file" @click="openFileDialog()"/>
                      </template>
                    </q-input>
                  </div>
                  <div class="col-12 col-sm-6">
                    <q-select filled square map-options emit-value :options="boolOptions" :label="lang.hideName"
                              v-model="activePlatformInfo.HideName">
                      <q-tooltip>{{ lang.showHidePlatform }}</q-tooltip>
                    </q-select>
                  </div>
                </div>
              </q-card-section>
            </q-card>

            <q-card flat bordered class="config-section">
              <q-card-section>
                <div class="section-title">{{ lang.menuConfig }}</div>
                <div class="row q-col-gutter-md">
                  <div class="col-12">
                    <q-input filled square :label="lang.platformResDir" class="open-dialog"
                             v-model="activePlatformInfo.RootPath">
                      <template v-slot:append>
                        <q-icon name="create_new_folder" @click="openDirectoryDialog('root')"/>
                      </template>
                    </q-input>
                  </div>
                  <div class="col-12">
                    <q-input filled square :label="lang.romDir" class="open-dialog"
                             v-model="activePlatformInfo.RomPath">
                      <template v-slot:append>
                        <q-icon name="create_new_folder" @click="openDirectoryDialog('rom')"/>
                      </template>
                      <q-tooltip>{{ lang.tipRomDir }}</q-tooltip>
                    </q-input>
                  </div>
                </div>
              </q-card-section>
            </q-card>

            <q-card flat bordered class="config-section">
              <q-card-section>
                <div class="section-title">{{ lang.romExt }}</div>
                <div class="row q-col-gutter-md items-center">
                  <div class="col">
                    <q-input filled square :label="lang.romExt"
                             v-model="addExtInput" @keydown.enter="addExt">
                      <template v-slot:append>
                        <q-btn round flat icon="add" @click="addExt"/>
                      </template>
                      <q-tooltip>{{ lang.tipRomExt }}</q-tooltip>
                    </q-input>
                  </div>
                  <div class="col-auto ext-actions">
                    <q-btn unelevated icon="assignment" @click="addNoExt()">
                      <q-tooltip>添加无扩展名支持</q-tooltip>
                    </q-btn>
                    <q-btn flat icon="delete" @click="clearRomExts()">
                      <q-tooltip>删除所有类型</q-tooltip>
                    </q-btn>
                  </div>
                </div>
                <div class="tag-label-wrapper q-gutter-sm q-mt-md">
                  <q-chip square removable v-for="(ext,index) in activePlatformInfo.RomExts"
                          @remove="removeExt(ext,index)">
                    {{ ext }}
                  </q-chip>
                </div>
              </q-card-section>
            </q-card>

            <div class="bottom-bar">
              <q-btn unelevated :label="lang.update" color="primary" class="update-btn" @click="updatePlatform()"/>
            </div>
          </q-tab-panel>

          <!-- 模拟器管理 -->
          <q-tab-panel name="simulator" class="config-panel">
            <simulators :platform="activePlatformId"/>
          </q-tab-panel>

          <!-- 平台简介 -->
          <q-tab-panel name="desc" class="config-panel">
            <q-card flat bordered class="config-section editor-section">
              <q-card-section class="editor-tip">
                <div class="editor-desc">
                  <q-badge outline align="middle" color="grey" label="ctrl" class="badge"/>
                  <span>+</span>
                  <q-badge outline align="middle" color="grey" label="shift" class="badge"/>
                  <span>+</span>
                  <q-badge outline align="middle" color="grey" label="v" class="badge"/>
                  <span>{{ lang.pastePlainText }}</span>
                </div>
              </q-card-section>
              <q-separator/>
              <q-card-section class="q-pa-none">
                <q-editor class="editor" square v-model="activePlatformInfo.Desc"
                          :definitions="editorDefinitions" :toolbar="EDITOR_TOOLBAR"/>
              </q-card-section>
            </q-card>
            <div class="bottom-bar">
              <q-btn unelevated :label="lang.update" color="primary" class="update-btn" @click="updateDesc()"/>
            </div>
          </q-tab-panel>
          <q-tab-panel name="rombase" class="config-panel">
            <rombase-alias :platform="activePlatformId"/>
          </q-tab-panel>
        </q-tab-panels>
      </div>
      <div v-else class="col right-wrapper relative-position">
        <div class="empty-state absolute-full flex flex-center column">
          <q-icon name="videogame_asset" size="64px" color="grey-7"/>
          <div class="empty-text">{{ lang.tipSelectPlatform }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

import {onMounted, ref} from 'vue'
import {Dialog} from 'quasar'
import {getPromptOpts} from 'components/dialog'
import {decodeApiData, deepClone, notify} from 'components/utils'
import Simulators from 'src/pages/configPlatform/Simulators.vue'
import RombaseAlias from 'src/pages/configPlatform/RombaseAlias.vue'
import {EDITOR_TOOLBAR} from "boot/constant";
import Sortable from 'sortablejs';
import {
  AddPlatform,
  DelPlatform,
  GetPlatformOriginal,
  OpenDirectoryDialog,
  OpenFileDialog,
  OpenFileDialogForEditor,
  UpdatePlatform,
  UpdatePlatformDesc,
  UpdatePlatformSort,
} from 'app/wailsjs/go/controller/Controller'
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);

const tab = ref('base')
const addExtInput = ref("")
const activePlatformId = ref(0)
const activePlatformIndex = ref(0)
const activePlatformInfo: any = ref(null)
const platformList: any = ref([])
//编辑器插入图片控件
const editorDefinitions = {
  image: {
    tip: 'insert image',
    icon: 'image',
    handler: editorUpdateImage
  }
}
const boolOptions = [
  {label: '否', value: 0},
  {label: '是', value: 1},
]
onMounted(() => {
  //平台
  platformDropSort();
  //加载平台
  GetPlatformOriginal().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      let jsonArr = resp.data
      jsonArr.forEach((item, index) => {
        jsonArr[index].RomExts = handleExts(item.RomExts);
      });
      platformList.value = jsonArr;
    }
  })
})

//点击平台
function changePlatform(id: number, index: number) {
  activePlatformId.value = id
  activePlatformIndex.value = index
  activePlatformInfo.value = deepClone(platformList.value[index]);
}

//添加平台
function addPlatform() {
  let opt = getPromptOpts(lang.value.inputPlatformName, "", lang.value.ok, false, "")
  Dialog.create(opt).onOk(resp => {
    if (resp.input == "") {
      return
    }
    let data = resp.input

    AddPlatform(data).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err == "") {
        let jsonObj = resp.data
        jsonObj.RomExts = handleExts(jsonObj.RomExts);
        platformList.value.push(jsonObj)
      }
    })
  })
}

//添加无扩展名支持
function addNoExt() {
  let len = activePlatformInfo.value.RomExts.length;
  if (len > 0 && activePlatformInfo.value.RomExts[len - 1] == "noext") {
    return;
  }
  activePlatformInfo.value.RomExts.push("noext");
}

//清空所有rom类型
function clearRomExts() {
  activePlatformInfo.value.RomExts = [];
}

//删除平台
function deletePlatform() {

  if (activePlatformId.value == 0) {
    notify("warn", lang.value.tipSelectPlatform)
    return;
  }

  let opt = getPromptOpts(lang.value.tipDelPlatform, "", lang.value.ok, false)
  Dialog.create(opt).onOk(() => {
    DelPlatform(activePlatformId.value).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err == "") {
        platformList.value.splice(activePlatformIndex.value, 1)
        activePlatformInfo.value = null
        activePlatformId.value = 0
      }
    })
  })
}

//更新平台信息
function updatePlatform() {
  let req = {
    id: activePlatformInfo.value.Id,
    name: activePlatformInfo.value.Name,
    icon: activePlatformInfo.value.Icon,
    tag: activePlatformInfo.value.Tag,
    romExts: activePlatformInfo.value.RomExts,
    rootPath: activePlatformInfo.value.RootPath,
    romPath: activePlatformInfo.value.RomPath.split(';'),
    hideName: activePlatformInfo.value.HideName,
  }

  UpdatePlatform(JSON.stringify(req)).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      activePlatformInfo.value = resp.data
      activePlatformInfo.value.RomExts = handleExts(activePlatformInfo.value.RomExts)
      platformList.value[activePlatformIndex.value] = activePlatformInfo.value
      notify("suc", lang.value.updateSuc)
    }
  })
}

//添加文件类型
function addExt() {
  if (addExtInput.value == "") {
    return
  }

  let exts = activePlatformInfo.value.RomExts
  let arr = addExtInput.value.split(",")
  arr.forEach((ext) => {
    ext = ext.trim()
    if (ext.charAt(0) == '.') {
      ext = ext.substring(1)
    }
    exts.push(ext);
  })

  exts = Array.from(new Set(exts));
  activePlatformInfo.value.RomExts = exts


  addExtInput.value = "";
}

//删除文件类型
function removeExt(ext: string, index: number) {
  activePlatformInfo.value.RomExts.splice(index, 1);
}

//选择文件
function openFileDialog() {

  OpenFileDialog("image").then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {
      activePlatformInfo.value.Icon = resp.data;
    }
  })

}

//选择目录
function openDirectoryDialog(opt: string) {
  console.log("openDirectoryDialog opt: ", opt);

  OpenDirectoryDialog().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      let pth = resp.data
      console.log("openDirectoryDialog pth: ", pth);
      activePlatformInfo.value.RomPath = activePlatformInfo.value.RomPath.trim();
      if (opt == "rom") {
        if (activePlatformInfo.value.RomPath != "") {
          pth = ";" + pth
        }
        activePlatformInfo.value.RomPath += pth.trim();
      } else {
        activePlatformInfo.value.RootPath = pth.trim();
        activePlatformInfo.value.RomPath = activePlatformInfo.value.RootPath + "/roms"
      }
    }
  })
}

//更新平台简介
function updateDesc() {

  UpdatePlatformDesc(activePlatformInfo.value.Id, activePlatformInfo.value.Desc).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      platformList.value[activePlatformIndex.value].Desc = activePlatformInfo.value.Desc
      notify("suc", lang.value.updateSuc)
    }
  })
}

//编辑器上传图片
function editorUpdateImage() {
  OpenFileDialogForEditor().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      activePlatformInfo.value.Desc += `<img src="` + resp.data + `">`
    }
  })
}

//扩展名替换
function handleExts(extstr: string) {
  if (extstr == "") {
    return [];
  }
  let exts = extstr.split(',');
  for (let i = 0; i < exts.length; i++) {
    exts[i] = exts[i].replace('.', '');
  }
  return exts;
}

//平台拖拽排序
function platformDropSort() {
  const el = document.querySelectorAll('.platform-list')[0] as HTMLElement;
  Sortable.create(el, {
    animation: 150,
    sort: true,
    disabled: false,
    // 结束拖拽后的回调函数
    onEnd: (evt: any) => {
      let currentRow = platformList.value.splice(evt.oldIndex, 1)[0];
      platformList.value.splice(evt.newIndex, 0, currentRow);
      //更新数据
      let req: any = [];
      platformList.value.forEach((item: any, index: number) => {
        platformList.value[index].Sort = index + 1
        req.push(item.Id);
        if (activePlatformInfo.value && activePlatformInfo.value.Id == item.Id) {
          activePlatformInfo.Sort = item.Sort
          activePlatformIndex.value = index
        }
      });

      UpdatePlatformSort(JSON.stringify(req)).then((result: string) => {
      })
    },
  });
}

</script>
<style scoped>
@import "src/css/manage.css";
@import "src/css/page/platform.css";
</style>
