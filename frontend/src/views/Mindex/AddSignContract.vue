<template>
    <!-- 发起签署合同页面 -->
    <el-form :model="AddData" :rules="rules" ref="AddData" size="small" label-width="130px">
        <!-- 头部开始 -->
        <el-row class="pur-top">
        <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
        <el-col :span="21" class="pur-right">
            <el-button size="small" @click="AddDatasubmitForm('AddData')" type="primary">发起签署</el-button>
            <el-button size="small" style="margin-left: 24px"  @click="Addopen">取消</el-button>
        </el-col>
        </el-row>
        <!-- 头部结束 -->
        <!-- 内容开始 -->
        <div class="pur-nav">
            <!-- 表单 -->
            <div class="pur-table">
                <el-row>
                    <!-- 任务信息 -->
                    <el-col :span="24">
                        <div class="see-top">任务信息</div>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:#CF1724;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">任务主题</label>
                            </el-col>
                            <el-form-item prop="topicName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.topicName" placeholder="请输入任务主题名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5" >
                                <span style="color:#CF1724;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">合同文件</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7"> 
                                    <el-upload
                                    :limit="1"
                                    :action="http + '/saas-common/chainEqbFlow/uploadFile'"
                                    :file-list="enclosurefileList"
                                    :headers="header"
                                    :on-remove="fileRemove"
                                    :on-preview="openfile"
                                    :before-upload="beforedemoUpload"
                                    :on-success="resfile"
                                    style="float:left;">
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
                            <span style="color:#CF1724;margin-left:3.5%;">* </span>
                            <span style="color:#666;font-size:12px;line-height:40px;">以下信息将用于电子签章服务注册，请谨慎填写！</span>
                        </div>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">发起方</label>
                            </el-col>
                            <el-form-item prop="launchType">
                                <el-col :span="7">  
                                    <el-radio v-model="AddData.launchType" :label="2">企业</el-radio>               
                                    <el-radio v-model="AddData.launchType" :label="0">个人</el-radio>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row v-if="AddData.launchType == 2">
                            <el-col :span="2.5">
                                <span style="color:#CF1724;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">发起方企业</label>
                            </el-col>
                            <el-form-item prop="etcName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.etcName" placeholder="请输入企业名称" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <!-- <el-row v-if="AddData.launchType == 2">
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">税号</label>
                            </el-col>
                            <el-form-item prop="thirdPartyUserId">
                                <el-col :span="7">
                                    <el-input v-model="AddData.thirdPartyUserId" placeholder="请输入税号" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row> -->
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:#CF1724;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">发起方姓名</label>
                            </el-col>
                            <el-form-item prop="participants">
                                <el-col :span="7">
                                    <el-input v-model="AddData.participants" placeholder="发起方姓名" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:#CF1724;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">发起方手机号</label>
                            </el-col>
                            <el-form-item prop="phone">
                                <el-col :span="7">
                                    <el-input v-model="AddData.phone" placeholder="请输入手机号" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                    <!-- 签署方信息 -->
                    <el-col :span="24">
                        <div class="see-top">签署方信息</div>
                        <el-row>
                            <el-col :span="2.5" >
                                <span style="color:#CF1724;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">参与人</label>
                            </el-col>
                            <el-form-item prop="signRole">
                                <el-col :span="7">  
                                    <el-radio v-model="AddData.signRole" :label="2" @change="prologistics($event)">企业</el-radio>               
                                    <el-radio v-model="AddData.signRole" :label="0" @change="prologistics($event)">个人</el-radio>
                                </el-col>
                                <!-- <el-col :span="7">
                                    <el-select v-model="AddData.signRole" style="width:100%;" @change="prologistics($event)">
                                        <el-option
                                        v-for="item in SelectcontractName"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                        </el-option>
                                    </el-select>
                                </el-col> -->
                            </el-form-item>
                        </el-row>
                        <!-- <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">发起方</label>
                            </el-col>
                            <el-form-item prop="launchType">
                                <el-col :span="7">  
                                    <el-radio v-model="AddData.launchType" :label="2">企业</el-radio>               
                                    <el-radio v-model="AddData.launchType" :label="0">个人</el-radio>
                                </el-col>
                            </el-form-item>
                        </el-row> -->
                        <!-- 企业 -->
                        <div v-if="AddData.signRole == 2" >
                            <!-- <div style="width:15px;height:20px;"></div> -->
                            <el-row v-for="(domain,index) in AddData.domains" :key="index">
                                <el-col :span="2.5">
                                    <span style="color:#CF1724;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">{{domain.label}}</label>
                                </el-col>
                                <el-form-item :prop="'domains.' + index + '.value'" :rules="{ required: true, message: '请输入' + domain.message + domain.id, trigger: 'blur' }">
                                    <el-col :span="7">
                                        <el-input v-model="domain.value" :placeholder="'请输入' + domain.label"></el-input>
                                    </el-col>
                                    <el-col :span="7" v-if="domain.labelone">
                                        <div class="removeDomainoneBox" @click.prevent="removeDomain(domain)">
                                            <i class="el-icon-remove"></i>
                                            <span> 删除</span>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                        </div>
                        <!-- 个人 -->
                        <div v-if="AddData.signRole == 0">
                            <el-row v-for="(domain,key) in AddData.domainsone" :key="key">
                                <el-col :span="2.5">
                                    <span style="color:#CF1724;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">{{domain.label}}</label>
                                </el-col>
                                <el-form-item :prop="'domainsone.' + key + '.value'" :rules="{ required: true, message: '请输入' + domain.message  + domain.id, trigger: 'blur' }">
                                    <el-col :span="7">
                                        <el-input v-model="domain.value" :placeholder="'请输入' + domain.label"></el-input>  
                                    </el-col>
                                    <el-col :span="2" v-if="domain.labelone">
                                        <div class="removeDomainoneBox" @click.prevent="removeDomainone(domain)">
                                            <i class="el-icon-remove"></i>
                                            <span> 删除</span>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                        </div>
                        <div>
                            <el-row>
                                <el-form-item>
                                    <el-col :span="7" v-if="AddData.signRole == 2">
                                        <div class="zjcyfspan" @click="addDomain">
                                            <i class="el-icon-circle-plus"></i><span > 增加参与方</span>
                                        </div>
                                    </el-col>
                                    <el-col :span="7" v-if="AddData.signRole == 0">
                                        <div class="zjcyfspan" @click="addDomainone">
                                            <i class="el-icon-circle-plus"></i><span > 增加参与方</span>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                        </div>

                    </el-col>
                </el-row>
            </div>
        </div>
      <!-- 文件预览弹窗 -->
      <el-dialog
        title="文件预览"
        :visible.sync="ImgVisible"
        width="44%"
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
        <!-- 内容结束 -->
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
            },
            // 新建签署表单数据
            AddData: {
                participants:'',
                // thirdPartyUserId:'',
                etcName:'',
                phone:'',
                launchType:2,           // 发起方类型
                domains: [],            // 企业
                domainsone:[],          // 个人
                signRole:2,             // 参与方类型
                
                field:'',               // 签署文件的集合id 
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
                // this.AddData.thirdPartyUserId = res.data.data.creditCode;
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
            // 签署方序号只增不减:回退后新加的签署方会与保留的重号,提交时被合并成一个
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
        // 参与人切换功能
        prologistics(id){
                if(id == 2){
                    this.AddData.domainsone = [];   // 个人 
                }
                if(id == 0){
                    this.AddData.domains = [];              // 企业
                }
        },

        /* 
        新建签署信息表单提交功能
        */ 
        // 新建签署信息表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    if(this.AddData.field){
                        if(this.AddData.domains.length > 0 || this.AddData.domainsone.length > 0){
                            if(this.AddData.launchType == 2){
                                this.AddData.orgInfo = {
                                    name:this.AddData.etcName,
                                    // thirdPartyUserId:this.AddData.thirdPartyUserId
                                }
                                this.AddData.userInfo = {
                                    name:this.AddData.participants,
                                    // thirdPartyUserId:this.AddData.phone
                                }
                            }else if(this.AddData.launchType == 0){
                                this.AddData.userInfo = {
                                    name:this.AddData.participants,
                                    // thirdPartyUserId:this.AddData.phone
                                }
                            }
                            var arrays = [];
                            if(this.AddData.domains.length > 0){
                                let a = {};
                                for(let i of new Set(this.AddData.domains.map(d => d.id))){ // 只遍历还在的签署方序号
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
                                for(let i of new Set(this.AddData.domainsone.map(d => d.id))){
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
                            console.log(this.AddData)
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
                                            that.$router.push({name:"SignContract"})
                                        }
                                    });
                                }
                            })
                        }else{
                            this.$message({
                                type: "warning",
                                message: "签署方信息不能少于一个签署方!"
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
            sessionStorage.removeItem('queryAll');
            this.$router.push({name:"SignContract"})
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
            //     sessionStorage.removeItem('queryAll');
            //     this.$router.push({name:"SignContract"})
            // });
        },
        // 附件文件格式验证
        beforedemoUpload(file) {
            let str = file.name;
            let index = str.lastIndexOf("\.");
            str = str.substring(index + 1,str.length)
            let types = [
                "png",
                "PNG",
                "jpg",
                "JPG",
                "pdf",
                "PDF",
                "doc",
                "DOC",
                "docx",
                "DOCX",
                'jpeg',
                "JPEG",
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
            if(response.code == 200 && response.data && response.data.field){
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
            }else{
                // 拿不到附件号的文件不能留在列表里:界面显示已上传,提交时却取不到附件号
                fileList.splice(fileList.indexOf(file), 1);
                this.$message.error(response.msg || '上传文件失败!');
            }
        },
        // 点击附件文件列表下载
        openfile(res){
            let url = this.http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + res.response.data.fileName;
            window.open(url)
        },
        // 删除附件文件
        fileRemove(file,fileList){
            this.AddData.field = ''; // 附件删完后不能带附件号提交:列表已空,取文件名会抛错
            this.enclosurefileList = fileList; // 列表与界面一致:原先删除接口失败时保存仍会带上
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
// /deep/ .el-upload-list__item {
//     width: 100%;
//     // background: red;
//     // position: relative;
// }
// 增加 删除样式开始
.removeDomainoneBox{
    width: 60px;
    height: 32px;
    margin-left:4%;
    cursor: pointer;
    i{
        font-size: 16px;
        color: #CF1724;
        line-height: 32px;
        vertical-align: middle;
    }
    span{
        font-size: 14px;
        color: #CF1724;
        line-height: 32px;
        font-weight: 600;
        vertical-align: middle;
    }
}
.zjcyfspan{
    width: 116px;
    height: 32px;
    text-align: center;
    cursor: pointer;
    border: 1px solid #CF1724;
    box-sizing: border-box;
    border-radius:5px;
    
    i{
        font-size: 14px;
        color: #CF1724;
        line-height: 32px;
    }
    span{
        font-size: 14px;
        color: #CF1724;
        line-height: 32px;
        font-weight: 600;
    }
}
// 增加 删除样式结束

.Addbuttonbox{
    cursor: pointer;
    width:90px;
    height: 30px;
    background: #FFFFFF;
    border: 1px solid #E4E4E4;
    border-radius: 4px;
    font-size: 14px;
    color: #666;
    line-height: 30px;
    text-align: left;
    padding: 0 15px;
    box-sizing: border-box;
    position: absolute;
    top:60%;
    left:30%;
}
.Addbuttonbox:hover{
    background: rgb(250,232,233);
    border: 1px solid #CF1724;
    color: #CF1724;
}
// 文件图片查看样式
.imgBox{
    width: 100%;
    height: 600px;
    overflow: hidden;
    overflow-y: scroll;
    img{
        width: 100%;
        // height:100%;
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
    cursor: pointer;
    width:100%;
    height: 32px;
    background: #FFFFFF;
    border: 1px solid #E4E4E4;
    border-radius: 4px;
    font-size: 14px;
    color: #666;
    line-height: 30px;
    text-align: left;
    padding: 0 15px;
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