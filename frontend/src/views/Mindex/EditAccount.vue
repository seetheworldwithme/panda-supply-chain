<template>
    <!-- 新建账户信息页面 -->
    <el-form :model="AddData" :rules="rules" ref="AddData" size="small" label-width="130px">
        <!-- 头部开始 -->
        <el-row class="pur-top">
        <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
        <el-col :span="21" class="pur-right">
            <el-button size="small" @click="AddDatasubmitForm('AddData')">保存</el-button>
            <el-button size="small" style="margin-left: 24px"  @click="Addopen">取消</el-button>
        </el-col>
        </el-row>
        <!-- 头部结束 -->
        <!-- 内容开始 -->
        <div class="pur-nav">
            <!-- 表单 -->
            <div class="pur-table">
                <el-row>
                    <!-- 基础信息 -->
                    <el-col :span="24">
                        <div class="see-top"><div></div>基本信息</div>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">账户类型</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-select v-model="AddData.ownerType" placeholder="请选择账户类型" style="width:100%;" @change="okAcout($event)">
                                        <el-option
                                        v-for="item in SelectcontractName"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                        </el-option>
                                    </el-select>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">{{titleAccount}}</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.accountName" :placeholder="'请输入' + titleAccount"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">{{titlenumber}}</label>
                            </el-col>
                            <el-form-item prop="accountNumber">
                                <el-col :span="7">
                                    <el-input v-model="AddData.accountNumber" :placeholder="'请输入' + titlenumber"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                </el-row>
            </div>
        </div>
        <!-- 内容结束 -->
    </el-form>
</template>
<script>
export default {
    data() {
        return {
            purTitle: "", // 标题
            titleAccount:'账户名称',
            titlenumber:'账户号码',
            // 表单验证
            rules: {
                accountNumber: [
                    { pattern:/^[1-9]\d*$/, message: '请输入正确的账号', trigger: 'blur'}
                ]
            },
            // 账户信息选择下拉框
            SelectcontractName:[
                { 
                    value: 1,
                    label: '公司账户'
                },
                { 
                    value: 2,
                    label: '对方账户'
                }
            ],  
            AddData: {
                ownerType:1,
                accountName: "",   
                accountNumber:""
            }, // 表单数据
        };
    },
    created() {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        // 接收主页单条数据
        this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));
    },
    mounted () {},
    methods: {
        //新建账户信息表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.api.Account.edit(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    sessionStorage.removeItem('queryAll');
                                    that.$router.push({name:"Account"})
                                }
                            });
                        }
                    })
                }
            });
        },
        // 取消保存表单
        Addopen() {
            sessionStorage.removeItem('queryAll');
            this.$router.push({name:"Account"})
        },
        okAcout(id){
            if(id == 1){
                this.titleAccount = "账户名称";
                this.titlenumber = "账户号码";
            }else if(id == 2){
                this.titleAccount = "对方账户";
                this.titlenumber = "对方账号";
            }
        }
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
    margin: 0 60px;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    padding:33px;
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
            margin-bottom:24px;
            div{
                width: 2px;
                height:15px;
                margin:15px 6px 0 0;
                background-color: #EA222E;
                float: left;
            }
        }
        .see-lable {
            padding: 15px 0;
            .el-input{
                .input{
                    width: 430px;
                }
            }
        }
    }
}
.AddPurIcon{
    width: 16px;
    height: 16px;
    border-radius: 50%;
    border: 1px solid #828282;
    color: #828282;
    font-size:12px;
    text-align: center;
    line-height: 16px;
    margin-top: 14px;
    margin-left: -25px;
}
.AddbuttonUpload{
    width: 100px;
    height: 32px;
    background: #FFFFFF;
    border: 1px solid #E4E4E4;
    border-radius: 4px;
    font-size: 14px;
    color: #666;
    line-height: 32px;
    text-align: center;
}
.AddbuttonUpload:hover{
    background: #e4e4e4;
    color: #333;
}
// 表单内容结束
</style>