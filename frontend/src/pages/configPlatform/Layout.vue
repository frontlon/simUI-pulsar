<template>
  <div class="manage-wrapper">
    <q-item style="--wails-draggable:drag">
      <q-item-section avatar @click="$router.back()">
        <q-btn flat size="lg" icon="arrow_back" class="vertical-middle"/>
      </q-item-section>
      <q-item-section>
        <h3>{{ lang.platformConfig }}</h3>
      </q-item-section>
    </q-item>
    <div class="row">

      <div class="col-3 left-wrapper">
        <div>
          <q-btn-group spread stretch>
            <q-btn unelevated size="md" :label="lang.addPlatform" @click="addPlatform()"/>
            <q-btn unelevated size="md" :label="lang.delPlatform" @click="deletePlatform()"/>
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
        <q-tabs v-model="tab" class="text-grey" active-color="primary" indicator-color="primary"
                align="justify">
          <q-tab name="base" :label="lang.platformConfig" tabindex="1"/>
          <q-tab name="simulator" :label="lang.simConfig" tabindex="2"/>
          <q-tab name="desc" :label="lang.platformDesc" tabindex="3"/>
          <q-tab name="rombase" :label="lang.rombaseAlias" tabindex="4"/>
        </q-tabs>
        <q-separator/>
        <q-tab-panels v-model="tab" keep-alive animated>
          <!-- 基本信息 -->
          <q-tab-panel name="base">
            <h6>{{ lang.platformBaseInfo }}</h6>
            <q-list>
              <q-item>
                <q-item-section>
                  <q-input filled square dense :label="lang.platformName" v-model="activePlatformInfo.Name"/>
                </q-item-section>
                <q-item-section>
                  <q-input filled square dense :label="lang.label" v-model="activePlatformInfo.Tag"/>
                </q-item-section>
                <q-item-section>
                  <q-input filled square dense :label="lang.platformIcon" class="open-dialog"
                           v-model="activePlatformInfo.Icon">
                    <template v-slot:append>
                      <q-icon name="upload_file" @click="openFileDialog()"/>
                    </template>
                  </q-input>
                </q-item-section>
                <q-item-section>
                  <q-select filled square dense map-options emit-value :options="boolOptions" :label="lang.hideName"
                            v-model="activePlatformInfo.HideName">
                    <q-tooltip>{{ lang.showHidePlatform }}</q-tooltip>
                  </q-select>
                </q-item-section>
              </q-item>
            </q-list>
            <h6>{{ lang.menuConfig }}</h6>
            <q-list>
              <q-item>
                <q-item-section>
                  <q-input filled square dense :label="lang.platformResDir" class="open-dialog"
                           v-model="activePlatformInfo.RootPath">
                    <template v-slot:append>
                      <q-icon name="create_new_folder" @click="openDirectoryDialog('root')"/>
                    </template>
                  </q-input>
                </q-item-section>
                <q-item-section>
                  <q-input filled square dense :label="lang.romDir" class="open-dialog"
                           v-model="activePlatformInfo.RomPath">
                    <template v-slot:append>
                      <q-icon name="create_new_folder" @click="openDirectoryDialog('rom')"/>
                    </template>
                    <q-tooltip>{{ lang.tipRomDir }}</q-tooltip>
                  </q-input>
                </q-item-section>
              </q-item>
            </q-list>

            <h6>{{ lang.romExt }}</h6>
            <q-list>
              <q-item>
                <q-item-section>
                  <q-input filled square dense :label="lang.romExt" style="width: 100%"
                           v-model="addExtInput" @keydown.enter="addExt">
                    <template v-slot:append>
                      <q-btn round dense flat icon="add" @click="addExt"/>
                    </template>
                    <q-tooltip>{{ lang.tipRomExt }}</q-tooltip>
                  </q-input>
                </q-item-section>
                <q-item-section side>
                  <q-btn square flat icon="assignment" @click="addNoExt()">
                    <q-tooltip>添加无扩展名支持</q-tooltip>
                  </q-btn>
                </q-item-section>
                <q-item-section side>
                  <q-btn square flat icon="delete" @click="clearRomExts()">
                    <q-tooltip>删除所有类型</q-tooltip>
                  </q-btn>
                </q-item-section>
              </q-item>
              <q-item>
                <div class="tag-label-wrapper q-gutter-xs">
                  <q-chip square removable size="sm" v-for="(ext,index) in activePlatformInfo.RomExts"
                          @remove="removeExt(ext,index)">
                    {{ ext }}
                  </q-chip>
                </div>
              </q-item>
            </q-list>
            <div class="q-gutter-sm bottom-bar">
              <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" @click="updatePlatform()"/>
            </div>
          </q-tab-panel>

          <!-- 模拟器管理 -->
          <q-tab-panel name="simulator">
            <simulators :platform="activePlatformId"/>
          </q-tab-panel>

          <!-- 平台简介 -->
          <q-tab-panel name="desc">
            <q-editor class="editor" square v-model="activePlatformInfo.Desc"
                      :definitions="editorDefinitions" :toolbar="EDITOR_TOOLBAR"/>
            <div class="q-gutter-sm bottom-bar">
              <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" @click="updateDesc()"/>
            </div>
          </q-tab-panel>
          <q-tab-panel name="rombase">
            <rombase-alias :platform="activePlatformId"/>
          </q-tab-panel>
        </q-tab-panels>
      </div>
      <div v-else class="col right-wrapper relative-position">
        <div class="absolute-full flex flex-center">{{ lang.tipSelectPlatform }}</div>
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
