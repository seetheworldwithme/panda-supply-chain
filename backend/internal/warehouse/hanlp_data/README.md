# HanLP 模型来源

`model.bin.gz` 从 HanLP `com.hankcs:hanlp-portable:portable-1.8.4` 的公开默认词典、双词频率、人物角色模型及字符表离线导出。该依赖 POM 声明 Apache License 2.0；模型随该库的许可声明分发。运行 Go 服务不需要 Java 或该 JAR。

导出器在生成模型前添加 `主焦` 和 `1/3焦` 两个煤种词条，以匹配原服务启动时的 `CustomDictionary.add`。对照词条由同一 JAR 的 `HanLP.segment` 生成，保存在 `golden.tsv`。

复现需要 JDK、同版本 JAR 和仓库内的 `tools/validation/ExportHanLP.java`：

```sh
javac -cp hanlp-portable-1.8.4.jar -d /tmp/hanlp-export tools/validation/ExportHanLP.java
java -cp /tmp/hanlp-export:hanlp-portable-1.8.4.jar ExportHanLP model.bin.gz golden.tsv
```

导出器反射读取 HanLP 模型内部结构，仅用于离线构建/验证；线上 Go 实现直接读取嵌入的 `model.bin.gz`。修改导出样例后，应重新生成黄金文件并运行 `go test ./internal/warehouse -run TestNativeHanLPMatchesOriginalPortable184`。
