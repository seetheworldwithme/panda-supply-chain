<template>
<div class="box">
    <el-empty style="padding:70% 0 0;box-sizing: border-box;" description="暂无数据！">
        <!-- <div>{{ceshiurl}}</div> -->
    </el-empty>
</div>
</template>
<script>
export default {
    data(){
        return {
            ceshiurl:'',
        }
    },
    created(){
                let a = window.location.href;
        
        sessionStorage.setItem('url',a);
        let url = sessionStorage.getItem('url');
        this.ceshiurl = url;
    }
}
</script>

<style lang="less" scoped>
 .box{
     width: 100%;
     height: 100vh;
     overflow: hidden;
 }
</style>