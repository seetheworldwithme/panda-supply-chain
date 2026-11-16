<template>
  <!-- 清洗数据页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur">刷新</el-button>
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch">
        </el-input>
        <el-select v-model="carriageName" size="small" @change="prologistics" placeholder="数据状态">
            <el-option
            v-for="item in Extended"
            :key="item.id"
            :label="item.name"
            :value="item.id">
            </el-option>
        </el-select>
      </el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav"
>
      <!-- 表格 -->
      <div class="pur-table"
      v-loading="conversionloading"
      element-loading-text = "转换中...">
        <el-table
          height="100%"
          border
          v-loading="loading"
          element-loading-text = "加载中..."
          :data="tableData"
          :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
          style="width: 100%;">
            <el-table-column prop="id" v-if="show">
            </el-table-column>
            <el-table-column prop="batchNumber" label="批次号" sortable :show-overflow-tooltip="true" width="100">
            </el-table-column>
            <el-table-column label="信用代码" :show-overflow-tooltip="true" prop="creditCode">
            </el-table-column>
            <el-table-column prop="orderName" label="订单昵称" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="orderCreateTime" label="订单创建时间" :show-overflow-tooltip="true" sortable width="130">
            </el-table-column>
            <el-table-column prop="orderNumber" label="订单编号" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="orderLoadingPoint" label="订单装车点" :show-overflow-tooltip="true" width="120">
            </el-table-column>
            <el-table-column prop="orderUnloadingPoint" label="订单卸货点" :show-overflow-tooltip="true" width="120">
            </el-table-column>
            <el-table-column prop="orderStartTime" label="订单开始时间" :show-overflow-tooltip="true" sortable width="130">
            </el-table-column>
            <el-table-column prop="orderEndTime" label="订单结束时间" :show-overflow-tooltip="true" sortable width="130">
            </el-table-column>
            <el-table-column prop="billNumber" label="运单编号" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="transTime" label="运单创建时间" :show-overflow-tooltip="true" sortable width="130">
            </el-table-column>
            <el-table-column prop="wbLoadingAddress" label="运单发货地址" :show-overflow-tooltip="true" width="130">
            </el-table-column>
            <el-table-column prop="wbUnloadingAddress" label="运单卸货地址" :show-overflow-tooltip="true" width="130">
            </el-table-column>
            <el-table-column prop="carNumber" label="接单车辆" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="subjectMatter" label="标的物" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="freightUnit" label="货运单价" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column label="货运金额" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.freightAmount">{{scope.row.freightAmount|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="原发吨数" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.originalTonnage">{{scope.row.originalTonnage|currency}} 吨</span>
                <span v-else>0.00 吨</span>
              </template>
            </el-table-column>
            <el-table-column label="实收吨数" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.tonnagePaid">{{scope.row.tonnagePaid|currency}} 吨</span>
                <span v-else>0.00 吨</span>
              </template>
            </el-table-column>
            <el-table-column label="应付金额" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.payMoney">{{scope.row.payMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="承运平台" :show-overflow-tooltip="true" width="120" prop="carriagePlatform">
              <template slot-scope="scope">
                <span v-if="scope.row.carriagePlatform">{{scope.row.carriagePlatform|logisticsplatform}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column prop="loadingProLnglog" label="发货地经纬度" :show-overflow-tooltip="true" width="130">
            </el-table-column>
            <el-table-column prop="unloadingProLnglog" label="卸货地经纬度" :show-overflow-tooltip="true" width="130">
            </el-table-column>
            <el-table-column label="开票状态" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.invoicingStatus == 1">待开票</span>
                <span v-if="scope.row.invoicingStatus == 2">已开票</span>
                <span v-if="scope.row.invoicingStatus == 3">发票作废</span>
                <span v-if="scope.row.invoicingStatus == 4">发票冲红</span>
              </template>
            </el-table-column>
            <el-table-column label="可开票金额" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.invoicingAmount">{{scope.row.invoicingAmount|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="可开票吨数" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.invoicingTonnage">{{scope.row.invoicingTonnage|currency}} 吨</span>
                <span v-else>0.00 吨</span>
              </template>
            </el-table-column>
            <el-table-column label="是否已清洗" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.cleanFlag == false">待转换</span>
                <span v-if="scope.row.cleanFlag == true">已转换</span>
              </template>
            </el-table-column>
            <el-table-column prop="createTime" label="创建时间" :show-overflow-tooltip="true">
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据!</p>
              </div>
          </div>
        </el-table>
      </div>
      <!-- 分页 -->
      <div class="pur-pages">
        <el-pagination
        background
        :page-sizes="[50, 100, 200, 300, 500]"
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
            Extended:[
                {
                    id:null,
                    name:'全部'
                },
                {
                    id:false,
                    name:'待转换'
                },
                {
                    id:true,
                    name:'已转换'
                }
            ],            // 网络货运平台数据
            purSearch: "",          // 搜索  
            carriageName:null,        // 数据状态
            dataAva:false,
            loading:false,
            conversionloading:false,
            current:1,              // 当前页数
            size:200,               // 显示条数
            total:0,                // 总条数
            delmul:[],              // 批量删除
            show:false,
            purTitle: "",           // 标题
            multipleSelection: [],  // 选中数据
            tableData: [],          // 表格数据
        };
    },
    created() {},
    mounted(){
      this.purTitle = this.$route.meta.title;
      this.logo = sessionStorage.getItem('companyName') // 企业名称
      this.getpurcon();
      
    },
    methods: {
      // 表格数据
      getpurcon(){
          this.loading = true;
          // 传参
          var listByPageData = {
          searchStr:this.purSearch,
          cleanFlag:this.carriageName,
          page:{
              current:this.current,
              size:this.size
          }
          }
          // 渲染表格
          this.api.logistihird.listByPage(listByPageData)
          .then(res=>{
            if(res.data.code == 200){
                this.loading = false; 
                if(res.data.data.records != ""){
                  this.tableData = res.data.data.records;
                  this.total = res.data.data.total;
                }else{
                  this.dataAva = true;
                }
            }
          })
      },
      // 搜索功能
      search() {
          this.tableData = [];
          this.getpurcon();
          this.current = 1;
      },
      // 数据状态搜索功能
      prologistics(){
          this.tableData = [];
          this.getpurcon();
          this.current = 1;
      },
      // 清洗数据
      openPur() {
        this.conversionloading = true; 
        this.api.logistihird.clearData()
        .then(res=>{
          if(res.data.code == 200){
            this.conversionloading = false; 
            this.$message({
              message: res.data.data,
              type: 'success'
            });
          }
        })
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
    .el-select{
      width: 11%;
      margin-right: 10px;
      float: right;
    }
    .el-input {
      width: 10%;
      float: right;
      margin-right: 10px;
    }
    .el-button {
      margin-right: 10px;
      float: right;
    }
    .el-date-picker{
      width: 10%;
      float: right;
      margin-right: 10px;
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
      left:0;
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
  // 暂无数据样式结束
  // 分页
  .pur-pages {
    width: 100%;
    height: 6%;
    padding: 6px 0 0 12px;
    box-sizing:border-box;
    span{
      font-size: 14px;
      line-height:33px;
      color: #333;
    }
    .el-pagination {
      float: right;
    }
  }
}
// 内容结束
</style>