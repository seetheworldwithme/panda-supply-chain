<template>
  <div id="app" @click="isTime">
    <router-view/>
  </div>
</template>
<script>
export default {
  data(){  
    return {
      lastTime: null,               // 最后一次点击的时间
      currentTime: null,            // 当前点击的时间
      timeOut: 360 * 60 * 1000      // 设置超时时间:30分钟
    }
  },
  created () {
    let levle = sessionStorage.getItem('level');
    if(this.$store.state.addRoutes.length <= 0 && this.$route.name !== 'Login'){
      // 企业路由刷新
      if(levle == 1){
        this.$store.dispatch("asynSetRoutes");
        if(this.$store.state.menus.length <= 0){
          this.$store.dispatch("asynSetMenus");
        }
      }
      if(levle == 2){
        this.$store.dispatch("asynSetplatformRoutes");
      }
      if(levle == 3 || levle == 4 || levle == 5 || levle == 6){
        this.$store.dispatch("asynSetgoverRoutes");
      }
    }
    this.lastTime = new Date().getTime();

    // this.$store.dispatch("AllBasics");
    // console.log(this.$store.state.Basics)
  },
  mounted(){
    let levle = sessionStorage.getItem('level');
    if(this.$store.state.addRoutes.length <= 0 && this.$route.name !== 'Login'){
      // 企业路由刷新
      if(levle == 1){
        this.$store.dispatch("asynSetRoutes");
        if(this.$store.state.menus.length <= 0){
          this.$store.dispatch("asynSetMenus");
        }
      }
      if(levle == 2){
        this.$store.dispatch("asynSetplatformRoutes");
      }
      if(levle == 3 || levle == 4 || levle == 5 || levle == 6){
        this.$store.dispatch("asynSetgoverRoutes");
      }
    }
  },
  methods: {
    isTime () {
      this.currentTime = new Date().getTime()                       // 记录这次点击的时间
      if (this.currentTime - this.lastTime > this.timeOut) {        // 判断上次最后一次点击的时间和这次点击的时间间隔是否大于30分钟
        if (sessionStorage.getItem('Authorization')) {              // 如果是登录状态
          sessionStorage.clear();
          this.$router.push({name:'Login'});
        }else {
          this.lastTime = new Date().getTime()
        }
      } else {
          this.lastTime = new Date().getTime()                      // 如果在30分钟内点击，则把这次点击的时间记录覆盖掉之前存的最后一次点击的时间
      }
    }
  }
  // mounted(){
  //   window.addEventListener('beforeunload', e => this.beforeunloadHandler(e))
  // },
  // methods: {
  //   beforeunloadHandler (e) {
  //     e = e || window.event
  //     if (e) {
  //       alert(111)
  //     }
  //   }
  // },
  // destroyed () {
  //   window.removeEventListener('beforeunload', e => this.beforeunloadHandler(e))
  // },
}
</script>
<style lang="less">
.el-table__body tr.current-row>td {
    background: #F6F6F6 !important;
}
* {
  margin: 0;
  padding: 0;
}
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}
::-webkit-scrollbar-thumb {
  background: #e4e4e4;
  border-radius: 60px;
}
.el-tooltip__popper {
    max-width: 50% !important;
  }
</style>
