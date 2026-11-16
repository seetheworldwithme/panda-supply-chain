<template>
  <!-- 对账记录页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <!-- <el-button size="small" @click="addPurs">新建</el-button> -->
        <!-- <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch"
        >
        </el-input> -->
        <!-- <el-date-picker
        size="small"
        @change="Times"
        v-model="TimeArrays"
        type="daterange"
        unlink-panels
        range-separator="~"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="yyyy-MM-dd"
        style="width:20%;float:right;margin-right:16px;">
        </el-date-picker> -->
        <el-date-picker
        @change="Times"
        size="small"
        v-model="TimeArrays"
        type="date"
        value-format="yyyyMMdd"
        placeholder="选择日期">
        </el-date-picker>
      </el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <!-- 
      v-loading="loading"
      element-loading-text = "加载中..." -->
    <div class="pur-nav">
      <!-- 表格 -->
      <div class="pur-table"
      v-loading="loading"
      element-loading-text = "加载中...">
        <el-table
          height="100%"
          border
          :data="tableData"
          :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
          style="width: 100%;">
            <el-table-column label="订单号" prop="orderNbr"  :show-overflow-tooltip="true" width="200">
            </el-table-column>
            <el-table-column prop="orderSubNbr" label="订单交易序号" :show-overflow-tooltip="true" width="120">
            </el-table-column>
            <el-table-column label="商户订单号" prop="mchOrderNbr"  :show-overflow-tooltip="true" width="180">
            </el-table-column>
            <el-table-column prop="mchTransNbr" label="商户交易流水号" :show-overflow-tooltip="true" width="180">
            </el-table-column>
            <el-table-column label="退款标志" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.bckFlag == 'Y'">已退款</span>
                <span v-if="scope.row.bckFlag == 'I'">处理中</span>
              </template>
            </el-table-column>
            <el-table-column label="订单类型" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.orderType == 'PAY'">消费类</span>
                <span v-if="scope.row.orderType == 'BCK'">退货类</span>
              </template>
            </el-table-column>
            <el-table-column label="订单子类" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.orderSubtype == 'DRT'">直接支付</span>
                <span v-if="scope.row.orderSubtype == 'GRT'">担保支付</span>
                <span v-if="scope.row.orderSubtype == 'BCK'">退货退款</span>
              </template>
            </el-table-column>
            <el-table-column label="订单金额" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.orderAmt">{{scope.row.orderAmt|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="订单币种" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.orderCcyNbr == '10'">人民币</span>
              </template>
            </el-table-column>
            <el-table-column label="订单状态" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.status == 'STR'">发起</span>
                <span v-if="scope.row.status == 'OPR'">已受理</span>
                <span v-if="scope.row.status == 'AUT'">待付款</span>
                <span v-if="scope.row.status == 'WCF'">待确认</span>
                <span v-if="scope.row.status == 'ING'">处理中（退款）</span>
                <span v-if="scope.row.status == 'DBT'">已扣款（退款）</span>
                <span v-if="scope.row.status == 'FIN'">完成</span>
              </template>
            </el-table-column>
            <el-table-column label="请求结果" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.rtnFlag == 'S'">成功</span>
                <span v-if="scope.row.rtnFlag == 'F'">失败</span>
                <span v-if="scope.row.rtnFlag == 'B'">退票</span>
                <span v-if="scope.row.rtnFlag == 'C'">取消</span>
                <span v-if="scope.row.rtnFlag == 'I'">处理中</span>
                <span v-if="scope.row.rtnFlag == 'W'">可疑</span>
              </template>
            </el-table-column>
            <el-table-column label="订单开始时间" :show-overflow-tooltip="true" width="160">
              <template slot-scope="scope">
                <span v-if="scope.row.orderStrTime">{{scope.row.orderStrTime|orderStrTime}}</span>
              </template>
            </el-table-column>
            <!-- <el-table-column label="订单完成时间" prop="orderEndTime"  :show-overflow-tooltip="true">
            </el-table-column> -->
            <el-table-column label="订单推送时间" :show-overflow-tooltip="true" width="160">
              <template slot-scope="scope">
                <span v-if="scope.row.orderPostTime">{{scope.row.orderPostTime|orderStrTime}}</span>
              </template>
            </el-table-column>
            <el-table-column label="订单无效时间" :show-overflow-tooltip="true" width="160">
              <template slot-scope="scope">
                <span v-if="scope.row.orderInvalidTime">{{scope.row.orderInvalidTime|orderStrTime}}</span>
              </template>
            </el-table-column>
            <!-- <el-table-column label="订单无效标志" prop="orderInvalidFlag"  :show-overflow-tooltip="true">
            </el-table-column> -->
            <el-table-column label="平台编号" prop="pltNbr"  :show-overflow-tooltip="true" width="120" >
            </el-table-column>
            <el-table-column label="商户编号" prop="mchNbr"  :show-overflow-tooltip="true" width="120">
            </el-table-column>
            <el-table-column label="商户简称" prop="mchAlias"  :show-overflow-tooltip="true" width="200">
            </el-table-column>
            <el-table-column label="收款方受理机构标识" prop="payeeIssrId" :show-overflow-tooltip="true" width="150">
            </el-table-column>
            <el-table-column label="收款方账户机构标识" prop="payeeAccIssrId" :show-overflow-tooltip="true" width="150">
            </el-table-column>
            <el-table-column label="收款方账号" prop="payeeAccNbr" :show-overflow-tooltip="true" width="150">
            </el-table-column>
            <el-table-column label="收款方账户名" prop="payeeAccName" :show-overflow-tooltip="true" width="200">
            </el-table-column>
            <el-table-column label="收款方行内外标志" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.payeeBankFlag == 'Y'">行内</span>
                <span v-if="scope.row.payeeBankFlag == 'N'">行外</span>
              </template>
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
              </div>
          </div>
        </el-table>
      </div>
      <!-- 分页 -->
      <!-- <div class="pur-pages">
        <el-pagination
        background
        :page-sizes="[100]"
        :page-size.sync="size"
        :pager-count="5"
        :current-page.sync="current"
        layout="sizes, prev, pager, next, jumper"
        @current-change='getpurcon'
        @size-change='getpurcon'>
        </el-pagination>
      </div> -->
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
    data() {
        return {
            purSearch: "",          // 搜索  
            TimeArrays:'',           // 时间搜索 
            dataAva:false,
            loading:false,
            current:1,              // 当前页数
            size:100,               // 显示条数
            purTitle: "",           // 标题
            tableData: [],          // 表格数据
        };
    },
    created() {},
    mounted(){
        this.purTitle = this.$route.meta.title;
        // 获取当前时间
        let date = new Date();
        let year = date.getFullYear();
        let month = date.getMonth() + 1;
        let day = date.getDate();
        if (month < 10) {
            month = "0" + month;
        }
        if (day < 10) {
            day = "0" + day;
        }
        this.TimeArrays = year + month + day;
        // this.TimeArrays = '20220328';
        this.getpurcon();
    },
    methods: {
        // 表格数据
        getpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
                orderDate:this.TimeArrays,
                pageNum:this.current,
                pageSize:this.size
            }
            // 渲染表格
            this.api.confirmPay.billInfo(listByPageData)
            .then(res=>{
              if(res.data.code == 200){
                  this.loading = false; 
                  if(res.data.data.body.billQryRtnDto != ""){
                    this.tableData = res.data.data.body.billQryRtnDto;
                  }else{
                    this.dataAva = true;
                  }
              }else if(res.data.code == 602){
                this.loading = false; 
                this.dataAva = true;
              }
            })
        },
        // 日期搜索
        Times() {
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
// 查看抽屉样式开始
.pur-tiao {
  width: 100%;
  height: 8px;
  background: #f5f5f5;
}
.pur-drawerNav {
  width: 100%;
  padding:24px 0 24px 24px;
  height: calc(100% - 8px);
  .el-row {
    height: calc(100vh - 200px);
    overflow-y: scroll;
  }
  .xia{
    width: 30px;
    margin: auto;
    font-size: 16px;
    color: #828282;
    line-height: 18px;
    transform: rotate(270deg);
  }
}
.see-top {
  height: 40px;
  font-weight: 500;
  font-size: 14px;
  line-height: 45px;
  border-bottom: 1px dashed #999999;
}
.see-lable {
  padding: 15px 0;
  label {
    display: inline-block;
    width: 18%;
    font-size: 14px;
    color: #666;
  }
  li {
    min-height: 40px;
    font-size: 14px;
    line-height: 40px;
    color: #333;
    label{
      float: left;
    }
    .see-libox{
      width: 76%;
      display: inline-block;
      span{
        color: #666;
      }
      div{
        display: inline-block;
        margin-right: 15px;
      }
    }
  }

}
// 文件样式
.see-file{
  padding: 15px 0 ;
  span {
    font-size: 14px;
    color: #666;
  }
  dd {
    display: flow-root;
    font-size: 14px;
    line-height: 34px;
    color: #333;
    a{
      text-decoration: none;
      color: #333;
    }
    .active{
      color: #EA222E;
    }
  }
}
// 查看抽屉样式结束
</style>