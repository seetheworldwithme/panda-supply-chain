<template>
  <!-- 邀请码生成页面 -->
    <div>
        <!-- 头部开始 -->
        <el-row class="pur-top">
        <el-col :span="3"
            ><span class="pur-size">{{ purTitle }}</span></el-col
        >
        </el-row>
        <!-- 头部结束 -->
        <!-- 邀请码生成页面 -->
        <el-form
            :model="AddData"
            ref="AddData"
            size="small"
            label-width="130px">
            <!-- 内容开始 -->
            <div class="pur-nav">
                <!-- 表单 -->
                <div class="pur-table">
                    <el-row>
                        <el-col :span="24">
                          <el-row>
                              <el-col :span="2">
                              <!-- <span style="color:red;">* </span> -->
                              <label
                                  style="
                                  color: #666;
                                  font-size: 14px;
                                  line-height: 30px;
                                  padding-left: 12px;
                                  "
                                  >生效时间</label
                              >
                              </el-col>
                              <el-form-item>
                              <el-col :span="7">
                                  <el-select
                                  v-model="AddData.effSeconds"
                                  placeholder="请选择生效时间"
                                  style="width: 100%"
                                  >
                                  <el-option
                                      v-for="item in SelectcontractName"
                                      :key="item.value"
                                      :label="item.label"
                                      :value="item.value"
                                  >
                                  </el-option>
                                  </el-select>
                              </el-col>
                              </el-form-item>
                          </el-row>
                          <el-row>
                              <el-col :span="2">
                              <!-- <span style="color: red">* </span> -->
                              <label style="color: #666; font-size: 14px; line-height: 30px;padding-left: 12px;">邀请码生成</label>
                              </el-col>
                              <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.enterpriseName">
                                      <span slot="suffix" @click="invCode" style="margin-right:16px;color:#CF1724;cursor: pointer;">随机生成</span>
                                    </el-input>
                                </el-col>
                              </el-form-item>
                          </el-row>

                        </el-col>
                    </el-row>
                </div>
            </div>
            <!-- 内容结束 -->
        </el-form>
    </div>
</template>
<script>
export default {
  data() {
    return {
      purTitle: "", // 标题
      SelectcontractName: [
        {
          value:60,
          label: "1分钟",
        },
        {
          value: 300,
          label: "5分钟",
        },
        {
          value:600,
          label: "10分钟",
        },
        {
          value:1800,
          label: "半小时",
        },
      ], // 时效选择下拉框
      AddData: {
        enterpriseName: "",
        effSeconds:"",  // 有效时间
        }, // 表单数据
    };
  },
  created() {
    this.purTitle = this.$route.meta.title;
  },
  methods: {
    invCode(){
      if(this.AddData.effSeconds != ''){
        this.api.authAll.createCode({effSeconds:this.AddData.effSeconds})
        .then(res=>{
          if(res.data.code == 200){
            this.AddData.enterpriseName = res.data.data;
          }
        })
      }else{
        this.$message({
          type: 'warning',
          message: "请先选择邀请码生效时间"
        });
      }
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
}
// 头部结束
// 表单内容开始
.pur-nav {
  height: calc(100vh - 126px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 33px;
  box-sizing: border-box;
  .pur-table {
    width: 100%;
    height: 100%;
    overflow-y: scroll;
    .see-top {
      height: 40px;
      font-weight: 600;
      font-size: 14px;
      line-height: 45px;
      margin-bottom: 24px;
      div {
        width: 2px;
        height: 15px;
        margin: 15px 6px 0 0;
        background-color: #ea222e;
        float: left;
      }
    }
    .see-lable {
      padding: 15px 0;
      .el-input {
        .input {
          width: 430px;
        }
      }
    }
  }
}
.AddPurIcon {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1px solid #828282;
  color: #828282;
  font-size: 12px;
  text-align: center;
  line-height: 16px;
  margin-top: 14px;
  margin-left: -25px;
}
.AddbuttonUpload {
  width: 100px;
  height: 32px;
  background: #ffffff;
  border: 1px solid #e4e4e4;
  border-radius: 4px;
  font-size: 14px;
  color: #666;
  line-height: 32px;
  text-align: center;
}
.AddbuttonUpload:hover {
  background: #e4e4e4;
  color: #333;
}
// 表单内容结束
</style>