<template>

  <q-banner>
    {{lang.tipRombaseAlias}}
  </q-banner>
  <q-list>
    <q-item>
      <q-item-section>
        <q-input filled square dense clearable style="width: 100%" :label="lang.baseOtherA" v-if="aliasMap"
                 v-model="aliasMap.OtherA"/>
      </q-item-section>
      <q-item-section>
        <q-input filled square dense clearable style="width: 100%" :label="lang.baseOtherB" v-if="aliasMap"
                 v-model="aliasMap.OtherB"/>
      </q-item-section>
      <q-item-section>
        <q-input filled square dense clearable style="width: 100%" :label="lang.baseOtherC" v-if="aliasMap"
                 v-model="aliasMap.OtherC"/>
      </q-item-section>
      <q-item-section>
        <q-input filled square dense clearable style="width: 100%" :label="lang.baseOtherD" v-if="aliasMap"
                 v-model="aliasMap.OtherD"/>
      </q-item-section>
    </q-item>
  </q-list>
  <div class="q-gutter-sm bottom-bar">
    <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" @click="updateAlias()"/>
  </div>

</template>

<script setup lang="ts">

import {defineProps, ref, watchEffect} from 'vue'
import {decodeApiData, notify} from 'components/utils'
import {GetRomBaseAlias, UpdateRomBaseAlias} from 'app/wailsjs/go/controller/Controller'
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";
const global = useGlobalStore();
const {lang} = storeToRefs(global);
const props = defineProps(['platform'])
const tab = ref('base')
const aliasMap = ref({})
const activePlatformId = ref(0)

watchEffect(() => {
  activePlatformId.value = props.platform
  GetRomBaseAlias(activePlatformId.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      aliasMap.value = resp.data
      if (aliasMap.value != null) {
        if (aliasMap.value.OtherA == lang.value.baseOtherA) {
          aliasMap.value.OtherA = ""
        }
        if (aliasMap.value.OtherB == lang.value.baseOtherB) {
          aliasMap.value.OtherB = ""
        }
        if (aliasMap.value.OtherC == lang.value.baseOtherC) {
          aliasMap.value.OtherC = ""
        }
        if (aliasMap.value.OtherD == lang.value.baseOtherD) {
          aliasMap.value.OtherD = ""
        }
      }
    }
  })
});

//更新资料项别名
function updateAlias() {
  let data = JSON.stringify(aliasMap.value)
  UpdateRomBaseAlias(activePlatformId.value, data).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
    } else {
      notify("err",  resp.err)
    }
  })
}
</script>
<style scoped>
@import "src/css/manage.css";
@import "src/css/page/platform.css";
</style>