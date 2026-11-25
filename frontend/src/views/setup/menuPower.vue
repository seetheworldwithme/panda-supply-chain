<template>
  <!-- 菜单管理页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur"  v-if="$route.meta.delete">删除</el-button>
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
            <el-table-column width="50" type="selection" fixed="left"> </el-table-column>
            <el-table-column prop="id" v-if="show">
            </el-table-column>
            <el-table-column v-if="$route.meta.edit"
              label="操作"
              width="100" fixed="left">
              <template slot-scope="scope">
                <el-button size="mini" @click="EditPurs(scope.row)">编辑</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="菜单名称" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="routeName" label="路由名称" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="menu" label="菜单标识" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="pname" label="父级菜单" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="isAuthy" label="是否为权限菜单" :show-overflow-tooltip="true">
                <template slot-scope="scope"> 
                    <span v-if="scope.row.isAuthy">是</span>
                    <span v-else>否</span>
                </template>
            </el-table-column>
            <el-table-column prop="hasFunChild" label="是否包含功能子菜单" :show-overflow-tooltip="true">
                <template slot-scope="scope"> 
                    <span v-if="scope.row.hasFunChild">是</span>
                    <span v-else>否</span>
                </template>
            </el-table-column>
            <el-table-column prop="openAuthy" label="是否开启权限控制" :show-overflow-tooltip="true">
                <template slot-scope="scope"> 
                    <span v-if="scope.row.openAuthy">是</span>
                    <span v-else>否</span>
                </template>
            </el-table-column>
            <el-table-column prop="disabled" label="是否禁用" :show-overflow-tooltip="true">
                <template slot-scope="scope"> 
                    <span v-if="scope.row.disabled">是</span>
                    <span v-else>否</span>
                </template>
            </el-table-column>
            <el-table-column prop="disabled" label="是否是子菜单" :show-overflow-tooltip="true">
                <template slot-scope="scope"> 
                    <span v-if="scope.row.isSon">是</span>
                    <span v-else>否</span>
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
        :current-page.sync="current"
        :pager-count="5"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @current-change='getpurcon'
        @size-change='getpurcon'>
        </el-pagination>
      </div>
        <!-- 新增弹窗 -->
        <el-dialog
        title="新增菜单"
        :visible.sync="dialogVisibledata"
        width="50%"
        :before-close="BuyerhandleClose">
            <div class="elTabBox">
                <!-- label-width="120px" -->
                <el-form :model="AdjustForm" :rules="rules" ref="AdjustForm" size="small"  label-position="left">  
                    <el-row>
                        <el-col :span="24" style="margin-top:1%;">
                            <el-form-item prop="pid" label="菜单类型" label-width="120px">
                                <el-select size="small" filterable v-model="AdjustForm.menuType" placeholder="请选择菜单类型" style="width:80%;">
                                    <el-option
                                    v-for="item in projectone"
                                    :key="item.value"
                                    :label="item.name"
                                    :value="item.id">
                                    </el-option>
                                </el-select>   
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item prop="pid" label="父级菜单" label-width="120px">
                                <el-select size="small" filterable v-model="AdjustForm.pid" placeholder="请选择父级菜单" style="width:80%;">
                                    <el-option
                                    v-for="item in projectName"
                                    :key="item.value"
                                    :label="item.name"
                                    :value="item.id">
                                    </el-option>
                                </el-select>   
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item prop="name" label="菜单名称" label-width="120px">
                                <el-input size="small" v-model="AdjustForm.name" placeholder="请输入菜单名称" style="width:80%;"></el-input>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="菜单标识" label-width="120px">
                                <el-input size="small" v-model="AdjustForm.menu" placeholder="请输入菜单标识" style="width:80%;"></el-input>  
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item prop="routeName" label="路由名称" label-width="120px">
                                <el-input size="small" v-model="AdjustForm.routeName" placeholder="请输入路由名称" style="width:80%;"></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否为权限菜单" label-width="180px">
                                <el-radio-group v-model="AdjustForm.isAuthy">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group>   
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否开启权限控制" label-width="180px">
                                <el-radio-group v-model="AdjustForm.openAuthy">   

                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否包含功能子菜单" label-width="180px">
                                <el-radio-group v-model="AdjustForm.hasFunChild">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否禁用" label-width="180px">
                                <el-radio-group v-model="AdjustForm.disabled">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否是子菜单" label-width="180px">
                                <el-radio-group v-model="AdjustForm.isSon">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="24" style="margin-top:1%;">
                            <el-form-item label="操作列表">
                                <div style="margin:44px 0"></div>
                                <el-checkbox :indeterminate="isIndeterminate" v-model="checkAll" @change="handleCheckAllChange" border>全选</el-checkbox>
                                <div style="margin:25px 0"></div>
                                <el-checkbox-group v-model="checkedCities" @change="handleCheckedCitiesChange">
                                    <el-checkbox v-for="(city,k) in cities" :label="city.id" :key="k" border>{{city.zhname}}</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <div class="adjusttop">
                        <el-button size="small" style="margin-left: 24px"  @click="dialogVisibledatass">取 消</el-button>
                        <el-button type="primary" size="small" @click="AddDatasubmitForm('AdjustForm')">确 定</el-button>
                    </div>
                </el-form>
            </div>
        </el-dialog>
        <!-- 编辑弹窗 -->
        <el-dialog
        title="编辑菜单"
        :visible.sync="AdjustVisibledata"
        width="50%"
        :before-close="editBuyerhandleClose">
            <div class="elTabBox">
                <el-form :model="editAdjustForm" :rules="editrules" ref="editAdjustForm" size="small" label-position="left">
                    <el-row>
                        <el-col :span="24" style="margin-top:1%;">
                            <el-form-item prop="pid" label="菜单类型" label-width="120px">
                                <el-select size="small" filterable v-model="editAdjustForm.menuType" placeholder="请选择菜单类型" style="width:80%;">
                                    <el-option
                                    v-for="item in projectone"
                                    :key="item.value"
                                    :label="item.name"
                                    :value="item.id">
                                    </el-option>
                                </el-select>   
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item prop="pid" label="父级菜单" label-width="120px">
                                <el-select size="small" filterable v-model="editAdjustForm.pid" placeholder="请选择父级菜单" style="width:80%;">
                                    <el-option
                                    v-for="item in editprojectName"
                                    :key="item.value"
                                    :label="item.name"
                                    :value="item.id">
                                    </el-option>
                                </el-select>   
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item prop="name" label="菜单名称" label-width="120px">
                                <el-input size="small" v-model="editAdjustForm.name" placeholder="请输入菜单名称" style="width:80%;"></el-input>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="菜单标识" label-width="120px">
                                <el-input size="small" v-model="editAdjustForm.menu" placeholder="请输入菜单标识" style="width:80%;"></el-input>  
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item prop="routeName" label="路由名称" label-width="120px">
                                <el-input size="small" v-model="editAdjustForm.routeName" placeholder="请输入路由名称" style="width:80%;"></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否为权限菜单" label-width="180px">
                                <el-radio-group v-model="editAdjustForm.isAuthy">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group>   
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否开启权限控制" label-width="180px">
                                <el-radio-group v-model="editAdjustForm.openAuthy">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否包含功能子菜单" label-width="180px">
                                <el-radio-group v-model="editAdjustForm.hasFunChild">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否禁用" label-width="180px">
                                <el-radio-group v-model="editAdjustForm.disabled">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="12" style="margin-top:1%;">
                            <el-form-item label="是否是子菜单" label-width="180px">
                                <el-radio-group v-model="editAdjustForm.isSon">
                                <el-radio :label="true">是</el-radio>
                                <el-radio :label="false">否</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                        </el-col>
                        <el-col :span="24" style="margin-top:1%;">
                            <el-form-item label="操作列表">
                                <!-- <div style="margin:44px 0"></div>
                                <el-checkbox :indeterminate="editisIndeterminate" v-model="editcheckAll" @change="edithandleCheckAllChange" border>全选</el-checkbox> -->
                                <div style="margin:45px 0"></div>
                                <el-checkbox-group v-model="editcheckedCities">
                                    <el-checkbox v-for="(city,k) in editcities" :label="city.id" :key="k" border>{{city.zhname}}</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <div class="adjusttop">
                        <el-button size="small" style="margin-left: 24px"  @click="AdjustVisibledatass">取 消</el-button>
                        <el-button type="primary" size="small" @click="editDatasubmitForm('editAdjustForm')">确 定</el-button>
                    </div>
                </el-form>
            </div>
        </el-dialog>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
  data() {
    return {
        projectone:[
            {
                name:'企业菜单',
                id:1
            },
            {
                name:'平台菜单',
                id:2
            },{
                name:'通用菜单',
                id:3
            }
        ],
        show:false,
        purTitle: "",                       // 标题
        dataAva:false,
        loading:false,
        current:1,                          // 当前页数
        size:50,                            // 显示条数
        total:0,                            // 总条数
        delmul:[],                          // 批量删除
        purSearch: "",                      // 搜索
        tableData: [],                      // 表格数据
        multipleSelection: [],              // 选中数据

        /*
        新增
        */
        dialogVisibledata:false,
        AdjustForm:{
            isAuthy:true,
            openAuthy:true,
            hasFunChild:true,
            disabled:false,
            // ids:[],
        },  
        rules: {
            menuType: [
                { required: true, message: '请选择菜单类型', trigger: 'change' }
            ],
            pid: [
                { required: true, message: '请选择父级菜单', trigger: 'change' }
            ],
            name: [
                { required: true, message: '请输入菜单名称', trigger: 'change' }
            ],
            routeName: [
                { required: true, message: '请输入路由名称', trigger: 'change' }
            ]
        },
        projectName:[],   // 父级菜单列表
        checkAll: false,
        checkedCities: [],
        cities:[],
        isIndeterminate: false,
        /*
        编辑
        */
        AdjustVisibledata:false,
        editAdjustForm:{},  
        editrules: {
            menuType: [
                { required: true, message: '请选择菜单类型', trigger: 'change' }
            ],
            pid: [
                { required: true, message: '请选择父级菜单', trigger: 'change' }
            ],
            name: [
                { required: true, message: '请输入菜单名称', trigger: 'change' }
            ],
            routeName: [
                { required: true, message: '请输入路由名称', trigger: 'change' }
            ]
        },
        editprojectName:[],   // 父级菜单列表
        editcheckedCities:[],
        editcities:[],
    };
  },
  created() {},
  mounted(){
    this.purTitle = this.$route.meta.title;             // 标题
    this.logo = sessionStorage.getItem('companyName')   // 企业名称
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
      this.api.authMenu.listByPage(listByPageData)
      .then(res=>{
        this.loading = false;
        if(res.data.code == 200){
          if(res.data.data.records != ""){
            this.tableData = res.data.data.records;
            this.total = res.data.data.total;
          }else{
            this.dataAva = true;
          }
        }
      })
    },
    //搜索功能
    search() {
      this.current = 1; // 先重置页码再查询,否则带着旧页码查不到数据
      this.tableData = [];
      this.getpurcon();
    },
    // 表格选中数据
    purChange(val) {
      this.multipleSelection = val;
      this.delmul = [];
      for(let i=0;i<this.multipleSelection.length;i++){
        this.delmul.push(this.multipleSelection[i].id);
      }
    },

    // 新建菜单信息弹窗(打开)
    addPurs(){
        this.dialogVisibledata = true;
        this.api.authMenu.selAuthys()
        .then(res=>{
            if(res.data.code == 200){
                this.cities = res.data.data;
            }
        });
        this.api.authMenu.list({is_authy:true})
        .then(res=>{
            if(res.data.code == 200){
                this.projectName = res.data.data;
            }
        });
    },
    // 新建菜单信息弹窗(关闭)
    BuyerhandleClose(){
        this.dialogVisibledata = false;
        this.checkAll = false;
        this.isIndeterminate = false;
        this.AdjustForm = {}
        this.checkedCities = []; // 不清的话下次新建菜单还勾着上一次的操作按钮,并随 ids 提交
    },
    // 新建菜单信息弹窗(取消)
    dialogVisibledatass(){
        this.dialogVisibledata = false;
        this.checkAll = false;
        this.isIndeterminate = false;
        this.AdjustForm = {}
        this.checkedCities = []; // 不清的话下次新建菜单还勾着上一次的操作按钮,并随 ids 提交
    },
    // 新建全选不全选逻辑
      handleCheckAllChange(val) {
        let arr = [];
        for(let i=0;i<this.cities.length;i++){
            arr.push(this.cities[i].id);
        }
        this.checkedCities = val ? arr : [];
        this.isIndeterminate = false;
      },
    handleCheckedCitiesChange(value) {
        let checkedCount = value.length;
        this.checkAll = checkedCount === this.cities.length;
        this.isIndeterminate = checkedCount > 0 && checkedCount < this.cities.length;
    },
    // 新建菜单信息表单提交
    AddDatasubmitForm(formName) {
        this.$refs[formName].validate((valid) => {
            if (valid) {
                this.AdjustForm.ids = this.checkedCities;
                this.api.authMenu.add(this.AdjustForm)
                .then(res=>{
                    if(res.data.code == 200){
                        this.$message({  
                            type: "success",
                            message: "新建成功!"
                        });
                        this.dialogVisibledata = false;
                        this.checkAll = false;
                        this.isIndeterminate = false;
                        this.AdjustForm = {}
                        this.checkedCities = []; // 同上
                        this.getpurcon();
                        this.$store.dispatch("asynSetMenus");
                    }
                })
            }
        });
    },


    // 编辑菜单信息弹窗(打开)
    EditPurs(val){
        this.editcheckedCities = [];
        this.AdjustVisibledata = true;
        this.api.authMenu.selAuthys()
        .then(res=>{
            if(res.data.code == 200){
                this.editcities = res.data.data;
            }
        });
        this.api.authMenu.list({is_authy:true})
        .then(res=>{
            if(res.data.code == 200){
                this.editprojectName = res.data.data;
            }
        });
        this.api.authMenu.getById({id:val.id})
        .then(res=>{
            if(res.data.code == 200){
                this.editAdjustForm = res.data.data;
                let a = res.data.data.menuAuths;
                for(let key in a){
                    this.editcheckedCities.push(a[key].authyId)
                }
            }
        }) 
    },
    // 编辑菜单信息弹窗(关闭)
    editBuyerhandleClose(){
        this.AdjustVisibledata = false;
        this.editAdjustForm = {}
    },
    // 编辑菜单信息弹窗(取消)
    AdjustVisibledatass(){
        this.AdjustVisibledata = false;
        this.editAdjustForm = {}
    },
    // 编辑全选不全选逻辑
    edithandleCheckAllChange(val) {
        let arr = [];
        for(let i=0;i<this.cities.length;i++){
            arr.push(this.cities[i].id);
        }
        this.ids = val ? arr : [];
        this.editisIndeterminate = false;
    },
    // 编辑提交
    editDatasubmitForm(formName){
        this.$refs[formName].validate((valid) => {
            if (valid) {
                this.editAdjustForm.ids = this.editcheckedCities;
                this.api.authMenu.edit(this.editAdjustForm)
                .then(res=>{
                    if(res.data.code == 200){
                        this.$message({  
                            type: "success",
                            message: "编辑成功!"
                        });
                        this.AdjustVisibledata = false;
                        this.editAdjustForm = {}
                        this.getpurcon();
                        this.$store.dispatch("asynSetMenus");
                    }
                })
            }
        });
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
              this.api.authMenu.del({ids:this.delmul})
              .then(res=>{
                if(res.data.code == 200){
                  this.$message({
                    type: 'success',
                    message: '删除成功!'
                  });
                  this.current = 1; // 先重置页码再查询,否则带着旧页码查不到数据
                  this.tableData = [];
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
// 弹窗样式
.elTabBox{
    padding: 0 30px 0 40px;
    box-sizing: border-box;
    position: relative;
    overflow: hidden;
    .adjusttop{
        margin-top:3%;
        float: right;
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
    margin-bottom:265px;
    .dataAva{
      position: absolute;
      top:75%;
      left: 25%;
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