<template>
  <!-- 菜单管理页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav">
        <!-- 表单 -->
        <el-form :model="AdjustForm" size="small" label-position="left" style="width:30%;float: left;">  
            <el-row>
                <el-col :span="24" style="margin-top:1%;">
                    <el-form-item label="是否开启招行支付功能" label-width="180px">
                        <el-radio-group v-model="AdjustForm.isOpenPay">
                        <el-radio :label="true">是</el-radio>
                        <el-radio :label="false">否</el-radio>
                        </el-radio-group>   
                    </el-form-item>
                </el-col>
                <el-col :span="24" style="margin-top:1%;">
                    <el-form-item label="是否开启企业注册功能" label-width="180px">
                        <el-radio-group v-model="AdjustForm.isOpenEtcRegist">
                        <el-radio :label="true">是</el-radio>
                        <el-radio :label="false">否</el-radio>
                        </el-radio-group> 
                    </el-form-item>
                </el-col>
            </el-row>
            <div class="adjusttop">
                <el-button type="primary" size="small" @click="editDatasubmitForm">保存全局配置</el-button>
            </div>
        </el-form>
            <!-- 表格 -->
            <div class="pur-table" style="width: 36%;float: left;">
                <p>小程序-现金红包</p>
                <el-table
                    v-loading="settlementloading"
                    element-loading-text = "加载中..."
                    height="80%"
                    border
                    :data="settlementtableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;border-bottom:0;overflow:visible !important;">
                    <el-table-column prop="id" v-if="show">
                    </el-table-column>
                    <el-table-column label="红包类型" :show-overflow-tooltip="true" width="200">
                        <template slot-scope="scope">
                          <span v-if="scope.row.receiveType == 1">首次登录</span>
                          <span v-if="scope.row.receiveType == 2">发布煤种并通过审核</span>
                          <span v-if="scope.row.receiveType == 3">完善名片资料</span>
                          <span v-if="scope.row.receiveType == 4">完善企业资料</span>
                          <span v-if="scope.row.receiveType == 5">分享转发</span>
                          <span v-if="scope.row.receiveType == 6">关注他人煤种</span>
                          <span v-if="scope.row.receiveType == 7">意见建议</span>
                          <span v-if="scope.row.receiveType == 8">价格更新</span>
                          <span v-if="scope.row.receiveType == 9">完善个人资料</span>
                          <span v-if="scope.row.receiveType == 10">发布视频并通过审核</span>
                          <span v-if="scope.row.receiveType == 11">分享转发视频</span>
                          <span v-if="scope.row.receiveType == 12">邀请新人注册</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="金额最大值（分）" :show-overflow-tooltip="true"  width="200">
                        <template slot-scope="scope">
                            <div v-if="scope.row.maxblone">
                                <el-input size="small" v-model="newmaxe" v-focus @blur="maxFcous(scope.row)"></el-input>
                            </div>
                            <div v-else>
                                <span v-if="scope.row.maxAmount">{{scope.row.maxAmount}} 分</span>
                                <i class="el-icon-edit" @click="maxlclick(scope.row)" style="color:#CF1724;float:right;margin-right:10px;line-height: 23px;cursor:pointer"></i>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column label="金额最小值（分）" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <div v-if="scope.row.sjhxblone">
                                <el-input size="small" v-model="newSjhxje" v-focus @blur="loseFcous(scope.row)"></el-input>
                            </div>
                            <div v-else>
                                <span v-if="scope.row.minAmount">{{scope.row.minAmount}} 分</span>
                                <i class="el-icon-edit" @click="cellclick(scope.row)" style="color:#CF1724;float:right;margin-right:10px;line-height: 23px;cursor:pointer"></i>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
                <div class="putbox">
                  <span>红包领取当日最大上限：</span>
                  <el-input v-model="redUplimitDays" placeholder="金额（单位：分）"></el-input><span> 分</span>
                </div>
                <div class="adjusttop">
                  <el-button type="primary" size="small" @click="moenyForm">更新红包配置</el-button>
                </div>
            </div>


    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
  data() {
    return {
        show:false,
        purTitle: "",                       // 标题
        /*
        编辑
        */
        AdjustForm:{},  

        settlementloading:false,
        settlementtableData: [],                        // 表格数据
        newSjhxje:'',                                   // 最小值input
        newmaxe:'',                                     // 最大值input
        redUplimitDays:100,              // 红包最大额度
    };
  },
  created() {},
  mounted(){
    this.purTitle = this.$route.meta.title;             // 标题
    this.logo = sessionStorage.getItem('companyName')   // 企业名称
    this.purget();
    this.settlementgetpurcon();
  },
  directives: {
        // 注册一个局部的自定义指令 v-focus
        focus: {
            // 指令的定义
            inserted: function (el) {
            // 聚焦元素
            el.querySelector('input').focus()
            }
        }
  },
  methods: {
        // 红包表格数据
        settlementgetpurcon(){
            this.settlementloading = true;
            // 渲染表格
            this.api.brelrge.list()
            .then(res=>{
                this.settlementloading = false;
                if(res.data.code == 200){
                    for(let i=0;i<res.data.data.length;i++){
                        res.data.data[i].sjhxblone = false;
                        res.data.data[i].maxblone = false;
                    }
                    this.settlementtableData = res.data.data;
                    console.log(res.data.data)
                    this.redUplimitDays = res.data.data[0].redUplimitDays;
                }
            })
        },
        // 最小值input失去焦点
        loseFcous(val){
            let eagstr = /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/;
            if(eagstr.test(this.newSjhxje) && this.newSjhxje*1 >= 10){
                if(val.maxAmount*1 > this.newSjhxje*1){
                    val.sjhxblone = false;
                    val.minAmount = this.newSjhxje *1;
                    this.newSjhxje = '';
                }else{
                  this.$message.error('输入的金额最大值小于金额最小值范围！');
                  val.minAmount = val.minAmount *1;
                  val.sjhxblone = false;
                }
            }else{
              this.$message.error('请输入正确金额或金额太小！');
              val.sjhxblone = false;
            }
        },
        // 最小值实际核销修改金额事件
        cellclick(row){
            row.sjhxblone = true;
        },

        // 最大值input失去焦点
        maxFcous(val){
            let eagstr = /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/;
            if(eagstr.test(this.newmaxe) && this.newmaxe*1 >= 10){
                if(val.minAmount*1 < this.newmaxe*1){
                    val.maxblone = false;
                    val.maxAmount = this.newmaxe *1;
                    this.newmaxe = '';
                }else{
                  this.$message.error('输入的金额最大值小于金额最小值范围！');
                  val.maxAmount = val.maxAmount *1;
                  val.maxblone = false;
                }
            }else{
              this.$message.error('请输入正确金额或金额太小！');
                val.maxblone = false;
            }
        },
        // 最大值实际核销修改金额事件
        maxlclick(row){
            row.maxblone = true;
        },
        moenyForm(){
          if(this.redUplimitDays*1 >= 100){
            for(let i=0;i<this.settlementtableData.length;i++){
              this.settlementtableData[i].redUplimitDays = this.redUplimitDays;
              if(this.settlementtableData[i].maxAmount > this.redUplimitDays){
                this.$message({  
                      type: "error",
                      message: "现金红包金额最大值超出红包领取当日最大上限，请重新设置!"
                });
                return 
              }
            }
            this.api.brelrge.updBatch(this.settlementtableData)
            .then(res=>{
              if(res.data.code == 200){
                this.$message({  
                      type: "success",
                      message: "小程序红包更新配置成功!"
                  });
                  this.settlementgetpurcon();
              }
            })
          }else{
            this.$message.error('请输入金额上限！');
          }
        },

    purget(){
        this.api.glconfig.list()
        .then(res=>{
            if(res.data.code == 200){
                this.AdjustForm = res.data.data[0]
            }
        })
    },
    // 编辑提交
    editDatasubmitForm(){
        this.api.glconfig.updateById(this.AdjustForm)
        .then(res=>{
            if(res.data.code == 200){
                this.$message({  
                    type: "success",
                    message: "全局配置设置成功!"
                });
                this.AdjustForm = {}
                this.purget();
            }
        })
    },
  },
};
</script>
<style lang="less" scoped>
.putbox{
  padding-top:20px;
  box-sizing: border-box;
  span{
    font-size: 16px;
    line-height: 30px;
    color: #333;
  }
  .el-input{
    width: 30%;
  }
}
.adjusttop{
  margin-top: 2%;
}
// 头部开始
.pur-top {
  width: 100%;
  height: 50px;
  padding: 9px 24px 9px 32px;
  background-color: #fff;
  box-sizing: border-box;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  border: 1px solid rgba(102, 102, 102, 0.15);
  .pur-size {
    font-size: 16px;
    line-height: 32px;
    color: #333;
    font-weight: 600;
  }
  .pur-right {
    .el-input {
      width: 200px;
      float: right;
    }
    .el-button {
      margin-left: 16px;
      float: right;
    }
  }
}
// 头部结束
// 内容开始
.pur-nav {
  width: calc(100% - 32px);
  height: calc(100vh - 142px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 8px 8px 8px 60px;
  .pur-table{
    height: 82%;
    p{
      font-size: 24px;
      line-height: 48px;
      color: #333;
      font-weight: 700;
    }

  } 
  // 暂无数据样式开始
  .dataAvailable{
    width: 100%;
    height: 245px;
    position: relative;
    margin-top:185px;
    margin-bottom:265px;
    .dataAva{
      position: absolute;
      top:75%;
      left: 25%;
      transform: translate(0, -50%);
      text-align: center;
      img{
        width: 371px;
        height: 200px;
        margin: auto;
      }
      p{
        font-size: 16px;
        line-height: 24px;
        color: #333;
      }
      .el-button{
        margin-bottom:100px;
      }
    }
  }
}
// 内容结束
</style>