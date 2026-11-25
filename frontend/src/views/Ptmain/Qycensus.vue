<template>
  <!-- 企业统计页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="Signingpur">签约企业</el-button>
        <el-button size="small" @click="entVip">企业会员</el-button>
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
        <el-select v-model="urrria" size="small" @change="prologistics">
          <el-option
          v-for="item in enterpriseSource"
          :key="item.value"
          :label="item.dkey"
          :value="item.dvalue">
          </el-option>
        </el-select>
        <el-select v-model="isMember" size="small" @change="memberChange">
          <el-option
            v-for="item in memberOption"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue"
          >
          </el-option>
        </el-select>
        <div class="selTit">会员：</div>
        <el-select v-model="isSign" size="small" @change="signChange">
          <el-option
            v-for="item in signOption"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue"
          >
          </el-option>
        </el-select>
        <div class="selTit">签约：</div>
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
          style="width: 100%;" @selection-change="purChange">
            <el-table-column type="selection"></el-table-column>
            <el-table-column prop="id" v-if="show"></el-table-column>
            <el-table-column
              label="操作" width="120" fixed="left" v-if="$route.meta.edit">
              <template slot-scope="scope">
                <el-button size="mini" @click="EditPurs(scope.row)" style="margin-left: 10px;">编辑</el-button>
                <el-button size="mini" @click="delPurs(scope.row)" style="margin-left: 10px;margin-top: 8px;">删除</el-button>
                <el-button v-if="scope.row.isSignEtc == true" size="mini" @click="Editisetc(scope.row)" style="margin-top: 8px;">取消签约</el-button>
                <el-button v-if="scope.row.isMemberEtc" size="mini" @click="cancelVip(scope.row)" style="margin-top: 8px;">取消会员</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="isSignEtc" label="签约企业" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.isSignEtc" style="color:#6FCF97;">已签约</span>
                <span v-if="scope.row.isSignEtc == false">未签约</span>
              </template>
            </el-table-column>
            <el-table-column label="企业会员" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.isMemberEtc" style="color:#D12B2B;">会员</span>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column label="是否认证企业" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.isAuthentication" style="color:#6FCF97;">已认证</span>
                <span v-else style="color:#D12B2B;">未认证</span>
              </template>
            </el-table-column>
            <el-table-column prop="enterpriseName" label="企业名称" :show-overflow-tooltip="true" width="220">
            </el-table-column>
            <el-table-column prop="proCount" label="煤种数量" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="scope" label="企业评级" :show-overflow-tooltip="true" width="180">
              <template slot-scope="scope">
                <div style="float: left;">
                  <el-tooltip class="item" effect="dark" placement="left-start">
                  <div slot="content">
                    <span v-if="scope.row.quaCertScore">认证企业资质{{scope.row.quaCertScore}}星,</span>
                    <span v-if="scope.row.basInfoScore">认证基础信息{{scope.row.basInfoScore}}星,</span>
                    <span v-if="scope.row.coalVerifyScore">煤种核实{{scope.row.coalVerifyScore}}星,</span>
                    <span v-if="scope.row.updShareScore">挂单价格更新和发布人分享{{scope.row.updShareScore}}星,</span>
                    <span v-if="scope.row.reportScore">煤种有化验单报告{{scope.row.reportScore}}星,</span>
                    <span v-if="scope.row.qualifiedScore">成交后评星{{scope.row.qualifiedScore}}星,</span>
                  </div>
                  <i class="el-icon-question"></i>
                </el-tooltip>
              </div>
                <el-rate
                  v-model="scope.row.sumScore"
                  disabled
                  show-score
                  text-color="#E2231A"
                  :colors="['#E2231A', '#E2231A', '#E2231A']"
                  score-template="{value}">
                </el-rate>
              </template>
            </el-table-column>
            <el-table-column prop="type" label="企业logo" :show-overflow-tooltip="true">
              <template  slot-scope="scope">
                <img v-if="scope.row.imageUrl" :src="scope.row.imageUrl" alt="" style="width:20px;height: 20px;">
                <img v-else :src="scope.row.logo" alt="" style="width:20px;height: 20px;">
              </template>
            </el-table-column>
            <el-table-column prop="regTime" label="注册时间" :show-overflow-tooltip="true" width="100">
            </el-table-column>
            <el-table-column prop="rateNumber" label="税号" :show-overflow-tooltip="true" width="180">
            </el-table-column>
            <el-table-column prop="legalPerson" label="法人姓名" :show-overflow-tooltip="true" width="90">
            </el-table-column>
            <el-table-column prop="regMoney" label="注册资本" :show-overflow-tooltip="true" width="90">
            </el-table-column>
            <el-table-column prop="recCap" label="实缴资本" :show-overflow-tooltip="true" width="90">
            </el-table-column>
            <el-table-column prop="econKind" label="企业类型" :show-overflow-tooltip="true" width="90">
            </el-table-column>
            <el-table-column prop="updatedDate" label="更新日期" :show-overflow-tooltip="true" width="90">
            </el-table-column>            
            <el-table-column prop="address" label="注册地址" :show-overflow-tooltip="true" width="90">
            </el-table-column>
            <el-table-column prop="scope" label="经营范围" :show-overflow-tooltip="true" width="90">
            </el-table-column>
            <el-table-column prop="proName" label="省" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="cityName" label="市" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="areaName" label="区 / 县" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="countryName" label="乡" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="legIdentity" label="法人身份证" :show-overflow-tooltip="true" width="170">
            </el-table-column>
            <el-table-column prop="shTime" label="上户时间" :show-overflow-tooltip="true" width="100">
            </el-table-column>
            <el-table-column prop="phone" label="联系电话" :show-overflow-tooltip="true" width="120">
            </el-table-column>
            <el-table-column prop="businesser" label="业务员" :show-overflow-tooltip="true" width="80">
            </el-table-column>
            <el-table-column prop="type" label="企业类别" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.type == '1'">一般纳税人</span>
                <span v-if="scope.row.type == '2'">小规模纳税人</span>
              </template>
            </el-table-column>
            <el-table-column prop="isNetwork" label="网络货运" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.isNetwork == true">是</span>
                <span v-if="scope.row.isNetwork == false">否</span>
              </template>
            </el-table-column>
            <el-table-column prop="rotationFlag" label="小程序轮播标志" :show-overflow-tooltip="true" width="130">
              <template slot-scope="scope">
                <span v-if="scope.row.rotationFlag == true">是</span>
                <span v-if="scope.row.rotationFlag == false">否</span>
              </template>
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据 !</p>
                <!-- <el-button size="small" @click="Importdailog">导入</el-button> -->
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
            show:false,
            purSearch: "",          // 搜索  
            http:'',
            dataAva:false,
            loading:false,
            current:1,              // 当前页数
            size:200,                // 显示条数
            total:0,                // 总条数
            purTitle: "",           // 标题
            tableData: [],          // 表格数据
            multipleSelection:[],
            // 下拉筛选
            Extended:[
              {
                  dkey:'已认证',
                  dvalue:true
              },
              {
                  dkey:'未认证',
                  dvalue:false
              }
            ], 
            carriageName:true,

            enterpriseSource:[
              {
                  dkey:'全部',
                  dvalue:0
              },
              {
                  dkey:'熊猫煤仓企业',
                  dvalue:1
              },
              {
                  dkey:'招商企业',
                  dvalue:2
              },
              {
                  dkey:'其他',
                  dvalue:3
              }
            ], 
            urrria:0,
            value:3.7, // 企业评级
            memberOption:[
               {
                  dkey:'全部',
                  dvalue:null
              },
              {
                  dkey:'会员',
                  dvalue:true
              },
              {
                  dkey:'非会员',
                  dvalue:false
              }
            ],
            isMember:null,
            signOption:[
               {
                  dkey:'全部',
                  dvalue:null
              },
              {
                  dkey:'已签约',
                  dvalue:1
              },
              {
                  dkey:'未签约',
                  dvalue:0
              }
            ],
            isSign:null,
        };
    },
    created() {},
    mounted(){
      this.purTitle = this.$route.meta.title;
      this.logo = sessionStorage.getItem('companyName') // 企业名称
      this.http  = this.$store.state.http;                 // IP地址
      this.getpurcon();
    },
    methods: {
      // 下拉查询
      prologistics(){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
      },
      // 表格选中数据
      purChange(val) {
        this.multipleSelection = val;
      },
      // 取消签约
      Editisetc(data){
        data.isSignEtc = false;
        this.api.enterData.updateById({
          id:data.id,
          isSignEtc:data.isSignEtc
        })
          .then(res=>{
            if(res.data.code == 200){
              this.getpurcon();
              this.$message({
                type: 'success',
                message: '取消签约企业成功!'
              });
            }
          })
      },
      // 签约企业
      Signingpur(){
        if(this.multipleSelection.length >0){
          let delmul = [];
          for(let i=0;i<this.multipleSelection.length;i++){
            delmul.push({id:this.multipleSelection[i].id,isSignEtc:true}) 
          }
          this.api.enterData.updBatch(delmul)
          .then(res=>{
            if(res.data.code == 200){
              this.getpurcon();
              this.$message({
                type: 'success',
                message: '签约企业成功!'
              });
            }
          })
          // console.log(delmul)
        }else{
          this.$message({
            type: 'warning',
            message: '请先选中企业在签约!'
          });
        }
      },
      // 企业会员
      entVip(){
        if(this.multipleSelection.length >0){
          let delmul = [];
          for(let i=0;i<this.multipleSelection.length;i++){
            delmul.push({id:this.multipleSelection[i].id,isMemberEtc:true}) 
          }
          this.api.enterData.updBatch(delmul)
          .then(res=>{
            if(res.data.code == 200){
              this.getpurcon();
              this.$message({
                type: 'success',
                message: '设置会员成功!'
              });
            }
          })
          // console.log(delmul)
        }else{
          this.$message({
            type: 'warning',
            message: '请先选中企业在设置会员!'
          });
        }
      },
      // 取消会员
      cancelVip (data){
        data.isMemberEtc = false;
        this.api.enterData.updateById({
          id:data.id,
          isMemberEtc:false
        })
          .then(res=>{
            if(res.data.code == 200){
              this.getpurcon();
              this.$message({
                type: 'success',
                message: '取消会员成功!'
              });
            }
          })
      },
      // 表格数据
      getpurcon(){
          this.loading = true;
          // 传参
          var listByPageData = {
            selEtcType:this.urrria,
            isAuthentication:this.carriageName,
            searchStr:this.purSearch,
            page:{
                current:this.current,
                size:this.size
            },
            isMemberEtc:this.isMember,
            isSignEtc:this.isSign
          }
          // 渲染表格
          this.api.enterData.all(listByPageData)
          .then(res=>{
            if(res.data.code == 200){
                this.loading = false; 
                if(res.data.data.records != ""){
                  res.data.data.records.forEach((item,index)=> {
                    item.logo =  this.http + '/saas-common/upload/imageShow?fileName=' + item.logo + '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
                  })
                  this.tableData = res.data.data.records
                  this.total = res.data.data.total
                }else{
                  this.dataAva = true;
                }
            }
          })
      },
      // 编辑企业统计
      EditPurs(row){
        this.api.enterData.byidData({id:row.id})
        .then(res=>{
          if(res.data.code == 200){
            let purdata = res.data.data;
            purdata.id = row.id;
            // purdata.etcs = [];
            sessionStorage.setItem('queryAll',JSON.stringify(purdata))
            this.$router.push({name:"EditQycensus"})
          }
        })
      },
      // 搜索功能
      search() {
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
      },
      // 删除企业信息
      delPurs(data){
          this.$confirm('此企业信息,删除后将不能恢复,请谨慎操作 是否继续?', '提示', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning'
          }).then(() => {
            this.api.enterData.onedel({id:data.id})
            .then(res=>{
              if(res.data.code == 200){
                this.$message({
                  type: 'success',
                  message: '删除成功!'
                });
                this.tableData = [];
                this.current = 1;
                this.getpurcon();
              }
            })
          }).catch(() => {
              this.$message({
                  type: 'info',
                  message: '已取消删除'
              });          
          });
        },
        memberChange(e){
          this.isMember = e
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        signChange(e){
          this.isSign = e
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        }
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
      width: 20%;
      float: right;
      margin-right: 10px;
    }
    .el-button {
        margin-left: 16px;
        float: right;
      }
      .el-select{
        margin-right: 16px;
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
      top:50%;
      left:15%;
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
.selTit {
  height: 32px;
  float: right;
  line-height: 32px;
  font-size: 14px;
}
// 查看抽屉样式结束
</style>