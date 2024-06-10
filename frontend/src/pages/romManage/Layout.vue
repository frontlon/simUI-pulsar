<template>

  <div>
    <q-item style="--wails-draggable:drag">
      <q-item-section avatar>
          <q-btn flat size="lg" icon="arrow_back" class="vertical-middle" @click="$router.back()"/>
      </q-item-section>
      <q-item-section>
        <h3>{{lang.romManage}}</h3>
      </q-item-section>
    </q-item>
  </div>
  <q-tabs v-model="tab" class="text-grey" active-color="primary" indicator-color="primary" align="justify">
    <q-tab name="rom" :label="lang.romManage"/>
    <q-tab name="base" :label="lang.rombaseManage"/>
    <q-tab name="simulator" :label="lang.simulatorManage"/>
    <q-tab name="repeat" :label="lang.romRepeatCheck"/>
    <q-tab name="unowned" :label="lang.unownedClear"/>
  </q-tabs>

  <q-separator/>
  <q-tab-panels v-model="tab" animated>

    <!-- ROM管理 -->
    <q-tab-panel name="rom">
      <romfile/>
    </q-tab-panel>

    <!-- 编辑资料 -->
    <q-tab-panel name="base">
      <rombase/>
    </q-tab-panel>

    <!-- 模拟器管理 -->
    <q-tab-panel name="simulator">
      <simulator/>
    </q-tab-panel>

    <!-- ROM重复检测 -->
    <q-tab-panel name="repeat">
      <repeat/>
    </q-tab-panel>

    <!-- 无效资源清理 -->
    <q-tab-panel name="unowned">
      <unowned/>
    </q-tab-panel>

  </q-tab-panels>

</template>

<script setup lang="ts">
import {ref} from 'vue'
import Rombase from 'src/pages/romManage/Rombase.vue'
import Romfile from 'src/pages/romManage/Romfile.vue'
import Simulator from 'src/pages/romManage/Simulator.vue'
import Unowned from 'src/pages/romManage/Unowned.vue'
import Repeat from 'src/pages/romManage/Repeat.vue'
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";
const global = useGlobalStore();
const {lang} = storeToRefs(global);
const tab = ref('rom')

</script>
<style scoped>
@import "src/css/romManage.css";
</style>