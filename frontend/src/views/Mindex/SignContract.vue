<template>
  <!-- 签署合同页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="4"
        ><span class="pur-size">{{ purTitle }} (剩余签章 {{signSurplusCount}} 份)</span>
      </el-col>
      <el-col :span="20" class="pur-right">
        <!-- <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button> -->
        <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">发起签署</el-button>
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
      </el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav"
      v-loading="loading"
      element-loading-text = "加载中...">
      <div class="pur-table">
        <!-- 表格 -->    
        <el-table
        style="width: 100%"
        height="100%"
        border
        :data="tableData"
        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
        @selection-change="purChange">
          <el-table-column type="selection" fixed="left" width="40"></el-table-column>
          <el-table-column prop="id" v-if="show">
          </el-table-column>  
          <!-- <el-table-column
            label="操作" width="100" fixed="left">
            <template slot-scope="scope">
              <el-button @click="Clicksee(scope.row)" size="mini" style="margin:0 0 5% 10px;">查看</el-button>
              <el-button size="mini" @click="EditPurs(scope.row)" v-if="$route.meta.edit">编辑</el-button>
            </template>
          </el-table-column> -->
          <el-table-column label="签署合同预览" :show-overflow-tooltip="true" >
            <template slot-scope="scope">
              <span v-if="scope.row.flowStatus == -1">待确认的</span>
              <span v-else @click="confselSign(scope.row)" style="color:#CF1724;cursor: pointer;">签署预览</span>
            </template>
          </el-table-column>
          <el-table-column prop="topicName" label="任务主题" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="contractName" label="合同名称" :show-overflow-tooltip="true" >
          </el-table-column>
          <el-table-column prop="launchType" label="发起方类型" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.launchType == '1'">个人</span>
              <span v-if="scope.row.launchType == '2'">企业</span>
            </template>
          </el-table-column>
          <el-table-column prop="etcName" label="发起方企业名称" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.etcName">{{scope.row.etcName}}</span>
              <span v-else> -- </span>
            </template>
          </el-table-column>
          <el-table-column prop="participants" label="发起方姓名" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="phone" label="发起方手机号" :show-overflow-tooltip="true">
          </el-table-column>
          <!-- <el-table-column prop="launchType" label="签署方类型" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.launchType == '1'">个人</span>
              <span v-if="scope.row.launchType == '2'">企业</span>
            </template>
          </el-table-column>
          <el-table-column prop="etcName" label="签署方企业名称" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.etcName">{{scope.row.etcName}}</span>
              <span v-else> -- </span>
            </template>
          </el-table-column>
          <el-table-column prop="participants" label="签署方姓名" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="phone" label="签署方手机号" :show-overflow-tooltip="true">
          </el-table-column> -->
          <el-table-column prop="flowStatus" label="签署状态" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.flowStatus == -1">待确认</span>
              <span v-if="scope.row.flowStatus == 0">草稿</span>
              <span v-if="scope.row.flowStatus == 1">签署中</span>
              <span v-if="scope.row.flowStatus == 2">完成</span>
              <span v-if="scope.row.flowStatus == 3">撤销</span>
              <span v-if="scope.row.flowStatus == 5">过期</span>
              <span v-if="scope.row.flowStatus == 7">拒签</span>
            </template>
          </el-table-column>
          <el-table-column label="签署方信息" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span @click="seeone(scope.row)" style="color:#CF1724;cursor: pointer;">查看</span>
            </template>
          </el-table-column>
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以发起签署试试</p>
                <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">发 起 签 署</el-button>
              </div>
          </div>
        </el-table>
      </div>
      <!-- 分页 -->
      <div class="pur-pages">
        <el-pagination
        background
        :page-sizes="[10, 30, 50, 100, 200]"
        :pager-count="5"
        :page-size.sync="size"
        :current-page.sync="current"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @current-change='getpurcon'
        @size-change='getpurcon'>
        </el-pagination>
      </div>
      <!-- 文件预览弹窗 -->
      <el-drawer
        title="合同预览"
        :visible.sync="ImgVisible"
        size="35%"
        :before-close="ImgClose">
            <div class="imgBox" v-loading="urlloading" element-loading-text="文件加载中">
              <iframe :src="srcurl" frameborder="0" width="675px" height="99.5%"></iframe>
            </div>  
            <div class="imgbuttons" v-if="oneData.flowStatus == 2">
              <!-- <el-button type="primary" icon="el-icon-printer" size="mini" @click="downPrint">打印</el-button> -->
              <el-button type="primary" icon="el-icon-download" style="float:right;" size="mini" @click="openfile">下载</el-button>
            </div>
      </el-drawer>
      <!-- 签署方信息 -->
      <el-dialog title="签署方信息" :visible.sync="dialogTable">
        <el-table :data="gridData" border>
          <el-table-column label="签署方类型" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.signRole == '0'">个人</span>
              <span v-if="scope.row.signRole == '2'">企业</span>
            </template>
          </el-table-column>
          <el-table-column label="签署方企业名称" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.signEtcName">{{scope.row.signEtcName}}</span>
              <span v-else> -- </span>
            </template>
          </el-table-column>
          <el-table-column prop="signName" label="签署方姓名" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="signPhone" label="签署方手机号" :show-overflow-tooltip="true">
          </el-table-column>
        </el-table>
      </el-dialog>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
  data() {
    return {
      signSurplusCount:null,  // 签章份数
      current:1,              // 当前页数
      size:50,                // 显示条数
      total:0,                // 总条数
      delmul:[],              // 批量删除
      show:false,
      purTitle: "",           // 模块标题
      purSearch: "",          // 搜索
      Extended:[
        {
          dkey:'我发起的',
          dvalue:1
        },
        {
          dkey:'我参与的',
          dvalue:2
        },
        {
          dkey:'待确认的',
          dvalue:3
        },
        {
          dkey:'已完成的',
          dvalue:4
        },
      ],            // 查询类型
      carriageName:1,
      tableData: [],          // 表格数据
      dataAva:false,
      loading:false,
      multipleSelection: [],  // 选中数据
      // 文件预览
      urlloading:true,       
      ImgVisible:false,
      srcurl:'',
      // openfileShow:false,
      oneData:'',   // 预览单条信息
      // 签署方信息
      dialogTable:false,
      gridData:[],      // 签署方表格信息
    };
  },
  created() {},
  mounted(){
    this.purTitle = this.$route.meta.title;
    this.api.enterData.byidData({id:sessionStorage.getItem('companyId')})
    .then(res=>{
      if(res.data.code == 200){
        this.signSurplusCount = res.data.data.signFreeCount*1 + res.data.data.signChargeCount*1;
      }
    })
    this.logo = sessionStorage.getItem('companyName') // 企业名称
    this.getpurcon();
  },
  methods: {
    // 查看签署方信息
    seeone(data){
      this.dialogTable = true
      this.api.eqbSingers.all({sys_flow_id:data.id})
      .then(res=>{
        if(res.data.code == 200){
          this.gridData = res.data.data;
          //  console.log(res.data.data)
        }
       
      })
    },
    // 查询签署人状态(合同预览)
    confselSign(data){
      // if(data.flowStatus == 2){
      //   this.openfileShow = true;
      // }
      // console.log(data)
      this.oneData = data;
      this.ImgVisible = true;
      if(this.carriageName != 3){
          this.api.eqbflows.selSignUrl({eqbFlowId:data.eqbFlowId})
          .then(res=>{
              if(res.data.code == 200){
                this.urlloading = false;
                this.srcurl = res.data.data;
                  // window.open()
              }
          })
      }
        // this.api.eqbflows.selSignPeople({eqbFlowId:data.eqbFlowId})
        // .then(res=>{
        //     if(res.data.code == 200){
        //         let arr = [];
        //         let day;
        //         for(let key in res.data.data){
        //             switch (res.data.data[key].signStatus) {
        //                 case 0:
        //                     day = "待签";
        //                     break;
        //                 case 1:
        //                     day = "未签";
        //                     break;
        //                 case 2:
        //                     day = "已签";
        //                     break;
        //                 case 3:
        //                     day = "待审批";
        //                     break;
        //                 case 4:
        //                     day = "拒签";
        //                     break;
        //             }
        //             arr.push('<div>' +res.data.data[key].signerName + ' : ' + day + '</div>')
        //         }
        //         let a = '<strong>' + arr.toString().replace(",","") + '</strong>';
        //         this.$notify({
        //             title: data.topicName,
        //             dangerouslyUseHTMLString: true,
        //             message: a,
        //             duration: 0
        //         });
        //     }
        // })
    },
    // 关闭签署合同预览弹窗
    ImgClose() {
        this.ImgVisible = false;
        this.srcurl = '';
        // this.openfileShow = false;
    },
    // 下载签署文件
    openfile(){
      this.api.eqbflows.vDownFile({eqbFlowId:this.oneData.eqbFlowId})
      .then(res=>{
        if(res.data.code == 200){
          window.open(res.data.data)
        }
      })
      // let url = this.http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&eqbFlowId=' + this.oneData.eqbFlowId;
      // window.open(url)
    },
    // 表格数据
    getpurcon(){
      this.loading = true;
      // 传参
      var listByPageData = {
        searchStr:this.purSearch,
        searchType:this.carriageName,
        page:{
          current:this.current,
          size:this.size
        }
      }
      // 渲染表格
      this.api.eqbflows.all(listByPageData)
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
    // 查询类型搜索功能
    prologistics(){
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },
    // 表格选中数据
    purChange(val) {
      this.multipleSelection = val;
      this.delmul = [];
      for(let i=0;i<this.multipleSelection.length;i++){
        this.delmul.push(this.multipleSelection[i].id) 
      }
    },
    //新建采购合同
    addPurs(){
        this.$router.push({name:"AddSignContract"})
    },
    // 编辑采购合同
    EditPurs(row){
    //   this.api.cgContract.one({id:row.id})
    //   .then(res=>{
    //     if(res.data.code == 200){
    //       let purdata = res.data.data;
    //       purdata.id = row.id;
    //       sessionStorage.setItem('queryAll',JSON.stringify(purdata))
    //       this.$router.push({name:"EditPurchaseContract"})
    //     }
    //   })
    },
 
    // 查看
    purSee(row) {
      this.drawer = true;
      this.activeName = 'first';
      this.seeId = row.id;
      if(row.orderId == null){
        this.orderId = 0;
      }else{
        this.orderId = row.orderId;
      }
      this.allone(row.id)
    },
    // 查看基础信息函数
    allone(id){
      this.api.cgContract.one({id:id})
      .then(res=>{
        if(res.data.code == 200){
          this.seeData = res.data.data;
        }
      })
      // 文件
      this.api.file.one({relationId:id})
      .then(res=>{
        this.filelist = res.data.data;
      })
    },
    // 删除功能
    openPur() {
        if(this.multipleSelection == ""){
          this.$message({
              type: 'warning',
              message: '请选择要删除的数据!'
          });
        }else{
          this.$confirm('此合同关联相关其他信息,删除后将不再关联,请谨慎操作 是否继续?', '提示', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning'
          }).then(() => {
            this.api.cgContract.del(this.delmul)
            .then(res=>{
              if(res.data.code == 200){
                this.$message({
                  type: 'success',
                  message: '删除成功!'
                });
                this.tableData = [];
                this.getpurcon();
                this.current = 1;
                if(this.tableData == []){
                  this.dataAva = true;
                }
              } 
            })
          }).catch(() => {
              this.$message({
                  type: 'info',
                  message: '已取消删除'
              });          
          });
        }
    },
  },
};
</script>
<style lang="less" scoped>
/deep/ .el-drawer__header {
  span{
    font-size: 18px;
    color: #333;
  }
    
  }
// 文件图片查看样式
.imgBox{
    width: 100%;
    height:100%;
    overflow: hidden;
    overflow-y: scroll;
    img{
        width: 100%;
    }
}
.imgbuttons{
  width:10%;
  height: 40px;
  position: absolute;
  top: 0;
  right:55%;
  bottom:92%;
  left: 0;
  margin: auto;
  z-index: 10000;
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
      width: 120px;
      float: right;
      margin-right:16px;
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
  position: relative;
  .pur-table{
    width:100%;
    height: 94%;
  } 
  // 暂无数据样式
  .dataAvailable{
    width: 100%;
    height: 245px;
    position: relative;
    margin-top:185px;
    margin-bottom:279px;
    .dataAva{
      position: absolute;
      top:75%;
      left:50%;
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
  // 分页样式
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