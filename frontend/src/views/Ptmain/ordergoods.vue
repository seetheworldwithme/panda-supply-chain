<!-- 订单记录页面 -->
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
              <el-table-column label="订单创建时间" :show-overflow-tooltip="true" prop="createTime" width="160">
              </el-table-column>
              <el-table-column label="煤仓号" :show-overflow-tooltip="true" prop="warehouseNb" width="100">
              </el-table-column>
              <el-table-column label="发布人手机号" :show-overflow-tooltip="true" prop="createPhone" width="140">
              </el-table-column>
              <el-table-column label="标的指标" :show-overflow-tooltip="true" width="300">
                <template slot-scope="scope">
                    <span v-if="scope.row.pstd!=0&&scope.row.pstd!=null">硫(St,d)%≤{{scope.row.pstd}},</span>
                    <span v-if="scope.row.pashcontent!=0&&scope.row.pashcontent!=null">灰分(Ad)%≤{{scope.row.pashcontent}},</span>
                    <span v-if="scope.row.pwatercontent!=0&&scope.row.pwatercontent!=null">水分(Mt)%≤{{scope.row.pwatercontent}},</span>
                    <span v-if="scope.row.pvolatiles!=0&&scope.row.pvolatiles!=null">挥发分(Vdaf)%≤{{scope.row.pvolatiles}},</span>
                    <span v-if="scope.row.pbondindex!=0&&scope.row.pbondindex!=null">粘结指数GR.I≥{{scope.row.pbondindex}},</span>
                </template>
              </el-table-column>
              <el-table-column label="参考价格（元）" :show-overflow-tooltip="true" prop="exFactoryPrice" width="160">
              </el-table-column>
              <el-table-column label="收货单位名称" :show-overflow-tooltip="true" prop="receivingCompany">
              </el-table-column>
              <el-table-column label="需求吨数" :show-overflow-tooltip="true" prop="demandQuantity">
              </el-table-column>
              <el-table-column label="收货人姓名" :show-overflow-tooltip="true" prop="receivingPeople">
              </el-table-column>
              <el-table-column label="收货地址" :show-overflow-tooltip="true" prop="receivingAddress">
              </el-table-column>
              <el-table-column label="收货电话" :show-overflow-tooltip="true" prop="receivingPhone">
              </el-table-column>
              <el-table-column label="在线留言" :show-overflow-tooltip="true" prop="leavelMsg">
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
                size:100,               // 显示条数
                total:0,                // 总条数
                show:false,
                purTitle: "",           // 标题
                purSearch: "",          // 搜索      
                tableData: [],          // 表格数据  
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
                    isPtSearch:true,
                    searchStr:this.purSearch,
                    page:{
                        current: this.current,
                        size: this.size
                    },
                }
                // 渲染表格
                this.api.order.all(listByPageData)
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

    