# JRebel & JetBrains License Server (Golang)

一个用 Golang 编写的 JRebel 和 JetBrains 产品许可证服务器。

## ⚠️ 免责声明

本项目仅供学习和研究目的。请勿用于任何商业用途或非法目的。

## 🚀 快速开始

### 使用 Docker（推荐）

```bash
# 从 Docker Hub 拉取
docker pull ruking001/jrebel-license-server:latest

# 运行容器
docker run -d \
  --name license-server \
  -p 8081:8081 \
  ruking001/jrebel-license-server:latest

# 查看日志
docker logs -f license-server
```

### 使用 Docker Compose
```bash
# 下载 docker-compose.yml
curl -O https://raw.githubusercontent.com/Ruk1ng001/JrebelBrainsLicenseServerforGo/main/docker-compose.yml

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

## 📦 支持的架构

- linux/amd64 (x86_64)
- linux/arm64 (aarch64)
Docker 会自动选择适合您系统的架构。

这是一个用Golang重写的JRebel和JetBrains产品许可证服务器。

## ⚠️ 免责声明

本项目仅供学习和研究目的。请勿用于任何商业用途或非法目的。使用本软件可能违反软件许可协议，请自行承担相关风险。

## 功能特性

- 支持JRebel 7.1及更早版本
- 支持JRebel 2018.1及更高版本
- 支持JetBrains产品激活
- 支持离线激活
- 优雅关闭
- 请求日志记录
- 可配置的私钥加载

## 安装

```bash
# 克隆仓库
git clone https://github.com/Ruk1ng001/JrebelBrainsLicenseServerforGo
cd JrebelBrainsLicenseServerforGo

# 下载依赖
go mod download

# 编译
go build -o license-server cmd/server/main.go
```

## 使用方法

### 基本使用

```bash
# 使用默认端口8081
./license-server

# 指定端口
./license-server -p 9000

# 指定私钥文件
./license-server -key /path/to/private.key
```

### 环境变量配置
```bash
export PORT=8081
export SERVER_GUID=your-guid
export PRIVATE_KEY_PATH=/path/to/key
./license-server
```

## API端点

### JRebel端点

- GET /jrebel/validate-connection - 验证连接
- POST /jrebel/leases - 获取租约
- DELETE /jrebel/leases/1 - 释放租约

### JetBrains端点

- GET /rpc/ping.action - Ping测试
- GET /rpc/obtainTicket.action - 获取ticket
- GET /rpc/releaseTicket.action - 释放ticket

## 许可证
本项目仅供学习使用，不提供任何担保。