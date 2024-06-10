<template>

</template>
<script lang="ts">
import {romEventKeyboard,clickEventKeyboard} from 'src/pages/playnite/RomListBar.vue'
import {openRunGameDialogByKeyboard} from 'src/pages/playnite/ContentBar.vue'

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

  if (event.key === 'Enter') {
    //event.preventDefault();
    openRunGameDialogByKeyboard()
  } else if (event.key == "ArrowUp") {
    //event.preventDefault();
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowDown") {
    // event.preventDefault();
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowLeft") {
    //event.preventDefault();
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowRight") {
    //event.preventDefault();
    romEventKeyboard(event.key)
  }
}

</script>

<style scoped>
</style>
