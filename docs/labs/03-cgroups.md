# 第3章：Cgroup v2 资源限制

## 背景知识
- cgroup v2 文件接口：cpu.max、memory.max、pids.max、cgroup.procs
- 需要 root 权限以写入控制文件

## 代码走读
- internal/cgroups/v2.go：创建组并写入限额，加入进程

## 练习题
- 设置 pids.max=2，尝试在容器内并发启动多个进程，观察失败现象
- 调整 memory.max，运行占内存程序，观察 OOM 行为

## 扩展思考
- 如何将 cgroup 统计（memory.current 等）加入 ps 输出
