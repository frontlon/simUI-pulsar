<template>
  <div class="row base-list-wrapper">
    <q-list>
      <q-item>
        <q-item-section>
          <q-input v-model="typeList" square filled type="textarea" :label="lang.baseType"/>
        </q-item-section>
        <q-item-section>
          <q-input v-model="yearList" square filled type="textarea" :label="lang.baseYear"/>
        </q-item-section>
        <q-item-section>
          <q-input v-model="countryList" square filled type="textarea" :label="lang.baseCountry"/>
        </q-item-section>
        <q-item-section>
          <q-input v-model="versionList" square filled type="textarea" :label="lang.baseVersion"/>
        </q-item-section>
      </q-item>
      <q-item>
        <q-item-section>
          <q-input v-model="producerList" square filled type="textarea" :label="lang.baseProducer"/>
        </q-item-section>
        <q-item-section>
          <q-input v-model="publisherList" square filled type="textarea" :label="lang.basePublisher"/>
        </q-item-section>
        <q-item-section>
          <q-input v-model="translateList" square filled type="textarea" :label="lang.baseTranslate"/>
        </q-item-section>
      </q-item>
    </q-list>
  </div>
  <div class="q-gutter-sm bottom-bar">
    <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" @click="updateEnum()"/>
  </div>

</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {useGlobalStore} from 'stores/globalData';
import {decodeApiData, notify} from "components/utils";
import {GetRomBaseEnum, UpdateRomBaseEnumByType,} from 'app/wailsjs/go/controller/Controller'
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);
const typeList = ref("")
const yearList = ref("")
const producerList = ref("")
const publisherList = ref("")
const countryList = ref("")
const translateList = ref("")
const versionList = ref("")

onMounted(() => {
  //读取全部枚举
  GetRomBaseEnum().then((result: string) => {
    let resp = decodeApiData(result)
    let data = resp.data
    typeList.value = data.type.join("\n");
    yearList.value = data.year.join("\n");
    producerList.value = data.producer.join("\n");
    publisherList.value = data.publisher.join("\n");
    countryList.value = data.country.join("\n");
    translateList.value = data.translate.join("\n");
    versionList.value = data.version.join("\n");
  })
})

//更新枚举数据
function updateEnum() {
  let type = typeList.value.split("\n");
  let year = yearList.value.split("\n");
  let producer = producerList.value.split("\n");
  let publisher = publisherList.value.split("\n");
  let country = countryList.value.split("\n");
  let translate = translateList.value.split("\n");
  let version = versionList.value.split("\n");

  UpdateRomBaseEnumByType("type", JSON.stringify(type)).then((result: string) => {
  })
  UpdateRomBaseEnumByType("year", JSON.stringify(year)).then((result: string) => {
  })
  UpdateRomBaseEnumByType("producer", JSON.stringify(producer)).then((result: string) => {
  })
  UpdateRomBaseEnumByType("publisher", JSON.stringify(publisher)).then((result: string) => {
  })
  UpdateRomBaseEnumByType("country", JSON.stringify(country)).then((result: string) => {
  })
  UpdateRomBaseEnumByType("translate", JSON.stringify(translate)).then((result: string) => {
  })
  UpdateRomBaseEnumByType("version", JSON.stringify(version)).then((result: string) => {
  })
  notify("suc", lang.value.updateSuc)
}

</script>
<style scoped>
@import "src/css/manage.css";

.base-list-wrapper .q-textarea {
  margin: 8px 4px;
}

.base-list-wrapper textarea {
  height: 300px;
}

</style>