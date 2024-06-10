//是否为开发模式
export const IS_DEV = false;

//模拟器常用启动参数
export const SIMULATOR_BOOT_PARAM_OPTIONS = [
    {label: 'Default', value: '{RomFullPath}',},
    {label: 'MAME', value: '{RomName}',},
    {label: 'Winkawaks', value: '{RomName}',},
    {label: 'FBA Shuffle', value: '-g {RomName} -w',},
    {label: 'Demul', value: '-run=dc -image={RomFullPath}',},
    {label: 'ePsXe', value: '-loadbin {RomFullPath} -nogui',},
    {label: 'Dolphin', value: '-e {RomFullPath}',},
    {label: 'OpenBOR', value: 'load {RomFullPath}',},
    {label: 'RetroArch', value: '-L cores\\xxx.dll {RomFullPath}',},
    {label: 'DosBox', value: '-conf {RomFullPath} -exit',},
    {label: 'DuckStation', value: '{RomFullPath} -batch',},
    {label: 'Nebula2', value: '{RomName}',},
    {
        label: 'NullDC',
        value: '-config nullDC:Emulator.Autostart=1 -config ImageReader:LoadDefaultImage=1 -config ImageReader:DefaultImage={RomFullPath}',
    },
    {label: 'Yabuse', value: '–iso={RomFullPath}',},
    {label: 'PCem', value: '–config {RomFullPath} -f',},
    {label: 'Cemu', value: '-g {RomFullPath} -f',},
];

//编辑器工具栏
export const EDITOR_TOOLBAR = [
    [
        'viewsource',
        {
            label: "",
            icon: "format_align_left",
            list: 'only-icons',
            options: ['left', 'center', 'right', 'justify']
        },
        {
            label: "",
            icon: "title",
            list: 'no-icons',
            options: ['p', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6']
        },
        'bold', 'italic', 'strike', 'underline', 'subscript', 'superscript',
        'code', 'quote', 'unordered', 'ordered', 'outdent', 'indent',
        'image', 'link', 'hr',
        'removeFormat', 'fullscreen'
    ],
]

//菜单图标大小
export const CONTEXT_ICON_SIZE = {
    width: '20px',
    height: '20px',
}