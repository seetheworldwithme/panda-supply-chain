<template>
  <!-- 云仓管理页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3">
        <span class="pur-size">{{ purTitle }}</span></el-col>
        <el-col :span="21" class="pur-right">
            <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button>
            <el-button size="small" @click="addPurs" v-if="$route.meta.add">新增</el-button>
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
          style="width: 100%;"
          @selection-change="purChange">
            <el-table-column type="selection"></el-table-column>
            <el-table-column prop="id" v-if="show">
            </el-table-column>
            <el-table-column label="煤仓仓号" :show-overflow-tooltip="true" prop="warehouseNb">
            </el-table-column>
            <el-table-column label="一级大类名称" :show-overflow-tooltip="true" prop="fstName">
            </el-table-column>
            <el-table-column label="二级小类名称" :show-overflow-tooltip="true" prop="secName">
            </el-table-column>
            <el-table-column label="三级小类名称" :show-overflow-tooltip="true" prop="threeName">
            </el-table-column>
            <el-table-column
              label="操作" width="80" fixed="left" v-if="$route.meta.edit">
              <template slot-scope="scope">
                <el-button @click="EditPurs(scope.row)" size="mini">编辑</el-button>
              </template>
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable">
              <div class="dataAva">
                <img src="@/assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
                <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新增</el-button>
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
      <!-- 新建仓号信息弹窗 -->
      <el-dialog
        :closeOnClickModal="false"
        title="新建仓号信息"
        :visible.sync="paymentShow"
        width="35%"
        :before-close="paymentShowClose">
          <el-form style="width:95%;" :model="AddData" :rules="paymentrules" ref="AddData" size="small" label-width="120px">
            <el-form-item label="煤仓仓号" prop="warehouseNb">
              <el-input v-model="AddData.warehouseNb" placeholder="请输入煤仓仓号"></el-input>
            </el-form-item>
            <el-form-item label="一级大类名称" prop="fstProId">
                <el-select @change="enterprises($event)" v-model="AddData.fstName" placeholder="请选择一级大类名称" style="width:100%;">
                    <el-option
                    v-for="item in provinceOptions"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="二级小类名称" prop="bigProId">
                <el-select @change="entisetwos($event)" v-model="AddData.secName" placeholder="请选择二级小类名称" style="width:100%;">
                    <el-option
                    v-for="item in cityOptions"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="三级小类名称" prop="threeProId">
                <el-select @change="citys($event)" v-model="AddData.threeName" placeholder="请选择三级小类名称" style="width:100%;" >
                    <el-option
                    v-for="item in options"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="paymentsubmitm('AddData')" style="float:right;margin:5% 0 0 2%;">确定</el-button>
              <el-button @click="paymentresetForm" style="float:right;margin:5% 0 0 2%;">取消</el-button>
            </el-form-item>
          </el-form>
      </el-dialog>
      <!-- 编辑仓号信息弹窗 -->
      <el-dialog
        title="编辑仓号信息"
        :closeOnClickModal="false"
        :visible.sync="editpaymentShow"
        width="35%"
        :before-close="editpaymentShowClose">
          <el-form style="width:95%;" :model="editData" :rules="editpaymentrules" ref="editData" size="small" label-width="120px">
            <el-form-item label="煤仓仓号" prop="warehouseNb">
              <el-input v-model="editData.warehouseNb" placeholder="请输入煤仓仓号"></el-input>
            </el-form-item>
            <el-form-item label="一级大类名称" prop="fstProId">
                <el-select @change="enterprises($event)" v-model="editData.fstName" placeholder="请选择一级大类名称" style="width:100%;">
                    <el-option
                    v-for="item in provinceOptions"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="二级小类名称" prop="bigProId">
                <el-select @change="entisetwos($event)" v-model="editData.secName" placeholder="请选择二级小类名称" style="width:100%;">
                    <el-option
                    v-for="item in cityOptions"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="三级小类名称" prop="threeProId">
                <el-select @change="citys($event)" v-model="editData.threeName" placeholder="请选择三级小类名称" style="width:100%;" >
                    <el-option
                    v-for="item in options"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="editpaymentsubmitm('editData')" style="float:right;margin:5% 0 0 2%;">确定</el-button>
              <el-button @click="editpaymentresetForm" style="float:right;margin:5% 0 0 2%;">取消</el-button>
            </el-form-item>
          </el-form>
      </el-dialog>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
    data() {
        return {
            // 车辆详情
            delistData:[],
            loading:false,
            current:1,              // 当前页数
            size:200,                // 显示条数
            total:0,                // 总条数
            delmul:[],              // 批量删除
            show:false,
            purTitle: "",           // 标题
            purSearch: "",          // 搜索        
            multipleSelection: [],  // 选中数据
            tableData: [],          // 表格数据  
            // 新建仓号信息弹窗
            paymentShow:false,
            AddData:{},
            paymentrules: {
                warehouseNb: [
                    { required: true, message: '请输入煤仓仓号', trigger: 'blur' },
                ],
                fstProId: [
                    { required: true, message: '请选择一级分类', trigger: 'blur' },
                ],
                bigProId: [
                    { required: true, message: '请选择二级分类', trigger: 'blur' },
                ],
                threeProId: [
                    { required: true, message: '请选择三级分类', trigger: 'blur' },
                ],
            },
            provinceOptions:[],             // 一级大类
            cityOptions: [],                // 二级小类
            options: [],                    // 三级小类
            // 仓号编辑
            editpaymentShow:false,
            editData:{},
            editpaymentrules: {
                warehouseNb: [
                    { required: true, message: '请输入煤仓仓号', trigger: 'blur' },
                ],
                fstProId: [
                    { required: true, message: '请选择一级分类', trigger: 'blur' },
                ],
                bigProId: [
                    { required: true, message: '请选择二级分类', trigger: 'blur' },
                ],
                threeProId: [
                    { required: true, message: '请选择三级分类', trigger: 'blur' },
                ],
            },
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
                    current: this.current,
                    size: this.size
                },
            }
            // 渲染表格
            this.api.ware.all(listByPageData)
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

        // 一级选中
        enterprises(e){
            this.cityOptions = [];               // 二级小类
            this.options = [];                    // 三级小类
            // this.AddData.threeName = null;
            let arr = e.split(",");
            this.AddData.secName = null;
            this.AddData.fstProId = arr[0]*1;
            this.AddData.fstName = arr[1];
            this.editData.fstProId = arr[0]*1;
            this.editData.fstName = arr[1];
            this.editData.secName = null;
            // 换了一级大类:二、三级一并清掉,否则显示为空、id 还是旧的,必填校验照样通过
            this.AddData.bigProId = this.editData.bigProId = null;
            this.AddData.threeProId = this.editData.threeProId = null;
            this.AddData.threeName = this.editData.threeName = null;
            this.list(arr[0]*1);
        },
        // 二级选中
        entisetwos(e){
            this.options = [];                    // 三级小类
            // this.AddData.threeName = null;
            let arr = e.split(",");
            this.AddData.bigProId = arr[0]*1;
            this.AddData.secName = arr[1];
            this.editData.bigProId = arr[0]*1;
            this.editData.secName = arr[1];
            // 换了二级:三级一并清掉(同上)
            this.AddData.threeProId = this.editData.threeProId = null;
            this.AddData.threeName = this.editData.threeName = null;
            this.twolists(arr[0]*1);
        },
        // 三级选中
        citys(e){
            console.log(e)
            let arr = e.split(",");
            this.AddData.threeProId = arr[0]*1;
            this.AddData.threeName = arr[1];
            this.editData.threeProId = arr[0]*1;
            this.editData.threeName = arr[1];
            // console.log(this.AddData)
        },
        // 云仓一级大类
        road(){
            this.provinceOptions = [];
            this.api.Prodect.list({parent_id:0,searchType:true})
            .then(res=>{
                if(res.data.code == 200){
                    this.provinceOptions = res.data.data;
                    
                }
            })
        },
        // 仓库产品二级小类类
        list(id){
            this.cityOptions = [];
            this.api.Prodect.list({parent_id:id,searchType:true})
            .then(res=>{
                if(res.data.code == 200){
                    this.cityOptions = res.data.data;
                }
            })
        },
        // 仓库产品三级小类类
        twolists(id){
            this.options = [];
            this.api.Prodect.list({parent_id:id,fst_pro_id:this.AddData.fstProId,searchType:true})
            .then(res=>{
                if(res.data.code == 200){
                    this.options = res.data.data;
                }
            })
        },
        // 新建仓号信息弹窗
        addPurs(){
            this.paymentShow = true;
            this.road();
        },
        // 新建仓号取消
        paymentShowClose(){
            this.paymentShow = false;
            this.AddData = {};
            this.provinceOptions = [];
            this.cityOptions = [];
            this.options = [];
        },
        // 取消新建仓号提交
        paymentresetForm() {
            this.paymentShow = false;
            this.AddData = {};
            this.provinceOptions = [];
            this.cityOptions = [];
            this.options = [];
        },
        // 新建仓号信息提交
        paymentsubmitm(formName) {
            this.$refs[formName].validate((valid) => {
            if (valid) {
                this.api.ware.add(this.AddData)
                .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                        type: 'success',
                        message: '保存成功!'
                    });
                    this.AddData = {}
                    this.tableData = [];
                    this.current = 1;
                    this.getpurcon();
                    this.paymentShow = false;
                }else if(res.data.code == 602){
                    this.$message({
                        type: 'warning',
                        message: res.data.msg
                    });
                }
                })
            }
            });
        },
        // 编辑仓号信息
        EditPurs(data){
            this.editpaymentShow = true;
            this.road();
            this.editData = JSON.parse(JSON.stringify(data));
            // this.api.ware.one({id:data.id})
            // .then(res=>{
            // if(res.data.code == 200){
            //     let ids = [];
            //     for(let key in res.data.data.carfences){
            //     ids.push(res.data.data.carfences[key].fenceId)
            //     }
            //     this.editData.ids = ids;
            // }
            // })
        },
        // 编辑仓号信息取消
        editpaymentShowClose(){
            this.editpaymentShow = false;
            this.editData = {}
            this.provinceOptions = [];
            this.cityOptions = [];
            this.options = [];
        },
        // 取消编辑仓号信息提交
        editpaymentresetForm() {
            this.editpaymentShow = false;
            this.editData = {}
            this.provinceOptions = [];
            this.cityOptions = [];
            this.options = [];
        },
        // 编辑仓号信息提交
        editpaymentsubmitm(formName) {
            this.$refs[formName].validate((valid) => {
            if (valid) {
                this.api.ware.edit(this.editData)
                .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                        type: 'success',
                        message: '编辑成功!'
                    });
                    this.editData = {}
                    this.tableData = [];
                    this.current = 1;
                    this.getpurcon();
                    this.editpaymentShow = false;
                }else if(res.data.code == 602){
                    // 原先这段挂在 if (valid) 上:表单校验不通过时 res 还不存在,直接抛 ReferenceError
                    this.$message({
                        type: 'warning',
                        message: res.data.msg
                    });
                }
                })
            }
            });
        },
        // 表格选中数据
        purChange(val) {
            this.multipleSelection = val;
            this.delmul = [];
            for(let i=0;i<this.multipleSelection.length;i++){
            this.delmul.push(this.multipleSelection[i].id) 
            }
        },
        // 删除功能
        openPur() {
            if(this.multipleSelection == ""){
                this.$message({
                    type: 'warning',
                    message: '请选择要删除的数据!'
                });
            }else{
                this.$confirm('此操作将永久删除该仓号信息, 是否继续?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                this.api.ware.del(this.delmul)
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
// 详情样式开始
.detilsBox{
  width: 100%;
  height: 100%;
  padding: 25px;
  box-sizing: border-box;
  h5{
    font-size: 16px;
    line-height: 58px;
    color: #333333;
  }
  .detilconent{
    width: 100%;
    span{
      display: inline-block;
      width: 33.33%;
      font-size: 14px;
      line-height: 24px;
      color: #333333;
    }
  }
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