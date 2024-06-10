<template>

</template>
<script lang="ts">
import {ref} from "vue";
import {cancelEventKeyboard, clickEventKeyboard, romEventKeyboard} from 'src/pages/tiny/RomListBar.vue'
import {openPlatformDialog} from 'src/pages/tiny/Platform.vue'
import {getMilliSeconds, notify} from "components/utils";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);
//初始化手柄事件
const gamePadTimeOut = ref(100) //有摇杆100，无摇杆3000
let intervalId: NodeJS.Timeout;
let lastActiveTime = 0 //手柄最后活动时间
//let activeTime = 5000 //这个时间不操作手柄进入待机模式
let winFocus = true //当前窗口是否有焦点
let timeoutConnect = 70 //插入手柄轮询时间
let timeoutDisConnect = 2000 //未插入手柄轮询时间
window.addEventListener('focus', function () {
  winFocus = true
  console.log('setInterval 窗口获得焦点');
});

window.addEventListener('blur', function () {
  winFocus = false
  console.log('setInterval onblur');
});


//启动手柄监控
function startInterval(e: GamepadEvent) {
  console.log("startInterval")
  intervalId = setInterval(function () {
    gamepadHandle(e)
  }, gamePadTimeOut.value);
}

//重启手柄定时器
function restartInterval(e: GamepadEvent, timeout: number, tip: boolean) {
  console.log("setInterval restartInterval", timeout)
  gamePadTimeOut.value = timeout
  clearInterval(intervalId)
  startInterval(e)
  if (tip) {
    if (gamePadTimeOut.value == timeoutDisConnect) {
      notify("info", lang.value.gamePadDisConnect)
    }/* else {
      notify("suc", "连接手柄成功")
    }*/
  }

}

function checkTime() {
  let curr = getMilliSeconds()
  if (curr - lastActiveTime < 100) {
    return false
  }
  lastActiveTime = curr
  return true
}

// 监听游戏手柄
export function initEventGamePad() {
  window.addEventListener("gamepadconnected", (e) => {
    startInterval(e)
  });
}

function gamepadHandle(e: GamepadEvent) {

  var gamepad = navigator.getGamepads()[e.gamepad.index];

  if (winFocus == false) {
    //失去焦点，降低轮询频次
    restartInterval(e, timeoutDisConnect, false)
    return
  }

  let timeout = gamepad == null ? timeoutDisConnect : timeoutConnect
  if (gamePadTimeOut.value != timeout) {
    //重启、关闭手柄
    restartInterval(e, timeout, true)
    return;
  }
  if (gamepad != null) {
    gamepadControlButton(gamepad.buttons)
    gamepadControlAxes(gamepad.axes)
  }

  /*let curr = getMilliSeconds()
  if (lastActiveTime > 0 && curr - lastActiveTime > activeTime) {
    //手柄长时间未活动，将此轮询频次
    console.log("asdfsadf",lastActiveTime)
    restartInterval(e, timeoutDisConnect, false)
    return
  }*/

}

//手柄方向
function gamepadControlAxes(axes: ReadonlyArray<number>) {
  let btn = ""
  if (axes[0] == -1) { //左1
    btn = "ArrowLeft"
  } else if (axes[0] == 1) { //右1
    btn = "ArrowRight"
  } else if (axes[1] == -1) { //上1
    btn = "ArrowUp"
  } else if (axes[1] == 1) { //下1
    btn = "ArrowDown"
  } else if (axes[2] == -1) { //左2
    btn = "ArrowLeft"
  } else if (axes[2] == 1) { //右2
    btn = "ArrowRight"
  } else if (axes[3] == -1) { //上2
    btn = "ArrowUp"
  } else if (axes[3] == 1) { //下2
    btn = "ArrowDown"
  }

  if (btn != "") {
    if (!checkTime()) {
      return
    }
    romEventKeyboard(btn)
  }

}

//手柄按钮
function gamepadControlButton(buttons: ReadonlyArray<GamepadButton>) {
  for (let i = 0; i < buttons.length; i++) {
    if (buttons[i].value == 1) {
      let btn = i.toString()
      if (btn != "") {
        if (!checkTime()) {
          return
        }
      }

      if (i == 0) { //摇杆1，A，点击激活
        console.log("摇杆1，A，点击激活");
        clickEventKeyboard()
      } else if (i == 1) { //取消
        console.log("摇杆2，B，取消");
        cancelEventKeyboard()
      } else if (i == 2) { //打开平台
        openPlatformDialog()
        console.log("摇杆3，X，打开平台");
      }else if (i == 3) { //摇杆4，Y，启动游戏
        //openPlatformDialog()
        //console.log("摇杆4，Y，启动游戏");
      } else if (i == 9) { //手柄START
        console.log("手柄START");
        clickEventKeyboard()
      }else if (i == 12) { //手柄START
        console.log("摇杆上");
        romEventKeyboard("ArrowUp")
      }else if (i == 13) { //手柄START
        console.log("摇杆下");
        romEventKeyboard("ArrowDown")
      }else if (i == 14) { //手柄START
        console.log("摇杆左");
        romEventKeyboard("ArrowLeft")
      }else if (i == 15) { //手柄START
        console.log("摇杆右");
        romEventKeyboard("ArrowRight")
      }
    }
  }
}


</script>

<style scoped>
</style>
