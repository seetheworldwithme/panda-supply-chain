<template>
	<view class="content">
        
		<view style="height: 500rpx">111</view>
        <button id="btn1" class="btn">打开App</button>
		<!-- 此处的APPID为移动应用APPID -->
		<wx-open-launch-app appid="wxb1bea91c99acb227" extinfo="" @launch="AppLaunch" @error="AppError">
            <button id="btn1" class="btn">打开App</button>
			<script type="text/wxtag-template">
				<style>.btn { margin-bottom:200rpx;line-height: 34rpx; font-size: 34rpx; color:rgba(255,255,255,1); background:rgba(33,186,189,1); border: none; outline: none;font-family:PingFangSC-Regular,PingFang SC;font-weight:400;}</style>
		  		<button id="btn1" class="btn">打开App</button>
		  		</script>
		</wx-open-launch-app>
	</view>
</template>

<script type="text/wxtag-template">
	import wx from 'weixin-js-sdk'
    import sha1 from 'js-sha1' //mian.js


	export default {
		data() {
			return {
				time:'',
				nums: '',
				signature:'',
				appId: ''
			}
		},
		onLoad() {

            // wx.config({
            //     debug: true, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
            //     appId: '', // 必填，公众号的唯一标识
            //     timestamp: '', // 必填，生成签名的时间戳
            //     nonceStr: '', // 必填，生成签名的随机串
            //     signature: '',// 必填，签名
            //     jsApiList: [
                    
            //     ] // 必填，需要使用的JS接口列表
            // });


            // let sha1 = require('js-sha1');//vue
            // this.signature = sha1(stringA)
            // var chars = ['0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z'];
            // var nums="";
            // for(var i=0;i<32;i++){
            //     var id = parseInt(Math.random()*35);
            //     nums+=chars[id];
            // }
            // this.nums = nums;
            // //得到日期随机数 timestamp
            // let timestamp = Math.round(new Date() / 1000)
            // this.time = timestamp.toString();
            // this.appid = 'wx0ff59ec3274e0a53';
			// let url = window.location.href;
			// console.log(url)
			// url = url.split('#')[0]
			// console.log(url)
		
            // wx.config({
            //     debug: true,
            //     appId: that.appid,  // 此处APPID为公众号的APPID
            //     timestamp: that.time, 
            //     nonceStr: that.nums, 
            //     signature: that.signature,
            //     jsApiList: ['onMenuShareTimeline'], //此处必须填写一个，不为null就行
            //     openTagList: ['wx-open-launch-app'] 
            // })
            // wx.ready(function () {
            //     console.log('ready')
            //     wx.checkJsApi({
            //         jsApiList: ['wx-open-launch-app'], // 校验跳转APP的标签是否可用
            //         success: function (res) {
            //         console.log('可用')
            //         console.log('config111')
            //         },
            //         fail: (err) => {
            //         console.log(err, '不可用')
            //         }
            //     })
            // }) 
            // wx.error(function(res){
            //     console.log("error：",res)
            // });

		},
		methods: {
			AppLaunch(e){
				console.log("success",e)
			},
			AppError(e){
				console.log("fail",e.detail)
			},
			// 监听error 函数
			handleErrorFn (e) {
			  // 跳转失败
			  console.log(JSON.stringify(e), '跳转失败')
			  download() // 这里跳转App下载页面
			},
			// 监听launch 函数
			handleLaunchFn (e) {
			  // 跳转成功
			  console.log(JSON.stringify(e), '跳转成功')
			}
		}
	}
</script>

<style>
	.content {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
        size: 24px;
        color: #333;
	}

	.logo {
		height: 200rpx;
		width: 200rpx;
		margin-top: 200rpx;
		margin-left: auto;
		margin-right: auto;
		margin-bottom: 50rpx;
	}

	.text-area {
		display: flex;
		justify-content: center;
	}
</style>

