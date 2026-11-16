<template>
    <!--  左侧菜单页面 -->  
    <el-container>
      <!-- 左侧导航子菜单开始 -->
      <el-aside class="menus-left" width="170px">
        <el-menu :router="true" :default-active="Application">
          <!-- 左侧一级菜单遍历显示 -->
          <el-menu-item
            v-for="(v, k) in routesSon"
            :key="k"
            :index="v.name"
            :route="{ name: v.name }"
            v-if="!v.children && v.meta.isShow == 0"
            class="menu-one">
            <i :class="v.meta.icon" style="padding-left: 16px"></i>
            <span slot="title">{{ v.meta.title }}</span>
          </el-menu-item>
          <!-- 左侧一级菜单禁用菜单 -->
          <el-menu-item class="menu-stop" v-else-if="v.meta.isShow == 2">
            <span slot="title" style="margin-left: 25px; color: #bfbfbf">{{
              v.meta.title
            }}</span>
          </el-menu-item>
          <!-- 左侧一级菜单 增删改查菜单 -->
          <el-menu-item v-else-if="v.meta.isShow == 1" v-show="false">
          </el-menu-item>
          <!-- 左侧一级下拉菜单显示 -->
          <el-submenu
            class="submenu"
            :index="v.name"
            v-else-if="v.meta.isShow == 0"
          >
            <template slot="title" class="submenu-bottom">
              <i :class="v.meta.icon"></i>
              <span slot="title">{{ v.meta.title }}</span>
            </template>
            <!-- 左侧二级菜单显示 -->
            <el-menu-item-group>
              <el-menu-item
                v-for="(val, is) in v.children"
                :key="is"
                :index="val.name"
                :route="{ name: val.name }"
                ><p style="margin-left: 40px">
                  {{ val.meta.title }}
                </p></el-menu-item
              >
            </el-menu-item-group>
          </el-submenu>
        </el-menu>
      </el-aside>
      <!-- 左侧导航子菜单结束 -->
    <el-main>
        <!-- 内容视口view渲染 -->
        <router-view></router-view>
    </el-main>
    </el-container>
</template>
<script>
export default {
    data() {
        return {
          Application:'Application',
          //导航子路由 
          routesSon: [],
        }
    }, 
    created(){
        //获取动态路由表
        let lengths = this.$router.options.routes.length - 1;
        if(this.$router.options.routes[lengths].children.length >= 2){
          this.routesSon = this.$router.options.routes[lengths].children[1].children;
        }else if(this.$router.options.routes[lengths].children.length < 2){
          this.routesSon = this.$router.options.routes[lengths].children[0].children;
        }
        //获取动态路由表
        if(this.$route.meta.routeName){
          this.Application = this.$route.meta.routeName;
          sessionStorage.setItem('routeName',this.MainPage); 
        }else{
          let newrou = this.$router.options.routes[lengths].children[0].redirect.name;
          this.routesSon.map(item=>{
            if(item.name == newrou){
              this.Application = item.children[0].name;
              sessionStorage.setItem('routeName',this.MainPage); 
              this.$router.push({ name: this.MainPage });
            }
          })
        }
        // this.Application = this.$route.meta.routeName;
        // console.log(this.$route)
        // sessionStorage.setItem('routeName',)
    }
}
</script>
<style  lang="less" scoped>
// 左侧菜单开始
.el-submenu [class^="el-icon-"] {
  color: #333 !important;
  font-size: 16px !important;
  font-weight: 500 !important;
}
.el-menu{
    border: 0;
}
    .menus-left {
    //  box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .8);
    border-right: 1px solid rgba(0, 0, 0, .06);
    background-color: #fff;
    height: calc(100vh - 60px);
    border-top: solid 1px #e6e6e6;
    // box-shadow: 2px 2px 10px rgba(128, 146, 192, 0.15);
    box-sizing: border-box;
    overflow-x: hidden;
    .el-menu-item {
        padding: 0 !important;
        span {
            color: #333;
        }
        i {
            font-weight: 400;
            font-size: 14px;
            color: #333;
        }
    }
    //禁用菜单样式开始
    .menu-stop{
      background-color: #fff;
    }
    .menu-stop:hover {
      background-color: #fff !important;
    }
    .menu-stop.is-active {
      background-color: #fff !important;
    }
    //禁用菜单样式结束
    //一级菜单样式开始
    .menu-one:hover {
      background-color: #fff !important;
    }
    .menu-one:hover span {
      color: #cf1724 !important;
    }
    .menu-one:hover i {
      color: #cf1724 !important;
    }
    .menu-one.is-active {
      background-color: #cf1724 !important;
      span {
        color: #fff !important;
      }
      i {
        color: #fff !important;
      }
    }
    //一级菜单样式结束
    //一级下拉菜单开始
    .submenu {
      //二级菜单开始
      .el-menu-item.is-active {
        background-color: #cf1724 !important;
        p {
          color: #fff !important;
        }
        i {
          color: #fff !important;
        }
      }
      //二级菜单结束
    }
    .submenu:hover span {
      color: #cf1724 !important;
    }
    .submenu:hover i {
      color: #cf1724 !important;
    }
    //一级下拉菜单结束
  }
//左侧菜单结束
// main内容开始
.el-main {
  background-color: #f5f5f5;
  padding: 0;
  overflow: hidden;
}
// mian内容结束
</style>