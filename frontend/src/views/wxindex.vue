<template>
  <div>
    <wx-open-launch-app id="launch-btn" appid="wx0ff59ec3274e0a53" class="wx-open-launch-app-container" extinfo="id=2&pid=3">
      <script type="text/wxtag-template">
        <style>
          .btn{
            color:red;
          }
        </style>
        <button class="btn">按钮</button>
      </script>
    </wx-open-launch-app>
  </div>
</template>
<script>
import wx from "weixin-js-sdk";
export default {
  data() {
   
  },
  mounted() {
    let url = "https://zsc.coalcloudtech.com/wxindex";
    this.getData(url);
  },
  methods: {
    getData(url) {
      this.api.wxChat
        .getWxConfig({
          url: url,
        })
        .then((res) => {
          if (res.data.code == 200) {
            console.log(res.data.data);
            wx.config({
              debug: true,
              appId: res.data.data.appId, // 此处APPID为公众号的APPID
              timestamp: res.data.data.timestamp,
              nonceStr: res.data.data.nonceStr,
              signature: res.data.data.signature,
              jsApiList: ['onMenuShareTimeline','onMenuShareAppMessage','updateAppMessageShareData', 'updateTimelineShareData'], //此处必须填写一个，不为null就行
              openTagList: ["wx-open-launch-app"],
            });
            wx.ready(function (sender) {
              alert("ready",JSON.stringify(sender));
              wx.checkJsApi({
                jsApiList: ["wx-open-launch-app"], // 校验跳转APP的标签是否可用
                success: function (res) {
                  // alert(JSON.stringify(res),'可用')
                  console.log("可用", res);
                  console.log("config111");
                  
                },
                fail: (err) => {
                  alert(JSON.stringify(err),'bu可用')
                  console.log(err, "不可用");
                },
              });
            });
            wx.error(function (err) {
              console.log("error：", err);
              alert(JSON.stringify(err),'wx.err')
            });
          }
        });
    },
  },
};
</script>
  
<style scoped>
.content{
  width: 100%;
  height: 100%;
}
  .wx-open-launch-app-container{
    display: block;
    width: 200px !important;
    height: 200px;
    background-color: skyblue;
  }
.btn-open-weapp {
  background: linear-gradient(to right, #ffd52e 0%, #ffef93 50%, #ffd52e 100%);
  border: 0;
  color: #424242;
  text-shadow: 0px 1px 1px #fff;
  border-radius: 50px;
  text-align: center;
  width: 200px;
  height: 160px;
  line-height: 35px;
  outline: none;
}
.goods {
  text-align: center;
  min-height: 100vh;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-box-pack: center;
  background: #0a045b;
  position: relative;
}
.login-button {
  margin-top: 20px;
  width: 100%;
  height: 50px;
  border-radius: 25px;
  border: none;
  font-weight: bold;
  color: white;
  background: #1d2473;
}

.active {
  margin-top: 15px;
  z-index: 2;
  background: rgba(255, 255, 255, 0.9);
  width: 80vw;
  margin-left: 10vw;
  border-radius: 10px;
  padding: 20px 20px 30px;
  box-sizing: border-box;
}

* {
  margin: 0;
  padding: 0;
  border: 0;
  overflow: hidden;
}

body {
  width: 100%;
  height: 100vh;
}

.content {
  width: 100%;
  height: calc(100vh - 150px);
  background: linear-gradient(#008b8a, #fff, #fff) no-repeat;
  overflow-x: hidden;
  overflow-y: scroll;
}

.header {
  width: 100%;
  text-align: center;
}

.header img {
  margin-top: 80px;
  width: 380px;
  height: 164px;
}

.header .tips {
  color: #fff;
  margin-top: 20px;
  font-size: 40px;
}

.button {
  text-align: center;
  position: fixed;
  bottom: 20px;
  margin-left: 2.5vw;
}

button {
  width: 95vw;
  color: #fff;
  height: 120px;
  font-size: 50px;
  border-radius: 15px;
}

.title {
  width: 95vw;
  height: 150px;
  border-radius: 15px;
  background-color: #fff;
  margin: 60px 2.5vw;
  display: flex;
  line-height: 150px;
  box-shadow: 0 0 2px 1px #dbdbdb;
}

.title img {
  width: 120px;
  height: 120px;
  margin-top: 15px;
  margin-left: 30px;
}

.title .desc {
  margin-left: 30px;
  font-size: 40px;
  color: rgb(0, 138, 139);
  font-weight: bold;
}

.card {
  position: relative;
  font-weight: bold;
  font-size: 40px;
  width: 95vw;
  margin-top: 40px;
  margin-left: 2.5vw;
  border-radius: 15px;
  background-color: #fff;
  box-shadow: 0 0 2px 1px #dbdbdb;
  margin-bottom: 30px;
}

.time {
  margin-left: 60px;
  height: 150px;
  line-height: 150px;
}

.time span:nth-child(2) {
  margin-left: 40px;
}

.card-item {
  width: 85vw;
  height: 220px;
  margin-left: 60px;
  display: flex;
}

.circle {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background-color: rgb(217, 223, 232);
  margin-top: 15px;
}

.line {
  position: absolute;
  width: 10px;
  height: 195px;
  background-color: rgb(217, 223, 232);
  margin-left: 10px;
}

.price {
  text-align: right;
  margin-right: 80px;
  margin-bottom: 40px;
  font-size: 46px;
  font-weight: normal;
}

.cus-button {
  width: 90vw;
  height: 120px;
  line-height: 120px;
  text-align: center;
  color: #fff;
  margin-left: 2.5vw;
  margin-bottom: 40px;
  background-color: rgb(0, 138, 139);
  border-radius: 15px;
  font-weight: normal;
}
#alert {
  display: none;
  position: relative;
  background: rgba(0, 0, 0, 0.5);
  width: 100vw;
  height: 100vh;
  z-index: 100;
}
.alert {
  position: absolute;
  width: 95vw;
  height: 580px;
  left: 2.5vw;
  top: 50%;
  margin-top: -300px;
  border-radius: 15px;
  background: #fff;
  box-shadow: 0 0 8px 4px #dbdbdb;
}
.alert .tips {
  text-align: center;
  font-size: 50px;
  font-weight: bold;
  height: 100px;
  line-height: 100px;
  border-bottom: 1px solid #9e9e9e;
}
.alert .cont {
  font-size: 50px;
  text-align: center;
  color: #8a8a8a;
  height: 150px;
  margin-top: 100px;
  margin-bottom: 100px;
}
.alert .btn {
  border-top: 1px solid #9e9e9e;
  font-size: 40px;
  display: flex;
  text-align: center;
  height: 128px;
}
.alert .btn div {
  width: 50vw;
  line-height: 128px;
}
.alert .btn div:nth-child(1) {
  border-right: 1px solid #9e9e9e;
  color: #8a8a8a;
}
.alert .btn div:nth-child(2) {
  color: rgb(0, 138, 139);
}

a {
  display: block;
  width: 100px;
  height: 100px;
  background-color: skyblue;
  line-height: 100px;
}
</style>
  
  