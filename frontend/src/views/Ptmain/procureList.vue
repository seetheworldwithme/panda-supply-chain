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
              <el-table-column label="煤种身份" :show-overflow-tooltip="true" prop="createTime" width="120">
                <template slot-scope="scope">
                    <span v-if="scope.row.wareIdentity == 1">洗煤厂</span>
                    <span v-else-if="scope.row.wareIdentity == 2">贸易商</span>
                    <span v-else-if="scope.row.wareIdentity == 3">焦化厂</span>
                    <span v-else-if="scope.row.wareIdentity == 4">经纪人</span>
                    <span v-else-if="scope.row.wareIdentity == 5">煤矿</span>
                    <span v-else-if="scope.row.wareIdentity == 6">钢厂</span>
                    <span v-else-if="scope.row.wareIdentity == 7">金融企业</span>
                    <span v-else-if="scope.row.wareIdentity == 8">平台经纪人</span>
                    <span v-else-if="scope.row.wareIdentity == 9">其他</span>
                    <span v-else>--</span>
                </template>
              </el-table-column>
              <el-table-column label="采购分类" :show-overflow-tooltip="true" prop="secName" width="120">
              </el-table-column>
              <el-table-column label="采购品种" :show-overflow-tooltip="true" prop="thName" width="120">
              </el-table-column>
              <el-table-column label="采购数量（吨）" :show-overflow-tooltip="true" prop="purchaseQuantity" width="120">
              </el-table-column>
              <el-table-column label="采购价（元/吨）" :show-overflow-tooltip="true" prop="purchasePrice" width="130">
              </el-table-column>
              <el-table-column label="采购有效期" :show-overflow-tooltip="true" prop="period" width="180">
                <template slot-scope="scope">
                    <span>{{scope.row.purchStartTime}}~{{scope.row.purchEndTime}}</span>
                </template>
              </el-table-column>
              <!-- <el-table-column label="长期采购执行有效期" :show-overflow-tooltip="true" prop="receivingPeople" width="180">
                <template slot-scope="scope">
                    <span v-if="scope.row.purchaseType == 2">{{scope.row.purchStartTime}}~{{scope.row.purchEndTime}}</span>
                    <span v-else>--</span>
                </template>
              </el-table-column> -->
              <el-table-column label="采购企业" :show-overflow-tooltip="true" prop="enterpriseName" width="150">
              </el-table-column>
              <el-table-column label="联系姓名" :show-overflow-tooltip="true" prop="coalRelName" width="100">
              </el-table-column>
              <el-table-column label="联系电话" :show-overflow-tooltip="true" prop="coalRelPhone" width="140">
              </el-table-column>
              <el-table-column label="标的指标" :show-overflow-tooltip="true" width="300">
                <template slot-scope="scope">
					<span v-if="scope.row.maxAshContent!=0&&scope.row.maxAshContent!=null">灰分(Ad)%≤{{scope.row.maxAshContent}},</span>
					<span v-if="scope.row.maxWaterContent!=0&&scope.row.maxWaterContent!=null">水分(Mt)%≤{{scope.row.maxWaterContent}},</span>
					<span v-if="scope.row.maxStd!=0&&scope.row.maxStd!=null">硫(St,d)%≤{{scope.row.maxStd}},</span>
					<span v-if="scope.row.maxVolatiles!=0&&scope.row.maxVolatiles!=null">挥发分(Vdaf)%≤{{scope.row.maxVolatiles}},</span>
					<span v-if="scope.row.minBondIndex!=0&&scope.row.minBondIndex!=null">粘结指数GR.I≥{{scope.row.minBondIndex}},</span>
					<span v-if="scope.row.minY!=0&&scope.row.minY!=null">Ymm≥{{scope.row.minY}},</span>
					<span v-if="scope.row.minCsr!=0&&scope.row.minCsr!=null">热态CSR≥{{scope.row.minCsr}},</span>
					<span v-if="scope.row.minHgi!=0&&scope.row.minHgi!=null">哈氏可磨指数（HGI）≥{{scope.row.minHgi}},</span>
					<span v-if="scope.row.lithofacies!=0&&scope.row.lithofacies!=null">岩相S≤{{scope.row.lithofacies}},</span>
					<span v-if="scope.row.recoveryAnalysis!=0&&scope.row.recoveryAnalysis!=null">回收率分析{{scope.row.recoveryAnalysis}},</span>
					<span v-if="scope.row.calorificValue!=0&&scope.row.calorificValue!=null">热值≥{{scope.row.calorificValue}}</span>
                </template>
              </el-table-column>
              <el-table-column label="交货地" :show-overflow-tooltip="true" prop="deliveryPlace">
              </el-table-column>
              <el-table-column label="运输方式" :show-overflow-tooltip="true" prop="transportModel">
              </el-table-column>
              <el-table-column label="付款方式" :show-overflow-tooltip="true" prop="paymentMethod">
              </el-table-column>
              <el-table-column label="煤质扣罚要求" :show-overflow-tooltip="true" prop="leavelMsg" width="120">
                <template slot-scope="scope">
                    <span v-if="scope.row.mzkfyq > 0" style="color: #E2231A;" @click="imgmap(scope.row.id)">详情</span>
                    <span v-else>--</span>
                </template>
              </el-table-column>
              <el-table-column label="发布时间" :show-overflow-tooltip="true" prop="createTime" width="150">
              </el-table-column>
              <el-table-column label="备注" :show-overflow-tooltip="true" prop="remark">
                <template slot-scope="scope">
                    <span v-if="scope.row.remark">{{ scope.row.remark }}</span>
                    <span v-else>--</span>
                </template>
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
          <el-pagination style="float:right;"
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
        <!-- 煤质扣罚要求信息弹窗 -->
        <el-dialog
          :closeOnClickModal="false"
          title="煤质扣罚要求"
          :visible.sync="paymentShow"
          width="30%"
          :before-close="paymentShowClose">
            <div class="imgmaps">
                <el-carousel indicator-position="none" height="650px">
                    <el-carousel-item v-for="(item,k) in AddData" :key="k">
                        <img style="width:100%;height:100%;" :src="http + '/saas-common/upload/imageShow?fileName=' + item.img + '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'" alt="">
                    </el-carousel-item>
                </el-carousel>
            </div>
        </el-dialog>
    </div>
  </template>
  <script>
  export default {
      data() {
            return {
                http:null,
                loading:false,
                current:1,              // 当前页数
                size:100,               // 显示条数
                total:0,                // 总条数
                show:false,
                purTitle: "",           // 标题
                purSearch: "",          // 搜索      
                tableData: [],          // 表格数据  
                // 煤质扣罚要求信息弹窗
                paymentShow:false,
                AddData:[],
            };
      },
      created() {
        this.http  = this.$store.state.http;                        // IP地址
      },
      mounted(){
          this.purTitle = this.$route.meta.title;
          this.getpurcon();
      },
      methods: {
            // 新建标签取消
            paymentShowClose(){
                this.paymentShow = false;
                this.getpurcon();
                this.AddData = [];
            },
            // 煤质扣罚要求
            imgmap(list){
                this.api.busrodect.one({id:list})
                .then(res=>{
                    if(res.data.code == 200){
                        res.data.data.imgListMap.map(item=>{
                            this.AddData.push({img:item.fileName})
                        })
                    }
                })
                this.paymentShow = true;
            },
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
                    effective: 1,
				    publishType:2
                }
                // 渲染表格
                this.api.busrodect.all(listByPageData)
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
                this.current = 1;
                this.getpurcon();
            },
      },
  };
  </script>
  <style lang="less" scoped>
  // 煤质扣罚要求样式
  .imgmaps{
    width: auto;
    height: 650px;
    // img{
    //     width: 100%;
    //     height: 650px;
    // }
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

    