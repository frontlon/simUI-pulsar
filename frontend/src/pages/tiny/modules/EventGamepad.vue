<template>

</template>
<script lang="ts">
import {cancelEventKeyboard, clickEventKeyboard, romEventKeyboard} from 'src/pages/tiny/RomListBar.vue'
import {openPlatformDialog} from 'src/pages/tiny/Platform.vue'
import {notify} from "components/utils";
import {createGamepadController} from "components/gamepadInput";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);

const controller = createGamepadController({
  handlers: {
    onNavigate: romEventKeyboard,
    onConfirm: clickEventKeyboard,
    onCancel: cancelEventKeyboard,
    onOpenPlatform: openPlatformDialog,
  },
  onDisconnect: () => notify("info", lang.value.gamePadDisConnect),
})

export function initEventGamePad() {
  controller.init()
}

</script>

<style scoped>
</style>
