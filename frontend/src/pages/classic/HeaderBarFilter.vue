<template>

  <div class="q-gutter-xs row filter-wrapper">

    <!-- 类型 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('BaseType')"
              :options="filters.base_type" v-model="activeType" dense="dense" popup-content-class="q-select-content"
              @update:model-value="changeFilter()">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- 发行日期 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('BaseYear')"
              popup-content-class="q-select-content" :options="filters.base_year" v-model="activeYear" dense="dense"
              @update:model-value="changeFilter()">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>
    <!-- 制作商 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('BaseProducer')"
              :options="filters.base_producer" v-model="activeProducer" dense="dense"
              @update:model-value="changeFilter()" popup-content-class="q-select-content">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- 发行商 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('BasePublisher')"
              :options="filters.base_publisher" v-model="activePublisher" dense="dense"
              @update:model-value="changeFilter()" popup-content-class="q-select-content">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- 地区 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('BaseCountry')"
              :options="filters.base_country" v-model="activeCountry" dense="dense"
              @update:model-value="changeFilter()" popup-content-class="q-select-content">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- 汉化者 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('BaseTranslate')"
              :options="filters.base_translate" v-model="activeTranslate" dense="dense"
              @update:model-value="changeFilter()" popup-content-class="q-select-content">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- 版本 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('BaseVersion')"
              :options="filters.base_version" v-model="activeVersion" dense="dense"
              @update:model-value="changeFilter()" popup-content-class="q-select-content">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- 评分 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('Score')"
              :options="filters.score" v-model="activeScore" dense="dense"
              @update:model-value="changeFilter()" popup-content-class="q-select-content">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- 通关状态 -->
    <q-select standout square v-if="!isEmpty(platformUi) && platformUi.RomListColumn.includes('Complete')"
              :options="filters.complete" v-model="activeComplete" dense="dense"
              @update:model-value="changeFilter()" popup-content-class="q-select-content">
      <template v-slot:option="scope">
        <q-item v-bind="scope.itemProps" dense="dense">
          <q-item-section>
            <q-item-label>{{ scope.opt.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
      <template v-slot:no-option>
        <q-item>
          <q-item-section class="select-no-option">{{ lang.notLabel }}</q-item-section>
        </q-item>
      </template>
    </q-select>

    <!-- 关键字 -->
    <q-input filled square :placeholder="lang.gameKeyword" v-model="activeKeyword" class="filter-input"
             @update:model-value="changeFilter()"/>
    <q-btn class="search-btn btn-primary" square unelevated color="primary" icon="filter_alt_off"
           @click="clearFilter()"/>
  </div>
</template>


<script setup lang="ts">
import {onMounted, ref, watch} from 'vue'
import {decodeApiData, isEmpty} from "components/utils";
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {GetFilter} from "app/wailsjs/go/controller/Controller";

const global = useGlobalStore();
const {
  config, lang, platformUi, activePlatform, activeKeyword, activeType, activeYear, activeProducer, activePublisher,
  activeCountry, activeTranslate, activeVersion, activeScore, activeComplete, activeLetter, scrollAreaRef
} = storeToRefs(global);
const filters: any = ref({})

onMounted(() => {
  createFilter();
})

//监听平台变更
watch(activePlatform, (newValue, oldValue) => {
  createFilter();
});

//创建标签
function createFilter() {
  //加载标签
  let platform = activePlatform.value == -1 ? 0 : activePlatform.value
  GetFilter(platform).then((result: string) => {
    let resp = decodeApiData(result)
    filters.value = resp.data
    setSelectDefault()
  })
}

//更改过滤器选项
function changeFilter() {
  activeLetter.value = "ALL"
  scrollAreaRef.value.setScrollPosition("vertical", 0) //滚动条回到最顶端
  global.incRomState()
}

//还原过滤器
function clearFilter() {
  setSelectDefault()
  global.incRomState()
}

function setSelectDefault() {
  activeType.value = filters.value.base_type[0];
  activeYear.value = filters.value.base_year[0];
  activeProducer.value = filters.value.base_producer[0];
  activePublisher.value = filters.value.base_publisher[0];
  activeCountry.value = filters.value.base_country[0];
  activeTranslate.value = filters.value.base_translate[0];
  activeVersion.value = filters.value.base_version[0];
  activeScore.value = filters.value.score[0];
  activeComplete.value = filters.value.complete[0];
  activeKeyword.value = "";
}

</script>

<style scoped>
@import "src/css/classic/common.css";
@import "src/css/classic/headerBarFilter.css";
</style>
