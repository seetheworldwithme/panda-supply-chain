<template>
    <div class="aymentyes">
        <div class="login-concent">
                <div class="img-lg">
                    <img src="../../src/assets/lg.svg" alt="">
                </div>
                <div class="img-size"><span>云管理</span></div>
        </div>
        <div class="PaymentyesBox">
            <i class="el-icon-success"></i>
            <h2> 支 付 成 功 ! </h2>
            <p> 温馨提示: 请返回启辰智链查看订单详情 </p>
        </div>
        <div class="copyright-bottom">
            <div class="copyright"><span>Copyright © 山西海云启辰科技有限公司 版权所有</span><p> 晋ICP备2022001434号</p></div>
        </div>
    </div>
</template>
<style scoped lang="less">
.aymentyes{
    width: 100%;
    height: 100vh;
    position: relative;
    overflow: hidden;
    .PaymentyesBox{
        width: 50%;
        height: 260px;
        padding: 2%;
        box-sizing: border-box;
        position: absolute;
        top: -25%;
        left: 0;
        right: 0;
        bottom: 0;
        margin: auto;
        text-align: center;
        i{
            font-size: 68px;
            display: inline-block;
            margin: auto;
            color: #6FCF97;
        }
        h2{
            font-size: 32px;
            line-height: 100px;
            color: #333;
            font-weight: 600;
        }
        p{
            font-size: 16px;
            color: #333;
        }
    }
}
// 头部开始
.login-concent{
    width: 100%;
    height: 76px;
    padding:8px 5%;
    box-sizing: border-box;
    border-bottom: 1px solid #e4e4e4;
    line-height: 60px;
    overflow: hidden;
    position: fixed;
    top: 0;
    left: 0;
    background-color: #fff;
    z-index: 10;
    .img-lg{
        float: left;
        img{
            width:48px;
            height: 48px;
            margin-top:6px;
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
</style>