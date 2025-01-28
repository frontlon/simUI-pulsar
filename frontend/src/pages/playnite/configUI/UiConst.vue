<script lang="ts">
import 'animate.css'
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";
import {ref} from "vue";

const global = useGlobalStore();
const {lang} = storeToRefs(global);

export const boolOptions = [
  {label: lang.value.no, value: 0},
  {label: lang.value.yes, value: 1},
];

export const clickAnimates = [
  {label: "none", value: ""},
  {label: "headShake", value: "animate__headShake"},
  {label: "shakeX", value: "animate__shakeX"},
  {label: "bounce", value: "animate__bounce"},
  {label: "shakeY", value: "animate__shakeY"},
  {label: "heartBeat", value: "animate__heartBeat"},
  {label: "pulse", value: "animate__pulse"},
  {label: "swing", value: "animate__swing"},
  {label: "wobble", value: "animate__wobble"},
  {label: "jello", value: "animate__jello"},
  {label: "rubberBand", value: "animate__rubberBand"},
  {label: "flash", value: "animate__flash"},
];

export const listStyleOptions = [
  {label: lang.value.module + '1', value: 1, image: "/images/list-style/1.png"},
  {label: lang.value.module + '2', value: 2, image: "/images/list-style/2.png"},
  {label: lang.value.module + '3', value: 3, image: "/images/list-style/3.png"},
  {label: lang.value.list, value: 4, image: "/images/list-style/4.png"},
];

export const thumbOptions = [
  {label: lang.value.thumb, value: "thumb"},
  {label: lang.value.snap, value: "snap"},
  {label: lang.value.poster, value: "poster"},
  {label: lang.value.packing, value: "packing"},
  {label: lang.value.title, value: "title"},
  {label: lang.value.cassette, value: "cassette"},
  {label: lang.value.icon, value: "icon"},
  {label: lang.value.gif, value: "gif"},
  {label: lang.value.background, value: "background"},
];

export const sizeOptions = [
  {label: lang.value.perRow1, value: 1},
  {label: lang.value.perRow2, value: 2},
  {label: lang.value.perRow3, value: 3},
  {label: lang.value.perRow4, value: 4},
  {label: lang.value.perRow5, value: 5},
  {label: lang.value.perRow6, value: 6},
  {label: lang.value.perRow7, value: 7},
];

export const marginOptions = [
  {label: lang.value.none, value: 'none'},
  {label: lang.value.xs, value: 'xs'},
  {label: lang.value.sm, value: 'sm'},
  {label: lang.value.md, value: 'md'},
  {label: lang.value.lg, value: 'lg'},
  {label: lang.value.xl, value: 'xl'},
];

export const sortOptions = [
  {label: lang.value.sortNameAsc, value: 1},
  {label: lang.value.sortNameDesc, value: 2},
  {label: lang.value.sortYearAsc, value: 3},
  {label: lang.value.sortYearDesc, value: 4},
  {label: lang.value.sortRatingAsc, value: 5},
  {label: lang.value.sortRatingDesc, value: 6},
  {label: lang.value.sortRunNumDesc, value: 7},
  {label: lang.value.sortLastRunTimeDesc, value: 8},
];

export const directionOptions = [
  {label: lang.value.auto, value: 0},
  {label: '4:3 (' + lang.value.horizontal + ')', value: 1},
  {label: '3:4 (' + lang.value.vertical + ')', value: 2},
  {label: '16:9 (' + lang.value.horizontal + ')', value: 3},
  {label: '9:16 (' + lang.value.vertical + ')', value: 4},
  {label: '3:2 (' + lang.value.horizontal + ')', value: 5},
  {label: '2:3 (' + lang.value.vertical + ')', value: 6},
  {label: '1:1', value: 7},

];

export const titleOptions = [
  {label: lang.value.alias, value: 1}, //别名
  {label: lang.value.filename, value: 2}, //文件名
];

export const listColumnOptions = [
  {value: 'BaseNameEn', label: lang.value.optNameEnLabel},
  {value: 'BaseNameJp', label: lang.value.optNameJpLabel},
  {value: 'BaseType', label: lang.value.baseType},
  {value: 'BaseYear', label: lang.value.baseYear},
  {value: 'BasePublisher', label: lang.value.basePublisher},
  {value: 'BaseProducer', label: lang.value.baseProducer},
  {value: 'BaseCountry', label: lang.value.baseCountry},
  {value: 'BaseTranslate', label: lang.value.baseTranslate},
  {value: 'BaseVersion', label: lang.value.baseVersion},
  {value: 'Score', label: lang.value.rating},
  {value: 'Complete', label: lang.value.complete},
  {value: 'Menu', label: lang.value.optMenuLabel},

];

export const fontsizeOptions = ['12px', '14px', '16px', '18px', '24px', '28px', '32px', '36px', '48px'];

export const fontFamilyOptions = ref([]);

export const backgroundRepeatOptions = [
  {value: 'cover', label: lang.value.cover},
  {value: 'repeat', label: lang.value.repeat},
  {value: 'no-repeat', label: lang.value.noRepeat},
];

export const backgroundFuzzyOptions = [
  {value: 0, label: lang.value.noFuzzy},
  {value: 1, label: '1%'},
  {value: 3, label: '3%'},
  {value: 5, label: '5%'},
  {value: 7, label: '7%'},
  {value: 9, label: '9%'},
  {value: 15, label: '15%'},
  {value: 30, label: '30%'},
  {value: 45, label: '45%'},
]
</script>
