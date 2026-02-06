# 第8章：测试与CI/CD

## 背景知识
- Go 单元测试与覆盖率统计
- GitHub Actions 工作流配置

## 代码走读
- .github/workflows/ci.yml：安装Go、运行make、覆盖率检查
- Makefile：test/cover/check-coverage 目标
- 设计测试：避开特权操作，多测纯函数与解析流程

## 练习题
- 编写更多解析与路径函数，提升覆盖率到 80%+
- 引入集成测试（root 环境下）验证 chroot/mount

## 扩展思考
- 如何在CI中使用特权容器运行集成测试
