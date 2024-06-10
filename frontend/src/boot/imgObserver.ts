// 定义一个全局的 Intersection Observer 实例
const observer = new IntersectionObserver(entries => {
    entries.forEach(entry => {

        const lazyImage = entry.target.querySelector('img') as HTMLImageElement;
        console.log(lazyImage);

        if (lazyImage) {
            const src = lazyImage.src;

            lazyImage.dataset.src = lazyImage.dataset.src || src;

            if (entry.isIntersecting) {
                lazyImage.src = lazyImage.dataset.src || '';
               // observer.unobserve(lazyImage);
            } else {
                lazyImage.src = '';
                //lazyImage.removeAttribute("src")

            }
        }

    });
});

// vue3 自定义指令
export default {
    mounted(el: HTMLElement) {
        observer.observe(el);
    },
    unmounted(el: HTMLElement) {
        observer.unobserve(el);
    }
};
