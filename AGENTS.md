# AGENTS.md — 熊猫云管理（panda-supply-chain，Java→Go 迁移）

> 供后续开发会话接续。更新于 2026-11-17。本文件只维护**云管理（供应链/老煤炭云）**相关内容；熊猫主工程（uni-app x 前端 + Go 后台）见 `/Volumes/coding/HBuilderxProject/panda/AGENTS.md`，服务器简称（「101」「152」「192」「后端机」）沿用其第 2 节定义：「后端机」= 101.42.30.102（`ssh root@101.42.30.102` 密码登录，密码徐先生持有，未存盘不入库）。

安全纪律（与主工程一致）：不将账号密码、token、私钥写入源码、文档或提交信息；后端机中间件一律绑 `127.0.0.1`，不得 `0.0.0.0` 监听；进程按 PID 操作，不用会自匹配的 `pkill -f`；老集群五台（.6 / .14 / .16.10 / .7 / .64.11）只读审计，不做写操作。

## 1. 仓库定位

- `origin = https://github.com/seetheworldwithme/panda-supply-chain.git`：老煤炭云 SaaS（云管理/供应链）的 Go 重写仓。`backend/` 为 Go 迁移源码，`docs/` 放迁移规格书（如 `BusProductController-Go迁移规格书.md`），`frontend/` 为老前端 Vue 源码。
- 迁移进度以 git log 为准（截至 2026-11-17：一/二/三期含合同 1fed41c、五期批2 票据 95 端点、六期批1 煤种域 153 端点、六期批2 会员积分红包域 184 端点，读对比与写验证全过）。
- 老煤炭云集群五台详表在主工程仓 `docs/服务器整理/服务器.md`（172.21.16.6 前端 / 49.232.228.168=.14 能力中心 / 172.21.16.10「数据库等」=研发底座+业务主库 / 82.157.64.231=.7 云物流 / 172.21.64.11 备份机；南京 119.45.238.116、10.206.0.7 已停摆）。老库核对结论见第 4 节。

## 2. 后端机上的云管理栈（101.42.30.102）

**Java 供应链栈已部署（2026-11-16，自老集群 .14/.16.10 迁移）**：9 个服务（common/auth/bill/capital/contract/logistics/warehouse + 双网关 8765/5765）跑在 `/home/ubuntu/saas`（复刻老机路径，`java -Dfile.encoding=utf-8 -Xmx512m -jar`，nohup，日志 `/root/panda/supply/logs/`）；中间件全在 `/root/panda/supply/`（软链至数据盘）：Nacos 1.4.1 standalone 8848（derby 配置数据原样迁移，9 个 yaml 已全量改写 172.21.16.10→127.0.0.1，原配置备份 `/root/panda/cfgs-orig/`）、RocketMQ 5.1.4（namesrv 9876 + broker 10911，brokerIP 已改 127.0.0.1，老 store 数据在 `/root/panda/supply/store`）、Redis 6.0.10 自编译（127.0.0.1:6379，requirepass 同老系统，AOF；**已停用机器自带的宝塔 Redis**，勿再启动）。

**老地址 172.21.16.10 兜底 = lo 别名（2026-11-16 定稿）**：common/warehouse 的 Redisson 把 `172.21.16.10:6379` 硬编码在类文件里（jarmod 只改 yml 改不到），现用 `ip addr add 172.21.16.10/32 dev lo` + Redis `bind 127.0.0.1 172.21.16.10` 兜底（早前的 iptables DNAT 方案被 Docker 重写 nat 表抹掉，已弃用）；**lo 别名重启即丢**，恢复脚本会重加，勿删别名与 Redis 双 bind；firewalld 已放行 8765/5765（来源仅 152.136.54.164，rich rule）。

**uploadfiles（13G，11,560 个文件）已迁完并核验（2026-11-17）**：双侧文件数 11560 相等、总字节数 13,379,309,501 分毫不差、逐文件大小清单比对 0 差异，拉取阶段 6 个 .done 标记齐全；老服务器零散 5 件（apiclient_key.pem/checkDiskSize.sh/wa.txt/zsc.html/arthas-boot.jar）经本机 Mac 中转迁完（字节数逐一核对），归档 `/mnt/datadisk0/saas/misc/`（与 `/home/ubuntu/` 原件字节级一致，原件保留防老路径引用）。

**已知路径问题（2026-11-17 起）：后端机→.14 方向 TCP 能握手但数据流卡死**（13G 传完后出现，scp/ssh 会话均挂起；Mac→.14 会话正常）——再从 .14 取东西须走本机中转，或先诊断该方向路径（疑似运营商/云间路径 MTU 或状态问题）。

**内存纪律（2026-11-16 抖死事故后定）**：8G 机器跑满 9 JVM + 中间件贴顶，已建 4G swap（fstab 持久，swappiness=10）；RocketMQ 必须压堆启动（namesrv 256m/broker 768m/proxy 256m，经 `JAVA_OPT_EXT`），默认大堆会爆；**proxy 必须等 broker 路由注册完（约 30s）再起**，否则建系统 topic `DefaultHeartBeatSyncerTopic` 失败秒退，proxy 不在线时拉 common/warehouse 也秒死（gRPC 依赖）。

**数据盘布局（`/mnt/datadisk0/saas/`，2026-11-17 定稿）**：`jars/`=9 个 Java 包（:9001-9102，双网关 8765/5765）、`supply/`=组件（Nacos 8848+RocketMQ 9876/10911/28081+Redis 6379）、`mysql-data/`=云管理 8 库（容器 `saas-mysql`，mysql:8.0.46，绑 127.0.0.1:**3307**，`--restart unless-stopped`）、`uploadfiles/`=附件 13G（运行资产）、`old-daily/`=旧每日备份与 `mysql-data.mariadb-old/`=旧 mariadb 数据目录副本（非运行资产，确认稳定后可清）、`misc/`=零散 5 件归档。日常操作经软链 `/home/ubuntu/saas`、`/root/panda/supply` 等落到数据盘。

**启停**：全部 nohup 启动，机器重启后不自动拉起（尚无 systemd）；整机一键恢复 `/root/panda/postboot-recovery.sh`（幂等，全栈顺序见脚本头，proxy 未起会拒拉服务）；只重启 Java 层用 `/root/panda/startjars.sh`（按 PID 操作，**不要 `pkill -f saas-`**——会匹配到自己的 ssh 命令行把会话一起杀掉）。saas-common 曾因 RocketMQ proxy gRPC 竞态启动失败，proxy 稳定后手动重拉即可。

## 3. 云管理数据库：saas-mysql（3307）

**数据库双容器分家（2026-11-17 徐先生定，自 mariadb:10.11 迁 mysql:8.0.46）**：同机熊猫库独立为 `pandamc-mysql`（3306，见主工程 AGENTS.md）；云管理 = 容器 `saas-mysql`（mysql:8.0.46，仅绑 127.0.0.1:**3307**，装 chain/chain2/zsc_chain/zjjg/bustrade/zjjgbak/zsc_bustrade/zsc_zjjg 共 8 库，datadir `/mnt/datadisk0/saas/mysql-data`，业务账号 admin 用 mysql_native_password 对齐老密码——**老 JDBC/mariadb 驱动必须 native 插件**）；数据源统一在 nacos `common-dev.yaml`（3 个 `jdbc:mariadb://127.0.0.1:3307/` URL，改前版本备份 `/mnt/datadisk0/backup/20261117/common-dev.before-3307.yaml`）；全量 SQL 备份在 `/mnt/datadisk0/backup/20261117/db-mysql8/`（`*.mysql8.sql` 为方言清洗版）。表数核对：chain=146、chain2=131、zsc_chain=137、zjjg=117，bustrade/zjjgbak/zsc_zjjg 本为空库，zsc_bustrade 为 0 表属正常（见第 4 节）。

MariaDB datadir 不能被 MySQL 8 直接挂载，再迁移必须走 SQL 导出导入；转储须做方言清洗：剥 TEXT/BLOB DEFAULT（错误 1101）、剥 >16 列复合索引（错误 1070）、varchar 挂 current_timestamp() 默认改 NULL（错误 1067）、生成列表导入后 ALTER 还原（错误 3105，qualification_versions 为例）。清洗脚本 `/root/panda/fixdump.py`。旧 mariadb:10.11 容器已不在（2026-11-17 实查 `docker ps -a` 仅剩双 mysql:8.0.46 容器）；如需回滚，用数据目录副本 `/mnt/datadisk0/saas/mysql-data.mariadb-old` + `mariadb:10.11` 镜像（仍在）重建容器。

## 4. 老库完整性核对与遗留定案（2026-11-17，徐先生定）

已对 .16.10 实测核对——业务库就是上述 8 个无多库，chain=146/chain2=131/zjjg=117/zsc_chain=137 表数与迁移副本逐一相符，bustrade/zjjgbak/zsc_zjjg 本为空库；唯一异常 zsc_bustrade 老库上 50 张表全部是引擎级孤儿表（错误 1932 Table doesn't exist in engine，迁移前即坏、老系统 zsc 服务同样读不出、SQL 层零可救、--force 强导 0 表 0 行），**徐先生定不做文件级抢救，就此了结**——后续会话勿再重复排查；核验用 admin 账号（SHOW DATABASES 仅显授权库），root 私有库未核；.14→.16.10 无免密 SSH。

## 5. 152 上的云管理资产（152.136.54.164）

- 旧 Java SaaS 两包 `saas-auth-zsc.jar`、`saas-warehouse-zsc.jar`（nohup 裸跑，无 systemd 管理，动前先确认进程），与 .14 微服务同家族。
- **2026-11-16 老煤炭云前端套件迁入 152**：.6 的 7 套前端产物在 `~/meiyu/web/supply/`（scmdist/xmmcdist/zxmmcdist/zscdist/zjjgdist/distrisk + validateFiles + applet），meiyu-nginx 的 cs.conf 尾部追加 81-87 端口 server 块（原配置备份 `cs.conf.bak-20261116`），`/saas/`→后端机:8765、`/zsc/`→:5765、`/risk/`→:7765（risk 服务未迁移暂 502）；目前为 IP+HTTP 明文形态，小程序端需等 HTTPS 域名。

## 6. 运维教训（2026-11-17）

① 停 Redis 必须 `redis-cli shutdown`——其进程标题会改写为 `redis-server 127.0.0.1:6379`，pkill 匹配不到配置路径（本次孤儿进程躲过三轮搬迁、一直在旧目录写 AOF，已优雅切换到数据盘）。
② rsync `--exclude='logs'` 匹配**任意层级**同名目录，find `-not -path './logs/*'` 只排顶层，混用会漏文件——核对此类搬迁必须用 `comm` 名单比对且「原目录独有文件=0」为准。
