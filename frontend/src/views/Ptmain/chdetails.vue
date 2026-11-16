<!-- 兑换明细页 -->
<template>
    <!-- 兑换明细页面 -->
    <div>
      <!-- 头部开始 -->
      <el-row class="pur-top">
        <el-col :span="3">
          <span class="pur-size">{{ purTitle }}</span></el-col>
          <el-col :span="21" class="pur-right">
              <el-input
              @keyup.enter.native="search"
              size="small"
              placeholder="请输入内容"
              prefix-icon="el-icon-search"
              v-model="purSearch">
              </el-input>
          </el-col>
      </el-row>
      <!-- 头部结束 -->
      <!-- 内容开始 -->
      <div class="pur-nav"
      v-loading="loading"
      element-loading-text = "加载中...">
        <!-- 表格 -->
        <div class="pur-table">
          <el-table
            height="100%"
            border
            :data="tableData"
            :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
            style="width: 100%;">
                <el-table-column prop="id" v-if="show"></el-table-column>
                <el-table-column label="用户手机号" prop="userPhone"></el-table-column>
                <el-table-column label="物流单号" prop="trackingNumber">
                    <template slot-scope="scope">
                        <span v-if="scope.row.trackingNumber">{{scope.row.trackingNumber}}</span>
                        <span v-else>--</span>
                    </template>
                </el-table-column>
                <el-table-column label="兑换说明" :show-overflow-tooltip="true" prop="exchangeInstructions">
                </el-table-column>
                <el-table-column label="兑换商品名称" :show-overflow-tooltip="true" prop="proName">
                </el-table-column>
                <el-table-column label="消耗积分数量" :show-overflow-tooltip="true" prop="consumeScore">
                </el-table-column>
                <el-table-column label="兑换商品名称" :show-overflow-tooltip="true" prop="proName">
                </el-table-column>
                <el-table-column label="收货人姓名" :show-overflow-tooltip="true" prop="consigneeName">
                </el-table-column>
                <el-table-column label="收货人电话" :show-overflow-tooltip="true" prop="consigneePhone">
                </el-table-column>
                <el-table-column label="收货人地址" :show-overflow-tooltip="true" prop="consigneeAddress">
                </el-table-column>
                <el-table-column label="创建时间" :show-overflow-tooltip="true" prop="createTime">
                </el-table-column>
              <el-table-column
                label="操作" width="80" fixed="left" v-if="$route.meta.edit">
                <template slot-scope="scope">
                  <el-button @click="EditPurs(scope.row)" size="mini">编辑</el-button>
                </template>
              </el-table-column>
                <!-- 暂无数据 -->
                <div slot="empty" class="dataAvailable">
                    <div class="dataAva">
                        <img src="@/assets/zanwushuju.png" alt="">
                        <p>暂无数据，您可以新建试试</p>
                        <!-- <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新增</el-button> -->
                    </div>
                </div>
          </el-table>
        </div>
        <!-- 分页 -->
        <div class="pur-pages">
          <el-pagination
          background
          :page-sizes="[10, 30, 50, 100, 200]"
          :page-size.sync="size"
          :pager-count="5"
          :current-page.sync="current"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @current-change='getpurcon'
          @size-change='getpurcon'>
          </el-pagination>
        </div>
        <!-- 编辑标签信息弹窗 -->
        <el-dialog
          title="编辑物流信息"
          :closeOnClickModal="false"
          :visible.sync="editpaymentShow"
          width="30%"
          :before-close="editpaymentShowClose">
            <el-form style="width:95%;" :model="editData" :rules="editpaymentrules" ref="editData" size="small" label-width="120px">
              <el-form-item label="物流单号" prop="trackingNumber">
                <el-input v-model="editData.trackingNumber" placeholder="请输入物流单号"></el-input>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="editpaymentsubmitm('editData')" style="float:right;margin:5% 0 0 2%;">确定</el-button>
                <el-button @click="editpaymentresetForm" style="float:right;margin:5% 0 0 2%;">取消</el-button>
              </el-form-item>
            </el-form>
        </el-dialog>
      </div>
      <!-- 内容结束 -->
    </div>
  </template>
  <script>
  export default {
      data() {
          return {
              loading:false,
              current:1,              // 当前页数
              size:50,                // 显示条数
              total:0,                // 总条数
              show:false,
              purTitle: "",           // 标题
              purSearch: "",          // 搜索      
              tableData: [],          // 表格数据  
              // 标签编辑
              editpaymentShow:false,
              editData:{},
              editpaymentrules: {
                trackingNumber: [
                      { required: true, message: '请输入物流单号', trigger: 'blur' },
                  ],
              },
          };
      },
      created() {},
      mounted(){
          this.purTitle = this.$route.meta.title;
          this.getpurcon();
      },
      methods: {
          // 表格数据
          getpurcon(){
              this.loading = true;
              // 传参
              var listByPageData = {
                  searchStr:this.purSearch,
                  page:{
                      current: this.current,
                      size: this.size
                  },
              }
              // 渲染表格
              this.api.busechlog.all(listByPageData)
              .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.tableData = res.data.data.records;
                    this.total = res.data.data.total;
                }
              })
          },
          // 搜索功能
          search() {
              this.tableData = [];
              this.getpurcon();
              this.current = 1;
          },
          // 编辑标签信息
          EditPurs(data){
              this.editpaymentShow = true;
              this.editData = JSON.parse(JSON.stringify(data));
          },
          // 编辑标签信息取消
          editpaymentShowClose(){
              this.editpaymentShow = false;
              this.getpurcon();
              this.editData = {}
          },
          // 取消编辑标签信息提交
          editpaymentresetForm() {
              this.editpaymentShow = false;
              this.getpurcon();
              this.editData = {}
          },
          // 编辑标签信息提交
          editpaymentsubmitm(formName) {
              this.$refs[formName].validate((valid) => {
              if (valid) {
                  this.api.busechlog.edit(this.editData)
                  .then(res=>{
                  if(res.data.code == 200){
                      this.$message({
                          type: 'success',
                          message: '编辑成功!'
                      });
                      this.editData = {}
                      this.tableData = [];
                      this.getpurcon();
                      this.current = 1;
                      this.editpaymentShow = false;
                  }
                  })
              }else if(res.data.code == 602){
                  this.$message({
                      type: 'warning',
                      message: res.data.msg
                  });
              }
              });
          },
      },
  };
  </script>
  <style lang="less" scoped>
  // 详情样式开始
  .detilsBox{
    width: 100%;
    height: 100%;
    padding: 25px;
    box-sizing: border-box;
    h5{
      font-size: 16px;
      line-height: 58px;
      color: #333333;
    }
    .detilconent{
      width: 100%;
      span{
        display: inline-block;
        width: 33.33%;
        font-size: 14px;
        line-height: 24px;
        color: #333333;
      }
    }
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
    padding: 8px;
    .pur-table{
      width:100%;
      height: 94%;
    } 
    // 暂无数据样式开始
    .dataAvailable{
      width: 100%;
      height: 245px;
      position: relative;
      margin-top:185px;
      margin-bottom:279px;
      .dataAva{
        position: absolute;
        top:75%;
        left: 50%;
        transform: translate(-50%, -50%);
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
    // 暂无数据样式结束
    .pur-pages {
      width: 100%;
      height: 6%;
      padding: 6px 0 0 12px;
      box-sizing:border-box;
      .el-pagination {
        float: right;
      }
    }
  }
  // 内容结束
  </style>