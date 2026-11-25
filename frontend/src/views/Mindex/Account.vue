<template>
  <!-- 账户维护页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col
      >
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button>
        <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新建</el-button>
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容" 
          prefix-icon="el-icon-search"
          v-model="purSearch"
        >
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
          style="width: 100%;"
          @selection-change="purChange">
            <el-table-column type="selection"></el-table-column>
            <el-table-column prop="id" v-if="show">
            </el-table-column>
            <el-table-column v-if="$route.meta.edit"
              label="操作"
              width="100" fixed="left">
              <template slot-scope="scope">
                <el-button @click="EditPurs(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
              </template>
            </el-table-column>
            <el-table-column label="账户类型" :show-overflow-tooltip="true">
                <template slot-scope="scope">
                    <span v-if="scope.row.ownerType == 1">公司账户</span>
                    <span v-if="scope.row.ownerType == 2">对方账户</span>
                </template>
            </el-table-column>
            <el-table-column label="账户名称" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                  <span v-if="scope.row.accountName">{{scope.row.accountName}}</span>
                  <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column label="账户号码" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                  <span v-if="scope.row.accountNumber">{{scope.row.accountNumber|cardCode}}</span>
                  <span v-else>---</span>
              </template>
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
                <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新建</el-button>
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
          dataAva:false,
          loading:false,
          current:1,              // 当前页数
          size:10,                // 显示条数
          total:0,                // 总条数
          delmul:[],              // 批量删除
          show:false,
          purTitle: "",           // 标题
          purSearch: "",          // 搜索        
          multipleSelection: [],  // 选中数据
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
          searchStr:this.purSearch,
          page:{
            current:this.current,
            size:this.size
          }
        }
        // 渲染表格
        this.api.Account.all(listByPageData)
        .then(res=>{
          this.loading = false;
          if(res.data.code == 200){
            if(res.data.data.records != ""){
              this.tableData = res.data.data.records
              this.total = res.data.data.total
            }else{
              this.dataAva = true;
            }
          }
        })
      },
      // 搜索功能
      search() {
        this.tableData = [];
        this.current = 1;
        this.getpurcon();
      },
      // 表格选中数据
      purChange(val) {
        this.multipleSelection = val;
        this.delmul = [];
        for(let i=0;i<this.multipleSelection.length;i++){
          this.delmul.push(this.multipleSelection[i].id) 
        }
      },
      // 新建账户信息
      addPurs(){
          this.$router.push({name:"AddAccount"})
      },
      // 编辑运单信息
      EditPurs(row){
        this.api.Account.one({id:row.id})
        .then(res=>{
          if(res.data.code == 200){
            let purdata = res.data.data;
            purdata.id = row.id;
            sessionStorage.setItem('queryAll',JSON.stringify(purdata));
            this.$router.push({name:"EditAccount"})
          }
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
            this.$confirm('此操作将永久删除该文件, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
              this.api.Account.del(this.delmul)
              .then(res=>{
                if(res.data.code == 200){
                  this.$message({
                    type: 'success',
                    message: '删除成功!'
                  });
                  this.tableData = [];
                  this.current = 1;
                  this.getpurcon();
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