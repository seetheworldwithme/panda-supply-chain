<template>
    <!-- 新建账户信息页面 -->   
    <el-form :model="AddData" :rules="rules" ref="AddData" size="small" label-width="130px">
        <!-- 头部开始 -->
        <el-row class="pur-top">
        <el-col :span="3">
            <span class="pur-size">{{ purTitle }}</span>
        </el-col>
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
                        <div class="see-top"><div>
                        </div>基本信息</div>
                        <el-row style="margin-bottom:20px;">
                            <el-col :span="12">
                                <p style="color:rgb(102, 102, 102);">计算规则：不含税金额 = 销项含税金额/1.13， 销项税额 = 不含税金额*0.13， 增值税额 = 销项税额 - 进项税额，
                                    城建税 = 增值税额*0.07， 教育费附加税额 = 增值税额*0.03， 地方教育费附加 = 增值税额*0.02， 印花税额 = 增值税额*0.03，
                                     缴税合计 = 增值税 + 城建税 + 教育费附加税额 + 地方教育费附加 + 印花税额 + 企业所得税 + 滞纳金</p>
                            </el-col>
                            <el-col :span="10"></el-col>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">所属月份</label>
                            </el-col> 
                            <el-form-item prop="month">
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.month"
                                    type="month"
                                    placeholder="请选择所属月份"
                                    value-format="yyyy-MM"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">销项含税金额</label>
                            </el-col>
                            <el-form-item prop="includMoney">
                                <el-col :span="7">
                                    <el-input @input="inputinclud" v-model="AddData.includMoney" placeholder="请输入销项含税金额（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">不含税金额</label>
                            </el-col>
                            <el-form-item prop="baseMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.baseMoney" placeholder="请输入不含税金额（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">销项税额</label>
                            </el-col>
                            <el-form-item prop="taxAmount">
                                <el-col :span="7">
                                    <el-input v-model="AddData.taxAmount" placeholder="请输入销项税额（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>

                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">进项税额</label>
                            </el-col>
                            <el-form-item prop="jxMoney">
                                <el-col :span="7">
                                    <el-input @input="inputjxMoney" v-model="AddData.jxMoney" placeholder="请输入进项税额（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">增值税额</label>
                            </el-col>
                            <el-form-item prop="zzMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.zzMoney" placeholder="请输入增值税额（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">城建税</label>
                            </el-col>
                            <el-form-item prop="cjMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.cjMoney" placeholder="请输入城建税（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">教育费附加税额</label>
                            </el-col>
                            <el-form-item prop="jyMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.jyMoney" placeholder="请输入教育费附加税额（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">地方教育费附加</label>
                            </el-col>
                            <el-form-item prop="dfjyMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.dfjyMoney" placeholder="请输入地方教育费附加（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">印花税额</label>
                            </el-col>
                            <el-form-item prop="yhMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.yhMoney" placeholder="请输入印花税额（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">企业所得税</label>
                            </el-col>
                            <el-form-item prop="qyMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.qyMoney" placeholder="请输入企业所得税（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">滞纳金</label>
                            </el-col>
                            <el-form-item prop="znjMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.znjMoney" placeholder="请输入滞纳金（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">缴税合计</label>
                            </el-col>
                            <el-form-item prop="sumMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.sumMoney" placeholder="请输入缴税合计（单位：元）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">备注</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.remark" placeholder="请输入备注"></el-input>
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
            // 表单验证
            rules: {
                sumMoney:[
                    { required: true, message: '请输入缴税合计金额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                znjMoney:[
                    { required: true, message: '请输入滞纳金', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                month:[
                    { required: true, message: '请输入所属月份', trigger: 'blur' },
                ],
                includMoney:[
                    { required: true, message: '请输入销项含税金额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                baseMoney:[
                    { required: true, message: '请输入不含税金额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                taxAmount:[
                    { required: true, message: '请输入销项税额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                jxMoney:[
                    { required: true, message: '请输入进项税额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                zzMoney:[
                    { required: true, message: '请输入增值税额', trigger: 'blur' }
                ],
                cjMoney:[
                    { required: true, message: '请输入城建税', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                jyMoney:[
                    { required: true, message: '请输入教育费附加税额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                dfjyMoney:[
                    { required: true, message: '请输入地方教育费附加', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                yhMoney:[
                    { required: true, message: '请输入印花税额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                qyMoney:[
                    { required: true, message: '请输入企业所得税', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
            }, 
            AddData: {
                includMoney:0,
                baseMoney:0,
                taxAmount:0,
                jxMoney:0,
                zzMoney:0,
                cjMoney:0,
                jyMoney:0,
                dfjyMoney:0,
                yhMoney:0,
                qyMoney:0,
                znjMoney:0,
                sumMoney:0,
            }, // 表单数据
        };
    },
    created() {
        // 模块标题
        this.purTitle = this.$route.meta.title;
    },
    mounted () {},
    methods: {
        //新建账户信息表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.api.billRate.add(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    that.$router.push({name:"Coststatistics"})
                                }
                            });
                        }
                    })
                }
            });
        },
        // 取消保存表单
        Addopen() {
            this.$router.push({name:"Coststatistics"})
        },
        inputinclud(){
            this.AddData.baseMoney = (Number(this.AddData.includMoney)/1.13).toFixed(2);
            this.AddData.taxAmount = (Number(this.AddData.baseMoney)*0.13).toFixed(2);
            this.inputjxMoney(); // 销项税额变了,增值税与附加税要跟着重算
        },
        inputjxMoney(){
            let val = this.AddData;
                this.AddData.zzMoney = (Number(val.taxAmount) - Number(val.jxMoney))>0?(Number(val.taxAmount) - Number(val.jxMoney)).toFixed(2):0;
                if(this.AddData.zzMoney != 0){
                    this.AddData.cjMoney = (Number(val.zzMoney)*0.07).toFixed(2);
                    this.AddData.jyMoney = (Number(val.zzMoney)*0.03).toFixed(2);
                    this.AddData.dfjyMoney = (Number(val.zzMoney)*0.02).toFixed(2);
                    this.AddData.yhMoney = (Number(val.zzMoney)*0.0003).toFixed(2);
                }else{
                    this.AddData.cjMoney = 0;
                    this.AddData.jyMoney = 0;
                    this.AddData.dfjyMoney = 0;
                    this.AddData.yhMoney = 0;
                }
        },
    },
    watch:{
        AddData:{
            handler(val){
                this.AddData.sumMoney = (val.zzMoney*1 + val.cjMoney*1 + val.jyMoney*1 + val.dfjyMoney*1 + val.yhMoney*1 + val.qyMoney*1 + val.znjMoney*1).toFixed(2);
            },
            deep:true
        }
    }
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