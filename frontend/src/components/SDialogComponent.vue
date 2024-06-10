<template>
  <!-- 请注意这里的dialogRef -->
  <q-dialog ref="dialogRef" :persistent="persistent" @hide="onDialogHide">
    <q-card class="q-dialog-plugin">
      <q-card-section class="row items-center q-pb-none q-pb-md">
        <div class="text-h6">{{ title }}</div>
        <q-space/>
        <q-btn icon="close" flat round @click="onCancelClick"/>
      </q-card-section>
      <q-card-section class="q-pt-none q-pb-none q-mb-xs">
        <div class="sub-title q-mb-xs" v-if="subTitle">
          <q-icon size="xs" name="sports_esports"/>
          {{ subTitle }}
        </div>
        <div class="message" v-if="message">{{ message }}</div>
      </q-card-section>

      <!-- 输入框 -->
      <q-card-section class="q-pt-none q-pb-none" v-if="input">
        <q-input autofocus dense v-model="input.value"></q-input>
      </q-card-section>

      <!-- toggle -->
      <q-card-section class="q-pt-none q-pb-none" v-if="toggle && toggle.value">
        <q-option-group :options="toggle.value.items" type="toggle" v-model="toggle.value.model"/>
      </q-card-section>

      <!-- 确定按钮 -->
      <q-card-actions class="q-pa-md" align="right">
        <q-btn color="primary" square size="md" :label="okLabel" @click="onOKClick()"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script>
import {useDialogPluginComponent} from 'quasar'
import {ref} from "vue";
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);

export default {
  props: {
    title: "",
    subTitle: "",
    message: "",
    persistent: false,
    okLabel: "",
    input: null, //输入框
    toggle: ref(null), //单选列表
  },

  emits: [
    // 必需的；需要指定一些事件，
    // 你的组件将通过useDialogPluginComponent()发出这些事件。
    ...useDialogPluginComponent.emits
  ],

  setup(props) {
    const {dialogRef, onDialogHide, onDialogOK, onDialogCancel} = useDialogPluginComponent()
    return {
      dialogRef,
      onDialogHide,
      lang,
      onOKClick() {
        let mod = {}
        //文本框
        if (props.input) {
          mod['input'] = props.input.value
        }
        if (props.toggle) {
          mod['toggle'] = props.toggle.value.model
        }

        onDialogOK(mod)
      },
      onCancelClick: onDialogCancel
    }
  }
}
</script>

<style scoped>
.message {
  background: var(--color-2);
  color: var(--color-9);
  font-size: 12px !important;
  padding: 10px !important;
}

.sub-title {
  font-size: 12px;
  color: var(--color-9);
}

</style>
