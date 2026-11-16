<template>
  <!-- 合同模板页面 -->
    <div>
        <!-- 头部开始 -->
        <el-row class="pur-top">
          <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
          <el-col :span="21" class="pur-right">
              <el-button type="info" size="small" @click="AddDatasubmitForm" v-if="$route.meta.add">保存</el-button>
          </el-col>
        </el-row>
        <!-- 头部结束 -->
        <!-- 内容开始 -->
        <div class="pur-nav">
          <div class="con-top">
              <b>提示 ：</b>
              <p>1.默认字段只有部分可选择是否启用，是否必填;</p>
              <p>2.字段的排序位置为模块表单的先后顺序位置;</p>
              <p>3.鼠标按住， 即可以调整排列顺序;</p>
              <p>4.单行输入在100字以内、多行输入需要在2000字以内。</p>
          </div>
          <el-button type="primary" class="con-button" v-if="$route.meta.add">+ 新建模板</el-button>
          <!-- 表格数据 -->
          <el-table
            :data="tableData"
            :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
            style="width: 100%; background-color: #f6f6f6 !important"
            @selection-change="purChange">
              <el-table-column type="selection"></el-table-column>
              <el-table-column prop="id" v-if="show">
              </el-table-column>
              <el-table-column prop="goodsName" label="模板名称" sortable>
              </el-table-column>
              <el-table-column prop="settlementTime" label="状态" sortable width="200">
              </el-table-column>
              <el-table-column label="操作">
                <template slot-scope="scope" style="color:#FB8138;">
                  <el-button type="text" size="small" v-if="$route.meta.edit">编辑</el-button>
                  <el-button @click="handleClick(scope.row)" type="text" size="mini" v-if="$route.meta.check">查看</el-button>
                  <el-button type="text" size="small">复制</el-button>
                  <el-button type="text" size="small" v-if="$route.meta.delete">删除</el-button>
                </template>
              </el-table-column>
              <el-table-column prop="unitPrice" label="排序" sortable>
              </el-table-column>
          
          </el-table>
        </div>
        <!-- 内容结束 -->
    </div>
</template>
<script>
export default {
  data() {
    return {
      show:false,
      purTitle: "", // 标题
      multipleSelection: [], // 表格选中数据
      tableData: [
            {
                id:1,
                goodsName: "采购合同",
                settlementTime: "生效",
                unitPrice: "",    
            },
            {
                id:2,
                goodsName: "销售合同",
                settlementTime: "未生效",
                unitPrice: "",
            },
            {
                id:2,
                goodsName: "供应链合同",
                settlementTime: "未生效",
                unitPrice: "",
            },
      ],  // 表格数据
    };
  },
  created() {
    this.purTitle = this.$route.meta.title;
  },
  methods: {
    AddDatasubmitForm(){
      // console.log("保存合同模板");
    },
    //表格
    purChange(val) {
      // console.log(val); // 表格选中数据
      this.multipleSelection = val;
    },
    //按钮事件
    handleClick(data){
      // console.log(data)
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
    .el-button {
        margin-left: 16px;
        float: right;
    }
  }
}
// 头部结束
// 内容开始
.pur-nav {
  // width: calc(100% - 110px);
  height: calc(100vh - 126px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 24px 32px;
  box-sizing: border-box;
  overflow-y: scroll;
  .con-top{
    width: 100%;
    padding: 16px 24px;
    background-color:#F2F2F2;
    box-sizing: border-box;
    b{
      display: inline-block;
      font-size: 14px;
      margin-bottom:20px;
      color: #333;
    }
    p{
      font-size: 14px;
      line-height: 25px;
      color: #666;
    }
  }
  .con-button{
    display: block;
    margin: 24px 0;
  }
}
// 内容结束
</style>