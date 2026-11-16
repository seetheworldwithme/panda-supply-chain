<template>
<!-- 应用权限页面 -->
  <div>
    <!-- tab选项卡 -->
    <div class="menu-list">
      <div
        class="totalTab"
        :class="{ active: item.id === sel }"
        v-for="item in tabs"
        :key="item.id"
        @click="select(item)"
        style="cursor: pointer"
      >
        <span>{{ item.label }}</span>
        <div
          class="stripnone"
          :style="item.width"
          :class="{ strip: item.id === sel }"
        ></div>
      </div>
    </div>

    <div class="parNavs">
        <div v-show="sel === 1">
            <dl class="card">
                <dd>
                    <img src="../../assets/Group.png" alt="">
                    <div class="card-size">
                        <span>项目管理</span>
                        <i class="el-icon-edit"></i>
                    </div>
                </dd>
                <dd>
                    <img src="../../assets/Group2.png" alt="">
                    <div class="card-size">
                        <span>业务管理</span>
                        <i class="el-icon-edit"></i>
                    </div>
                </dd>
                <dd>
                    <img src="../../assets/Union.png" alt="">
                    <div class="card-size">
                        <span>统计管理</span>
                        <i class="el-icon-edit"></i>
                    </div>
                </dd>
            </dl>
        </div>
        <div v-show="sel === 2">
            <dl class="card">

                <dd>
                    <img src="../../assets/Group2.png" alt="">
                    <div class="card-size">
                        <span>项目管理</span>
                        <i class="el-icon-edit"></i>
                    </div>
                </dd>
                <dd>
                    <img src="../../assets/Group.png" alt="">
                    <div class="card-size">
                        <span>统计管理</span>
                        <i class="el-icon-edit"></i>
                    </div>
                </dd>
                <dd>
                    <img src="../../assets/Union.png" alt="">
                    <div class="card-size">
                        <span>业务管理</span>
                        <i class="el-icon-edit"></i>
                    </div>
                </dd>
            </dl>
        </div>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      sel: 1, // 选项卡默认显示tab
      tabs: [
        { label: "已开启应用", id: 1, width: "width:80px;" },
        { label: "未开启应用", id: 2, width: "width:80px;" },
      ],
    };
  },
  methods: {
    //选项卡
    select(item) {
      this.sel = item.id;
    },
  },
};
</script>
<style lang="less" scoped>
// 选项卡开始
.menu-list {
  width: 100%;
  height: 50px;
  background-color: #fff;
  box-shadow: rgba(128, 146, 192, 0.15) 2px 0px 10px;
}
.totalTab {
  width: 7.6%;
  float: left;
  font-size: 16px;
  line-height: 26px;
  font-weight: 650;
  color: #333;
  text-align: center;
  margin: 9px 0 0 35px;
}
.totalTab.active {
  color: #ea222e;
}
.stripnone {
  height: 3px;
  background-color: #ea222e;
  margin: 11px auto 0;
  display: none;
}
.strip {
  display: block;
}
.parNavs {
  width: 100%;
  height: calc(100vh - 110px);
  box-shadow: rgba(128, 146, 192, 0.15) 2px 0 10px 0 inset; //边框内阴影
  background-color: #fff;
  box-sizing: border-box;
  .card {
    dd {
      width: 220px;
      height: 124px;
      box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
      margin: 32px 0 0 32px;
      float: left;
      padding: 27px 0;
      box-sizing: border-box;
      img{
        display: block;
        margin: auto;
      }
      .card-size{
          margin-top: 16px;
        text-align: center;
        font-size: 14px;
        span{
            color: #333;
        }
        i{
            margin-left: 8px;
            color: rgb(130,130,130);
        }
      }
    }
  }
}
//选项卡结束
//卡片内容样式开始

//卡片内容样式结束
</style>