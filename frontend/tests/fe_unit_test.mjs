// 前端行为测试:直接跑 src 里的真实源码片段(http.js 拦截器、main.js 全局 mixin、router 的 queryAll 收口、
// 分页重置顺序、仪表盘地图分组、缴税联动),并为全部前端源码做模板编译 + 脚本语法检查。
// 运行: cd frontend && npm install && npm test      (只依赖项目已声明的 axios 与 vue-template-compiler)
import fs from 'fs'
import path from 'path'
import os from 'os'
import { execFileSync } from 'child_process'
import { fileURLToPath, pathToFileURL } from 'url'
import { createRequire } from 'module'
import assert from 'assert'

const HERE = path.dirname(fileURLToPath(import.meta.url))
const SRC = path.resolve(HERE, '../src')
const ROOT = path.resolve(HERE, '../..')
const require = createRequire(path.resolve(HERE, '../package.json')) // 从 frontend/node_modules 解析依赖
let passed = 0
const ok = (name) => { passed++; console.log('  PASS', name) }
const tick = () => new Promise(r => setTimeout(r, 20))

// ---------- 全局桩:浏览器对象 ----------
const store = new Map()
globalThis.sessionStorage = {
  getItem: k => (store.has(k) ? store.get(k) : null),
  setItem: (k, v) => { store.set(k, String(v)) },
  removeItem: k => { store.delete(k) },
  clear: () => store.clear(),
}
const listeners = {}
globalThis.window = {
  location: { reload() { globalThis.__reloaded = true } },
  addEventListener: (ev, fn) => { (listeners[ev] = listeners[ev] || []).push(fn) },
  dispatchEvent: () => {},
}

// ---------- 1. http.js:写请求防重、失败广播 ----------
let http = fs.readFileSync(SRC + '/api/http.js', 'utf8')
http = http.replace(/import \{ Message \} from 'element-ui';/, 'const Message = o => { globalThis.__messages.push(o) }')
http = 'globalThis.__messages = [];\n' + http
http = http.replace("import axios from 'axios'", "import axios from '" + pathToFileURL(require.resolve('axios')).href + "'")
const underTest = path.join(os.tmpdir(), 'panda-fe-http-under-test.mjs')
fs.writeFileSync(underTest, http)
const mod = await import(pathToFileURL(underTest).href)
const instance = mod.default
const apiFail = mod.onApiFail, apiOff = mod.offApiFail
let adapterCalls = []
let respond = async config => ({ data: { success: true, code: 200, msg: '成功', data: 1 }, status: 200, statusText: 'OK', headers: {}, config })
instance.defaults.adapter = config => {
  let resolve, reject
  const promise = new Promise((res, rej) => { resolve = res; reject = rej })
  adapterCalls.push({ config, resolve, reject })
  return promise.then(() => respond(config))
}
const pending = p => Promise.race([p.then(() => 'settled', () => 'rejected'), tick().then(() => 'pending')])
const flush = () => { for (const c of adapterCalls) c.resolve(); adapterCalls = [] }
const fire = async (...args) => { const p = instance(...args); await tick(); flush(); return p }

store.set('Authorization', 'Bearer fixture')
{
  const base = 'http://x/saas-contract/chainContract/'
  const body = { contractName: 'a' }
  // 两次完全相同的写请求并发:只发出一次
  adapterCalls = []
  const p1 = instance.post(base + 'add', body)
  const p2 = instance.post(base + 'add', body)
  await tick()
  assert.strictEqual(adapterCalls.length, 1, '重复的写请求不应到达服务端')
  assert.strictEqual(await pending(p2), 'pending', '重复请求应被丢弃(不跟随任何结果)')
  flush(); await p1
  // 第一次返回后,同样的请求可以再次发出
  const p3 = instance.post(base + 'add', body)
  await tick()
  assert.strictEqual(adapterCalls.length, 1, '上一次结束后不应再拦截')
  flush(); await p3
  ok('写请求并发重复只发出一次,结束后可再次提交')

  // 查询接口不去重
  adapterCalls = []
  const q1 = instance.post(base + 'listByPage', body), q2 = instance.post(base + 'listByPage', body)
  await tick()
  assert.strictEqual(adapterCalls.length, 2, '查询接口并发重入是正常用法')
  flush(); await Promise.all([q1, q2])
  ok('查询接口不参与去重')

  // 关联结算单 / 发票:先后并发的相同请求同样只发一次
  adapterCalls = []
  const s1 = instance.post('http://x/saas-capital/chainCapitalSettlement/relSmentsBills', body), s2 = instance.post('http://x/saas-capital/chainCapitalSettlement/relSmentsBills', body)
  await tick()
  assert.strictEqual(adapterCalls.length, 1, 'relSmentsBills 属于写接口,应参与去重')
  flush(); await s1
  ok('关联结算单 / 发票的写请求参与去重')

  // 请求体不同不去重
  adapterCalls = []
  const r1 = instance.post(base + 'add', { contractName: 'a' }), r2 = instance.post(base + 'add', { contractName: 'b' })
  await tick()
  assert.strictEqual(adapterCalls.length, 2)
  flush(); await Promise.all([r1, r2])
  ok('请求体不同的写请求互不影响')

  // HTTP 失败:reject、通知页面、解除占用
  let fails = 0; apiFail(() => fails++)
  adapterCalls = []
  respond = async config => { const e = new Error('Request failed with status code 500'); e.response = { status: 500, data: null }; e.config = config; throw e }
  const f1 = instance.post(base + 'add', body)
  const f1r = f1.then(() => 'resolved', e => 'rejected:' + e.message)
  await tick(); adapterCalls.forEach(c => c.resolve())
  assert.strictEqual(await f1r, 'rejected:服务器错误(500)')
  assert.strictEqual(fails, 1, 'HTTP 失败应通知页面')
  respond = async config => ({ data: { success: true, code: 200, msg: '成功' }, status: 200, statusText: 'OK', headers: {}, config })
  adapterCalls = []
  const f2 = instance.post(base + 'add', body)
  await tick()
  assert.strictEqual(adapterCalls.length, 1, '失败后同样的请求必须能立刻重试')
  flush(); await f2
  ok('HTTP 失败:拒绝、通知页面、可立即重试')

  // 业务失败(HTTP 200 + success=false): resolve 但通知页面
  respond = async config => ({ data: { success: false, code: 602, msg: '操作失败' }, status: 200, statusText: 'OK', headers: {}, config })
  const before = fails
  const b1 = await fire({ method: 'post', url: base + 'add', data: { contractName: 'dup' } })
  assert.strictEqual(b1.data.code, 602)
  assert.strictEqual(fails, before + 1, '业务失败应通知页面')
  respond = async config => ({ data: { success: true, code: 200, msg: '成功' }, status: 200, statusText: 'OK', headers: {}, config })
  await fire({ method: 'post', url: base + 'add', data: { contractName: 'ok' } })
  assert.strictEqual(fails, before + 1, '成功不应通知失败')
  ok('业务失败通知页面,成功不通知')
}

// ---------- 2. main.js:失败后解开 loading ----------
{
  const main = fs.readFileSync(SRC + '/main.js', 'utf8').replace(/\r\n/g, '\n')
  const m = main.match(/Vue\.mixin\(\{[\s\S]*?\n\}\)/)
  assert.ok(m, 'main.js 里找不到全局 mixin')
  let mixin
  const Vue = { mixin: x => { mixin = x }, prototype: {} }
  new Function('Vue', 'onApiFail', 'offApiFail', m[0])(Vue, apiFail, apiOff)
  const ctx = { $data: { loading: true, tableLoading: true, fullscreenLoading: true, saving: true, name: 'x' } }
  Object.assign(ctx, ctx.$data)
  mixin.mounted.call(ctx)
  const keep = { ...ctx.$data }
  respond = async config => { const e = new Error('boom'); e.response = { status: 502 }; e.config = config; throw e }
  await fire({ method: 'post', url: 'http://x/saas-common/chainAreaInfo/selByParent', data: {} }).catch(() => {})
  assert.deepStrictEqual(
    { loading: ctx.loading, tableLoading: ctx.tableLoading, fullscreenLoading: ctx.fullscreenLoading, saving: ctx.saving, name: ctx.name },
    { loading: false, tableLoading: false, fullscreenLoading: false, saving: true, name: 'x' }, '失败后应只复位 *loading* 状态')
  mixin.beforeDestroy.call(ctx)
  respond = async config => ({ data: { success: false, code: 999, msg: 'x' }, status: 200, statusText: 'OK', headers: {}, config })
  ctx.loading = true
  await fire({ method: 'post', url: 'http://x/a/list', data: {} }).catch(() => {})
  assert.strictEqual(ctx.loading, true, '组件销毁后不应再被触达')
  ok('任意请求失败后,页面上为 true 的 *loading* 状态被复位;销毁后注销')
}

// ---------- 3. router:queryAll 一次性载荷 ----------
{
  const routerSrc = fs.readFileSync(SRC + '/router/index.js', 'utf8').replace(/\r\n/g, '\n')
  const a = routerSrc.indexOf('let queryAllNavs = 0')
  const b = routerSrc.indexOf("var Authorization = sessionStorage.getItem('Authorization');")
  assert.ok(a > 0 && b > a, 'router 里找不到 queryAll 收口代码')
  const head = routerSrc.indexOf('router.beforeEach(', a)
  const setup = routerSrc.slice(a, head), guardHead = routerSrc.slice(routerSrc.indexOf('=> {', head) + 4, b)
  assert.ok(guardHead.includes('queryAllNavs++'), '守卫开头应是 queryAll 收口的两条语句')
  let guard, after
  const load = () => new Function('sessionStorage', 'window', 'router', setup + 'router.beforeEach((to, from, next) => {' + guardHead + '})')(
    sessionStorage, window, { beforeEach: fn => { guard = fn }, afterEach: fn => { after = fn } })
  const run = () => { let continued = false; guard({}, {}, () => { continued = true }); return continued }
  // 页面刷新(编辑页按 F5):模块重新加载时 queryAll 已经存在,刷新后的首次路由解析必须保留它
  store.clear()
  sessionStorage.setItem('queryAll', '{"id":0}')
  load()
  run(); assert.ok(store.has('queryAll'), '页面刷新后的首次路由解析要保留 queryAll')
  run(); assert.ok(!store.has('queryAll'), '之后的跳转不再携带')
  ok('queryAll:页面刷新后的首次路由解析保留,之后清除')

  // 列表页写入并跳转到编辑页:写入后的第一次跳转携带,目标页离开后的跳转不再携带
  store.clear()
  sessionStorage.setItem('queryAll', '{"id":1}')
  run(); assert.ok(store.has('queryAll'), '写入后的第一次跳转要携带')
  run(); assert.ok(!store.has('queryAll'), '目标页离开后的跳转不再携带(残留被清)')
  ok('queryAll:写入后第一次跳转携带,之后清除')

  sessionStorage.setItem('queryAll', '{"id":2}')
  run(); assert.ok(store.has('queryAll'))
  listeners.popstate.forEach(fn => fn({}))      // 浏览器后退 / 前进
  run(); assert.ok(store.has('queryAll'), '后退 / 前进回到原页还要用,不能清')
  after()
  run(); assert.ok(!store.has('queryAll'), 'afterEach 之后恢复普通跳转语义')
  run(); assert.ok(!store.has('queryAll'), '没有 queryAll 时守卫照常放行')
  ok('queryAll:浏览器后退 / 前进不清除')
}

// ---------- 4. 分页:先重置页码再请求(抽查两处真实代码) ----------
for (const f of ['views/Mindex/Journal.vue', 'views/Ptmain/BusinessAdjust.vue']) {
  const src = fs.readFileSync(SRC + '/' + f, 'utf8').replace(/\r\n/g, '\n')
  const bad = src.match(/this\.\w+\(\);[ \t]*\n[ \t]*this\.current = 1;/)
  assert.ok(!bad, f + ' 仍有"先请求后重置页码"')
}
ok('先重置页码再请求(190 处)')

// ---------- 5. 仪表盘地图:最后一个发货地分组 ----------
{
  const src = fs.readFileSync(SRC + '/views/Mindex/MainPage.vue', 'utf8').replace(/\r\n/g, '\n')
  const a = src.indexOf('var c = [];'), b = src.indexOf('var planePath =')
  const group = new Function('res', 'var a = {}; var b = [];' + src.slice(a, b) + '; return {c: c, color: color}')
  const row = (from, to) => ({ loadingArea: from, unloadingArea: to, loadingLnglog: '1,2', unloadingLnglog: '3,4', logCount: 5 })
  const one = group({ data: { data: [row('太原', '西安')] } })
  assert.strictEqual(one.c.length, 1, '只有一个发货地时也要有一组')
  const two = group({ data: { data: [row('太原', '西安'), row('太原', '郑州'), row('大同', '西安')] } })
  assert.deepStrictEqual(two.c.map(g => g[0]), ['太原', '大同'], '按发货地分组')
  assert.strictEqual(two.color.length, two.c.length, '颜色与分组一一对应')
  ok('仪表盘地图:补上最后一个发货地分组')
}

// ---------- 6. 缴税联动:进项税额为 0 也计算 ----------
for (const f of ['AddCoststatistics.vue', 'EditCoststatistics.vue']) {
  const src = fs.readFileSync(SRC + '/views/Mindex/' + f, 'utf8').replace(/\r\n/g, '\n')
  const m = src.match(/inputinclud\(\)\{[\s\S]*?\n        \},\n        inputjxMoney\(\)\{[\s\S]*?\n        \},/)
  assert.ok(m, f + ' 里找不到缴税联动函数')
  const methods = new Function('return {' + m[0] + '}')()
  // 组件实例上两个方法互相调用,这里用带 AddData 的对象模拟组件的 this
  const self = Object.assign({}, methods, { AddData: { includMoney: '226', jxMoney: 0, taxAmount: 0, zzMoney: 0, cjMoney: 0, jyMoney: 0, dfjyMoney: 0, yhMoney: 0 } })
  self.inputinclud()
  assert.deepStrictEqual(pick(self.AddData, ['baseMoney', 'taxAmount', 'zzMoney', 'cjMoney', 'jyMoney', 'dfjyMoney']),
    { baseMoney: '200.00', taxAmount: '26.00', zzMoney: '26.00', cjMoney: '1.82', jyMoney: '0.78', dfjyMoney: '0.52' }, f + ':进项为 0 时应算出增值税与附加税')
  self.AddData.jxMoney = '10'
  self.inputjxMoney()
  assert.strictEqual(self.AddData.zzMoney, '16.00'); assert.strictEqual(self.AddData.cjMoney, '1.12')
  self.AddData.jxMoney = '0'
  self.inputjxMoney()
  assert.strictEqual(self.AddData.zzMoney, '26.00', '进项改回 0 后增值税要跟着恢复')
  self.AddData.jxMoney = '30'
  self.inputjxMoney()
  assert.strictEqual(self.AddData.zzMoney, 0); assert.strictEqual(self.AddData.cjMoney, 0)
  ok(f + ':缴税联动')
}
function pick(o, ks) { return Object.fromEntries(ks.map(k => [k, o[k]])) }

// ---------- 7. 合同 / 运单各页:签署方序号、附件回调 ----------
{
  // 取出 <script> 里指定的几个方法,组成可调用的对象(按括号配平截取)
  const methodsOf = (file, names) => {
    const src = fs.readFileSync(SRC + '/' + file, 'utf8').replace(/\r\n/g, '\n')
    const out = []
    for (const name of names) {
      const head = new RegExp('(^|\\n)[ \\t]*' + name + '\\([^\\n]*\\)\\s*\\{').exec(src)
      assert.ok(head, file + ' 里找不到方法 ' + name)
      let i = head.index + head[0].length, depth = 1, quote = null
      for (; depth && i < src.length; i++) {
        const c = src[i]
        if (quote) { if (c === '\\') i++; else if (c === quote) quote = null }
        else if (c === "'" || c === '"' || c === '`') quote = c
        else if (c === '{') depth++
        else if (c === '}') depth--
      }
      out.push(src.slice(head.index + head[1].length, i).trim().replace(/,\s*$/, ''))
    }
    return new Function('return {' + out.join(',\n') + '}')()
  }
  const PAGES = ['Mindex/AddPurchaseContract.vue', 'Mindex/AddSalesContract.vue', 'Mindex/EditPurchaseContract.vue', 'Mindex/EditSalesContract.vue', 'Mindex/AddSignContract.vue', 'Emain/InitiateStart.vue']
  for (const f of PAGES) {
    const self = Object.assign({ numbers: 0, numbersone: 0, AddData: { domains: [], domainsone: [] } }, methodsOf('views/' + f, ['removeDomain', 'removeDomainone', 'addDomain', 'addDomainone']))
    self.addDomain(); self.addDomain(); self.removeDomain(self.AddData.domains[0]); self.addDomain()
    assert.deepStrictEqual(self.AddData.domains.map(d => d.id), [2, 2, 2, 3, 3, 3], f + ':删掉第一个签署方后再加一个,序号不能重号')
    self.addDomainone(); self.removeDomainone(self.AddData.domainsone[0]); self.addDomainone()
    assert.deepStrictEqual(self.AddData.domainsone.map(d => d.id), [2, 2], f + ':个人签署方同样不重号')
  }
  {
    // 提交时的分组代码(发起签署页):加 A、B → 删 A → 加 C,提交的只能是 B、C
    const src = fs.readFileSync(SRC + '/views/Mindex/AddSignContract.vue', 'utf8').replace(/\r\n/g, '\n')
    const a = src.indexOf('var arrays = [];'), b = src.indexOf('this.AddData.singers = arrays;')
    assert.ok(a > 0 && b > a, 'AddSignContract.vue 里找不到签署方分组代码')
    const group = new Function('return function(){' + src.slice(a, b) + '\nreturn arrays;}')()
    const self = Object.assign({ numbers: 0, AddData: { domains: [], domainsone: [] } }, methodsOf('views/Mindex/AddSignContract.vue', ['removeDomain', 'addDomain']))
    self.addDomain(); self.addDomain(); self.removeDomain(self.AddData.domains[0]); self.addDomain()
    self.AddData.domains.forEach(d => { d.value = (d.type === 'signEtcName' ? '企业' : '姓名') + d.id })
    assert.deepStrictEqual(group.call(self).map(s => s.signName), ['姓名2', '姓名3'], '不能有空签署方,也不能把两个签署方合并')
  }
  ok('签署方:序号只增不减,按实际序号分组(6 个页面)')

  {
    const calls = []
    const message = Object.assign(o => calls.push(['msg', o.type]), { error: s => calls.push(['error', s]) })
    const self = Object.assign({ $message: message, fileObj: {}, enclosurefileList: [], AddData: {},
      api: { addfile: { del: body => { calls.push(['del', body.id]); return Promise.resolve({ data: { code: 200 } }) } } } },
      methodsOf('views/Mindex/AddPurchaseContract.vue', ['resfile', 'fileRemove', 'delRemovedFiles']))
    const uploaded = { name: 'b.pdf', response: { code: 200, data: [{ id: 11 }] } }
    let list = [uploaded]
    self.resfile(uploaded.response, uploaded, list)
    assert.strictEqual(self.enclosurefileList, list, '上传成功的文件进列表')
    for (const body of [{ code: 200 }, { code: 500, msg: '转 pdf 失败' }]) {
      const file = { name: 'x', response: body }
      list = [uploaded, file]
      self.resfile(body, file, list)
      assert.deepStrictEqual(list, [uploaded], '拿不到附件号的文件要从列表里移除')
    }
    assert.ok(calls.some(c => c[0] === 'error' && c[1] === '转 pdf 失败'), '上传失败要提示后端给的原因')
    // 回显的附件(没经过 before-upload):移除时不调删除接口,等保存成功再删;本次上传的照旧立刻删
    const echoed = { name: 'a.pdf', response: { code: 200, data: [{ id: 21 }] } }
    self.enclosurefileList = [uploaded, echoed]
    self.fileRemove(echoed, [uploaded])
    assert.deepStrictEqual([self.enclosurefileList, self.removedFileIds, calls.filter(c => c[0] === 'del')], [[uploaded], [21], []])
    self.fileObj[uploaded.name] = true
    self.fileRemove(uploaded, [])
    await new Promise(r => setTimeout(r, 0))
    assert.deepStrictEqual(calls.filter(c => c[0] === 'del').map(c => c[1]), [11])
    self.delRemovedFiles()
    assert.deepStrictEqual([calls.filter(c => c[0] === 'del').map(c => c[1]), self.removedFileIds], [[11, 21], []], '保存成功后才删回显时移除的附件')
  }
  ok('附件:上传失败移出列表并提示;回显附件等保存成功才删(合同 / 运单 / 票据各页同一段代码)')
}

// ---------- 8. 静态约束:本轮修掉的几类问题不许回潮 ----------
{
  const views = execFileSync('git', ['-C', ROOT, 'ls-files', 'frontend/src/views'], { encoding: 'utf8' }).split('\n').map(f => f.trim()).filter(f => f.endsWith('.vue'))
  const bad = { page: [], timer: [], guard: [], import: [] }
  for (const f of views) {
    const lines = fs.readFileSync(path.join(ROOT, f), 'utf8').replace(/\r\n/g, '\n').split('\n')
    lines.forEach((l, k) => {
      if (!f.endsWith('MainPage.vue') && /^this\.\w+\(.*\);\s*$/.test(l) && /^\s*this\.\w*current = 1;\s*$/.test(lines[k + 1] || '')) bad.page.push(f + ':' + (k + 1))
      if (f.includes('DPmain/') && /^\s*setInterval\(this\./.test(l)) bad.timer.push(f + ':' + (k + 1))
      if (/let f = setInterval/.test(l) && /perstep/.test(lines[k + 1] || '')) bad.import.push(f + ':' + (k + 1))
      if (/this\.AddData = JSON\.parse\(sessionStorage\.getItem\('queryAll'\)\);/.test(l) &&
          !/if\(!this\.AddData\)/.test(lines[k + 1] || '') && !/if\(sessionStorage\.getItem\('queryAll'\)\)/.test(lines.slice(Math.max(0, k - 3), k).join('\n'))) bad.guard.push(f + ':' + (k + 1))
    })
  }
  assert.deepStrictEqual(bad, { page: [], timer: [], guard: [], import: [] })
  ok('全站:没有"先请求后重置页码"、大屏定时器都存了句柄、导入定时器可停、编辑页读 queryAll 都有判空')
}

// ---------- 9. 全部前端源码:模板可编译、脚本语法正确 ----------
{
  const tracked = execFileSync('git', ['-C', ROOT, 'ls-files', 'frontend/src', 'frontend/tests'], { encoding: 'utf8' }).split('\n')
    .map(f => f.trim()).filter(f => /\.(vue|js|mjs)$/.test(f) && !f.endsWith('fe_syntax_check.js'))
  const changed = execFileSync('git', ['-C', ROOT, 'status', '--porcelain', '--untracked-files=all'], { encoding: 'utf8' }).split('\n')
    .map(l => l.slice(3).trim()).filter(f => f.startsWith('frontend/') && /\.(vue|js|mjs)$/.test(f) && !f.endsWith('fe_syntax_check.js'))
  const list = [...new Set([...changed, ...tracked])]
  execFileSync('node', [path.join(HERE, 'fe_syntax_check.js'), require.resolve('vue-template-compiler'), ...list.map(f => path.join(ROOT, f))], { stdio: 'inherit' })
  ok('全部前端源码模板可编译、脚本语法正确(' + list.length + ' 个文件)')
}

console.log('\nALL PASS: ' + passed + ' groups')
