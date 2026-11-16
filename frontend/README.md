# dev
WangYanFeng 2021-2-26
## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

### 项目布局

```
serve-admin
src -- 源码目录
├── api -- axios网络请求定义
├── assets -- 静态图片资源文件
├── components -- 通用组件封装
├── router -- vue-router路由配置
├── store -- vuex的状态管理
└── views -- 前端页面
    └── Mindex -- 企业首页侧导航菜单所有页面
        ├── Account.vue -- 账户维护页面
        ├── AddAccount.vue -- 新建账户信息页面
        ├── Addcollect.vue -- 新建地址维护页面
    └──  Ptmain -- 平台首页侧导航菜单所有页面
    └──  setup -- 企业设置侧导航菜单所有页面
    ├── CommonView.vue -- 首页侧导航栏
    ├── AsidsMain.vue -- 设置侧导航栏
    ├── Forget.vue -- 忘记密码页
    ├── Login.vue -- 登录页面
    ├── Main.vue -- 首页
    ├── NotFont.vue -- 404页面
    ├── Policy.vue -- 隐私政策页面
    ├── Register.vue -- 立即注册页面
    ├── Service.vue -- 平台服务协议页面
    ├── Xzmain.vue -- 个地方政府行政区划页面
    └── NotFont.vue -- 404页面
```