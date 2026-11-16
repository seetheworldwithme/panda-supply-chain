<template>
  <div style="width: 100%; height: 100vh">
    <div ref="bpmn" style="width: 100%; height: 100%;"></div>
    <div ref="panel" class="panel"></div>
    <button @click="showSave = true" class="saveBtn" style="margin-left: -50px;">保存</button>
    <button @click="loadData" class="saveBtn" style="margin-left: 50px;">重置</button>
    <el-dialog title="保存流程图" :visible.sync="showSave">
      <el-form :model="model">
        <el-form-item label="key">
          <el-input v-model="model.key" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="name">
          <el-input v-model="model.name" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="showSave = false">取 消</el-button>
        <el-button type="primary" @click="save">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
  
  // 这个只能预览
  // import BpmnViewer from 'bpmn-js/lib/Viewer'
  import Modeler from 'bpmn-js/lib/Modeler'
  // 面板
  import propertiesPanelModule from 'bpmn-js-properties-panel'
  import propertiesProviderModule from 'bpmn-js-properties-panel/lib/provider/bpmn'
  // 适应activiti
  import activiti from '@/activiti.json'
  // 汉化
  import customTranslate from '@/customTranslate/customTranslate.js'
  
  export default {
    data() {
      return {
        showSave: false, model: {id: '', key: '', name: '', xml: ''}, modeler: null
      }
    },
    mounted() {
      this.modeler = new Modeler({
        container: this.$refs.bpmn, propertiesPanel: {parent: this.$refs.panel},
        additionalModules: [
          propertiesPanelModule, propertiesProviderModule,
          {translate: ['value', customTranslate]}
        ],
        moddleExtensions: {
          activiti: activiti
        }
      })
      this.loadData()
    },
    methods: {
      loadData() {
        this.model.id = this.$route.query.id
        new Promise(resolve => {
          if (this.model.id) {
            this.$apis.act.act.getModel.req({id: this.model.id}).then(res => {
              if (res.data && res.data.metaInfo) {
                this.model = res.data
              }
              resolve(this.model.metaInfo)
            })
          } else {
            resolve()
          }
        }).then(xml => {
          xml ? this.modeler.importXML(xml) : this.modeler.createDiagram()
        })
      },
      save() {
        this.modeler.saveXML().then(xml => {
          this.model.xml = xml.xml;
          console.log(this.model)
          this.$apis.act.act.saveModel.req(this.model).then(res => {
            console.log(res)
            this.$message({type: 'success', message: '成功'})
            this.$router.push({path: '/bpmn/modelEdit', query: {id: res.data.id}})
            this.loadData()
          })
        })
        /*this.modeler.saveSVG().then(b => {
          console.log(b)
        })*/
      },
    }
  }
</script>

<style scoped lang="less">
  /* 基础样式 */
  @import '~bpmn-js/dist/assets/diagram-js.css';
  @import '~bpmn-js/dist/assets/bpmn-font/css/bpmn.css';
  @import '~bpmn-js/dist/assets/bpmn-font/css/bpmn-codes.css';
  @import '~bpmn-js/dist/assets/bpmn-font/css/bpmn-embedded.css';
  /* 右侧面板样式 */
  @import '~bpmn-js-properties-panel/dist/assets/bpmn-js-properties-panel.css';
  
  .saveBtn {
    position: absolute;
    top: 80%;
    left: 50%;
    padding: 12px 16px;
    border-radius: 10px;
    font-size: 18px;
    cursor: pointer;
  }
  
  .panel {
    position: absolute;
    box-sizing: border-box;
    width: 400px;
    height: 100%;
    padding: 10px 15px;
    top: 80px;
    right: 15px;
    overflow: auto;
  }
</style>