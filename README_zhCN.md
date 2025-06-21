# CAN-Bridge 项目文档

[English](README.md)

![# CAN Bridge](.github/assets/banner.jpg)

## 🔧项目概述

**CAN-Bridge** 是一个基于 Golang 开发的硬件适配网桥服务，专门用于支持使用 CAN（Controller Area Network）协议的设备进行高效稳定的数据通信。

本项目旨在提供一个简单易用的 HTTP API 接口，允许用户通过网络请求实现 CAN 消息的发送、接收、接口设置管理及设备状态监控，同时支持动态配置多个 CAN 接口。

## ⚙️功能特性

### ✨新增特性

* **接口设置管理器 (InterfaceSetupManager)**：

  * 自动执行 `ip link set can0 up type can bitrate 1000000`
  * 支持可配置的比特率、采样点、重启超时
  * 提供重试机制和错误处理
  * 支持接口状态查询和验证

* **增强的配置系统**：

  * 支持命令行参数和环境变量配置
  * 提供配置验证以确保比特率等参数合法
  * 包含详细的使用说明

* **完整的接口管理 API**：

  * 配置管理 API（`GET /api/setup/config`, `PUT /api/setup/config`）
  * 接口操作 API（设置、关闭、重置、状态查询）
  * 批量操作 API（一次性设置或关闭所有接口）

### 📦主要功能

* **动态接口配置**：支持命令行或环境变量动态配置多个 CAN 接口（例如 `can0`, `can1`）。
* **HTTP API 服务**：提供易于调用的 RESTful API 接口。
* **消息发送与验证**：支持发送标准 CAN 消息，并进行数据长度及接口有效性验证。
* **健康检查与自动恢复**：内置接口健康监测，自动重启失效接口。
* **实时监控与统计**：提供接口状态、消息统计、发送成功率、错误记录和平均延迟监控。

### 🛠️程序特性

* ✅ 自动初始化：程序启动时自动配置 CAN 接口
* ✅ 重试机制：设置失败时自动重试
* ✅ 状态监控：实时监控接口状态和错误统计
* ✅ 优雅关闭：程序退出时自动关闭接口
* ✅ 依赖注入：便于测试和扩展
* ✅ 错误处理：完善的错误处理和日志记录

## 📦安装与使用

支持裸金属安装、Docker 容器环境使用、从源码构建。

具体安装方式参考[安装文档](docs/install_zhCN.md)。

### 💡使用示例

**基本使用**

```bash
./can-bridge -can-ports can0,can1
```

**设置端口**

```bash
./can-bridge -port 5260
```

**禁用自动设置（通过 API 手动管理）**

```bash
./can-bridge -auto-setup=false
```

**自定义比特率**

```bash
./can-bridge -can-ports can0 -bitrate 500000
```

**采样点**

```bash
./can-bridge -sample-point 0.875
```

**重启超时**

```bash
./can-bridge -restart-ms 100
```

**重试次数**

```bash
./can-bridge -setup-retry 3
```

**启用服务发现**

```bash
./can-bridge -enable-finder=true
```

**服务发现间隔**

```bash
./can-bridge -finder-interval 5
```

**启用健康检查**

```bash
./can-bridge -enable-healthcheck=true
```

**通过 API 设置接口**

```bash
curl -X POST localhost:5260/api/setup/interfaces/can0 \
  -H "Content-Type: application/json" \
  -d '{"bitrate": 500000, "withRetry": true}'
```

## 🌐 API 文档

### 📍 基础路径

`http://localhost:5260/api`

### ⭐ 状态与监控

用于获取系统、接口的状态、健康信息和性能指标。

- `GET /api/status`: 获取完整的系统状态，包括正常运行时间、看门狗状态和所有接口的详细信息。
- `GET /api/interfaces`: 获取已配置和活动的接口列表。
- `GET /api/interfaces/:name/status`: 获取指定接口的详细状态。
- `GET /api/health`: 获取系统健康状况摘要。
- `GET /api/metrics`: 获取用于外部监控系统（如 Prometheus）的详细指标。

### ✉️ 消息发送

- `POST /api/can`: 发送一条 CAN 消息。请求体需要包含 CAN 消息的详细信息（如 ID, Data 等）。

### 🔧 接口设置管理 

用于动态配置、启动、停止和管理 CAN 接口。

**配置管理**：

- `GET /api/setup/config`: 获取当前的接口设置配置（如默认比特率、采样点等）。
- `PUT /api/setup/config`: 更新接口设置的全局配置。

**单个接口操作**：

- `GET /api/setup/available`: 获取操作系统上所有可用的 CAN 接口列表。
- `POST /api/setup/interfaces/{name}`: 根据配置设置并启动指定的 CAN 接口。
- `DELETE /api/setup/interfaces/{name}`: 关闭并拆除指定的 CAN 接口。
- `POST /api/setup/interfaces/{name}/reset`: 重置（先关闭再启动）指定的 CAN 接口。
- `GET /api/setup/interfaces/{name}/state`: 获取指定接口的当前状态（是否已设置、配置详情等）。

**批量接口操作**：

- `POST /api/setup/interfaces/setup-all`: 批量设置所有已配置的或请求中指定的接口。
- `POST /api/setup/interfaces/teardown-all`: 批量关闭并拆除所有已配置的接口。

### 📡 消息监听与获取

用于从 CAN 总线上实时捕获、查看和管理消息。

**监听控制**：

- `POST /api/messages/:interface/listen/start`: 在指定接口上开始监听 CAN 消息。
- `POST /api/messages/:interface/listen/stop`: 在指定接口上停止监听 CAN 消息。
- `GET /api/messages/:interface/listen/status`: 获取指定接口的当前监听状态。
- `GET /api/messages/listen/status`: 获取所有接口的监听状态汇总。

**消息获取**：

- `GET /api/messages/:interface`: 获取指定接口已缓存的所有消息。支持通过 `id` 参数进行过滤。
- `GET /api/messages/:interface/recent`: 获取指定接口最近收到的 N 条消息（可通过 `count` 参数指定数量）。
- `GET /api/messages`: 以接口为单位，获取所有接口缓存的所有消息。

**消息管理与统计**：

- `GET /api/messages/:interface/statistics`: 获取指定接口的消息统计信息（如接收总数、错误数等）。
- `DELETE /api/messages/:interface`: 清除指定接口的消息缓存。
- `GET /api/messages/statistics`: 获取所有接口的全局消息统计信息。
- `DELETE /api/messages`: 清除所有接口的消息缓存。

## 🚀性能优化与稳定性

* 支持消息发送重试机制，确保数据传输可靠性。
* 使用互斥锁（Mutex）确保多线程安全性。
* 实时监测接口健康状态并进行自动恢复。

## 📝日志与调试

日志采用标准输出，格式友好，包含清晰的错误提示和运行状态信息。

## 📦部署建议

建议使用 systemd 或 Docker 容器化进行部署，确保服务长期稳定运行。

## 🤝贡献指南

欢迎提交 Issue 和 Pull Request 来帮助项目完善和优化。

## 📄许可证

本项目使用 Apache-2.0 license 许可证。
