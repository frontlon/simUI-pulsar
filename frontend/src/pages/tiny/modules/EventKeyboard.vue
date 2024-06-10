<template>

</template>
<script lang="ts">
import {openRunGameDialogByKeyboard, romEventKeyboard} from 'src/pages/tiny/RomListBar.vue'


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
  //event.preventDefault();
  if (event.key === 'Enter') {
    openRunGameDialogByKeyboard()
  } else if (event.key == "ArrowUp") {
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowDown") {
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowLeft") {
    romEventKeyboard(event.key)
  } else if (event.key == "ArrowRight") {
    romEventKeyboard(event.key)
  }
}

</script>

<style scoped>
</style>
