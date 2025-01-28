<template>
  <q-dialog filled square ref="dialogRef">
    <q-card flat filled square class="q-dialog-plugin wrapper">
      <q-card-section>
        <img style="width: 50%" src="/images/copyright.png"/>
      </q-card-section>
      <q-card-section class="selectable">
        <h5>copyright</h5>
        <p>info</p>
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