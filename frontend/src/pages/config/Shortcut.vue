<template>
  <q-card flat bordered class="config-section">
    <q-card-section>
      <div class="section-title">{{ lang.shortcut }}</div>

      <q-list v-if="shortcuts.length > 0" class="shortcut-list">
        <q-item v-for="(item,index) in shortcuts" :key="item.Id ?? index" class="shortcut-item">
          <q-item-section>
            <q-input filled square :label="lang.shortcutName" v-model="item.Name"/>
          </q-item-section>
          <q-item-section>
            <q-input filled square :label="lang.shortcutPath" v-model="item.Path">
              <template v-slot:append>
                <q-icon name="upload_file" class="open-dialog" @click="openFileDialog(index)"/>
              </template>
            </q-input>
          </q-item-section>
          <q-item-section side>
            <q-btn flat icon="delete" @click="delShortcut(index)"/>
          </q-item-section>
        </q-item>
      </q-list>
    </q-card-section>
  </q-card>

  <div class="bottom-bar">
    <q-btn unelevated :label="lang.add" color="primary" outline @click="addShortcut()"/>
    <q-btn unelevated :label="lang.update" color="primary" class="update-btn" @click="updateShortcut()"/>
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
const shortcuts: any = ref([])


onMounted(() => {
  GetShortcuts(false).then((result: string) => {
    let resp = decodeApiData(result)
    shortcuts.value = resp.data;
  })

  shortcutDropSort()
})

function addShortcut() {
  shortcuts.value.push({})
}

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

function delShortcut(idx: number) {
  shortcuts.value.splice(idx, 1)
}

function openFileDialog(index: number) {
  OpenFileDialog("app").then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "" && resp.data != "") {
      shortcuts.value[index].Path = resp.data;
    }
  })
}

function shortcutDropSort() {
  const el = document.querySelectorAll('.shortcut-list')[0] as HTMLElement;
  if (!el) return;
  Sortable.create(el, {
    animation: 150,
    sort: true,
    disabled: false,
    onEnd: (evt: any) => {
      let currentRow = shortcuts.value.splice(evt.oldIndex, 1)[0];
      shortcuts.value.splice(evt.newIndex, 0, currentRow);
      let req: any = [];
      shortcuts.value.forEach((item, index) => {
        shortcuts.value[index].Sort = index + 1
        req.push(item.Id);
      });

      UpdateShortcutSort(JSON.stringify(req)).then(() => {
      })
    },
  });
}

</script>
<style scoped>
@import "src/css/page/config.css";
</style>
