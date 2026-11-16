import Vue from 'vue'
import 'element-ui/lib/theme-chalk/index.css';
import App from './App.vue'
import 'u-reset.css'
import router from './router'
import Filters from './filters'
import ElementUI from 'element-ui'
import api from './api'
import store from './store'
import Vuex from 'vuex'

import promptMessages from '../src/components/messages/index'
import NumberCounter from '../src/components/NumberCounter'
import 'lib-flexible/flexible'
import '../theme/index.css'
import '../src/assets/icon/iconfont.css'
import VueParticles from 'vue-particles'

// import 'video.js/dist/video-js.css'
Vue.use(VueParticles)


Vue.config.ignoredElements = ['wx-open-launch-weapp']
// 图片放大缩小
import Viewer from 'v-viewer'
import 'viewerjs/dist/viewer.css'
Vue.use(Viewer);
Viewer.setDefaults({
  toolbar: false
});
Vue.use(ElementUI, { zIndex: 1000 })

Vue.prototype.api = api
Vue.config.productionTip = false
Vue.use(ElementUI)
Vue.use(Filters)
Vue.use(promptMessages)
Vue.use(NumberCounter)
Vue.use(Vuex)

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app')


