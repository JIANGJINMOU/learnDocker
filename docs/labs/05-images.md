# 第5章：镜像格式与拉取

## 背景知识
- Docker save 格式：manifest.json + 多个 layer.tar
- OCI 镜像与 Registry v2（课后拓展）

## 代码走读
- ImportDockerSaveTar：解析 manifest，解出层并写入 metadata.json

## 练习题
- 使用 docker save busybox:latest 生成 tar，并导入本地
- 设计自己的 rootfs（FROM scratch + ADD）并运行

## 扩展思考
- 如何实现直接从 registry 拉取（token、manifest、多架构支持）
