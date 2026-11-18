# 六期收官接续计划

**Goal:** 完成 warehouse 第三批及 common 剩余接口，保留 Java 实际行为，按批验收后推送 main。
**Architecture:** 保留现有 warehouse 元数据模板与业务钩子，common 剩余使用独立包，避免重注册。
**Tech Stack:** Go net/http、MySQL、Java Spring 对照、Python 对比、ego-browser。

1. 已核对 git pull --ff-only：main 包含四期 46a1dfe、五期 254544c、六期前两批 832ba6e/05ccafb；保留六个未提交 warehouse 文件与未跟踪 AGENTS.md。
2. 恢复 Go 8799，对 phase6-cases-b3.json 跑 Java 8765/Go 8799 对比；将失败响应作为回归用例，逐处读 Java 修复。
3. 独立盘点 warehouse 未注册 bespoke；补 Xmmc、VR、Wechat 等实际缺口。生成器与不含凭据用例迁入仓库，解决 /tmp 丢失不可复现问题。
4. 对照全部 common 控制器与既有路由差集，实现剩余服务；不得以未验证的恒999替代业务。
5. 每批 go test ./...、读JSON对比、写入999777xx→DB断言→清理后提交与 push。失败不得声明完成。
6. 重跑历史对比/写验证，浏览器双账号全菜单冒烟及截图；恢复临时前端与网关端口。
7. 更新三表、迁移记忆与已知差异；仅在证据覆盖后登记迁移完成。不部署远端。
