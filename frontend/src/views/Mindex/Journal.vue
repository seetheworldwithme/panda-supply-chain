<template>
  <!-- 日志明细页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch"
        >
        </el-input>
        <el-date-picker
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
        </el-date-picker>
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
            <el-table-column prop="userName" label="用户名称" :show-overflow-tooltip="true" width="120">
            </el-table-column>
            <el-table-column prop="createTime" label="创建时间" :show-overflow-tooltip="true" width="160">
            </el-table-column>
            <el-table-column prop="operationType" label="操作类型" :show-overflow-tooltip="true" width="120">
                <template slot-scope="scope">
                <span v-if="scope.row.operationType">{{scope.row.operationType|operationType}}</span>
                </template>
            </el-table-column>
            <el-table-column prop="operation" label="操作内容" :show-overflow-tooltip="true">
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无操作日志数据！</p>
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
            purTitle:'',
            purSearch: "",          // 搜索  
            dataAva:false,
            loading:false,
            current:1,              // 当前页数
            size:50,                // 显示条数
            total:0,                // 总条数
            tableData: [],          // 表格数据
            TimeArrays:[],          // 时间搜索
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
            startTime:this.TimeArrays != null?this.TimeArrays[0]:"",
            endTime:this.TimeArrays != null?this.TimeArrays[1]:"",
            page:{
                current:this.current,
                size:this.size
            }
            }
            // 渲染表格
            this.api.commonlog.all(listByPageData)
            .then(res=>{
              if(res.data.code == 200){
                  this.loading = false; 
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
        // 时间搜索
        Times() {
            this.tableData = [];
            this.current = 1;
            this.getpurcon();
        },
    },
};
</script>
<style lang="less" scoped>
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
    .el-input {
      width: 16%;
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
      left:25%;
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