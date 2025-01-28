import {Loading, Notify} from "quasar";
import * as ctl from "app/wailsjs/go/controller/Controller";
import {SetTheme} from "app/wailsjs/go/controller/Controller";
import {useGlobalStore} from "stores/globalData";

const global = useGlobalStore();


//深拷贝
export function deepClone(obj: any) {
    return JSON.parse(JSON.stringify(obj))
}

//检查变量是否为空
export function isEmpty(data: any) {

    if (data == undefined || data == null) {
        return true
    }
    switch (typeof (data)) {
        case "number":
            return data == 0
        case "string":
            return data == "" || data == "0"
        case "boolean":
            return data
        case "object":
            if (Array.isArray(data)) { //数组
                return data.length == 0
            } else { //对象
                return Object.keys(data).length == 0
            }
    }
    return true
}

//调用go服务 同步
export function decodeApiData(resp: string) {
    //console.log(resp)
    return JSON.parse(resp);
}

//加载框
export function loading(opt: string, text: string = "") {
    if (opt == "hide") {
        Loading.hide()
        return
    }
    Loading.show({
        group: 'default',
        boxClass: 'bg-grey-2 text-grey-9',
        spinnerColor: 'primary',
        message: text,
    })
}

//提示框
export function notify(opt: string, text: string) {
    let type = "";
    switch (opt) {
        case "suc": //成功
            type = "positive";
            break;
        case "warn": //警告
            type = "warning";
            break;
        case "err": //错误
            type = "negative";
            break;
        case "info": //信息
            type = "info";
            break;
        default:
            return;
    }

    Notify.create({
            type: type,
            message: text,
            timeout: 1500,
            position: "top",
            progress: true,
        }
    )
}

//根据rgb颜色，读取15个同色
export function getGradientColors(rgbString: string) {
    let rgbValues = rgbString.substring(4, rgbString.length - 1).split(',');
    let rgbColor = rgbValues.map(value => parseInt(value));
    let r = rgbColor[0];
    let g = rgbColor[1];
    let b = rgbColor[2];
    let numSteps: number = 15
    const step = 255 / numSteps;
    const gradientColors = [];
    let newR = 0, newG = 0, newB = 0

    for (let i = 0; i < numSteps; i++) {
        let mR = Math.max(20, Math.floor((255 - r) / step))
        let mG = Math.max(20, Math.floor((255 - g) / step))
        let mB = Math.max(20, Math.floor((255 - b) / step))
        newR = Math.round(r + (i * mR));
        newG = Math.round(g + (i * mG));
        newB = Math.round(b + (i * mB));
        newR = newR > 255 ? 255 : newR < 0 ? 0 : newR
        newG = newG > 255 ? 255 : newG < 0 ? 0 : newG
        newB = newB > 255 ? 255 : newB < 0 ? 0 : newB
        gradientColors.push(`rgb(${newR}, ${newG}, ${newB})`);
    }

    return gradientColors;
}

export function getDate(timestamp: number) {
    let date = new Date(timestamp * 1000);
    let year = date.getFullYear();
    let month = date.getMonth() + 1;
    let day = date.getDate();
    return `${year}-${month}-${day}`;
}

export function getDateTime(timestamp: number) {
    let date = new Date(timestamp * 1000);
    let year = date.getFullYear();
    let month = date.getMonth() + 1;
    let day = date.getDate();
    let hours = date.getHours();
    const minutes = date.getMinutes();
    const seconds = date.getSeconds();
    hours = hours < 10 ? `0${hours}` : hours;
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

//读取当前时间戳
export function getTimestamp() {
    return Math.floor(new Date().getTime() / 1000)
}

//读取当前毫秒时间戳
export function getMilliSeconds() {
    return new Date().getTime();
}


//首字母转小写
export function firstLetterToLower(str: string) {
    if (str == "") {
        return ""
    }
    return str.charAt(0).toLowerCase() + str.slice(1);
}


//调用go服务 同步
export async function callSrv(funcName: keyof typeof ctl, ...args: any[]) {

    const func = ctl[funcName];
    if (typeof func !== 'function') {
        console.log("Method " + funcName + " does not exist")
        return null;
    }

    const result = await func(...args);
    console.log("api:", funcName, args, result);
    const data = JSON.parse(result);
    if (data.err != "") {
        //@todo 错误处理
        return null;
    }
    return data.data;
}

//从路径中读取文件名，不包含扩展名
export function getFileName(p: string) {
    if (isEmpty(p)) {
        return ""
    }
    return p.match(/\/([^/]+)\.\w+$/)[1];
}

//从路径中读取文件扩展名
export function getFileExt(filePath: string) {
    if(filePath == ""){
        return ""
    }
    const parts = filePath.split('.');
    const extension = '.' + parts[parts.length - 1];
    return extension
}

// wails 路径解码
export function wailsPathDecode(p: string) {
    if (p == "") {
        return ""
    }

    if (!p.startsWith("/ASSET/")) {
        return p
    }

    p = p.replace("/ASSET/", "")
    p = p.replace(".asset", "")
    p = p.replace(/_/g, "=");

    return decodeURIComponent(escape(window.atob(p)));
}

export function gotoTheme(theme: string) {
    SetTheme(theme).then(() => {
        switch (theme) {
            case "Default":
                global.goto("/default")
                break
            case "Playnite":
                global.goto("/playnite")
                break
            case "Tiny":
                global.goto("/tiny")
                break
        }
    })
}

//加载字体
export function loadFont(fontName,fontType, fontUrl) {
    return new Promise((resolve, reject) => {
        const style = document.createElement('style');

        const fontFace = `
            @font-face {
                font-family: '${fontName}';
                src: url('${fontUrl}') format('${fontType}');
            }
        `;

        console.log("loadFont",fontFace)

        style.appendChild(document.createTextNode(fontFace));

        document.head.appendChild(style);

        // 检查字体是否加载完成
        const fontFaceSet = new FontFace(fontName, `url('${fontUrl}')`);
        fontFaceSet.load().then(() => {
            document.fonts.add(fontFaceSet);
            resolve(); // 字体加载完成
        }).catch(err => {
            reject(err); // 字体加载失败
        });
    });
}

//调用go服务 异步
/*export async function callSrvAsync(funcName: keyof typeof ctl, ...args: any[]) :Promise<string> {

    if (IS_DEV) {
        return "";
    }

    const func = ctl[funcName];
    if (typeof func !== 'function') {
        console.log("Method" + funcName + "does not exist")
        return "";
    }
    return func(...args)
}

*/