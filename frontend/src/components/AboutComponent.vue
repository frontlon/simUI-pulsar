<template>
  <q-dialog filled square ref="dialogRef">
    <q-card flat filled square class="q-dialog-plugin wrapper">
      <q-card-section>
        <img src="/images/ico.png"/>
      </q-card-section>
      <q-card-section class="selectable">
        <h5>SIMUI Pulsar v{{versionNo}}</h5>
        <p>编译时间：{{buildTime}}</p>
        <p>网站：<a target="_blank" href="http://www.simui.net/">www.simui.net</a></p>
        <p>QQ群：<a target="_blank"
                   href="https://qm.qq.com/cgi-bin/qm/qr?k=1RxGE9izpXHFnS48ocXzjwpqu6yQUSkb&jump_from=webapi">851347939</a>
        </p>
        <p>邮箱：front_diablo@163.com</p>
        <h6>特别鸣谢</h6>
        <p>春华秋实｜孟虎｜睡神林克｜wood.exe｜alexxyz2｜<br>冰澤｜Vanisper｜以及广大网友 ...</p>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import {useDialogPluginComponent} from 'quasar'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
const {dialogRef} = useDialogPluginComponent()
const global = useGlobalStore();
const {buildTime,versionNo} = storeToRefs(global);

</script>

<style scoped>
.wrapper {
  width: 400px;
  text-align: center;
  border: 1px solid var(--color-3);
}

h5 {
  margin: 0 0 10px 0;
  font-size: 18px;
}

h6 {
  margin:0;
  font-size: 16px;
  color: var(--color-10);
  font-weight: bold;
}

p {
  color: var(--color-7);
  line-height: 2em;
  margin: 0;
}

a,a:visited {
  color: var(--q-primary);
}

</style>