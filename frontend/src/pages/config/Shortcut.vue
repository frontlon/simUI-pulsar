<template>
  <q-list class="shortcut-list">
    <q-item v-for="(item,index) in shortcuts" :key="item.Id">
      <q-item-section>
        <q-input filled square dense :label="lang.shortcutName" v-model="item.Name"/>
      </q-item-section>
      <q-item-section>
        <q-input filled square dense :label="lang.shortcutPath" v-model="item.Path">
          <template v-slot:append>
            <q-icon name="upload_file" class="open-dialog" @click="openFileDialog(index)"/>
          </template>
        </q-input>
      </q-item-section>
      <q-item-section side>
        <q-btn flat square dense icon="delete" @click="delShortcut(index)"/>
      </q-item-section>
    </q-item>
  </q-list>
  <div class="q-gutter-sm bottom-bar">
    <q-btn flat square :label="lang.add" size="lg" padding="md" color="primary" @click="addShortcut()"/>
    <q-btn flat square :label="lang.update" size="lg" padding="md" color="primary" @click="updateShortcut()"/>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {decodeApiData, notify} from "components/utils";
import Sortable from 'sortablejs';
import {GetShortcuts, OpenFileDialog, UpdateShortcut, UpdateShortcutSort} from 'app/wailsjs/go/controller/Controller'
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);
const config = ref(null)
const shortcuts: any = ref([])


onMounted(() => {
  //读取快捷工具
  GetShortcuts(false).then((result: string) => {
    let resp = decodeApiData(result)
    shortcuts.value = resp.data;
  })

  //快捷工具排序
  shortcutDropSort()
})

//添加快捷工具
function addShortcut() {
  shortcuts.value.push({})
}

//更新快捷工具
function updateShortcut() {
  UpdateShortcut(JSON.stringify(shortcuts.value)).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      notify("suc", lang.value.updateSuc)
    } else {
      notify("err", resp.err)
    }
  })
}

//删除快捷工具
function delShortcut(idx: number) {
  shortcuts.value.splice(idx, 1)
}

//选择文件
function openFileDialog(index: number) {
  OpenFileDialog("app").then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {
      shortcuts.value[index].Path = resp.data;
    }
  })
}

//快捷软件拖拽排序
function shortcutDropSort() {
  const el = document.querySelectorAll('.shortcut-list')[0] as HTMLElement;
  Sortable.create(el, {
    animation: 150,
    sort: true,
    disabled: false,
    // 结束拖拽后的回调函数
    onEnd: (evt: any) => {
      let currentRow = shortcuts.value.splice(evt.oldIndex, 1)[0];
      shortcuts.value.splice(evt.newIndex, 0, currentRow);
      //更新数据
      let req: any = [];
      shortcuts.value.forEach((item, index) => {
        shortcuts.value[index].Sort = index + 1
        req.push(item.Id);
      });

      UpdateShortcutSort(JSON.stringify(req)).then((result: string) => {
      })
    },
  });
}

</script>
<style scoped>
@import "src/css/manage.css";
</style>