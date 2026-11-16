// const path = require('path');
// const HtmlWebpackPlugin = require('html-webpack-plugin');           // 自动生成html模板，以便在多个入口文件时，不用手动修改引入的js
// const webpack = require('webpack');
// const VueLoaderPlugin = require('vue-loader/lib/plugin');

module.exports = {
    // 准生产 ygls   生产 ssygl
    assetsDir: 'ssygl',                                            // 静态资源目录 (js, css, img, fonts)
    publicPath: '/',                                                // 根域上下文目录
    // chainWebpack: config => {
    //     //全局下的title
    //     config.plugin('html').tap(args => { 
    //         args[0].title = '海云煤炭资源'
    //         return args
    //      })
    // },
    
    // outputDir: 'hello', // 打包名称
};