<template>
  <div :class="{ 'drop-highlight': isDragOver[typ] }" style="border: 1px solid #f00;width: 100px;height: 100px;"
       @dragover.prevent="handleDragOver($event,'over','thumb')"
       @dragleave="handleDragOver($event,'leave','thumb')"
       @drop.prevent="handleDrop($event,'thumb')">
    将文件拖拽到此处上传

    <q-img v-for="item in images" :src="item" style="width: 100px;height:100px;"></q-img>

  </div>
</template>
<script setup lang="ts">


//图集拖拽上传
import {callSrv, decodeApiData, notify} from "components/utils";
import {CreateRomResByBase64, GetBaseConfig, GetConfig} from "app/wailsjs/go/controller/Controller";

import {onMounted, ref} from "vue";
import axios from "axios";

const isDragOver: any = ref({})
const images: any = ref([])
let uploadServer = ""
onMounted(async () => {
  //加载配置
  // GetConfig().then((result: string) => {
  //   let resp = decodeApiData(result)
  //   uploadServer = resp.data.UploadServer
  // })

  if (uploadServer == "") {
    console.log("aaaaaa", uploadServer)

    uploadServer = await callSrv("StartUploadServer")
    console.log("uploadServer", uploadServer)
  }


})


function handleDragOver(event: any, opt: string, type: string) {
  isDragOver.value[type] = opt == "over";
}

function handleDrop(event: any, type: string) {
  isDragOver.value[type] = false;
  const files = event.dataTransfer.files;
  if (files.length > 0) {
    const file = files[0];

    uploadFile(file)

  }
}

function uploadFile(file) {

  const formData = new FormData();
  formData.append('file', file);
  formData.append('id', "1561");
  formData.append('type', "snap");

  axios.post(uploadServer+'/uploadFile', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  }).then(response => {
    console.log('上传成功', response.data);
  }).catch(error => {
    console.error('上传失败', error);
  });












}

</script>

<style>

.drop-highlight {
  background: yellow;
}
</style>
