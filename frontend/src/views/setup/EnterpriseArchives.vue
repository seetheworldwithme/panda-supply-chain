<template>
    <!-- 企业档案页面 -->
    <el-form :model="AddData" ref="AddData" size="small" label-width="130px">
        <!-- 头部开始 -->
        <el-row class="pur-top">
        <el-col :span="3"
            ><span class="pur-size">{{ purTitle }}</span></el-col
        >
        <el-col :span="21" class="pur-right">
            <el-button
            plain
            type="primary"
            size="small"
            @click="AddDatasubmitForm('AddData')" v-if="$route.meta.edit">更新</el-button>
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
                <div class="see-top">
                    <div></div>
                    基础信息
                </div>
                <el-row>
                    <el-col :span="2">
                    <!-- <span style="color: red">* </span> -->
                    <label
                        style="
                        color: #666;
                        font-size: 14px;
                        line-height: 30px;
                        padding-left: 12px;
                        "
                        >企业名称</label
                    >
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-input
                        v-model="AddData.enterpriseName"
                        placeholder="请输入企业名称"
                        ></el-input>
                    </el-col>
                    </el-form-item>
                </el-row>
                <el-row>
                    <el-col :span="2">
                        <label style="color: #666;font-size: 14px;line-height: 30px;padding-left: 12px;">企业轮播图</label>
                    </el-col>
                    <el-form-item>
                        <el-col :span="10">
                            <el-upload
                            :limit="2"
                            :action="http +'/saas-common/upload/batch'"
                            :headers="header"
                            name='file'
                            list-type="picture-card"
                            :file-list="enfileList"
                            :on-remove="handleRemove"
                            :on-success="handSuccess"
                            :before-upload="beforeUpload">
                                <i slot="default" class="el-icon-plus"></i>
                                <div slot="file" slot-scope="{file}">
                                    <img v-if="file.url"
                                        class="imagesurl"
                                        :src="file.url" alt=""
                                    >
                                    <!-- <img v-else :src="'data:image/jpg;base64,' + url" alt=""> -->
                                    <span class="el-upload-list__item-actions">
                                        <span
                                        class="el-upload-list__item-preview"
                                        @click="handlePictureCardPreview(file)"
                                        >
                                        <i class="el-icon-zoom-in"></i>
                                        </span>
                                        <span
                                        v-if="!disabled"
                                        class="el-upload-list__item-delete"
                                        @click="handleRemove(file)"
                                        >
                                        <i class="el-icon-delete"></i>
                                        </span>
                                    </span>
                                </div>
                                <div slot="tip" style="color:rgb(102, 102, 102);">最多上传3张轮播图</div>
                            </el-upload>
                            <el-dialog :visible.sync="dialogVisible">
                            <img width="100%" :src="dialogImageUrl" alt="">
                            </el-dialog>


                            <!-- <el-upload
                                :limit="3"
                                class="img-uploader"
                                :file-list="enclosurefileList"
                                :action="http +'/saas-common/upload/batch'"
                                :headers="header"
                                :show-file-list="false"
                                name='file'
                                :on-success="handleQrSuccess"
                                :before-upload="beforelogoUpload">
                                <img v-if="qrUrl" :src="qrUrl" class="avatar"/>
                                <i v-else class="el-icon-plus avatar-uploader-icon"></i>
                            </el-upload> -->
                        </el-col>
                    </el-form-item>
                </el-row>
                <el-row>
                    <el-col :span="2">
                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">企业地区</label>
                    </el-col>
                    <el-form-item>
                        <el-col :span="7">
                            <el-select @change="enterprises($event)" v-model="AddData.provinceCode" placeholder="省" style="width:24%;margin-right:1.33%;">
                                <el-option
                                v-for="item in provinceOptions"
                                :key="item.code"
                                :label="item.name"
                                :value="item.code">
                                </el-option>
                            </el-select>
                            <el-select v-model="AddData.cityCode" placeholder="市" style="width:24%;margin-right:1.33%;" @change="citys($event)">
                                <el-option
                                v-for="item in cityOptions"
                                :key="item.value"
                                :label="item.name"
                                :value="item.code">
                                </el-option>
                            </el-select>
                            <el-select v-model="AddData.areaCode" placeholder="县 / 区" style="width:24%;margin-right:1.33%;" @change="countys($event)">
                                <el-option
                                v-for="item in countyOptions"
                                :key="item.value"
                                :label="item.name"
                                :value="item.code">
                                </el-option>
                            </el-select>
                            <el-select v-model="AddData.countyCode" placeholder="乡" style="width:24%;" @change="villages($event)">
                                <el-option
                                v-for="item in villageOptions"
                                :key="item.value"
                                :label="item.name"
                                :value="item.code">
                                </el-option>
                            </el-select>
                        </el-col>
                    </el-form-item>
                </el-row>
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
                        >企业性质</label
                    >
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-select
                        v-model="AddData.enterpriseType"
                        placeholder="请选择企业性质"
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
                    <!-- <span style="color:red;">* </span> -->
                    <label
                        style="
                        color: #666;
                        font-size: 14px;
                        line-height: 30px;
                        padding-left: 12px;
                        "
                        placeholder="请输入税号"
                        >税号</label
                    >
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-input v-model="AddData.rateNumber"></el-input>
                    </el-col>
                    </el-form-item>
                </el-row>
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
                        placeholder="请输入经营地址"
                        >经营地址</label
                    >
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-input v-model="AddData.address"></el-input>
                    </el-col>
                    </el-form-item>
                </el-row>
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
                        placeholder="请输入法人姓名"
                        >法人</label
                    >
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-input v-model="AddData.legalPerson"></el-input>
                    </el-col>
                    </el-form-item>
                </el-row>
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
                        placeholder="请输入实际控制人姓名"
                        >实际控制人</label
                    >
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-input v-model="AddData.controllerPerson"></el-input>
                    </el-col>
                    </el-form-item>
                </el-row>
                <el-row>
                    <el-col :span="2">
                    <label style="color: #666;font-size: 14px;line-height: 30px;padding-left: 12px;" placeholder="请输入主要联系人">主要联系人</label>
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-input v-model="AddData.contactsPerson"></el-input>
                    </el-col>
                    </el-form-item>
                </el-row>
                <el-row>
                    <el-col :span="2">
                        <span style="color:red;">* </span>
                        <label style="color:#666;font-size:14px;line-height:30px;">企业类别</label>
                    </el-col>
                    <el-form-item>
                        <el-col :span="7">
                            <el-select v-model="AddData.type" placeholder="请选择企业类别" style="width:100%;">
                                <el-option
                                v-for="item in typeOptions"
                                :key="item.id"
                                :label="item.name"
                                :value="item.id">
                                </el-option>
                            </el-select>
                        </el-col>
                    </el-form-item>
                </el-row>
                <el-row>
                    <el-col :span="2">
                        <label style="color: #666;font-size: 14px;line-height: 30px;padding-left: 12px;">营业执照</label>
                    </el-col>
                    <el-form-item>
                        <el-col :span="7">
                            <el-upload
                                class="img-uploader"
                                :file-list="enclosurefileList"
                                :action="http +'/saas-common/upload/batch'"
                                :headers="header"
                                :show-file-list="false"
                                name='file'
                                :on-success="handleQrSuccess"
                                :before-upload="beforelogoUpload">
                                <img v-if="qrUrl" :src="qrUrl" class="avatar"/>
                                <i v-else class="el-icon-plus avatar-uploader-icon"></i>
                            </el-upload>
                        </el-col>
                    </el-form-item>
                </el-row>
                </el-col>
                <!-- 生产经营信息 -->
                <el-col :span="24">
                <div class="see-top"><div></div>产品</div>
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
                        >设计产品</label
                    >
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-input
                        v-model="AddData.coalMine"
                        placeholder="主力产品类型+主力产品指标"
                        ></el-input>
                    </el-col>
                    </el-form-item>
                </el-row>
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
                        >实际产量</label
                    >
                    </el-col>
                    <el-form-item>
                    <el-col :span="7">
                        <el-input
                        v-model="AddData.annualOutput"
                        placeholder="主力产品类型+主力产品指标"
                        ></el-input>
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
            http:'',
            header:{},                      // 上传请求头
            qrUrl:'',
            purTitle: "", // 标题
            typeOptions:[
                {
                    id:"1",
                    name:'一般纳税人'
                },
                {
                    id:"2",
                    name:'小规模纳税人'
                },
            ],
            // 企业性质选择下拉框
            SelectcontractName: [
                {
                value: "1",
                label: "国企",
                },
                {
                value: "2",
                label: "民营企业",
                },
                {
                value: "3",
                label: "央企",
                },
                {
                value: "4",
                label: "外企",
                },
            ], 
            // 行政区划
            provinceOptions:[],             // 省
            cityOptions: [],                // 市
            countyOptions: [],              // 县 区
            villageOptions: [],             // 乡
            // 表单数据
            AddData: {
                filelists:[],
            },
            enclosurefileList:[], 
            enfileList:[],   // 轮播图上传列表
            dialogImageUrl: '',
            dialogVisible: false,
            disabled: false,
            // list:[],
            // url:''
        };
    },
    created() {
        
        this.header = {Authorization:sessionStorage.getItem('Authorization')}  // 文件上传请求头
        this.http = this.$store.state.http;                 // IP地址
        this.purTitle = this.$route.meta.title;
        this.api.enterData.byidData({id:sessionStorage.getItem('companyId')})
        .then(res=>{
            if(res.data.code == 200){
                this.AddData = res.data.data;
                // 企业地区
                this.api.areaInfo.all({parentCode:0})
                .then(res=>{
                    if(res.data.code == 200){
                        this.provinceOptions = res.data.data;
                    }
                })
                if(this.AddData.cityCode != null || this.AddData.cityCode != ""){
                    // 市
                    this.api.areaInfo.all({parentCode:this.AddData.provinceCode})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.cityOptions = res.data.data;
                        }
                    })
                    // 县 区
                    this.api.areaInfo.all({parentCode:this.AddData.cityCode})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.countyOptions = res.data.data;
                        }
                    })
                    // 乡
                    this.api.areaInfo.all({parentCode:this.AddData.areaCode})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.villageOptions = res.data.data;
                        }
                    })
                }
            }
        })
        this.api.file.imageShow({relation_id:sessionStorage.getItem('companyId')})
        .then(res=>{
            if(res.data.code == 200){
                this.qrUrl = res.data.data.yyzz;
                for(let i in res.data.data.lbts){
                    this.enfileList.push(res.data.data.lbts[i])
                }
            }
        })
    },
    mounted(){},
    methods: {
        // 删除轮播
        handleRemove(file) {
            let arr = [];
            this.enfileList.forEach((value,index,array)=>{
                if(value.id == file.id){
                    array.splice(value,1)
                }else{
                    arr.push(value)
                }
                
            })
            this.enfileList = arr;
        },
        handlePictureCardPreview(file) {
            this.dialogImageUrl = file.url;
            this.dialogVisible = true;
        },
        // 轮播图上传成功
        handSuccess(res, file) {
            if(res.code == 200){
                // this.$message({
                //     message:'上传轮播图成功',
                //     type:'success',
                // })
                this.enfileList.push(res.data[0]);
            }else{
                this.$message({
                    message:'上传轮播图失败,请重新上传!',
                    type:'error'
                })
            }
        },
        // 轮播图认证
        beforeUpload(file) {
            var types = [
                "image/jpeg",
                "image/png", 
                "image/jpg"
            ];
            const isImg = types.indexOf(file.type) == -1;
            const isLt2M = file.size / 1024 / 1024 < 2;
            if (isImg) {
                this.$message.error("上传轮播图图片只能是 JPG PNG JPEG 格式!");
            }
            if (!isLt2M) {
                this.$message.error("上传轮播图图片大小不能超过 2MB!");
            }
            return !isImg && isLt2M;
        },
        // 更新企业信息表单提交
        AddDatasubmitForm(formName){
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.AddData.id = sessionStorage.getItem('companyId');
                    this.AddData.imagelist = this.enfileList;
                    let arr = [];
                    for(let ki in this.enfileList){
                        arr.push(this.enfileList[ki].id);
                    }
                    
                    this.AddData.ids = arr;
                    this.api.enterData.updateById(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                type: "success",
                                message: "更新企业信息成功!",
                            });
                        }
                    })
                }
            });
        },
        // 营业执照上传
        handleQrSuccess(res, file) {
            if(res.code == 200){
                this.$message({
                    message:'上传营业执照成功',
                    type:'success',
                })
                this.AddData.filelist = [{id:res.data[0].id}]
                this.qrUrl = res.data[0].url
            }else{
                this.$message({
                    message:'上传营业执照失败,请重新上传!',
                    type:'error'
                })
            }
        },
        // 封装文件流函数
        // ImgUrl(res){
        //     this.qrUrl = "";
        //     const dataInfo = res;
        //     let reader = new window.FileReader()
        //     // 使用readAsArrayBuffer读取文件, result属性中将包含一个 ArrayBuffer对象以表示所读取文件的数据
        //     reader.readAsArrayBuffer(dataInfo)
        //     var that = this;
        //     reader.onload = function (e) {
        //         const result = e.target.result
        //         const contentType = result.type
        //         // 生成blob图片,需要参数(字节数组, 文件类型)
        //         const blob = new Blob([result], { type: contentType })
        //         // 使用 Blob 创建一个指向类型化数组的URL, URL.createObjectURL是new Blob文件的方法,可以生成一个普通的url,可以直接使用,比如用在img.src上
        //         const url = window.URL.createObjectURL(blob);
        //         // console.log(url) // 产生一个类似 blob:d3958f5c-0777-0845-9dcf-2cb28783acaf 这样的URL字符串
        //         that.qrUrl = url
        //     }
        // },
        // 营业执照图片格式认证
        beforelogoUpload(file) {
            var types = [
                "image/jpeg",
                "image/png", 
                "image/jpg"
            ];
            const isImg = types.indexOf(file.type) == -1;
            const isLt3M = file.size / 1024 / 1024 /1024 < 3;
            if (isImg) {
                this.$message.error("上传营业执照图片只能是 JPG PNG JPEG 格式!");
            }
            if (!isLt3M) {
                this.$message.error("上传营业执照图片大小不能超过 3MB!");
            }
            return !isImg && isLt3M;
        },
        // 选择省
        enterprises(id){
            let resultArr = this.provinceOptions.filter((item)=>{
                return item.code == id;
            });
            this.AddData.proName = resultArr[0].name;
            this.api.areaInfo.all({parentCode:id})
            .then(res=>{
                if(res.data.code == 200){
                    this.AddData.cityCode = '';                 // 市
                    this.AddData.areaCode = '';                 // 县 区
                    this.AddData.countyCode = '';               // 乡
                    this.cityOptions = res.data.data;
                }
            })
        },
        // 选择市
        citys(id){
            let resultArr = this.cityOptions.filter((item)=>{
                return item.code == id;
            });
            this.AddData.cityName = resultArr[0].name;
            this.api.areaInfo.all({parentCode:id})
            .then(res=>{
                if(res.data.code == 200){
                    this.AddData.areaCode = [];                 // 县 区
                    this.AddData.countyCode = '';               // 乡
                    this.countyOptions = res.data.data;
                }
            })
        },
        // 选择县 区
        countys(id){
            let resultArr = this.countyOptions.filter((item)=>{
                return item.code == id;
            });
            this.AddData.areaName = resultArr[0].name;
            this.api.areaInfo.all({parentCode:id})
            .then(res=>{
                if(res.data.code == 200){
                    this.AddData.countyCode = '';               // 乡
                    this.villageOptions = res.data.data;
                }
            })
        },
        // 选择乡
        villages(id){
            let resultArr = this.villageOptions.filter((item)=>{
                return item.code == id;
            });
            this.AddData.countryName = resultArr[0].name;
        }
    },
};
</script>
<style lang="less" scoped>
// 上传图片样式开始
.imagesurl{
    width:100%;
    height: 94px;
    margin-top:18%;
  
}
// .imagesurl .el-upload {
//     width:160px;
//     height: 94px;
//     cursor: pointer;
//     float: left;
//     overflow: hidden;
// }
.img-uploader{
    width: 80px;
    height: 106px;
}
.img-uploader .el-upload {
    width: 80px;
    height: 106px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
}
.img-uploader .el-upload:hover {
    border-color: #409eff;
}
[class*=" avatar-uploader-icon"], [class^=avatar-uploader-icon]{
    border: 1px dashed #d9d9d9;
    border-radius: 4px;
    font-size: 28px;
    color: #8c939d;
    width: 80px;
    height: 106px;
    line-height: 106px;
    text-align: center;
}
.avatar {
    width: 80px;
    height: 106px;
    display: block;
}
// 上传图片样式结束
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
// 表单内容开始
.pur-nav {
  height: calc(100vh - 126px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 33px 0 33px 33px;
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