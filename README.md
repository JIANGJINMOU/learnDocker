# ContainerEdu (cede)

简化版容器引擎教学项目。通过 Go 语言实现容器核心机制与 CLI，覆盖命名空间、cgroup v2、OverlayFS 镜像分层、镜像导入与构建、可插拔网络插件等主题，配套完整的实验手册与课堂讲义。

## 特性
- 子命令：run / build / ps / pull / net
- 隔离：UTS / PID / NET / MNT 命名空间
- 资源：cgroup v2（cpu.max / memory.max / pids.max）
- 存储：OverlayFS（lower/upper/work）
- 镜像：导入 docker save tar、FROM scratch + ADD 构建
- 网络：插件架构，示例 bridge0（veth + bridge + NAT）

## 目录结构
- cmd/cede：CLI 与运行时（Linux 下 run/init 生效）
- internal/images：镜像导入
- internal/overlay：OverlayFS 准备与卸载
- internal/cgroups：cgroup v2 限额应用
- internal/state：容器状态持久化与 ps
- internal/plugins：网络与存储插件注册器与示例
- internal/netpool：IP 池持久化分配与释放
- docs/：实验手册、讲义、Quiz、评估问卷
- scripts/：演示与覆盖率脚本

## 快速开始（Ubuntu 22.04）

```bash
sudo apt update
sudo apt install -y golang docker.io iproute2 iptables util-linux

make build

docker pull busybox:latest
docker save -o busybox.tar busybox:latest

sudo bin/cede pull --tar busybox.tar --name busybox
sudo bin/cede run --image busybox --net bridge0 --cmd /bin/sh -c "hostname && ip addr && ip route"
sudo bin/cede ps

sudo bin/cede net ls
sudo bin/cede net config --cidr 10.0.0.0/24 --gateway 10.0.0.1
sudo bin/cede net release --id <容器ID>
```

演示脚本（含 pull 与 run）也可使用：

```bash
bash scripts/demo.sh
```

## 构建与测试

```bash
go fmt ./...
go vet ./...
go test ./... -covermode=atomic
```

在 CI（Ubuntu）环境下可执行：

```bash
bash scripts/coverage.sh
```

## 运行平台说明
- run/init 等容器运行逻辑仅在 Linux 有效（命名空间与挂载调用）
- Windows/macOS 可运行解析、导入与测试；网络插件在非 Linux 下使用 stub

## 许可协议
MIT License（见 LICENSE）

## 教学资料
- 实验手册：docs/labs
- 课堂讲义：docs/slides
- 小测题：docs/quiz
- 评估问卷：docs/evaluation_survey

