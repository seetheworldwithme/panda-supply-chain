import Vue from 'vue'
import 'element-ui/lib/theme-chalk/index.css';
import App from './App.vue'
import 'u-reset.css'
import router from './router'
import Filters from './filters'
import ElementUI from 'element-ui'
import api from './api'
import { onApiFail, offApiFail } from './api/http'
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
// 任意请求失败后,把页面上仍为 true 的 *loading* 状态复位:
// 页面普遍写成 loading=true → 请求 → 只在 code==200 时复位,HTTP 500 或业务失败后遮罩 / 保存锁永不解除。
Vue.mixin({
  mounted() {
    if (this._apiFailHandler) return
    this._apiFailHandler = () => {
      Object.keys(this.$data).forEach(key => {
        if (/loading/i.test(key) && this[key] === true) this[key] = false
      })
    }
    onApiFail(this._apiFailHandler)
  },
  beforeDestroy() {
    if (this._apiFailHandler) { offApiFail(this._apiFailHandler); this._apiFailHandler = null }
  }
})
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


