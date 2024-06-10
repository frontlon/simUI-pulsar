<template>
  <q-dialog v-model="showDialog" transition-show="scale" transition-hide="scale">
    <q-card style="min-width: 250px">
      <q-card-section class="row items-center q-pb-none q-mb-md">
        <div class="text-h6">{{ lang.moveModule }}</div>
        <q-space/>
        <q-btn icon="close" flat round v-close-popup/>
      </q-card-section>

      <q-card-section class="q-pt-none">
        <div class="rom-name q-mb-xs" v-if="romInfo">
          <q-icon size="xs" name="sports_esports"/>
          {{ romInfo.Name }}
        </div>

        <q-list padding class="move-list">
          <q-item clickable active-class="submenu-active" @click="activeMenu = '/'" :active="activeMenu == '/'"
                  :disable="'/' == oldPath">
            <q-item-section side>
              <q-avatar icon="folder" color="primary" size="xs"/>
            </q-item-section>
            <q-item-section>
              <q-item-label lines="1">{{lang.notCate}}</q-item-label>
            </q-item-section>
          </q-item>
          <div v-for="item in menuList">
            <q-item clickable active-class="submenu-active" @click="activeMenu = item.Path"
                    :active="activeMenu == item.Path" :disable="item.Path == oldPath">
              <q-item-section side>
                <q-avatar icon="folder" color="primary" size="xs"/>
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ item.Name }}</q-item-label>
              </q-item-section>
            </q-item>
            <div v-if="item.SubMenu" v-for="sub in item.SubMenu">
              <q-item class="sub-tree" clickable active-class="submenu-active" @click="activeMenu = sub.Path"
                      :active="activeMenu == sub.Path" :disable="sub.Path == oldPath ">
                <q-item-section side>
                  <q-avatar icon="folder" color="primary" size="xs"/>
                </q-item-section>
                <q-item-section>
                  <q-item-label lines="1">{{ sub.Name }}</q-item-label>
                </q-item-section>
              </q-item>
            </div>
          </div>
        </q-list>
      </q-card-section>
      <q-card-actions align="right" class="text-primary">
        <q-btn color="primary" square size="md" :label="lang.ok" :disable="activeMenu == null" v-close-popup
               class="update-btn" @click="moveModule()"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script lang="ts">

import {ref} from 'vue'
import {useRoute} from "vue-router";
import {GetMenuList, MoveRomLink} from "app/wailsjs/go/controller/Controller";
import {decodeApiData} from "components/utils";
import {callback} from 'pages/classic/context/Context.vue';
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";
const global = useGlobalStore();
const {lang} = storeToRefs(global);
const route = useRoute();
const tab = ref('base')
const romInfo: any = ref(null)
const romId: any = ref(null)
const romIndex: any = ref(null)
const platformId: any = ref(null)
const showDialog: any = ref(false)
const menuList: any = ref([])
const activeMenu: any = ref(null);
const oldPath: any = ref(null);

export function openMoveModuleDialog(index: number,detail) {
  console.log("openMoveModule", index, detail)
  showDialog.value = true
  romInfo.value = detail;
  romId.value = detail.Id;
  romIndex.value = index;
  platformId.value = detail.Platform;
  oldPath.value = detail.Menu;

  GetMenuList(platformId.value).then((result: string) => {
    let resp = decodeApiData(result)
    console.log(detail.Platform)
    menuList.value = resp.data
    console.log(menuList.value)
  })

}

function moveModule() {
  MoveRomLink(romId.value, activeMenu.value).then((result: string) => {
    let resp = decodeApiData(result)
    if (resp.err == "") {
      callback("moveRomLink", romIndex.value, activeMenu.value)
    }
  })

}

export default {
  setup() {
    return {
      route,
      tab,
      lang,
      romId,
      romInfo,
      oldPath,
      platformId,
      showDialog,
      menuList,
      activeMenu,
      moveModule,
    };
  }
}

</script>

<style scoped>
@import "src/css/manage.css";

.move-list .q-item {
  padding: 0 10px;
}

.sub-tree {
  padding-left: 30px !important;
}

.submenu-active {
  color: var(--color-15);
  background: var(--q-primary);
}

.rom-name {
  font-size: 12px;
  color: var(--color-9);
}
</style>
