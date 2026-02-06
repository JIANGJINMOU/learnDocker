# 第4章：OverlayFS 与镜像分层

## 背景知识
- 写时复制：lowerdir 只读层 + upperdir 可写层 + workdir
- 叠加顺序：上层优先覆盖下层

## 代码走读
- internal/overlay/overlay.go：Prepare/Unmount 实现
- internal/images/import.go：从 docker save tar 解出多层

## 练习题
- 在容器中修改文件，观察 upperdir 中的变化
- 删除 lowerdir 文件，验证白out 表现（高级）

## 扩展思考
- 如何支持多镜像的多 lower 叠加（镜像拼接）
