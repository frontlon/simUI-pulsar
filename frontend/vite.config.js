import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()]
})

/*
export index ({
    build:{
        rollupOptions:{
            external(source, importer, isResolved) {
                console.log(source, importer. isResolved);
                // 判断是否是/assets-’开头的路径，如果是则认为是外部依赖
                console.log("inininininiini")
                console.log(source)
                return source.startswith("assets-");

            }
        }
    }
})*/
