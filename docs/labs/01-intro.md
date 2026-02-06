# 第1章：容器核心概念与项目总览

## 背景知识
- 容器与虚拟机区别：进程隔离 vs 硬件虚拟化
- Linux Namespace：PID/UTS/NET/MNT 等隔离机制
- Cgroup v2：资源配额与限制
- UnionFS（OverlayFS）：写时复制与镜像分层

## 关键数据结构
- 镜像元数据：images/<name>/metadata.json
- 容器状态：containers/<id>.json
- Overlay 挂载参数：lowerdir/upperdir/workdir

## 代码走读
- cmd/cede/main.go：CLI入口与子命令分发
- cmd/cede/runtime.go：run/init 流程与命名空间
- internal/images/import.go：docker save tar 导入
- internal/overlay/overlay.go：overlayfs挂载
- internal/cgroups/v2.go：cgroup v2基础限额
- internal/state/state.go：容器状态持久化

## 练习题
- 对比 PID/UTS/NET 的隔离效果，写出验证命令
- 画出 run 子命令的执行时序图

## 扩展思考
- 为什么选择在 exec.Command 上设置 Cloneflags 而不是 unshare 命令
- 如何将 NET 命名空间接入 bridge，实现容器网络互联
