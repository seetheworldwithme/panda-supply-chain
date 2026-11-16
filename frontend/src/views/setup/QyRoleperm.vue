<template>
  <!-- 角色权限页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3">
        <span class="pur-size">{{ purTitle }}</span>
      </el-col>
      <el-col :span="21" class="pur-right">
          <el-button type="primary" size="small" @click="AddDatasubmitForm">保 存</el-button>
      </el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav">
        <!-- 角色列表 -->
        <div class="nav-left">
          <div class="RoleList-title"> 
            <span>角色列表</span>
            <el-button style="float:right;padding:3px 0" type="text" @click="RoleListVisible = true">+ 新增角色</el-button>
          </div>
          <div class="RoleListBOX">
            <ul class="RoleList">
              <li v-for="(k,index) in rolelist" :key="index" :class="[currentIdx==index?'itemBg':'',currentCheckIdx==index?'checkedItemBg':'']" @mouseenter="enter(index)" @mouseleave="leave()" @click="currentCheckIdxclick(index,k.id)">
                <span>{{k.name}}</span>
                <template  v-if="currentIdx==index? showIcon = true:showIcon = false || currentCheckIdx==index? showIcon = true:showIcon = false">
                  <i class="el-icon-delete" @click="DelClick(k.id)"></i>
                  <i class="el-icon-edit" @click="EditRoleListVisibles(k.id)"></i>
                </template>
              </li>
            </ul>
          </div>

        </div>
        <!-- 新增角色弹窗 -->
        <el-dialog title="新增角色" :visible.sync="RoleListVisible" width="20%" :before-close="RoleListClose">
            <el-form :model="RoleListData" :rules="RoleListRules" ref="RoleListData" label-width="80px" size="small">
              <el-form-item label="角色名称" prop="name" style="margin-bottom:8%;">
                <el-input v-model="RoleListData.name" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item style="margin-bottom:0;">
                <el-button type="primary" size="small" @click="RoleListsubmitForm('RoleListData')" style="float:right;">保存</el-button>
                <el-button @click="RoleListVisibleclick" size="small" style="float:right;margin-right:5%;">取消</el-button>
              </el-form-item>
            </el-form>
        </el-dialog>
        <!-- 编辑角色弹窗 -->
        <el-dialog title="编辑角色" :visible.sync="EditRoleListVisible" width="20%" :before-close="EditRoleListClose">
            <el-form :model="EditRoleListData" :rules="EditRoleListRules" ref="EditRoleListData" label-width="80px" size="small">
              <el-form-item label="角色名称" prop="name" style="margin-bottom:8%;">
                <el-input v-model="EditRoleListData.name" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item style="margin-bottom:0;">
                <el-button type="primary" size="small" @click="EditRoleListsubmitForm('EditRoleListData')" style="float:right;">保存</el-button>
                <el-button @click="EditRoleListVisibleclick" size="small" style="float:right;margin-right:5%;">取消</el-button>
              </el-form-item>
            </el-form>
        </el-dialog>
        <!-- 角色权限 -->
        <div class="nav-right">
          <div class="clearfixstyle">
            <span>角色权限</span>
            <span style="float:right;">操作权限</span>
          </div>
          <div class="powerBox">
            <el-tree
            :data="treedata"
            show-checkbox
            node-key="id"
            ref="addtree"
            :props="defaultProps" 
            @check-change="handleNodeClick"
            :indent="30"
            default-expand-all
            >
              <span class="custom-tree-node" slot-scope="{ node, data }" style="width:100%;display:flex;justify-content: space-between;">
                  <span>{{ node.label }}</span>
                  <span style="margin-right:1%;">
                    <el-checkbox-group v-model="data.button" size="small">
                      <el-checkbox border label="导入" v-if="data.import !==-1"></el-checkbox>
                      <el-checkbox border label="导出" v-if="data.export !==-1"></el-checkbox>
                      <el-checkbox border label="新增" v-if="data.add !==-1" name="button"></el-checkbox>
                      <el-checkbox border label="查看" v-if="data.check !==-1" name="button"></el-checkbox>
                      <el-checkbox border label="编辑" v-if="data.edit !==-1" name="button"></el-checkbox>
                      <el-checkbox border label="删除" v-if="data.delete !==-1" name="button"></el-checkbox>
                      <el-checkbox border label="上传" v-if="data.upload !==-1" name="button"></el-checkbox>
                      <el-checkbox border label="文件删除" v-if="data.fileDelete !==-1" name="button"></el-checkbox>
                    </el-checkbox-group>
                  </span>
              </span>
            </el-tree>
          </div>
        </div>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
    data() {
        return {
          purTitle: "",   // 标题
          /*
          新增角色功能
          */
          RoleListVisible: false,   // 新增角色弹窗开关
          // 新增角色表单
          RoleListData: {
            name: '',
            roleType:3
          },
          // 新建角色表单验证
          RoleListRules: {
              name: [
                  { required: true, message: '请输入角色名称', trigger: 'blur' },
              ]
          },
          /*
          编辑角色功能
          */
          EditRoleListVisible: false,   // 编辑角色弹窗开关
          // 编辑角色表单
          EditRoleListData: {},
          // 编辑角色表单验证
          EditRoleListRules: {
              name: [
                  { required: true, message: '请输入角色名称', trigger: 'blur' },
              ]
          },
          rolelist:[],    // 获取角色列表
          currentIdx:null,
          currentCheckIdx:null,
          showIcon:true,
          treedata:[],    // 菜单列表
          defaultProps: {
              label: 'title'
          },
          roleID:null,
          newArr:[],      // 选中的菜单按钮
        };
    },
    created() {
      // let levle = sessionStorage.getItem('level');
      // // 企业路由刷新
      // if(levle == 1){
      //   if(this.$store.state.menus.length <= 0){
          
      //   }
      //   this.$store.dispatch("asynSetRoutes");
      // } 
      this.purTitle = this.$route.meta.title;   // 标题
      this.treedata = this.$store.state.menus;
      this.roleListAll();
    },
    mounted(){
      // this.$store.dispatch("asynSetMenus");
    },
    methods: {
      // 选中树形结构回调
      handleNodeClick(data, checked, indeterminate) {
        if(!data.children && checked && !indeterminate){
          if(data.import == 1){
            data.button.push('导入');
          }
          if(data.export == 1){
            data.button.push('导出');
          }
          if(data.add == 1){
            data.button.push('新增');
          }
          if(data.check == 1){
            data.button.push('查看');
          }
          if(data.edit == 1){
            data.button.push('编辑');
          }
          if(data.delete == 1){
            data.button.push('删除');
          }
          if(data.upload == 1){
            data.button.push('上传');
          }
          if(data.fileDelete == 1){
            data.button.push('文件删除');
          }
        }else{
          data.button = [];
        }
      },
      /*
      权限列表
      */
      // 保存角色权限
      AddDatasubmitForm(){
        let result=this.$refs.addtree.getCheckedNodes(false,true);
        let roleId = this.roleID;
        if(roleId == null){
          this.$message({
              type: "warning",
              message: "请给权限赋予角色!",
          });
        }else{
          let newArrs = [];
          let menus =  [];
          // console.log(result)
          result.map(function(value,index){
            // if(!value.children){
              let obj = {roleId:roleId,authoritys:[]};
              for(let i in value.button){
                switch (value.button[i]) {
                    case '导入':
                        obj.authoritys.push('import');
                        break;
                    case '导出':
                        obj.authoritys.push('export');
                        break;
                    case '新增':
                        obj.authoritys.push('add');
                        break;
                    case '查看':
                        obj.authoritys.push('check');
                        break;
                    case '编辑':
                        obj.authoritys.push('edit');
                        break; 
                    case '删除':
                        obj.authoritys.push('delete');
                        break; 
                    case '上传':
                        obj.authoritys.push('upload');
                        break; 
                    case '文件删除':
                        obj.authoritys.push('fileDelete');
                        break; 
                        
                }
              }
              obj.authoritys = [...new Set(obj.authoritys)];
              if(obj.authoritys.length != 0){
                obj.authoritys = obj.authoritys.toString();
                obj.menu = value.menu;
                obj.menuId = value.id;
                menus.push(obj);
                let newobjone = {
                  roleId:roleId,
                  menuId:value.id
                }
                newArrs.push(newobjone);
              }else{
                let newobj = {
                  roleId:roleId,
                  menuId:value.id
                }
                newArrs.push(newobj);
              } 
            // }
          })
          let role = {
            roleId:roleId,
            authRoles:menus,
            menuRoles:newArrs
          }
          this.api.roleList.saveEtcAuths(role)
          .then(res=>{
            if(res.data.code == 200){
                this.$message({
                    type: "success",
                    message: "设置角色权限成功!",
                    // duration:500,  
                    // onClose(){
                    //     this.roleListAll();          // 刷新
                    //     this.RoleListVisible = false;  
                    //     this.RoleListData = {};
                    // }
                });
            }
          })
        }
      },

      // 选中角色列表中的角色id
      currentCheckIdxclick(index,id){
        this.$store.dispatch("asynPtQySetMenus");
        this.treedata = this.$store.state.menus;
        this.currentCheckIdx = index;
        this.roleID = id;
        // let f = [{id:3}];
        // this.$refs.addtree.setCheckedNodes(f);
        // console.log('角色列表选中id为' + id);
        // 根据角色id回显权限列表
        this.api.roleList.selAuthyByRole({roleId:id})
        .then(res=>{
          if(res.data.code == 200){
            // console.log(res.data.data)
            let obj = {};
            let buttonArr = [];
            let newArr = [];
            for(let key in res.data.data.authRoles){
              newArr = [];
              buttonArr = [];
              let arr = res.data.data.authRoles[key].authoritys.split(',');
              for(let i in arr){
                  switch (arr[i]) {
                      case 'import':
                          buttonArr.push('导入')
                          break;
                      case 'export':
                          buttonArr.push('导出')
                          break;
                      case 'add':
                          buttonArr.push('新增')
                          break;
                      case 'check':
                          buttonArr.push('查看')
                          break;
                      case 'edit':
                          buttonArr.push('编辑')
                          break; 
                      case 'delete':
                          buttonArr.push('删除')
                          break; 
                      case 'upload':
                          buttonArr.push('上传')
                          break; 
                      case 'fileDelete':
                          buttonArr.push('文件删除')
                          break; 
                          
                  }
              }
              obj[res.data.data.authRoles[key].menuId] = buttonArr;
            }
            for(let kev in res.data.data.menuRoles){
              if(res.data.data.menuRoles[kev].isSon){
                let o = {
                  id:res.data.data.menuRoles[kev].id,
                  button:[]
                }
                if(obj[res.data.data.menuRoles[kev].id]){
                  o = {
                    id:res.data.data.menuRoles[kev].id,
                    button:obj[res.data.data.menuRoles[kev].id]
                  }
                }
                newArr.push(o);
              }
            }
            // console.log(newArr)
            this.newArr = [];
            this.newArr = newArr;
            this.$refs.addtree.setCheckedNodes(newArr);
            this.setmenudata(this.treedata);
          }
        })
      },
      setmenudata(treedata){
          for(var i in treedata){
              var data=this.newArr.filter(v=>{
                  return v.id == treedata[i].id;
              });
              if(data.length>0){
                  treedata[i].button=data[0].button;
              }
              if(treedata[i].children){
                  this.setmenudata(treedata[i].children)
              }
          }
      },
      enter(idx){
        this.currentIdx = idx;
      },
      leave(){
        this.currentIdx =null;
      },
      // 企业角色全部列表
      roleListAll(){
        this.api.roleList.AllQy({role_type:3})
        .then(res=>{
          if(res.data.code == 200){
            this.rolelist = res.data.data;
          }
        })
      },
      /*
      新增角色功能
      */
      // 取消新增角色
      RoleListVisibleclick(){
        this.RoleListVisible = false;
        this.RoleListData.name = "";
      },
      RoleListClose(){
        this.RoleListData.name = "";
        this.RoleListVisible = false;
      },
      // 新增角色提交表单 
      RoleListsubmitForm(formName) {
          this.$refs[formName].validate((valid) => {
              if (valid) {
                  this.api.roleList.addQy(this.RoleListData)
                  .then(res=>{
                      if(res.data.code == 200){
                          let that = this
                          that.$message({
                              type: "success",
                              message: "新增企业角色成功!",
                              duration:500,  
                              onClose(){
                                  that.roleListAll();          // 刷新角色列表
                                  that.RoleListVisible = false;  
                                  that.RoleListData = {};
                              }
                          });
                      }
                  })
              }
          });
      },

      /*
      编辑角色功能
      */
      EditRoleListVisibles(id){
        this.EditRoleListVisible = true;
        this.api.roleList.oneQy({id:id})
        .then(res=>{
          if(res.data.code == 200){
            this.EditRoleListData = res.data.data;
          }
        })
      },
      // 取消编辑角色
      EditRoleListVisibleclick(){
        this.EditRoleListVisible = false;
        this.EditRoleListData = {};
      },
      EditRoleListClose(){
        this.EditRoleListData = {};
        this.EditRoleListVisible = false;
      },
      // 编辑角色提交表单 
      EditRoleListsubmitForm(formName) {
          this.$refs[formName].validate((valid) => {
              if (valid) {
                  this.api.roleList.editQy(this.EditRoleListData)
                  .then(res=>{
                      if(res.data.code == 200){
                          let that = this
                          that.$message({
                              type: "success",
                              message: "编辑角色成功!",
                              duration:500,  
                              onClose(){
                                  that.roleListAll();             // 刷新角色列表
                                  that.EditRoleListVisible = false;  
                              }
                          });
                      }
                  })
              }
          });
      },
      // 删除角色
      DelClick(id){
          this.$confirm('此操作将删除角色及菜单一切权限,请谨慎操作 是否继续?', '提示', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning'
          }).then(() => {
            this.api.roleList.delQy({id:id})
            .then(res=>{
              if(res.data.code == 200){
                this.$message({
                  type: 'success',
                  message: '删除角色成功!'
                });
                this.roleListAll();
              }
            })
          }).catch(() => {
              this.$message({
                  type: 'info',
                  message: '已取消角色删除'
              });          
          });
      },
    },
};
</script>
<style lang="less" scoped>
.custom-tree-node{
  span{
    line-height: 60px;
  }
}
// 默认颜色
.el-tree /deep/ .el-tree-node .el-tree-node__content{
  width: 100%;
  height: 60px;
  background: #fff;
  // margin: 1px 0;
  color: #333;
  padding: 0 0 0 30px;
  box-sizing: border-box;
}
// 选中颜色
.el-tree /deep/ .is-current>.el-tree-node__content{
  background: #FFECEC !important;
  color: #333 !important;
}
.el-tree /deep/ .el-tree-node__label {
  // color: aqua;
  display: inline-block;
  width: 100%;
}
// .el-tree /deep/ .is-current>.el-tree-node__content:after{
//   content: "\2022";
//   color: rgb(155, 218, 97);
//   font-size: 40px;
//   position: relative;
//   right: 0;
//   bottom: 0;
//   cursor: pointer;
// }
.el-tree /deep/ .el-tree-node__children .el-tree-node__content{
  width: 100%;
  height:60px;
  background:#fff !important;
}
.el-tree /deep/ .el-tree-node__children .is-current .el-tree-node__content {
  background: #F6F6F6 !important;
  color: #CF1724 !important;
}
.el-tree /deep/ .el-tree-node__children .el-tree-node__content{
  color: #333 !important;
}
// .el-tree /deep/ .el-tree-node__children .is-current .el-tree-node__content:after{
//   content: "\2022";
//   font-size: 40px;
//   position: relative;
//   right: 10px;
//   bottom: 0;
//   cursor: pointer;
//   color: #333 !important;
// }
.el-tree .el-tree-node__expand-icon{
  color: #fff !important;
}
.el-tree-node__expand-icon.is-leaf{
  color: transparent !important;
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
    .el-button {
      margin-left: 16px;
      float: right;
    }
  }
}
// 头部结束
// 内容开始
.pur-nav {
  width: calc(100% - 16px);
  height: calc(100vh - 126px);
  margin: 8px;
  box-sizing: border-box;
  display:flex;
  flex-wrap: nowrap;
  .nav-left{
    width: 17.5%;
    height: 100%;
    background: #fff;
    margin-right:0.5%;
    padding: 0 4px;
    box-sizing: border-box;
    .RoleList-title{
      width: 100%;
      height: 60px;
      border-bottom:1px solid #e4e4e4;
      padding: 0 15px;
      box-sizing: border-box;
      font-size: 14px;
      line-height: 60px;
      font-weight: 600;
      color: #333;
      .el-button{
        line-height: 55px;
      }
    }
    // 角色列表
    .RoleListBOX{
        width: 100%;
        height: calc(100% - 60px);
        padding: 20px 0;
        box-sizing: border-box;
        overflow-y: scroll;
      .RoleList{  
        li{
          width: 100%;
          padding: 0 20px;
          box-sizing: border-box;
          position: relative;
          cursor: pointer;
          span{
            font-size: 14px;
            line-height:48px;
            color: #333;
          }
          i{
            float: right;
            font-size: 16px;
            line-height:48px;
            margin: 0 2%;
            cursor: pointer;
          }
        }
      }
    }
    .itemBg{
      background-color:#FFECEC; 
      i{
        color:#CF1724;
      }
    }
    .checkedItemBg{
      background-color:#FFECEC; 
      i{
        color:#CF1724;
      }
    } 
  }
  // 角色权限
  .nav-right{
    width: 82%;
    height: 100%;
    padding: 0 4px;
    box-sizing: border-box;
    background: #fff;
    .clearfixstyle{
      width: 100%;
      height: 60px;
      background: #fff;
      border-bottom:1px solid #e4e4e4;
      padding: 0 25% 0 15px;
      box-sizing: border-box;
      font-size: 14px;
      line-height: 60px;
      font-weight: 600;
      color: #333;
    }
    .powerBox{
      width: 100%;
      height: calc(100% - 60px);
      font-size: 16px;
      color: #333;
      background: #fff;
      font-weight: 400;
      overflow: hidden;
      overflow-y: scroll;
      padding-right:10%;
      box-sizing: border-box;
      .el-tree{
        padding: 20px 0;
        box-sizing: border-box;
      }
    }
  }
}
// 内容结束
</style>