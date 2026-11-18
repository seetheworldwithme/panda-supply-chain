# 六期 Upload 与外围验收

2026-11-17，本地 Java 8765 / Go 8799，chain 库，仅临时 fixture，未调用腾讯 OCR 或部署远程。

- `tools/migration/cmp-common-external.py`：18/18 安全只读场景响应 HTTP 状态及排序后的全量 JSON/文本一致，覆盖 Upload 缺文件、缺参、过期分享、OCR 缺参/缺文件、APK 缺文件、固定路径下载缺文件。
- `tools/migration/verify-common-external.py`：11 项正向断言通过（Java/Go 各4项图片、视频、下载二进制、关联文件；Go 内存分享令牌、uploadFiles 文件及DB、bigUpload 合并文件及DB）。测试行99977781及新上传行、生成文件均清理。
- `external_test.go` 5项单测覆盖固定下载缺失文本、视频Range语义、分享令牌、Java下载缓冲尾部bug、网关鉴权。均通过。
- 脚本凭据读取 `CMP_PASSWORD` 环境变量，无凭据写入仓库。

## Java 实证怪癖

1. `FileUtil.download` 每次写1024字节而非实际读取长度；最后一块补零或残留上一块尾部。Go保留该行为。
2. `getVideo` Range忽略请求终点，返回从起点至文件末；`play` 仅返回最多1MiB但Content-Range仍标末字节。Go保留。
3. `preShareFIle` 取Redis字面键`credit`而非请求credit；`getShareToken`则写UUID键。因此新发分享令牌仍提示601过期，Go保留。
4. `ActivitiController.java` 整文件注释，无生效工作流路由，不能按旧提示增加虚假实现。
5. Java网关对上述外围全部要求token，包含名为“未登录上传”的upload/upload，Go补authGuard。

## 已验证的兼容范围与差异

- 本地 Java 有效id获取分享token返回HTTP500（Redis连接不可用）。Go按任务允许使用进程内带过期表，返回UUID，重启不保留。TTL当前24小时，后续环境配置需与commenConfig.overHours对齐。
- 本地 Java 有效txt的uploadFiles返回200空数组（配置上传目录不可用），Go使用PANDA_UPLOAD_DIR（默认data/uploads）可写路径，已验证文件与DB。该环境差异不纳入18/18一致性统计。
- uploadFiles 同步生成压缩图、源文件压缩关联和 Excel 预览图。Go 用 LibreOffice 和 pdftoppm 渲染 Excel，输出字节和 Java Spire 不同。Go 的文档上传保持本机 Java Doc2PdfUtil 的行为，保存源文件记录后终止转换并返回空列表。运行远程适配前需要确认实际部署的文档转换组件。
- `tools/migration/verify-common-media.py` 验证本地生成压缩关联、Excel 预览图和文档失败后的存储行为。
- OCR有效type 1/2/3需要腾讯服务。Go不复制Java硬编码凭据，支持PANDA_OCR_ADAPTER_URL接收{imageBase64,type}、返回Java JsonResult；未配置返回701/ExternalServiceUnavailable。仅缺参/缺文件与Java一致，真实识别尚未验收。非法type按Java空字符串成功。
- 普通下载/预览依赖现存磁盘附件，PANDA_APK_FILE、PANDA_TL_SFTP_DIR、PANDA_BIG_UPLOAD_DIR分别覆盖APK/SFTP/分片路径。测试仅用临时本地小文件，未覆盖历史13G附件或所有MIME。

2026-11-18 重新启动本地 Java 对照环境后，外围兼容样本重新验证为 18/18。
