<template>
    <!-- e签宝响应式页面 -->
    <div>
        <div class="casket">
            <div class="casketlog">
                <!-- 头部 -->
                <div class="login-concent">
                    <div class="img-lg">
                        <img src="../../assets/lg.svg" alt="">
                    </div>
                    <div class="img-size"><span>云管理</span></div>
                    <div class="goTu" @click="commands"><i class="el-icon-alituichu"></i><span> 退出</span></div>
                </div>
                <div class="conentherard">
                    <el-select v-model="carriageName" @change="prologistics">
                        <el-option
                        v-for="item in Extended"
                        :key="item.value"
                        :label="item.dkey"
                        :value="item.dvalue">
                        </el-option>
                    </el-select>
                    <el-input
                    @keyup.enter.native="search"
                    placeholder="请输入内容"
                    prefix-icon="el-icon-search"
                    v-model="purSearch">
                    </el-input>
                </div>
                <div class="conent" v-if="!dataAva">
                    <ul class="col" v-for="(v,k) in tableData" :key="k" @click="eqburl(v)">
                        <li style="width:80%;">
                            <b>签署主题名称：</b>
                            <span>{{v.topicName}}</span>
                        </li>
                        <li style="width:20%;" v-if="carriageName == 2 || carriageName == 4">
                            <el-button type="primary" size="small" @click.stop="eqburl(v)">进入</el-button>
                        </li>
                        <li>
                            <b>合同名称：</b>
                            <span>{{v.contractName}}</span>
                        </li>
                        <li>
                            <b>发起方企业名称：</b>
                            <span>{{v.etcName}}</span>
                        </li>
                        <li style="width:100%;">
                            <b>发起方姓名：</b>
                            <span>{{v.participants}}</span>
                        </li>
                        <li style="width:100%;">
                            <b>发起方手机号：</b>
                            <span>{{v.phone}}</span>
                        </li>
                        <li style="width:50%;">
                            <b>流程状态：</b>
                            <span v-if="v.flowStatus == -1">待确认</span>
                            <span v-if="v.flowStatus == 0">草稿</span>
                            <span v-if="v.flowStatus == 1">签署中</span>
                            <span v-if="v.flowStatus == 2">完成</span>
                            <span v-if="v.flowStatus == 3">撤销</span>
                            <span v-if="v.flowStatus == 5">过期</span>
                            <span v-if="v.flowStatus == 7">拒签</span>
                        </li>
                        <li style="width:50%;">
                            <b>发起方类型：</b>
                            <span v-if="v.launchType == '1'">个人</span>
                            <span v-if="v.launchType == '2'">企业</span>
                        </li>
                        <li style="width:100%;" v-if="v.flowStatus != -1">
                            <b>签署人状态：</b>
                            <span style="color:#CF1724;" plain @click.stop="confselSign(v)">查询签署</span>
                        </li>
                        <li style="padding-top:2%;box-sizing: border-box;margin:0;" v-if="carriageName == 3">
                            <el-button type="primary" size="small" @click.stop="confirmgo(v)">确认</el-button>
                        </li>
                    </ul>
                </div>
                <!-- 暂无数据 -->
                <div class="conentno" v-if="dataAva">
                    <div class="dataAva">
                        <img src="../../assets/zanwushuju.png" alt="">
                        <p>暂 无 数 据 !</p>
                    </div>
                </div>
                <!-- 升级企业信息 -->
                <div class="conentbottom">
                    <!-- <div class="Fsize">
                        <span style="float:left;">注册新企业</span>
                        <span @click="opengoqe">加入已有企业</span>
                        企业名称
                    </div> -->
                </div>
                <!-- 底部 -->
                <div class="copyright-bottom">
                    <div class="copyright"><span>Copyright © 2021 山西海云启辰科技有限公司 版权所有</span><p> 赣ICP备2021002370号</p></div>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
export default {
    data(){
        return {
            purSearch: "",          // 搜索
            Extended:[
                {
                    dkey:'我发起的',
                    dvalue:1
                },
                {
                    dkey:'我参与的',
                    dvalue:2
                },
                {
                    dkey:'待确认的',
                    dvalue:3
                },
                {
                    dkey:'已完成的',
                    dvalue:4
                },
            ],            // 查询类型
            carriageName:3,
            current:1,              // 当前页数
            size:50,                // 显示条数
            total:0,                // 总条数
            tableData: [],          // 表格数据
            dataAva:false,          // 无数据状态
        }
    },
    created(){
        this.getpurcon();
    },
    methods:{
        // 表格数据
        getpurcon(){
            this.dataAva = false;
            // 传参
            var listByPageData = {
                searchStr:this.purSearch,
                searchType:this.carriageName,
                page:{
                    current:this.current,
                    size:this.size
                }
            }
            // 渲染表格
            this.api.eqbflows.all(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    if(res.data.data.records != ""){
                        this.tableData = res.data.data.records;
                        this.total = res.data.data.total;
                    }else{
                        this.dataAva = true;
                    }
                }
            })
        },
        // 搜索功能
        search() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 查询类型搜索功能
        prologistics(){
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        }, 
        // 加入企业
        opengoqe(){
            this.$prompt('企业邀请码向管理员联系申请', '加入企业', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            inputValue: '请输入邀请码',
            closeOnClickModal:false,
            }).then(({ value }) => {
                this.$message({
                    type: 'success',
                    message: '你的邀请码是: ' + value
                });
            }).catch(() => {
                // this.$message({
                //     type: 'info',
                //     message: '取消输入'
                // });       
            });
        },
        // 确认发起签署
        confirmgo(data){
            let obj = {
                id:data.id,
                companyId:data.companyId,
                field:data.field,
                topicName:data.topicName,
                launchPerId:data.launchPerId,
                launchPubId:data.launchPubId
            }
            this.api.eqbflows.confirmFlow(obj)
            .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                        type: "success",
                        message: "确认成功!",
                    });
                    this.tableData = [];
                    this.getpurcon();
                    this.current = 1;
                }
            })
        },
        // e签宝点击链接
        eqburl(data){
            console.log(data)
            if(this.carriageName != 3){
                this.api.eqbflows.selSignUrl({eqbFlowId:data.eqbFlowId})
                .then(res=>{
                    if(res.data.code == 200){
                        window.open(res.data.data)
                    }
                })
            }
        },
        // 查询签署人状态
        confselSign(data){
            this.api.eqbflows.selSignPeople({eqbFlowId:data.eqbFlowId})
            .then(res=>{
                if(res.data.code == 200){
                    let arr = [];
                    let day;
                    for(let key in res.data.data){
                        switch (res.data.data[key].signStatus) {
                            case 0:
                                day = "待签";
                                break;
                            case 1:
                                day = "未签";
                                break;
                            case 2:
                                day = "已签";
                                break;
                            case 3:
                                day = "待审批";
                                break;
                            case 4:
                                day = "拒签";
                                break;
                        }
                        // 签署方姓名是用户输入,拼进 HTML 前必须转义(否则对方点"查询签署"即执行脚本)
                        arr.push('<div>' + String(res.data.data[key].signerName == null ? '' : res.data.data[key].signerName).replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c])) + ' : ' + day + '</div>')
                    }
                    let a = '<strong>' + arr.join('') + '</strong>';
                    this.$notify({
                        title: data.topicName,
                        dangerouslyUseHTMLString: true,
                        message: a,
                        duration: 0
                    });
                }
            })
        },
        // 退出
        commands(){
            sessionStorage.clear();
            this.$router.push({name:'ELogin'});
        }
    }
}
</script>
<style lang="less" scoped>
// 头部开始
.login-concent{
    width: 100%;
    height: 70px;
    padding:8px 5%;
    box-sizing: border-box;
    border-bottom: 1px solid #e4e4e4;
    line-height: 48px;
    overflow: hidden;
    background-color: #fff;
    .img-lg{
        float: left;
        img{
            width:48px;
            height: 48px;
        }
    }
    .img-size{
        float: left;
        margin-left:24px;
        span{
            font-size: 24px;
            font-weight: 600;
            color: #333;
        }
    }
    .goTu{
        float: right;
        cursor: pointer;
        i{
            font-size: 18px;
            color:#333;
            font-weight: 600;
        }
        span{
            font-size: 18px;
            color: #333;
            font-weight: 600;
        }
    }
}
// 头部结束
// 版权开始
.copyright-bottom{
  width: 100%;
  height:60px;
  padding-top:10px;
  box-sizing: border-box;
  background-color:#F6F6F6;
  position: fixed;
  bottom: 0;
  left: 0;
  border-top: 1px solid #e4e4e4;
  .copyright{
    width:auto;
    height: 27px;
    margin: auto;
    text-align: center;
    span{
      font-size: 12px;
      line-height: 24px;
      color: #828282;
    }
    p{
      font-size: 12px;
      color: #828282;
    }
  }
}
// 版权结束
// 暂无数据开始
.conentno{
    width:100%;
    height: calc(100vh - 190px);
    padding: 32px 60px 60px;
    box-sizing: border-box;
    background: #F6F6F6;
    position: relative;
    .dataAva{
        width:370px;
        height: 240px;
        position: absolute;
        top: 0;
        left: 0;
        bottom: 0;
        right: 0;
        margin: auto;
        text-align: center;
        img{
            width: 370px;
            height: 200px;
        }
        p{
            font-size: 16px;
            line-height: 40px;
            color: #333;
        }
    }

}
// 暂无数据结束
// 中间内容开始
.conentherard{
    width: 100%;
    height: 60px;
    padding:10px 60px;
    box-sizing:border-box;
    background: #F6F6F6;
    .el-input {
      width:20%;
      float: left;
    }
    .el-select{
      width: 10%;
      float: left;
      margin-right:16px;
    }
}
.conentbottom{
    width: 100%;
    height: 40px;
    padding:0 60px;
    box-sizing:border-box;
    background: #F6F6F6;
    .Fsize{
        width: 20%;
        height: 40px;
        margin: auto;
        text-align: center;
        font-size: 18px;
        font-weight: bold;
        color: #333;
        line-height: 40px;
        span{
            cursor: pointer;
            float: right;
        }
    }
}
.conent{
    width:100%;
    height: calc(100vh - 230px);
    padding: 32px 60px 60px;
    box-sizing: border-box;
    background: #F6F6F6;
    display: flex;
    flex-wrap: wrap;
    justify-content: flex-start;
    align-content:flex-start;
    overflow-x: scroll;
    .col{
        cursor: pointer;
        width: 24%;
        background: #fff;
        margin-bottom:1%;
        margin-right:1%;
        border-radius:6px;
        padding: 16px;
        box-sizing: border-box;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        align-content:flex-start;
        li{
            width: 100%;
            margin-bottom:1.5%;
            b{
                font-size: 14px;
                color: #333;
                line-height: 30px;
            }
            span{
                font-size: 14px;
                color: #333;
                line-height: 30px;
            }
            .el-button{
                float: right;
            }
        }
    }
}
// 中间内容结束
// 屏幕尺寸大于900小于1500
@media screen and (max-width:1600px) and (min-width:1100px) {
    .conentbottom{
        width: 100%;
        height: 40px;
        padding:0 60px;
        box-sizing:border-box;
        background: #F6F6F6;
        .Fsize{
            width: 30%;
            height: 40px;
            margin: auto;
        }
    }
    .conentherard{
        width: 100%;
        height: 60px;
        padding:10px 50px;
        box-sizing:border-box;
        background: #F6F6F6;
        .el-input {
            width:30%;
            float: left;
        }
        .el-select{
            width: 15%;
            float: left;
            margin-right:16px;
        }
    }
    .conent{
        width:100%;
        height: calc(100vh - 230px);
        padding: 32px 50px 60px;
        box-sizing: border-box;
        background: #F6F6F6;
        display: flex;
        flex-wrap: wrap;
        justify-content: flex-start;
        align-content:flex-start;
        overflow-x: scroll;
        .col{
            width: 32%;
            background: #fff;
            border-radius:6px;
            padding: 16px;
            box-sizing: border-box;
            display: flex;
            flex-wrap: wrap;
            position: relative;
            .colgo{
                width: 10%;
                height: 48px;
                background: red;

            }
            li{
                width: 100%;
                margin-bottom:1.5%;
                b{
                    font-size: 14px;
                    color: #333;
                    line-height: 30px;
                }
                span{
                    font-size: 14px;
                    color: #333;
                    line-height: 30px;
                }
                .el-button{
                    float: right;
                }
            }
        }
    }
}
// 屏幕尺寸大于750小于900
@media screen and (max-width:1100px) and (min-width:750px) {
    .conentbottom{
        width: 100%;
        height: 40px;
        padding:0 60px;
        box-sizing:border-box;
        background: #F6F6F6;
        .Fsize{
            width: 35%;
            height: 40px;
            margin: auto;
        }
    }
    .conentherard{
        width: 100%;
        height: 60px;
        padding:10px 25px;
        box-sizing:border-box;
        background: #F6F6F6;
        .el-input {
            width:35%;
            float: left;
        }
        .el-select{
            width: 20%;
            float: left;
            margin-right:16px;
        }
    }
    .conent{
        width:100%;
        height: calc(100vh - 230px);
        padding: 18px 25px 60px;
        box-sizing: border-box;
        background: #F6F6F6;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        align-content:flex-start;
        .col{
            width: 49%;
            background: #fff;
            border-radius:6px;
            padding: 16px;
            box-sizing: border-box;
            display: flex;
            flex-wrap: wrap;
            li{
                width: 100%;
                margin-bottom:1.5%;
                b{
                    font-size: 14px;
                    color: #333;
                    line-height: 30px;
                }
                span{
                    font-size: 14px;
                    color: #333;
                    line-height: 30px;
                }
                .el-button{
                    float: right;
                }
            }
        }
    }
}
// 屏幕尺寸小于750
@media screen and (max-width:750px) {
    .conentbottom{
        width: 100%;
        height: 40px;
        padding:0 60px;
        box-sizing:border-box;
        background: #F6F6F6;
        .Fsize{
            width: 70%;
            height: 40px;
            margin: auto;
        }
    }
    .conentherard{
        width: 100%;
        height: 60px;
        padding:10px 20px;
        box-sizing:border-box;
        background: #F6F6F6;
        .el-input {
            width:50%;
            float: left;
        }
        .el-select{
            width: 35%;
            float: left;
            margin-right:16px;
        }
    }
    .conent{
        width:100%;
        height: calc(100vh - 230px);
        padding: 14px 20px 60px;
        box-sizing: border-box;
        background: #F6F6F6;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        align-content:flex-start;
        .col{
            width: 100%;
            background: #fff;
            margin-bottom:2%;
            border-radius:6px;
            padding: 16px;
            box-sizing: border-box;
            display: flex;
            flex-wrap: wrap;
            justify-content: space-between;
            align-content:flex-start;
            li{
                width: 100%;
                margin-bottom:1.5%;
                b{
                    font-size: 14px;
                    color: #333;
                    line-height: 30px;
                }
                span{
                    font-size: 14px;
                    color: #333;
                    line-height: 30px;
                }
                .el-button{
                    float: right;
                }
            }
        }
    }
}
</style>


