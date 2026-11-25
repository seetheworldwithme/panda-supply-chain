// 前端源码语法检查:模板能否编译(vue-template-compiler)+ 脚本块语法(node --check)。
// 用法: node tests/fe_syntax_check.js <vue-template-compiler 的绝对路径> <文件...>
// 由 tests/fe_unit_test.mjs 调用;也可手工运行。只读,临时文件写到系统临时目录。
const fs = require('fs')
const os = require('os')
const path = require('path')
const { execFileSync } = require('child_process')

const compiler = require(process.argv[2])
const files = process.argv.slice(3)
let bad = 0
for (const f of files) {
  const src = fs.readFileSync(f, 'utf8')
  if (f.endsWith('.vue')) {
    const sfc = compiler.parseComponent(src)
    const res = compiler.compile(sfc.template ? sfc.template.content : '')
    if (res.errors.length) { bad++; console.log('TEMPLATE ERROR', f, res.errors.slice(0, 3)) }
    if (!sfc.script) continue
    const out = path.join(os.tmpdir(), 'panda-fe-' + Buffer.from(f).toString('base64').replace(/[^A-Za-z0-9]/g, '') + '.mjs')
    fs.writeFileSync(out, sfc.script.content)
    try { execFileSync('node', ['--check', out], { stdio: 'pipe' }) } catch (e) { bad++; console.log('SCRIPT ERROR', f, String(e.stderr).slice(0, 400)) }
    fs.unlinkSync(out)
  } else {
    try { execFileSync('node', ['--check', f], { stdio: 'pipe' }) } catch (e) { bad++; console.log('SCRIPT ERROR', f, String(e.stderr).slice(0, 400)) }
  }
}
console.log('checked ' + files.length + ' files: template/script errors = ' + bad)
process.exit(bad ? 1 : 0)
