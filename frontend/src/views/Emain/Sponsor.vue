<template>
<div>
    <div class="box" v-if="TaburlShow">
        <el-empty style="padding:70% 0 0;box-sizing: border-box;" description="暂无数据！">
            <!-- <div>{{ceshiurl}}</div> -->
        </el-empty>
    </div>
    <!-- e签宝响应式签署列表页面 -->
    <div class="casketlog" v-else>
        <div class="conentherard">
            <el-select v-model="carriageName" @change="prologistics">
                <el-option
                v-for="item in Extended"
                :key="item.value"
                :label="item.dkey"
                :value="item.dvalue">
                </el-option>
            </el-select>
            <el-button type="primary" @click="launchGo">发起签署</el-button>
            <el-input
            @keyup.enter.native="search"
            placeholder="请输入内容"
            prefix-icon="el-icon-search"
            v-model="purSearch">
            </el-input>
            <!-- <div>{{ceshiurl}}</div> -->
        </div>
        <div class="conent" v-if="!dataAva">
            <ul class="col" v-for="(v,k) in tableData" :key="k">
                <li style="width:70%;overflow:hidden;text-overflow:ellipsis;white-space: nowrap;" v-if="carriageName == 4">
                    <b>签署主题名称：</b>
                    <span>{{v.topicName}}</span>
                </li>
                <li style="width:30%;" v-if="carriageName == 4">
                    <el-button type="primary" size="mini" @click="eqburl(v)">合同预览</el-button>
                </li>
                <li style="width:100%;overflow:hidden;text-overflow:ellipsis;white-space: nowrap;" v-else>
                    <b>签署主题名称：</b>
                    <span>{{v.topicName}}</span>
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
                <el-empty description="暂无数据！">
                    <!-- <div>{{ceshiurl}}</div> -->
                </el-empty>
                <!-- <img src="../../assets/zanwushuju.png" alt="">
                <p>暂 无 数 据 !</p> -->
            </div>
        </div>
        <!-- 底部 -->
        <!-- <div class="copyright-bottom">
            <div class="copyright"><span>Copyright © 2021 山西海云启辰科技有限公司 版权所有</span><p> 赣ICP备2021002370号</p></div>
        </div> -->
        <!-- 复制文本弹框 -->
        <!-- <el-dialog :visible.sync="dialogTableVisible" width="80%">
            <div class="tishibox"><i class="el-icon-success"></i>复制链接成功，请在浏览器打开！</div>
        </el-dialog> -->
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
            // dialogTableVisible: false,
            TaburlShow:null,
            ceshiurl:'',
        }
    },
    created(){
        // this.GetRequest();  // 存储手机号 类型   
        // window.location.search(设置或获取 href 属性中跟在问号后面的部分)  window.location.href
        // + '?account=admin&level=1';
        // window.location.reload();
        // let a = window.location.href;
        // sessionStorage.setItem('url',a);
        // let url = sessionStorage.getItem('url');
        let url = window.location.href;
        // this.ceshiurl = url;
        // console.log(a)
        if(!sessionStorage.getItem('Authorization')){
            if(url.indexOf('?') == -1){
                console.log('不存在账号');
                this.TaburlShow = true;
                // this.$router.push({name:"Unknown"});
            }else{
                console.log('存在账号')
                let pos = url.indexOf('?');
                let result = url.substring(pos);
                let theRequest = new Object();
                let str = result.substr(1);
                let strs = str.split("&");
                for(var i=0;i<strs.length;i++) {
                    theRequest[strs[i].split("=")[0]] = unescape(strs[i].split("=")[1]);
                }
                // console.log(theRequest)
                if(theRequest.level == 1){
                    this.TaburlShow = false;
                    console.log('有权限')
                    // sessionStorage.clear();
                    this.api.wxChat.wxChat({account:theRequest.account})
                    .then(res=>{
                        if(res.data.code == 200){
                            sessionStorage.setItem('Authorization','Bearer ' + res.data.data.token);
                            sessionStorage.setItem('companyName',res.data.data.companyName)
                            sessionStorage.setItem('username',res.data.data.username)
                            sessionStorage.setItem('level',res.data.data.level)
                            sessionStorage.setItem('path','')
                            sessionStorage.setItem('pathID','')
                            sessionStorage.setItem('companyId',res.data.data.companyId)
                            this.getpurcon();
                        }
                    })
                }else{
                    console.log('无权限')
                    this.TaburlShow = true;
                    // this.$router.push({name:"Unknown"})
                }
            }
        }else{
            this.getpurcon();
        }
    },
    methods:{
        // 发起签署
        launchGo(){
            this.$router.push({name:"InitiateStart"})
        },
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
                this.api.eqbflows.selSignUrl({eqbFlowId:data.eqbFlowId})
                .then(res=>{
                    if(res.data.code == 200){
                        // let url = res.data.data;
                        // let oInput = document.createElement('input')
                        // oInput.value = url
                        // document.body.appendChild(oInput)
                        // oInput.select() // 选择对象
                        // document.execCommand("Copy") // 执行浏览器复制命令
                        // this.dialogTableVisible = true;
                        // oInput.remove();
                        window.location.replace(res.data.data)
                        // window.loaction.href = res.data.data;
                        // window.open(res.data.data)
                    }
                })
            
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
                        arr.push('<div>' +res.data.data[key].signerName + ' : ' + day + '</div>')
                    }
                    let a = '<strong>' + arr.toString().replace(",","") + '</strong>';
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
.casketlog{
    width: 100%;
    height: 100vh;
    overflow: hidden;
    background: #F6F6F6;
}
// 复制文本样式
.tishibox{
    width: 100%;
    text-align: center;
    font-size: 14px;
    color: #333;
    i{
        color: rgb(65,184,131);
    }
}
// 暂无数据开始
.conentno{
    width:100%;
    height: calc(100vh - 120px);
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
    height: 120px;
    padding:10px 20px;
    box-sizing:border-box;
    .el-input {
        width:20%;
        float: left;
        margin:5px 0;
    }
    .el-button{
        width: 10%;
        float: left;
        margin:5px 1%;
    }
    .el-select{
        width: 10%;
        float: left;
        margin:5px 0;
    }
}
.conent{
    width:100%;
    height: calc(100vh - 40px);
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
            overflow: hidden;
            text-overflow:ellipsis;
            white-space: nowrap;
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
    .conentherard{
        width: 100%;
        height: 60px;
        padding:10px 50px;
        box-sizing:border-box;
        background: #F6F6F6;
        .el-input {
            width:30%;
            float: left;
            margin:5px 0;
        }
        .el-button{
            width: 20%;
            float: left;
            margin:5px 1%;
        }
        .el-select{
            width: 20%;
            float: left;
            margin:5px 0;
        }
    }
    .conent{
        width:100%;
        height: calc(100vh - 40px);
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
                overflow: hidden;
                text-overflow:ellipsis;
                white-space: nowrap;
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
    .conentherard{
        width: 100%;
        height: 60px;
        padding:10px 25px;
        box-sizing:border-box;
        background: #F6F6F6;
        .el-input {
            width:35%;
            float: left;
            margin:5px 0;
        }
        .el-button{
            width: 20%;
            float: left;
            margin:5px 2%;
        }
        .el-select{
            width: 20%;
            float: left;
            margin:5px 0;
        }
    }
    .conent{
        width:100%;
        height: calc(100vh - 100px);
        padding: 0 25px 0;
        margin-top:18px;
        box-sizing: border-box;
        background-color: #F6F6F6;
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
                overflow: hidden;
                text-overflow:ellipsis;
                white-space: nowrap;
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
    .conentherard{
        width: 100%;
        height: 120px;
        padding:10px 20px;
        box-sizing:border-box;
        background: #F6F6F6;
        .el-select{
            width: 49%;
            float: left;
            margin:5px 0;
        }
        .el-button{
            width: 49%;
            float: left;
            margin:5px 0 5px 2%;
        }
        .el-input {
            width:100%;
            float: left;
            margin:5px 0;
        }
    }
    .conent{
        width:100%;
        height: calc(100vh - 120px);
        padding: 0 20px 0;
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
                overflow: hidden;
                text-overflow:ellipsis;
                white-space: nowrap;
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


