<template>

  <router-view v-slot="{ Component }" :key="$route.path">
    <transition appear enterActiveClass="animate__animated animate__fadeIn">
      <keep-alive>
        <component :is="Component"/>
      </keep-alive>
    </transition>
  </router-view>
</template>


<script setup lang="ts">
import {Dark} from 'quasar'

// 设置状态
Dark.set(true); // 或者false或者"auto"

</script>
