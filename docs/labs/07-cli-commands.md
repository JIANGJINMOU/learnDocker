# 第7章：CLI 命令实现（run/build/ps/pull）

## 背景知识
- CLI 参数解析与子命令设计

## 代码走读
- main.go：命令分发、参数校验
- runtime.go：具体实现与辅助函数

## 练习题
- 扩展 run 支持 -e 环境变量注入、-v 目录挂载（提示：mount bind）
- 为 ps 增加 cgroup 统计列

## 扩展思考
- 如何将命令行改造成 REST API + CLI 客户端
