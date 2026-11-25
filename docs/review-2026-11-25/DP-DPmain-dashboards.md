# DP：DPmain 大屏 8 页（定向检查）

> 日期：2026-11-26　范围：`frontend/src/views/DPmain/` 下 `Firstscreen`、`Secondscreen`、`Thirdscreen`、`Maps`、`NjMain`、`warehouse`、`Dpmain`、`Welcome` 共 8 个 `.vue`。

## 方法与限制

这一组做的是**定向检查，不是逐行通读**。原计划的逐行评审任务被会话中断，无法恢复，所以改由主会话按类扫描脚本块并读上下文：

- 定时器（`setInterval` / `setTimeout`）、`window` 监听、`echarts.init` 与 `beforeDestroy` 的配对：全部命中行逐条读了创建与销毁两段；
- 对接口返回值取 `.length` / 下标 / 做算术、日期计算、HTML 注入面：按模式扫并统计。

**没有看的部分**：各图表的 option 配置、统计接口字段的业务口径、全部 `<style>`。结论里的「已核对」只对上面列出的类别成立。

## 已修（本轮）

| 级别 | 位置 | 问题 | 修法 |
| --- | --- | --- | --- |
| 🔴 | `Firstscreen:182-204`、`Secondscreen:263-314`、`NjMain:219-234`、`Thirdscreen:188`、`Maps:51-52`、`warehouse:368` | `setInterval(this.formatDate, 1000)` 这类调用的返回值直接丢弃；`beforeDestroy` 里写的是 `clearInterval(this.formatDate)`——传的是函数，不是定时器号，**一个也清不掉**。大屏靠 `routername` 每 60 秒换一屏，每换一屏就多出 1–9 个永不停止的定时器（每秒刷时间、每 2 秒滚动、每 10–30 秒轮询统计接口），挂一夜的大屏会有上百个定时器同时在刷接口 | 返回值存进 `this.timers`，`beforeDestroy` 里统一 `clearInterval` |
| 🟡 | `Firstscreen`、`Secondscreen`、`Thirdscreen`、`Maps`、`NjMain` | `echarts.init` 的实例都是局部变量，离开页面从不 `dispose`；`window.addEventListener('resize', () => chart.resize())` 是匿名监听，无法移除 | 销毁时按容器上的 `_echarts_instance_` 标记逐个 `dispose`；监听改成具名保存，销毁时 `removeEventListener` |
| 🟡 | `Thirdscreen:1010`、`:1302` | 两个图表轮播定时器都存在 `this.timer` 上，前一个的句柄被覆盖，销毁时清不掉 | 第二个改存 `this.timerthree` |
| 🟡 | `warehouse:369`、`:575` | `interone` / `inter` 两个滚动定时器销毁时不清 | `beforeDestroy` 补清 |
| 🟡 | `Maps.vue:286-301` | 地图分组只在「切换发货地」时入组，最后一组（或只有一个发货地时）丢失；接口无数据返回 `null` 时 `.length` 抛错。与 `Mindex/MainPage.vue` 是同一段代码 | 循环后补最后一组；`Array.isArray` 判空 |
| 🟢 | `Firstscreen`、`NjMain`、`Secondscreen` | `(res.data.data.xxx / 100000000)`：接口无数据时 `data` 是 `null`，回调里抛 `TypeError` | `(res.data.data || {})` |

## 未修

- `Thirdscreen:484` 还有一处 `resize` 监听是匿名的（监听体里除了 `resize()` 还接着起定时器，形状不同，没有套用同一段修法）。图表销毁后窗口尺寸变化时仍会调用已销毁的实例，只在控制台报错。
- `warehouse.vue` 有 4 处把写死的签名拼进图片 / 下载地址，属全仓模式，统一见 `docs/review.md` §8.4。
- `Dpmain.vue` 的 `secede()` 退出登录不做整页刷新。大屏账号独立（`account_type=4`），换账号的场景很少，未动。
- 图表 option、统计口径没有评审。

## 已核对（仅限上面的类别）

- 8 个页面都没有 `v-html` / `innerHTML`。
- 当月最后一天用 `new Date(year, month, 0).getDate()` 计算，没有写死。
- `Dpmain.vue`、`Welcome.vue` 没有定时器和图表。

## 回归检查

`frontend/tests/fe_unit_test.mjs` 第 8 组有一条静态约束：`DPmain/` 下不允许再出现不保存返回值的 `setInterval(this.xxx, n)`。
