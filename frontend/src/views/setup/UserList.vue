<template>
<!-- 用户列表页面 -->
    <div> 
        <!-- tab选项卡 -->
        <div class="menu-list">
            <div class="totalTab"
                :class="{ active: item.id === sel }"
                v-for="item in tabs"
                :key="item.id"
                @click="select(item)" style="cursor: pointer;">
            <span>{{ item.label }}</span>
            <div class="stripnone" :style="item.width"
             :class="{ strip: item.id === sel }"></div>
        </div>
        </div>
        <div class="parNavs">
            <div v-show="sel === 1">
                <!-- 待授权表格 -->
                <div class="pur-nav"
                v-loading="loading"
                element-loading-text = "加载中...">
                    <div class="pur-table">
                        <el-table
                        height="100%"
                        :data="UnauthorizedData"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;">
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column label="操作" fixed="right">
                        <template slot-scope="scope">
                            <el-button type="primary" plain size="mini" @click="UnautAdopt(scope.row.id)">通过</el-button>
                            <el-button type="info" size="mini" @click="UnautRefuse(scope.row.id)">拒绝</el-button>
                        </template>
                        </el-table-column>
                        <el-table-column label="手机号" prop="phone" :show-overflow-tooltip="true" width="140">
                        </el-table-column>
                        <el-table-column label="姓名" prop="username" :show-overflow-tooltip="true" width="120">
                        </el-table-column>
                        <el-table-column label="账号" prop="account" :show-overflow-tooltip="true" width="200">
                        </el-table-column>
                        <!-- 暂无数据 -->
                        <div slot="empty" class="dataAvailable">
                            <div class="dataAva">
                                <img src="../../assets/zanwushuju.png" alt="">
                                <p>暂无数据！</p>
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
                        @current-change='Unauthorized'
                        @size-change='Unauthorized'>
                        </el-pagination>
                    </div> 
                </div> 
            </div>
            <div v-show="sel === 2">
                <!-- 已授权表格 -->
                <div class="pur-nav"
                v-loading="loadingOne"
                element-loading-text = "加载中...">
                    <div class="pur-table">
                        <el-table
                        height="100%"
                        :data="AuthorizedData"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;">
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column label="操作" fixed="right" width="200">
                            <template slot-scope="scope">
                                <el-button type="primary" plain size="small" @click="AuthorAdopt(scope.row.id,scope.row.level)" v-if="$route.meta.edit">分配</el-button>
                                <el-button size="small" @click="AuthorRefuse(scope.row.id)" v-if="$route.meta.delete">删除</el-button>
                            </template>
                        </el-table-column>
                        <el-table-column label="手机号" prop="phone" :show-overflow-tooltip="true" width="140">
                        </el-table-column>
                        <el-table-column label="姓名" prop="username" :show-overflow-tooltip="true" width="120">
                        </el-table-column>
                        <el-table-column label="账号" prop="account" :show-overflow-tooltip="true" width="150">
                        </el-table-column>
                        <el-table-column label="角色">
                            <template slot-scope="scope">
                                <el-tag v-for="(val,kal) in scope.row.roleNames" :key="kal" type="warning" size="small" style="margin:0 0.8% 0 0;">{{val}}</el-tag>
                            </template>
                        </el-table-column>
                        <!-- 暂无数据 -->
                        <div slot="empty" class="dataAvailable">
                            <div class="dataAva">
                                <img src="../../assets/zanwushuju.png" alt="">
                                <p>暂无数据！</p>
                            </div>
                        </div>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        background
                        :page-sizes="[10, 30, 50, 100, 200]"
                        :page-size.sync="sizeOne"
                        :pager-count="5"
                        :current-page.sync="currentOne"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="totalOne"
                        @current-change='Authorized'
                        @size-change='Authorized'>
                        </el-pagination>
                    </div> 
                </div> 
                <!-- 分配角色弹窗 -->
                <el-dialog title="分配角色" :visible.sync="EditRoleListVisible" width="20%" :before-close="EditRoleListClose">
                    <el-form :model="EditRoleListData" size="small">
                    <el-form-item style="margin-bottom:8%;">
                        <el-select
                        v-model="EditRoleListData.roleLists"
                        multiple
                        collapse-tags
                        style="margin-left: 20px;"
                        placeholder="请选择角色">
                            <el-option
                            v-for="item in rolelist"
                            :key="item.id"
                            :label="item.name"
                            :value="item.id">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item style="margin-bottom:0;">
                        <el-button type="primary" size="small" @click="EditRoleListsubmitForm('EditRoleListData')" style="float:right;">确定</el-button>
                        <el-button @click="EditRoleListVisibleclick" size="small" style="float:right;margin-right:5%;">取消</el-button>
                    </el-form-item>
                    </el-form>
                </el-dialog>
            </div>
            <!-- 未通过 -->
            <div v-show="sel === 3">
                <!-- 未通过表格 -->
                <div class="pur-nav"
                v-loading="loadingTwo"
                element-loading-text = "加载中...">
                    <div class="pur-table">
                        <el-table
                        height="100%"
                        :data="InvalidData"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;">
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column label="操作" fixed="right">
                        <template slot-scope="scope">
                            <el-button type="info" size="small" @click="InvalidRefuse(scope.row.id)">恢复</el-button>
                        </template>
                        </el-table-column>
                        <el-table-column label="手机号" prop="phone" :show-overflow-tooltip="true" width="140">
                        </el-table-column>
                        <el-table-column label="姓名" prop="username" :show-overflow-tooltip="true" width="120">
                        </el-table-column>
                        <el-table-column label="账号" prop="account" :show-overflow-tooltip="true" width="200">
                        </el-table-column>
                        <!-- 暂无数据 -->
                        <div slot="empty" class="dataAvailable">
                            <div class="dataAva">
                                <img src="../../assets/zanwushuju.png" alt="">
                                <p>暂无数据！</p>
                            </div>
                        </div>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        background
                        :page-sizes="[10, 30, 50, 100, 200]"
                        :page-size.sync="sizeTwo"
                        :pager-count="5"
                        :current-page.sync="currentTwo"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="totalTwo"
                        @current-change='Invalid'
                        @size-change='Invalid'>
                        </el-pagination>
                    </div> 
                </div> 
            </div>
        </div>
    </div>
</template>
<script>
export default {
    data(){
        return {
            sel: 1,   // 选项卡默认显示tab
            tabs: [
                { label: "待授权", id: 1 ,width:"width:48px;"},
                { label: "已授权", id: 2 ,width:"width:48px;"},
                { label: "未通过", id: 3 ,width:"width:48px;"},
            ],
            show:false,  //隐藏表格id
            loading:false,          // 待授权数据状态
            current:1,              // 待授权数据当前页数
            size:10,                // 待授权数据显示条数
            total:0,                // 待授权数据总条数
            UnauthorizedData:[],    // 待授权数据
            loadingOne:false,       // 已授权数据状态
            currentOne:1,           // 已授权数据当前页数
            sizeOne:10,             // 已授权数据显示条数
            totalOne:0,             // 已授权数据总条数
            AuthorizedData:[],      // 已授权数据
            loadingTwo:false,       // 未通过状态
            currentTwo:1,           // 未通过前页数
            sizeTwo:10,             // 未通过显示条数
            totalTwo:0,             // 未通过总条数
            InvalidData:[],         // 未通过数据
            /*
            分配角色弹窗
            */ 
            userID:'',                      // 选中用户id
            EditRoleListVisible:false,      // 分配角色弹窗开关     
            EditRoleListData: {
                roleLists:[],
            },           // 分配角色表单
            rolelist:[],                    // 获取角色列表
        }
    },
    created() {
        this.Unauthorized();        // 待授权
    },
    methods: {
        /*
        分配角色功能
        */
        AuthorAdopt(id){
            this.EditRoleListVisible = true;
            this.userID = id;
            // 分配角色
            this.api.roleList.AllQy()
            .then(res=>{
                if(res.data.code == 200){
                    this.rolelist = res.data.data;
                    // 回显分配好的角色
                    this.api.roleUserList.selUserRoles({userId:this.userID})
                    .then(res=>{
                        if(res.data.code == 200){
                            for(let key in res.data.data){
                                this.EditRoleListData.roleLists.push(res.data.data[key].roleId);
                            }
                        }
                    })
                }
            })

        },
        // 取消分配角色
        EditRoleListVisibleclick(){
            this.EditRoleListVisible = false;
            this.EditRoleListData.roleLists = [];
        },
        EditRoleListClose(){
            this.EditRoleListData.roleLists = [];
            this.EditRoleListVisible = false;
        },
        // 分配角色提交表单 
        EditRoleListsubmitForm() {
            let roleObj = {
                id:this.userID,
                userRoles:[]
            };
            for(let key in this.EditRoleListData.roleLists){
                let obj = {
                    userId:this.userID,
                    roleId:this.EditRoleListData.roleLists[key]
                };
                roleObj.userRoles.push(obj) 
            }
            this.api.roleUserList.assignRole(roleObj)
            .then(res=>{
                if(res.data.code == 200){
                    var that = this;
                    that.$message({
                        type: "success",
                        message: "分配角色成功!",
                        duration:500,  
                        onClose(){
                            that.EditRoleListData.roleLists = [];             // 清空角色列表
                            that.EditRoleListVisible = false;  
                            that.Authorized();
                            that.sel = 2;
                        }
                    });
                }
            })
        },

        //选项卡
        select(item) {
            if(item.id == 1){
                // 待授权
                this.sel = item.id;
                this.Unauthorized();
            }
            if(item.id == 2){
                // 已授权
                this.sel = item.id;
                this.Authorized();
            }
            if(item.id == 3){
                // 未通过
                this.sel = item.id;
                this.Invalid();
            }   
        },
        // 待授权接口
        Unauthorized(){
            this.loading = true;
            // 传参
            let listByPageData = {
                accountType:1,
                handleResult:0,
                searchStr:'',
                page:{
                    current:this.current,
                    size:this.size
                }
            }
            this.api.authAll.listByPage(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    this.loading = false;
                    this.UnauthorizedData = res.data.data.records;
                    this.total = res.data.data.total;
                }
            })
        },
        // 已授权接口
        Authorized(){
            this.loadingOne = true;
            // 传参
            let listByPageData = {
                accountType:1,
                handleResult:1,
                searchStr:'',
                page:{
                    current:this.currentOne,
                    size:this.sizeOne
                }
            }
            this.api.authAll.listByPage(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    this.loadingOne = false;
                    for(let key in res.data.data.records){
                        if(res.data.data.records[key].roleNames == null){
                            res.data.data.records[key].roleNames = ['未分配角色'];
                        }else{
                            let arr = res.data.data.records[key].roleNames;
                            res.data.data.records[key].roleNames = arr.split(',');
                        }
                    }
                    this.AuthorizedData = res.data.data.records;
                    this.totalOne = res.data.data.total;
                }
            })
        },
        // 未通过接口
        Invalid(){
            this.loadingTwo = true;
            // 传参
            let listByPageData = {
                accountType:1,
                handleResult:2,
                searchStr:'',
                page:{
                    current:this.currentTwo,
                    size:this.sizeTwo
                }
            }
            this.api.authAll.listByPage(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    this.loadingTwo = false;
                    this.InvalidData = res.data.data.records;
                    this.totalTwo = res.data.data.total;
                }
            })
        },
        
        // 恢复账号
        InvalidRefuse(id){
            this.$confirm('此操作将恢复该账号及权限, 请谨慎操作，是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.api.authAll.updateById({id:id,handleResult:0})
                .then(res=>{
                    if(res.data.code == 200){
                        this.sel = 3;
                        this.Invalid();
                    }
                })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消恢复账号！'
                });          
            });

        },
        // 待授权通过
        UnautAdopt(id){
            this.api.authAll.updateById({id:id,handleResult:1})
            .then(res=>{
                if(res.data.code == 200){
                    this.sel = 1;
                    this.Unauthorized();
                }
            })
        },
        // 待授权拒绝
        UnautRefuse(id){
            this.api.authAll.updateById({id:id,handleResult:2})
            .then(res=>{
                if(res.data.code == 200){
                    this.sel = 1;
                    this.Unauthorized();
                }
            })
        },
        // 删除账号
        AuthorRefuse(id){
            this.$confirm('此操作将永久删除该账号, 请谨慎操作，是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.api.authAll.deleteById({id:id})
                .then(res=>{
                    if(res.data.code == 200){
                        this.sel = 2;
                        this.Authorized();
                    }
                })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除账号！'
                });          
            });
        },
    }
}
</script>
<style lang="less" scoped>

// 选项卡开始
.menu-list{
    width: 100%;
    height: 50px;
    background-color: #fff;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    border: 1px solid rgba(102, 102, 102, 0.15);
}
.totalTab {
    width: 7.6%;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    color: #333;
    text-align: center;
    margin:9px 0 0 35px;
}
.totalTab.active {
    color: #ea222e;
}
.stripnone{
    height: 3px;
    background-color: #ea222e;
    margin:11px auto 0;
    display: none;
}
.strip {
    display: block;
}
.pur-nav {
  width: calc(100% - 32px);
  height: calc(100vh - 144px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 8px 8px 8px 8px;
  overflow: hidden;
  .pur-table{
    width:100%;
    height: 94%;
  } 
}
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

//选项卡结束
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
//选项卡结束
</style>