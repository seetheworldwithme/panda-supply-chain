<template>
  <el-container>
    <el-header style="height:auto;">
      <!-- header导航菜单开始 -->
      <el-menu
        :default-active="activeIndex"
        class="menus-top"
        mode="horizontal"
        background-color="#CF1724"
        text-color="#fff"
        active-text-color="#FFA6AC"
        :router="true"
      >
        <!-- 个人中心 -->
        <el-dropdown class="menus-right" @command="oldLoginDome">
          <span class="dropdown-link">
            {{ username }}<i class="el-icon-arrow-down el-icon--right"></i>
          </span>
          <el-dropdown-menu slot="dropdown">
            <!-- <el-dropdown-item
              command="a"
              icon="el-icon-aliicon-personal"
              style="width: 100px; text-align: left;font-size:14px;"
              > 个人资料</el-dropdown-item> -->
            <el-dropdown-item
              command="b"
              @click.native="dialogFormVisible = true"
              icon="el-icon-aliicon-wangjimima"
              style="width: 100px; text-align: left"
              > 修改密码</el-dropdown-item
            >
            <el-dropdown-item
              command="c"
              icon="el-icon-alituichu"
              style="width: 100px; text-align: left"
              >  退 出</el-dropdown-item
            >
          </el-dropdown-menu>
        </el-dropdown>
        <!-- 标题 -->
        <label class="logosize">{{ logo }}</label>
        <div style="float:right;">
          <!-- 帮助开始 -->
          <!-- <el-dropdown class="drop-help" @command="oldLoghelp">
            <span class="drop-link"><i class="el-icon-alihelp"></i>帮 助</span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="a">常见问题</el-dropdown-item>
              <el-dropdown-item command="b">帮助文档</el-dropdown-item>
              <el-dropdown-item command="c">更新日志</el-dropdown-item>
              <el-dropdown-item command="d">反馈意见</el-dropdown-item>
              <el-dropdown-item command="e">客服中心</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown> -->
          <!-- 帮助结束 -->
          <!-- 消息开始 -->
          <!-- <div class="menus-news" @click="oldlognews">
            <span class="news-link"><i class="el-icon-alixiaoxi1"></i>消 息</span>
          </div> -->
          <!-- 消息结束 -->
          <!-- 顶部路由菜单循环遍历 -->
          <!-- @click="newRoutes(v.children)" -->
          <el-menu-item
            v-for="(v, i) in routes"
            :key="i"
            :index="v.name"
            :route="{ name: v.name }"
            style="border:0 !important;">
            <i :class="v.meta.icon" style="color: #fff; font-size: 16px;vertical-align:middle;"></i>{{v.meta.title}}
          </el-menu-item>
        </div>
      </el-menu>
      <!-- header导航菜单结束 -->
    </el-header>
    <!-- 内容区域开始 -->
    <el-main>
      <!-- 内容视口view渲染 -->
      <router-view></router-view>
    </el-main>
    <!-- 内容区域结束 -->
    <!-- 消息一级抽屉 -->
    <el-drawer
      size="350px"
      :visible.sync="drawer"
      :with-header="false"
      :append-to-body="true"
    >
      <div class="concentBox">
        <!-- tab选项卡 -->
        <div class="menu-list">
          <div
            class="totalTabs"
            :class="{ active: item.id === sel }"
            v-for="item in tabs"
            :key="item.id"
            @click="select(item)"
          >
            <span>{{ item.label }}</span>
            <div
              class="stripnone"
              :style="item.width"
              :class="{ strip: item.id === sel }"
            ></div>
          </div>
        </div>
        <div class="text">
          <div class="textAll" v-show="sel === 1">
            <span class="onsize" @click="emptynews">一键已读1</span>
            <ul class="newsAll">
              <li
                v-for="(v, kitme) in Fullnotice"
                :key="kitme"
                :class="{ Active: yesactive == kitme }"
                @click="active(kitme)"
              >
                <div class="new-left">
                  <span>{{ v.title }}管理 <el-badge :value="v.value" /></span>
                  <p v-if="v.value == ''">有0条新的合同通知</p>
                  <p v-else>有{{ v.value }}条新的合同通知</p>
                </div>
                <div class="new-right">{{ v.time }}</div>
              </li>
            </ul>
          </div>
          <div class="textAll" v-show="sel === 2">
            <span class="onsize">一键已读2</span>
            <!-- <ul class="newsAll">
              <li v-for="(v,k) in Fullnotice" :key="k" :class="{Active: yesactive == k}" @click="active(k)">
                <div class="new-left">
                  <span>{{v.title}}管理 <el-badge class="mark" :value="v.value"/></span>
                  <p v-if="v.value ==''">有0条新的合同通知</p>
                  <p v-else>有{{v.value}}条新的合同通知</p>
                </div>
                <div class="new-right">{{v.time}}</div>
              </li>
            </ul> -->
          </div>
        </div>
      </div>
    </el-drawer>
    <!-- 消息二级抽屉 -->
    <el-drawer
      size="300px"
      :append-to-body="true"
      :visible.sync="innerDrawer"
      :with-header="false"
      :modal="false"
      style="
        margin-right: 350px !important;
        box-shadow: 2px 0 10px rgba(102, 102, 102, 0.15);
      "
    >
      <dl class="Secon">
        <dd v-for="(v, itme) in Secon" :key="itme">
          <p>{{ v.title }}</p>
          <el-button
            v-for="(val, kls) in v.button"
            :key="kls"
            size="mini"
            :label="val.lable"
            >{{ val.value }}</el-button
          >
        </dd>lllllllllll
        <dd v-for="(v, kindex) in Secondate" :key="kindex">
          <p>{{ v.title }}</p>
          <el-button
            v-for="(val, kd) in v.button"
            :key="kd"
            size="mini"
            :label="val.lable"
            >{{ val.value }}</el-button
          >
        </dd>
        <dd v-for="(v, index) in Seconfrom" :key="index">
          <p>{{ v.title }}</p>
          <el-button
            v-for="(val, ka) in v.button"
            :key="ka"
            size="mini"
            :label="val.lable"
            >{{ val.value }}</el-button
          >
        </dd>
      </dl>
      <ul class="seNav">
        <span class="senavsize"  @click="emptynav">一键消息已读</span>
        <li v-for="(valu,vindex) in senavNotice" :key="vindex" :class="{ Active: Navyesactive == vindex }" @click="Navactive(vindex)">
          <div class="seNav-com">
            <i class="el-icon-chat-dot-round"></i>
            <el-badge :is-dot="valu.dot" ><b>{{valu.title}}</b></el-badge>
            <span>{{valu.time}}</span>
          </div>
          <p>{{valu.value }}</p>
        </li>
      </ul>
    
    </el-drawer>
    <!-- 修改密码弹窗 -->
    <el-dialog title="修改密码" :visible.sync="dialogFormVisible" width="35%" :before-close="BuyerhandleClose">
      <div style="padding:0 30px;">
        <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="100px" label-position="left">
          <el-form-item v-if="accountType == true" label="选择账号" prop="id">
            <el-select v-model="ruleForm.id" placeholder="请选择账号类型" style="width:100%;" size="small">
                <el-option
                v-for="item in projectName"
                :key="item.value"
                :label="item.label"
                :value="item.value">
                </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="原密码" prop="password">
            <el-input size="small" v-model="ruleForm.password" show-password placeholder="如忘记原密码,请在登录页使用手机验证码登录,即可重置密码."></el-input>
          </el-form-item>
          <el-form-item label="新密码" prop="newPassword">
            <el-input size="small" v-model="ruleForm.newPassword" show-password placeholder="请输入新密码"></el-input>
          </el-form-item>
          <el-form-item label="确认新密码" prop="checkPass">
            <el-input size="small" v-model="ruleForm.checkPass" show-password placeholder="请再次确认新密码"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitForm('ruleForm')" size="small" style="float:right;">确 定</el-button>
            <el-button @click="resetForm('ruleForm')" size="small" style="float:right;margin-right:10px;">取 消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </el-container>
</template>
<script>
export default {
  data() {
    return {
      accountType:false,      // 选择账号
      // 修改密码
      ruleForm: {
        id:'',  // 账号类型
        password: '',         // 原密码
        newPassword: '',      // 新密码
        checkPass: '',        // 二次新密码
      },
      rules: {
        id: [
          { required: true, message: '请选择账号', trigger: 'blur'}
        ],
        password: [
          { required: true, message: '请输入原密码', trigger: 'blur'},
          { pattern: /(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,16}$/, message: '请输入正确格式的原密码', trigger: 'blur'},
        ],
        newPassword: [
          { required: true, message: '请输入新密码', trigger: 'blur'},
          { pattern: /(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,16}$/, message: '密码设置为 8 到 16 位至少包含数字跟字母', trigger: 'blur'},
        ],
        checkPass: [
          { required: true, message: '请再次输入新密码', trigger: 'blur'},
          { pattern: /(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,16}$/, message: '密码设置为 8 到 16 位至少包含数字跟字母', trigger: 'blur'},
        ]
      },
      // 账号类型
      projectName:[], 
      dialogFormVisible: false,
      activeIndex:'Mindex',
      logo: "",
      username: "",
      //导航路由
      routes: [],
      drawer: false, // 抽屉开关
      innerDrawer: false, //里面抽屉
      // activenewsName:'first', // tab
      sel: 1, // 选项卡默认显示tab
      tabs: [
        { label: "全部通知", id: 1, width: "width:65px;" },
        { label: "未读", id: 2, width: "width:33px;" },
      ],
      Fullnotice: [
        {
          title: "合同",
          value: "12",
          time: "2021-03-01 23:58",
        },
        {
          title: "物流",
          value: "5",
          time: "2021-03-02 08:58",
        },
        {
          title: "收支",
          value: "2",
          time: "2021-01-09 05:58",
        },
        {
          title: "往来",
          value: "",
          time: "2021-02-10 12:58",
        },
      ], // 消息通知
      yesactive: "-1",
      Secon: [
        {
          title: "状态筛选器",
          button: [
            {
              label: "1",
              value: "全部",
            },
            {
              label: "2",
              value: "未读",
            },
            {
              label: "3",
              value: "已读",
            },
          ],
        },
      ],
      Secondate: [
        {
          title: "时间筛选器",
          button: [
            {
              label: "1",
              value: "全部",
            },
            {
              label: "2",
              value: "今日",
            },
            {
              label: "3",
              value: "本周",
            },
            {
              label: "4",
              value: "本月",
            },
          ],
        },
      ],
      Seconfrom: [
        {
          title: "表单筛选器",
          button: [
            {
              label: "1",
              value: "全部",
            },
            {
              label: "2",
              value: "智能助手",
            },
          ],
        },
      ],
      senavNotice: [
        {
          title: "新增合同订单",
          dot: true,
          time: "1分钟前",
          value:"刘娜 新增了 一个合同HTCG.0000001202102220001"
        },
        {
          title: "新增合同订单",
          dot: false,
          time: "10分钟前",
          value:"刘娜 新增了 一个合同HTCG.0000001202102220001"
        },
        {
          title: "新增合同订单",
          dot: true,
          time: "2021-01-09 05:58",
          value:"刘娜 新增了 一个合同HTCG.0000001202102220001"
        }
      ], // 子消息通知
      Navyesactive:'-1',
    };
  },
  created() {
    // this.$store.dispatch("AllBasics");
    //获取动态路由表
    let lengths = this.$router.options.routes.length - 1;
    this.routes = this.$router.options.routes[lengths].children;
    this.activeIndex = this.$route.matched[1].name;
    //回显企业名称
    this.logo = sessionStorage.getItem('companyName');
    //回显个人中心标题
    this.username = sessionStorage.getItem('username');
    // 是否显示账号类型
    // console.log(sessionStorage.getItem('accountType'))
    if(sessionStorage.getItem('level') == 2 && sessionStorage.getItem('accountTypes') == 2){
      this.accountType = true;
      // 查询非企业账号列表
      this.api.authAll.ptlist()
      .then(res=>{
        if(res.data.code == 200){
          for(let key in res.data.data){
            let obj = {
              value: res.data.data[key].id,
              label: res.data.data[key].username + ' - ' + res.data.data[key].account
            }
            this.projectName.push(obj)
          } 
        }
      })
    }else{
      this.accountType = false;
    }
    // 全部基础维护数据
    this.api.basicslist.list() 
    .then(res=>{
      if(res.data.code == 200){
        var newdata = res.data.data;
        var obj = {};                 // 全部基础维护信息
        var price_type = {};          // 价格类型
        var payment_method = {};      // 付款方式
        var transport_mode = {};      // 运输模式
        var carrier = {};             // 承运工具
        var payment_type = {};        // 支付方式
        var pro_type = {};            // 产品大类
        var pro__middle_type = {};    // 产品小类
        var pro__small_type = {};     // 标的指标
        var logistics_platform = {};  // 物流平台
        var file_relation = {};       // 文件关联类型
        var purpose = {};             // 用途
        // console.log(newdata)
        for(let i=0;i<newdata.length;i++){
          if(newdata[i].dtype == "price_type"){
              price_type[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "payment_method"){
              payment_method[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "transport_mode"){
              transport_mode[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "carrier"){
              carrier[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "payment_type"){
              payment_type[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "pro_type"){
              pro_type[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "pro__middle_type"){
              pro__middle_type[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "pro__small_type"){
              pro__small_type[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "logistics_platform"){
              logistics_platform[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "file_relation"){
              file_relation[newdata[i].dvalue] = newdata[i].dkey
          }else if(newdata[i].dtype == "purpose"){
              purpose[newdata[i].dvalue] = newdata[i].dkey
          }
        }
        obj.price_type = price_type
        obj.payment_method = payment_method
        obj.transport_mode = transport_mode
        obj.carrier = carrier
        obj.payment_type = payment_type
        obj.pro_type = pro_type
        obj.pro__middle_type = pro__middle_type
        obj.pro__small_type = pro__small_type
        obj.logistics_platform = logistics_platform
        obj.file_relation = file_relation
        obj.purpose = purpose
        sessionStorage.setItem('data',JSON.stringify(obj))
      }
    })
  },
  beforeUpdate() {},
  methods: {
    // 修改密码提交
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          if(this.ruleForm.checkPass != this.ruleForm.newPassword){
            this.$message({
              message: '两次新密码不一致，请重新设置',
              type: 'warning'
            });
          }else{
            this.api.updatePwd.updatePwd(this.ruleForm)
            .then(res=>{
              if(res.data.code == 999){
                  this.$message({
                    message: '修改失败,原密码错误',
                    type: 'error'
                  });
              }else if(res.data.code == 200){
                this.$message({
                    type: 'success',
                    message: "修改密码成功,请重新登录",
                    duration:500,
                    onClose:()=>{
                      sessionStorage.clear()
                      this.$router.push({name:'Login'});
                    }
                });
              }
            })

          }
        }
      });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
      this.dialogFormVisible = false;
    },
    BuyerhandleClose(){
      this.dialogFormVisible = false;
      this.ruleForm = {}
    },
    //消息回调
    oldlognews() {
      // this.news = 'oldnews'
      this.drawer = true;
    },
    //消息回调 选项卡
    select(item) {
      this.sel = item.id;
      if(item.id == 2){
        this.innerDrawer = false;
      }
    },
    emptynews() {
      for (let i = 0; i < this.Fullnotice.length; i++) {
        this.Fullnotice[i].value = "";
      }
    },
    //消息回调 选中
    active(index) {
      this.yesactive = index;
      this.innerDrawer = true;
    },
    //子消息回调
    emptynav(){
      for (let i = 0; i < this.senavNotice.length; i++) {
        this.senavNotice[i].dot = false;
      }
    },
    //消息回调 选中
    Navactive(index) {
      this.Navyesactive = index;
    },
    // 帮助回调
    oldLoghelp(command) {
      if (command === "a") {
        window.open("http://www.baidu.com");
      }
      if (command === "b") {
        window.open("http://www.taobao.com");
      }
      if (command === "c") {
        window.open("http://www.jd.com");
      }
      if (command === "d") {
        window.open("http://www.baidu.com");
      }
      if (command === "e") {
        window.open("http://www.taobao.com");
      }
    },
    //个人中心回调
    oldLoginDome(command) {
      if (command === "a") {
      }
      if (command === "b") {
      }
      if (command === "c") {
        sessionStorage.clear();
        this.$router.push({name:'Login'});
        this.$router.go(0);
      }
    },
  },
};
</script>
<style lang="less" scoped>
.boxss{
    width: 100%;
    height: calc(100vh - 60px);
    background-color: #fdfd;
    position: relative;
    .hbox{
        width: 300px;
        height: 30px;
        position: absolute;
        top: 0;
        left: 0;
        bottom: 0;
        right:0;
        margin: auto;
    }
}
.logosize{
  outline:none;
}
// 选项卡
.concentBox {
  width: 100%;
  position: relative;
}
.totalTabs {
  width: 64px;
  float: left;
  font-size: 16px;
  line-height: 26px;
  font-weight: 650;
  color: #333;
  text-align: center;
  margin: 31px 0 0 35px;
  cursor: pointer;
}
.totalTabs.active {
  color: #cf1724;
}
.stripnone {
  height: 3px;
  background-color: #cf1724;
  margin: 7px auto 16px;
  display: none;
}
.strip {
  display: block;
}
.text {
  position: absolute;
  left: 0;
  top: 100px;
  width: 100%;
  .onsize {
    font-size: 14px;
    color: #cf1724;
    position: absolute;
    top: -67px;
    right: 24px;
    cursor: pointer;
  }
  // 消息列表样式
  .newsAll {
    li {
      width: 100%;
      overflow: hidden;
      padding: 12px 24px;
      cursor: pointer;
      .new-left {
        width: 50%;
        float: left;
        color: #333;
        span {
          font-size: 18px;
          font-weight: 600;
        }
        p {
          font-size: 14px;
          line-height: 30px;
        }
      }
      .new-right {
        float: right;
        width: 50%;
        color: #828282;
        line-height: 60px;
      }
    }
    li:hover {
      border-left: 3px solid #cf1724;
      box-shadow: 2px 2px 10px rgba(128, 146, 192, 0.15);
      transition: all 0.1s;
    }
    .Active {
      border-left: 3px solid #cf1724;
      box-shadow: 2px 2px 10px rgba(128, 146, 192, 0.15);
    }
  }
}
// 子消息列表样式
.Secon {
  dd {
    padding: 0 16px 16px;
    border-bottom: 1px solid #e0e0e0;
    p {
      font-size: 12px;
      color: #4f4f4f;
      line-height: 35px;
    }
    .el-button {
      margin-bottom: 8px;
    }
  }
}
.seNav {
  position: relative;
  .senavsize{
    font-size: 14px;
    color: #cf1724;
    position: absolute;
    top: -50px;
    right: 8px;
    cursor: pointer;
  }
  li {
    width: 100%;
    overflow: hidden;
    padding: 16px 12px 8px;
    cursor: pointer;
    .seNav-com {
      i {
        font-size: 16px;
        color: #333;
      }
      b {
        font-size: 16px;
        color: #333;
        margin-left: 20px;
      }
      span {
        font-size: 12px;
        float: right;
        color: #828282;
        margin-right: 25px;
      }
    }
    p {
      font-size: 12px;
      line-height: 20px;
      color: #333;
      margin:10px 0 0 35px;
    }
  }
  li:hover {
    border-left: 3px solid #cf1724;
    box-shadow: 2px 2px 10px rgba(128, 146, 192, 0.15);
    transition: all 0.1s;
  }
  .Active {
      border-left: 3px solid #cf1724;
      box-shadow: 2px 2px 10px rgba(128, 146, 192, 0.15);
  }
}
// header 开始
.el-menu{
  border: 0;
}
.el-header {
  padding: 0;
  label {
    font-size: 14px;
    line-height: 60px;
    color: #fff;
    float: left;
  }
  .el-submenu {
    float: right;
  }
  .el-menu-item {
    float: left;
  }
  .menus-top {
    width: 100%;
    padding: 0 20px;
    box-sizing: border-box;
    background-color: #cf1724;
    .el-menu-item {
      text-align: center;
    }
    .menus-right {
      float: right;
      margin: 10px 0 0 90px;
      outline: none;
      cursor: pointer;
      .dropdown-link {
        display: inline-block;
        font-size: 14px;
        line-height: 40px;
        color: #fff;
      }
    }
    .drop-help {
      width: 100px;
      height: 60px;
      text-align: center;
      float: right;
      cursor: pointer;
      outline:none;
      .drop-link {
        font-size: 14px;
        line-height: 60px;
        vertical-align:middle; 
        color: #fff;
        i{
          font-size: 16px;
          vertical-align:middle; 
          margin-right: 8px;
        }
      }
    }
    .drop-help:hover {
      background-color: rgba(255, 166, 172, 0.3);
    }
    .menus-news {
      width: 100px;
      height: 60px;
      text-align: center;
      float: right;
      cursor: pointer;
      outline:none;
      .news-link {
        font-size: 14px;
        line-height: 60px;
        color: #fff;
        i{
          font-size: 16px;
          margin-right: 8px;
          vertical-align:middle; 
        }
      }
    }
    .menus-news:hover {
      background-color: rgba(255, 166, 172, 0.3);
    }
  }
  .el-menu--horizontal .el-menu-item:hover {
    background-color: rgba(255, 255, 255, 0.2) !important;
  }
  .el-menu-item.is-active {
    border-bottom: 0 !important;
    // background-color: #c10b18 !important;
    i {
      color: #ffa6ac !important;
    }
  }
}
// header 结束

// main内容开始
.el-main {
  background-color: #f5f5f5;
  padding: 0;
  overflow: hidden;
}
// mian内容结束
</style>