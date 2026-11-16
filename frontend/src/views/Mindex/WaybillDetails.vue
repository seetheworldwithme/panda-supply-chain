<template>
  <!-- 运单详情页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur">返回订单列表</el-button>
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
            <!-- <el-table-column
              label="操作"
              width="90" fixed="left">
              <template slot-scope="scope">
                <el-button @click="Clicksee(scope.row)" size="mini">查 看</el-button>
              </template>
            </el-table-column> -->
            <el-table-column prop="billNumber" label="运单编号" sortable :show-overflow-tooltip="true" width="200">
            </el-table-column>
            <el-table-column label="时间" :show-overflow-tooltip="true" width="180" prop="transTime">
              <!-- <template slot-scope="scope">
                <span v-if="scope.row.transTime">{{scope.row.transTime|Time}}</span>
              </template> -->
            </el-table-column>
            <el-table-column prop="carNumber" label="车号" :show-overflow-tooltip="true" width="100">
            </el-table-column>
            <el-table-column prop="loadingAddress" label="装货地" :show-overflow-tooltip="true" width="260">
            </el-table-column>
            <el-table-column prop="loadingLng" label="装货地经纬度" :show-overflow-tooltip="true" width="300">
            </el-table-column>
            <el-table-column prop="unloadingPoint" label="卸货地" :show-overflow-tooltip="true" width="260">
            </el-table-column>
            <el-table-column prop="unloadingLng" label="卸货地经纬度" :show-overflow-tooltip="true" width="300">
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
            <el-table-column label="开票状态" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.invoicingStatus == 1">待开票</span>
                <span v-if="scope.row.invoicingStatus == 2">已开票</span>
                <span v-if="scope.row.invoicingStatus == 3">发票作废</span>
                <span v-if="scope.row.invoicingStatus == 4">发票冲红</span>
              </template>
            </el-table-column>
            <el-table-column prop="subjectMatter" label="标的物" :show-overflow-tooltip="true" width="100">
            </el-table-column>
            <el-table-column prop="freightUnit" label="运费单价" :show-overflow-tooltip="true" width="100">
            </el-table-column>
            <el-table-column label="运费金额" :show-overflow-tooltip="true" width="120" prop="freightAmount">
              <template slot-scope="scope">
                <span v-if="scope.row.freightAmount">{{scope.row.freightAmount|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="原发吨数" :show-overflow-tooltip="true" prop="originalTonnage">
              <template slot-scope="scope">
                <span v-if="scope.row.originalTonnage">{{scope.row.originalTonnage}} 吨</span>
              </template>
            </el-table-column>
            <el-table-column label="实收吨数" :show-overflow-tooltip="true" prop="tonnagePaid">
              <template slot-scope="scope">
                <span v-if="scope.row.tonnagePaid">{{scope.row.tonnagePaid}} 吨</span>
              </template>
            </el-table-column>
            <el-table-column label="应付金额" :show-overflow-tooltip="true" width="120" prop="payMoney">
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
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
                <!-- <el-button size="small" @click="Importdailog" v-if="$route.meta.import">导入</el-button> -->
              </div>
          </div>
        </el-table>
      </div>
      <!-- 抽屉 查看功能 -->
      <el-drawer :title="logo" :visible.sync="drawer" size="700px">
        <div class="pur-tiao"></div>
        <div class="pur-drawerNav">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="基本信息" name="first">
              <el-row>
                <el-col :span="24" v-if="seeData">
                    <div class="see-top">基础信息</div>
                    <ul class="see-lable">
                        <li>
                            <label>系统编码：</label>
                            <span>{{seeData.systemCode}}</span>
                        </li>
                        <li>
                            <label>运单编号：</label>
                            <span v-if="seeData.billNumber">{{seeData.billNumber}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>时间：</label>
                            <span>{{seeData.transTime}}</span>
                        </li>
                        <li>
                            <label>车号：</label>
                            <span>{{seeData.carNumber}}</span>
                        </li>
                        <li>
                            <label>装货地：</label>
                            <span>{{seeData.loadingAddress}}</span>
                        </li>
                        <li>
                            <label>卸货地：</label>
                            <span>{{seeData.unloadingPoint}}</span>
                        </li>
                        <li>
                            <label>标的物：</label>
                            <span>{{seeData.subjectMatter}}</span>
                        </li>
                        <li>
                            <label>运费单价：</label>
                            <span v-if="seeData.freightUnit">{{seeData.freightUnit|currency}} 元</span>
                            <span v-if="seeData.freightUnit == 0 || seeData.freightUnit==''">---</span>
                        </li>
                        <li>
                            <label>运费金额：</label>
                            <span v-if="seeData.freightAmount">{{seeData.freightAmount|currency}} 元</span>
                            <span v-if="seeData.freightAmount == 0 || seeData.freightAmount==''">---</span>
                        </li>
                        <li>
                            <label>原发吨数：</label>
                            <span>{{seeData.originalTonnage}} 吨</span>
                        </li>
                        <li>
                            <label>实收吨数：</label>
                            <span>{{seeData.tonnagePaid}} 吨</span>
                        </li>
                        <li>
                            <label>应付金额：</label>
                            <span v-if="seeData.payMoney">{{seeData.payMoney|currency}} 元</span>
                            <span v-if="seeData.payMoney == 0 || seeData.payMoney==''">---</span>
                        </li>
                        <li>
                            <label>承运平台：</label>
                            <span v-if="seeData.carriagePlatform">{{seeData.carriagePlatform|logisticsplatform}}</span>
                        </li>
                    </ul>
                </el-col>
              </el-row>
            </el-tab-pane>
            <el-tab-pane label="动态" name="second">
                <el-row>
                    <el-col :span="24">
                        <div class="see-top" style="margin-bottom:24px;">动态信息</div>
                        <el-timeline>
                            <el-timeline-item v-for="v,i in timeData" :key="i" :timestamp="v.createtime" placement="top">
                                <el-card>
                                    <h4>{{v.operation}}</h4>
                                </el-card>
                            </el-timeline-item>
                        </el-timeline>
                    </el-col>
              </el-row>
            </el-tab-pane>
            <!-- 团队模块 V.0.0.2 上线  -->
            <!-- <el-tab-pane label="团队" name="third">团队</el-tab-pane> -->
            <el-tab-pane label="展示图" name="chart">
              <el-row>
                <el-col :span="24">
                  <div class="chart-title">单笔业务流程图</div>
                  <div class="chart-leftImg" @click="zstJump">
                    <img :src="zstimg" alt="">
                    <div class="htl">合同 ( {{searsData.contractNum}} ) </div>
                    <div class="js">结算 ( {{searsData.settNum}} ) </div>
                    <div class="wl">订单 ( {{searsData.orderNum}} ) </div>
                    <div class="zj">资金 ( {{searsData.detNum}} ) </div>
                    <div class="pjl">票据流 ( {{searsData.billNum}} ) </div>
                  </div>
                  <div class="chart-rightImg">
                    <img :src="slhyimg" alt="">
                  </div>
                  <div class="chart-title" style="margin:80px 0 24px 0;">单笔业务数据统计</div>
                  <ul class="chart-data">
                    <li>
                      <div>
                        <i class="el-icon-alihetong1" style="color: #2272EA;"></i>
                        <span>合同</span>
                      </div>
                      <div>
                        <span v-if="searsData.contractMoney"> {{searsData.contractMoney|currency}} 元 </span>
                        <span v-else> 0.00 元</span>
                      </div>
                      <div>
                        <span> {{searsData.contractCount}} 吨 </span>
                      </div>
                    </li>
                    <li>
                      <div>
                        <i class="el-icon-aliwuliu2" style="color: #3FBF9C;"></i>
                        <span>运单</span>
                      </div>
                      <div>
                        <span> {{searsData.loMoney}} </span>
                        
                      </div>
                      <div>
                        <span> {{searsData.loCount}} 吨 </span>
                      </div>
                    </li>
                    <li>
                      <div>
                        <i class="el-icon-alizijin" style="color:#F54D36;"></i>
                        <span>结算</span>
                      </div>
                      <div>
                        <span v-if="searsData.settMoney"> {{searsData.settMoney|currency}} 元 </span>
                        <span v-else> 0.00 元</span>
                      </div>
                      <div>
                        <span> {{searsData.settCount}} 吨 </span>
                      </div>
                    </li>
                    <li>
                      <div>
                        <i class="el-icon-aliticket-fill" style="color: #8092C0;"></i>
                        <span>票据</span>
                      </div>
                      <div>
                        <span v-if="searsData.billMoney"> {{searsData.billMoney|currency}} 元 </span>
                        <span v-else> 0.00 元</span>
                      </div>
                      <div>
                        <span> {{searsData.billCount}} 吨 </span>
                      </div>
                    </li>
                    <li>
                      <div>
                        <i class="el-icon-aliqianbao11" style="color:#F59D36;"></i>
                        <span>流水</span>
                      </div>
                      <div>
                        <span v-if="searsData.detMoney"> {{searsData.detMoney|currency}} 元 </span>
                        <span v-else> 0.00 元</span>
                      </div>
                      <div>
                        <span> {{searsData.detCount}} </span>
                      </div>
                    </li>
                  </ul>
                </el-col>
              </el-row>
            </el-tab-pane>
          </el-tabs>
          <div class="xia">《</div>
        </div>
      </el-drawer>
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
import slhy from '../../assets/slhy.svg';
import zst from '../../assets/zst.svg';
export default {
    data() {
        return {
            // 展示图
            pathIDS:[],             // 四流合同集合
            zstimg:zst,
            slhyimg:slhy,
            seeId:'',               // 查看单条信息的id 
            searsData:{},           // 展示图数据

            Extended:[],            // 网络货运平台数据
            purSearch: "",          // 搜索  
            purSearchtree:'',       // 运单编号搜索
            TimeArrays:[],           // 时间搜索 
            carriageName:'',        // 网络货运平台搜索
            purSearchone: "",       // 发货地搜索
            purSearchtwo: "",       // 卸货地搜索
            tags:[],
            http:'',
            dataAva:false,
            loading:false,
            current:1,              // 当前页数
            size:200,               // 显示条数
            total:0,                // 总条数
            delmul:[],              // 批量删除
            show:false,
            purTitle: "",           // 标题
            multipleSelection: [],  // 选中数据
            drawer: false,          // 查看抽屉
            logo: "",               // 抽屉企业名称
            activeName: "first",    // 查看tab栏
            tableData: [],          // 表格数据
            seeData:{},             // 查看数据
            filelist:[],            // 文件列表
            // 查看...动态
            timeData:[
                // {
                //   createtime:'2021-01-02',
                //   operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
                // }
            ], 
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
            let Did = JSON.parse(sessionStorage.getItem('DetailsId'));
            // 传参
            var listByPageDatatwo = {
                orderId: Did,
                partnerType: 2,
                searchStr: '',
                page: {
                    current: this.current,
                    size: this.size
                }
            }
            this.api.waybill.all(listByPageDatatwo)
            .then(res=>{
                if(res.data.code == 200){
                    this.loading = false; 
                    if(res.data.data.logistics != ""){
                        this.tableData = res.data.data.records;
                        this.total = res.data.data.total;
                    }else{
                        this.dataAva = true;
                    }
                }
            })
        },
        // 查看基础信息函数
        allone(id){
            this.api.waybill.one({id:id})
            .then(res=>{
              if(res.data.code == 200){
                  this.seeData = res.data.data;
              }
            })
        },
        // 查看四流合一图函数
        sears(id){
          this.api.listByPage.sear({loId:id})
          .then(res=>{
            if(res.data.code == 200){
              this.searsData = res.data.data;
              let arr = res.data.data.contractIds;
              this.pathIDS = [];
              for(let key in arr){
                this.pathIDS.push(arr[key].id)
              }
            }
          })
        },
        // 查看
        purSee(index, row) {
            this.drawer = true;
            this.activeName = 'first';
            this.seeId = row.id;
            this.allone(row.id)
        },
        // 查看展示图
        Clicksee(row){
          this.drawer = true;
          this.activeName = 'chart';
          this.seeId = row.id;
          this.sears(row.id)
        },
        // 跳转四流合一
        zstJump(){
          this.$router.push({name:"Relation",params: {paths:'Relation',id:this.pathIDS}});
        },
        changeActive(e) {
            e.currentTarget.className = 'active';
        },
        removeActive(e) {
            e.currentTarget.className = '';
        },
        // 抽屉
        handleClick(tab) {
          // console.log(tab)
          // 判断tab栏选中状态，请求相应的接口
          if (tab.label == "基本信息") {
            this.allone(this.seeId);
          }
          if (tab.label == "动态") {
            // console.log(2);
          }
          if (tab.label == "展示图") {
            this.sears(this.seeId);
          }
        },
        // 返回按钮
        openPur() {
            sessionStorage.removeItem('DetailsId')
            this.$router.push({name:"Waybill"})
        },
    },
};
</script>
<style lang="less" scoped>
// 展示图样式开始
.chart-title{
  width: 160px;
  height: 30px;
  font-size: 18px;
  line-height: 30px;
  padding-left:18px;
  color: #fff;
  background-color:#FB8138;
  border-radius:0 15px 15px 0;
  margin: 24px 0;
}
.chart-leftImg{
  width: 50%;
  display: inline-block;
  padding-left: 20px;
  box-sizing: border-box;
  position: relative;
  cursor: pointer;
  img{
    width: 234px;
  }
  div{
    font-size:14px;
    line-height: 30px;
    color: #fff;
    text-align: center;
  }
  .htl{
    background-color:#2272EA;
    position:absolute;
    top:20%;
    left: 42%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
  .js{
    background-color:#F54D36;
    position: absolute;
    left:15%;
    top:61%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
  .wl{
    background-color:#48CAA6;
    position: absolute;
    left:69%;
    top:61%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
  .zj{
    background-color:#F59D36;
    position: absolute;
    left:12%;
    top:106%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
  .pjl{
    background-color:#8092C0;
    position: absolute;
    left:66%;
    top:106%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
}
.chart-rightImg{
  width: 50%;
  display: inline-block;
  padding-left:20px;
  box-sizing: border-box; 
  img{
    width:90%;
  }
}
.chart-data{
  li{
    overflow: hidden;
    width: 100%;
    padding-left:50px;
    box-sizing: border-box;
    font-weight: 600;
    div{
      float: left;
      width: 20%;
      font-size: 14px;
      color: #333;
      line-height: 30px;
      margin:0 40px 20px 0;
      i{
        font-size: 18px;
        line-height: 30px;
        vertical-align: middle;
        margin-right:6px;
      }
      span{
        font-size: 16px;
        line-height: 30px;
        color: #333;
        vertical-align: middle;
      }
    }
  }
}
// 展示图样式结束
// 导入样式开始
.importone{
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding:50px 75px 0;
  box-sizing: border-box;
  .importone-bottom{
    float: right;
    margin-right: -75px;
    margin-top: 52px;
  }
  a{
    margin-left: 16%;
    font-size: 14px;
    color:#CF1724;
    text-decoration:none;
  }
  .importone-nav{
    padding-top:54px;
    box-sizing: border-box;
    span{
      font-size: 16px;
      line-height: 38px;
      color: #666666;
    }
    p{
      font-size: 14px;
      line-height: 30px;
      color: #666;
    }
  }
}
.importtwo{
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding:50px 75px 0;
  box-sizing: border-box;
  .importtwo-Box{
    width: 130px;
    margin: auto;
    text-align: center;
    span{
      line-height: 50px;
      color:#CF1724;
    }
  }
}
.importstree{
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding:50px 75px 0;
  box-sizing: border-box;
  i{
    font-size: 18px;
    color: #CF1724;
  }
  span{
    font-size: 16px;
    color: #999;
    margin-left:5px;
    line-height: 32px;
  }
  p{
    font-size: 14px;
    color: #666;
    margin-left:30px;
  }
  .stree-buton{
    margin-left:24px;
  }
  .stree-ok{
    margin-top: 109px;
    .el-button{
      float: right;
      margin-right: -75px;
    }
  }
}
.importfour{
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding:50px 75px 0;
  box-sizing: border-box;
  .importfour-Box{
    width: auto;
    margin: auto;
    text-align: center;
    span{
      display: block;
      line-height: 50px;
      color:#999;
    }
    .err-box{
        width: 100%;
        .el-alert{
            float: left;
            width: 19%;
            margin: 1% 1% 1% 0;   
        }
    }
  }
.error-size{
    width: 100%;
    height: 200px;
    border: 1px solid #E4E4E4;
    border-radius: 4px;
    padding: 8px;
    box-sizing: border-box;
    margin-top:15px;
    overflow-y: scroll;
    text-align: left ;
    p{
        font-size: 14px;
        line-height: 24px;
        color: #999;
    }
}
.four-ok{
    margin-top: 30px;
    .el-button{
      float: right;
      margin-right: -75px;
    }
  }
}
// 导入样式结束
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
  padding: 24px 0 24px 24px;
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
  padding: 12px 0;
  label {
      display: inline-block;
    width: 18%;
    font-size: 14px;
    color: #666;
  }
  li {
    height: 40px;
    font-size: 14px;
    line-height: 40px;
    color: #333;
    a {
      span {
        color: #333;
        
      }
    }
  }
}
// 文件样式
.see-file{
  padding: 15px 0 ;
  overflow: hidden;
  span {
    font-size: 14px;
    color: #666;
  }
  dd {
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