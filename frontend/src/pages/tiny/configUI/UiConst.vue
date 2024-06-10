<script lang="ts">
import 'animate.css'
import {useGlobalStore} from "stores/globalData";
import {storeToRefs} from "pinia";

const global = useGlobalStore();
const {lang} = storeToRefs(global);

export const boolOptions = [
  {label: lang.value.no, value: 0},
  {label: lang.value.yes, value: 1},
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

export const titleOptions = [
  {label: lang.value.alias, value: 1}, //别名
  {label: lang.value.filename, value: 2}, //文件名
];

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

export const fontsizeOptions = ['16px', '18px', '24px', '28px', '32px', '36px', '48px'];

</script>
