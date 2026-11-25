<template>
  <div>
    <!-- tab选项卡 -->
    <div class="menu-list">
      <div
        class="totalTab"
        :class="{ active: item.id === sel }"
        v-for="item in tabs"
        :key="item.id"
        @click="select(item)"
        style="cursor: pointer"
      >
        <span>{{ item.label }}</span>
        <div
          class="stripnone"
          :style="item.width"
          :class="{ strip: item.id === sel }"
        ></div>
      </div>
    </div>
    <div class="parNavs">
      <!-- 发货地 -->
      <div v-show="sel === 1">
        <div class="parNavstrip">
          <el-col :span="13" class="par-right">
            <el-button size="small" @click="delethair" v-if="$route.meta.delete">删除</el-button>
            <el-button size="small" style="margin-left: 24px" @click="addhair" v-if="$route.meta.add"
              >新建</el-button>
            <el-input
              @keyup.enter.native="hairsearch"
              size="small"
              placeholder="请输入内容"
              prefix-icon="el-icon-search"
              v-model="hairpurSearch"
            >
            </el-input>
          </el-col>
          <!-- 发货地表格 -->
          <div class="pur-nav"
          v-loading="loading"
          element-loading-text = "加载中...">
            <div class="pur-table">
                <el-table
                  height="100%"
                  border
                  :data="hairtableData"
                  :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                  style="width:100%;"
                  @selection-change="hairChange"
                >
                  <el-table-column type="selection"></el-table-column>
                  <el-table-column v-if="$route.meta.edit"
                    label="操作"
                    width="100">
                    <template slot-scope="scope">
                      <el-button @click="Edithair(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
                    </template>
                  </el-table-column>
                  <el-table-column prop="id" v-if="show"> </el-table-column>
                  <el-table-column prop="type" v-if="show"> </el-table-column>
                  <el-table-column prop="addressName" label="地址名称">
                  </el-table-column>
                  <el-table-column prop="createTime" label="创建时间" sortable>
                  </el-table-column>
                  <!-- 暂无数据 -->
                  <div slot="empty" class="dataAvailable" v-if="dataAva">
                    <div class="dataAva">
                      <img src="../../assets/zanwushuju.png" alt="" />
                      <p>暂无数据，您可以新建合同试试</p>
                      <el-button
                        size="small"
                        style="margin-left: 24px"
                        @click="addhair" v-if="$route.meta.add"
                        >新建</el-button
                      >
                    </div>
                  </div>
                </el-table>
            </div>
            <!-- 分页 -->
            <div class="par-pages">
              <el-pagination
              background
              :page-sizes="[10, 30, 50, 100, 200]"
              :page-size.sync="hairsize"
              :current-page.sync="haircurrent"
              :pager-count="5"
              layout="total, sizes, prev, pager, next, jumper"
              :total="hairtotal"
              @current-change='hairgetpurcon'
              @size-change='hairgetpurcon'>
              </el-pagination>
            </div>
          </div>
        </div>
      </div>
       <!-- 收货地 -->
      <div v-show="sel === 2">
        <div class="parNavstrip">
          <el-col :span="13" class="par-right">
            <el-button size="small" @click="deletcollect" v-if="$route.meta.delete">删除</el-button>
            <el-button
              size="small"
              style="margin-left: 24px"
              @click="addhair"
               v-if="$route.meta.add">新建</el-button
            >
            <el-input
              @keyup.enter.native="collectsearch"
              size="small"
              placeholder="请输入内容"
              prefix-icon="el-icon-search"
              v-model="collectpurSearch"
            >
            </el-input>
          </el-col>
          <!-- 收货地表格 -->
          <div class="pur-nav"
          v-loading="loadinghair"
          element-loading-text = "加载中...">
            <div class="pur-table">
              <el-table
                height="100%"
                border
                :data="collecttableData"
                :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                style="width: 100%;"
                @selection-change="collectChange"
              >
                <el-table-column type="selection"> </el-table-column>
                <el-table-column v-if="$route.meta.edit"
                  label="操作"
                  width="100">
                  <template slot-scope="scope">
                    <el-button @click="Editcollect(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
                  </template>
                </el-table-column>
                <el-table-column prop="id" v-if="show"> </el-table-column>
                <el-table-column prop="type" v-if="show"> </el-table-column>
                <el-table-column prop="addressName" label="地址名称">
                </el-table-column>
                <el-table-column prop="createTime" label="创建时间" sortable>
                </el-table-column>
                <!-- 暂无数据 -->
                <div slot="empty" class="dataAvailable" v-if="dataAvahair">
                  <div class="dataAva">
                    <img src="../../assets/zanwushuju.png" alt="" />
                    <p>暂无数据，您可以新建合同试试</p>
                    <el-button
                      size="small"
                      style="margin-left: 24px"
                      @click="addhair"
                       v-if="$route.meta.add">新建</el-button
                    >
                  </div>
                </div>
              </el-table>
            </div>
            <!-- 分页 -->
            <div class="par-pages">
              <el-pagination
              background
              :page-sizes="[10, 30, 50, 100, 200]"
              :page-size.sync="collectsize"
              :pager-count="5"
              :current-page.sync="collectcurrent"
              layout="total, sizes, prev, pager, next, jumper"
              :total="collecttotal"
              @current-change='collectgetpurcon'
              @size-change='collectgetpurcon'>
              </el-pagination>
            </div>
          </div>
        </div>
      </div>
      <!-- 装货地 -->
      <div v-show="sel === 3">
        <div class="parNavstrip">
          <el-col :span="13" class="par-right">
            <el-button size="small" @click="deletpretend" v-if="$route.meta.delete">删除</el-button>
            <el-button
              size="small"
              style="margin-left: 24px"
              @click="addhair"
               v-if="$route.meta.add">新建</el-button
            >
            <el-input
              @keyup.enter.native="pretendsearch"
              size="small"
              placeholder="请输入内容"
              prefix-icon="el-icon-search"
              v-model="pretendpurSearch"
            >
            </el-input>
          </el-col>
          <!-- 装货地点表格 -->
          <div class="pur-nav"
              v-loading="loadingpretend"
              element-loading-text = "加载中...">
            <div class="pur-table">
              <el-table
                height="100%"
                border
                :data="pretendtableData"
                :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                style="width: 100%;"
                @selection-change="pretendChange">
                <el-table-column type="selection"> </el-table-column>
                <el-table-column v-if="$route.meta.edit"
                  label="操作"
                  width="100">
                  <template slot-scope="scope">
                    <el-button @click="Editpretend(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
                  </template>
                </el-table-column>
                <el-table-column prop="id" v-if="show"> </el-table-column>
                <el-table-column prop="type" v-if="show"> </el-table-column>
                <el-table-column prop="addressName" label="地址名称">
                </el-table-column>
                <el-table-column prop="createTime" label="创建时间" sortable>
                </el-table-column>
                <!-- 暂无数据 -->
                <div slot="empty" class="dataAvailable" v-if="dataAvapretend">
                  <div class="dataAva">
                    <img src="../../assets/zanwushuju.png" alt="" />
                    <p>暂无数据，您可以新建合同试试</p>
                    <el-button
                      size="small"
                      style="margin-left: 24px"
                      @click="addhair"
                       v-if="$route.meta.add">新建</el-button
                    >
                  </div>
                </div>
              </el-table>
            </div>
            <!-- 分页 -->
            <div class="par-pages">
              <el-pagination
              background
              :page-sizes="[10, 30, 50, 100, 200]"
              :page-size.sync="pretendsize"
              :pager-count="5"
              :current-page.sync="pretendcurrent"
              layout="total, sizes, prev, pager, next, jumper"
              :total="pretendtotal"
              @current-change='pretendgetpurcon'
              @size-change='pretendgetpurcon'>
              </el-pagination>
            </div>
          </div>
        </div>
      </div>
      <!-- 卸货地 -->
      <div v-show="sel === 4">
        <div class="parNavstrip">
          <el-col :span="13" class="par-right">
            <el-button size="small" @click="deletdischarge" v-if="$route.meta.delete">删除</el-button>
            <el-button
              size="small"
              style="margin-left: 24px"
              @click="addhair"
               v-if="$route.meta.add">新建</el-button
            >
            <el-input
              @keyup.enter.native="dischargesearch"
              size="small"
              placeholder="请输入内容"
              prefix-icon="el-icon-search"
              v-model="dischargepurSearch"
            >
            </el-input>
          </el-col>
          <!-- 卸货地点表格 -->
          <div class="pur-nav"
            v-loading="loadingpretendcharge"
            element-loading-text = "加载中...">
            <div class="pur-table">
              <el-table
                height="100%"
                border
                :data="dischargetableData"
                :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                style="width: 100%;"
                @selection-change="dischargeChange"
              >
                <el-table-column type="selection"> </el-table-column>
                <el-table-column v-if="$route.meta.edit"
                  label="操作"
                  width="100">
                  <template slot-scope="scope">
                    <el-button @click="Editdischarge(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
                  </template>
                </el-table-column>
                <el-table-column prop="id" v-if="show"> </el-table-column>
                <el-table-column prop="type" v-if="show"> </el-table-column>
                <el-table-column prop="addressName" label="地址名称">
                </el-table-column>
                <el-table-column prop="createTime" label="创建时间" sortable>
                </el-table-column>
                <!-- 暂无数据 -->
                <div slot="empty" class="dataAvailable" v-if="dataAvacharge">
                  <div class="dataAva">
                    <img src="../../assets/zanwushuju.png" alt="" />
                    <p>暂无数据，您可以新建合同试试</p>
                    <el-button
                      size="small"
                      style="margin-left: 24px"
                      @click="addhair"
                       v-if="$route.meta.add">新建</el-button
                    >
                  </div>
                </div>
              </el-table>
            </div>
            <!-- 分页 -->
            <div class="par-pages">
              <el-pagination
              background
              :page-sizes="[10, 30, 50, 100, 200]"
              :page-size.sync="dischargesize"
              :pager-count="5"
              :current-page.sync="dischargecurrent"
              layout="total, sizes, prev, pager, next, jumper"
              :total="dischargetotal"
              @current-change='dischargegetpurcon'
              @size-change='dischargegetpurcon'>
              </el-pagination>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      show: false,
      sel: '', // 选项卡默认显示tab
      tabs: [
        { label: "发货地", id: 1, width: "width:48px;" },
        { label: "收货地", id: 2, width: "width:48px;" },
        { label: "装货地", id: 3, width: "width:64px;" },
        { label: "卸货地", id: 4, width: "width:64px;" },
      ],

      /*
      发货地数据
      */
      dataAva:false,
      loading:false,
      haircurrent:1,                         // 当前页数
      hairsize:10,                           // 显示条数
      hairtotal:0,                           // 总条数
      hairdelmul:[],                         // 批量删除
      hairpurSearch: "",                     // 发货地搜索
      hairAddressval: [],                    // 表格选中数据
      hairtableData: [],                     // 发货地表格数据

      /*
      收货地数据
      */
      dataAvahair:false,
      loadinghair:false,
      collectcurrent:1,                      // 当前页数
      collectsize:10,                        // 显示条数
      collecttotal:0,                        // 总条数
      collectdelmul:[],                      // 批量删除
      collectpurSearch: "",                  // 收货地搜索
      collectAddressval: [],
      collecttableData: [],                  // 收货地表格数据

      /*
      装货地点数据
      */
      dataAvapretend:false,
      loadingpretend:false,
      pretendcurrent:1,                       // 当前页数
      pretendsize:10,                         // 显示条数
      pretendtotal:0,                         // 总条数
      pretenddelmul:[],                       // 批量删除
      pretendpurSearch: "",                   // 装货地点搜索
      pretendAddressval: [],
      pretendtableData: [],                   // 装货地点表格数据

      /*
      卸货地点数据
      */
      dataAvacharge:false,
      loadingpretendcharge:false,
      dischargecurrent:1,                     // 当前页数
      dischargesize:10,                       // 显示条数
      dischargetotal:0,                       // 总条数
      dischargedelmul:[],                     // 批量删除
      dischargepurSearch: "",                 // 卸货地点搜索
      dischargeAddressval: [],
      dischargetableData: [],                 // 卸货地点表格数据
    };
  },
  mounted(){
    if(this.$route.query.sel == 1){
      this.sel = 1;
      this.hairgetpurcon();
    }else if(this.$route.query.sel == 2){
      this.sel = 2;
      this.collectgetpurcon();
    }else if(this.$route.query.sel == 3){
      this.sel = 3;
      this.pretendgetpurcon();
    }else if(this.$route.query.sel == 4){
      this.sel = 4;
      this.dischargegetpurcon();
    }else{
      this.sel = 1;
      this.hairgetpurcon();
    }
  },
  methods: {
    //选项卡
    select(item) {
      this.hairpurSearch = "";
      this.collectpurSearch = "";
      this.pretendpurSearch = "";
      this.dischargepurSearch = "";
      this.sel = item.id;
      if(this.sel == 1){
        this.hairgetpurcon();
      }else if(this.sel == 2){
        this.collectgetpurcon();
      }else if(this.sel == 3){
        this.pretendgetpurcon();
      }else if(this.sel == 4){
        this.dischargegetpurcon();
      }
    },
    //新建维护地址
    addhair() {
      this.$router.push({ name: "Addhair" });
    },

    /*
    发货地方法
    */
    hairgetpurcon(){
      this.loading = true;
      // 传参
      var listByPageData = {
        type:2,
        searchStr:this.hairpurSearch,
        page:{
          current:this.haircurrent,
          size:this.hairsize
        }
      }
      // 渲染表格
      this.api.address.all(listByPageData)
      .then(res=>{
        this.loading = false;
        if(res.data.code == 200){
          if(res.data.data.records != ""){
            this.hairtableData = res.data.data.records;
            this.hairtotal = res.data.data.total;
          }else{
            this.dataAva = true;
          }
        }
      })
    },
    //发货地搜索功能
    hairsearch() {
      this.hairtableData = [];
      this.haircurrent = 1;
      this.hairgetpurcon();
    },
    hairChange(val) {
      // 发货地表格选中数据val
      this.hairAddressval = val;
      this.hairdelmul = [];
      for(let i=0;i<this.hairAddressval.length;i++){
          this.hairdelmul.push(this.hairAddressval[i].id) 
      }
    },
    //发货地 删除功能
    delethair() {
      if (this.hairAddressval.length == []) {
        this.$message({
          type: "warning",
          message: "请选择要删除的数据!",
        });
      } else {
        this.$confirm("此操作将永久删除该文件, 是否继续?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        })
          .then(() => {
            this.api.address.del(this.hairdelmul)
            .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                    type: 'success',
                    message: '删除成功!'
                    });
                    this.hairtableData = [];
                    this.haircurrent = 1;
                    this.hairgetpurcon();
                    if(this.hairtableData == []){
                      this.dataAva = true;
                    }
                }
            })
          })
          .catch(() => {
            this.$message({
              type: "info",
              message: "已取消删除",
            });
          });
      }
    },
    // 编辑发货地
    Edithair(row) {
      this.api.address.one({id:row.id})
      .then(res=>{
          if(res.data.code == 200){
          let purdata = res.data.data;
          purdata.id = row.id;
          sessionStorage.setItem('queryAll',JSON.stringify(purdata));
          this.$router.push({name:"Edithair"})
          }
      })
    },

    /*
    收货地方法
    */
    collectgetpurcon(){
      this.loadinghair = true;
      // 传参
      var listByPageData = {
        type:1,
        searchStr:this.collectpurSearch,
        page:{
          current:this.collectcurrent,
          size:this.collectsize
        }
      }
      // 渲染表格
      this.api.address.all(listByPageData)
      .then(res=>{
        this.loadinghair = false;
        if(res.data.code == 200){
          if(res.data.data.records != ""){
            this.collecttableData = res.data.data.records;
            this.collecttotal = res.data.data.total;
          }else{
            this.dataAvahair = true;
          }

        }
      })
    },
    //收货地搜索功能
    collectsearch() {
      this.collecttableData = [];
      this.collectcurrent = 1;
      this.collectgetpurcon();
    },
    collectChange(val) {
      this.collectAddressval = val;
      this.collectdelmul = [];
      for(let i=0;i<this.collectAddressval.length;i++){
          this.collectdelmul.push(this.collectAddressval[i].id) 
      }
    },
    //收货地 删除功能
    deletcollect() {
      if (this.collectAddressval.length == []) {
        this.$message({
          type: "warning",
          message: "请选择要删除的数据!",
        });
      } else {
        this.$confirm("此操作将永久删除该文件, 是否继续?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        })
          .then(() => {
            this.api.address.del(this.collectdelmul)
            .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                    type: 'success',
                    message: '删除成功!'
                    });
                    this.collecttableData = [];
                    this.collectcurrent = 1;
                    this.collectgetpurcon();
                    if(this.collecttableData == []){
                        this.dataAvahair = true;
                    }
                }
            })
          })
          .catch(() => {
            this.$message({
              type: "info",
              message: "已取消删除",
            });
          });
      }
    },
    // 编辑收货地
    Editcollect(row) {
      this.api.address.one({id:row.id})
      .then(res=>{
          if(res.data.code == 200){
          let purdata = res.data.data;
          purdata.id = row.id;
          sessionStorage.setItem('queryAll',JSON.stringify(purdata));
          this.$router.push({name:"Editcollect"})
          }
      })
    },

    /*
    装货地点方法
    */
    pretendgetpurcon(){
      this.loadingpretend = true;
      // 传参
      var listByPageData = {
        type:3,
        searchStr:this.pretendpurSearch,
        page:{
          current:this.pretendcurrent,
          size:this.pretendsize
        }
      }
      // 渲染表格
      this.api.address.all(listByPageData)
      .then(res=>{
        this.loadingpretend = false;
        if(res.data.code == 200){
          if(res.data.data.records != ""){
            this.pretendtableData = res.data.data.records;
            this.pretendtotal = res.data.data.total;
          }else{
            this.dataAvapretend = true;
          }

        }
      })
    },
    //装货地点搜索功能
    pretendsearch() {
      this.pretendtableData = [];
      this.pretendcurrent = 1;
      this.pretendgetpurcon();
    },
    pretendChange(val) {
      this.pretendAddressval = val;
      this.pretenddelmul = [];
      for(let i=0;i<this.pretendAddressval.length;i++){
          this.pretenddelmul.push(this.pretendAddressval[i].id) 
      }
    },
    //装货地点 删除功能
    deletpretend() {
      if (this.pretendAddressval.length == []) {
        this.$message({
          type: "warning",
          message: "请选择要删除的数据!",
        });
      } else {
        this.$confirm("此操作将永久删除该文件, 是否继续?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        })
          .then(() => {
            this.api.address.del(this.pretenddelmul)
            .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                    type: 'success',
                    message: '删除成功!'
                    });
                    this.pretendtableData = [];
                    this.pretendcurrent = 1;
                    this.pretendgetpurcon();
                    if(this.pretendtableData == []){
                        this.dataAvapretend = true;
                    }
                }
            })
          })
          .catch(() => {
            this.$message({
              type: "info",
              message: "已取消删除",
            });
          });
      }
    },
    // 编辑装货地点
    Editpretend(row) {
      this.api.address.one({id:row.id})
      .then(res=>{
          if(res.data.code == 200){
          let purdata = res.data.data;
          purdata.id = row.id;
          sessionStorage.setItem('queryAll',JSON.stringify(purdata));
          this.$router.push({name:"Editpretend"})
          }
      })
    },

    /*
    卸货地点方法
    */
    dischargegetpurcon(){
      this.loadingpretendcharge = true;
      // 传参
      var listByPageData = {
        type:4,
        searchStr:this.dischargepurSearch,
        page:{
          current:this.dischargecurrent,
          size:this.dischargesize
        }
      }
      // 渲染表格
      this.api.address.all(listByPageData)
      .then(res=>{
        this.loadingpretendcharge = false;
        if(res.data.code == 200){
          if(res.data.data.records != ""){
            this.dischargetableData = res.data.data.records;
            this.dischargetotal = res.data.data.total;
          }else{
            this.dataAvacharge = true;
          }

        }
      })
    },
    //卸货地点搜索功能
    dischargesearch() {
      this.dischargetableData = [];
      this.dischargecurrent = 1;
      this.dischargegetpurcon();
    },
    dischargeChange(val) {
      this.dischargeAddressval = val;
      this.dischargedelmul = [];
      for(let i=0;i<this.dischargeAddressval.length;i++){
          this.dischargedelmul.push(this.dischargeAddressval[i].id) 
      }
    },
    //卸货地点 删除功能
    deletdischarge() {
      if (this.dischargeAddressval.length == []) {
        this.$message({
          type: "warning",
          message: "请选择要删除的数据!",
        });
      } else {
        this.$confirm("此操作将永久删除该文件, 是否继续?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        })
          .then(() => {
            this.api.address.del(this.dischargedelmul)
            .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                    type: 'success',
                    message: '删除成功!'
                    });
                    this.dischargetableData = [];
                    this.dischargecurrent = 1;
                    this.dischargegetpurcon();
                    if(this.dischargetableData == []){
                        this.dataAvacharge = true;
                    }
                }
            })
          })
          .catch(() => {
            this.$message({
              type: "info",
              message: "已取消删除",
            });
          });
      }
    },
    // 编辑卸货地点
    Editdischarge(row) {
      this.api.address.one({id:row.id})
      .then(res=>{
          if(res.data.code == 200){
          let purdata = res.data.data;
          purdata.id = row.id;
          sessionStorage.setItem('queryAll',JSON.stringify(purdata));
          this.$router.push({name:"Editdischarge"})
          }
      })
    },
  },
};
</script>
<style lang="less" scoped>
// 选项卡开始
.menu-list {
  width: 100%;
  height: 48px;
  background-color: #fff;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  border: 1px solid rgba(102, 102, 102, 0.15);
}
.totalTab {
  width: 7.6%;
  float: left;
  font-size: 16px;
  line-height: 26px;
  font-weight: 650;
  color: #333;
  text-align: center;
  margin: 9px 0 0 35px;
}
.totalTab.active {
  color: #ea222e;
}
.stripnone {
  height: 3px;
  background-color: #ea222e;
  margin: 11px auto 0;
  display: none;
}
.strip {
  display: block;
}
.pur-nav {
  width: calc(100% - 32px);
  height: calc(100vh - 144px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 8px 8px 8px 8px;
  .pur-table{
    width:100%;
    height: 94%;
  } 
  .par-pages {
    width: 100%;
    height: 6%;
    padding: 6px 0 0 12px;
    box-sizing:border-box;
    .el-pagination {
      float: right;
    }
  }
}
//选项卡结束
// 暂无数据样式开始
.dataAvailable {
  width: 100%;
  height: 245px;
  position: relative;
  margin-top:185px;
  margin-bottom:279px;
  .dataAva {
    position: absolute;
    top: 75%;
    left: 50%;
    transform: translate(-50%, -50%);
    text-align: center;
    img {
      width: 371px;
      height: 200px;
      margin: auto;
    }
    p {
      font-size: 16px;
      line-height: 24px;
      color: #333;
    }
    .el-button {
      margin-bottom: 100px;
    }
  }
}
// 暂无数据样式结束
.parNavstrip {
  position: relative;
  .par-right {
    position: absolute;
    top: -50px;
    right: 24px;
    .el-input {
      width: 200px;
      float: right;
    }
    .el-button {
      margin-left: 16px;
      float: right;
    }
  }
}
</style>