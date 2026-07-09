<template>

</template>
<script lang="ts">
import {clickGameEvent, romEventKeyboard, runActiveGame} from 'src/pages/classic/ContentBar.vue'
import {notify} from "components/utils";
import {clickMenuEvent, menuEventKeyboard} from "pages/classic/LeftBarMenu.vue";
import {clickPlatformEvent, platformEventKeyboard} from "pages/classic/LeftBarPlatform.vue";
import {createGamepadController} from "components/gamepadInput";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);

const controller = createGamepadController({
  handlers: {
    onNavigate: romEventKeyboard,
    onConfirm: () => {
      clickGameEvent()
      clickMenuEvent()
      clickPlatformEvent()
    },
    onRunGame: runActiveGame,
    onPlatformNav: platformEventKeyboard,
    onMenuNav: menuEventKeyboard,
  },
  onDisconnect: () => notify("info", lang.value.gamePadDisConnect),
})

export function initEventGamePad() {
  controller.init()
}

</script>

<style scoped>
</style>
