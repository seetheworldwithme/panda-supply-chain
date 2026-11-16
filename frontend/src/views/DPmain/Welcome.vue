<template>
    <div class="box">
        <div class="TimesAll">
            <div class="buttonbox" @click="BackOff">
                <i class="iconfont el-icon-alireturn"></i>
            </div>
        </div>
        <div class="boxone" @click="dialogFormVisible = true">
            <div class="boxoneh1">{{conent.content1}}</div>
            <div class="boxb">{{conent.content2}}</div>
            <div class="boxp">
                <p>{{conent.content3}}</p>
                <span>{{conent.content4}}</span>
            </div>
        </div>
        <el-dialog title="汇报内容" :visible.sync="dialogFormVisible" width="45%">
            <el-form :model="form" label-width="120px">
                <el-form-item label="标题:">
                    <el-input v-model="form.content1" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="内容1:">
                    <el-input v-model="form.content2" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="内容2:">
                    <el-input v-model="form.content3" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="公司名称:">
                    <el-input v-model="form.content4" autocomplete="off"></el-input>
                </el-form-item>
            </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button type="primary" @click="dialogFormVisibless">确 定</el-button>
        </div>
        </el-dialog>
    </div>
</template>
<script>
export default{
    data (){
        return {
            conent:'',
            dialogFormVisible: false,
            form: {
                id:1,
                content1: '',
                content2: '',
                content3: '',
                content4: '',
            },
        }
    },
    created(){
        this.api.config.getById({id:1})
        .then(res=>{
            if(res.data.code == 200){
                this.conent = res.data.data;
                this.form = res.data.data;
            }
        })
    },
    methods:{
        // 退出大屏
        BackOff(){
            this.$router.push({name:'Dpmain'})
        },
        dialogFormVisibless(){
            this.api.config.edit(this.form)
            .then(res=>{
                if(res.data.code == 200){
                    this.conent = res.data.data;
                    this.dialogFormVisible = false;
                }
            })
        }
    }
}
</script>

<style lang="less" scoped>
.box{
    width: 100%;
    height: 100vh;
    background-image:url(../../assets/rlhy.png);
    background-size: 100% 100%;
    background-repeat: no-repeat;
    position: relative;
    overflow: hidden;
}
.TimesAll{
    position: absolute;
    top: 3%;
    right: 2%;
    z-index: 80;
    text-align: left;
    span{
        font-size: 20px;
        color: #fff;
    }
    .buttonbox{
        width: 60px;
        padding-right:5%;
        box-sizing: border-box;
        float: right;
        text-align: right;
        font-size: 22px;
        color: #fff;
        cursor: pointer;
    }
}
.boxone{
    width: 100%;
    height: 100%;
    padding:0 36px;
    box-sizing: border-box;
    // background: #FFECBB;
    text-align: center;
    cursor: pointer;
    position: relative;

}
// 屏幕尺寸大于1900
@media screen and (min-width:1900px) {
    .boxoneh1{
        width: 100%;
        font-size: 160px;
        line-height:260px;
        color: #FFECBB;
        font-weight: bold;
        text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
    }
    .boxb{
        width: 100%;
        font-size: 108px;
        line-height:125px;
        color: #FFECBB;
        font-weight: bold;
        text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
        cursor: pointer;
        margin-top:8%;
    }
    .boxp{
        width: 100%;
        height: 200px;
        text-align: center;
        position: absolute;
        bottom: 0;
        left:0;
        p{
            font-size: 96px;
            line-height:120px;
            color: #FFECBB;
            font-weight: bold;
            text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
            font-weight: 600;
        }
        span{
            font-size: 64px;
            font-weight: bold;
            color: #FFECBB;
            line-height: 80px;
            text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
        }
    }
}
// 屏幕尺寸大于1370小于1900
@media screen and (max-width:1900px) and (min-width:1370px) {
    .boxoneh1{
        width: 100%;
        font-size: 128px;
        line-height:260px;
        color: #FFECBB;
        font-weight: bold;
        text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
    }
    .boxb{
        width: 100%;
        font-size: 64px;
        line-height:80px;
        color: #FFECBB;
        font-weight: bold;
        text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
        cursor: pointer;
        margin-top:5%;
    }
    .boxp{
        width: 100%;
        height: 110px;
        text-align: center;
        position: absolute;
        bottom: 0;
        left:0;
        p{
            font-size: 48px;
            line-height:60px;
            color: #FFECBB;
            font-weight: bold;
            text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
            font-weight: 600;
        }
        span{
            font-size: 44px;
            font-weight: bold;
            color: #FFECBB;
            line-height: 50px;
            text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
        }
    }
}
// 屏幕尺寸小于1370
@media screen and (max-width:1370px) {
    .boxoneh1{
        width: 100%;
        font-size: 128px;
        line-height:260px;
        color: #FFECBB;
        font-weight: bold;
        text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
    }
    .boxb{
        width: 100%;
        font-size: 64px;
        line-height:80px;
        color: #FFECBB;
        font-weight: bold;
        text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
        cursor: pointer;
        margin-top:3%;
    }
    .boxp{
        width: 100%;
        height: 110px;
        text-align: center;
        position: absolute;
        bottom: 0;
        left:0;
        p{
            font-size: 48px;
            line-height:60px;
            color: #FFECBB;
            font-weight: bold;
            text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
            font-weight: 600;
        }
        span{
            font-size: 40px;
            font-weight: bold;
            color: #FFECBB;
            line-height: 50px;
            text-shadow: -1px -1px 5px black, 0 0 10px yellow, 0 0 8px yellow;
        }
    }
}
</style>
