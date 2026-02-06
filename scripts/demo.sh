#!/usr/bin/env bash
set -euo pipefail

echo "== ContainerEdu Demo on Ubuntu 22.04 =="
echo "1) Build cede"
make build
echo "2) Prepare busybox image via docker save"
docker pull busybox:latest
docker save -o busybox.tar busybox:latest
echo "3) Import image into ContainerEdu"
sudo bin/cede pull --tar busybox.tar --name busybox
echo "4) Run container with isolated PID/UTS/NET namespaces"
sudo bin/cede run --image busybox --cmd /bin/sh -c "hostname && sleep 1 && echo PID=\$\$ && ip link"
