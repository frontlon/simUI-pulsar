<template>

</template>
<script lang="ts">
import {createCache} from 'src/pages/classic/HeaderBarTool.vue'
import {manageEventKeyboard, romEventKeyboard, runActiveGame} from 'src/pages/classic/ContentBar.vue'
import {platformEventKeyboard} from 'src/pages/classic/LeftBarPlatform.vue'
import {menuEventKeyboard} from 'src/pages/classic/LeftBarMenu.vue'

//绑定键盘事件
export function addKeyboardEvent() {
  document.addEventListener('keydown', keyboardOpt)
}

//解绑键盘事件
export function removeKeyboardEvent() {
  document.removeEventListener('keydown', keyboardOpt, false);
}

//键盘事件操作
function keyboardOpt(event) {

  if (event.target.matches("input")) {
    return;
  } else if (event.target.matches("textarea")) {
    return;
  } else if (event.target.classList.contains("q-editor__content")) {
    return;
  }

  if (event.key === 'F5') {
    createCache('platform');
  } else if (event.ctrlKey && event.key === 'Enter') {
    runActiveGame()
  } else if (event.shiftKey && event.key == "ArrowUp") {
    event.preventDefault();
    platformEventKeyboard(event.key)
  } else if (event.shiftKey && event.key == "ArrowDown") {
    event.preventDefault();
    platformEventKeyboard(event.key)
  } else if (event.ctrlKey && event.key == "ArrowUp") {
    event.preventDefault();
    menuEventKeyboard(event.key)
  } else if (event.ctrlKey && event.key == "ArrowDown") {
    event.preventDefault();
    menuEventKeyboard(event.key)
  } else if (event.key == "ArrowUp") {
    event.preventDefault();
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowDown") {
    event.preventDefault();
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowLeft") {
    event.preventDefault();
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowRight") {
    event.preventDefault();
    romEventKeyboard(event.key)
  } else if (event.ctrlKey && event.shiftKey && event.key.toLowerCase() === 'd') { //删除rom
    event.preventDefault();
    manageEventKeyboard("deleteRom")
  } else if (event.ctrlKey && event.shiftKey && event.key.toLowerCase() === 'r') { //rom重命名
    event.preventDefault();
    manageEventKeyboard("renameFile")
  }  else if (event.ctrlKey && event.shiftKey && event.key == 'F2') { //rom重命名
    event.preventDefault();
    manageEventKeyboard("renameFile")
  } else if (event.ctrlKey && event.shiftKey && event.key.toLowerCase() === 'c') { //没用，防止触发ctrl+c
    event.preventDefault();
  } else if (event.ctrlKey && event.key.toLowerCase() === 'h') { //设为隐藏
    event.preventDefault();
    manageEventKeyboard("setHide")
  } else if (event.ctrlKey && event.key.toLowerCase() === 'r') { //修改别名
    event.preventDefault();
    manageEventKeyboard("renameAlias")
  } else if (event.key === 'F2') {
    event.preventDefault();
    manageEventKeyboard("renameAlias")
  } else if (event.ctrlKey && event.key.toLowerCase() === 'c') { //复制模块
    event.preventDefault();
    manageEventKeyboard("openCopyModuleDialog")
  } else if (event.ctrlKey && event.key.toLowerCase() === 'v') { //移动模块
    event.preventDefault();
    manageEventKeyboard("openMoveModuleDialog")
  } else if (event.ctrlKey && event.key.toLowerCase() === 'd') { //删除模块
    event.preventDefault();
    manageEventKeyboard("deleteRomLink")
  } else if (event.ctrlKey && event.key.toLowerCase() === 't') { //编辑展示图
    event.preventDefault();
    manageEventKeyboard("openThumbsDialog")
  } else if (event.ctrlKey && event.key.toLowerCase() === 's') { //编辑子游戏
    event.preventDefault();
    manageEventKeyboard("openSubGameDialog")
  } else if (event.ctrlKey && event.key.toLowerCase() === 'e') { //编辑资料
    event.preventDefault();
    manageEventKeyboard("openBaseInfoDialog")
  }

}

</script>

<style scoped>
</style>