<template>
  <!-- 上传文件页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur" v-if="$route.meta.fileDelete">文件删除</el-button>
        <el-button size="small" @click="Importdailog" v-if="$route.meta.upload">上传</el-button>
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch"
        >
        </el-input>
      </el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav"
      v-loading="loading"
      element-loading-text = "加载中...">
      <!-- 表格 -->
      <div class="pur-table">
        <el-table
          height="100%"
          border
          :data="tableData"
          :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
          style="width: 100%;"
          @selection-change="purChange">
            <el-table-column type="selection" fixed="left"> </el-table-column>
            <el-table-column type="index" label="序号" width="50"></el-table-column>
            <el-table-column label="文件名称" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span style="cursor: pointer;color:#CF1724;" v-if="scope.row.fileRealName" @click="scopImg(scope.row)">{{scope.row.fileRealName}}</span>
              </template>
            </el-table-column>
            <el-table-column prop="fileDescription" label="文件备注" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="createTime" label="创建时间" :show-overflow-tooltip="true">
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以上传文件试试</p>
                <el-button size="small" @click="Importdailog">上传</el-button>
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
        @current-change='getpurcon'
        @size-change='getpurcon'>
        </el-pagination>
      </div>
      <!-- 上传弹窗功能 -->
      <el-dialog
        title="新建文件上传"
        :close-on-click-modal="false"
        :visible.sync="ImportVisible"
        width="30%"
        :before-close="ImporthandleClose">
            <el-form :model="AddData" ref="AddData" size="small" :rules="rules">
                <el-row>
                  <el-col :span="4">
                      <span style="color:red;">* </span>
                      <label style="color:#666;font-size:14px;line-height:40px;">上传文件</label>
                  </el-col> 
                  <el-col :span="18">
                    <el-form-item>
                    <el-upload
                    ref="upload"
                    :action="http + '/saas-common/upload/uploadFiles'"
                    :file-list="enclosurefileList"
                    :headers="header"
                    :on-remove="fileRemove"
                    :before-upload="beforedemoUpload"
                    :on-success="resfile">
                        <div class="AddbuttonUpload"><span>请选择附件</span></div>
                    </el-upload>
                    </el-form-item>
                  </el-col>  
                </el-row>
                <el-row>
                  <el-col :span="4">
                      <span style="color:red;">* </span>
                      <label style="color:#666;font-size:14px;line-height:40px;">文件备注</label>
                  </el-col> 
                  <el-col :span="18">
                    <el-form-item prop="fileDescription">
                      <el-input v-model="AddData.fileDescription" placeholder="请输入文件备注"></el-input>
                    </el-form-item>
                  </el-col>
                </el-row>
            </el-form>
            <div slot="footer">
              <!-- <el-button @click="Importble" size="small">取 消</el-button> -->
              <el-button type="primary" @click="AddDatasubmitForm('AddData')" size="small">确 定</el-button>
            </div>
      </el-dialog>
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
            <div class="imgbuttons">
              <el-button type="primary" icon="el-icon-printer" size="mini" @click="downPrint">打印</el-button>
              <el-button type="primary" icon="el-icon-download" style="float:right;" size="mini" @click="openfile">下载</el-button>
            </div>
      </el-dialog>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
import pdf from 'vue-pdf'
import CMapReaderFactory from 'vue-pdf/src/CMapReaderFactory.js' // 加载中文的包
export default {
  components: {pdf},
    data() {
        return {
            purSearch: "",          // 搜索  
            http:'',
            dataAva:false,
            loading:false,
            current:1,              // 当前页数
            size:50,                // 显示条数
            total:0,                // 总条数
            delmul:[],              // 批量删除
            show:false,
            purTitle: "",           // 标题
            multipleSelection: [],  // 选中数据
            tableData: [],          // 表格数据

            ImportVisible:false,                    // 上传打开弹窗
            header:{},                              // 上传请求头
            fileObj:{},                             // 文件格式验证
            enclosurefileList: [],                  // 文件展示列表
            // 上传表单数据
            AddData: {}, 
            rules: {
              region: [
                { required: true, message: '', trigger: 'change' }
              ],
              fileDescription: [
                { required: true, message: '请填写文件备注', trigger: 'blur' }
              ],
            },

            // 图片预览
            ImgVisible:false,
            inmages:[],
            numPages:'',
            src:'',
            path:'',//pdf的地址，例如：/testFile.pdf
            geshi:'',
            itemshow:null,
            loadingimg:true,
            loadingimgss:true,
            oneData:'',     // 单条数据
        };
    },
    created() { },
    mounted(){
        this.purTitle = this.$route.meta.title;
        this.logo = sessionStorage.getItem('companyName') // 企业名称
        // 文件上传请求头
        this.header = {Authorization:sessionStorage.getItem('Authorization')};
        this.http  = this.$store.state.http;                 // IP地址
        this.getpurcon();
    },
    methods: {
        // 点击查看图片文件
        scopImg(data){
          console.log(data)
            this.oneData = data;
            this.ImgVisible = true;
            this.geshi = data.fileName.substring(data.fileName.length - 4);
            let arrs = ['.pdf','.PDF','docx','DOCX','.doc','.DOC'];
            if(!arrs.includes(this.geshi)){
              this.itemshow = false;
              this.api.file.preFIle({file_id:data.id})
              .then(res=>{
                  if(res.data.code == 200){
                    this.loadingimg = false;
                    this.inmages = res.data.data;
                  }
              })
            }else{
              this.itemshow = true;
              this.src = pdf.createLoadingTask({
                url:this.http + '/saas-common/upload/previewStream?fileName=' + data.fileName + '&&AuthorizationId=wxda0b1c468b83bc99&&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f' ,     // 正式环境用这个！
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

        // 表格数据
        getpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
                searchStr:this.purSearch,
                fileType:5,
                effective:true,
                page:{
                    current:this.current,
                    size:this.size
                }
            }
            // 渲染表格
            this.api.addfile.all(listByPageData)
            .then(res=>{
              if(res.data.code == 200){
                  this.loading = false; 
                  if(res.data.data.records != ""){
                    console.log(res.data.data)
                    this.tableData = res.data.data.records
                    this.total = res.data.data.total
                  }else{
                    this.dataAva = true;
                  }
              }
            })
        },
        // 搜索功能
        search() {
            this.tableData = [];
            this.current = 1;
            this.getpurcon();
        },
        // 表格选中数据
        purChange(val) {
            this.multipleSelection = val;
            this.delmul = [];
            for(let i=0;i<this.multipleSelection.length;i++){
            this.delmul.push(this.multipleSelection[i].id) 
            }
        },
        // 打开弹窗
        Importdailog(){
            this.ImportVisible = true;
            this.enclosurefileList = [];
            this.inmages = [];
        },
        // 保存表单提交
        AddDatasubmitForm(formName) {
          this.$refs[formName].validate((valid) => {
            if (valid) {
              if(this.$refs.upload.uploadFiles.some(f => f.status !== 'success')){
                  // 还有文件没传完:此时确定只提交已传完的,其余的传完也不会被关联
                  return this.$message({ type: 'warning', message: '文件还在上传,请等传完再确定' });
              }
              // 文件上传
              var arr = [];
              this.enclosurefileList.map((itme,index)=>{
                  arr.push({id:itme.response.data[0].id,fileDescription:this.AddData.fileDescription,effective:true,});
              });
              if(arr.length != 0){
                this.api.addfile.updBthById(arr)
                .then(res=>{
                    if(res.data.code == 200){
                        var that = this
                        that.$message({
                            type: "success",
                            message: "上传文件合同成功!"
                        });
                        that.ImportVisible = false;
                        that.AddData = {};
                        that.getpurcon()
                    }
                })
              }else{
                this.$message({
                    type: 'warning',
                    message: '请上传合同附件!'
                });
              }
            }
          });
        },
        // 取消导入
        // Importble(){
        //     this.ImportVisible = false;
        //     this.fileForm.importnav = "";
        //     this.$refs.upload.clearFiles();
        // },
        // 关闭弹窗
        ImporthandleClose() {
            this.ImportVisible = false;
            this.AddData = {};
            this.getpurcon();
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
                "xlsx",
                "XLSX",
                "xls",
                "XLS",
            ];
            const status = types.indexOf(str) != -1;
            this.fileObj[file.name] = true;
            if(!status) {
                this.fileObj[file.name] = false;
                this.$message.error("上传文件只能是 PNG JPG PDF DOCX DOC XLSX XLS JPEG 格式,请重新上传!");
            }
            return status;
        },
        // 附件上传文件成功
        resfile(response, file, fileList){
            if(response.code == 200 && response.data && response.data.length){ // 被服务端拒收的文件不能留在列表里
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
        // 点击文件下载
        openfile(){
            let url = this.http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + this.oneData.fileName;
            window.open(url)
        },
        // 打印文件
        downPrint(){
         let url =  this.http + '/saas-common/upload/previewStream?fileName=' + this.oneData.fileName + '&&AuthorizationId=wxda0b1c468b83bc99&&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f';     // 正式环境用这个！
          window.open(url)
        },
        // 删除附件文件
        fileRemove(file,fileList){
            if(this.fileObj[file.name]){
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
            }
        },

        // 删除功能
        openPur() { 
            if(this.multipleSelection == ""){
                this.$message({
                  type: 'warning',
                  message: '请选择要删除的数据!'
                });
            }else{
                this.$confirm('此操作将永久删除该文件, 是否继续?', '提示', {
                  confirmButtonText: '确定',
                  cancelButtonText: '取消',
                  type: 'warning'
                }).then(() => {
                  this.api.addfile.alldel({ids:this.delmul})
                  .then(res=>{
                      if(res.data.code == 200){
                      this.$message({
                        type: 'success',
                        message: '删除成功!'
                      });
                      this.tableData = [];
                      this.current = 1;
                      this.getpurcon();
                      if(this.tableData == []){
                        this.dataAva = true;
                      }
                      }
                  })
                })
            }
        },
    },
};
</script>
<style lang="less" scoped>
// 文件图片查看样式
.imgBox{
    width: 100%;
    height: 600px;
    overflow: hidden;
    overflow-y: scroll;
    position: relative;
    img{
      width: 100%;
    }
}
.imgbuttons{
  width:24%;
  height: 40px;
  position: absolute;
  top: 0;
  right: 0;
  bottom: -90%;
  left: 0;
  margin: auto;
  z-index: 10000;
}
// 文件图片查看样式结束
.AddbuttonUpload{
    width: 100%;
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
}
.AddbuttonUpload:hover{
    background: #e4e4e4;
    color: #333;
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
    .el-select{
      width: 11%;
      margin-right: 10px;
      float: right;
    }
    .el-input {
      width: 10%;
      float: right;
      margin-right: 10px;
    }
    .el-button {
      margin-right: 10px;
      float: right;
    }
    .el-date-picker{
      width: 10%;
      float: right;
      margin-right: 10px;
    }
  }
}
// 头部结束
// 内容开始
.pur-nav {
  width: calc(100% - 32px);
  height: calc(100vh - 142px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 8px;
  .pur-table{
    width:100%;
    height: 94%;
  } 
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
      left:25%;
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
  // 分页
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
}
// 内容结束
</style>