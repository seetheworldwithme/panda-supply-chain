<!-- 红包记录页面 -->
<template>
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
              <el-select v-model="carriageName" size="small" @change="prologistics">
                    <el-option
                    v-for="item in Extended"
                    :key="item.value"
                    :label="item.dkey"
                    :value="item.dvalue">
                    </el-option>
                </el-select>
              <el-date-picker
              clearable
              size="small"
              @change="usrTie"
              v-model="usrTimArray"
              type="daterange"
              unlink-panels
              range-separator="~"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="yyyy-MM-dd"
              style="width:20%;float:right;margin-right: 10px;">
              </el-date-picker>
              <div class="ljnum">累计领取金额：{{nums/100}}元</div>
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
              <el-table-column prop="id" v-if="show">
              </el-table-column>
              <!-- <el-table-column label="姓名" :show-overflow-tooltip="true" prop="phone">
              </el-table-column> -->
              <el-table-column label="手机号" :show-overflow-tooltip="true" prop="userPhone" width="130">
              </el-table-column>
              <el-table-column label="领取金额" :show-overflow-tooltip="true" width="100">
                <template slot-scope="scope">
                    <span v-if="scope.row.receiveAmount">{{ scope.row.receiveAmount*1/100 }} 元</span>
                </template>
              </el-table-column>
              <el-table-column label="领取途径" :show-overflow-tooltip="true" width="200">
                <template slot-scope="scope">
                    <span v-if="scope.row.receiveType == 1">首次登录</span>
                    <span v-if="scope.row.receiveType == 2">发布煤种并审核通过</span>
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
              <el-table-column label="领取时间" :show-overflow-tooltip="true" prop="createTime">
              </el-table-column>
            <!-- 暂无数据 -->
            <div slot="empty" class="dataAvailable">
                <div class="dataAva">
                  <img src="@/assets/zanwushuju.png" alt="">
                  <p>暂无数据！</p>
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
                size:200,               // 显示条数
                total:0,                // 总条数
                delmul:[],              // 批量删除
                show:false,
                purTitle: "",           // 标题
                purSearch: "",          // 搜索      
                tableData: [],          // 表格数据  
                // 下拉筛选
                Extended:[
                    {
                        dkey:'全部',
                        dvalue:0
                    },
                    {
                        dkey:'首次登录',
                        dvalue:1
                    },
                    {
                        dkey:'发布煤种并审核通过',
                        dvalue:2
                    },
                    {
                        dkey:'完善名片资料',
                        dvalue:3
                    },
                    {
                        dkey:'完善企业资料',
                        dvalue:4
                    },
                    {
                        dkey:'分享转发',
                        dvalue:5
                    },
                    {
                        dkey:'关注他人煤种',
                        dvalue:6
                    },
                    {
                        dkey:'意见建议',
                        dvalue:7
                    },
                    {
                        dkey:'价格更新',
                        dvalue:8
                    },
                    {
                        dkey:'完善个人资料',
                        dvalue:9
                    },
                    {
                        dkey:'发布视频并通过审核',
                        dvalue:10
                    },
                    {
                        dkey:'分享转发视频',
                        dvalue:11
                    },
                    {
                        dkey:'邀请新人注册',
                        dvalue:12
                    },
                ], 
                carriageName:0,
                usrTimArray:[],             // 时间搜索
                nums:0,   // 累计金额
          };
      },
      created() {},
      mounted(){
          this.purTitle = this.$route.meta.title;
          this.usrTie(); // 进页面就要查累计金额:原先只在改日期时查,一进页面显示 0 元
      },
      methods: {
            // 日期搜搜
            usrTie() {
              this.nums = 0;
              this.tableData = [];
              this.current = 1;
              this.getpurcon();
              this.api.brelord.totoalMoney({
                searchStr:'',
                startTime:this.usrTimArray != null?this.usrTimArray[0]:"",
                endTime:this.usrTimArray != null?this.usrTimArray[1]:"",
              })
              .then(res=>{
                if(res.data.code == 200){
                  this.nums = res.data.data*1;
                }
              })
            },
            // 下拉查询
            prologistics(){
                this.nums = 0;
                this.tableData = [];
                this.current = 1;
                this.getpurcon();
            },
            // 表格数据
            getpurcon(){
                    this.loading = true;
                    // 传参
                    var listByPageData = {
                        receiveType:this.carriageName==0?null:this.carriageName,
                        searchStr:this.purSearch,
                        startTime:this.usrTimArray != null?this.usrTimArray[0]:"",
                        endTime:this.usrTimArray != null?this.usrTimArray[1]:"",
                        page:{
                            current: this.current,
                            size: this.size
                        },
                    }
                    // 渲染表格
                    this.api.brelord.all(listByPageData)
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
              this.nums = 0;
              this.tableData = [];
              this.current = 1;
              this.getpurcon();
          },
      },
  };
  </script>
  <style lang="less" scoped>
  // 头部开始
  .pur-top {
    width: 100%;
    height: 50px;
    padding: 9px 24px 9px 32px;
    background-color: #fff;
    box-sizing: border-box;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    border: 1px solid rgba(102, 102, 102, 0.15);
    .ljnum{
      float: right;
      padding-right:16px;
      font-size: 18px;
      line-height: 32px;
    }
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
      .el-select{
        width: 11%;
        margin-right: 10px;
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
    }
  }
  // 内容结束
  </style>

    