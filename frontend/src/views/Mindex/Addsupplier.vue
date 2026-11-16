<template>
    <!-- 新建供应商档案页面 -->
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
                        <div class="see-top"><div></div>基础信息</div>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">企业名称</label>
                            </el-col>
                            <el-form-item prop="enterpriseName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.enterpriseName" placeholder="请输入企业全称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">企业性质</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-select v-model="AddData.enterpriseNature" placeholder="请选择企业性质" style="width:100%;">
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
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;" placeholder="请输入统一信用代码">统一信用代码</label>
                            </el-col>
                            <el-form-item prop="creditCode">
                                <el-col :span="7">
                                    <el-input v-model="AddData.creditCode"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;" placeholder="请输入经营地址">经营地址</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.businessAddress"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                    <!-- 生产经营信息 -->
                    <el-col :span="24">
                        <div class="see-top"><div></div>生产经营信息</div>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">煤矿</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.coalMine" placeholder="请输入煤矿名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">年产量</label>
                            </el-col>
                            <el-form-item prop="annualOutput">
                                <el-col :span="7">
                                    <el-input v-model="AddData.annualOutput" placeholder="吨"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">煤场</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.coalYard" placeholder="请输入煤场名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">产品大类选择</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="2" v-for="(v,k) in AddproBigType" :key="k">
                                <el-radio v-model="AddData.pro1" :label="v.dvalue" @change="getValue()">{{v.dkey}}</el-radio>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">产品小类选择</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-select filterable v-model="AddData.pro2" placeholder="请选择产品小类" style="width:100%;" @change="proSmaclick($event)">
                                        <el-option
                                        v-for="item in SelectproSmaType"
                                        :key="item.dvalue"
                                        :label="item.dkey"
                                        :value="item.dvalue"
                                        >
                                        </el-option>
                                    </el-select>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row v-if="pIDshow">
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">产品名称</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.proRemark" placeholder="请输入产品名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">标的指标</label>
                            </el-col>
                            <el-form-item>
                                    <el-col :span="2" style="margin-top:5px;">
                                        <el-input v-model="AddData.sulfur" style="width:64px;"></el-input>
                                        <p style="margin-right:10px;">硫分St.d%</p>
                                    </el-col>
                                    <el-col :span="2" style="margin-top:5px;">
                                        <el-input v-model="AddData.volatilization" style="width:64px;"></el-input>
                                        <p style="margin-right:10px;">挥发分Vdaf%</p>
                                    </el-col>
                                    <el-col :span="2" style="margin-top:5px;">
                                        <el-input v-model="AddData.water" style="width:64px;"></el-input>
                                        <p style="margin-right:10px;">水分Mt%</p>
                                    </el-col>
                                    <!-- <el-col :span="2" v-for="(v, index) in AddData.list" :key="index" style="margin-top:5px;">
                                         <el-input v-model="v.dValue" style="width:64px;"></el-input>
                                        <p style="margin-right:10px;">{{v.dkey}}</p>
                                    </el-col> -->
                                    <el-col :span="1">
                                        <el-tooltip effect="dark" content="标的指标需要在基础数据里设置" placement="bottom">
                                            <div class="AddPurIcon">?</div>
                                        </el-tooltip>
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
            /* 
            新建供应商档案数据
            */
            pIDshow:false,  // 产品名称显示隐藏
            purTitle: "", // 标题
            AddproBigType: [], // 产品大类选择
            SelectproSmaType: [],  // 小类选择下拉框
            SelectcontractName:[
                { 
                    value: '1',
                    label: '国企'
                },
                { 
                    value: '2',
                    label: '民营企业'
                },
                { 
                    value: '3',
                    label: '央企'
                },
                { 
                    value: '4',
                    label: '外企'
                }
            ],  // 企业性质选择下拉框
            // 表单验证
            rules: {
                enterpriseName: [
                    { required: true, message: '请输入企业名称', trigger: 'change' }
                ],
                annualOutput: [
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '请输入正确的数量格式', trigger: 'blur'}
                ],
                creditCode: [
                    { pattern:/(^(?:(?![IOZSV])[\dA-Z]){2}\d{6}(?:(?![IOZSV])[\dA-Z]){10}$)|(^\d{15}$)/, message: '请输入正确的统一信用代码', trigger: 'blur'}
                ],
            },
            Extended:'',  // 扩展信息菜单
            AddData: {
                enterpriseName: "",
                enterpriseNature: "",
                creditCode: "",
                businessAddress: "",
                coalMine: "",
                annualOutput: "",
                coalYard: "",
                pro1: "",
                pro2: "",
                // 标的指标
                sulfur: "",
                volatilization: "",
                water: "",
                list:[],    
            }, // 表单数据
        };
    },
    created() {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        //基础维护数据
        this.api.basics.all({template_id:1})
        .then(res=>{
            if(res.data.code == 200){
                this.Extended = res.data.data;
                this.AddproBigType = res.data.data.pro_type;    // 产品大类选择
            }
        })
    },
    mounted () {},
    methods: {
        /* 
        新建供应商档案表单提交功能
        */ 
        // 大类目前的选中值
        getValue(){
            this.SelectproSmaType = [];
            var pro__middle_typeID = ""; 
                // 获取大类id
                for(let i=0;i<this.Extended.pro_type.length;i++){
                    this.SelectproSmaType = [];
                    if(this.Extended.pro_type[i].dvalue == this.AddData.pro1){
                        pro__middle_typeID = this.Extended.pro_type[i].dictId  // 获取产品大类id
                    }
                }
                // 获取小类数据
                for(let k=0;k<this.Extended.pro__middle_type.length;k++){
                    var proArrs = this.Extended.pro__middle_type[k].dparentId.split(',');
                    for(let i=0;i<proArrs.length;i++){
                        if(proArrs[i] == ","){
                            this.splice(proArrs[i]);  
                        }
                        if(proArrs[i].indexOf(pro__middle_typeID) != -1){
                            this.AddData.pro2 = '';  // 产品小类为空
                            this.SelectproSmaType.push(this.Extended.pro__middle_type[k]); // 产品小类选择
                        }
                    }
                }
                // console.log(pro__middle_typeID)
                // console.log(this.SelectproSmaType)
                // console.log(this.Extended.pro__middle_type)
        },
        // 小类目前的选中值
        proSmaclick(pID){
                // 选中小类其他显示产品名称
                if(pID == 20){
                    this.pIDshow = true;
                }else{
                    this.pIDshow = false;   
                }
                var proid = "";                     // 下拉框选中value对应的id
                if(this.AddData.pro2 == ''){
                    this.AddData.list = [];
                }else{
                    // 获取小类id
                    for(let i=0;i<this.Extended.pro__middle_type.length;i++){
                        if(this.Extended.pro__middle_type[i].dvalue == pID){
                            this.AddData.list = [];
                            proid = this.Extended.pro__middle_type[i].dictId;
                        }
                    }
                    // console.log(proid)
                    // // 获取标的指标数据
                    // for(let k=0;k<this.Extended.pro__small_type.length;k++){
                    //     var proArr = this.Extended.pro__small_type[k].dparentId.split(',');
                    //     for(let i=0;i<proArr.length;i++){
                    //         if(proArr[i] == ","){
                    //             this.splice(proArr[i]);  
                    //         }
                    //         if(proArr[i].indexOf(proid) != -1 && this.Extended.pro__small_type[k].dextend == true){
                    //             this.AddData.list = [];
                    //             // 标的指标扩展字段
                    //             this.AddData.list.push(
                    //                 {
                    //                     dictId:this.Extended.pro__small_type[k].id,     // 字典id
                    //                     dkey:this.Extended.pro__small_type[k].dkey,     // lab名称
                    //                     dvalue:'',                                      // name名称
                    //                     dtype:'pro__small_type',                        // 扩展类型
                    //                 }
                    //             );         
                    //         }
                    //     }
                    // }
                }
        },
        //新建供应商档案表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.AddData.partnerType = 1;
                    this.api.supplier.add(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    that.$router.push({name:"Partner",query:{sel: 1}})
                                }
                            });
                        }
                    })
                }
            });
        },
        // 取消保存表单
        Addopen() {
            this.$router.push({name:"Partner",query:{sel: 1}})
            // this.$confirm("检测到表单已有内容填写，是否保存为草稿?", "提示", {
            //     confirmButtonText: "确定",
            //     cancelButtonText: "取消",
            //     type: "warning",
            // })
            // .then(() => {
            //     let that = this
            //     that.$message({
            //         type: "success",
            //         message: "保存草稿成功!",
            //         duration:1000,  
            //         onClose(){
            //             that.$router.push({name:"Partner",query:{sel: 1}})
            //         }
            //     });
            // })
            // .catch(() => {
            //     this.$router.push({name:"Partner",query:{sel: 1}})
            // });
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
            // border-left: 2px solid #EA222E;
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