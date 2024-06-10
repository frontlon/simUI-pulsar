import {QDialogOptions} from 'quasar'
import SDialogComponent from "components/SDialogComponent.vue";
import {ref} from "vue";

//弹出dialog
export function getPromptOpts(title: string, msg: string, okLabel: string, persistent: boolean, val: any = null, subTitle = ""): QDialogOptions {
    let data = {
        component: SDialogComponent,
        componentProps: {
            title: title,
            subTitle: subTitle,
            message: msg,
            persistent: persistent,
            okLabel: okLabel,
        }
    }
    if (val !== null) {
        data.componentProps.input = ref(val)
    }
    return data
}
