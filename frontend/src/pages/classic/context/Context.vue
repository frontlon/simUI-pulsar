<!-- 右键菜单 -->
<template>
  <base-info/> <!-- 编辑资料 -->
  <move-module/> <!-- 移动模块 -->
  <copy-module/> <!-- 复制模块 -->
  <sub-game/> <!-- 子游戏配置 -->
  <thumbs/> <!-- 图片资源编辑 -->
  <share-rom/> <!-- 分享ROM -->
</template>

<script lang="ts">
import {h, ref} from 'vue';
import {Dialog} from 'quasar';
import {getPromptOpts} from 'components/dialog'
import {decodeApiData, isEmpty, loading, notify} from "components/utils";
import {CONTEXT_ICON_SIZE} from "boot/constant";
import BaseInfo, {openBaseInfoDialog} from 'pages/classic/context/BaseInfo.vue'
import MoveModule, {openMoveModuleDialog} from 'pages/classic/context/MoveModule.vue'
import CopyModule, {openCopyModuleDialog} from 'pages/classic/context/CopyModule.vue'
import SubGame, {openSubGameDialog} from 'pages/classic/context/SubGame.vue'
import Thumbs, {openThumbsDialog} from 'pages/classic/context/Thumbs.vue'
import ShareRom, {openShareRomDialog} from 'pages/classic/context/ShareRom.vue'

import {
  DeleteRom,
  DeleteRomLink,
  GetSubGameList,
  OpenFolder,
  RenameRomFile,
  RenameRomLink,
  RunGame,
  SetFavorite,
  SetHide,
} from "app/wailsjs/go/controller/Controller";
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from "pinia";
import '@imengyu/vue3-context-menu/lib/vue3-context-menu.css'
import ContextMenu from '@imengyu/vue3-context-menu';

const global = useGlobalStore();
const {config, lang, callbackOpts, simulatorMap} = storeToRefs(global);
const targetIndex: any = ref(0)
const targetId: any = ref(0)
const romInfo: any = ref(null)
const modelRenameFile = ref("")

export function showContextMenu(event: any, index: number, id: number, info: any) {

  event.preventDefault();

  targetId.value = id;
  targetIndex.value = index;
  romInfo.value = info;
  console.log('Right-clicked', id);

  //加载我的最爱
  let favorite;
  if (romInfo.value.Favorite == 1) { //取消最爱
    favorite = {
      label: lang.value.cancelFavorite + "(C+f)",
      icon: h('img', {src: '/images/context/fav_0.png', style: CONTEXT_ICON_SIZE}),
      onClick: () => setFavorite(0)
    }
  } else {
    favorite = {
      label: lang.value.setFavorite + "(C+f)",
      icon: h('img', {src: '/images/context/fav_1.png', style: CONTEXT_ICON_SIZE}),
      onClick: () => setFavorite(1)
    }
  }

  let contextData = getContextMenuData(event)

  //加载模拟器数据
  let simulator = getContextMenuDataSimulator()
  contextData.items = [favorite, ...contextData.items, simulator]

  //加载游戏列表
  let gameList: any = [];
  let master = {
    label: lang.value.run + romInfo.value.Name,
    icon: h('img', {src: '/images/context/app.png', style: CONTEXT_ICON_SIZE}),
    onClick: () => runGame(targetId.value)

  }
  gameList.push(master)

  //读取子游戏
  GetSubGameList(id).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data.length > 0) {
      resp.data.forEach(function (item: any, index: number) {
        let s = {
          label: lang.value.run + item.Name,
          icon: h('img', {src: '/images/context/app.png', style: CONTEXT_ICON_SIZE}),
          onClick: () => runGame(item.Id)
        }
        gameList.push(s);
      })
    }

    gameList[gameList.length - 1]['divided'] = true;
    contextData.items = [...gameList, ...contextData.items]
    ContextMenu.showContextMenu(contextData)

  })

}

//右键菜单
function getContextMenuData(event: any) {

  return {
    x: event.x,
    y: event.y,
    customClass: "default dark",
    items: [
      {
        label: lang.value.editRombase + "(C+e)",
        icon: h('img', {src: '/images/context/baseinfo.png', style: CONTEXT_ICON_SIZE}),
        onClick: () => openBaseInfoDialog(targetId.value, targetIndex.value)
      },
      {
        label: lang.value.editThumbs + "(C+t)",
        icon: h('img', {src: '/images/context/thumbs.png', style: CONTEXT_ICON_SIZE}),
        onClick: () => openThumbsDialog(targetIndex.value, romInfo.value)
      },
      {
        label: lang.value.editSubGame + "(C+s)",
        divided: true,
        icon: h('img', {src: '/images/context/sub.png', style: CONTEXT_ICON_SIZE}),
        onClick: () => openSubGameDialog(targetIndex.value, romInfo.value)
      },
      {
        label: lang.value.fileOperation,
        divided: true,
        icon: h('img', {src: '/images/context/files.png', style: CONTEXT_ICON_SIZE}),
        children: [
          {
            label: lang.value.editAlias + "(C+r | F2)",
            divided: true,
            icon: h('img', {src: '/images/context/alias.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => renameAlias()
          },
          {
            label: (romInfo.value.Hide == 1 ? lang.value.setShow : lang.value.setHide) + "(C+h)",
            divided: true,
            icon: h('img', {src: '/images/context/hide_1.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => setHide()
          },
          {
            label: lang.value.copyModule + "(C+c)",
            icon: h('img', {src: '/images/context/copy.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openCopyModuleDialog(targetIndex.value, romInfo.value)
          },
          {
            label: lang.value.moveModule + "(C+v)",
            icon: h('img', {src: '/images/context/move.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openMoveModuleDialog(targetIndex.value, romInfo.value)
          },
          {
            label: lang.value.delModule + "(C+d)",
            divided: true,
            icon: h('img', {src: '/images/context/unlink.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => deleteRomLink()
          },
          {
            label: lang.value.renameRom + "(C+S+r | C+S+F2)",
            icon: h('img', {src: '/images/context/rename.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => renameFile()
          },
          {
            label: lang.value.delRom + "(C+S+d)",
            divided: true,
            icon: h('img', {src: '/images/context/delete.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => deleteRom()
          },
          {
            label: lang.value.outputShare,
            icon: h('img', {src: '/images/context/output.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openShareRomDialog(targetId.value, romInfo.value)
          },
        ]
      },
      {
        label: lang.value.locationRomDir,
        icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
        onClick: () => openFolder("rom", targetId.value)
      },
      {
        label: lang.value.locationRomLinkDir,
        icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
        onClick: () => openFolder("slnk", targetId.value)
      },
      {
        label: lang.value.locationPlatformDir,
        icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
        onClick: () => openFolder("platform", targetId.value)
      },
      {
        label: lang.value.locationResDir,
        icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
        children: [
          {
            label: lang.value.thumb,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("thumb", targetId.value)
          },
          {
            label: lang.value.snap,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("snap", targetId.value)
          },
          {
            label: lang.value.poster,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("poster", targetId.value)
          },
          {
            label: lang.value.packing,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("packing", targetId.value)
          },
          {
            label: lang.value.title,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("title", targetId.value)
          },
          {
            label: lang.value.cassette,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("cassette", targetId.value)
          },
          {
            label: lang.value.icon,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("icon", targetId.value)
          },
          {
            label: lang.value.gif,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("gif", targetId.value)
          },
          {
            label: lang.value.background,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("background", targetId.value)
          },
          {
            label: lang.value.video,
            icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
            onClick: () => openFolder("video", targetId.value)
          }
        ]
      }
    ]
  }
}

//模拟器菜单
function getContextMenuDataSimulator() {
  let simulator: any = {}
  if (!isEmpty(simulatorMap.value[romInfo.value.Platform])) {
    let childrens: any = []
    simulatorMap.value[romInfo.value.Platform].forEach(function (item: any, index: number) {
      let s = {
        label: item.Name,
        icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
        onClick: () => openFolder("sim", item.Id)
      }
      childrens.push(s);
    })

    simulator = {
      label: lang.value.locationSimDir,
      icon: h('img', {src: '/images/context/folder.png', style: CONTEXT_ICON_SIZE}),
      children: childrens,
    }
  }

  return simulator
}

//运行游戏
function runGame(id: number) {
  RunGame(id, 0).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }
  })
}

//修改别名
function renameAlias() {
  let opt = getPromptOpts(lang.value.editAlias, lang.value.tipRenameAlias, lang.value.ok, false, romInfo.value.Name)

  Dialog.create(opt).onOk(data => {
    let val = data.input

    if (val == "") {
      notify("err", lang.value.tipAliasIsNotEmpty)
      return
    }

    RenameRomLink(targetId.value, val).then((result: string) => {
      let resp = decodeApiData(result)
      if (resp.err == "") {
        callback("renameRomLink", targetIndex.value, val)
      }
    })
  })

}

//修改文件名
function renameFile() {
  let opt = getPromptOpts(lang.value.renameRom, lang.value.tipRenameRomFile, lang.value.ok, false, romInfo.value.RomName)

  Dialog.create(opt).onOk(resp => {
    if (resp.input == "") {
      notify("err", lang.value.tipFileNameIsNotEmpty)
      return
    }
    let data = resp.input
    console.log(targetId.value, data)
    RenameRomFile(targetId.value, data).then((result: string) => {
    })
  })
}

//设置隐藏
function setHide() {
  let hide = romInfo.value.Hide == 1 ? 0 : 1;
  console.log("hide", romInfo.value.Hide, hide)
  SetHide(targetId.value, hide).then((result: string) => {
    let resp = decodeApiData(result)
    romInfo.value.Hide = hide;
    console.log(resp.data)
    callback("romHide", targetIndex.value, hide)
  })
}


//设置喜爱
function setFavorite(fav) {
  console.log("setFavorite", romInfo.value.Favorite, fav)
  SetFavorite(targetId.value, fav).then((result: string) => {
    let resp = decodeApiData(result)
    romInfo.value.Favorite = fav;
    console.log(resp.data)
    notify("suc", lang.value.operSuc)
    callback("setFavorite", targetIndex.value, fav)
  })
}

//删除rom
function deleteRom() {
  let opt = getPromptOpts(lang.value.delRom, lang.value.tipDelRom, lang.value.ok, false, null, romInfo.value.Name)
  opt.componentProps.toggle = ref({
    model: [],
    items: [
      {label: lang.value.andDelSubRom, value: '1', color: 'primary'},
      {label: lang.value.andDelRomRes, value: '2', color: 'green'},
    ]
  })

  Dialog.create(opt).onOk((resp) => {

    let data = resp.toggle

    loading("show", lang.value.deleteing)
    let subGame = data.includes('1') ? 1 : 0;
    let res = data.includes('2') ? 1 : 0;

    DeleteRom(targetId.value, subGame, res).then((result: string) => {
      let resp = decodeApiData(result)
      console.log(resp.data)
      callback("deleteRom", targetIndex.value)
      loading("hide")
    })
  })
}

//删除链接
function deleteRomLink() {
  let opt = getPromptOpts(lang.value.delModule, lang.value.tipDelModule, lang.value.ok, false, null, romInfo.value.Name)
  Dialog.create(opt).onOk(() => {
    DeleteRomLink(targetId.value).then((result: string) => {
      let resp = decodeApiData(result)
      console.log(resp.data)
      callback("deleteRomLink", targetIndex.value)
    })
  })
}


//返回值给父组件
export function callback(opt: string, index: number, val: any = "") {
  console.log("callback", opt, index, val)
  callbackOpts.value = {
    //"id": Date.now(),
    "index": index,
    "opt": opt,
    "data": val,
  }
}

//打开文件夹
function openFolder(type: string, id: number) {

  OpenFolder(type, id).then((result: string) => {
  })
}

//键盘打开窗口
export function keybordOpenWin(opt: string, index: number, info) {
  targetId.value = info.Id;
  targetIndex.value = index;
  romInfo.value = info;

  switch (opt) {
    case "deleteRom": //删除rom
      deleteRom()
      break;
    case "renameFile": //rom重命名
      renameFile()
      break;
    case "setHide": //设为隐藏
      setHide()
      break;
    case "setFavorite": //设为喜爱
      setFavorite(1)
      break;
    case "renameAlias": //修改别名
      renameAlias()
      break;
    case "openCopyModuleDialog": //复制模块
      openCopyModuleDialog(targetIndex.value, romInfo.value)
      break;
    case "openMoveModuleDialog": //移动模块
      openMoveModuleDialog(targetIndex.value, romInfo.value)
      break;
    case "deleteRomLink": //删除模块
      deleteRomLink()
      break;
    case "openBaseInfoDialog": //编辑资料
      openBaseInfoDialog(targetId.value, targetIndex.value)
      break;
    case "openThumbsDialog": //编辑图片
      openThumbsDialog(targetIndex.value, romInfo.value)
      break;
    case "openSubGameDialog": //编辑子游戏
      openSubGameDialog(targetIndex.value, romInfo.value)
      break;
  }
}

export default {
  components: {ShareRom, Thumbs, SubGame, CopyModule, MoveModule, BaseInfo},
  setup() {
    return {
      ContextMenu,
      targetId,
      modelRenameFile,
      romInfo,
      renameAlias,
      renameFile,
      setHide,
      deleteRom,
      deleteRomLink,
      openFolder,
      keybordOpenWin,
    };
  }
}
</script>


<style scoped>
@import "src/css/classic/contentContext.css";

.submenu-active {
  color: var(--color-15);
  background: var(--q-primary);
  background: #f00;
}

</style>

