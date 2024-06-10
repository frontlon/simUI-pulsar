<template>
  <q-dialog v-if="showDialog" persistent filled square ref="dialogRef">
    <q-card flat filled square class="q-dialog-plugin wrapper">

      <q-card-section>
        <h5>发现新版本</h5>

        <q-markup-table flat style="background:none">
          <tbody>
          <tr>
            <td class="text-right">{{ lang.currentVersion }}</td>
            <td class="text-left">{{ versionNo }}</td>
          </tr>
          <tr>
            <td class="text-right">{{ lang.newVersion }}</td>
            <td class="text-left">{{ props.data.Version }}</td>
          </tr>
          <tr>
            <td class="text-right">{{ lang.publishDate }}</td>
            <td class="text-left">{{ props.data.Date }}</td>
          </tr>
          <tr>
            <td class="text-right">{{ lang.upgradeContent }}</td>
            <td class="text-left"><a href="http://www.simui.net/" target="_blank">[{{ lang.gotoWebsite }}]</a></td>
          </tr>
          </tbody>
        </q-markup-table>
      </q-card-section>
      <q-card-section>
        <q-btn square :label="lang.upgradeNot" class="upgrade-btn" v-close-popup></q-btn>
        <q-btn v-if="props.auto == 1" square :label="lang.upgradeJump" class="upgrade-btn" @click="jump()"></q-btn>
        <q-btn square :label="lang.upgradeStart" class="upgrade-btn" color="primary" @click="start()"></q-btn>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import {Dialog, useDialogPluginComponent} from 'quasar'
import {useGlobalStore} from 'stores/globalData';
import {storeToRefs} from 'pinia';
import {defineProps, onMounted, ref} from "vue";
import {JumpUpgrade, DownloadNewVersion, InstallUpgrade} from "app/wailsjs/go/controller/Controller";
import {decodeApiData, loading, notify} from "components/utils";
import {getPromptOpts} from "components/dialog";

const props = defineProps(['data', 'auto'])
const {dialogRef} = useDialogPluginComponent()
const global = useGlobalStore();
const {versionNo, lang} = storeToRefs(global);
const showDialog = ref(true)
onMounted(() => {
  console.log("onMounted upda", props.data.Version)
})

//跳过这个版本
function jump() {
  JumpUpgrade(props.data.Version).then((result) => {
    showDialog.value = false
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
    }
  })
}

//开始升级
function start() {
  loading("show", lang.value.upgradeDownloading)
  DownloadNewVersion(props.data.Url).then((result) => {
    loading("hide")
    let resp = decodeApiData(result)
    if (resp.err != "") {
      notify("err", resp.err)
      return
    }

    let zipPath = resp.data

    //提示安装更新
    let opt = getPromptOpts(lang.value.upgradeDownSuc, lang.value.upgradeInstallMsg, lang.value.ok, true)
    Dialog.create(opt).onOk(() => {
      InstallUpgrade(props.data.Version, zipPath).then((result) => {
        let resp = decodeApiData(result)
        if (resp.err != "") {
          notify("err", resp.err)
          return
        }
      })
    })
  })
}

</script>

<style scoped>
.wrapper {
  width: 450px;
  text-align: center;
  border: 1px solid var(--color-3);
}

h5 {
  margin: 0 0 10px 0;
  font-size: 18px;
}

.upgrade-btn {
  margin: 0 5px
}

a, a:visited {
  color: var(--q-primary);
}

</style>
