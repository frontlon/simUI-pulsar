<template>

  <div class="wrapper">
    <q-banner class="bg-primary text-white">
      {{ lang.tipRomfile }}
    </q-banner>
    <q-list>
      <q-item>
        <q-item-section>
          <q-select dense filled square emit-value map-options :label="lang.selectPlatform" class="mt"
                    v-if="!isEmpty(platformList)" :options="platformList" v-model="activePlatform"
                    popup-content-class="q-select-content-min"
                    @update:model-value="changePlatform()"/>
        </q-item-section>
        <q-item-section>
          <q-select dense filled square emit-value map-options v-if="activePlatform" :options="menuOptions"
                    v-model="activeMenu" class="mt" popup-content-class="q-select-content-min"
                    @update:model-value="changeMenu()"
                    :label="lang.selectMenu"/>
        </q-item-section>
      </q-item>
    </q-list>

    <q-table
        v-if="romList && !isEmpty(listColumns)"
        flat bordered dense separator="cell" row-key="Id" selection="multiple" :pagination="initialPagination"
        :rows-per-page-options="[20,50,100,200,500,900,0]"
        :selected-rows-label="Id" v-model:selected="selected" :rows="romList" :columns="listColumns"
    >
      <template v-slot:top="props">
        <div class="col-2 q-table__title">{{ lang.romList }}</div>
        <q-space/>
        <q-btn-group outline style="margin-right: 10px">
          <q-btn size="md" color="primary" :label="lang.copyModule" :disable="isEmpty(selected)"
                 @click="clickOptButton('linkCopy')"></q-btn>
          <q-btn size="md" color="primary" :label="lang.moveModule" :disable="isEmpty(selected)"
                 @click="clickOptButton('linkMove')"></q-btn>
          <q-btn size="md" color="primary" :label="lang.delModule" :disable="isEmpty(selected)"
                 @click="clickOptButton('linkDel')"></q-btn>
        </q-btn-group>
        <q-btn-group outline>
          <q-btn size="md" color="negative" :label="lang.delRom" :disable="isEmpty(selected)"
                 @click="clickOptButton('romDel')"></q-btn>
        </q-btn-group>
      </template>


      <template v-slot:body-cell-RomName="props">

        <q-td key="RomName" :props="props">
          {{ props.row.RomName }}
          <q-popup-edit v-model="props.row.RomName" v-slot="scope">
            <q-input v-model="props.row.RomName" dense autofocus :label="lang.filename" @keyup.enter="scope.set"
                     :placeholder="props.row.RomName" :error-message="lang.tipFileNameIsNotEmpty"
                     :error="props.row.RomName== ''"
                     @change="renameRomfile(props.row)"/>
          </q-popup-edit>
        </q-td>
      </template>
      <template v-slot:body-cell-Name="props">
        <q-td key="Name" :props="props">
          {{ props.row.Name }}
          <q-popup-edit v-model="props.row.Name" v-slot="scope">
            <q-input v-model="props.row.Name" dense autofocus :label="lang.alias" @keyup.enter="scope.set"
                     :placeholder="props.row.Name" :error-message="lang.tipAliasIsNotEmpty" :error="props.row.Name== ''"
                     @change="renameLinkfile(props.row)"/>
          </q-popup-edit>
        </q-td>
      </template>
    </q-table>
  </div>

  <!-- 目录选择框 -->
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 250px">
      <q-card-section>
        <div class="text-h6">{{ currentOpt == 'linkCopy' ? lang.copyModule : lang.moveModule }}</div>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <q-list padding class="move-list">
          <q-item clickable dense active-class="menu-active" @click="activeMenuDialog = '/'"
                  :active="activeMenuDialog == '/'">
            <q-item-section side>
              <q-avatar icon="folder" color="primary" size="xs"/>
            </q-item-section>
            <q-item-section>
              <q-item-label lines="1">{{ lang.notCate }}</q-item-label>
            </q-item-section>
          </q-item>
          <div v-for="item in menuList">
            <q-item clickable dense active-class="menu-active" @click="activeMenuDialog = item.Path"
                    :active="activeMenuDialog == item.Path">
              <q-item-section side>
                <q-avatar icon="folder" color="primary" size="xs"/>
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ item.Name }}</q-item-label>
              </q-item-section>
            </q-item>
            <div v-if="item.SubMenu" v-for="sub in item.SubMenu">
              <q-item class="sub-tree" dense clickable active-class="menu-active" @click="activeMenuDialog = sub.Path"
                      :active="activeMenuDialog == sub.Path">
                <q-item-section side>
                  <q-avatar icon="folder" color="primary" size="xs"/>
                </q-item-section>
                <q-item-section>
                  <q-item-label lines="1">{{ sub.Name }}</q-item-label>
                </q-item-section>
              </q-item>
            </div>
          </div>
        </q-list>
      </q-card-section>
      <q-card-actions align="right" class="text-primary">
        <q-btn flat :label="lang.value.cancel" v-close-popup/>
        <q-btn flat :label="lang.value.ok" :disable="activeMenuDialog == null" v-close-popup @click="clickDialogOpt()"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue';

import {decodeApiData, isEmpty, loading, notify} from "components/utils";
import {
  BatchCopyRomLink,
  BatchDeleteRom,
  BatchDeleteRomLink,
  BatchMoveRomLink,
  GetGameList,
  GetMenuList,
  GetPlatform,
  RenameRomFile,
  RenameRomLink
} from "app/wailsjs/go/controller/Controller";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";
import {getPromptOpts} from "components/dialog";
import {Dialog} from "quasar";

const global = useGlobalStore();
const {config, lang} = storeToRefs(global);
const menuOptions: any = ref([{value: "", label: lang.value.all}, {value: "/", label: lang.value.notCate}])
const menuList: any = ref(null)
const showDialog: any = ref(false)
const platformList: any = ref([])
const romList: any = ref(null)
const activePlatform: any = ref(null)
const activeMenu: any = ref("")
const activeMenuDialog: any = ref("")
const selected = ref([])
const currentOpt = ref("")
const listColumns: any = ref([
  {name: "Id", label: lang.value.filename, field: "Id", sortable: true},
  {name: "RomName", label: lang.value.filename, field: "Name", sortable: true},
  {name: "Name", label: lang.value.alias, field: "Name", sortable: true},
  {name: "Menu", label: lang.value.belongMenu, field: "Menu", sortable: true},
])

let initialPagination = {descending: false, page: 0, rowsPerPage: 20}

onMounted(() => {
  //读取全部平台数据
  GetPlatform().then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      let jsonArr = resp.data
      jsonArr.forEach(item => {
        platformList.value.push({value: item.Id, label: item.Name});
      });
    }
  })

})

//点击切换平台
function changePlatform() {
  console.log("changePlatform")

  //创建菜单
  GetMenuList(activePlatform.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      menuList.value = resp.data
      let jsonArr = resp.data
      jsonArr.forEach(item => {
        menuOptions.value.push({value: item.Path, label: item.Name});
      });
    }
  })
  //读取rom列表
  getRomList()
}

//点击切换菜单
function changeMenu() {
  getRomList()
}

//读取rom列表
function getRomList() {
  selected.value = [];
  let isLike = activeMenu.value == "/" ? 0 : 1
  var req = {
    "platform": activePlatform.value,
    "catname": activeMenu.value,
    "showSubGame": 0,
    "page": -1, //加载全部数据
    "simpleModel": "simple",
    "catnameLike": isLike,
  };
  var request = JSON.stringify(req);
  GetGameList(request).then((result: string) => {
    let resp = decodeApiData(result)
    romList.value = resp.data;
  })
}

//点击操作按钮
function clickOptButton(type: string) {

  currentOpt.value = type
  console.log(type, selected.value)

  let roomIds: number[] = []
  selected.value.forEach((item, index) => {
    roomIds.push(item.Id);
  });

  switch (type) {
    case "linkCopy": //复制模块
      showDialog.value = true
      break
    case "linkMove": //移动模块
      showDialog.value = true
      break
    case "linkDel": //删除模块
      let linkOpt = getPromptOpts(lang.value.tipDelModule2, "", lang.value.ok, false)
      Dialog.create(linkOpt).onOk(() => {
        BatchDeleteRomLink(roomIds).then((result: string) => {
          let resp = decodeApiData(result)
          if (resp.err != "") {
            notify("err", resp.err)
            return
          }
          getRomList()
          notify("suc", lang.value.operSuc)
        })
      })
      break
    case "romDel": //删除ROM
      let opt = getPromptOpts(lang.value.delRom, lang.value.tipDelRom, lang.value.ok, false)
      opt.componentProps.toggle = ref({
        model: [],
        items: [
          {label: lang.value.andDelSubRom, value: '1', color: 'parmary'},
          {label: lang.value.andDelRomRes, value: '2', color: 'parmary'},
        ]
      })
      Dialog.create(opt).onOk((resp) => {
        let data = resp.toggle
        loading("show", lang.value.operating)
        let subGame = data.includes('1') ? 1 : 0;
        let res = data.includes('2') ? 1 : 0;
        BatchDeleteRom(roomIds, subGame, res).then((result: string) => {
          let resp = decodeApiData(result)
          console.log(resp.data)
          getRomList()
          notify("suc", lang.value.operSuc)
          loading("hide")
        })
      })
      break
  }
}

//移动或复制
function clickDialogOpt() {

  let roomIds: number[] = []
  selected.value.forEach((item, index) => {
    roomIds.push(item.Id);
  });

  switch (currentOpt.value) {
    case "linkCopy": //复制模块
      BatchCopyRomLink(roomIds, activeMenuDialog.value).then((result: string) => {
        let resp = decodeApiData(result)
        console.log(resp.data)
        if (resp.err != "") {
          notify("err", resp.err)
          return
        }
        getRomList()
        notify("suc", lang.value.operSuc)
      })
      break
    case "linkMove": //移动模块
      BatchMoveRomLink(roomIds, activeMenuDialog.value).then((result: string) => {
        let resp = decodeApiData(result)
        console.log(resp.data)
        if (resp.err != "") {
          notify("err", resp.err)
          return
        }
        getRomList()
        notify("suc", lang.value.operSuc)
      })
      break
  }
}


//修改rom文件名
function renameRomfile(item: any) {

  RenameRomFile(item.Id, item.RomName).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    getRomList()
    notify("suc", lang.value.operSuc)
  })

}

//修改模块名称
function renameLinkfile(item: any) {
  RenameRomLink(item.Id, item.Name).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
    getRomList()
    notify("suc", lang.value.operSuc)
  })
}

</script>
<style scoped>

.wrapper {
  width: 60%;
  margin: 0 auto;
}

.move-list .q-item {
  padding: 0 10px;
}

.sub-tree {
  padding-left: 30px !important;
}

.menu-active {
  background: var(--q-primary);
}
</style>
