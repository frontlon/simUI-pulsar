import {defineStore} from 'pinia';

let filterDefault = {"label": "", "value": ""}
export const useGlobalStore = defineStore('globalData', {
    state: () => ({
        rootPath: "", //当前项目根目录
        theme: "", //当前主题
        activePlatformDesc: "", //平台描述
        activePlatform: -1, //平台id
        activeMenu: "", //菜单
        activeRom: 0, //活动ROM
        activeRomDetail: {}, //活动ROM详细信息
        activeFocus: [0, -1], //当前焦点位置 1平台 2菜单 3ROM 4子游戏 5模拟器 (index)
        activeType: filterDefault, //过滤器 游戏类型
        activeYear: filterDefault, //过滤器 发型年份
        activeProducer: filterDefault, //过滤器 制作商
        activePublisher: filterDefault,//过滤器 发行商
        activeCountry: filterDefault,//过滤器  国家
        activeTranslate: filterDefault,//过滤器 汉化组
        activeVersion: filterDefault,//过滤器  版本
        activeScore: filterDefault,//过滤器 评分
        activeComplete: filterDefault,//过滤器 通关状态
        activeKeyword: "",//过滤器 关键字搜索
        activeLetter: "ALL",//过滤器 字母搜索
        menuLike: 0, //是否使用菜单模糊查询
        romState: 0, //rom列表更新状态
        romCount: 0, //当前列表rom数量
        config: {} as any, //通用配置
        lang: {}, //语言定义
        langList: [] as any, //语言列表
        buildTime: "", //应用构建时间
        versionNo: "", //当前版本号
        platformUi: {}, //当前平台界面设置
        rombaseAlias: {},//资料项别名配置
        callbackOpts: {},//监测操作回调
        simulatorMap: {},//平台模拟器列表
        scrollAreaRef: null, //滚动区域ref
        uploadServer: "", //文件上传服务
        logo: "", //侧边栏默认平台图片
    }),
    actions: {
        goto(path: string, query: any = {}, isFresh = 1) {
            this.router.push({path: path, query: query}).then(() => {
                // 重新刷新页面
                if (isFresh == 1) {
                    location.reload()
                }
            })
        },
        //变更rom列表刷新状态
        incRomState() {
            this.romState++;
        },
        //初始化
        initData(data: any) {
            let menu = JSON.parse(data.Config.Menu)
            this.config = data.Config;
            this.lang = data.Lang;
            this.langList = data.LangList;
            this.buildTime = data.BuildTime;
            this.versionNo = data.VersionNo;
            this.activePlatform = data.Config.Platform;
            this.activeMenu = menu[0];
            this.menuLike = menu[1];
            this.rootPath = data.RootPath
            this.activeRom = 0;
            this.uploadServer = data.UploadServer;
            this.logo = data.logo;
        },
        //清空过滤器
        clearFilter() {
            this.activeType = filterDefault;
            this.activeYear = filterDefault;
            this.activeProducer = filterDefault;
            this.activePublisher = filterDefault;
            this.activeCountry = filterDefault;
            this.activeTranslate = filterDefault;
            this.activeVersion = filterDefault;
            this.activeScore = filterDefault;
            this.activeComplete = filterDefault;
            this.activeKeyword = "";
        }
    }
});
