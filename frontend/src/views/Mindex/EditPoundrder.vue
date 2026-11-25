<template>
    <!--新建磅单页面 -->
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
                        <div class="see-top">基本信息</div>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">票号</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.billNumber" placeholder="请输入票号"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">车牌号</label>
                            </el-col>
                            <el-form-item prop="licensePlate">
                                <el-col :span="7">
                                    <el-input v-model="AddData.licensePlate" placeholder="请输入车牌号"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">煤种</label>
                            </el-col>
                            <el-form-item prop="coalType">
                                <el-col :span="7">
                                    <el-input v-model="AddData.coalType" placeholder="请输入煤种"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">毛重</label>
                            </el-col>
                            <el-form-item prop="roughWeight">
                                <el-col :span="7">
                                    <el-input @input="inputinclud" v-model="AddData.roughWeight" placeholder="请输入毛重(吨)"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">皮重</label>
                            </el-col>
                            <el-form-item prop="tare">
                                <el-col :span="7">
                                    <el-input @input="inputinclud" v-model="AddData.tare" placeholder="请输入皮重(吨)"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">净重</label>
                            </el-col>
                            <el-form-item prop="netWeight">
                                <el-col :span="7">
                                    <el-input v-model="AddData.netWeight" placeholder="请输入净重(吨)"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">开票单位</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.issuingOffice" placeholder="请输入开票单位"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">原煤供煤方</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.rawCoalSupplier" placeholder="请输入原煤供煤方"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>

                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收套打单位</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.sstdCompany" placeholder="请输入实收套打单位"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">进厂时间</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.inTime"
                                    type="datetime"
                                    placeholder="请选择进厂时间"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收进厂时间</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.ssInTime"
                                    type="datetime"
                                    placeholder="请选择实收进厂时间"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">合同号</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.contractNo" placeholder="请输入合同号"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">购买方</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.purchaser" placeholder="请输入购买方"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收皮重</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.ssTare" placeholder="请输入实收皮重(吨)"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">皮重时间</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.tareTime"
                                    type="datetime"
                                    placeholder="请选择皮重时间"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收毛重</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.ssRoughWeight" placeholder="请输入实收毛重(吨)"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">毛重时间</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.roughWeightTime"
                                    type="datetime"
                                    placeholder="请选择毛重时间"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">装煤司机</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.coalLoadingDriver" placeholder="请输入装煤司机"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">装煤时间</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.coalLoadingTime"
                                    type="datetime"
                                    placeholder="请选择装煤时间"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收净重</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.ssNetWeight" placeholder="请输入实收净重(吨)"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">开票时间</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.invoicingTime"
                                    type="datetime"
                                    placeholder="请选择开票时间"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收套打时间</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.sstdTime"
                                    type="datetime"
                                    placeholder="请选择实收套打时间"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">上传时间</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.uploadTime"
                                    type="datetime"
                                    placeholder="请选择上传时间"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收状态</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.ssStatus" placeholder="请输入实收状态"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">打印状态</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.printStatus" placeholder="请输入打印状态"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">时间间隔</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.timeIntervalMinute" placeholder="请输入时间间隔（分钟）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">是否补录</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.addOrNot" placeholder="请输入是否补录"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">回皮人</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.repatriator" placeholder="请输入回皮人"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">称重人</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.weigher" placeholder="请输入称重人"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收回皮人</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.ssRepatriator" placeholder="请输入实收回皮人"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收称重人</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.ssWeigher" placeholder="请输入实收称重人"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">作废原因</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.reasonCancellation" placeholder="请输入作废原因"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">补录原因</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.reasonSupplementary" placeholder="请输入补录原因"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                    <!-- 附件信息 -->
                    <el-col :span="24">
                        <div class="see-top">附件信息</div>
                        <!-- 上传磅单附件 -->
                        <el-row>
                            <el-col :span="2" >
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">上传磅单附件</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7"> 
                                    <el-upload
                                    :action="http + '/saas-common/upload/batch'"
                                    :file-list="enclosurefileList"
                                    :headers="header"
                                    :on-remove="fileRemove"
                                    :on-preview="openfile"
                                    :before-upload="beforedemoUpload"
                                    :on-success="resfile">
                                        <div class="AddbuttonUpload"><span>请选择附件</span></div>
                                    </el-upload>
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
            新建销项发票数据
            */  
            http:'',
            header:{},                      // 上传请求头
            fileObj:{},                     // 附件文件格式验证
            enclosurefileList: [],          // 附件文件列表
            purTitle: "", // 标题
            // 表单验证
            rules: {
                licensePlate: [
                    { required: true, message: '请输入车牌号', trigger: 'change' }
                ],
                coalType: [
                    { required: true, message: '请输入煤种', trigger: 'change' }
                ],
                tare: [
                    { required: true, message: '请输入皮重', trigger: 'blur' }
                ],
                roughWeight: [
                    { required: true, message: '请输入毛重', trigger: 'blur' }
                ],
                netWeight: [
                    { required: true, message: '请输入净重', trigger: 'blur' }
                ],
            },
            // 新建表单数据
            AddData: {
                filelist:[],  // 磅单文件的集合id 
                roughWeight:0,
                tare:0,
            },
        };
    },
    created() {},
    mounted () {
        
        // 模块标题
        this.purTitle = this.$route.meta.title;
        // 接收主页单条数据
        this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));
        if(!this.AddData){
            // 直接输入网址 / 后退进入时没有待编辑的磅单:回列表页,否则下面取字段就抛错
            return this.$router.replace({name:'Poundrder'});
        }
        // 文件上传请求头
        this.header = {Authorization:sessionStorage.getItem('Authorization')}
        this.http  = this.$store.state.http;                 // IP地址
        // 回显文件列表
        this.api.file.one({relationId:JSON.parse(sessionStorage.getItem('queryAll')).id})
        .then(res=>{
            var newarr = [];
            var newobj;
            for(let i=0;i<res.data.data.length;i++){
                newobj = {
                    name:res.data.data[i].name,
                    response:{
                        code:200,
                        data:[res.data.data[i]]
                    }
                }
                newarr.push(newobj);
            }
            for(let i=0;i<newarr.length;i++){
                // 附件文件内容
                this.enclosurefileList.push(newarr[i]); 
            }
        })
    },
    methods: {  
        inputinclud(){
            if(this.AddData.roughWeight*1 - this.AddData.tare*1 < 0){
                this.AddData.netWeight = 0;
            }else{
                this.AddData.netWeight = (this.AddData.roughWeight*1 - this.AddData.tare*1).toFixed(2)
            }
        },
        //编辑磅单表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    // 文件上传
                    var arr = [];
                    this.enclosurefileList.map((itme,index)=>{
                        arr.push({id:itme.response.data[0].id});
                    }).join(',')
                    this.AddData.filelist = arr;
                    this.api.Weighing.edit(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                    this.delRemovedFiles();
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存磅单成功!",
                                duration:500,  
                                onClose(){
                                    that.$router.push({name:"Poundrder"})
                                }
                            });
                        }
                    })
                }
            });
        },
        // 取消保存表单
        Addopen() {
            this.$router.push({name:"Poundrder"})
        },
        // 附件文件格式验证
        beforedemoUpload(file) {
            let str = file.name;
            let index = str.lastIndexOf("\.");
            str = str.substring(index + 1,str.length)
            let types = [
                "png", 
                "jpg",
                "pdf",
                "doc",
                "docx",
                "xlsx",
                "xls"
            ];
            const status = types.indexOf(str) != -1;
            this.fileObj[file.name] = true;
            if(!status) {
                this.fileObj[file.name] = false;
                this.$message.error("上传文件只能是 PNG JPG PDF DOCX DOC XLSX XLS 格式,请重新上传!");
            }
            return status;
        },
        // 附件上传文件成功
        resfile(response, file, fileList){
            // HTTP 200 + code≠200 或 data 为空(例如后端转 pdf 失败)时拿不到附件号,文件不能留在列表里
            if(response.code == 200 && response.data && response.data.length){
                this.$message({
                    type: 'success',
                    message: '上传文件成功!'
                });
                this.enclosurefileList = fileList;
            }else{
                fileList.splice(fileList.indexOf(file), 1);
                this.$message.error(response.msg || '上传文件失败!');
            }
        },
        // 点击附件文件列表下载
        openfile(res){
            let url = this.http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + res.response.data[0].fileName;
            window.open(url)
        },
        // 删除附件文件
        // 保存成功后才删编辑期间移除的回显附件(原先一移除就删,取消编辑后记录指向已删文件)
        delRemovedFiles(){
            (this.removedFileIds || []).forEach(id => this.api.addfile.del({id}));
            this.removedFileIds = [];
        },
        fileRemove(file,fileList){
            this.enclosurefileList = fileList; // 列表与界面一致:原先删除接口失败时提交仍会带上
            if(!this.fileObj[file.name]){
                // 回显的附件没有经过 before-upload:先记下,等保存成功再删
                this.removedFileIds = (this.removedFileIds || []).concat(file.response.data[0].id);
                return;
            }
            return  this.api.addfile.del({id:file.response.data[0].id})
                        .then(res=>{
                            if(res.data.code == 200){
                                this.$message({
                                    type: 'warning',
                                    message: '删除文件成功!'
                                });
                                this.enclosurefileList = fileList;
                            }
                        })
        },
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
.smbox{
    width:100%;
    min-height:78px;
    border: 1px solid #DCDFE6;
    border-radius: 4px;
    padding:8px;
    box-sizing: border-box;
}
.details{
   width:100%;
   height:32px;
   border: 1px solid #DCDFE6;
   border-radius: 4px;
   .details-nav{
       background-color:#F5F7FA;
        width: calc(100% - 73px);
        height: 32px;
        float: left;
        overflow: hidden;
        .details-tag{
            width: 44%;
            height: 32px;
            display: inline-block;
            .tag-box{
                width: 100%;
                display:inline-block;
                height: 32px;
                color: #C4C4C4;
                .el-tag{
                    color:#333;
                    width: 80%;
                    overflow: hidden;           // 盒子溢出隐藏
                    text-overflow:ellipsis;     // 文字溢出显示省略号
                    white-space: nowrap;        // 文字不换行
                }
                b{
                    width: 20%;
                    display: inline-block;
                }
            }
        }
        .details-conent{
            width:56%;
            text-align: right;
            height: 32px;
            float: right;
            padding-right: 4px;
            box-sizing: border-box;
            overflow: hidden;           // 盒子溢出隐藏
            text-overflow:ellipsis;     // 文字溢出显示省略号
            white-space: nowrap;        // 文字不换行
        }
   }  
   .details-size{
       background-color: #F5F7FA;
        border-left: 1px solid #DCDFE6;
        height:100%;
        text-align: center;
        width:73px;
        box-sizing: border-box;
        color:#333;
        cursor: pointer;
        float: right;
        span{
            width: 100%;
            height: 32px;
            font-size: 13px;
        }
   }
}
.pur-nav {
    height: calc(100vh - 110px);
    background-color: #fff;
    margin: 0 60px;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    padding:33px 0 33px 33px;
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
            border-bottom: 1px dashed #999999;
            margin-bottom:24px;
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
// 关联合同弹窗样式开始
.el-dialog__body{
    padding: 0;
}
.Wdialog{
    width: 100%;
    border: 1px solid #E0E0E0;
    margin-top: 16px;
}
// 选项卡
.concentBox {
    width: 100%; 
    position: relative;
}
.textbutton{
    width:56px;
    height: 32px;
    position: absolute;
    right:25px;
    bottom: 10px;
}
.totalTab {
    width: 9%;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    color: #333;
    text-align: center;
    margin:9px 0 0 35px;
}
.textlast{
    width: 50%;
    height: 46px;
    padding: 9px 17px 5px 24px;
    float: right;
    span{
        font-size: 16px;
        line-height: 32px;
        color: #EA222E;
        font-weight: 600;
        cursor: pointer;
        float: right;
        margin-right:24px;
    }
    .el-input{
        float: right;
    }
}
.totalTab.active {
    color: #333;
}
.stripnone{
    height: 3px;
    background-color: #ea222e;
    margin:7px auto 16px;
    display: none;
}
.strip {
    display: block;
}
.text {
  width: 100%;
}
  // 分页
.pur-pages {
    width: 100%;
    height: 32px;
    margin: 148px 0 9px 9px;
    .el-pagination {
      float: left;
    }
    .el-button{
        float: right;
        margin-right: 24px;
    }
}
  // 新建表单
.CompanyData{
    width: 100%;
    padding: 8px 16px;
    position: absolute;
    left: 0;
    bottom:115px;
}
// 关联合同弹窗样式结束
</style>