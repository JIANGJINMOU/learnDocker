# 第2章：命名空间实现与容器生命周期

## 背景知识
- SysProcAttr.Cloneflags 创建新命名空间
- init 进程职责：挂载 /proc、chroot、执行用户命令

## 关键数据结构
- ContainerState：记录 PID、镜像名、命令、挂载点

## 代码走读
- runContainer：构造 overlay 根文件系统、启动带 Cloneflags 的子进程
- childInit：chroot 到 rootfs，挂载 proc，执行命令

## 练习题
- 修改 hostname 并验证 UTS 隔离
- 在容器内运行 ps，验证 PID 隔离

## 扩展思考
- 容器停止后的清理流程（卸载 overlay、删除 cgroup）
