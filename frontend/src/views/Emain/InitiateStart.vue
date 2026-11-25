<template>
    <!-- 小程序发起签署合同页面 -->
    <el-form class="pur-nav" :model="AddData" :rules="rules" ref="AddData" size="small" label-width="120px">
        <!-- 表单 -->
        <div class="pur-table">
            <el-row>
                <!-- 任务信息 -->
                <el-col :span="24">
                    <div class="see-top">任务信息</div>
                    <el-row>
                        <el-col :span="8">
                            <span style="color:#CF1724;">* </span>
                            <label style="color:#666;font-size:14px;line-height:30px;">任务主题</label>
                        </el-col>
                        <el-form-item prop="topicName">
                            <el-col :span="24">
                                <el-input v-model="AddData.topicName" placeholder="请输入任务主题名称"></el-input>
                            </el-col>
                        </el-form-item>
                        
                    </el-row>
                    <el-row>
                        <el-col :span="8" >
                            <span style="color:#CF1724;">* </span>
                            <label style="color:#666;font-size:14px;line-height:30px;">合同文件</label>
                        </el-col>
                        <el-form-item>
                            <el-col :span="24"> 
                                <el-upload
                                :limit="1"
                                :action="http + '/saas-common/chainEqbFlow/uploadFile'"
                                :file-list="enclosurefileList"
                                :headers="header"
                                :on-remove="fileRemove"
                                :on-preview="openfile"
                                :before-upload="beforedemoUpload"
                                :on-success="resfile"
                                >
                                    <div class="AddbuttonUpload"><span>添加签署文件</span></div>
                                </el-upload>
                                <div class="Addbuttonbox" v-if="fileShow" @click="scopImg"><span>在线预览</span></div>
                            </el-col>
                        </el-form-item>
                    </el-row>
                </el-col>
                <!-- 发起方信息 -->
                <el-col :span="24">
                    <div class="see-top">发起方信息
                        <span style="color:#CF1724;margin-left:0.5%;"> * </span>
                        <span style="color:#666;font-size:10px;line-height:40px;">以下信息用于签署发起，请谨慎填写！</span>
                    </div>
                    <el-row>
                        <el-col :span="8">
                            <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">发起方</label>
                        </el-col>
                        <el-form-item prop="launchType">
                            <el-col :span="24">                 
                                <el-radio v-model="AddData.launchType" :label="0">个人</el-radio>
                                <el-radio v-model="AddData.launchType" :label="2">企业</el-radio>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row v-if="AddData.launchType == 2">
                        <el-col :span="8">
                            <span style="color:#CF1724;">* </span>
                            <label style="color:#666;font-size:14px;line-height:30px;">发起方企业</label>
                        </el-col>
                        <el-form-item prop="etcName">
                            <el-col :span="24">
                                <el-input v-model="AddData.etcName" placeholder="请输入企业名称" disabled></el-input>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="8">
                            <span style="color:#CF1724;">* </span>
                            <label style="color:#666;font-size:14px;line-height:30px;">发起方姓名</label>
                        </el-col>
                        <el-form-item prop="participants">
                            <el-col :span="24">
                                <el-input v-model="AddData.participants" placeholder="发起方姓名" disabled></el-input>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="8">
                            <span style="color:#CF1724;">* </span>
                            <label style="color:#666;font-size:14px;line-height:30px;">发起方手机号</label>
                        </el-col>
                        <el-form-item prop="phone">
                            <el-col :span="24">
                                <el-input v-model="AddData.phone" placeholder="请输入手机号" disabled></el-input>
                            </el-col>
                        </el-form-item>
                    </el-row>
                </el-col>
                <!-- 签署方信息 -->
                <el-col :span="24">
                    <div class="see-top">签署方信息</div>
                    <el-row>
                        <el-col :span="9" >
                            <span style="color:#CF1724;">* </span>
                            <label style="color:#666;font-size:14px;line-height:30px;">参与人</label>
                        </el-col>
                        <el-form-item prop="signRole">
                            <el-col :span="10">
                                <el-select v-model="AddData.signRole" style="width:100%;" @change="prologistics($event)">
                                    <el-option
                                    v-for="item in SelectcontractName"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value">
                                    </el-option>
                                </el-select>
                            </el-col>
                            <el-col :span="12" v-if="logistics">
                                <span style="margin-left:10%;color:#CF1724;cursor: pointer;" @click="addDomain">+ 增加参与方</span>
                            </el-col>
                            <el-col :span="12" v-if="bliLogistics">
                                <span style="margin-left:10%;color:#CF1724;cursor: pointer;" @click="addDomainone">+ 增加参与方</span>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <!-- 企业 -->
                    <div v-if="logistics" >
                        <!-- <div style="width:15px;height:20px;"></div> -->
                        <el-row v-for="(domain,index) in AddData.domains" :key="index">
                            <el-col :span="9">
                                <span style="color:#CF1724;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">{{domain.label}}</label>
                            </el-col>
                            <el-form-item :prop="'domains.' + index + '.value'" :rules="{ required: true, message: '请输入' + domain.message + domain.id, trigger: 'blur' }">
                                <el-col :span="21">
                                    <el-input v-model="domain.value" :placeholder="'请输入' + domain.label" style="width:100%;"></el-input>
                                </el-col>
                                <el-col :span="1" v-if="domain.labelone">
                                    <i class="el-icon-delete" @click.prevent="removeDomain(domain)" style="margin-left:20%;cursor: pointer;font-size: 20px;color: #CF1724;line-height: 32px;"></i>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </div>
                    <!-- 个人 -->
                    <div v-if="bliLogistics">
                        <el-row v-for="(domain,key) in AddData.domainsone" :key="key">
                            <el-col :span="9">
                                <span style="color:#CF1724;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">{{domain.label}}</label>
                            </el-col>
                            <el-form-item :prop="'domainsone.' + key + '.value'" :rules="{ required: true, message: '请输入' + domain.message  + domain.id, trigger: 'blur' }">
                                <el-col :span="21">
                                    <el-input v-model="domain.value" :placeholder="'请输入' + domain.label"></el-input>  
                                </el-col>
                                <el-col :span="1" v-if="domain.labelone">
                                    <i class="el-icon-delete" @click.prevent="removeDomainone(domain)" style="margin-left:20%;cursor: pointer;font-size: 20px;color: #CF1724;line-height: 32px;"></i>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </div>
                </el-col>
            </el-row>
        </div>
        <div class="pur-right">
            <el-button style="width:30%;"  @click="Addopen">取 消</el-button>
            <el-button @click="AddDatasubmitForm('AddData')" type="primary" style="width:65%;float:right;">发 起 签 署</el-button>
            
        </div>
        <!-- 文件预览弹窗 -->
        <el-dialog
            title="文件预览"
            :visible.sync="ImgVisible"
            width="95%"
            :before-close="ImgClose">
                <div class="imgBox" v-if="itemshow == true" v-loading="loadingimgss" element-loading-text="文件加载中">
                <!-- <iframe :src="url" frameborder="0" width="800px" height="600px"></iframe> -->
                <div v-for="i in numPages" :key="i">
                    <pdf :src="src" :page="i"></pdf>
                </div>
                </div>  
                <viewer class="imgBox" :images="inmages" v-else v-loading="loadingimg" element-loading-text="文件加载中"> 
                <!-- // photo 一定要一个数组，否则报错 -->
                    <img
                    v-for="(src,index) in inmages"
                    :src="src.filePath"
                    :key="index">
                </viewer>
        </el-dialog>
    </el-form>
</template>
<script>
import pdf from 'vue-pdf'
import CMapReaderFactory from 'vue-pdf/src/CMapReaderFactory.js' // 加载中文的包
export default {
    components: {pdf},
    data() {
        return {
            /* 
            新建订单信息数据
            */
            purTitle: "",                           // 标题
            http:'',                                // 文件上传地址
            header:{},                              // 上传请求头
            fileObj:{},                             // 文件格式验证
            enclosurefileList: [],                  // 文件展示列表
            logistics:true,                         // 企业默认显示
            bliLogistics:false,                     // 个人默认隐藏
            // 参与人选择下拉框
            SelectcontractName:[
                { 
                    value: 0,
                    label: '个人'
                },
                { 
                    value: 2,
                    label: '企业'
                }
            ],  
            numbers:0,
            numbersone:0,
            // 表单验证
            rules: {
                systemCode: [
                    { required: true, message: '请输入系统编码', trigger: 'change' }
                ],
                // signName: [
                //     { required: true, message: '请输入姓名', trigger: 'change' }
                // ],
                // signPhone: [
                //     { required: true, message: '请输入手机号', trigger: 'change' }
                // ],
                // signEtcName: [
                //     { required: true, message: '请输入企业名称', trigger: 'change' }
                // ],
                billNumber: [
                    { pattern: /(?=.*\d)(?=.*\W+)(?=.*[A-Z])(?!.*\n).*$/, message: '订单编号格式不正确', trigger: 'blur'}
                ],
                tonnagePaid: [
                    { required: true, message: '请输入实收吨数', trigger: 'blur' },
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '格式不正确', trigger: 'blur'}
                ],
            },
            // 新建签署表单数据
            AddData: {
                participants:'',
                thirdPartyUserId:'',
                etcName:'',
                phone:'',
                launchType:2,          // 发起方类型
                domains: [],
                domainsone:[],
                signRole:2,         // 参与人
                field:'',           // 签署文件的集合id 
                singers:[]
            }, 
            // 在线预览
            fileShow:false,
            ImgVisible:false,
            faileData:'',
            inmages:[],
            numPages:'',
            src:'',
            path:'',//pdf的地址，例如：/testFile.pdf
            geshi:'',
            itemshow:null,
            loadingimg:true,
            loadingimgss:true,
        };
    },
    created() {
        this.api.eqbflows.selLanucnInfo()
        .then(res=>{
            if(res.data.code == 200){
                this.AddData.etcName = res.data.data.companyName;
                this.AddData.thirdPartyUserId = res.data.data.creditCode;
                this.AddData.participants = res.data.data.username;
                this.AddData.phone = res.data.data.phone;
            }
        })
    },
    mounted () {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        // 文件上传请求头
        this.header = {Authorization:sessionStorage.getItem('Authorization')};
        this.http  = this.$store.state.http;                 // IP地址
        // this.AddData.domains.push(
        //     {
        //         id:this.numbers,
        //         signRole:2,
        //         label:'姓名',
        //         value: '',
        //         labelone:'删除',
        //         type:'signName',
        //         message:'姓名'
        //     },{
        //         id:this.numbers,
        //         signRole:2,
        //         label:'手机号',
        //         value: '',
        //         labelone:'',
        //         type:'signPhone',
        //         message:'手机号'
        //     },{
        //         id:this.numbers,
        //         signRole:2,
        //         label:'企业名称',
        //         value: '',
        //         labelone:'',
        //         type:'signEtcName',
        //         message:'企业名称'
        //     }
        // );
        // this.numbers = this.AddData.domains.length/3;
        // this.numbersone = this.AddData.domainsone.length/2;
        // this.AddData.domainsone.push(
        //     {
        //         id:this.numbersone,
        //         signRole:0,
        //         label:'姓名',
        //         value: '',
        //         labelone:'删除',
        //         type:'signName',
        //         message:'姓名'
        //     },{
        //         id:this.numbersone,
        //         signRole:0,
        //         label:'手机号',
        //         value: '',
        //         labelone:'',
        //         type:'signPhone',
        //         message:'手机号'
        //     }
        // );
        // this.addDomain();
    },
    methods: {
        scopImg(){
            this.ImgVisible = true;
            this.geshi = this.faileData.fileName.substring(this.faileData.fileName.length - 4);
            let arrs = ['.pdf','.PDF','docx','DOCX','.doc','.DOC'];
            if(!arrs.includes(this.geshi)){
              this.itemshow = false;
              this.api.file.preFIle({file_id:this.faileData.id})
              .then(res=>{
                  if(res.data.code == 200){
                    this.loadingimg = false;
                    this.inmages = res.data.data;
                  }
              })
            }else{
              this.itemshow = true;
              this.src = pdf.createLoadingTask({
                url:this.http + '/saas-common/upload/previewStream?fileName=' + this.faileData.fileName + '&&AuthorizationId=wxda0b1c468b83bc99&&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f' ,     // 正式环境用这个！
                CMapReaderFactory
              })
              // 让所有页数一次性加载完，否则就只会加载第一页
              this.src.promise.then(pdf => {
                  this.loadingimgss = false;
                  this.numPages = pdf.numPages;
              }).catch(() => {})
            }
        },
        // 关闭查看图片弹窗
        ImgClose() {
            this.ImgVisible = false;
            this.src = '';
            this.inmages = [];
            this.numPages = '';
            this.path = '';//pdf的地址，例如：/testFile.pdf
            this.geshi = '';
            this.loadingimg = true;
            this.loadingimgss = true;
        },

        // 企业动态增加input
        removeDomain(item) {
            this.numbers--
            var index = this.AddData.domains.indexOf(item);
            if (index !== -1) {
            this.AddData.domains.splice(index,3);
            }
        },
        addDomain() {
            this.numbers++;
            this.AddData.domains.push(
                {
                    id:this.numbers,
                    signRole:2,
                    label:'签署方姓名' + this.numbers,
                    value: '',
                    labelone:'删除',
                    type:'signName',
                    message:'签署方姓名'
                },{
                    id:this.numbers,
                    signRole:2,
                    label:'签署方手机号' + this.numbers,
                    value: '',
                    labelone:'',
                    type:'signPhone',
                    message:'签署方手机号'
                },{
                    id:this.numbers,
                    signRole:2,
                    label:'签署方企业名称' + this.numbers,
                    value: '',
                    labelone:'',
                    type:'signEtcName',
                    message:'签署方企业名称'
                }
            );
        },
        // 个人动态增加input
        removeDomainone(item) {
            this.numbersone--;
            var index = this.AddData.domainsone.indexOf(item);
            if (index !== -1) {
            this.AddData.domainsone.splice(index,2);
            }
        },
        addDomainone() {
            // this.numbersone = this.numbersone - 1
            this.numbersone++;
            this.AddData.domainsone.push(
                {
                    id:this.numbersone,
                    signRole:0,
                    label:'签署方姓名' + this.numbersone,
                    value: '',
                    labelone:'删除',
                    type:'signName',
                    message:'签署方姓名'
                },{
                    id:this.numbersone,
                    signRole:0,
                    label:'签署方手机号' + this.numbersone,
                    value: '',
                    labelone:'',
                    type:'signPhone',
                    message:'签署方手机号'
                }
            );
        },

        /* 
        新建签署信息表单提交功能
        */ 
       // 参与人切换功能
       prologistics(id){
            if(id == 2){
                this.logistics = true;          // 企业显示
                this.bliLogistics = false;      // 个人隐藏
            }
            if(id == 0){
                this.logistics = false;         // 企业隐藏
                this.bliLogistics = true;       // 个人显示
            }
       },
        // 新建签署信息表单提交
        AddDatasubmitForm(formName) {
            const now = Date.now();
            if (this._submitLock && now - this._submitLock < 3000) return; // 防重:双击会发起两条签署流程
            this._submitLock = now;
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    if(this.AddData.field){
                        if(this.AddData.domains.length > 0 || this.AddData.domainsone.length > 0){
                            if(this.AddData.launchType == 2){
                                this.AddData.orgInfo = {
                                    name:this.AddData.etcName,
                                    thirdPartyUserId:this.AddData.thirdPartyUserId
                                }
                                this.AddData.userInfo = {
                                    name:this.AddData.participants,
                                    thirdPartyUserId:this.AddData.phone
                                }
                            }else if(this.AddData.launchType == 0){
                                this.AddData.userInfo = {
                                    name:this.AddData.participants,
                                    thirdPartyUserId:this.AddData.phone
                                }
                            }
                            var arrays = [];
                            if(this.AddData.domains.length > 0){
                                let a = {};
                                for(let i=1;i<=this.numbers;i++){
                                    a = {}
                                    for(let key in this.AddData.domains){
                                        if(this.AddData.domains[key].type == 'signName'){
                                            this.AddData.domains[key].signName = this.AddData.domains[key].value
                                        }
                                        if(this.AddData.domains[key].type == 'signPhone'){
                                            this.AddData.domains[key].signPhone = this.AddData.domains[key].value
                                        }
                                        if(this.AddData.domains[key].type == 'signEtcName'){
                                            this.AddData.domains[key].signEtcName = this.AddData.domains[key].value
                                        }
                                        if(i === this.AddData.domains[key].id){
                                            Object.assign(a,this.AddData.domains[key]);
                                        }
                                    }
                                    arrays.push(a);
                                }
                            }
                            if(this.AddData.domainsone.length > 0){
                                let b ={};
                                for(let i=1;i<=this.numbersone;i++){
                                    b = {}
                                    for(let key in this.AddData.domainsone){
                                        if(this.AddData.domainsone[key].type == 'signName'){
                                            this.AddData.domainsone[key].signName = this.AddData.domainsone[key].value
                                        }
                                        if(this.AddData.domainsone[key].type == 'signPhone'){
                                            this.AddData.domainsone[key].signPhone = this.AddData.domainsone[key].value
                                        }
                                        if(i === this.AddData.domainsone[key].id){
                                            Object.assign(b,this.AddData.domainsone[key]);
                                        }
                                    }
                                    arrays.push(b)
                                }
                            }
                            this.AddData.singers = arrays;
                            this.AddData.contractName = this.enclosurefileList[0].name;
                            
                            let newobj = {
                                signEtcName:this.AddData.etcName,
                                signPhone:this.AddData.phone,
                                signName:this.AddData.participants,
                                signRole:this.AddData.launchType,
                            }
                            this.AddData.singers.push(newobj);
                            // console.log(this.AddData)
                            this.api.eqbflows.launchFlow(this.AddData)
                            .then(res=>{
                                if(res.data.code == 200){
                                    var that = this
                                    that.$message({
                                        type: "success",
                                        message: "发起签署合同成功!",
                                        duration:500,  
                                        onClose(){
                                            // sessionStorage.removeItem('queryAll');
                                            that.$router.push({name:"Sponsor"})
                                        }
                                    });
                                }
                            })
                        }else{
                            this.$message({
                                type: "warning",
                                message: "签署方信息不能少于两个签署方!"
                            });
                        }
                    }else{
                        this.$message({
                            type: "warning",
                            message: "上传文件不能为空!"
                        });
                    }
                }
            });
        },
        // 取消保存表单
        Addopen() {
            // sessionStorage.removeItem('queryAll');
            this.$router.push({name:"Sponsor"});
            // this.$confirm("检测到表单已有内容填写，是否保存为草稿?", "提示", {
            //     confirmButtonText: "确定",
            //     cancelButtonText: "取消",
            //     type: "warning",
            // })
            // .then(() => {
            //     var that = this
            //     that.$message({
            //         type: "success",
            //         message: "保存草稿成功!",
            //         duration:1000,  
            //         onClose(){
            //             sessionStorage.removeItem('queryAll');
            //             that.$router.push({name:"SignContract"})
            //         }
            //     });
            // })
            // .catch(() => {

            // });
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
                'jpeg',
            ];
            const status = types.indexOf(str) != -1;
            this.fileObj[file.name] = true;
            if(!status) {
                this.fileObj[file.name] = false;
                this.$message.error("上传文件只能是 PNG JPG PDF DOCX DOC JPEG 格式,请重新上传!");
            }
            return status;
        },
        // 附件上传文件成功
        resfile(response, file, fileList){
            if(response.code == 200){
                this.fileShow = true;
                this.$message({
                    type: 'success',
                    message: '上传文件成功!'
                });
                this.faileData = response.data;
                this.enclosurefileList = fileList;
                // 文件上传
                // console.log(this.enclosurefileList)
                var arr = [];
                this.enclosurefileList.map((itme,index)=>{
                    arr.push({id:itme.response.data.field});
                }).join(',');
                this.AddData.field = arr[0].id;
            }
        },
        // 点击附件文件列表下载
        openfile(res){
            let url = this.http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + res.response.data.fileName;
            window.open(url)
        },
        // 删除附件文件
        fileRemove(file,fileList){
            if(this.fileObj[file.name]){
                return  this.api.addfile.del({id:file.response.data.id})
                        .then(res=>{
                            if(res.data.code == 200){
                                this.$message({
                                    type: 'warning',
                                    message: '删除文件成功!'
                                });
                                this.fileShow = false;
                                this.enclosurefileList = fileList;
                            }
                        })
            }
        },
    },
    computed: {},
    watch: {}
};
</script>
<style lang="less" scoped>
// 在线预览样式开始
.Addbuttonbox{
    cursor: pointer;
    width:74px;
    height: 30px;
    background: #FFFFFF;
    border: 1px solid #E4E4E4;
    border-radius: 4px;
    font-size: 14px;
    color: #666;
    line-height: 30px;
    text-align: left;
    padding: 0 6px;
    box-sizing: border-box;
    position: absolute;
    top:0;
    right:0;
}
.Addbuttonbox:hover{
    background: rgb(250,232,233);
    border: 1px solid #CF1724;
    color: #CF1724;
}
// 在线预览样式结束
// 文件图片查看样式开始
.imgBox{
    width: 100%;
    height: 400px;
    overflow: hidden;
    overflow-y: scroll;
    img{
        width: 100%;
    }
}
// 文件图片查看样式结束
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
            width: 50%;
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
            width:50%;
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
    width: 100%;
    min-height: 100vh;
    background-color: #F6F6F6;
    padding:14px;
    box-sizing: border-box;
    .pur-table {
        width: 100%;
        min-height:calc(100vh - 102px);
        border-radius:4px;
        box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
        background: #fff;
        padding:0 14px 0;
        // margin-top: 14px;
        box-sizing: border-box;
        // overflow-y: scroll;
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
    .pur-right{
        width: 100%;
        height: 60px;
        padding-top:20px;
        box-sizing: border-box;
        .el-button{
            float: left;
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
    cursor: pointer;
    width:100%;
    height: 30px;
    background: #FFFFFF;
    border: 1px solid #E4E4E4;
    border-radius: 4px;
    font-size: 14px;
    color: #666;
    line-height: 30px;
    text-align: left;
    padding: 0 8px;
    box-sizing: border-box;
}
.AddbuttonUpload:hover{
    background: rgb(250,232,233);
    border: 1px solid #CF1724;
    color: #CF1724;
}
// 表单内容结束
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
</style>