<template>
    <!-- 暂无数据 -->
    <div class="main">
        <div class="error">
            <img src="../assets/zanwushuju.png" alt="">
            <p>暂无数据，请在左侧菜单栏操作！</p>
        </div>
    </div>

</template>
<script>
    export default{
        methods:{
         
        }
    }
</script>

<style lang="less" scoped>
.main{
    width: 100%;
    height: 100vh;
    overflow: hidden;
    position: relative;
    .error{
        width:353px;
        height: 200px;
        position: absolute;
        top: -20%;
        bottom:0;
        left: 0;
        right: 0;
        margin: auto;
    }
    img{
        width: 353px;
        height: 200px;
    }
    p{
        font-size: 22px;
        // font-weight: 600;
        text-align: center;
        line-height: 68px;
        color: #333;
    }
    h2{
        text-align: center;
        color: #333;
        cursor: pointer;
    }
}
</style>