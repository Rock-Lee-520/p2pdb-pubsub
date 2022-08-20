import Vue from "vue";
import App from "./App.vue";
import LemonIMUI from "lemon-imui";
import "lemon-imui/dist/index.css";
Vue.use(LemonIMUI);
Vue.config.productionTip = false;

new Vue({
    render: (h) => h(App)
}).$mount("#app");
