<template>
    <!-- 关联运单信息页面 -->
    <div>
        <!-- 头部开始 -->
        <el-row class="pur-top">
        <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
        <el-col :span="21" class="pur-right">
            <el-button size="small" @click="AddDatasubmitForm">保存</el-button>
            <el-button size="small" style="margin-left: 24px"  @click="Addopen">取消</el-button>
        </el-col>
        </el-row>
        <!-- 头部结束 -->
        <!-- 内容开始 -->
        <div class="pur-nav"
        v-loading="loading"
        element-loading-text = "加载中...">
            <div class="searchs">
                <el-input
                @keyup.enter.native="search"
                size="small"
                placeholder="请输入内容"
                prefix-icon="el-icon-search"
                v-model="purSearch"
                >
                </el-input>
                <el-input
                @keyup.enter.native="searchtree"
                size="small"
                placeholder="运单编号"
                prefix-icon="el-icon-search"
                v-model="purSearchtree">
                </el-input>
                <el-date-picker
                @change="DataTime"
                size="small"
                v-model="keyupTime"
                type="date"
                value-format="yyyy-MM-dd"
                placeholder="运单创建时间">
                </el-date-picker>
                <el-input
                @keyup.enter.native="searchone"
                size="small"
                placeholder="发货地搜索"
                prefix-icon="el-icon-search"
                v-model="purSearchone"
                >
                </el-input>
                <el-input
                @keyup.enter.native="searchtwo"
                size="small"
                placeholder="卸货地搜索"
                prefix-icon="el-icon-search"
                v-model="purSearchtwo"
                >
                </el-input>
            </div>
            <!-- 表格 -->
            <div class="pur-table">
                <el-table
                    height="100%"
                    border
                    :data="tableData"
                    :row-key="getRowKey"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;"
                    ref="multipleTable"
                    @selection-change="purChange">
                    <el-table-column type="selection" :reserve-selection="true"> </el-table-column>
                    <el-table-column prop="id" v-if="show">
                    </el-table-column>
                    <el-table-column prop="billNumber" label="运单编号" sortable :show-overflow-tooltip="true" width="260">
                    </el-table-column>
                    <el-table-column prop="loadingAddress" label="发货地" :show-overflow-tooltip="true" width="260">
                    </el-table-column>
                    <el-table-column prop="unloadingPoint" label="卸货地" :show-overflow-tooltip="true" width="260">
                    </el-table-column>
                    <el-table-column prop="subjectMatter" label="标的物" :show-overflow-tooltip="true">
                    </el-table-column>
                    <el-table-column prop="freightUnit" label="运费单价" :show-overflow-tooltip="true">
                    </el-table-column>
                    <el-table-column label="运费金额" :show-overflow-tooltip="true" width="140">
                        <template slot-scope="scope">
                        <span v-if="scope.row.freightAmount">{{scope.row.freightAmount|currency}} 元</span>
                        <span v-if="scope.row.freightAmount==0 || scope.row.freightAmount==''">---</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="吨数" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                        <span>{{scope.row.subjectCount}} 吨</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="原发吨数" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                        <span>{{scope.row.originalTonnage}} 吨</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="实收吨数" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                        <span>{{scope.row.tonnagePaid}} 吨</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="transTime" label="交易时间" sortable :show-overflow-tooltip="true" width="180">
                    </el-table-column>
                    <!-- 暂无数据 -->
                    <div slot="empty" class="dataAvailable" v-if="dataAva">
                        <div class="dataAva">
                        <img src="../../assets/zanwushuju.png" alt="">
                        <p>暂无数据，您可以新建试试</p>
                        <!-- <el-button size="small" @click="Importdailog">导入</el-button> -->
                        </div>
                    </div>
                </el-table>
            </div>
            <!-- 分页 -->
            <div class="pur-pages">
                <!-- <span>合计：运费金额 {{newcontractNum|currency}} 元 , 原发吨数 {{yuansumPrice|currency}} 吨 , 实收吨数 {{newsumPrice|currency}} 吨 </span> -->
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
            /* 
            关联运单信息数据
            */
            show:false,
            purTitle: "",                   // 标题
            dataAva:false,
            loading:false,
            purSearch: "",          // 搜索
            purSearchtree:'',       // 运单编号搜索
            keyupTime:'',           // 运单创建时间搜索
            purSearchone: "",       // 发货地搜索
            purSearchtwo: "",       // 卸货地搜索
            current:1,              // 当前页数
            size:200,                // 显示条数
            total:0,                // 总条数
            multipleSelection: [],  // 选中数据
            tableData: [],          // 表格数据
            // logisticsAll:[],        // 选中表格数据
        };
    },
    created() {},
    mounted () {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        this.getpurcon();

        

        // this.onOldAjglId();
    },
    methods: {
        // 回显表格选中数据
        // onOldAjglId(){
        //     console.log(this.logisticsAll,this.tableData)
        // },
        // 表格数据
        getpurcon(){
            this.loading = true; 
            let listByPageData;
            var logisticsAll = JSON.parse(sessionStorage.getItem('queryAll'));
            if(JSON.parse(sessionStorage.getItem('queryAlls'))){
                let arrAll = JSON.parse(sessionStorage.getItem('queryAlls'));
                // 传参
                listByPageData = {
                    orderId:arrAll.id,
                    notBindOrder:true,
                    searchStr:this.purSearch,
                    billNumber:this.purSearchtree,
                    billData:this.keyupTime,
                    loadingAddress:this.purSearchone,
                    unloadingPoint:this.purSearchtwo,
                    page:{
                        current:this.current,
                        size:this.size
                    }
                }
            }else{
                // 传参
                listByPageData = {
                    notBindOrder:true,
                    searchStr:this.purSearch,
                    billNumber:this.purSearchtree,
                    billData:this.keyupTime,
                    loadingAddress:this.purSearchone,
                    unloadingPoint:this.purSearchtwo,
                    page:{
                        current:this.current,
                        size:this.size
                    }
                }
            }
            // 渲染表格
            this.api.waybill.all(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    this.loading = false; 
                    if(res.data.data.records != ""){
                        this.tableData = res.data.data.records;
                        this.total = res.data.data.total;
                        if (logisticsAll.length > 0) {
                            for (var i = 0; i < logisticsAll.length; i++) {
                                for (var j = 0; j < this.tableData.length; j++) {
                                    if (logisticsAll[i].id === this.tableData[j].id) {//items里面包含有的都回显打勾                 
                                        this.$refs.multipleTable.toggleRowSelection(this.tableData[j],true);
                                    }
                                }
                            }
                        }
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
        // 发货地搜索功能
        searchone() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 卸货地搜索功能
        searchtwo() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 运单编号搜索功能
        searchtree() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 日期搜索
        DataTime() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 表格选中数据
        purChange(val) {
            // console.log(this.multipleSelection)
            // this.multipleSelection = [];
            this.multipleSelection = val;
        },
        getRowKey(row) {
            return row.id;
        },
        // 保存关联信息
        AddDatasubmitForm() {
            let that = this;
            that.$message({
                type: "success",
                message: "保存关联信息成功!",
                duration:500,  
                onClose(){
                   
                    sessionStorage.removeItem('queryAll');
                    sessionStorage.setItem('queryAll',JSON.stringify(that.multipleSelection))
                    that.$router.push({name:"EditWaybill"})
                }
            });
        },
        // 取消保存关联信息
        Addopen() {
            this.$router.push({name:"EditWaybill"})
        },
    },
    computed: {},
    watch: {}
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
    .pur-size {
        font-size: 16px;
        line-height: 32px;
        color: #333;
        font-weight: 600;
    }
    .pur-right {
        .el-button {
            margin-left: 16px;
            float: right;
        }
    }
}
// 头部结束
// 表单内容开始

    
.pur-nav {
  height: calc(100vh - 110px);
  background-color: #fff;
  margin:0 60px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 17px;
  box-sizing: border-box;
  .searchs{
      height: 50px;
    .el-input {
        width: 12%;
        margin:7px 16px 0 0 ;
    }
  }
  .pur-table{
    width:100%;
    height: 88%;
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
      left:35%;
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
</style>