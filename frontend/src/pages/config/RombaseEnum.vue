<template>
  <q-card flat bordered class="config-section">
    <q-card-section>
      <div class="section-title">{{ lang.enum }}</div>
      <div class="row q-col-gutter-md enum-grid">
        <div class="col-12 col-sm-6 col-lg-3 enum-field">
          <q-input v-model="typeList" square filled type="textarea" :label="lang.baseType"/>
        </div>
        <div class="col-12 col-sm-6 col-lg-3 enum-field">
          <q-input v-model="yearList" square filled type="textarea" :label="lang.baseYear"/>
        </div>
        <div class="col-12 col-sm-6 col-lg-3 enum-field">
          <q-input v-model="countryList" square filled type="textarea" :label="lang.baseCountry"/>
        </div>
        <div class="col-12 col-sm-6 col-lg-3 enum-field">
          <q-input v-model="versionList" square filled type="textarea" :label="lang.baseVersion"/>
        </div>
        <div class="col-12 col-sm-6 col-lg-4 enum-field">
          <q-input v-model="producerList" square filled type="textarea" :label="lang.baseProducer"/>
        </div>
        <div class="col-12 col-sm-6 col-lg-4 enum-field">
          <q-input v-model="publisherList" square filled type="textarea" :label="lang.basePublisher"/>
        </div>
        <div class="col-12 col-sm-6 col-lg-4 enum-field">
          <q-input v-model="translateList" square filled type="textarea" :label="lang.baseTranslate"/>
        </div>
      </div>
    </q-card-section>
  </q-card>

  <div class="bottom-bar">
    <q-btn unelevated :label="lang.update" color="primary" class="update-btn" @click="updateEnum()"/>
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

function updateEnum() {
  let type = typeList.value.split("\n");
  let year = yearList.value.split("\n");
  let producer = producerList.value.split("\n");
  let publisher = publisherList.value.split("\n");
  let country = countryList.value.split("\n");
  let translate = translateList.value.split("\n");
  let version = versionList.value.split("\n");

  UpdateRomBaseEnumByType("type", JSON.stringify(type)).then(() => {
  })
  UpdateRomBaseEnumByType("year", JSON.stringify(year)).then(() => {
  })
  UpdateRomBaseEnumByType("producer", JSON.stringify(producer)).then(() => {
  })
  UpdateRomBaseEnumByType("publisher", JSON.stringify(publisher)).then(() => {
  })
  UpdateRomBaseEnumByType("country", JSON.stringify(country)).then(() => {
  })
  UpdateRomBaseEnumByType("translate", JSON.stringify(translate)).then(() => {
  })
  UpdateRomBaseEnumByType("version", JSON.stringify(version)).then(() => {
  })
  notify("suc", lang.value.updateSuc)
}

</script>
<style scoped>
@import "src/css/page/config.css";
</style>
