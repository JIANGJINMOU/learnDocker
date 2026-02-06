# 第6章：网络与存储驱动（可选拓展）

## 背景知识
- NET 命名空间隔离；veth、bridge、iptables/NAT 基本原理
- 存储驱动：OverlayFS、Device Mapper、Btrfs 简述

## 代码走读
- 当前示例新增可选网络插件参数 --net，默认提供 bridge0 示例插件：
  - 创建/启用 cede0 bridge（10.0.0.1/24）
  - 创建 veth pair，将容器端放入 NET 命名空间并命名为 eth0
  - 分配 10.0.0.X/24 地址与默认路由指向 10.0.0.1
  - 开启 NAT（iptables MASQUERADE）

## 练习题
- 使用 ip link 在容器内查看网络设备
- 尝试使用 iproute2 创建 veth pair 并接入 cede0
- 使用 --net bridge0 运行两个容器，验证互通（ping）与出外网（curl）
 - 观察 netpool.json 中的分配结果，尝试释放并重新分配

## 扩展思考
- 设计可插拔网络/存储插件接口（interface + 注册器模式）
- 如何进行 IP 池管理与冲突检测；如何处理容器重启后的地址保留
 - NAT 规则的幂等性与重复追加的处理策略
